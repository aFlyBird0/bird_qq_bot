package kaoyanScore

import (
	"bytes"
	"fmt"
	"time"

	"github.com/Mrs4s/MiraiGo/client"
	"github.com/Mrs4s/MiraiGo/message"
	"github.com/parnurzeal/gorequest"
	"github.com/sirupsen/logrus"

	"bird_qq_bot/utils"
)

func (m *kaoyanScore) AnalyseByGroupTrigger(c *client.QQClient, msg *message.GroupMessage) {
	if time.Now().Sub(m.lastUpdateTime) < 15*time.Second {
		tooOftenHint := "查询太频繁啦！"
		groupMsg := &message.SendingMessage{}
		groupMsg.Append(message.NewText(tooOftenHint))
		c.SendGroupMessage(msg.GroupCode, groupMsg)
		return
	}
	// 分析和存储成绩
	m.AnalyseAndSave(c)
	// 发送成绩的webserver消息
	m.sendWebserverMsgToGroup(c, msg.GroupCode)
	// 发送成绩的图片消息
	if scoreAnalyse, ok := msgFinalMap.Load(msg.GroupCode); ok {
		m.sendGroupImgMsgFromStr(c, msg.GroupCode, scoreAnalyse.(string))
	}
}

func (m *kaoyanScore) sendWebserverMsgToGroup(c *client.QQClient, groupCode int64) {
	if m.webserver.remoteURL == "" && m.webserver.localPort == "" {
		return
	}
	url := fmt.Sprintf("%s?group=%v", m.webserver.displayURL, groupCode)
	hint := "分数如下，每10分钟自动更新，每次发送关键词立即更新: "
	groupMsg := &message.SendingMessage{}
	groupMsg.Append(message.NewText(hint + url))
	c.SendGroupMessage(groupCode, groupMsg)
}

func (m *kaoyanScore) sendGroupImgMsgFromStr(c *client.QQClient, groupCode int64, msg string) {
	var buf bytes.Buffer
	err := utils.String2PicWriter(msg, m.fontPath, &buf)

	reader := bytes.NewReader(buf.Bytes())
	source := message.Source{
		SourceType: message.SourceGroup,
		PrimaryID:  groupCode,
	}
	imgMsg, err := c.UploadImage(source, reader)
	if err != nil {
		logrus.WithError(err).Error("upload img failed")
		return
	}
	c.SendGroupMessage(groupCode, message.NewSendingMessage().Append(imgMsg))
}

func (m *kaoyanScore) AnalyseAndSave(c *client.QQClient) {
	m.lastUpdateTime = time.Now()
	for group, msg := range m.generateScoreAnalyse(c) {
		updateTimeStr := m.lastUpdateTime.Format("2006-01-02 15:04:05")
		scoreAnalyseMsg := "最后更新于:" + updateTimeStr + "\n\n" + msg
		// 本地服务器和远程服务器都存一份
		m.saveGroupScoreToLocalWebServer(group, scoreAnalyseMsg)
		m.saveGroupScoreToRemoteWebServer(group, scoreAnalyseMsg)
	}
}

func (m *kaoyanScore) saveGroupScoreToLocalWebServer(groupCode int64, msg string) {
	msgFinalMap.Store(groupCode, msg)
}

func (m *kaoyanScore) saveGroupScoreToRemoteWebServer(groupCode int64, msg string) {
	url := m.webserver.remoteURL
	if url != "" {
		if _, _, errs := gorequest.New().Post(url).Send(map[string]any{
			"group": groupCode,
			"msg":   msg,
		}).End(); errs != nil {
			logger.Errorf("保存考研分数到远程服务器<%s>失败: %+v\n", url, errs)
		}
	}
}

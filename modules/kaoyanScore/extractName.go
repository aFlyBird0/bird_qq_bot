package kaoyanScore

import (
	"bird_qq_bot/utils"
	"fmt"
	"github.com/Mrs4s/MiraiGo/client"
	"github.com/Mrs4s/MiraiGo/message"
	"strings"
	"time"
)

func (m *kaoyanScore) updateAdminList(c *client.QQClient) {
	m.adminList = make([]int64, 0)
	if err := c.ReloadGroupList(); err != nil {
		fmt.Printf("ReloadGroupList error: %v\n", err)
	}
	for _, group := range c.GroupList {
		members, err := c.GetGroupMembers(group)
		if err != nil {
			logger.Error("获取群成员列表失败", err)
		} else {
			for _, member := range members {
				if member.Permission == 1 || member.Permission == 2 {
					m.adminList = append(m.adminList, member.Uin)
				}
			}
		}
	}
	logger.Info("更新管理员列表成功", m.adminList)
}

func (m *kaoyanScore) getGroupCardNames(groups []*client.GroupInfo) (cardsMap map[int64][]string) {
	// 键是群号，值是群名片列表
	cardsMap = make(map[int64][]string)
	for _, group := range groups {
		// 只筛查需要统计的群
		if !utils.InInt64(group.Code, m.AllowGroupList) {
			continue
		}
		cardsOneGroup := make([]string, 0)
		for _, member := range group.Members {
			if member.CardName != "" {
				cardsOneGroup = append(cardsOneGroup, member.CardName)
			}
		}
		cardsMap[group.Code] = cardsOneGroup
	}
	return
}

type UinCardName struct {
	Uin      int64
	CardName string
}

func (m *kaoyanScore) FindInvalidCardName(c *client.QQClient) {
	if err := c.ReloadGroupList(); err != nil {
		logger.Errorf("ReloadGroupList error: %v", err)
	}
	cardsMap := make(map[int64][]UinCardName)
	for _, group := range c.GroupList {
		// 只筛查需要统计的群
		if !utils.InInt64(group.Code, m.AllowGroupList) {
			continue
		}
		cardsOneGroup := make([]UinCardName, 0)
		for _, member := range group.Members {
			if member.CardName != "" {
				cardsOneGroup = append(cardsOneGroup, UinCardName{Uin: member.Uin, CardName: member.CardName})
			}
		}
		cardsMap[group.Code] = cardsOneGroup
	}
	zhuanHint := "请按要求更改群名片，改为「计专」或「软专」"
	xueHint := "请按要求更改群名片，改为「计学」或「软学」"
	for groupCode, cards := range cardsMap {
		for _, card := range cards {
			if strings.HasPrefix(card.CardName, "专") {
				fmt.Printf("专：%d %s\n", groupCode, card)
				msgSend := message.SendingMessage{}
				atDisplay := "@" + card.CardName
				msgSend.Append(message.NewAt(card.Uin, atDisplay))
				msgSend.Append(message.NewText(zhuanHint))
				c.SendGroupMessage(groupCode, &msgSend)
			}
			if strings.HasPrefix(card.CardName, "学") {
				fmt.Printf("学：%d %s\n", groupCode, card)
				msgSend := message.SendingMessage{}
				atDisplay := "@" + card.CardName
				msgSend.Append(message.NewAt(card.Uin, atDisplay))
				msgSend.Append(message.NewText(xueHint))
				c.SendGroupMessage(groupCode, &msgSend)
			}
			time.Sleep(time.Second * 1)
		}
	}
}

// FindCSAcademicStudent 找到群内的计算机学硕学生
func (m *kaoyanScore) FindCSAcademicStudent(c *client.QQClient) []UinCardName {
	if err := c.ReloadGroupList(); err != nil {
		logger.Errorf("ReloadGroupList error: %v", err)
	}
	cardsMap := make(map[int64][]UinCardName)
	for _, group := range c.GroupList {
		// 只筛查需要统计的群
		if !utils.InInt64(group.Code, m.AllowGroupList) {
			continue
		}
		cardsOneGroup := make([]UinCardName, 0)
		for _, member := range group.Members {
			if member.CardName != "" {
				cardsOneGroup = append(cardsOneGroup, UinCardName{Uin: member.Uin, CardName: member.CardName})
			}
		}
		cardsMap[group.Code] = cardsOneGroup
	}
	var csAcademicStudent []UinCardName
	filter := CSAcademicFilter{}
	for _, cards := range cardsMap {
		for _, card := range cards {
			if filter.Filter(card.CardName) {
				csAcademicStudent = append(csAcademicStudent, UinCardName{Uin: card.Uin, CardName: card.CardName})
			}
		}
	}
	return csAcademicStudent
}

// GenMailsFromUins 通过群名片生成群发邮箱收件人
func GenMailsFromUins(uins []UinCardName) (mails string) {
	mailSlice := make([]string, 0, len(uins))
	for _, uin := range uins {
		mailSlice = append(mailSlice, fmt.Sprintf("%v@qq.com", uin.Uin))
	}
	mails = strings.Join(mailSlice, ";")

	return
}

package kaoyanScore

import (
	"fmt"
	"github.com/Mrs4s/MiraiGo/client"
	"github.com/Mrs4s/MiraiGo/message"
)

func (m *kaoyanScore) Calculate(c *client.QQClient) {
	//m.updateGroupMembers(c)
	//fmt.Printf("群列表: %+v\n", m.Groups)
	//m.updateAdminList(c)
	if err := c.ReloadGroupList(); err != nil {
		logger.Error("ReloadGroupList error: %v\n", err)
	}
	//fmt.Printf("群列表: %+v\n", c.GroupList)
	cardsMap := m.getGroupCardNames(c.GroupList)
	//fmt.Printf("cardsMap: %+v\n", cardsMap)
	for groupCode, cards := range cardsMap {
		if len(cards) == 0 {
			continue
		}
		scoresMap := GetScoreMap(cards, CSAcademicFilter{}, CSProfessionalFilter{},
			SEAcademicFilter{}, SEProfessionalFilter{}, Russia{}, Japan{})
		GroupCodeScoreMap := make(map[string][]ScoreGroup)
		for typ, scores := range scoresMap {
			logger.Infof("%s: %+v\n", typ, scores)
			// 为每个类型分组
			GroupCodeScoreMap[typ] = GroupScores(scores)
		}
		msg := "考研分数段统计来啦！\n"
		for typ, scoreGroups := range GroupCodeScoreMap {
			counts := 0
			for _, v := range scoreGroups {
				counts += v.Len()
			}
			msg += fmt.Sprintf("%s(共%v个分数)\n", typ, counts)
			for _, scoreGroup := range scoreGroups {
				msg += fmt.Sprintf("%v: %v人\n", scoreGroup.Describe(), scoreGroup.Len())
			}
		}
		msg += "以上结果通过群名片分析而得，存在一定误差，如误识别实验室门牌号为考研分数，仅供参考。\n"
		msg += m.tailMsg
		logger.Info("拼接得到的考研分数排名:\n %v\n", msg)
		groupMsg := &message.SendingMessage{}
		groupMsg.Append(message.NewText(msg))
		c.SendGroupMessage(groupCode, groupMsg)
	}
}

func (m *kaoyanScore) CalculateByTrigger(c *client.QQClient, msg *message.GroupMessage) {
	m.Calculate(c)
}

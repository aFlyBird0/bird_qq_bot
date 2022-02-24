package kaoyanScore

import (
	"fmt"
	"github.com/Mrs4s/MiraiGo/client"
	"github.com/Mrs4s/MiraiGo/message"
)

func (m *kaoyanScore) Calculate(c *client.QQClient) {
	if err := c.ReloadGroupList(); err != nil {
		logger.Error("ReloadGroupList error: %v\n", err)
	}
	//fmt.Printf("群列表: %+v\n", c.GroupList)
	cardsMap := m.getGroupCardNames(c.GroupList)
	//fmt.Printf("cardsMap: %+v\n", cardsMap)
	filters := append([]ScoreFilter{}, CSAcademicFilter{}, CSProfessionalFilter{},
		SEAcademicFilter{}, SEProfessionalFilter{}, Russia{}, Japan{})

	// 发送分数段信息
	for groupCode, cards := range cardsMap {
		if len(cards) == 0 {
			continue
		}
		scoresMap := GetScoreMap(cards, filters...)
		GroupCodeScoreMap := make(map[ScoreFilter][]ScoreGroup)
		for filter, scores := range scoresMap {
			//logger.Infof("%s: %+v\n", filter.Name(), scores)
			// 为每个类型分组
			GroupCodeScoreMap[filter] = GroupScoresEachTen(scores)
		}
		msg := "考研分数段统计来啦！\n"
		for _, filter := range filters {
			if scoreGroups, ok := GroupCodeScoreMap[filter]; ok {
				counts := 0
				for _, v := range scoreGroups {
					counts += v.Len()
				}
				msg += fmt.Sprintf("%s(共%v个分数)\n", filter.Name(), counts)
				sum := 0 // 每段依次累加人数
				for _, scoreGroup := range scoreGroups {
					sum += scoreGroup.Len()
					msg += fmt.Sprintf("%v: %v人(累计%v)\n", scoreGroup.Describe(), scoreGroup.Len(), sum)
				}
			}
		}
		msg += "以上结果通过群名片分析而得，存在一定误差，仅供参考。\n"
		msg += m.tailMsg
		logger.Info("拼接得到的考研分数排名:\n %v\n", msg)
		groupMsg := &message.SendingMessage{}
		groupMsg.Append(message.NewText(msg))
		c.SendGroupMessage(groupCode, groupMsg)
	}

	// 发送某些过密分数段分布详情
	for groupCode, cards := range cardsMap {
		if len(cards) == 0 {
			continue
		}
		scoresMap := GetScoreMap(cards, filters...)
		GroupCodeScoreMap := make(map[ScoreFilter][]ScoreGroupCount)
		for filter, scores := range scoresMap {
			// 密集的分数段
			denseGroups := ScoreGroupCountList(GroupScoresFromCount(CountScores(scores))).FilterByCount(10)
			GroupCodeScoreMap[filter] = denseGroups
		}
		msg := "过密分数段分布来啦！\n"
		for _, filter := range filters {
			if scoreGroupCounts, ok := GroupCodeScoreMap[filter]; ok {
				msg += filter.Name() + " 过密分数段分布\n"
				for _, scoreGroupCount := range scoreGroupCounts {
					msg += fmt.Sprintf("%s(共%v人)\n", scoreGroupCount.Describe(), scoreGroupCount.Count())
					for _, scoreCount := range scoreGroupCount.Scores {
						msg += scoreCount.Describe() + "  "
					}
					msg += "\n"
				}
			}
		}
		logger.Info("拼接得到的过密分数排名:\n %v\n", msg)
		fmt.Println(msg)
		groupMsg := &message.SendingMessage{}
		groupMsg.Append(message.NewText(msg))
		c.SendGroupMessage(groupCode, groupMsg)
	}
}

func (m *kaoyanScore) CalculateByTrigger(c *client.QQClient, msg *message.GroupMessage) {
	m.Calculate(c)
}

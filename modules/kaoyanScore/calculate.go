package kaoyanScore

import (
	"fmt"
	"github.com/Mrs4s/MiraiGo/client"
	"github.com/Mrs4s/MiraiGo/message"
	"strconv"
	"time"
)

// Calculate 这里消息的拼接和发送揉在一起了，有空再拆
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
		msg := "过密分数段分布来啦！(某段达到10人及以上）\n"
		for _, filter := range filters {
			if scoreGroupCounts, ok := GroupCodeScoreMap[filter]; ok {
				if len(scoreGroupCounts) == 0 {
					continue
				}
				msg += "\n" + filter.Name() + " 过密分数段分布\n"
				for _, scoreGroupCount := range scoreGroupCounts {
					msg += fmt.Sprintf("【%s】(共%v人)\n", scoreGroupCount.Describe(), scoreGroupCount.Count())
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

type groupCode = int64

func (m *kaoyanScore) CalculatePure(c *client.QQClient) map[groupCode][2]string {
	if err := c.ReloadGroupList(); err != nil {
		logger.Error("ReloadGroupList error: %v\n", err)
	}
	//fmt.Printf("群列表: %+v\n", c.GroupList)
	cardsMap := m.getGroupCardNames(c.GroupList)
	//fmt.Printf("cardsMap: %+v\n", cardsMap)
	filters := append([]ScoreFilter{}, CSAcademicFilter{}, CSProfessionalFilter{},
		SEAcademicFilter{}, SEProfessionalFilter{}, Russia{}, Japan{})

	result := make(map[groupCode][2]string)
	for group, cards := range cardsMap {

		// 发送分数段信息
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
		msg1 := "考研分数段统计来啦！\n"
		for _, filter := range filters {
			if scoreGroups, ok := GroupCodeScoreMap[filter]; ok {
				counts := 0
				for _, v := range scoreGroups {
					counts += v.Len()
				}
				msg1 += fmt.Sprintf("%s(共%v个分数)\n", filter.Name(), counts)
				sum := 0 // 每段依次累加人数
				for _, scoreGroup := range scoreGroups {
					sum += scoreGroup.Len()
					msg1 += fmt.Sprintf("%v: %v人(累计%v)\n", scoreGroup.Describe(), scoreGroup.Len(), sum)
				}
			}
		}
		msg1 += "以上结果通过群名片分析而得，存在一定误差，仅供参考。\n"
		msg1 += m.tailMsg
		logger.Info("拼接得到的考研分数排名:\n %v\n", msg1)

		// 发送某些过密分数段分布详情
		GroupCodeScoreCountMap := make(map[ScoreFilter][]ScoreGroupCount)
		for filter, scores := range scoresMap {
			// 密集的分数段
			denseGroups := ScoreGroupCountList(GroupScoresFromCount(CountScores(scores))).FilterByCount(10)
			GroupCodeScoreCountMap[filter] = denseGroups
		}
		msg2 := "过密分数段分布来啦！(某段达到10人及以上）\n"
		for _, filter := range filters {
			if scoreGroupCounts, ok := GroupCodeScoreCountMap[filter]; ok {
				if len(scoreGroupCounts) == 0 {
					continue
				}
				msg2 += "\n" + filter.Name() + " 过密分数段分布\n"
				for _, scoreGroupCount := range scoreGroupCounts {
					msg2 += fmt.Sprintf("【%s】(共%v人)\n", scoreGroupCount.Describe(), scoreGroupCount.Count())
					for _, scoreCount := range scoreGroupCount.Scores {
						msg2 += scoreCount.Describe() + "  "
					}
					msg2 += "\n"
				}
			}
		}
		logger.Info("拼接得到的过密分数排名:\n %v\n", msg2)
		fmt.Println(msg2)
		result[group] = [2]string{msg1, msg2}
	}

	return result
}

func (m *kaoyanScore) CalculateAndSave(c *client.QQClient) {
	msgMap := m.CalculatePure(c)
	m.lastUpdateTime = time.Now()
	for group, msgs := range msgMap {
		updateTimeStr := m.lastUpdateTime.Format("2006-01-02 15:04:05")
		msgFinalMap.Store(group, "最后更新于:"+updateTimeStr+"\n\n"+msgs[0]+"\n"+msgs[1]+"\n")
	}
}

func (m *kaoyanScore) CalculateByGroupTrigger(c *client.QQClient, msg *message.GroupMessage) {
	if time.Now().Sub(m.lastUpdateTime) < 15*time.Second {
		tooOftenHint := "查询太频繁啦！"
		groupMsg := &message.SendingMessage{}
		groupMsg.Append(message.NewText(tooOftenHint))
		c.SendGroupMessage(msg.GroupCode, groupMsg)
		return
	}
	m.CalculateAndSave(c)
	url := "http://score.kaoyan.aflybird.cn/score?group=" + strconv.FormatInt(msg.GroupCode, 10)
	hint := "分数如下，每10分钟自动更新，每次发送关键词立即更新: "
	groupMsg := &message.SendingMessage{}
	groupMsg.Append(message.NewText(hint + url))
	c.SendGroupMessage(msg.GroupCode, groupMsg)
	fmt.Println(hint + url)
}

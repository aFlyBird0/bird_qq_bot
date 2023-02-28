package kaoyanScore

import (
	"fmt"

	"golang.org/x/exp/constraints"

	"github.com/Mrs4s/MiraiGo/client"
)

type groupCode = int64

// generateScoreAnalyse 生成消息段分析结果
func (m *kaoyanScore) generateScoreAnalyse(c *client.QQClient) map[groupCode]string {
	// todo: 重构，不是一下子返回整个消息，而是返回关键的数据结构，再用tpl渲染
	if err := c.ReloadGroupList(); err != nil {
		logger.Error("ReloadGroupList error: %v\n", err)
	}
	//fmt.Printf("群列表: %+v\n", c.GroupList)
	cardsMap := m.getGroupCardNames(c.GroupList)
	//fmt.Printf("cardsMap: %+v\n", cardsMap)
	filters := append([]ScoreFilter{}, CSAcademicFilter{}, CSProfessionalFilter{},
		SEAcademicFilter{}, SEProfessionalFilter{}, Russia{}, Japan{})

	result := make(map[groupCode]string)
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
		msg1 := "---考研分数段统计来啦！---\n\n"
		scoresTotal := 0
		for _, filter := range filters {
			if scoreGroups, ok := GroupCodeScoreMap[filter]; ok {
				for _, v := range scoreGroups {
					scoresTotal += v.Len()
				}
			}
		}
		msg1 += fmt.Sprintf("所有专业一共统计到%v个分数\n\n", scoresTotal)
		for _, filter := range filters {
			if scoreGroups, ok := GroupCodeScoreMap[filter]; ok {
				counts := 0
				scoreSum := 0
				scoresFlat := make([]float32, 0)
				for _, v := range scoreGroups {
					counts += v.Len()
					for _, score := range v.scores {
						scoreSum += score
						scoresFlat = append(scoresFlat, float32(score))
					}
				}
				avg := float32(scoreSum) / float32(counts)
				mid := getMidNum(scoresFlat)

				msg1 += fmt.Sprintf("%s(共%v个分数，均分%.2f, 中位数%.1f)\n", filter.Name(), counts, avg, mid)
				sum := 0 // 每段依次累加人数
				for _, scoreGroup := range scoreGroups {
					sum += scoreGroup.Len()
					msg1 += fmt.Sprintf("%v: %v人(累计%v)\n", scoreGroup.Describe(), scoreGroup.Len(), sum)
				}
			}
		}
		msg1 += "\n\n以上结果通过群名片分析而得，存在一定误差，仅供参考。\n"
		//msg1 += m.headMsgInWebserver + "\n"
		//logger.Info("拼接得到的考研分数排名:\n %v\n", msg1)

		// 发送某些过密分数段分布详情
		GroupCodeScoreCountMap := make(map[ScoreFilter][]ScoreGroupCount)
		for filter, scores := range scoresMap {
			// 密集的分数段
			denseGroups := ScoreGroupCountList(GroupScoresFromCount(CountScores(scores))).FilterByCount(10)
			GroupCodeScoreCountMap[filter] = denseGroups
		}
		msg2 := "---过密分数段分布来啦！(某段达到10人及以上）---\n"
		for _, filter := range filters {
			if scoreGroupCounts, ok := GroupCodeScoreCountMap[filter]; ok {
				if len(scoreGroupCounts) == 0 {
					continue
				}
				msg2 += "\n" + filter.Name() + " 过密分数段分布\n"
				for _, scoreGroupCount := range scoreGroupCounts {
					msg2 += fmt.Sprintf("【%s】(共%v人)\n", scoreGroupCount.Describe(), scoreGroupCount.Count())
					count := 0
					for _, scoreCount := range scoreGroupCount.Scores {
						if count > 0 && count%5 == 0 {
							msg2 += "\n"
						}
						msg2 += scoreCount.Describe() + "  "
						count += 1
					}
					msg2 += "\n"
				}
			}
		}
		//logger.Info("拼接得到的过密分数排名:\n %v\n", msg2)
		//fmt.Println(msg2)
		msg := msg1 + msg2
		logger.Infof("拼接得到的考研分数排名:\n %v\n", msg)
		result[group] = msg
	}

	return result
}

// 从给定的增序的数组中求中位数
func getMidNum[T constraints.Float | constraints.Integer](nums []T) T {
	if len(nums) == 0 {
		return 0
	}
	if len(nums)%2 == 0 {
		return (nums[len(nums)/2-1] + nums[len(nums)/2]) / 2
	}
	return nums[len(nums)/2]
}

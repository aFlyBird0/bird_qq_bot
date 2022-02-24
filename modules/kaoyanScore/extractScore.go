package kaoyanScore

import (
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

type ScoreFilter interface {
	Name() string
	Filter(nickname string) bool
}

func ExtractScore(nickname string) (score int, ok bool) {
	regThreeNum := regexp.MustCompile("\\d{3}")
	scoreStr := regThreeNum.FindString(nickname)
	if scoreStr == "" {
		return 0, false
	}
	score, err := strconv.Atoi(scoreStr)
	if err != nil {
		return 0, false
	}
	if score < 200 || score > 500 {
		return 0, false
	}
	return score, true
}

// GetScoreMap 传入昵称列表和过滤规则，返回每个规则命中的分数列表
func GetScoreMap(nicknames []string, filters ...ScoreFilter) map[ScoreFilter][]int {
	scoreMap := make(map[ScoreFilter][]int)
	for _, nickname := range nicknames {
		for _, filter := range filters {
			if filter.Filter(nickname) {
				score, ok := ExtractScore(nickname)
				if ok {
					scoreMap[filter] = append(scoreMap[filter], score)
				}
			}
		}
	}
	// 分数排序
	for _, scores := range scoreMap {
		sort.Slice(scores, func(i, j int) bool {
			return scores[i] > scores[j]
		})
	}
	return scoreMap
}

type ScoreRank struct {
	Name        string       // 名称，如计算机专硕
	ScoreGroups []ScoreGroup // 分数段
}

type ScoreGroup struct {
	Min    int   // 最小分数
	Max    int   // 最大分数
	scores []int // 分数列表
}

func (s ScoreGroup) Describe() string {
	return fmt.Sprintf("%v - %v", s.Min, s.Max)
}

func (s ScoreGroup) Len() int {
	return len(s.scores)
}

// GroupScoresEachTen 将分数分组，严格十分一段
func GroupScoresEachTen(scores []int) (scoreGroups []ScoreGroup) {
	if len(scores) == 0 {
		return
	}
	// 分数段起始
	var start int = scores[0] / 10 * 10
	var end int = start + 9
	scoresOneGroup := make([]int, 0)
	for _, score := range scores {
		// 在同一分数段内
		if score >= start && score <= end {
			scoresOneGroup = append(scoresOneGroup, score)
		} else {
			// 不在同一分数段内
			scoreGroups = append(scoreGroups, ScoreGroup{start, end, scoresOneGroup})
			scoresOneGroup = make([]int, 0)
			start = score / 10 * 10
			end = start + 9
			scoresOneGroup = append(scoresOneGroup, score)
		}
	}
	// 别忘了最后一段
	scoreGroups = append(scoreGroups, ScoreGroup{start, end, scoresOneGroup})
	return scoreGroups
}

// ScoreCount 记录每个分数的数量
type ScoreCount struct {
	Score int
	Count int
}

type ScoreCountList []ScoreCount

func CountScores(scores []int) []ScoreCount {
	// 哨兵
	scores = append(scores, -1)
	scoreCounts := make([]ScoreCount, 0)
	scoreCountNow := ScoreCount{Score: scores[0]}
	countNow := 1
	for i, score := range scores[:len(scores)-1] {
		if score == scores[i+1] {
			countNow++
		} else {
			scoreCountNow.Count = countNow
			scoreCounts = append(scoreCounts, scoreCountNow)
			scoreCountNow = ScoreCount{Score: score}
			countNow = 1
		}
	}
	return scoreCounts
}

func (s ScoreCount) Describe() string {
	return fmt.Sprintf("%v分: %v人", s.Score, s.Count)
}

func (s ScoreCountList) Filter(remain func(s ScoreCount) bool) ScoreCountList {
	var result ScoreCountList
	for _, scoreCount := range s {
		if remain(scoreCount) {
			result = append(result, scoreCount)
		}
	}
	return result
}

func (s ScoreCountList) FilterByCount(count int) ScoreCountList {
	return s.Filter(func(s ScoreCount) bool {
		return s.Count >= count
	})
}

type ScoreGroupCount struct {
	Min    int          // 最小分数
	Max    int          // 最大分数
	Scores []ScoreCount // 分数列表
}

func GroupScoresFromCount(scores []ScoreCount) []ScoreGroupCount {
	if len(scores) == 0 {
		return nil
	}
	scoreGroupCounts := make([]ScoreGroupCount, 0)
	// 分数段起始
	var start int = scores[0].Score / 10 * 10
	var end int = start + 9
	scoresOneGroup := make([]ScoreCount, 0)
	for _, scoreCount := range scores {
		// 在同一分数段内的 ScoreCount
		if scoreCount.Score >= start && scoreCount.Score <= end {
			scoresOneGroup = append(scoresOneGroup, scoreCount)
		} else {
			// 不在同一分数段内
			scoreGroupCounts = append(scoreGroupCounts, ScoreGroupCount{start, end, scoresOneGroup})
			start = scoreCount.Score / 10 * 10
			end = start + 9
			scoresOneGroup = make([]ScoreCount, 0)
		}
	}
	// 别忘了最后一段
	scoreGroupCounts = append(scoreGroupCounts, ScoreGroupCount{start, end, scoresOneGroup})
	return scoreGroupCounts
}

type ScoreGroupCountList []ScoreGroupCount

func (s *ScoreGroupCount) Describe() string {
	return fmt.Sprintf("%v - %v 分段明细", s.Min, s.Max)
}

func (s *ScoreGroupCount) Count() (count int) {
	for _, scoreCount := range s.Scores {
		count += scoreCount.Count
	}
	return
}

func (list ScoreGroupCountList) Filter(remain func(one ScoreGroupCount) bool) ScoreGroupCountList {
	var result ScoreGroupCountList
	for _, one := range list {
		if remain(one) {
			result = append(result, one)
		}
	}
	return result
}

func (list ScoreGroupCountList) FilterByCount(count int) ScoreGroupCountList {
	return list.Filter(func(one ScoreGroupCount) bool {
		return one.Count() >= count
	})
}

// CSAcademicFilter 计算机学硕
type CSAcademicFilter struct {
}

func (C CSAcademicFilter) Name() string {
	return "计算机学硕"
}

func (C CSAcademicFilter) Filter(nickname string) bool {
	reg1 := regexp.MustCompile("计")
	reg2 := regexp.MustCompile("机")
	reg3 := regexp.MustCompile("学")
	reg4 := regexp.MustCompile("计科")
	// 计科 || 计学 || 机学
	return ((reg1.MatchString(nickname) || reg2.MatchString(nickname)) && reg3.MatchString(nickname)) || reg4.MatchString(nickname) || strings.HasPrefix(nickname, "学")
}

// CSProfessionalFilter 计算机专硕
type CSProfessionalFilter struct {
}

func (C CSProfessionalFilter) Name() string {
	return "计算机专硕"
}

func (C CSProfessionalFilter) Filter(nickname string) bool {
	reg1 := regexp.MustCompile("计")
	reg2 := regexp.MustCompile("机")
	reg3 := regexp.MustCompile("专")
	return ((reg1.MatchString(nickname) || reg2.MatchString(nickname)) && reg3.MatchString(nickname)) || strings.HasPrefix(nickname, "专")
}

// SEAcademicFilter 软件工程学硕
type SEAcademicFilter struct {
}

func (S SEAcademicFilter) Name() string {
	return "软件工程学硕"
}

func (S SEAcademicFilter) Filter(nickname string) bool {
	reg1 := regexp.MustCompile("软")
	reg2 := regexp.MustCompile("学")
	return reg1.MatchString(nickname) && reg2.MatchString(nickname)
}

// SEProfessionalFilter 软件工程专硕
type SEProfessionalFilter struct {
}

func (S SEProfessionalFilter) Name() string {
	return "软件工程专硕"
}

func (S SEProfessionalFilter) Filter(nickname string) bool {
	reg1 := regexp.MustCompile("软")
	reg2 := regexp.MustCompile("专")
	return reg1.MatchString(nickname) && reg2.MatchString(nickname)
}

// Russia 中俄
type Russia struct {
}

func (r Russia) Name() string {
	return "中俄"
}

func (r Russia) Filter(nickname string) bool {
	reg := regexp.MustCompile("中俄")
	return reg.MatchString(nickname)
}

// Japan 中日
type Japan struct {
}

func (j Japan) Name() string {
	return "中日"
}

func (j Japan) Filter(nickname string) bool {
	reg := regexp.MustCompile("中日")
	return reg.MatchString(nickname)
}

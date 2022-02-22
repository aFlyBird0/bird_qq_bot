package kaoyanScore

import (
	"fmt"
	"regexp"
	"sort"
	"strconv"
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

func GetScoreMap(nicknames []string, filters ...ScoreFilter) map[string][]int {
	scoreMap := make(map[string][]int)
	for _, nickname := range nicknames {
		for _, filter := range filters {
			if filter.Filter(nickname) {
				score, ok := ExtractScore(nickname)
				if ok {
					scoreMap[filter.Name()] = append(scoreMap[filter.Name()], score)
				}
			}
		}
	}
	// 排序
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

func GroupScores(scores []int) (scoreGroups []ScoreGroup) {
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
	return ((reg1.MatchString(nickname) || reg2.MatchString(nickname)) && reg3.MatchString(nickname)) || reg4.MatchString(nickname)
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
	return (reg1.MatchString(nickname) || reg2.MatchString(nickname)) && reg3.MatchString(nickname)
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

package utils

import (
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().Unix())
}

// GetSomeRandNumber 获取 num 个随机数
// 形参：num:随机数数量，start: 起始数，end：结束数，包含开始不包含结束
// 返回：nums：随机数切片
func GetSomeRandNumber(num, start, end int) (nums []int) {
	i := 0
	for i < num {
		// 循环生成随机数，若与已生成不重复，就加入
		// 注意这里生成的数是不包含边界的
		numTemp := rand.Intn(end-start) + start
		//st := time.Now().Unix()
		//numTemp := int(st % int64(end-start)) + start
		//numTemp := int(time.Now().UnixNano()%int64(end-start) + int64(start))
		//time.Sleep(1)
		//fmt.Println(numTemp)
		if !contains(nums, numTemp) {
			nums = append(nums, numTemp)
			i++
		}
	}
	return nums
}

func contains(nums []int, target int) (exist bool) {
	for _, value := range nums {
		if target == value {
			return true
		}
	}
	return false
}

func GetOneRandNum(start, end int) (randNum int) {
	return rand.Intn(end-start) + start
}

package kaoyanScore

import (
	"bird_qq_bot/utils"
	"fmt"
	"github.com/Mrs4s/MiraiGo/client"
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

package noCopy

import (
	"github.com/Mrs4s/MiraiGo/client"
	"sync"
)

// 循环消息队列，记录最近的 msgTraceBackNum 条消息
type groupMsg struct {
	msgMap   map[int64][]string
	indexMap map[int64]int
	*sync.Mutex
}

func NewGroupMsg() *groupMsg {
	g := groupMsg{
		msgMap:   make(map[int64][]string),
		indexMap: make(map[int64]int),
	}
	g.Mutex = &sync.Mutex{}
	return &g
}

// 判断用户消息是否属于复读，并更新消息队列
func (g *groupMsg) isMsgRepeat(groupCode int64, msg string, same msgCompareFunc) bool {
	g.Lock()
	defer g.Unlock()
	if _, ok := g.msgMap[groupCode]; !ok {
		g.msgMap[groupCode] = make([]string, msgTraceBackNum)
		g.indexMap[groupCode] = 0
	}

	// 遍历循环消息队列，计算前面已经出现过的相同的消息数量
	appearedTimes := 0
	for _, v := range g.msgMap[groupCode] {
		if same(v, msg) {
			appearedTimes++
			if appearedTimes >= maxRepeat {
				break
			}
		}
	}

	// 更新消息队列记录
	g.msgMap[groupCode][g.indexMap[groupCode]] = msg
	g.indexMap[groupCode]++
	if g.indexMap[groupCode] == 10 {
		g.indexMap[groupCode] = 0
	}
	return appearedTimes >= maxRepeat
}

// 消息对比函数签名
// 这么写的好处在于以后方便修改消息相等逻辑，例如忽略表情差异
type msgCompareFunc func(string, string) bool

// 严格对比
func strictCompare(msg1, msg2 string) bool {
	return msg1 == msg2
}

func (n *noCopy) updateGroupMembers(c *client.QQClient) {
	//n.Groups = make(map[int64]*client.GroupInfo)
	groups, err := c.GetGroupList()
	if err != nil {
		logger.Error("获取群列表失败", err)
	}
	for _, group := range groups {
		members, err := c.GetGroupMembers(group)
		if err != nil {
			logger.Error("获取群成员列表失败", err)
		} else {
			group.Members = members
		}
		n.Lock()
		n.Groups[group.Code] = group
		n.Unlock()
	}
}

func (n *noCopy) getGroupFromCache(groupCode int64) *client.GroupInfo {
	n.Lock()
	logger.Info("获取群信息", n.Groups)
	g := n.Groups[groupCode]
	n.Unlock()
	return g
}

func in(array []string, src string) bool {
	for _, s := range array {
		if s == src {
			return true
		}
	}
	return false
}

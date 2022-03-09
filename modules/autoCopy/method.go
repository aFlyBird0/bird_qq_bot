package noCopy

import (
	"github.com/Mrs4s/MiraiGo/message"
	"strconv"
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
			// 本次消息也要计算在内，所以减一
			if appearedTimes >= maxRepeat-1 {
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
	return appearedTimes >= maxRepeat-1
}

func (g *groupMsg) reset(groupCode int64) {
	g.Lock()
	defer g.Unlock()
	delete(g.msgMap, groupCode)
	delete(g.indexMap, groupCode)
	g.msgMap[groupCode] = make([]string, msgTraceBackNum)
	g.indexMap[groupCode] = 0
}

// 消息对比函数签名
// 这么写的好处在于以后方便修改消息相等逻辑，例如忽略表情差异
type msgCompareFunc func(string, string) bool

// 严格对比
func strictCompare(msg1, msg2 string) bool {
	return msg1 == msg2
}

// GroupMessageWrapper 包装一层 GroupMessage， 为了重写 ToString 方法
type GroupMessageWrapper struct {
	*message.GroupMessage
}

func NewGroupMessageWrapper(msg *message.GroupMessage) *GroupMessageWrapper {
	return &GroupMessageWrapper{msg}
}

func (msg *GroupMessageWrapper) ToString() (res string) {
	for _, elem := range msg.Elements {
		switch e := elem.(type) {
		case *message.TextElement:
			res += e.Content
		case *message.FaceElement:
			res += "[" + e.Name + "]"
		case *message.MarketFaceElement:
			res += "[" + e.Name + "]"
		case *message.GroupImageElement:
			res += "[Image: " + e.ImageId + "]"
		case *message.AtElement:
			res += e.Display
		case *message.RedBagElement:
			res += "[RedBag:" + e.Title + "]"
		case *message.ReplyElement:
			res += "[Reply:" + strconv.FormatInt(int64(e.ReplySeq), 10) + "]"
		}
	}
	return
}

func in(array []string, src string) bool {
	for _, s := range array {
		if s == src {
			return true
		}
	}
	return false
}

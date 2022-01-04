package goodNight

import (
	"github.com/parnurzeal/gorequest"
	"time"
)

const tianXingWanAnApi = "http://api.tianapi.com/wanan/index"

type apiResp struct {
	Code     int    `json:"code"`
	Msg      string `json:"msg"`
	NewsList []struct {
		Content string `json:"content"`
	} `json:"newslist"`
}

func getNightMsg(apiKey string) string {
	resp := apiResp{}
	url := tianXingWanAnApi + "?key=" + apiKey
	_, _, errors := gorequest.New().Timeout(time.Second * 10).Get(url).EndStruct(&resp)
	if errors != nil {
		logger.Errorf("获取情话失败：%v", errors)
		return ""
	}
	if resp.Code == 200 && len(resp.NewsList) > 0 {
		return resp.NewsList[0].Content
	} else {
		logger.Errorf("获取晚安消息失败：%v: %v", resp.Code, resp.Msg)
		return ""
	}
}

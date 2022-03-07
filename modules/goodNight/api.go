package goodNight

import (
	"fmt"
	"github.com/parnurzeal/gorequest"
	"time"
)

const (
	tianXingWanAnUrl   = "http://api.tianapi.com/wanan/index"
	tianXingTianGouUrl = "http://api.tianapi.com/tiangou/index"
)

type Client struct {
	apiKey string
}

type API struct {
	URL  string
	Name string
}

var (
	WanAnAPI = API{
		URL:  tianXingWanAnUrl,
		Name: "晚安",
	}
	TianGouAPI = API{
		URL:  tianXingTianGouUrl,
		Name: "舔狗日记",
	}
)

type apiResp struct {
	Code     int    `json:"code"`
	Msg      string `json:"msg"`
	NewsList []struct {
		Content string `json:"content"`
	} `json:"newslist"`
}

func NewClient(apiKey string) *Client {
	return &Client{
		apiKey: apiKey,
	}
}

func (c *Client) getFirstMsg(api API) (string, error) {
	resp := apiResp{}
	url := api.URL + "?key=" + c.apiKey
	_, _, errArr := gorequest.New().Timeout(time.Second * 10).Get(url).EndStruct(&resp)
	err := fmt.Errorf("%v接口请求失败：%v", api.Name, errArr)
	if len(errArr) > 0 {
		return "", err
	}
	if resp.Code == 200 && len(resp.NewsList) > 0 {
		return resp.NewsList[0].Content, nil
	} else {
		return "", err
	}
}

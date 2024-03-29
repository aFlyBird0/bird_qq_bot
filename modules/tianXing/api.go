package tianXing

import (
	"fmt"
	"time"

	"github.com/parnurzeal/gorequest"
)

const (
	baseUrl   = "http://api.tianapi.com/"
	urlSuffix = "/index"
)

type Client struct {
	apiKey string
}

type API struct {
	URL  string
	Name string
}

type apiResp struct {
	Code     int    `json:"code"`
	Msg      string `json:"msg"`
	NewsList []struct {
		Content string `json:"content"`
	} `json:"newslist"`
}

var (
	wanAnAPI = API{
		URL:  buildUrl("wanan"),
		Name: "晚安",
	}
	tianGouAPI = API{
		URL:  buildUrl("tiangou"),
		Name: "舔狗日记",
	}
	zaoAnAPI = API{
		URL:  buildUrl("zaoan"),
		Name: "早安",
	}
	healthTipAPI = API{
		URL:  buildUrl("healthtip"),
		Name: "健康小贴士",
	}
	sayloveAPI = API{
		URL:  buildUrl("saylove"),
		Name: "土味情话",
	}
)

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

func buildUrl(path string) string {
	return baseUrl + path + urlSuffix
}

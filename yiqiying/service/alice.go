package service

import (
	"encoding/json"
	"fmt"
	"github.com/gookit/goutil/arrutil"
	"github.com/yangzhenrui/finance/credential"
	"github.com/yangzhenrui/finance/util"
	"github.com/yangzhenrui/finance/yiqiying/context"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

const (
	// GetCloseInfoUrl 结账信息接口
	GetCloseInfoUrl = "https://openapi.17win.com/gateway/openyqdz/alice/closeInfo/getCloseInfo"
)

type Alice struct {
	*context.Context
}

func NewAlice(ctx *context.Context) *Alice {
	return &Alice{ctx}
}

func (c *Alice) setHeader(signature string, httpRequest *http.Request) {
	httpRequest.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	httpRequest.Header.Set("version", c.Context.Version)
	httpRequest.Header.Set("timestamp", strconv.FormatInt(c.Context.Timestamp, 10))
	httpRequest.Header.Set("appKey", c.Context.Config.AppKey)
	httpRequest.Header.Set("signature", signature)
	httpRequest.Header.Set("xReqNonce", c.Context.XReqNonce)
}

type GetCloseInfoRequest struct {
	CustomerIds []int64 `json:"customerIds"` // 企业ID，最大支持500个
}

type GetCloseInfoResponse struct {
	Head util.CommonError   `json:"head"`
	Body []GetCloseInfoList `json:"body"`
}

type GetCloseInfoList struct {
	CompanyId      string `json:"companyId"`
	CustomerId     string `json:"customerId"`
	AccountSetId   string `json:"accountSetId"`
	MaxClosePeriod string `json:"maxClosePeriod"`
	CreatePeriod   string `json:"createPeriod"`
}

// GetCloseInfo 结账信息接口
func (c *Alice) GetCloseInfo(req GetCloseInfoRequest) (result GetCloseInfoResponse, err error) {
	customerIdsOfString := arrutil.AnyToString(req.CustomerIds)
	customerIdsOfString = strings.ReplaceAll(customerIdsOfString, "[", "")
	customerIdsOfString = strings.ReplaceAll(customerIdsOfString, "]", "")
	postData := url.Values{}
	postData.Add("customerIds", customerIdsOfString)

	httpRequest, err := http.NewRequest("POST", GetCloseInfoUrl, strings.NewReader(postData.Encode()))

	//httpRequest, err := http.NewRequest("POST", GetCloseInfoUrl, reader)
	c.SignatureHandle = credential.NewDefaultSignature(nil, c.AppKey, c.AppSecret, nil, nil, nil, nil, c.Timestamp, c.Version, c.XReqNonce, credential.CacheKeyYiQiYingPrefix, c.Cache)
	signature, err := c.GetSignature()
	if err != nil {
		return
	}

	c.setHeader(signature, httpRequest)
	client := &http.Client{}
	response, err := client.Do(httpRequest)
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return
	}
	if result.Head.Status != "Y" || result.Head.Code != "00000000" {
		err = fmt.Errorf("结账信息数据出错,%v,%v 出错代码为(%v)", result.Head.Msg, result.Head.Description, result.Head.Code)
		return
	}
	return
}

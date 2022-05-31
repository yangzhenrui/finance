package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/yangzhenrui/finance/credential"
	"github.com/yangzhenrui/finance/util"
	"github.com/yangzhenrui/finance/yiqiying/context"
	"io/ioutil"
	"net/http"
	"strconv"
)

const (
	// GetCloseInfo 结账信息接口
	GetCloseInfoUrl = "https://openapi.17win.com/gateway/openyqdz/alice/closeInfo/getCloseInfo"
)

type Alice struct {
	*context.Context
}

func NewAlice(ctx *context.Context) *Alice {
	return &Alice{ctx}
}

func (c *Alice) setHeader(signature string, httpRequest *http.Request) {
	httpRequest.Header.Set("Content-Type", "application/json")
	httpRequest.Header.Set("version", c.Context.Version)
	httpRequest.Header.Set("timestamp", strconv.FormatInt(c.Context.Timestamp, 10))
	httpRequest.Header.Set("appKey", c.Context.Config.AppKey)
	httpRequest.Header.Set("signature", signature)
	httpRequest.Header.Set("xReqNonce", c.Context.XReqNonce)
}

type GetCloseInfoRequest struct {
	CustomerIds []string `json:"customerIds"` // 企业ID，最大支持500个
}

type GetCloseInfoResponse struct {
	Head                     util.CommonError         `json:"head"`
	GetCloseInfoResponseBody GetCloseInfoResponseBody `json:"body"`
}

type GetCloseInfoResponseBody struct {
	GetCloseInfoList []GetCloseInfoList `json:"customerList"`
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
	aliceReq, err := json.Marshal(&req)
	reader := bytes.NewReader(aliceReq)
	httpRequest, err := http.NewRequest("POST", GetCloseInfoUrl, reader)
	c.SignatureHandle = credential.NewDefaultSignature(nil, c.AppKey, c.AppSecret, c.CustomerId, nil, nil, nil, c.Timestamp, c.Version, c.XReqNonce, credential.CacheKeyYiQiYingPrefix, c.Cache)
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

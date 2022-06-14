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
	// QueryCustomersUrl 查询客户信息接口
	QueryCustomersUrl = "https://openapi.17win.com/gateway/openyqdz/manage/customer/queryCustomers"

	// AddCustomerUrl 新增客户接口
	AddCustomerUrl = "https://openapi.17win.com/gateway/openyqdz/manage/customer/addCustomer"

	// BatchAssignRolesUrl 批量派工接口
	BatchAssignRolesUrl = "https://openapi.17win.com/gateway/openyqdz/manage/customer/batchAssignRoles"

	// UpdateCustomerUrl 更新客户信息接口
	UpdateCustomerUrl = "https://openapi.17win.com/gateway/openyqdz/manage/customer/updateCustomer"
)

type Customer struct {
	*context.Context
}

func NewCustomer(ctx *context.Context) *Customer {
	return &Customer{ctx}
}

type QueryCustomersRequest struct {
	CustomerIds          []string             `json:"customerIds,omitempty"`
	PageNo               int                  `json:"pageNo" form:"pageNo"`
	PageSize             int                  `json:"pageSize" form:"pageSize"`
	CustomerLikeCriteria CustomerLikeCriteria `json:"customerLikeCriteria,omitempty"`
}

type CustomerLikeCriteria struct {
	Name     string `json:"name,omitempty"`
	FullName string `json:"fullName,omitempty"`
	TaxNo    string `json:"taxNo,omitempty"`
}

type QueryCustomersResponse struct {
	Head                       util.CommonError           `json:"head"`
	QueryCustomersResponseBody QueryCustomersResponseBody `json:"body"`
}

type QueryCustomersResponseBody struct {
	Pager        Pager          `json:"pager"`
	Total        int            `json:"total"`
	CustomerList []CustomerList `json:"customerList"`
}

type CustomerList struct {
	CustomerId       string        `json:"customerId"`
	IndustryCategory string        `json:"industryCategory"`
	IndustryType     string        `json:"industryType"`
	Name             string        `json:"name"`
	FullName         string        `json:"fullName"`
	CustomerNo       string        `json:"customerNo"`
	TaxNo            string        `json:"taxNo"`
	LocationCode     string        `json:"locationCode"`
	TaxType          string        `json:"taxType"`
	Level            int           `json:"level"`
	Status           int           `json:"status"`
	CustomerType     int           `json:"customerType"`
	Address          string        `json:"address"`
	DepartmentId     int           `json:"departmentId"`
	AccountList      []AccountList `json:"accountList"`
}

type AccountList struct {
	AccountId        string `json:"accountId"`
	RelationShipType int    `json:"relationshipType"`
	LoginName        string `json:"loginName"`
}

type Pager struct {
	CurrentPage int `json:"currentPage"`
	PageSize    int `json:"pageSize"`
}

func (c *Customer) setHeader(signature string, httpRequest *http.Request) {
	httpRequest.Header.Set("Content-Type", "application/json")
	httpRequest.Header.Set("version", c.Context.Version)
	httpRequest.Header.Set("timestamp", strconv.FormatInt(c.Context.Timestamp, 10))
	httpRequest.Header.Set("appKey", c.Context.Config.AppKey)
	httpRequest.Header.Set("signature", signature)
	httpRequest.Header.Set("xReqNonce", c.Context.XReqNonce)
}

// QueryCustomers 查询客户信息
func (c *Customer) QueryCustomers(req QueryCustomersRequest) (result QueryCustomersResponse, err error) {
	customersReq, err := json.Marshal(&req)
	reader := bytes.NewReader(customersReq)
	httpRequest, err := http.NewRequest("POST", QueryCustomersUrl, reader)
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
		err = fmt.Errorf("查询三方平台客户信息,%v,%v 出错代码为(%v)", result.Head.Msg, result.Head.Description, result.Head.Code)
		return
	}
	return
}

type AddCustomerRequest struct {
	CustomerName      string `json:"customerName"`
	OperatorLoginName string `json:"operatorLoginName"`
}

type AddCustomerResponse struct {
	Head util.CommonError `json:"head"`
	Body string           `json:"body"`
}

// AddCustomer 新增客户
func (c *Customer) AddCustomer(req AddCustomerRequest) (result AddCustomerResponse, err error) {
	customersReq, err := json.Marshal(&req)
	reader := bytes.NewReader(customersReq)
	httpRequest, err := http.NewRequest("POST", AddCustomerUrl, reader)
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
		err = fmt.Errorf("新增客户失败,%v,%v 出错代码为(%v)", result.Head.Msg, result.Head.Description, result.Head.Code)
		return
	}
	return
}

type BatchAssignRolesRequest struct {
	CustomerIdList     []string             `json:"customerIdList"`
	OperatorLoginName  string               `json:"operatorLoginName"`
	RoleAssignmentList []RoleAssignmentList `json:"roleAssignmentList"`
}

type RoleAssignmentList struct {
	RelationShipType int      `json:"relationshipType"` // 类型, 1: 服务顾问, 2:其他服务人员, 3:税务会计,4:财务会计, 5:审核会计,6:收款负责人，7：客户经理，8：开票员
	LoginNameList    []string `json:"loginNameList"`    // 人员列表(若为空,则表示删除)
}

type BatchAssignRolesResponse struct {
	Head util.CommonError `json:"head"`
	Body string           `json:"body"`
}

// BatchAssignRoles 批量派工
func (c *Customer) BatchAssignRoles(req BatchAssignRolesRequest) (result BatchAssignRolesResponse, err error) {
	customersReq, err := json.Marshal(&req)
	reader := bytes.NewReader(customersReq)
	httpRequest, err := http.NewRequest("POST", BatchAssignRolesUrl, reader)
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
		err = fmt.Errorf("三方平台批量分配出错,%v,%v 出错代码为(%v)", result.Head.Msg, result.Head.Description, result.Head.Code)
		return
	}
	return
}

type UpdateCustomerRequest struct {
	CustomerId        string `json:"customerId"`
	OperatorLoginName string `json:"operatorLoginName"`
	CustomerNo        string `json:"customerNo"`
	Name              string `json:"name"`
	FullName          string `json:"fullName"`
	TaxNo             string `json:"taxNo"`
	IndustryCategory  string `json:"industryCategory"`
	IndustryType      string `json:"industryType"`
	LocationCode      string `json:"locationCode"`
}

type UpdateCustomerResponse struct {
	Head util.CommonError `json:"head"`
	Body string           `json:"body"`
}

// UpdateCustomer 更新客户信息
func (c *Customer) UpdateCustomer(req UpdateCustomerRequest) (result UpdateCustomerResponse, err error) {
	customersReq, err := json.Marshal(&req)
	reader := bytes.NewReader(customersReq)
	httpRequest, err := http.NewRequest("POST", UpdateCustomerUrl, reader)
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
		err = fmt.Errorf("更新客户信息失败,%v,%v 出错代码为(%v)", result.Head.Msg, result.Head.Description, result.Head.Code)
		return
	}
	return
}

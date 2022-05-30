package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/google/go-querystring/query"
	"github.com/yangzhenrui/finance/credential"
	"github.com/yangzhenrui/finance/util"
	"github.com/yangzhenrui/finance/yiqiying/context"
	"io/ioutil"
	"net/http"
	"strconv"
)

const (
	// QueryAccountBalanceSheetUrl 科目余额表接口
	QueryAccountBalanceSheetUrl = "https://openapi.17win.com/gateway/openyqdz/finance/sheetController/queryAccountBalanceSheet"

	// SelectAssetsDebtSheetUrl 资产负债表接口
	SelectAssetsDebtSheetUrl = "https://openapi.17win.com/gateway/openyqdz/finance/sheetController/selectAssetsDebtSheet"

	// SelectIncomeSheetUrl 利润表接口
	SelectIncomeSheetUrl = "https://openapi.17win.com/gateway/openyqdz/finance/sheetController/selectIncomeSheet"

	// GetMonthCashFlowsStatementSheetUrl 现金流量表接口
	GetMonthCashFlowsStatementSheetUrl = "https://openapi.17win.com/gateway/openyqdz/finance/sheetController/getMonthCashFlowsStatement"

	// SelectQuarterIncomeSheetUrl 利润表季报接口
	SelectQuarterIncomeSheetUrl = "https://openapi.17win.com/gateway/openyqdz/finance/sheetController/selectQuarterIncomeSheet"

	// GetAllYearMonthFinancialPositionStatementSheetUrl 资产负债表全年接口
	GetAllYearMonthFinancialPositionStatementSheetUrl = "https://openapi.17win.com/gateway/openyqdz/finance/sheetController/getAllYearMonthFinancialPositionStatement"

	// GetAllYearMonthIncomeStatementSheetUrl 利润表全年接口
	GetAllYearMonthIncomeStatementSheetUrl = "https://openapi.17win.com/gateway/openyqdz/finance/sheetController/getAllYearMonthIncomeStatement"

	// GetAllYearMonthCashFlowsStatementSheetUrl 现金流量表全年接口
	GetAllYearMonthCashFlowsStatementSheetUrl = "https://openapi.17win.com/gateway/openyqdz/finance/sheetController/getAllYearMonthCashFlowsStatement"
)

type Finance struct {
	*context.Context
}

func NewFinance(ctx *context.Context) *Finance {
	return &Finance{ctx}
}

func (c *Finance) setHeader(signature string, httpRequest *http.Request) {
	httpRequest.Header.Set("Content-Type", "application/json")
	httpRequest.Header.Set("version", c.Context.Version)
	httpRequest.Header.Set("timestamp", strconv.FormatInt(c.Context.Timestamp, 10))
	httpRequest.Header.Set("appKey", c.Context.Config.AppKey)
	httpRequest.Header.Set("signature", signature)
	httpRequest.Header.Set("xReqNonce", c.Context.XReqNonce)
}

type QueryAccountBalanceSheetRequest struct {
	CustomerId          string `json:"customerId"`
	BeginPeriod         string `json:"beginPeriod"`
	EndPeriod           string `json:"endPeriod"`
	PageNo              int    `json:"pageNo"`
	PageSize            int    `json:"pageSize"`
	BeginTitleCode      string `json:"beginTitleCode,omitempty"`
	EndTitleCode        string `json:"endTitleCode,omitempty"`
	TitleLevel          int    `json:"titleLevel,omitempty"` // 从1到6
	ShowTitle           bool   `json:"showTitle,omitempty"`
	ShowAssistant       bool   `json:"showAssistant,omitempty"`
	ShowYearAccumulated bool   `json:"showYearAccumulated,omitempty"`
	ShowQuantity        bool   `json:"showQuantity,omitempty"`
	FcurCode            string `json:"fcurCode,omitempty"`      // 科目启用外币时的外币编码
	AssistantType       string `json:"assistantType,omitempty"` // "c", "客户"；"s", "供应商"；"i", "存货"；"p", "项目"；"d", "部门"；"e", "员工"
	InventoryType       string `json:"inventoryType,omitempty"` // 10：库存商品  20：原材料  30：委托加工物资  40：周转材料  50：劳务或服务  90：未分类
	AssistantId         int    `json:"assistantId,omitempty"`
	ShowEndBalance0     bool   `json:"showEndBalance0,omitempty"`
	FirstAccountTitle   bool   `json:"firstAccountTitle,omitempty"`
}

type QueryAccountBalanceSheetResponse struct {
	Head                       util.CommonError                     `json:"head"`
	QueryCustomersResponseBody QueryAccountBalanceSheetResponseBody `json:"body"`
}

type QueryAccountBalanceSheetResponseBody struct {
	List []QueryAccountBalanceSheetList `json:"list"`
}

type QueryAccountBalanceSheetList struct {
	TitleId                 int     `json:"titleId"`
	TitleIsLast             bool    `json:"titleIsLast"`
	AssistantType           string  `json:"assistantType"`
	AssistantId             int     `json:"assistantId"`
	TitleCode               string  `json:"titleCode"`
	TitleName               string  `json:"titleName"`
	TitleFullName           string  `json:"titleFullName"`
	PTitleCode              string  `json:"PTitleCode"`
	Level                   int     `json:"level"`
	Type                    string  `json:"type"`
	Unit                    string  `json:"unit"`
	FcurCode                string  `json:"fcurCode"`
	Specification           string  `json:"specification"`
	AssistantName           string  `json:"assistantName"`
	BeginDirection          int     `json:"beginDirection"`
	BeginAmount             float64 `json:"beginAmount"`
	BeginQuantity           float64 `json:"beginQuantity"`
	BeginUnitPrice          float64 `json:"beginUnitPrice"`
	BeginDebit              float64 `json:"beginDebit"`
	BeginCredit             float64 `json:"beginCredit"`
	BeginDebitFcur          float64 `json:"beginDebitFcur"`
	BeginCreditFcur         float64 `json:"beginCreditFcur"`
	OccurredDebit           float64 `json:"occurredDebit"`
	OccurredDebitQuantity   float64 `json:"occurredDebitQuantity"`
	OccurredDebitFcur       float64 `json:"occurredDebitFcur"`
	OccurredCredit          float64 `json:"occurredCredit"`
	OccurredCreditQuantity  float64 `json:"occurredCreditQuantity"`
	OccurredCreditFcur      float64 `json:"occurredCreditFcur"`
	YearAccumulatedDebit    float64 `json:"yearAccumulatedDebit"`
	YearAccumulatedCredit   float64 `json:"yearAccumulatedCredit"`
	EndDebit                float64 `json:"endDebit"`
	EndCredit               float64 `json:"endCredit"`
	EndDebitFcur            float64 `json:"endDebitFcur"`
	EndCreditFcur           float64 `json:"endCreditFcur"`
	EndDirection            int     `json:"endDirection"`
	EndAmount               float64 `json:"endAmount"`
	EndQuantity             float64 `json:"endQuantity"`
	EndUnitPrice            float64 `json:"endUnitPrice"`
	OccurredDebitUnitPrice  float64 `json:"occurredDebitUnitPrice"`
	OccurredCreditUnitPrice float64 `json:"occurredCreditUnitPrice"`
}

// QueryAccountBalanceSheet 科目余额表接口
func (c *Finance) QueryAccountBalanceSheet(req QueryAccountBalanceSheetRequest) (result QueryAccountBalanceSheetResponse, err error) {
	financeReq, err := json.Marshal(&req)
	reader := bytes.NewReader(financeReq)
	httpRequest, err := http.NewRequest("POST", QueryAccountBalanceSheetUrl, reader)
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
		err = fmt.Errorf("queryCustomers error : errcode=%v , errmsg=%v, errdesc=%v", result.Head.Code, result.Head.Msg, result.Head.Description)
		return
	}
	return
}

type SelectAssetsDebtSheetRequest struct {
	CustomerId     string `json:"customerId" form:"customerId" url:"customerId"`
	AccountPeriod  string `json:"accountPeriod" form:"accountPeriod" url:"accountPeriod"`
	ReclassifyFlag string `json:"reclassifyFlag,omitempty" form:"reclassifyFlag" url:"reclassifyFlag"`
}

type SelectAssetsDebtSheetResponse struct {
	Head util.CommonError            `json:"head"`
	Body []SelectAssetsDebtSheetList `json:"body"`
}

type SelectAssetsDebtSheetList struct {
	AccountTitleName        string  `json:"accountTitleName"`        // 科目名称
	Row                     int     `json:"row"`                     // 行次
	Number                  int     `json:"number"`                  // 序号
	Level                   int     `json:"level"`                   // 级次
	PRowNum                 int     `json:"PRowNum"`                 // 父级序号
	Warn                    bool    `json:"warn"`                    // 是否是异常数据
	ShowLine                int     `json:"showLine"`                // 是否显示行号
	BalanceEnd              float64 `json:"balanceEnd"`              // 期末余额
	YearBeginBalance        float64 `json:"yearBeginBalance"`        // 年初余额
	FomularDetail           string  `json:"fomularDetail"`           // 公式注释说明
	LimitFomularDetail      string  `json:"limitFomularDetail"`      // 限定性注释说明
	NonLimitFomularDetail   string  `json:"nonLimitFomularDetail"`   // 非限定性注释说明
	LimitOccurreAmount      float64 `json:"limitOccurreAmount"`      // 限定性本月发生额
	NonLimitOccurreAmout    float64 `json:"nonLimitOccurreAmout"`    // 非限定性本月发生额
	LimitYearAccumulated    float64 `json:"limitYearAccumulated"`    // 限定性本年累计发生额
	NonLimitYearAccumulated float64 `json:"nonLimitYearAccumulated"` // 非限定性本年累计发生额
	YearAccumulated         float64 `json:"yearAccumulated"`         // 本年累计发生额（小企业-利-月季年）（企业制度-利-月季年）
	QuarterOne              float64 `json:"quarterOne"`              // 第一季度
	QuarterTwo              float64 `json:"quarterTwo"`              // 第二季度
	QuarterThree            float64 `json:"quarterThree"`            // 第三季度
	QuarterFour             float64 `json:"quarterFour"`             // 第四季度
	OccurredAmount          float64 `json:"occurredAmount"`          // 本月发生额 （小企业-利-月） （企业制度-利-月）
	PreYearAccumulated      float64 `json:"preYearAccumulated"`      // 上一年累计发生额
	AmountOfLocalPeriod     float64 `json:"amountOfLocalPeriod"`     // 本期发生额（企业准则-利-月季年） （小企业-利-季）（企业制度-利-季年）
	AmountOfPrePeriod       float64 `json:"amountOfPrePeriod"`       // 上期发生额 （小企业-利-年）
}

// SelectAssetsDebtSheet 资产负债表接口
func (c *Finance) SelectAssetsDebtSheet(req SelectAssetsDebtSheetRequest) (result SelectAssetsDebtSheetResponse, err error) {
	uriArr, _ := query.Values(req)
	hasRf := 1
	if uriArr.Get("reclassifyFlag") == "" {
		uriArr.Del("reclassifyFlag")
		hasRf = 0
	}
	url := fmt.Sprintf("%v?%v", SelectAssetsDebtSheetUrl, uriArr.Encode())
	httpRequest, err := http.NewRequest("GET", url, nil)

	c.SignatureHandle = credential.NewDefaultSignature(&req.AccountPeriod, c.AppKey, c.AppSecret, &req.CustomerId, nil, &req.ReclassifyFlag, nil, c.Timestamp, c.Version, c.XReqNonce, credential.CacheKeyYiQiYingPrefix, c.Cache)
	signature, err := c.GetSignature()
	if err != nil {
		return
	}

	// 请求头设置
	c.setHeader(signature, httpRequest)
	httpRequest.Header.Set("customerId", req.CustomerId)
	httpRequest.Header.Set("accountPeriod", req.AccountPeriod)
	if hasRf == 1 {
		httpRequest.Header.Set("reclassifyFlag", req.ReclassifyFlag)
	}

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
		err = fmt.Errorf("queryCustomers error : errcode=%v , errmsg=%v, errdesc=%v", result.Head.Code, result.Head.Msg, result.Head.Description)
		return
	}
	return
}

type SelectIncomeSheetRequest struct {
	CustomerId    string `json:"customerId" form:"customerId" url:"customerId"`
	AccountPeriod string `json:"accountPeriod" form:"accountPeriod" url:"accountPeriod"`
}

type SelectIncomeSheetResponse struct {
	Head util.CommonError        `json:"head"`
	Body []SelectIncomeSheetList `json:"body"`
}

type SelectIncomeSheetList struct {
	AccountTitleName        string  `json:"accountTitleName"`        // 科目名称
	Row                     int     `json:"row"`                     // 行次
	Number                  int     `json:"number"`                  // 序号
	Level                   int     `json:"level"`                   // 级次
	PRowNum                 int     `json:"PRowNum"`                 // 父级序号
	Warn                    bool    `json:"warn"`                    // 是否是异常数据
	ShowLine                int     `json:"showLine"`                // 是否显示行号
	BalanceEnd              float64 `json:"balanceEnd"`              // 期末余额
	YearBeginBalance        float64 `json:"yearBeginBalance"`        // 年初余额
	FomularDetail           string  `json:"fomularDetail"`           // 公式注释说明
	LimitFomularDetail      string  `json:"limitFomularDetail"`      // 限定性注释说明
	NonLimitFomularDetail   string  `json:"nonLimitFomularDetail"`   // 非限定性注释说明
	LimitOccurreAmount      float64 `json:"limitOccurreAmount"`      // 限定性本月发生额
	NonLimitOccurreAmout    float64 `json:"nonLimitOccurreAmout"`    // 非限定性本月发生额
	LimitYearAccumulated    float64 `json:"limitYearAccumulated"`    // 限定性本年累计发生额
	NonLimitYearAccumulated float64 `json:"nonLimitYearAccumulated"` // 非限定性本年累计发生额
	YearAccumulated         float64 `json:"yearAccumulated"`         // 本年累计发生额（小企业-利-月季年）（企业制度-利-月季年）
	QuarterOne              float64 `json:"quarterOne"`              // 第一季度
	QuarterTwo              float64 `json:"quarterTwo"`              // 第二季度
	QuarterThree            float64 `json:"quarterThree"`            // 第三季度
	QuarterFour             float64 `json:"quarterFour"`             // 第四季度
	OccurredAmount          float64 `json:"occurredAmount"`          // 本月发生额 （小企业-利-月） （企业制度-利-月）
	PreYearAccumulated      float64 `json:"preYearAccumulated"`      // 上一年累计发生额
	AmountOfLocalPeriod     float64 `json:"amountOfLocalPeriod"`     // 本期发生额（企业准则-利-月季年） （小企业-利-季）（企业制度-利-季年）
	AmountOfPrePeriod       float64 `json:"amountOfPrePeriod"`       // 上期发生额 （小企业-利-年）
}

// SelectIncomeSheet 利润表接口
func (c *Finance) SelectIncomeSheet(req SelectIncomeSheetRequest) (result SelectIncomeSheetResponse, err error) {
	uriArr, _ := query.Values(req)
	url := fmt.Sprintf("%v?%v", SelectIncomeSheetUrl, uriArr.Encode())
	httpRequest, err := http.NewRequest("GET", url, nil)

	c.SignatureHandle = credential.NewDefaultSignature(&req.AccountPeriod, c.AppKey, c.AppSecret, &req.CustomerId, nil, nil, nil, c.Timestamp, c.Version, c.XReqNonce, credential.CacheKeyYiQiYingPrefix, c.Cache)
	signature, err := c.GetSignature()
	if err != nil {
		return
	}

	// 请求头设置
	c.setHeader(signature, httpRequest)
	httpRequest.Header.Set("customerId", req.CustomerId)
	httpRequest.Header.Set("accountPeriod", req.AccountPeriod)

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
		err = fmt.Errorf("queryCustomers error : errcode=%v , errmsg=%v, errdesc=%v", result.Head.Code, result.Head.Msg, result.Head.Description)
		return
	}
	return
}

type GetMonthCashFlowsStatementSheetRequest struct {
	CustomerId    string `json:"customerId" form:"customerId" url:"customerId"`
	AccountPeriod string `json:"accountPeriod" form:"accountPeriod" url:"accountPeriod"`
}

type GetMonthCashFlowsStatementSheetResponse struct {
	Head util.CommonError                      `json:"head"`
	Body []GetMonthCashFlowsStatementSheetList `json:"body"`
}

type GetMonthCashFlowsStatementSheetList struct {
	Name                  string  `json:"name"`                  // 项目名称
	Line                  int     `json:"line"`                  // 行次（用于界面显示）
	RowNum                int     `json:"rowNum"`                // 行号（区别于行次，没有行次的行，也会有行号）
	PRowNum               int     `json:"PRowNum"`               // 父级行号（可根据些参数构建层级结构）
	YearAccumulatedAmount float64 `json:"yearAccumulatedAmount"` // 本年累计金额（小企业会计准则）
	CurrentMonthAmount    float64 `json:"currentMonthAmount"`    // 本月金额（小企业会计准则）
	CurrentPeriodAmount   float64 `json:"currentPeriodAmount"`   // 本期金额（企业会计准则）
	PrevPeriodAmount      float64 `json:"prevPeriodAmount"`      // 上期金额（企业会计准则）
	PrevYearAmount        float64 `json:"prevYearAmount"`        // 上年累计发生额
}

// GetMonthCashFlowsStatementSheet 现金流量表接口
func (c *Finance) GetMonthCashFlowsStatementSheet(req GetMonthCashFlowsStatementSheetRequest) (result GetMonthCashFlowsStatementSheetResponse, err error) {
	financeReq, err := json.Marshal(&req)
	reader := bytes.NewReader(financeReq)
	httpRequest, err := http.NewRequest("GET", GetMonthCashFlowsStatementSheetUrl, reader)
	signature, err := c.GetSignature()
	if err != nil {
		return
	}

	// 请求头设置
	c.setHeader(signature, httpRequest)
	httpRequest.Header.Set("customerId", req.CustomerId)
	httpRequest.Header.Set("accountPeriod", req.AccountPeriod)

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
		err = fmt.Errorf("queryCustomers error : errcode=%v , errmsg=%v, errdesc=%v", result.Head.Code, result.Head.Msg, result.Head.Description)
		return
	}
	return
}

type SelectQuarterIncomeSheetRequest struct {
	CustomerId    string `json:"customerId" form:"customerId" url:"customerId"`
	AccountPeriod string `json:"accountPeriod" form:"accountPeriod" url:"accountPeriod"`
}

type SelectQuarterIncomeSheetResponse struct {
	Head util.CommonError               `json:"head"`
	Body []SelectQuarterIncomeSheetList `json:"body"`
}

type SelectQuarterIncomeSheetList struct {
	AccountTitleName        string  `json:"accountTitleName"`        // 科目名称
	Row                     int     `json:"row"`                     // 行次
	Number                  int     `json:"number"`                  // 序号
	Level                   int     `json:"level"`                   // 级次
	PRowNum                 int     `json:"PRowNum"`                 // 父级序号
	Warn                    bool    `json:"warn"`                    // 是否是异常数据
	ShowLine                int     `json:"showLine"`                // 是否显示行号
	BalanceEnd              float64 `json:"balanceEnd"`              // 期末余额
	YearBeginBalance        float64 `json:"yearBeginBalance"`        // 年初余额
	FomularDetail           string  `json:"fomularDetail"`           // 公式注释说明
	LimitFomularDetail      string  `json:"limitFomularDetail"`      // 限定性注释说明
	NonLimitFomularDetail   string  `json:"nonLimitFomularDetail"`   // 非限定性注释说明
	LimitOccurreAmount      float64 `json:"limitOccurreAmount"`      // 限定性本月发生额
	NonLimitOccurreAmout    float64 `json:"nonLimitOccurreAmout"`    // 非限定性本月发生额
	LimitYearAccumulated    float64 `json:"limitYearAccumulated"`    // 限定性本年累计发生额
	NonLimitYearAccumulated float64 `json:"nonLimitYearAccumulated"` // 非限定性本年累计发生额
	YearAccumulated         float64 `json:"yearAccumulated"`         // 本年累计发生额（小企业-利-月季年）（企业制度-利-月季年）
	QuarterOne              float64 `json:"quarterOne"`              // 第一季度
	QuarterTwo              float64 `json:"quarterTwo"`              // 第二季度
	QuarterThree            float64 `json:"quarterThree"`            // 第三季度
	QuarterFour             float64 `json:"quarterFour"`             // 第四季度
	OccurredAmount          float64 `json:"occurredAmount"`          // 本月发生额 （小企业-利-月） （企业制度-利-月）
	PreYearAccumulated      float64 `json:"preYearAccumulated"`      // 上一年累计发生额
	AmountOfLocalPeriod     float64 `json:"amountOfLocalPeriod"`     // 本期发生额（企业准则-利-月季年） （小企业-利-季）（企业制度-利-季年）
	AmountOfPrePeriod       float64 `json:"amountOfPrePeriod"`       // 上期发生额 （小企业-利-年）
}

// SelectQuarterIncomeSheet 利润表季报接口
func (c *Finance) SelectQuarterIncomeSheet(req SelectQuarterIncomeSheetRequest) (result SelectQuarterIncomeSheetResponse, err error) {
	uriArr, _ := query.Values(req)
	url := fmt.Sprintf("%v?%v", SelectQuarterIncomeSheetUrl, uriArr.Encode())
	httpRequest, err := http.NewRequest("GET", url, nil)

	c.SignatureHandle = credential.NewDefaultSignature(&req.AccountPeriod, c.AppKey, c.AppSecret, &req.CustomerId, nil, nil, nil, c.Timestamp, c.Version, c.XReqNonce, credential.CacheKeyYiQiYingPrefix, c.Cache)
	signature, err := c.GetSignature()
	if err != nil {
		return
	}

	// 请求头设置
	c.setHeader(signature, httpRequest)
	httpRequest.Header.Set("customerId", req.CustomerId)
	httpRequest.Header.Set("accountPeriod", req.AccountPeriod)

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
		err = fmt.Errorf("queryCustomers error : errcode=%v , errmsg=%v, errdesc=%v", result.Head.Code, result.Head.Msg, result.Head.Description)
		return
	}
	return
}

type GetAllYearMonthFinancialPositionStatementSheetRequest struct {
	CustomerId     string `json:"customerId" form:"customerId" url:"customerId"`
	AccountPeriod  string `json:"accountPeriod" form:"accountPeriod" url:"accountPeriod"`
	ReclassifyFlag string `json:"reclassifyFlag" form:"reclassifyFlag" url:"reclassifyFlag"`
}

type GetAllYearMonthFinancialPositionStatementSheetResponse struct {
	Head util.CommonError                                     `json:"head"`
	Body []GetAllYearMonthFinancialPositionStatementSheetList `json:"body"`
}

type GetAllYearMonthFinancialPositionStatementSheetList struct {
	AccountTitleName        string  `json:"accountTitleName"`        // 科目名称
	Row                     int     `json:"row"`                     // 行次
	Number                  int     `json:"number"`                  // 序号
	Level                   int     `json:"level"`                   // 级次
	PRowNum                 int     `json:"PRowNum"`                 // 父级序号
	Warn                    bool    `json:"warn"`                    // 是否是异常数据
	ShowLine                int     `json:"showLine"`                // 是否显示行号
	BalanceEnd              float64 `json:"balanceEnd"`              // 期末余额
	YearBeginBalance        float64 `json:"yearBeginBalance"`        // 年初余额
	FomularDetail           string  `json:"fomularDetail"`           // 公式注释说明
	LimitFomularDetail      string  `json:"limitFomularDetail"`      // 限定性注释说明
	NonLimitFomularDetail   string  `json:"nonLimitFomularDetail"`   // 非限定性注释说明
	LimitOccurreAmount      float64 `json:"limitOccurreAmount"`      // 限定性本月发生额
	NonLimitOccurreAmout    float64 `json:"nonLimitOccurreAmout"`    // 非限定性本月发生额
	LimitYearAccumulated    float64 `json:"limitYearAccumulated"`    // 限定性本年累计发生额
	NonLimitYearAccumulated float64 `json:"nonLimitYearAccumulated"` // 非限定性本年累计发生额
	YearAccumulated         float64 `json:"yearAccumulated"`         // 本年累计发生额（小企业-利-月季年）（企业制度-利-月季年）
	QuarterOne              float64 `json:"quarterOne"`              // 第一季度
	QuarterTwo              float64 `json:"quarterTwo"`              // 第二季度
	QuarterThree            float64 `json:"quarterThree"`            // 第三季度
	QuarterFour             float64 `json:"quarterFour"`             // 第四季度
	OccurredAmount          float64 `json:"occurredAmount"`          // 本月发生额 （小企业-利-月） （企业制度-利-月）
	PreYearAccumulated      float64 `json:"preYearAccumulated"`      // 上一年累计发生额
	AmountOfLocalPeriod     float64 `json:"amountOfLocalPeriod"`     // 本期发生额（企业准则-利-月季年） （小企业-利-季）（企业制度-利-季年）
	AmountOfPrePeriod       float64 `json:"amountOfPrePeriod"`       // 上期发生额 （小企业-利-年）
}

// GetAllYearMonthFinancialPositionStatementSheet 资产负债表全年接口
func (c *Finance) GetAllYearMonthFinancialPositionStatementSheet(req GetAllYearMonthFinancialPositionStatementSheetRequest) (result GetAllYearMonthFinancialPositionStatementSheetResponse, err error) {
	uriArr, _ := query.Values(req)
	hasRf := 1
	if uriArr.Get("reclassifyFlag") == "" {
		uriArr.Del("reclassifyFlag")
		hasRf = 0
	}
	url := fmt.Sprintf("%v?%v", GetAllYearMonthFinancialPositionStatementSheetUrl, uriArr.Encode())
	httpRequest, err := http.NewRequest("GET", url, nil)

	c.SignatureHandle = credential.NewDefaultSignature(&req.AccountPeriod, c.AppKey, c.AppSecret, &req.CustomerId, nil, &req.ReclassifyFlag, nil, c.Timestamp, c.Version, c.XReqNonce, credential.CacheKeyYiQiYingPrefix, c.Cache)
	signature, err := c.GetSignature()
	if err != nil {
		return
	}

	// 请求头设置
	c.setHeader(signature, httpRequest)
	httpRequest.Header.Set("customerId", req.CustomerId)
	httpRequest.Header.Set("accountPeriod", req.AccountPeriod)
	if hasRf == 1 {
		httpRequest.Header.Set("reclassifyFlag", req.ReclassifyFlag)
	}
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
		err = fmt.Errorf("queryCustomers error : errcode=%v , errmsg=%v, errdesc=%v", result.Head.Code, result.Head.Msg, result.Head.Description)
		return
	}
	return
}

type GetAllYearMonthIncomeStatementSheetRequest struct {
	CustomerId    string `json:"customerId" form:"customerId" url:"customerId"`
	AccountPeriod string `json:"accountPeriod" form:"accountPeriod" url:"accountPeriod"`
}

type GetAllYearMonthIncomeStatementSheetResponse struct {
	Head util.CommonError                          `json:"head"`
	Body []GetAllYearMonthIncomeStatementSheetList `json:"body"`
}

type GetAllYearMonthIncomeStatementSheetList struct {
	AccountTitleName        string  `json:"accountTitleName"`        // 科目名称
	Row                     int     `json:"row"`                     // 行次
	Number                  int     `json:"number"`                  // 序号
	Level                   int     `json:"level"`                   // 级次
	PRowNum                 int     `json:"PRowNum"`                 // 父级序号
	Warn                    bool    `json:"warn"`                    // 是否是异常数据
	ShowLine                int     `json:"showLine"`                // 是否显示行号
	BalanceEnd              float64 `json:"balanceEnd"`              // 期末余额
	YearBeginBalance        float64 `json:"yearBeginBalance"`        // 年初余额
	FomularDetail           string  `json:"fomularDetail"`           // 公式注释说明
	LimitFomularDetail      string  `json:"limitFomularDetail"`      // 限定性注释说明
	NonLimitFomularDetail   string  `json:"nonLimitFomularDetail"`   // 非限定性注释说明
	LimitOccurreAmount      float64 `json:"limitOccurreAmount"`      // 限定性本月发生额
	NonLimitOccurreAmout    float64 `json:"nonLimitOccurreAmout"`    // 非限定性本月发生额
	LimitYearAccumulated    float64 `json:"limitYearAccumulated"`    // 限定性本年累计发生额
	NonLimitYearAccumulated float64 `json:"nonLimitYearAccumulated"` // 非限定性本年累计发生额
	YearAccumulated         float64 `json:"yearAccumulated"`         // 本年累计发生额（小企业-利-月季年）（企业制度-利-月季年）
	QuarterOne              float64 `json:"quarterOne"`              // 第一季度
	QuarterTwo              float64 `json:"quarterTwo"`              // 第二季度
	QuarterThree            float64 `json:"quarterThree"`            // 第三季度
	QuarterFour             float64 `json:"quarterFour"`             // 第四季度
	OccurredAmount          float64 `json:"occurredAmount"`          // 本月发生额 （小企业-利-月） （企业制度-利-月）
	PreYearAccumulated      float64 `json:"preYearAccumulated"`      // 上一年累计发生额
	AmountOfLocalPeriod     float64 `json:"amountOfLocalPeriod"`     // 本期发生额（企业准则-利-月季年） （小企业-利-季）（企业制度-利-季年）
	AmountOfPrePeriod       float64 `json:"amountOfPrePeriod"`       // 上期发生额 （小企业-利-年）
}

// GetAllYearMonthIncomeStatementSheet 利润表全年接口
func (c *Finance) GetAllYearMonthIncomeStatementSheet(req GetAllYearMonthIncomeStatementSheetRequest) (result GetAllYearMonthIncomeStatementSheetResponse, err error) {
	uriArr, _ := query.Values(req)
	url := fmt.Sprintf("%v?%v", GetAllYearMonthIncomeStatementSheetUrl, uriArr.Encode())
	httpRequest, err := http.NewRequest("GET", url, nil)

	c.SignatureHandle = credential.NewDefaultSignature(&req.AccountPeriod, c.AppKey, c.AppSecret, &req.CustomerId, nil, nil, nil, c.Timestamp, c.Version, c.XReqNonce, credential.CacheKeyYiQiYingPrefix, c.Cache)
	signature, err := c.GetSignature()
	if err != nil {
		return
	}

	// 请求头设置
	c.setHeader(signature, httpRequest)
	httpRequest.Header.Set("customerId", req.CustomerId)
	httpRequest.Header.Set("accountPeriod", req.AccountPeriod)

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
		err = fmt.Errorf("queryCustomers error : errcode=%v , errmsg=%v, errdesc=%v", result.Head.Code, result.Head.Msg, result.Head.Description)
		return
	}
	return
}

type GetAllYearMonthCashFlowsStatementSheetRequest struct {
	CustomerId    string `json:"customerId" form:"customerId" url:"customerId"`
	AccountPeriod string `json:"accountPeriod" form:"accountPeriod" url:"accountPeriod"`
}

type GetAllYearMonthCashFlowsStatementSheetResponse struct {
	Head util.CommonError                             `json:"head"`
	Body []GetAllYearMonthCashFlowsStatementSheetList `json:"body"`
}

type GetAllYearMonthCashFlowsStatementSheetList struct {
	Name                  string  `json:"name"`                  // 项目名称
	Line                  int     `json:"line"`                  // 行次（用于界面显示）
	RowNum                int     `json:"rowNum"`                // 行号（区别于行次，没有行次的行，也会有行号）
	PRowNum               int     `json:"PRowNum"`               // 父级行号（可根据些参数构建层级结构）
	YearAccumulatedAmount float64 `json:"yearAccumulatedAmount"` // 本年累计金额（小企业会计准则）
	CurrentMonthAmount    float64 `json:"currentMonthAmount"`    // 本月金额（小企业会计准则）
	CurrentPeriodAmount   float64 `json:"currentPeriodAmount"`   // 本期金额（企业会计准则）
	PrevPeriodAmount      float64 `json:"prevPeriodAmount"`      // 上期金额（企业会计准则）
	PrevYearAmount        float64 `json:"prevYearAmount"`        // 上年累计发生额
}

// GetAllYearMonthCashFlowsStatementSheet 现金流量表全年接口
func (c *Finance) GetAllYearMonthCashFlowsStatementSheet(req GetAllYearMonthCashFlowsStatementSheetRequest) (result GetAllYearMonthCashFlowsStatementSheetResponse, err error) {
	uriArr, _ := query.Values(req)
	url := fmt.Sprintf("%v?%v", GetAllYearMonthCashFlowsStatementSheetUrl, uriArr.Encode())
	httpRequest, err := http.NewRequest("GET", url, nil)

	c.SignatureHandle = credential.NewDefaultSignature(&req.AccountPeriod, c.AppKey, c.AppSecret, &req.CustomerId, nil, nil, nil, c.Timestamp, c.Version, c.XReqNonce, credential.CacheKeyYiQiYingPrefix, c.Cache)
	signature, err := c.GetSignature()
	if err != nil {
		return
	}

	// 请求头设置
	c.setHeader(signature, httpRequest)
	httpRequest.Header.Set("customerId", req.CustomerId)
	httpRequest.Header.Set("accountPeriod", req.AccountPeriod)

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
		err = fmt.Errorf("queryCustomers error : errcode=%v , errmsg=%v, errdesc=%v", result.Head.Code, result.Head.Msg, result.Head.Description)
		return
	}
	return
}

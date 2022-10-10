package service

import (
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
	// GetTaxListUrl 查询税种信息接口
	GetTaxListUrl = "https://openapi.17win.com/gateway/openyqdz/tax/info/getTaxList"

	// GetReportUrl 查询税种报表数据接口
	GetReportUrl = "https://openapi.17win.com/gateway/openyqdz/tax/info/getReport"
)

type Tax struct {
	*context.Context
}

func NewTax(ctx *context.Context) *Tax {
	return &Tax{ctx}
}

func (c *Tax) setHeader(signature string, httpRequest *http.Request) {
	httpRequest.Header.Set("Content-Type", "application/json")
	httpRequest.Header.Set("version", c.Context.Version)
	httpRequest.Header.Set("timestamp", strconv.FormatInt(c.Context.Timestamp, 10))
	httpRequest.Header.Set("appKey", c.Context.Config.AppKey)
	httpRequest.Header.Set("signature", signature)
	httpRequest.Header.Set("xReqNonce", c.Context.XReqNonce)
}

type GetTaxListRequest struct {
	CustomerId string `json:"customerId" form:"customerId" url:"customerId"`
	Period     string `json:"period" form:"period" url:"period"`
	TaxCode    string `json:"taxCode" form:"taxCode" url:"taxCode"`
}

type GetTaxListResponse struct {
	Head util.CommonError `json:"head"`
	Body []GetTaxList     `json:"body"`
}

type GetTaxList struct {
	TaxName                      string  `json:"taxName"`                      // 税种名称
	TaxCode                      string  `json:"taxCode"`                      // 税种code
	CurrentSales                 float64 `json:"currentSales"`                 // 本期销售额
	InputTax                     float64 `json:"inputTax"`                     // 进项税额
	FinalTaxCredit               float64 `json:"finalTaxCredit"`               // 期末留抵税额
	CumulativeSales              float64 `json:"cumulativeSales"`              // 累计销售收入
	PrepaymentAmount             float64 `json:"prepaymentAmount"`             // 本期预缴税额
	TotalAmount                  float64 `json:"totalAmount"`                  // 资产总额
	NetProfit                    float64 `json:"netProfit"`                    // 本期净利润
	OutOfBusinessCumulativeSales float64 `json:"outOfBusinessCumulativeSales"` // 营业外本年累计销售收入: 小企业会计准则:取利润表“营业外收入”“ 本年累计金额 ”栏次金额； 企业会计准则:取利润表“ 营业外收入”“本期金额”栏次金额； 企业会计制度：取利润表“营业外收入”“本年累计数”栏次金额；
	PostDate                     string  `json:"postDate"`                     // 申报时间
	DeclarationState             string  `json:"declarationState"`             // 申报状态
	CreditTax                    float64 `json:"creditTax"`                    // 应补退税额
	CurrentSalesRevenue          float64 `json:"currentSalesRevenue"`          // 本期销售收入
	IncomeTotalAmount            float64 `json:"incomeTotalAmount"`            // 收入总额
	ProfitTotalAmount            float64 `json:"profitTotalAmount"`            // 收入总额
	IncomeReliefTax              float64 `json:"incomeReliefTax"`              // 减免所得税额
	DeadLine                     string  `json:"deadLine"`                     // 申报期限
	Period                       string  `json:"period"`                       // 所属期
	TotalNetProfit               float64 `json:"totalNetProfit"`               // 累计净利润
	PreNetProfit                 float64 `json:"preNetProfit"`                 // 累计净利润
	PayStatus                    int     `json:"payStatus"`                    // 缴款状态(1:申报成功、2:申报成功 未缴款、3:缴款成功、4:部分缴款、5:申报成功 已缓缴)
	PeriodBegin                  string  `json:"periodBegin"`                  // 所属期起 格式：yyyy-MM-dd。 例子：2020-10-01
	PeriodEnd                    string  `json:"periodEnd"`                    // 所属期止 格式：yyyy-MM-dd。 例子：2020-10-01
	PayAmt                       float64 `json:"payAmt"`                       // 缴税金额
	CopyTaxState                 string  `json:"copyTaxState"`                 // 抄税状态
	ClearCardState               string  `json:"clearCardState"`               // 清卡状态
	Fjss                         []Fjss  `json:"fjss"`                         // 附加税应补退税额
}

// Fjss 附加税应补退税额
type Fjss struct {
	Zsxmdm string  `json:"zsxmdm"` // 征收项目代码
	Name   string  `json:"name"`   // 附加税名称
	Value  float64 `json:"value"`  // 附加税值
}

// GetTaxList 查询税种信息接口
func (c *Tax) GetTaxList(req GetTaxListRequest) (result GetTaxListResponse, err error) {
	uriArr, _ := query.Values(req)
	hasTc := 1
	if uriArr.Get("taxCode") == "" {
		uriArr.Del("taxCode")
		hasTc = 0
	}
	url := fmt.Sprintf("%v?%v", GetTaxListUrl, uriArr.Encode())
	httpRequest, err := http.NewRequest("GET", url, nil)

	c.SignatureHandle = credential.NewDefaultSignature(nil, c.AppKey, c.AppSecret, &req.CustomerId, &req.Period, nil, &req.TaxCode, c.Timestamp, c.Version, c.XReqNonce, credential.CacheKeyYiQiYingPrefix, c.Cache)
	signature, err := c.GetSignature()
	if err != nil {
		return
	}

	// 请求头设置
	c.setHeader(signature, httpRequest)
	if c.CustomerId != nil {
		httpRequest.Header.Set("customerId", req.CustomerId)
	}
	if c.Period != nil {
		httpRequest.Header.Set("period", req.Period)
	}
	if hasTc == 1 {
		httpRequest.Header.Set("taxCode", req.TaxCode)
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
		err = fmt.Errorf("查询税种信息出错,%v,%v 出错代码为(%v)", result.Head.Msg, result.Head.Description, result.Head.Code)
		return
	}
	return
}

type GetReportRequest struct {
	CustomerId string `json:"customerId" url:"customerId"`
	Period     string `json:"period" url:"period"`
	TaxCode    string `json:"taxCode" url:"taxCode"`
}

type GetReportResponse struct {
	Head util.CommonError `json:"head"`
	Body []GetReport      `json:"body"`
}

type GetReport struct {
	PeriodBegin   string         `json:"periodBegin"` //
	PeriodEnd     string         `json:"periodEnd"`
	PostDate      string         `json:"postDate"`
	TaxTypeEnum   string         `json:"taxTypeEnum"`
	ReportDTOMap  []ReportDTOMap `json:"reportDTOMap"`
	OtherParamMap OtherParamMap  `json:"otherParamMap"`
}

type ReportDTOMap struct {
	OT_ATTACH_SEASON OT_ATTACH_SEASON `json:"OT_ATTACH_SEASON"`
}

type OT_ATTACH_SEASON struct {
	Head     OasHead       `json:"head"`
	BodyList []OasBodyList `json:"bodyList"`
}

type OasHead struct {
	CustomerId         string `json:"customerId"`
	DeclarationStateId string `json:"declarationStateId"`
	Period             string `json:"period"`
	ReportSource       int    `json:"reportSource"`
	ReadOnly           int    `json:"readOnly"`
	TaxpayerNo         string `json:"taxpayerNo"`
	TaxpayerName       string `json:"taxpayerName"`
	FillDateShow       string `json:"fillDateShow"`
	AmountUnit         string `json:"amountUnit"`
	DeclareDateShow    int    `json:"declareDateShow"`
	TaxPeriodShow      string `json:"taxPeriodShow"`
	TaxCode            string `json:"taxCode"`
	ReportId           string `json:"reportId"`
	PeriodBegin        string `json:"periodBegin"`
	PeriodEnd          string `json:"periodEnd"`
	TemplateId         string `json:"templateId"`
}

type OasBodyList struct {
	Id                      int     `json:"id"`
	CustomerId              string  `json:"customerId"`
	DeclarationStateId      int     `json:"declarationStateId"`
	Period                  string  `json:"period"`
	ReportSource            int     `json:"reportSource"`
	ItemCode                string  `json:"itemCode"`
	ItemIndex               int     `json:"itemIndex"`
	ItemLineShow            int     `json:"itemLineShow"`
	ItemLine                int     `json:"itemLine"`
	ItemLineChar            int     `json:"itemLineChar"`
	ReadOnly                string  `json:"readOnly"`
	DisplayStyle            int     `json:"displayStyle"`
	Indent                  int     `json:"indent"`
	PeriodStartAmount       float64 `json:"periodStartAmount"`
	PeriodAmount            float64 `json:"periodAmount"`
	PeriodShouldMinusAmount float64 `json:"periodShouldMinusAmount"`
	PeriodActualMinusAmount float64 `json:"periodActualMinusAmount"`
	PeriodEndAmount         float64 `json:"periodEndAmount"`
	PeriodReductionAmount   float64 `json:"periodReductionAmount"`
	ItemLineShowForDeclare  string  `json:"itemLineShowForDeclare"`
}

type OtherParamMap struct {
}

// GetReport 查询税种报表数据接口
func (c *Tax) GetReport(req GetReportRequest) (result GetReportResponse, err error) {
	uriArr, _ := query.Values(req)
	url := fmt.Sprintf("%v?%v", GetReportUrl, uriArr.Encode())
	httpRequest, err := http.NewRequest("GET", url, nil)

	c.SignatureHandle = credential.NewDefaultSignature(nil, c.AppKey, c.AppSecret, &req.CustomerId, &req.Period, nil, &req.TaxCode, c.Timestamp, c.Version, c.XReqNonce, credential.CacheKeyYiQiYingPrefix, c.Cache)
	signature, err := c.GetSignature()
	if err != nil {
		return
	}

	// 请求头设置
	c.setHeader(signature, httpRequest)
	httpRequest.Header.Set("customerId", req.CustomerId)
	httpRequest.Header.Set("period", req.Period)
	httpRequest.Header.Set("taxCode", req.TaxCode)
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
		err = fmt.Errorf("查询税种报表数据出错,%v,%v 出错代码为(%v)", result.Head.Msg, result.Head.Description, result.Head.Code)
		return
	}
	return
}

package yiqiying

import (
	uuid "github.com/satori/go.uuid"
	"github.com/yangzhenrui/finance/credential"
	"github.com/yangzhenrui/finance/yiqiying/config"
	"github.com/yangzhenrui/finance/yiqiying/context"
	service2 "github.com/yangzhenrui/finance/yiqiying/service"
	"strings"
	"time"
)

// YiQiYing 亿企赢
type YiQiYing struct {
	ctx *context.Context
}

// NewYiQiYing 实例化亿企赢API
func NewYiQiYing(cfg *config.Config) *YiQiYing {
	timestamp := time.Now().UnixNano() / 1e6
	xReqNonce := strings.Replace(uuid.NewV4().String(), "-", "", -1)
	version := "1.0.0"
	//signatureHandle := credential.NewDefaultSignature(cfg.AppKey, cfg.AppSecret, timestamp, version, xReqNonce, credential.CacheKeyYiQiYingPrefix, cfg.Cache)

	ctx := &context.Context{
		Config: cfg,
		//SignatureHandle: signatureHandle,
		Timestamp: timestamp,
		XReqNonce: xReqNonce,
		Version:   version,
	}
	return &YiQiYing{ctx}
}

// SetSignatureHandle 自定义signature获取方式
func (yqy *YiQiYing) SetSignatureHandle(signatureHandle credential.SignatureHandle) {
	yqy.ctx.SignatureHandle = signatureHandle
}

// GetContext get Context
func (yqy *YiQiYing) GetContext() *context.Context {
	return yqy.ctx
}

// GetCustomers 客户信息
func (yqy *YiQiYing) GetCustomers() *service2.Customer {
	return service2.NewCustomer(yqy.ctx)
}

// GetAlice 结账信息
func (yqy *YiQiYing) GetAlice() *service2.Alice {
	return service2.NewAlice(yqy.ctx)
}

// GetFinance 资金信息
func (yqy *YiQiYing) GetFinance() *service2.Finance {
	return service2.NewFinance(yqy.ctx)
}

// GetTax 税种信息
func (yqy *YiQiYing) GetTax() *service2.Tax {
	return service2.NewTax(yqy.ctx)
}

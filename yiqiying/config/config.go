// Package config 亿企赢config配置
package config

import "github.com/yangzhenrui/finance/cache"

// Config .config for 亿企赢
type Config struct {
	AppKey    string `json:"appKey"`    // appid
	AppSecret string `json:"appSecret"` // appsecret
	Cache     cache.Cache
}

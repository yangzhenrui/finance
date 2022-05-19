package context

import (
	"github.com/yangzhenrui/finance/credential"
	"github.com/yangzhenrui/finance/yiqiying/config"
)

// Context struct
type Context struct {
	*config.Config
	credential.SignatureHandle
	Version     string `json:"version"`
	Timestamp   int64  `json:"timestamp"`
	ContentType string `json:"contentType"`
	XReqNonce   string `json:"XReqNonce"`
}

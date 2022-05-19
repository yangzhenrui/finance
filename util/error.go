package util

// CommonError 亿企赢返回的通用错误json
type CommonError struct {
	Code        string `json:"code"`
	Description string `json:"description"`
	Msg         string `json:"msg"`
	Time        string `json:"time"`
	Status      string `json:"status"`
}

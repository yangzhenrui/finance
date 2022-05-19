package credential

// SignatureHandle signature 接口
type SignatureHandle interface {
	GetSignature() (signature string, err error)
}

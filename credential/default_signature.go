package credential

import (
	"encoding/base64"
	"fmt"
	"github.com/chmike/hmacsha256"
	"github.com/gookit/goutil/strutil"
	"github.com/yangzhenrui/finance/cache"
	"strconv"
	"sync"
)

const (
	// CacheKeyYiQiYingPrefix 亿企赢cache key前缀
	CacheKeyYiQiYingPrefix = "go_yiqiying_"
)

// Signature 默认signature 获取
type Signature struct {
	appKey          string
	appSecret       string
	timestamp       int64
	version         string
	xReqNonce       string
	cacheKeyPrefix  string
	cache           cache.Cache
	accessTokenLock *sync.Mutex
}

// NewDefaultSignature new NewDefaultSignature
func NewDefaultSignature(appKey string, appSecret string, timestamp int64, version string, xReqNonce string, cacheKeyPrefix string, cache cache.Cache) SignatureHandle {
	if cache == nil {
		panic("cache is ineed")
	}
	return &Signature{
		appKey:          appKey,
		appSecret:       appSecret,
		timestamp:       timestamp,
		version:         version,
		xReqNonce:       xReqNonce,
		cache:           cache,
		cacheKeyPrefix:  cacheKeyPrefix,
		accessTokenLock: new(sync.Mutex),
	}
}

// GetSignature 获取signature,先从cache中获取，没有则重新生成
func (s *Signature) GetSignature() (signature string, err error) {
	// 先从cache中取
	signatureCacheKey := fmt.Sprintf("%s_signature_%s", s.cacheKeyPrefix, s.appKey)
	if val := s.cache.Get(signatureCacheKey); val != nil {
		return val.(string), nil
	}

	// 生成signature
	//signatureRet := util.Signature(s.appKey, s.appSecret, strconv.FormatInt(s.timestamp, 10), s.version, s.xReqNonce)
	mergeStr := fmt.Sprintf("%s%s%s%s%s", s.appKey, s.appSecret, strconv.FormatInt(s.timestamp, 10), s.version, s.xReqNonce)
	encodedStr := strutil.URLEncode(mergeStr)

	digest := hmacsha256.Digest(nil, strutil.ToBytes(s.appSecret), strutil.ToBytes(encodedStr))
	signatureRet := base64.StdEncoding.EncodeToString(digest)
	err = s.cache.Set(signatureCacheKey, signatureRet, 86400) // 1天过期
	if err != nil {
		return
	}
	signature = signatureRet
	return
}

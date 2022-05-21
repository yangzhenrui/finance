package credential

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"github.com/yangzhenrui/finance/cache"
	"net/url"
	"strconv"
)

const (
	// CacheKeyYiQiYingPrefix 亿企赢cache key前缀
	CacheKeyYiQiYingPrefix = "go_yiqiying_"
)

// Signature 默认signature 获取
type Signature struct {
	accountPeriod  *string
	appKey         string
	appSecret      string
	customerId     *string
	period         *string
	reclassifyFlag *string
	taxCode        *string
	timestamp      int64
	version        string
	xReqNonce      string
	cacheKeyPrefix string
	cache          cache.Cache
}

// NewDefaultSignature new NewDefaultSignature
func NewDefaultSignature(accountPeriod *string, appKey string, appSecret string, customerId *string, period *string, reclassifyFlag *string, taxCode *string, timestamp int64, version string, xReqNonce string, cacheKeyPrefix string, cache cache.Cache) SignatureHandle {
	if cache == nil {
		panic("cache is ineed")
	}

	var s Signature
	if accountPeriod != nil {
		s.accountPeriod = accountPeriod
	}
	s.appKey = appKey
	s.appSecret = appSecret
	if customerId != nil {
		s.customerId = customerId
	}
	if reclassifyFlag != nil {
		s.reclassifyFlag = reclassifyFlag
	}
	if period != nil {
		s.period = period
	}
	if taxCode != nil {
		s.taxCode = taxCode
	}

	s.timestamp = timestamp
	s.version = version
	s.xReqNonce = xReqNonce
	s.cache = cache
	s.cacheKeyPrefix = cacheKeyPrefix

	return &s
	//return &Signature{
	//	appKey:         appKey,
	//	appSecret:      appSecret,
	//	timestamp:      timestamp,
	//	version:        version,
	//	xReqNonce:      xReqNonce,
	//	cache:          cache,
	//	cacheKeyPrefix: cacheKeyPrefix,
	//}
}

// GetSignature 获取signature,先从cache中获取，没有则重新生成
func (s *Signature) GetSignature() (signature string, err error) {
	// 先从cache中取
	//signatureCacheKey := fmt.Sprintf("%s_signature_%s", s.cacheKeyPrefix, s.appKey)
	//if val := s.cache.Get(signatureCacheKey); val != nil {
	//	return val.(string), nil
	//}

	// 生成signature
	var accountPeriod, customerId, reclassifyFlag, period, taxCode string = "", "", "", "", ""
	if s.accountPeriod != nil {
		accountPeriod = *s.accountPeriod
	}
	if s.customerId != nil {
		customerId = *s.customerId
	}
	if s.reclassifyFlag != nil {
		reclassifyFlag = *s.reclassifyFlag
	}
	if s.taxCode != nil {
		taxCode = *s.taxCode
	}
	if s.period != nil {
		period = *s.period
	}
	mergeStr := fmt.Sprintf("%s%s%s%s%s%s%s%s%s%s", accountPeriod, s.appKey, s.appSecret, customerId, period, reclassifyFlag, taxCode, strconv.FormatInt(s.timestamp, 10), s.version, s.xReqNonce)
	encodedStr := url.QueryEscape(mergeStr)
	hash := hmac.New(sha256.New, []byte(s.appSecret))
	hash.Write([]byte(encodedStr))
	signData := hash.Sum(nil)
	sign := base64.StdEncoding.EncodeToString(signData)
	signature = sign

	return
}

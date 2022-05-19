package finance

import (
	log "github.com/sirupsen/logrus"
	"github.com/yangzhenrui/finance/cache"
	"github.com/yangzhenrui/finance/yiqiying"
	yiqiyingConfig "github.com/yangzhenrui/finance/yiqiying/config"
	"os"
)

func init() {
	// Log as JSON instead of the default ASCII formatter.
	log.SetFormatter(&log.TextFormatter{})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stdout)

	// Only log the warning severity or above.
	log.SetLevel(log.DebugLevel)
}

// Finance struct
type Finance struct {
	cache cache.Cache
}

// NewFinance init
func NewFinance() *Finance {
	return &Finance{}
}

// SetCache 设置cache
func (f *Finance) SetCache(cahce cache.Cache) {
	f.cache = cahce
}

// GetYiQiYing 获取亿企赢的实例
func (f *Finance) GetYiQiYing(cfg *yiqiyingConfig.Config) *yiqiying.YiQiYing {
	if cfg.Cache == nil {
		cfg.Cache = f.cache
	}
	return yiqiying.NewYiQiYing(cfg)
}

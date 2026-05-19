package gf

import (
	"errors"
	"time"

	"github.com/patrickmn/go-cache"
)

var goCahce *cache.Cache
func init() {
	goCahce = cache.New(5*time.Minute, 60*time.Second)
}

func SetGoCacheData(key string, data interface{}, expire time.Duration) {
	goCahce.Set(key, data, expire*time.Minute)
}

func GetGoCacheData(key string) (interface{}, error) {
	value, err := goCahce.Get(key)
	if !err {
		return nil, errors.New("获取缓存失败")
	}
	return value, nil
}

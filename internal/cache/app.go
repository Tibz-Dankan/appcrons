package cache

import (
	"context"
	"fmt"

	"github.com/Tibz-Dankan/keep-active/internal/config"
)

var redisClient = config.RedisClient()
var ctx = context.Background()

type AppCache struct {
	AppId         string `json:"appId"`
	DurationStart string `json:"durationStart"`
	DurationEnd   string `json:"durationEnd"`
	ReqInterval   string `json:"requestInterval"`
	LastReqTime   string `json:"lastReqTime,omitempty"`
	NextReqTime   string `json:"nextReqTime,omitempty"`
}

func (ac *AppCache) Write(appCache AppCache) error {

	app := map[string]string{
		"appId":         appCache.AppId,
		"durationStart": appCache.DurationStart,
		"durationEnd":   appCache.DurationEnd,
		"lastReqTime":   appCache.LastReqTime,
		"nextReqTime":   appCache.NextReqTime,
	}

	cacheId := "appid-" + ac.AppId
	for k, v := range app {
		err := redisClient.HSet(ctx, cacheId, k, v).Err()
		if err != nil {
			return err
		}
	}

	return nil
}

func (ac *AppCache) Read(appId string) map[string]string {
	cacheId := "appid-" + ac.AppId

	appCache := redisClient.HGetAll(ctx, cacheId).Val()
	fmt.Println(appCache)

	return appCache
}

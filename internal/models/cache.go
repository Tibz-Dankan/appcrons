package models

import (
	"context"

	"github.com/Tibz-Dankan/keep-active/internal/config"
)

var redisClient = config.RedisClient()
var ctx = context.Background()

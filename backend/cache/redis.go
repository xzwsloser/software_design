package cache

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/xzwsloser/software_design/backend/utils"
)

type RedisHanlder struct {
	Client 		*redis.Client
	Ctx    		context.Context
	CancelFunc  context.CancelFunc
}

var (
	redisHandler = new(RedisHanlder)
)

func InitRedisClient() {
	address := fmt.Sprintf("%s:%d", 
							utils.GetRedisConfig().Addr,
							utils.GetRedisConfig().Port)
    client := redis.NewClient(&redis.Options{
		Addr: address,
		Password: "",
		DB: int(utils.GetRedisConfig().DB),
	})

	redisHandler.Client = client
	redisHandler.Ctx, redisHandler.CancelFunc = context.WithCancel(context.Background())
}

func (r *RedisHanlder) Cancel() {
	r.CancelFunc()
}


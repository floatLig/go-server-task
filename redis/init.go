package redis

import (
	"context"
	"fmt"

	config "shopee.com/zeliang-entry-task/const"

	"github.com/go-redis/redis/v8"
)

//
//var Client *redis.Client = nil
//
//func InitRedis() error {
//	redisAddress := fmt.Sprintf("%s:%s", config.RedisHost, config.RedisPort)
//	Client = redis.NewClient(&redis.Options{
//		Addr:     redisAddress,
//		Password: config.RedisPassword,
//		DB:       config.RedisDb,
//		PoolSize: config.RedisPoolSize,
//	})
//	_, err := Client.Ping().Result()
//	if err != nil {
//		return err
//	}
//
//	go mysqlWriteToRedis()
//	return nil
//}

var (
	Ctx    = context.Background()
	Client *redis.Client
)

func InitRedis() error {
	redisAddress := fmt.Sprintf("%s:%s", config.RedisHost, config.RedisPort)
	rdb := redis.NewClient(&redis.Options{
		Addr:     redisAddress,
		Password: config.RedisPassword,
		DB:       config.RedisDb,
		PoolSize: config.RedisPoolSize,
	})

	_, err := rdb.Ping(Ctx).Result()
	if err != nil {
		fmt.Println("初始化 redis 失败: %s", err)
		return err
	}
	Client = rdb
	go mysqlWriteToRedis()

	return nil
}

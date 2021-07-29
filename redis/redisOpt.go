package redis

import (
	"encoding/json"
	"fmt"

	config "shopee.com/zeliang-entry-task/const"

	"shopee.com/zeliang-entry-task/model"
	"shopee.com/zeliang-entry-task/mysql"
)

func mysqlWriteToRedis() {
	fmt.Println("开始写redis")
	for userInfo := range mysql.UserChannel {
		err := SaveUser(&userInfo)
		if err != nil {
			fmt.Println("xx==> redis hset userinfo fail, err:", err.Error())
			continue
		}
	}
}

func SaveUserByToken(token string, userInfo *model.UserInfo) error {
	userKey := config.UserKeyPrefix + userInfo.Username
	Client.Set(Ctx, token, userKey, config.RedisExpire)

	result, _ := Client.Exists(Ctx, userKey).Result()
	if result == 0 {
		err := SaveUser(userInfo)
		return err
	}
	return nil
}

func SaveUser(info *model.UserInfo) error {
	data, err := json.Marshal(info)
	if err != nil {
		return err
	}
	userKey := config.UserKeyPrefix + info.Username
	Client.Set(Ctx, userKey, data, config.RedisExpire)
	return nil
}

func GetUserKeyByToken(token string) string {
	userKey := Client.Get(Ctx, token).Val()
	return userKey
}

func GetUser(userKey string) (*model.UserInfo, error) {
	value, err := Client.Get(Ctx, userKey).Bytes()
	if err != nil {
		return nil, err
	}

	var userInfo model.UserInfo
	err = json.Unmarshal(value, &userInfo)
	if err != nil {
		return nil, err
	}
	return &userInfo, nil
}

func Del(key string) {
	Client.Del(Ctx, key)
}

func UpdateExpire(key string) {
	Client.Expire(Ctx, key, config.RedisExpire)
}

func Exists(key string) int64 {
	return Client.Exists(Ctx, key).Val()
}

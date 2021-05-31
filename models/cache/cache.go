package cache

import (
	"encoding/json"
	"kumparantes/databases"
	"time"
)

func Set(key string, value interface{}) bool {
	js, _ := json.Marshal(value)
	str := string(js)
	return SetKey(key, str, 0)
}

func SetKey(key string, value interface{}, expired int64) bool {
	exp := time.Duration(expired) * time.Second
	client := databases.App.RedisConfig
	client.Set(key, value, exp)
	return true
}

func Clear(keys ...string) {
	client := databases.App.RedisConfig
	client.Del(keys...)
}

func GetAll(key string) (val []string) {
	// var cursor uint64
	client := databases.App.RedisConfig
	result := client.Keys(key).Val()
	for _, v := range result {
		data, _ := client.Get(v).Result()

		val = append(val, data)
	}
	return val
}

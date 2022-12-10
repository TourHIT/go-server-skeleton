package util

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"os"
	"time"
)

var (
	myRedis *redis.Client
	ctx     = context.Background()
)

func init() {

}

//初始化redis客户端
func InitRedis() {
	myRedis = redis.NewClient(&redis.Options{
		Addr: os.Getenv("REDIS_HOST"),
	})
	_, err := myRedis.Ping(context.Background()).Result()
	if err != nil {
		fmt.Println("Redis connect ping failed, err:", err)
		return
	}
	fmt.Println("Redis connect succeeded")
}

//关闭redis客户端
func Close() {
	myRedis.Close()
	myRedis = nil
	fmt.Println("Redis close")
}

func SetRedis(key string, value string, t int64) bool {
	if myRedis == nil {
		return false
	}
	expire := time.Duration(t) * time.Second
	if err := myRedis.Set(ctx, key, value, expire).Err(); err != nil {
		return false
	}
	return true
}

func GetRedis(key string) string {
	if myRedis == nil {
		return ""
	}
	result, err := myRedis.Get(ctx, key).Result()
	if err != nil {
		return ""
	}
	return result
}

func DelRedis(key string) bool {
	_, err := myRedis.Del(ctx, key).Result()
	if err != nil {
		fmt.Println(err)
		return false
	}
	return true
}

func ExpireRedis(key string, t int64) bool {
	// 延长过期时间
	expire := time.Duration(t) * time.Second
	if err := myRedis.Expire(ctx, key, expire).Err(); err != nil {
		fmt.Println(err)
		return false
	}
	return true
}

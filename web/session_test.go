package web

// import (
// 	"fmt"
// 	"git.lcgc.work/platform/kelp/redis"
// 	"testing"
// 	"time"
// )

// func TestRedisSession(t *testing.T) {
// 	redis.AddRedis("redis_session", "127.0.0.1:6379", 0)
// 	redis := redis.UseRedis("redis_session")
// 	duration := 60 * time.Second
// 	redisSessionPool := &RedisSessionPool{
// 		redis:    redis,
// 		duration: duration,
// 	}
// 	session := NewSession("token", duration)
// 	session.Set("test", "test")
// 	fmt.Println(session)
// 	redisSessionPool.Set("token", session)
// 	fmt.Println(redisSessionPool.Get("token"))
// 	redisSessionPool.Del("token")
// 	redisSessionPool.Get("token")
// }

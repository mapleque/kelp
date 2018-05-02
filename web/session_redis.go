package web

import (
	"encoding/json"
	"time"
)

type RedisSessionPool struct {
	redis    RedisClient
	duration time.Duration
}

type RedisClient interface {
	Get(string) (string, error)
	Set(key, value string, expiration time.Duration) error
	Del(key string) error
	Expire(key string, expiration time.Duration) error
}

func (this *Server) UseRedisSession(duration time.Duration, redis RedisClient) {
	redisSessionPool := &RedisSessionPool{
		redis:    redis,
		duration: duration,
	}
	this.UseSession(redisSessionPool)
}
func (this *RedisSessionPool) Get(token string) *Session {
	metaStr, err := this.redis.Get(token)
	if err != nil {
		return this.add(token)
	}
	meta := make(map[string]interface{})
	if err := json.Unmarshal([]byte(metaStr), &meta); err != nil {
		panic("json decode error: session{token:" + token + "};" + err.Error())
	}
	session := &Session{token: token, meta: meta, duration: this.duration}
	session.Refresh()
	return session
}

func (this *RedisSessionPool) Set(token string, session *Session) {
	metaBytes, err := json.Marshal(session.meta)
	if err != nil {
		panic(err)
	}
	if err := this.redis.Set(token, string(metaBytes), this.duration); err != nil {
		panic(err)
	}
}

func (this *RedisSessionPool) Del(token string) {
	if err := this.redis.Del(token); err != nil {
		panic(err)
	}
}

func (this *RedisSessionPool) add(token string) *Session {
	session := NewSession(token, this.duration)
	this.Set(token, session)
	return session
}

func (this *RedisSessionPool) Refresh(token string) {
	if err := this.redis.Expire(token, this.duration); err != nil {
		panic(err)
	}
}

package web

import (
	"errors"
	"time"
)

const (
	_SESSION_META_KEY        = "session"
	_SESSION_SERVER_META_KEY = "session_server"
)

type Session struct {
	token    string                 `json:"token"`
	meta     map[string]interface{} `json:"meta"`
	expired  time.Time
	duration time.Duration
}

type iSessionPool interface {
	Get(string) *Session
	Set(string, *Session)
	Del(string)
}

type _SessionServer struct {
	pool iSessionPool
}

func (this *Server) UseSession(sessionPool iSessionPool) {
	this.session = &_SessionServer{}
	this.session.setPool(sessionPool)
	this.Use(func(c *Context) {
		c.metaInternal[_SESSION_SERVER_META_KEY] = this.session
		c.Next()
		if session, ok := c.metaInternal[_SESSION_META_KEY]; ok {
			s := session.(*Session)
			this.session.pool.Set(s.token, s)
		}
	})
}

func (this *_SessionServer) setPool(pool iSessionPool) {
	this.pool = pool
}

// 支持自己实现的session方法
func (this *Context) NewSession(token string, duration time.Duration) {
	session := NewSession(token, duration)
	this.metaInternal[_SESSION_META_KEY] = session
}

func (this *Context) StartSession(token string) error {
	sessionServer, ok := this.metaInternal[_SESSION_SERVER_META_KEY]
	if !ok {
		return errors.New("this server dose not use any session server")
	}
	session := sessionServer.(*_SessionServer).pool.Get(token)
	this.metaInternal[_SESSION_META_KEY] = session
	return nil
}

func (this *Context) GetSession(key string) (interface{}, error) {
	session, ok := this.metaInternal[_SESSION_META_KEY]
	if !ok {
		return nil, errors.New("you should start session before get")
	}
	return session.(*Session).Get(key), nil
}

func (this *Context) SetSession(key string, value interface{}) error {
	session, ok := this.metaInternal[_SESSION_META_KEY]
	if !ok {
		return errors.New("you should start session before set")
	}
	session.(*Session).Set(key, value)
	return nil
}

func (this *Context) DestroySession(token string) error {
	sessionServer, ok := this.metaInternal[_SESSION_SERVER_META_KEY]
	if !ok {
		return errors.New("this server dose not use any session server")
	}
	sessionServer.(*_SessionServer).pool.Del(token)
	delete(this.metaInternal, _SESSION_META_KEY)
	return nil
}

func NewSession(token string, duration time.Duration) *Session {
	return &Session{
		token:    token,
		meta:     make(map[string]interface{}),
		expired:  time.Now().Add(duration),
		duration: duration,
	}
}

func (this *Session) Get(key string) interface{} {
	if value, ok := this.meta[key]; ok {
		return value
	}
	return nil
}

func (this *Session) Set(key string, value interface{}) {
	this.meta[key] = value
}

func (this *Session) Refresh() {
	this.expired = time.Now().Add(this.duration)
}

func (this *Session) IsExpired(t time.Time) bool {
	return t.After(this.expired)
}

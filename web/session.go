package web

import (
	"crypto/md5"
	"encoding/hex"
	"math/rand"
	"strconv"
	"time"
)

type SessionPool struct {
	pool map[string]*Session
}

type Session struct {
	key  string
	stub map[string]interface{}
	time time.Time
}

var sessionPool *SessionPool

const (
	SESSION_COOKIE_KEY = "KELP_SESSION"
	SESSION_EXPIRED    = 3 * time.Hour
)

func init() {
	sessionPool = &SessionPool{make(map[string]*Session)}
	go sessionGC()
}

func sessionGC() {
}

func SessionStart(context *Context) {
	session_key := context.GetCookie(SESSION_COOKIE_KEY)
	ok := true
	if len(session_key) > 0 {
		_, ok = sessionPool.pool[session_key]
	}
	if len(session_key) == 0 || !ok {
		session_key = generalSessionKey()
		context.SetCookie(SESSION_COOKIE_KEY, session_key, SESSION_EXPIRED)
		session := &Session{
			key:  session_key,
			stub: make(map[string]interface{}),
			time: time.Now(),
		}
		sessionPool.pool[session_key] = session
	} else {
		sessionPool.pool[session_key].time = time.Now()
	}
}

func sessionDestroy(context *Context) {
	session_key := context.GetCookie(SESSION_COOKIE_KEY)
	if len(session_key) == 0 {
		return
	}
	delete(sessionPool.pool, session_key)
}

func (context *Context) SessionDestroy() {
	sessionDestroy(context)
}

func (context *Context) SessionGet(key string) interface{} {
	session_key := context.GetCookie(SESSION_COOKIE_KEY)
	if len(session_key) == 0 {
		return nil
	}
	session, ok := sessionPool.pool[session_key]
	if !ok {
		return nil
	}
	return session.get(key)
}

func (context *Context) SessionSet(key string, value interface{}) {
	session_key := context.GetCookie(SESSION_COOKIE_KEY)
	if len(session_key) == 0 {
		log.Warn("empty session key may session not start yet")
		return
	}
	session, ok := sessionPool.pool[session_key]
	if !ok {
		log.Warn("can not find session in poll may session destroyed")
		return
	}
	session.set(key, value)
}

func (context *Context) SessionDel(key string) {
	session_key := context.GetCookie(SESSION_COOKIE_KEY)
	if len(session_key) == 0 {
		log.Warn("empty session key may session not start yet")
		return
	}
	session, ok := sessionPool.pool[session_key]
	if !ok {
		log.Warn("can not find session in poll may session destroyed")
		return
	}
	session.del(key)
}

func (session *Session) set(key string, value interface{}) {
	session.stub[key] = value
}

func (session *Session) get(key string) interface{} {
	return session.stub[key]
}

func (session *Session) del(key string) {
	delete(session.stub, key)
}

func generalSessionKey() string {
	seed := rand.Intn(10000)
	randStr := strconv.Itoa(seed)
	h := md5.New()
	h.Write([]byte(randStr))
	bs := h.Sum(nil)
	return hex.EncodeToString(bs)
}

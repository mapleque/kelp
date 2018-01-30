package web

import (
	"sync"
	"time"
)

type MemSessionPool struct {
	pool       map[string]*Session
	mux        *sync.RWMutex
	duration   time.Duration
	gcInterval time.Duration
}

func (this *Server) UseMemSession(duration time.Duration, gcInterval time.Duration) {
	memSessionPool := &MemSessionPool{
		pool:       make(map[string]*Session),
		mux:        new(sync.RWMutex),
		duration:   duration,
		gcInterval: gcInterval,
	}
	memSessionPool.startGC()
	this.UseSession(memSessionPool)
}

func (this *MemSessionPool) startGC() {
	ticker := time.NewTicker(this.gcInterval)
	go func() {
		for t := range ticker.C {
			for token, session := range this.pool {
				if session.IsExpired(t) {
					this.Del(token)
				}
			}
		}
	}()
}

func (this *MemSessionPool) Del(token string) {
	this.mux.Lock()
	defer this.mux.Unlock()
	delete(this.pool, token)
}

func (this *MemSessionPool) Get(token string) *Session {
	session, ok := this.pool[token]
	if ok {
		session.Refresh()
		return session
	}
	return this.add(token)
}

func (this *MemSessionPool) Set(token string, session *Session) {
	this.mux.Lock()
	defer this.mux.Unlock()
	this.pool[token] = session
}

func (this *MemSessionPool) add(token string) *Session {
	session := NewSession(token, this.duration)
	this.Set(token, session)
	return session
}

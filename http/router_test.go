package http

import (
	"testing"
)

func TestFind(t *testing.T) {
	s := New("")
	r0 := s.Group("/0")
	r0.Handle("", "/1", func() {})
	r0.Handle("", "/2", func() {}, func() {})
	r1 := s.Group("/1").Use(func() {})
	r1.Handle("", "/1", func() {})
	r1.Handle("", "/2", func() {}, func() {})

	assertFindNum(t, s.router.find("/1/1"), 2)
	assertFindNum(t, s.router.find("/1/2"), 3)
	assertFindNum(t, s.router.find("/0/1"), 1)
	assertFindNum(t, s.router.find("/0/2"), 2)
}

func assertFindNum(t *testing.T, r *Router, num int) {
	if len(r.handlerChain) != num {
		t.Error("handler num should be", num, "but", len(r.handlerChain))
	}
}

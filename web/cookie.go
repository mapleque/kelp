package web

import (
	"net/http"
	"time"
)

func (this *Context) GetCookie(key string) (string, error) {
	cookie, err := this.Request.Cookie(key)
	if err != nil {
		return "", err
	}
	return cookie.Value, nil
}

func (this *Context) SetCookie(key, value string, duration time.Duration) {
	now := time.Now()
	cookie := &http.Cookie{
		Name:    key,
		Value:   value,
		Path:    "/",
		Expires: now.Add(duration),
	}
	this.Request.AddCookie(cookie)
	http.SetCookie(this.ResponseWriter, cookie)
}

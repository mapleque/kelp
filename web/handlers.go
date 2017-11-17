package web

import (
	"time"
)

func RecoveryHandler(c *Context) {
	defer func() {
		if err := recover(); err != nil {
			log.Error("[panic]", err)
			c.DieWithHttpStatus(500)
		}
	}()
	c.Next()
}

func LogHandler(c *Context) {
	start := time.Now()
	path := c.Request.URL.Path
	raw := c.Request.URL.RawQuery

	c.Next()

	end := time.Now()
	latency := end.Sub(start)
	method := c.Request.Method
	status := c.Status
	httpStatus := c.HttpStatus
	resp := string(c.Response)

	if raw != "" {
		path = path + "?" + raw
	}

	log.Info(
		end.Format("2006/01/02 15:04:05"),
		latency,
		method,
		path,
		status,
		httpStatus,
		resp,
	)
}

func SessionWithCookieHandler(cookieSessionKey string, duration time.Duration) HandlerFunc {
	return func(c *Context) {
		token, err := c.GetCookie(cookieSessionKey)
		if err != nil || len(token) < 1 {
			token = RandMd5()
		}
		if err := c.StartSession(token); err != nil {
			c.Error(-1, err)
			return
		}
		c.SetCookie(cookieSessionKey, token, duration)
		c.Next()
	}
}

func SessionWithCookieDestoryHandler(cookieSessionKey string) HandlerFunc {
	return func(c *Context) {
		token, err := c.GetCookie(cookieSessionKey)
		if err == nil && len(token) > 0 {
			c.DestroySession(token)
		}
		c.SetCookie(cookieSessionKey, token, -1*time.Second)
		c.Next()
	}
}

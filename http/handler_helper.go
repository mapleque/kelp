package http

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// TraceHandler
// If you want start trace, use this on your root router
func TraceHandler(c *Context) {
	c.Request.Header.Set("Kelp-Traceid", RandMd5())
	c.Next()
}

func RecoveryHandler(c *Context) {
	defer func() {
		if err := recover(); err != nil {
			log.Error("[panic]", err)
			c.DieWithHttpStatus(500)
		}
	}()
	c.Next()
}

func str(v string) string {
	if v == "" {
		return "-"
	}
	return v
}

func LogHandler(c *Context) {
	start := time.Now()
	path := c.Request.URL.Path
	raw := c.Request.URL.RawQuery
	ips := c.Request.Header.Get("X-Forwarded-For")
	ip := ""
	if ips != "" {
		ip = strings.Split(ips, ",")[0]
	}
	if ip == "" {
		ip = c.Request.Header.Get("X-Real-Ip")
	}
	if ip == "" {
		ip = c.Request.RemoteAddr
	}

	c.Next()

	traceId := c.Request.Header.Get("Kelp-Traceid")
	uuid := c.Request.Header.Get("uuid")
	end := time.Now()
	latency := end.Sub(start)
	method := c.Request.Method
	resp := string(c.Response)
	if len(resp) > 500 {
		resp = fmt.Sprintf("response is too large (with %d bytes, head is %s)", len(resp), resp[0:100]+"...")
	}
	req := string(c.Body())
	if raw != "" {
		path = path + "?" + raw
	}

	log.Log(
		"REQ",
		ip, // remote ip
		end.Format("2006/01/02 15:04:05"),
		latency.Nanoseconds()/int64(time.Millisecond),
		str(method),
		str(path),
		str(traceId), // trace id
		str(uuid),    // uuid
		`"""`+str(req)+`"""`,
		`"""`+str(resp)+`"""`,
	)
}

func TokenAuthorization(token string) HandlerFunc {
	return func(c *Context) {
		auth := c.Request.Header.Get("Authorization")
		if auth != token {
			log.Error("[authorization failed]", auth)
			c.DieWithHttpStatus(401)
		} else {
			c.Next()
		}
	}
}

func SignCheck(sign string) HandlerFunc {
	return func(c *Context) {
		token := c.QueryDefault("token", "")
		timestamp, err := strconv.ParseInt(c.QueryDefault("timestamp", "0"), 10, 64)
		if err != nil {
			c.DieWithHttpStatus(401)
			return
		}
		if timestamp == 0 {
			if !Sha1Verify([]byte(sign), c.Body(), []byte(token), 5) {
				c.DieWithHttpStatus(401)
				return
			}
		} else {
			if !Sha1VerifyTimestamp([]byte(sign), c.Body(), []byte(token), 5, timestamp) {
				c.DieWithHttpStatus(401)
				return
			}
		}
		c.Next()
	}
}

package http

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

// TraceHandler
// If you want start trace, use this on your root router
func TraceHandler(c *Context) {
	c.Request.Header.Set("Kelp-Traceid", randMd5())
	c.Next()
}

func RecoveryHandler(c *Context) {
	defer func() {
		if err := recover(); err != nil {
			log.Log("ERROR", "[panic]", err)
			c.DieWithHttpStatus(500)
		}
	}()
	c.Next()
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

func str(v string) string {
	if v == "" {
		return "-"
	}
	return v
}

func randMd5() string {
	timestamp := []byte(strconv.FormatInt(time.Now().Unix(), 10))
	prefix := []byte(strconv.Itoa(rand.Intn(10000)))
	surfix := []byte(strconv.Itoa(rand.Intn(10000)))
	seed := bytes.Join([][]byte{prefix, timestamp, surfix}, []byte(""))

	h := md5.New()
	h.Write(seed)
	data := h.Sum(nil)
	dst := make([]byte, hex.EncodedLen(len(data)))
	hex.Encode(dst, data)
	return string(dst)
}

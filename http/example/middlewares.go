package example

import (
	"../../http"
)

func (this *Server) Auth(c *http.Context) {
	// 获取basic auth凭证
	username, password, ok := c.Request.BasicAuth()
	if ok {
		// check username and password
		if username == "kelp" && password == "kelp" {
			c.Next()
		}
	}
	// 直接返回无权限
	c.DieWithHttpStatus(401)
}

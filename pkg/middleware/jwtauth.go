package middleware

import (
	"github.com/bryant-rh/cm/pkg/jwt"
	"github.com/bryant-rh/cm/pkg/util"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

//基于JWT的认证中间件
func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 客户端携带Token有三种方式 1.放在请求头 2.放在请求体 3.放在URI
		// 这里假设Token放在URI的Authorization中，并使用Bearer开头

		authHeader := c.Request.Header.Get("Authorization")
		msg := ""
		//auth := c.Query("auth")
		//authHeader := auth
		if authHeader == "" {
			msg = "请求头中auth为空"
			util.ReturnMsg(c, http.StatusUnauthorized, "", msg)
			c.Abort()
			return
		}
		// 按空格分割
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			msg = "请求头中auth格式有误"
			util.ReturnMsg(c, http.StatusUnauthorized, "", msg)
			c.Abort()
			return
		}
		// parts[1]是获取到的tokenString，我们使用之前定义好的解析JWT的函数来解析它
		mc, err := jwt.ParseToken(parts[1])
		//mc, err := jwt.ParseToken(authHeader)
		if err != nil {
			msg = fmt.Sprintf("%s", err)
			util.ReturnMsg(c, http.StatusUnauthorized, "", msg)
			c.Abort()
			return
		}

		// 将当前请求的username信息保存到请求的上下文c上
		c.Set("username", mc.Username)
		c.Next()
		// 后续的处理函数可以用过c.Get("username")来获取当前请求的用户信息
	}
}

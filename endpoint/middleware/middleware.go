/**
* @Auther:gy
* @Date:2020/11/23 16:20
 */

package middleware

import (
	"bytes"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"time"
)

// CORSMiddleware 跨域请求中间件
func CORSMiddleware() gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "PATCH"},
		AllowHeaders:     []string{"Accept, Accept-Encoding, Authorization, Content-Type, DNT, Origin, User-Agent, X-CSFRTOKEN, X-Requested-With"},
		AllowCredentials: true,
		MaxAge:           time.Second * time.Duration(7200),
	})
}
func RequestParam() gin.HandlerFunc {
	return func(c *gin.Context) {
		data, _ := c.GetRawData()
		c.Set("reqBody", data)
		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(data))
	}
}

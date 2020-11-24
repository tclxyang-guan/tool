/**
* @Auther:gy
* @Date:2020/11/23 16:23
 */
package endpoint

import (
	"github.com/gin-gonic/gin"
	"transfDoc/conf"
	"transfDoc/endpoint/api/build"
	"transfDoc/endpoint/middleware"
)

func InitRouter() *gin.Engine {
	r := gin.New()

	//跨域访问
	r.Use(middleware.CORSMiddleware())
	//请求是否生成在线showDoc文档
	if conf.GetConfig().ShowDocOpen {
		r.Use(middleware.RequestParam())
	}
	smartyBasic := r.Group("/test/v1")
	smartyBasic.POST("/build", build.BuildCreate)      //创建
	smartyBasic.GET("/builds", build.BuildPage)        //分页查询
	smartyBasic.PUT("/build/:id", build.BuildUpdate)   //整体修改
	smartyBasic.PATCH("/build/:id", build.BuildUpdate) //分页查询
	return r
}

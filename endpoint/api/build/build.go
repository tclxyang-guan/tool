/**
* @Auther:gy
* @Date:2020/11/23 16:23
 */

package build

import (
	"github.com/gin-gonic/gin"
	"transfDoc/endpoint/common"
	"transfDoc/models"
	"transfDoc/pkg/logging"
)

//添加楼栋
func BuildCreate(c *gin.Context) {
	req := models.Build{} //为了方便直接使用的models struct
	rsp := common.Rsp{}
	c.BindJSON(&req)
	logging.Infof("build create req:%+v", req)

	rsp.ReplySuccess(c, "")
	return
}

//楼栋修改
func BuildUpdate(c *gin.Context) {
	req := map[string]interface{}{}
	rsp := common.Rsp{}
	c.BindJSON(&req)
	logging.Infof("BuildUpdate req:%+v", req)

	rsp.ReplySuccess(c, "")
	return
}

// 查询楼栋信息分页
func BuildPage(c *gin.Context) {
	rsp := common.Rsp{}
	build1 := models.Build{}
	build1.ID = 1
	build1.BuildName = "a"
	build2 := models.Build{}
	build2.ID = 2
	build2.BuildName = "b"
	builds := []models.Build{build1, build2}
	rsp.ReplySuccess(c, map[string]interface{}{
		"list":  builds,
		"count": 10,
	})
	return
}

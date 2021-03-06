/**
* @Auther:gy
* @Date:2020/11/23 16:23
 */
package tool

import (
	"github.com/gin-gonic/gin"
	"github.com/guanyang-lab/tool/utils/e"
	"net/http"
)

type ListRsp struct {
	List  interface{} `json:"list" comment:"列表数据"`
	Count int         `json:"count" comment:"总条数"`
}

type Rsp struct {
	Code int         `json:"code" comment:"状态码 0成功"`
	Msg  string      `json:"msg" comment:"失败信息"`
	Data interface{} `json:"data,omitempty" comment:"返回数据"`
}

func (r Rsp) GetMessage() string {
	return e.GetMsg(r.Code)
}

// 回复参数错误消息
func (r Rsp) ReplyInvalidParam(c *gin.Context) {
	r.Code = e.INVALID_PARAMS
	r.Msg = e.GetMsg(r.Code)
	ReplyJson(r, c)
}

// 回复成功， data为返回值
func (r Rsp) ReplySuccess(c *gin.Context, data interface{}) {
	r.Code = e.SUCCESS
	r.Msg = e.GetMsg(r.Code)
	r.Data = data
	ReplyJson(r, c)
}

// 回复操作失败给前端，msg为失败原因
func (r Rsp) ReplyFailOperation(c *gin.Context, msg string) {
	r.Code = e.FAILE_TO_CREATE_OP
	r.Msg = msg
	ReplyJson(r, c)
}

// 回复自定义code 自定义消息
func (r Rsp) Reply(c *gin.Context, code int, msg string) {
	r.Code = code
	r.Msg = msg
	ReplyJson(r, c)
}
func ReplyJson(r Rsp, c *gin.Context) {
	c.Set("rsp", r)
	go GetClient().Upload(c)
	c.JSON(http.StatusOK, r)
}

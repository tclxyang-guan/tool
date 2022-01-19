/**
* @Auther:gy
* @Date:2020/11/23 16:23
 */
package example

import "gitee.com/yanggit123/tool"

const (
	SUCCESS            = 0
	INVALID_PARAMS     = 501
	FAILE_TO_CREATE_OP = 502
)

var MsgFlags = map[int]string{
	SUCCESS:            "ok",
	INVALID_PARAMS:     "请求参数错误",
	FAILE_TO_CREATE_OP: "操作失败:",
}

// GetMsg get error information based on Code
func GetMsg(code int) string {
	msg, ok := MsgFlags[code]
	if ok {
		return msg
	}
	return ""
}

//私人返回前端响应体例子
type Rsp struct {
	Code int         `json:"code" comment:"状态码 0成功"`
	Msg  string      `json:"msg" comment:"失败信息"`
	Data interface{} `json:"data,omitempty" comment:"返回数据"`
}

func (r Rsp) GetMessage() string {
	return GetMsg(r.Code)
}

// 回复参数错误消息
func (r Rsp) ReplyInvalidParam(c interface{}) {
	r.Code = INVALID_PARAMS
	r.Msg = GetMsg(r.Code)
	tool.ReplyJson(c, r)
}

// 回复成功， data为返回值
func (r Rsp) ReplySuccess(c interface{}, data interface{}) {
	r.Code = SUCCESS
	r.Msg = GetMsg(r.Code)
	r.Data = data
	tool.ReplyJson(c, r)
}

// 回复成功， 自定义code
func (r Rsp) ReplySuccessCode(c interface{}, code int, data interface{}) {
	r.Code = code
	r.Data = data
	tool.ReplyJson(c, r)
}

// 回复操作失败给前端，msg为失败原因
func (r Rsp) ReplyFailOperation(c interface{}, msg string) {
	r.Code = FAILE_TO_CREATE_OP
	r.Msg = msg
	tool.ReplyJson(c, r)
}

// 回复自定义code 自定义消息
func (r Rsp) Reply(c interface{}, code int, msg string) {
	r.Code = code
	r.Msg = msg
	tool.ReplyJson(c, r)
}

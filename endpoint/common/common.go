/**
* @Auther:gy
* @Date:2020/11/23 16:23
 */
package common

import (
	"github.com/gin-gonic/gin"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"net/http"
	"reflect"
	"transfDoc/utils/e"

	zhongwen "github.com/go-playground/locales/zh"
	zh_translations "github.com/go-playground/validator/v10/translations/zh"
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

// 回复token验证失败消息
func (r Rsp) ReplyTokenError(c *gin.Context) {
	r.Code = e.ERROR_AUTH_CHECK_TOKEN_FAIL
	r.Msg = e.GetMsg(r.Code)
	ReplyJson(r, c)
}
func ReplyJson(r Rsp, c *gin.Context) {
	c.Set("rsp", r)
	go DeferShowDoc(c)
	c.JSON(http.StatusOK, r)
}

var Validate *validator.Validate
var Trans ut.Translator

func init() {
	zh := zhongwen.New()
	uni := ut.New(zh, zh)
	Trans, _ = uni.GetTranslator("zh")

	Validate = validator.New()
	Validate.RegisterTagNameFunc(func(field reflect.StructField) string {
		label := field.Tag.Get("label")
		if label == "" {
			return field.Name
		}
		return label
	})
	zh_translations.RegisterDefaultTranslations(Validate, Trans)
}

func TransError(err error) string {
	s := err.(validator.ValidationErrors).Translate(Trans)
	for _, value := range s {
		return value
	}
	return "参数错误"
}

/**
* @Auther:gy
* @Date:2020/11/23 16:23
 */
package common

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/spf13/cast"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"transfDoc/conf"
	"transfDoc/pkg/logging"
	"transfDoc/utils"
	"transfDoc/utils/e"

	zhongwen "github.com/go-playground/locales/zh"
	zh_translations "github.com/go-playground/validator/v10/translations/zh"
)

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
	go r.DeferShowDoc(c)
	c.JSON(http.StatusOK, r)
}

// 回复成功， data为返回值
func (r Rsp) ReplySuccess(c *gin.Context, data interface{}) {
	r.Code = e.SUCCESS
	r.Msg = e.GetMsg(r.Code)
	r.Data = data
	go r.DeferShowDoc(c)
	c.JSON(http.StatusOK, r)
}

// 回复操作失败给前端，msg为失败原因
func (r Rsp) ReplyFailOperation(c *gin.Context, msg string) {
	r.Code = e.FAILE_TO_CREATE_OP
	r.Msg = msg
	go r.DeferShowDoc(c)
	c.JSON(http.StatusOK, r)
}

// 打印接口请求参数，path为接口路径，p为参数
func LogParams(path string, p interface{}) {
	logging.Infof("Request of %s : %v", path, p)
}

// 回复token验证失败消息
func (r Rsp) ReplyTokenError(c *gin.Context) {
	r.Code = e.ERROR_AUTH_CHECK_TOKEN_FAIL
	r.Msg = e.GetMsg(r.Code)
	go r.DeferShowDoc(c)
	c.JSON(http.StatusOK, r)
}

func (r Rsp) DeferShowDoc(c *gin.Context) {
	if !conf.GetConfig().ShowDocOpen {
		return
	}
	key := c.Request.Header.Get("dockey")
	if key == "" {
		return
	}
	method := c.Request.Method
	data := ShowDocMap[key]
	usdp := UploadShowDocParam{}
	usdp.ApiKey = conf.GetConfig().ApiKey
	usdp.ApiToken = conf.GetConfig().ApiToken
	usdp.CatName = data.CatName
	usdp.PageTitle = data.PageTitle
	usdp.SNumber = data.SNumber
	reqbytes, _ := c.Get("reqBody")
	reqParam := cast.ToString(reqbytes)
	if strings.ToUpper(method) == "GET" {
		reqParam = "同rul"
	}
	respBody, _ := json.Marshal(r)
	respParam := string(respBody)
	respParam = strings.ReplaceAll(respParam, "{", "{\r\n\t")
	respParam = strings.ReplaceAll(respParam, "}", "\t\r\n}\r\n")
	respParam = strings.ReplaceAll(respParam, ",", ",\r\n\t")
	resp := DatamapGenerateResp(data.Resp)
	if data.Resp != nil && data.Resp.Data != nil {
		resp += "\r\n"
		resp += DatamapGenerateResp(data.Resp.Data)
	}

	usdp.PageContent = "" +
		"##### 简要描述\r\n" +
		"\r\n" +
		"- " + usdp.PageTitle + "\r\n" +
		"\r\n" +
		"##### 请求URL\r\n" +
		"- ` " + c.Request.URL.String() + " `\r\n" +
		"\r\n" +
		"##### 请求方式\r\n" +
		"- " + method + "\r\n" +
		"\r\n" +
		"##### 参数\r\n" +
		"\r\n" + DatamapGenerateReq(data.Req) +
		"\r\n" +
		"\r\n" +
		"##### 请求示例\r\n" +
		"\r\n" +
		"```\r\n" +
		reqParam + "\r\n" +
		"```\r\n" +
		"##### 返回示例\r\n" +
		"\r\n" +
		"```\r\n" +
		stringTransf(string(respBody)) + "\r\n" +
		"```\r\n" +
		"\r\n" +
		"##### 返回参数说明\r\n" +
		"\r\n" +
		resp + "\r\n" +
		"\r\n" +
		"##### 备注\r\n" +
		"\r\n"
	//调用showdoc接口
	b, err := json.Marshal(usdp)
	bs, err := utils.NewRequest("post", conf.GetConfig().ShowDocUrl, bytes.NewReader(b), "")
	if err != nil {
		//return
	}
	var m map[string]interface{}
	json.Unmarshal(bs, &m)
	if cast.ToInt(m["error_code"]) != 0 {
		fmt.Println(m["error_message"])
	}
}

type Session struct {
	UserId        string // 用户id
	UserName      string // 用户名
	CommunityId   string // 社区id
	CommunityName string // 社区名
	NeighborId    string // 小区id
	NeighborName  string // 小区名
	UserType      int    // 用户类型  1，社区管理员 2，社区住户
}

// 获取用户基本信息
func GetSession(c *gin.Context) *Session {
	s := &Session{}
	s.UserId = c.GetHeader("uid")
	s.UserName = c.GetHeader("uname")
	s.CommunityId = c.GetString("cid")
	s.CommunityName = c.GetHeader("cname")
	s.NeighborId = c.GetHeader("nid")
	s.NeighborName = c.GetHeader("nname")
	s.UserType, _ = strconv.Atoi(c.GetHeader("utype"))
	return s
}

// 获取用户id
func GetUserID(c *gin.Context) string {
	return c.GetHeader("uid")
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

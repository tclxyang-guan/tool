/**
* @Auther:gy
* @Date:2021/3/9 20:01
 */

package example

import (
	"fmt"
	"gitee.com/yanggit123/tool"
	"github.com/gin-gonic/gin"
	"github.com/kataras/iris/v12"
	"strconv"
	"testing"
)

//请求绑定必须使用tool包中的BindJSON或者BindJSONNotWithValidate
//若不想使用rsp.go 请使用ReplyJson或者ReplyCodeJson
func TestGinResp(t *testing.T) {
	r := gin.New()
	tool.EnableDoc(1, tool.Doc{
		Url:      "https://www.showdoc.cc/server/api/item/updateByApi",
		ApiKey:   "06f5a93cbfe1c481d0fc20cf7d6692b92034679355",
		ApiToken: "061f631ecd83adaef9d91ed3616f7db71466087183",
	}, tool.Gin{})
	/*EnableDoc(r, 2, Doc{
		DocDir: "./test",
	})*/
	r.POST("/user/create", userCreate)
	r.GET("/user/detail", userDetail)
	r.POST("/user/update", userUpdate)
	r.POST("/user/delete", userDelete)
	r.POST("/user/page", userPage)

	r.POST("/role/create", roleCreate)
	r.Run("0.0.0.0:8080")
}
func TestIrisResp(t *testing.T) {
	app := iris.New()
	tool.EnableDoc(1, tool.Doc{
		Url:      "https://www.showdoc.cc/server/api/item/updateByApi",
		ApiKey:   "06f5a93cbfe1c481d0fc20cf7d6692b92034679355",
		ApiToken: "061f631ecd83adaef9d91ed3616f7db71466087183",
	}, tool.Iris12{})
	app.Post("/user/create", userIrisCreate)
	app.Get("/user/detail", userIrisDetail)
	app.Post("/user/update", userIrisUpdate)
	app.Post("/user/delete", userIrisDelete)
	app.Post("/user/page", userIrisPage)
	app.Run(iris.Addr("0.0.0.0:8080"))
}

type docModel struct {
	ID        uint    `gorm:"primary_key;comment:'数据id'" json:"id" req:"-"`
	CreatedAt string  `json:"created_at" gorm:"type:varchar(30);comment:创建时间" req:"-"`
	UpdatedAt string  `json:"updated_at" gorm:"type:varchar(30);comment:修改时间" req:"-"`
	DeletedAt *string `gorm:"type:varchar(30);default:null;comment:'删除时间'" json:"deleted_at" resp:"-" req:"-"`
}
type u struct {
	docModel
	Name string `json:"name" comment:"用户名称" validate:"required"`
}

/*
使用公共响应体
*/
func userCreate(c *gin.Context) {
	var (
		rsp Rsp
		req u
	)
	err := tool.BindJSON(c, &req)
	if err != nil {
		rsp.ReplyFailOperation(c, err.Error())
		return
	}
	//获取参数成功进行service、db操作
	//操作完成返回
	rsp.ReplySuccess(c, "用户创建成功")
}

//因为修改也需要进行参数校验所以接收还是采用struct
func userUpdate(c *gin.Context) {
	var (
		rsp Rsp
		req u
	)
	err := tool.BindJSON(c, &req)
	if err != nil {
		rsp.ReplyFailOperation(c, err.Error())
		return
	}
	//获取参数成功进行service、db操作
	//操作完成返回
	rsp.ReplySuccess(c, "用户修改成功")
}

func userDetail(c *gin.Context) {
	var (
		rsp Rsp
	)
	id, _ := strconv.Atoi(c.Query("id"))
	//获取参数成功进行service、db操作
	//操作完成返回
	u := u{}
	u.ID = uint(id)
	u.Name = "123"
	rsp.ReplySuccess(c, u)
}

type deteleReq struct {
	ID uint `json:"id" comment:"数据id" validate:"required,gt=0"`
}

func userDelete(c *gin.Context) {
	var (
		rsp Rsp
		req deteleReq
	)
	err := tool.BindJSON(c, &req)
	if err != nil {
		rsp.ReplyFailOperation(c, err.Error())
		return
	}
	//获取参数成功进行service、db操作
	//操作完成返回
	rsp.ReplySuccess(c, "删除成功")
}

type pageReq struct {
	Page     int `json:"page" comment:"当前页大小" validate:"required,gt=0"`
	PageSize int `json:"page_size" comment:"每页大小" validate:"required"`
}

func userPage(c *gin.Context) {
	var (
		rsp Rsp
		req pageReq
	)
	err := tool.BindJSON(c, &req)
	if err != nil {
		rsp.ReplyFailOperation(c, err.Error())
		return
	}
	//获取参数成功进行service、db操作
	data, count, err := []u{{Name: "zhang"}, {Name: "sss"}}, 2, nil
	//操作完成返回
	//分页返回的数据需使用工具包提供的ListRsp  或者返回数据中不能包含interface  interface无法反射
	rsp.ReplySuccess(c, map[string]interface{}{
		"list":  data,
		"count": count,
	})
}

/*
使用自定义响应体
*/
type myRsp struct {
	Code string      `json:"code" comment:"状态码"`
	Msg  interface{} `json:"msg" comment:"错误信息"`
	Data interface{} `json:"data" comment:"响应数据"`
}
type role struct {
	Name string `json:"name" comment:"名称" validate:"required"`
	Desc string `json:"desc" comment:"描述"`
}

func roleCreate(c *gin.Context) {
	var (
		rsp myRsp
		req role
	)
	err := tool.BindJSON(c, &req)
	if err != nil {
		rsp.Code = "2001"
		rsp.Msg = err.Error()
		tool.ReplyJson(c, rsp)
		return
	}
	rsp.Code = "0"
	rsp.Msg = "角色创建成功"
	tool.ReplyJson(c, rsp)
}
func userIrisCreate(c iris.Context) {
	var (
		rsp Rsp
		req u
	)
	err := tool.BindJSON(c, &req)
	if err != nil {
		rsp.ReplyFailOperation(c, err.Error())
		fmt.Println(err.Error())
		return
	}
	fmt.Println(c.Values().Get("req"))
	//获取参数成功进行service、db操作
	//操作完成返回
	rsp.ReplySuccess(c, "用户创建成功")
}

//因为修改也需要进行参数校验所以接收还是采用struct
func userIrisUpdate(c iris.Context) {
	var (
		rsp Rsp
		req u
	)
	err := tool.BindJSON(c, &req)
	if err != nil {
		rsp.ReplyFailOperation(c, err.Error())
		return
	}
	//获取参数成功进行service、db操作
	//操作完成返回
	rsp.ReplySuccess(c, "用户修改成功")
}

func userIrisDetail(c iris.Context) {
	var (
		rsp Rsp
	)
	id, _ := strconv.Atoi(c.URLParam("id"))
	//获取参数成功进行service、db操作
	//操作完成返回
	u := u{}
	u.ID = uint(id)
	u.Name = "123"
	rsp.ReplySuccess(c, u)
}

func userIrisDelete(c iris.Context) {
	var (
		rsp Rsp
		req deteleReq
	)
	err := tool.BindJSON(c, &req)
	if err != nil {
		rsp.ReplyFailOperation(c, err.Error())
		return
	}
	//获取参数成功进行service、db操作
	//操作完成返回
	rsp.ReplySuccess(c, "删除成功")
}

func userIrisPage(c iris.Context) {
	var (
		rsp Rsp
		req pageReq
	)
	err := tool.BindJSON(c, &req)
	if err != nil {
		rsp.ReplyFailOperation(c, err.Error())
		return
	}
	//获取参数成功进行service、db操作
	data, count, err := []u{{Name: "zhang"}, {Name: "sss"}}, 2, nil
	//操作完成返回
	//分页返回的数据需使用工具包提供的ListRsp  或者返回数据中不能包含interface  interface无法反射
	rsp.ReplySuccess(c, map[string]interface{}{
		"list": data, "count": count})
}

package tool

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/url"
	"strconv"
)

type Gin struct {
}

func (g Gin) BindJSON(c, req interface{}) error {
	ctx, ok := c.(*gin.Context)
	if !ok {
		return errors.New("类型错误")
	}
	return ctx.BindJSON(req)
}
func (g Gin) Set(c interface{}, key string, value interface{}) {
	ctx, ok := c.(*gin.Context)
	if !ok {
		return
	}
	ctx.Set(key, value)
}
func (g Gin) JSON(c, data interface{}) {
	ctx, ok := c.(*gin.Context)
	if !ok {
		return
	}
	ctx.JSON(http.StatusOK, data)
}
func (g Gin) GetDocParam(ctx interface{}) *docParam {
	c, ok := ctx.(*gin.Context)
	if !ok {
		return nil
	}
	pageTitle := c.Request.Header.Get("title")
	if pageTitle == "" {
		fmt.Println("not found title in header")
		return nil
	}
	pageTitle, _ = url.PathUnescape(pageTitle)
	//接口没有文件夹和顺序也没关系
	catName := c.Request.Header.Get("dir")
	if pageTitle != "" {
		catName, _ = url.PathUnescape(catName)
	}
	sNumber, _ := strconv.Atoi(c.Request.Header.Get("number"))
	req, _ := c.Get("req")
	resp, _ := c.Get("rsp")
	method := c.Request.Method
	return &docParam{
		ApiKey:      cli.Doc.ApiKey,
		ApiToken:    cli.Doc.ApiToken,
		CatName:     catName,
		PageTitle:   pageTitle,
		PageContent: "",
		SNumber:     sNumber,
		Url:         c.Request.URL.String(),
		Req:         req,
		Resp:        resp,
		Method:      method,
	}
}

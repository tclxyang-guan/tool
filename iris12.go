package tool

import (
	"errors"
	"fmt"
	"github.com/kataras/iris/v12"
	"net/url"
	"strconv"
)

type Iris12 struct {
}

func (g Iris12) BindJSON(c, req interface{}) error {
	ctx, ok := c.(iris.Context)
	if !ok {
		return errors.New("类型错误")
	}
	return ctx.ReadJSON(req)
}
func (g Iris12) Set(c interface{}, key string, value interface{}) {
	ctx, ok := c.(iris.Context)
	if !ok {
		return
	}
	ctx.Values().Set(key, value)
}
func (g Iris12) JSON(c, data interface{}) {
	ctx, ok := c.(iris.Context)
	if !ok {
		return
	}
	ctx.JSON(data)
}
func (g Iris12) GetDocParam(ctx interface{}) *docParam {
	c, ok := ctx.(iris.Context)
	if !ok {
		return nil
	}
	pageTitle := c.Request().Header.Get("title")
	if pageTitle == "" {
		fmt.Println("not found title in header")
		return nil
	}
	pageTitle, _ = url.PathUnescape(pageTitle)
	//接口没有文件夹和顺序也没关系
	catName := c.Request().Header.Get("dir")
	if pageTitle != "" {
		catName, _ = url.PathUnescape(catName)
	}
	sNumber, _ := strconv.Atoi(c.Request().Header.Get("number"))
	req := c.Values().Get("req")
	resp := c.Values().Get("rsp")
	method := c.Request().Method
	return &docParam{
		ApiKey:      cli.Doc.ApiKey,
		ApiToken:    cli.Doc.ApiToken,
		CatName:     catName,
		PageTitle:   pageTitle,
		PageContent: "",
		SNumber:     sNumber,
		Url:         c.Request().URL.String(),
		Method:      method,
		Req:         req,
		Resp:        resp,
	}
}

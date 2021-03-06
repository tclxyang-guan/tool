/**
* @Auther:gy
* @Date:2020/11/23 16:23
 */

package tool

import (
	"bytes"
	"container/list"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"strings"
)

var cli *client

type client struct {
	DocOpen int     `comment:"0不生成文档 1生成ShowDoc 2生成markDown"`
	ShowDoc ShowDoc `comment:"生成markDown路径地址"`
}
type ShowDoc struct {
	Url      string
	ApiKey   string
	ApiToken string
}

func EnableShowdoc(c *gin.Engine, docOpen int, showDoc ShowDoc) *client {
	if docOpen != 0 {
		c.Use(RequestParam())
	}
	return &client{docOpen, showDoc}
}
func RequestParam() gin.HandlerFunc {
	return func(c *gin.Context) {
		data, _ := c.GetRawData()
		fmt.Println(string(data))
		c.Set("reqBody", data)
		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(data))
	}
}
func BindJSON(c *gin.Context, req interface{}) error {
	err := c.BindJSON(req)
	if err != nil {
		return err
	}
	//验证参数
	if err := Validate.Struct(req); err != nil {
		return errors.New(TransError(err))
	}
	if cli.DocOpen != 0 {
		c.Set("req", req)
	}
	return err
}
func BindJSONNotWithValidate(c *gin.Context, req interface{}) error {
	err := c.BindJSON(req)
	if cli.DocOpen != 0 {
		c.Set("req", req)
	}
	return err
}
func GetClient() *client {
	return cli
}
func (cli *client) Upload(c *gin.Context) error {
	if cli.DocOpen == 0 {
		return nil
	} else if cli.DocOpen == 1 {
		return cli.uploadShowDoc(c)
	}
	return nil
}
func (cli *client) uploadShowDoc(c *gin.Context) error {
	pageTitle := c.Request.Header.Get("title")
	if pageTitle == "" {
		return errors.New("not found title in header")
	}
	pageTitle, _ = url.PathUnescape(pageTitle)
	//接口没有文件夹和顺序也没关系
	catName := c.Request.Header.Get("dir")
	if pageTitle != "" {
		catName, _ = url.PathUnescape(catName)
	}
	sNumber, _ := strconv.Atoi(c.Request.Header.Get("number"))
	req, _ := c.Get("req")
	reqBody, _ := c.Get("reqBody")
	method := c.Request.Method
	usdp := UploadShowDocParam{}
	usdp.ApiKey = cli.ShowDoc.ApiKey
	usdp.ApiToken = cli.ShowDoc.ApiToken
	usdp.CatName = catName
	usdp.PageTitle = pageTitle
	usdp.SNumber = sNumber
	reqParam := string(reqBody.([]byte))
	if strings.ToUpper(method) == "GET" || strings.ToUpper(method) == "DELETE" {
		reqParam = "同rul"
	}
	var resp, transfBody string
	rsp, ok := c.Get("rsp")
	if ok {
		respBody, _ := json.Marshal(rsp)
		respParam := string(respBody)
		respParam = strings.ReplaceAll(respParam, "{", "{\r\n\t")
		respParam = strings.ReplaceAll(respParam, "}", "\t\r\n}\r\n")
		respParam = strings.ReplaceAll(respParam, ",", ",\r\n\t")
		transfBody = stringTransf(string(respBody))
		resp = DatamapGenerateResp(rsp)
	}

	data, ok := rsp.(Rsp)
	if ok {
		resp += DatamapGenerateResp(data.Data)
		list, ok := data.Data.(ListRsp)
		if ok {
			resp += DatamapGenerateResp(list.List)
		}
		list1, ok := data.Data.(*ListRsp)
		if ok {
			resp += DatamapGenerateResp(list1.List)
		}
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
		DatamapGenerateReq(req) +
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
		transfBody + "\r\n" +
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
	if err != nil {
		return errors.New("marshal showdoc err:" + err.Error())
	}
	var m map[string]interface{}
	err = newRequest("post", cli.ShowDoc.Url, bytes.NewReader(b), "", &m)
	if err != nil {
		return errors.New(err.Error())
	}
	if m["error_code"].(float64) != 0 {
		return errors.New(m["error_message"].(string))
	}
	return nil
}

type UploadShowDocParam struct {
	ApiKey      string `json:"api_key"`
	ApiToken    string `json:"api_token"`
	CatName     string `json:"cat_name"`     //文件夹
	PageTitle   string `json:"page_title"`   //标题
	PageContent string `json:"page_content"` //内容
	SNumber     int    `json:"s_number"`     //排序
}

//生成请求的表
func DatamapGenerateReq(model interface{}) string {
	if model == nil {
		return "无"
	}
	str := "|参数名|必选|类型|说明|\r\n"
	str += "|:----    |:---|:----- |-----   |\r\n"
	if value, ok := model.(string); ok {
		if value == "" {
			return "无"
		}
		return str + value
	}
	elem := reflect.TypeOf(model)
	list1 := list.New()
	list1.PushBack(elem)
	str1 := ""
	m := map[string]bool{} //防止重复建相同的表结构
	for list1.Len() > 0 {
		str1 += "\r\n"
		e := list1.Remove(list1.Front()).(reflect.Type)
		for {
			if strings.HasPrefix(e.String(), "[]") || strings.HasPrefix(e.String(), "*") { //如果是切片需处理
				e = e.Elem()
			} else {
				break
			}
		}
		if strings.HasPrefix(e.String(), "map") {
			continue
		}
		tbl := e.Name() + "\r\n\r\n"
		tbl += reqRecursive(str, e, list1, m)
		if !strings.HasSuffix(tbl, str) { //防止重复建相同的表结构
			str1 += tbl
		}
	}
	return str1
}

func reqRecursive(str string, elem reflect.Type, list1 *list.List, m map[string]bool) string {
	if m[elem.String()] { //防止重复建相同的表结构
		return str
	} else {
		m[elem.String()] = true
	}
	for j := 0; j < elem.NumField(); j++ {
		isNeed := elem.Field(j).Tag.Get("req")
		if isNeed == "-" {
			continue
		}
		validate := elem.Field(j).Tag.Get("validate")
		require := "否"
		if strings.Contains(validate, "require") {
			require = "是"
		}
		gorm := elem.Field(j).Tag.Get("gorm")
		comment := ""
		gormList := strings.Split(gorm, ";")
		for _, v := range gormList {
			comment = elem.Field(j).Tag.Get("comment")
			if comment == "" && strings.Contains(v, "comment") {
				comment = strings.Replace(strings.Split(v, ":")[1], "'", "", -1)
			}
		}
		if comment == "" {
			comment = "暂无,若需要联系开发者"
		}
		json := elem.Field(j).Tag.Get("json")
		if elem.Field(j).Anonymous { //是镶嵌结构体
			str += reqRecursive("", elem.Field(j).Type, list1, m)
			continue
		}
		str += "|" + strings.ReplaceAll(json, ",omitempty", "") + "|" + require + "|" + elem.Field(j).Type.String() + "|" + comment + "|\r\n"
		if strings.Contains(elem.Field(j).Type.String(), ".") {
			list1.PushBack(elem.Field(j).Type)
			continue
		}
	}
	return str
}

//生成响应的表
func DatamapGenerateResp(model interface{}) string {
	if model == nil {
		return "无"
	}

	str := "|参数名|类型|说明|\r\n"
	str += "|:----    |:----- |-----   |\r\n"
	if value, ok := model.(string); ok {
		if value == "" {
			return "无"
		}
		return str + value
	}
	elem := reflect.TypeOf(model)
	list1 := list.New()
	list1.PushBack(elem)
	str1 := ""
	m := map[string]bool{} //防止重复建相同的表结构
	for list1.Len() > 0 {
		str1 += "\r\n"
		e := list1.Remove(list1.Front()).(reflect.Type)
		for {
			if strings.HasPrefix(e.String(), "[]") || strings.HasPrefix(e.String(), "*") { //如果是切片需处理
				e = e.Elem()
			} else {
				break
			}
		}
		if strings.HasPrefix(e.String(), "map") {
			continue
		}
		tbl := e.Name() + "\r\n\r\n"
		tbl += respRecursive(str, e, list1, m)
		if !strings.HasSuffix(tbl, str) { //防止重复建相同的表结构
			str1 += tbl
		}
	}
	return str1
}
func respRecursive(str string, elem reflect.Type, list1 *list.List, m map[string]bool) string {
	if m[elem.String()] { //防止重复建相同的表结构
		return str
	} else {
		m[elem.String()] = true
	}
	for i := 0; i < elem.NumField(); i++ {
		isNeed := elem.Field(i).Tag.Get("resp")
		if isNeed == "-" {
			continue
		}
		json := elem.Field(i).Tag.Get("json")
		comment := elem.Field(i).Tag.Get("comment")
		if comment == "" {
			gorm := elem.Field(i).Tag.Get("gorm")
			gormList := strings.Split(gorm, ";")
			for _, v := range gormList {
				if strings.Contains(v, "comment") {
					comment = strings.Replace(strings.Split(v, ":")[1], "'", "", -1)
				}
			}
		}
		if comment == "" {
			comment = "暂无,若需要联系开发者"
		}
		if elem.Field(i).Anonymous { //是否镶嵌结构体
			str += respRecursive("", elem.Field(i).Type, list1, m)
			continue
		}
		str += "|" + strings.ReplaceAll(json, ",omitempty", "") + "|" + elem.Field(i).Type.String() + "|" + comment + "|\r\n"
		if strings.Contains(elem.Field(i).Type.String(), ".") {
			list1.PushBack(elem.Field(i).Type)
			continue
		}
	}
	return str
}
func stringTransf(str string) string {
	index := 0
	newstr := ""
	start, end := 0, 0
	ind := 0
	for i, v := range str {
		if str[i] == '"' {
			ind++
			if ind%2 == 1 {
				start = i + 1
			} else {
				end = i
				newstr += `"` + str[start:end] + `"`
				continue
			}
		}
		if ind%2 == 1 {
			continue
		}
		if str[i] == '{' {
			newstr += "{\r\n"
			index++
			for i := 0; i < index; i++ {
				newstr += "\t"
			}
		} else if str[i] == ',' {
			newstr += ",\r\n"
			for i := 0; i < index; i++ {
				newstr += "\t"
			}
		} else if str[i] == '}' {
			index--
			newstr += "\r\n"
			for i := 0; i < index; i++ {
				newstr += "\t"
			}
			newstr += "}"
		} else {
			newstr += string([]byte{byte(v)})
		}
	}
	return newstr
}
func newRequest(method, url string, body io.Reader, ContentType string, data interface{}) (err error) {
	method = strings.ToUpper(method)
	if method == "POST" {
		client := &http.Client{}
		req, err := http.NewRequest(method, url, body)
		if err != nil {
			return err
		}
		if ContentType == "" {
			req.Header.Set("content-type", "application/json; charset=utf-8")
		} else {
			req.Header.Set("content-type", ContentType)
		}
		resp, err := client.Do(req)
		if err != nil {
			return err
		}
		defer resp.Body.Close()
		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		return json.Unmarshal(b, data)
	} else if method == "GET" {
		resp, err := http.Get(url)
		if err != nil {
			return err
		}
		defer resp.Body.Close()
		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		return json.Unmarshal(b, data)
	}
	return errors.New("请求方式错误")
}

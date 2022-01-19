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
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"reflect"
	"strings"
	"text/template"
)

//可通过修改模板字符串来变更模板
var ShowDocTemplateStr = "##### 简要描述\n\n- {{.title}}\n\n##### 请求URL\n- ` {{.url}} `\n\n##### 请求方式\n- {{.method}}\n\n##### 参数\n\n{{.req}}\n\n\n##### 请求示例\n\n```\n{{.reqExample}}\n```\n##### 返回示例\n\n```\n{{.respExample}}\n```\n\n##### 返回参数说明\n\n{{.resp}}\n\n\n\n##### 备注\n{{.remark}}"
var MarkdownTemplateStr = "####{{.title}}\n#####请求地址\n```text\n{{.url}}\n```\n#####请求方式\n```text\n{{.method}}\n```\n#####请求参数\n{{.req}}\n#####请求示例\n```text\n{{.reqExample}}\n```\n#####返回示例\n```text\n{{.respExample}}\n```\n#####返回参数说明\n{{.resp}}\n#####备注\n{{.remark}}"

var cli *client

type client struct {
	DocOpen int    `comment:"0不生成文档 1生成ShowDoc 2生成markDown"`
	Doc     Doc    `comment:"doc配置"`
	Option  Option `json:"option" comment:"自己框架的实现"`
	tmp     *template.Template
}

type Doc struct {
	Url      string `comment:"showdoc 项目设置->开放api->下面的开放api连接文档"`
	ApiKey   string `comment:"showdoc 项目设置->开放api->api_key"`
	ApiToken string `comment:"showdoc 项目设置->开放api->api_token"`
	DocDir   string `comment:"Markdown 项目文档文件夹 绝对路径"`
}

var DefaultOption = Gin{}

//目前只支持 gin和iris框架
func EnableDoc(docOpen int, doc Doc, option Option) {
	var err error
	cli = &client{DocOpen: docOpen, Doc: doc}
	if option != nil {
		cli.Option = option
	} else {
		cli.Option = DefaultOption
	}
	if docOpen == 1 {
		cli.tmp, err = template.New("showdoc").Parse(ShowDocTemplateStr)
		if err != nil {
			panic("showdoc模板错误")
		}
	} else {
		cli.tmp, err = template.New("markdown").Parse(MarkdownTemplateStr)
		if err != nil {
			panic("markdown模板错误")
		}
	}

}
func BindJSON(c interface{}, req interface{}) error {
	cli.Option.BindJSON(c, req)
	if cli.DocOpen != 0 {
		cli.Option.Set(c, "req", req)
	}
	//验证参数.BindJSON()
	if err := Validate.Struct(req); err != nil {
		return errors.New(TransError(err))
	}
	return nil
}
func BindJSONNotWithValidate(c interface{}, req interface{}) error {
	cli.Option.BindJSON(c, req)
	if cli.DocOpen != 0 {
		cli.Option.Set(c, "req", req)
	}
	return nil
}
func saveDoc(param *docParam) {
	var request, reqParam string
	if strings.ToUpper(param.Method) == "GET" || strings.ToUpper(param.Method) == "DELETE" {
		reqParam = "同rul"
	} else {
		if param.Req != nil {
			b, _ := json.MarshalIndent(param.Req, "", "\t")
			reqParam = string(b)
			vt := reflect.TypeOf(param.Req)
			vv := reflect.ValueOf(param.Req)
			request = DataGenerate(1, data{
				Value: vv,
				Type:  vt,
			})
		}
	}
	var resp, respParam string
	if param.Resp != nil {
		respBody, _ := json.MarshalIndent(param.Resp, "", "\t")
		respParam = string(respBody)
		vt := reflect.TypeOf(param.Resp)
		vv := reflect.ValueOf(param.Resp)
		resp = DataGenerate(2, data{
			Value: vv,
			Type:  vt,
		})
	}
	buf := &bytes.Buffer{}
	m1 := map[string]interface{}{
		"title":       param.PageTitle,
		"url":         param.Url,
		"method":      param.Method,
		"req":         request,
		"reqExample":  reqParam,
		"respExample": respParam,
		"resp":        resp,
		"remark":      "",
	}
	err := cli.tmp.Execute(buf, m1)
	if err != nil {
		fmt.Println("cli.Doc.tmp.Execute err:" + err.Error())
		return
	}
	param.PageContent = buf.String()
	if cli.DocOpen == 1 {
		b, _ := json.Marshal(param)
		var m map[string]interface{}
		err = newRequest("post", cli.Doc.Url, bytes.NewReader(b), "", &m)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		if m["error_code"].(float64) != 0 {
			fmt.Println(m["error_message"])
		}
	} else {
		err = os.MkdirAll(cli.Doc.DocDir+"/"+param.CatName, os.ModePerm)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		file, err := os.OpenFile(cli.Doc.DocDir+"/"+param.CatName+"/"+param.PageTitle+".md", os.O_CREATE|os.O_TRUNC, os.ModePerm)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		defer file.Close()
		file.Write(buf.Bytes())
	}

	return
}

type docParam struct {
	ApiKey      string `json:"api_key"`
	ApiToken    string `json:"api_token"`
	CatName     string `json:"cat_name"`     //文件夹
	PageTitle   string `json:"page_title"`   //标题
	PageContent string `json:"page_content"` //内容
	SNumber     int    `json:"s_number"`     //排序

	Url    string      `json:"-"`
	Method string      `json:"-"`
	Req    interface{} `json:"-"`
	Resp   interface{} `json:"-"`
}

func getStrPrefix(t int) string {
	if t == 1 { //req
		str := "|参数名|必选|类型|说明|\r\n"
		str += "|:----:|:---:|:-----:|:-----:|\r\n"
		return str
	} else { //resp
		str := "|参数名|类型|说明|\r\n"
		str += "|:----:|:-----:|:-----:|\r\n"
		return str
	}
}
func getStrJson(t int, jsonStr, require, tStr, comment string) string {
	if t == 1 { //req
		return "|" + jsonStr + "|" + require + "|" + tStr + "|" + comment + "|\r\n"
	} else { //resp
		return "|" + jsonStr + "|" + tStr + "|" + comment + "|\r\n"
	}
}

//生成响应的表
func DataGenerate(t int, d data) string {
	str := getStrPrefix(t)
	list1 := list.New()
	list1.PushBack(d)
	str1 := ""
	m := map[string]bool{} //防止重复建相同的表结构
	for list1.Len() > 0 {
		str1 += "\r\n"
		e := list1.Remove(list1.Front()).(data)
		for {
			if e.Type.Kind() == reflect.Ptr {
				e.Type = e.Type.Elem()
				if e.Value.Kind() == reflect.Ptr {
					e.Value = e.Value.Elem()
				}
			}
			if e.Type.Kind() == reflect.Slice {
				e.Type = e.Type.Elem()
				if e.Value.Len() > 0 {
					e.Value = e.Value.Index(0)
				}
			} else {
				break
			}
		}
		if e.Type.Kind() != reflect.Struct && e.Type.Kind() != reflect.Map {
			continue
		}
		tbl := e.Type.Name() + "\r\n\r\n"
		if e.Type.Kind() == reflect.Map {
			tbl += respMap(str, e, list1)
			str1 += tbl
			continue
		}
		tbl += Recursive(str, t, e, list1, m)
		if !strings.HasSuffix(tbl, str) { //防止重复建相同的表结构
			str1 += tbl
		}
	}
	return str1
}

func Recursive(str string, t int, d data, list1 *list.List, m map[string]bool) string {
	elem := d.Type
	if !d.Anonymous {
		if m[elem.String()] { //防止重复建相同的表结构
			return str
		} else {
			m[elem.String()] = true
		}
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
		validate := elem.Field(i).Tag.Get("validate")
		require := "否"
		if strings.Contains(validate, "require") {
			require = "是"
		}
		if elem.Field(i).Anonymous { //是否镶嵌结构体
			if !d.Value.IsZero() {
				if d.Value.Kind() == reflect.Slice {
					str += Recursive("", t, data{
						Type:      elem.Field(i).Type,
						Value:     reflect.ValueOf(nil),
						Anonymous: true,
					}, list1, m)
				} else {
					str += Recursive("", t, data{
						Type:      elem.Field(i).Type,
						Value:     d.Value.Field(i),
						Anonymous: true,
					}, list1, m)
				}
			}
			continue
		}
		if elem.Field(i).Type.Kind() == reflect.Interface {
			if !d.Value.IsZero() {
				d1 := data{}
				if d.Value.Kind() == reflect.Slice {
					d1.Type = elem.Field(i).Type
					d1.Value = reflect.ValueOf(nil)
				} else {
					d1.Type = d.Value.Field(i).Elem().Type()
					d1.Value = d.Value.Field(i).Elem()
				}
				list1.PushBack(d1)
			}
		}
		str += getStrJson(t, strings.ReplaceAll(json, ",omitempty", ""), require, elem.Field(i).Type.String(), comment)
		if elem.Field(i).Type.Kind() == reflect.Struct || elem.Field(i).Type.Kind() == reflect.Ptr || elem.Field(i).Type.Kind() == reflect.Slice {
			list1.PushBack(data{
				Type:  elem.Field(i).Type,
				Value: d.Value,
			})
			continue
		}
	}
	return str
}
func respMap(str string, d data, list1 *list.List) string {
	r := d.Value.MapRange()
	for r.Next() {
		if r.Value().IsNil() {
			continue
		}
		str += "|" + r.Key().String() + "|" + r.Value().Elem().Type().String() + "|暂无若需要请联系开发人员|\r\n"
		if r.Value().Elem().Type().Kind() == reflect.Slice || r.Value().Elem().Type().Kind() == reflect.Struct || r.Value().Elem().Type().Kind() == reflect.Ptr {
			list1.PushBack(data{
				Type:  r.Value().Elem().Type(),
				Value: r.Value().Elem(),
			})
			continue
		}
		if r.Value().Elem().Type().Kind() == reflect.Interface && !r.Value().IsNil() {
			d1 := data{
				Type:  r.Value().Elem().Type(),
				Value: r.Value().Elem(),
			}
			list1.PushBack(d1)
		}
	}
	return str
}

type data struct {
	Type      reflect.Type
	Value     reflect.Value
	Anonymous bool `json:"anonymous" comment:"是否嵌入结构体"`
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

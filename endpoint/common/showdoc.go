/**
* @Auther:gy
* @Date:2020/11/23 16:23
 */

package common

import (
	"reflect"
	"strings"
	"transfDoc/models"
)

var ShowDocMap = map[string]ShowDocData{
	"BuildCreate": {
		&models.Build{},
		&Rsp{Data: nil},
		"楼栋",
		"楼栋创建接口",
		0,
	}, //示例
	"BuildPage": {
		nil,
		&Rsp{Data: &models.Build{}},
		"楼栋",
		"楼栋分页接口",
		0,
	}, //示例
	"BuildUpdate": {
		&models.Build{},
		&Rsp{Data: nil},
		"楼栋",
		"楼栋修改接口",
		0,
	}, //示例
	"BuildUpdateStatus": {
		nil,
		&Rsp{Data: nil},
		"楼栋",
		"楼栋修改状态接口",
		0,
	}, //示例
}

type ShowDocData struct {
	Req       interface{}
	Resp      *Rsp
	CatName   string //文件夹
	PageTitle string //标题
	SNumber   int    //排序
}
type UploadShowDocParam struct {
	ApiKey      string `json:"api_key"`
	ApiToken    string `json:"api_token"`
	CatName     string `json:"cat_name"`     //文件夹
	PageTitle   string `json:"page_title"`   //标题
	PageContent string `json:"page_content"` //内容
	SNumber     int    `json:"s_number"`     //排序
}

func reqRecursive(str string, typ reflect.Type) string {
	if !strings.Contains(typ.String(), ".") {
		return ""
	}
	for j := 0; j < typ.NumField(); j++ {
		isNeed := typ.Field(j).Tag.Get("req")
		if isNeed == "-" {
			continue
		}
		if strings.Contains(typ.Field(j).Type.String(), ".") {
			str += reqRecursive("", typ.Field(j).Type)
			continue
		}
		validate := typ.Field(j).Tag.Get("validate")
		require := "否"
		if strings.Contains(validate, "require") {
			require = "是"
		}
		gorm := typ.Field(j).Tag.Get("gorm")
		comment := ""
		gormList := strings.Split(gorm, ";")
		for _, v := range gormList {
			comment = typ.Field(j).Tag.Get("comment")
			if comment == "" && strings.Contains(v, "comment") {
				comment = strings.Replace(strings.Split(v, ":")[1], "'", "", -1)
			}
		}
		if comment == "" {
			comment = "暂无,若需要联系开发者"
		}
		str += "|" + strings.ReplaceAll(typ.Field(j).Tag.Get("json"), ",omitempty", "") + "|" + require + "|" + typ.Field(j).Type.String() + "|" + comment + "|\r\n"
	}
	return str
}

//生成请求的表
func DatamapGenerateReq(model interface{}) string {
	if model == nil {
		return "无(若有请求示例，按示例走)"
	}
	elem := reflect.TypeOf(model).Elem()
	str := "|参数名|必选|类型|说明|\r\n"
	str += "|:----    |:---|:----- |-----   |\r\n"

	return reqRecursive(str, elem)
}

//生成响应的表
func DatamapGenerateResp(model interface{}) string {
	if model == nil {
		return "无"
	}
	elem := reflect.TypeOf(model).Elem()
	str := "|参数名|类型|说明|\r\n"
	str += "|:----    |:----- |-----   |\r\n"

	return respRecursive(str, elem)
}
func respRecursive(str string, elem reflect.Type) string {
	if !strings.Contains(elem.String(), ".") {
		return ""
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
		if strings.Contains(elem.Field(i).Type.String(), ".") {
			//判断是否以空表头结尾
			if strings.HasSuffix(str, "|参数名|类型|说明|\r\n|:----    |:----- |-----   |\r\n") {
				str = strings.ReplaceAll(str, "|参数名|类型|说明|\r\n|:----    |:----- |-----   |\r\n", "")
			} else {
				str += "|" + strings.ReplaceAll(json, ",omitempty", "") + "|" + elem.Field(i).Type.String() + "|" + comment + "|\r\n"
			}
			str1 := "\r\n"
			str1 += "|参数名|类型|说明|\r\n"
			str1 += "|:----    |:----- |-----   |\r\n"
			str += respRecursive(str1, elem.Field(i).Type)
			continue
		}
		str += "|" + strings.ReplaceAll(json, ",omitempty", "") + "|" + elem.Field(i).Type.String() + "|" + comment + "|\r\n"
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

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

//Req "|字段一|是|int|说明一|\r\n|字段二|是|int|说明二|\r\n"  可直接传入字符串也可为结构体指针
//Rsp "|字段一|int|说明一|\r\n|字段二|int|说明二|\r\n" 可直接传入字符串也可为结构体指针
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
	Req       interface{} //请求参数
	Resp      *Rsp        //返回参数
	CatName   string      //文件夹
	PageTitle string      //标题
	SNumber   int         //排序
}
type UploadShowDocParam struct {
	ApiKey      string `json:"api_key"`
	ApiToken    string `json:"api_token"`
	CatName     string `json:"cat_name"`     //文件夹
	PageTitle   string `json:"page_title"`   //标题
	PageContent string `json:"page_content"` //内容
	SNumber     int    `json:"s_number"`     //排序
}

func reqRecursive(str string, elem reflect.Type) string {
	if !strings.Contains(elem.String(), ".") {
		return ""
	}
	if strings.HasPrefix(elem.String(), "*") || strings.Contains(elem.String(), "[]") { //如果是指针无法转义或者切片
		return ""
	}
	for j := 0; j < elem.NumField(); j++ {
		isNeed := elem.Field(j).Tag.Get("req")
		if isNeed == "-" {
			continue
		}
		if strings.Contains(elem.Field(j).Type.String(), ".") {
			str += reqRecursive("", elem.Field(j).Type)
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
		str += "|" + strings.ReplaceAll(elem.Field(j).Tag.Get("json"), ",omitempty", "") + "|" + require + "|" + elem.Field(j).Type.String() + "|" + comment + "|\r\n"
	}
	return str
}

//生成请求的表
func DatamapGenerateReq(model interface{}) string {
	if model == nil {
		return "无"
	}
	str := "|参数名|必选|类型|说明|\r\n"
	str += "|:----    |:---|:----- |-----   |\r\n"
	if value, ok := model.(string); ok {
		return str + value
	}
	elem := reflect.TypeOf(model).Elem()

	return reqRecursive(str, elem)
}

//生成响应的表
func DatamapGenerateResp(model interface{}) string {
	if model == nil {
		return "无"
	}

	str := "|参数名|类型|说明|\r\n"
	str += "|:----    |:----- |-----   |\r\n"
	if value, ok := model.(string); ok {
		return str + value
	}
	elem := reflect.TypeOf(model).Elem()
	return respRecursive(str, elem)
}
func respRecursive(str string, elem reflect.Type) string {
	if !strings.Contains(elem.String(), ".") {
		return ""
	}
	if strings.HasPrefix(elem.String(), "*") || strings.Contains(elem.String(), "[]") { //如果是指针无法转义或者切片
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

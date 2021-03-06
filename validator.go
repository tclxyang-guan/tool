/**
* @Auther:gy
* @Date:2021/3/5 19:44
 */

package tool

import (
	zhongwen "github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	zh_translations "github.com/go-playground/validator/v10/translations/zh"
	"reflect"
	"strings"
)

var Validate *validator.Validate
var Trans ut.Translator

func EnableValidate() {
	zh := zhongwen.New()
	uni := ut.New(zh, zh)
	Trans, _ = uni.GetTranslator("zh")

	Validate = validator.New()
	Validate.RegisterTagNameFunc(func(field reflect.StructField) string {
		label := field.Tag.Get("comment")
		if label == "" {
			gorm := field.Tag.Get("gorm")
			gormList := strings.Split(gorm, ";")
			for _, v := range gormList {
				if strings.Contains(v, "comment") {
					comment := strings.Replace(strings.Split(v, ":")[1], "'", "", -1)
					comments := strings.Split(comment, " ")
					if len(comments) > 0 {
						label = comments[0]
					}
				}
			}
			if label == "" {
				return field.Name
			}
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

/**
* @Auther:gy
* @Date:2020/11/23 16:23
 */
package utils

import (
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

const (
	loginUrl  = ""
	createUrl = ""
)

//body = strings.NewReader(string(by)
func NewRequest(method, url string, body io.Reader, ContentType string) (by []byte, err error) {
	method = strings.ToUpper(method)
	if method == "POST" {
		client := &http.Client{}
		req, err := http.NewRequest(method, url, body)
		if err != nil {
			return nil, err
		}
		if ContentType == "" {
			req.Header.Set("content-type", "application/json; charset=utf-8")
		} else {
			req.Header.Set("content-type", ContentType)
		}
		resp, err := client.Do(req)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()
		return ioutil.ReadAll(resp.Body)
	} else if method == "GET" {
		resp, err := http.Get(url)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()
		return ioutil.ReadAll(resp.Body)
	}
	return nil, errors.New("请求方式错误")
}

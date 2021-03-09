/**
* @Auther:gy
* @Date:2021/3/9 20:01
 */

package tool

import (
	"fmt"
	"testing"
)

type testModel struct {
	Name string `json:"name" comment:"名称"`
}

func TestDatamapGenerateResp(t *testing.T) {
	fmt.Println(DatamapGenerateResp(testModel{}))
	fmt.Println(DatamapGenerateResp(map[string]interface{}{}))
	fmt.Println(DatamapGenerateResp([]string{}))
}

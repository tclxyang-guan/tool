/**
* @Auther:gy
* @Date:2021/3/9 20:06
 */

package tool

import (
	"fmt"
	"testing"
)

type user struct {
	Model
	Name string
}

func TestEnableMysql(t *testing.T) {
	db, err := EnableMysql(MysqlConf{
		Address:         "",
		Username:        "",
		Password:        "",
		DbName:          "test",
		Prefix:          "t_",
		MaxOpenConns:    1,
		MaxIdleConns:    1,
		ConnMaxLifetime: 10,
	})
	fmt.Println(err)
	u := user{}
	fmt.Println(db.First(&u))
}

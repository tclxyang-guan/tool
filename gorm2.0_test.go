/**
* @Auther:gy
* @Date:2021/3/9 20:06
 */

package tool

import (
	"fmt"
	"testing"
)

func TestEnableMysql2(t *testing.T) {
	db, err := EnableMysql2(MysqlConf{
		Address:         "",
		Username:        "",
		Password:        "",
		DbName:          "",
		Prefix:          "",
		MaxOpenConns:    1,
		MaxIdleConns:    1,
		ConnMaxLifetime: 10,
	})
	fmt.Println(err)
	db.AutoMigrate(&user{})
	u := user{}
	u.Name = "test"
	db.Create(&u)
	db.Unscoped().Where("name=?", "test").Delete(&user{})
}

package tool

import (
	"fmt"
	"reflect"
	"testing"
)

type user struct {
	role
	Name string `json:"name" comment:"用户名称"`
	//	Role *role `json:"role" comment:"角色信息"`
	//	Role role `json:"role" comment:"角色信息"`
	//	Role []role `json:"role" comment:"角色信息"`
	//	Role []*role `json:"role" comment:"角色信息"`
	//	Role *[]role `json:"role" comment:"角色信息"`
	Role interface{} `json:"role" comment:"角色信息"`
}
type role struct {
	Age int `json:"age" comment:"年龄"`
}

func TestDataGenerate(t *testing.T) {
	ds := []data{}
	u1 := &user{}
	u1.Role = &[]role{}
	u2 := user{}
	u3 := []user{}
	u4 := []*user{}
	u5 := &[]user{}
	ds = append(ds, data{
		Type:  reflect.TypeOf(u1),
		Value: reflect.ValueOf(u1),
	})
	ds = append(ds, data{
		Type:  reflect.TypeOf(u2),
		Value: reflect.ValueOf(u2),
	})
	ds = append(ds, data{
		Type:  reflect.TypeOf(u3),
		Value: reflect.ValueOf(u3),
	})
	ds = append(ds, data{
		Type:  reflect.TypeOf(u4),
		Value: reflect.ValueOf(u4),
	})
	ds = append(ds, data{
		Type:  reflect.TypeOf(u5),
		Value: reflect.ValueOf(u5),
	})
	for i, d := range ds {
		fmt.Println(i)
		fmt.Println(DataGenerate(2, d))
	}
}

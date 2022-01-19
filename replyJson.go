/**
* @Auther:gy
* @Date:2021/8/20 14:40
 */

package tool

/*
返回json数据到前台 配合showdoc使用
*/
// c context
//data 返回的json数据
func ReplyJson(c, data interface{}) {
	cli.Option.Set(c, "rsp", data)
	if cli.DocOpen != 0 {
		param := cli.Option.GetDocParam(c)
		if param != nil {
			go saveDoc(param)
		}
	}
	cli.Option.JSON(c, data)
}

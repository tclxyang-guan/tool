/**
* @Auther:gy
* @Date:2020/11/21 17:32
 */

package utils

//获取最后一个/位置 j为0时无正斜杠
func GetLastLine(str string) (j int) {
	l := len(str)
	for i := l - 1; i > 0; i-- {
		if str[i] == '/' {
			j = i
			return
		}
	}
	return
}

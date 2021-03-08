package tool

var MsgFlags = map[int]string{
	SUCCESS:            "ok",
	INVALID_PARAMS:     "请求参数错误",
	FAILE_TO_CREATE_OP: "操作失败:",
}

// GetMsg get error information based on Code
func GetMsg(code int) string {
	msg, ok := MsgFlags[code]
	if ok {
		return msg
	}
	return ""
}

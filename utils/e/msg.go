package e

var MsgFlags = map[int]string{
	SUCCESS:        "ok",
	ERROR:          "fail",
	INVALID_PARAMS: "请求参数错误",
	DATABASE_ERROR: "数据库异常",

	ERROR_AUTH_CHECK_TOKEN_FAIL:    "Token鉴权失败",
	ERROR_AUTH_CHECK_TOKEN_TIMEOUT: "Token已超时",
	ERROR_AUTH_TOKEN:               "Token生成失败",
	ERROR_AUTH:                     "Token错误",
	FAILE_TO_GET_OPENID:            "无法获取openid",
	ERROR_LOGIN_FAIL:               "登录失败",
	FAILE_TO_CREATE_OP:             "操作失败:",
	NOT_FOUND_RECORD:               "记录未找到",
	CACHE_ERROR:                    "缓存查询失败",
	WS_ERROR:                       "WebSocket连接失败",
}

// GetMsg get error information based on Code
func GetMsg(code int) string {
	msg, ok := MsgFlags[code]
	if ok {
		return msg
	}

	return MsgFlags[ERROR]
}

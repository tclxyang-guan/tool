package tool

type Option interface {
	BindJSON(c, req interface{}) error
	Set(c interface{}, key string, value interface{})
	JSON(c, data interface{})
	GetDocParam(ctx interface{}) *docParam
}

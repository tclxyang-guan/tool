# transfDoc
golang集成showdoc 调取接口自动生成在线文档

1、在showdoc文档中获取apikey、apitoken、showdocurl写入配置文件中

若公有showdoc  showdocurl: https://www.showdoc.cc/server/api/item/updateByApi

若私有showdoc  showdocurl: http://xx.com/server/index.php?s=/api/item/updateByApi

2、在路由层添加前置拦截获取请求体

r.Use(middleware.RequestParam())

3、在包装的响应方法中调用DeferShowDoc方法

4、配置ShowDocMap此项目采用的请求头中传入dockey:key(该key为map中的key),获取showdoc的一些基础信息，具体请看ShowDocData结构体

5、若请求参数不需要某些字段请使用 req:"-",若返回参数中不需要某些字段请使用resp:"-",若必填validate:"required"(结合validate参数校验),说明则是 common:"aa"或者`gorm:"comment:'数据id'"，
参数名使用json:"id"。若无json则不会在请求参数、返回参数中显示该字段,具体可拉取代码运行参照build


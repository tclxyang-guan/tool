# transfDoc
golang集成showdoc 调取接口自动生成在线文档

1、在showdoc文档中获取apikey、apitoken、showdocurl写入配置文件中

若公有showdoc  showdocurl: https://www.showdoc.cc/server/api/item/updateByApi

若私有showdoc  showdocurl: http://xx.com/server/index.php?s=/api/item/updateByApi

2、在路由层添加前置拦截获取请求体(一定记得发布线上注释掉)

r.Use(middleware.RequestParam())

3、在包装的响应方法中调用DeferShowDoc方法

4、配置ShowDocMap此项目采用的请求头中传入dir(文件夹) title(接口名称),获取showdoc的一些基础信息，具体请看ShowDocData结构体。

5、若请求参数不需要某些字段请使用 req:"-",若返回参数中不需要某些字段请使用resp:"-",若必填validate:"required"(结合validate参数校验),说明则是 common:"aa"或者`gorm:"comment:'数据id'"，
参数名使用json:"id"。若无json则不会在请求参数、返回参数中显示该字段,具体可拉取代码运行参照build

6、因为请求头无法传入中文所以需要进行转码 postman右键有个EncodeUrlComponent
例如：

curl --location --request POST '127.0.0.1:8080/test/v1/build' \
--header 'dir: test' \
--header 'title: %E6%A5%BC%E6%A0%8B%E5%88%9B%E5%BB%BA' \
--header 'Content-Type: application/json' \
--data-raw '{
    "build_name":"楼栋1",
    "organization_id":1,
    "community_id":1
}'

生成好的文档路径
https://www.showdoc.com.cn/1263702264386833?page_id=6345620950068093

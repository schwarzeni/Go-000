# Go 工程化

[toc]

## 作业

> 按照自己的构想，写一个项目满足基本的目录结构和工程，代码需要包含对数据层、业务层、API 注册，以及 main 函数对于服务的注册和启动，信号处理，使用 Wire 构建依赖。可以使用自己熟悉的框架。

目录结构：

```txt
.
├── cmd
│      └── myapp
│          └── myapp.go  // 程序主入口
├── go.mod
├── go.sum
├── internal
│        └── app
│            └── myapp
│                ├── dao
│                │     └── dao.go  // DAO 层
│                ├── dto
│                │     └── dto.go  // DTO 数据对象 
│                ├── errcode 
│                │     └── errcode.go // 业务错误码
│                ├── model  
│                │     └── article.go // 数据表映射对象
│                ├── server 
│                │     ├── handler.go
│                │     ├── server.go  // HTTP 服务
│                │     ├── wire.go
│                │     └── wire_gen.go
│                └── service 
│                    └── service.go   // Service 层
└── pkg
    ├── appgroup
    │        └── appgroup.go // 基于errgroup 的封装，用于控制多应用的生命周期
    └── db
        └── gorm.go // gorm 初始化
```

## 目录结构

https://github.com/golang-standards/project-layout/blob/master/README_zh.md

https://github.com/go-kratos/service-layout

## 数据库orm

https://github.com/facebook/ent

## 依赖注入，控制反转

- 方便测试
- 单次初始化与复用

`NewSVC(IDAO)`

https://blog.golang.org/wire

## 生命周期

[https://github.com/go-kratos/kratos/blob/v2/app.go](https://github.com/go-kratos/kratos/blob/v2/app.go)

各个服务的启动与终止由 kit 来完成

```go
type Hook struct {
	OnStart func(context.Context) error
	OnStop  func(context.Context) error
}
```



----

其他：

go io 接口

Code Review，依赖注入，https://www.youtube.com/watch?v=ifBUfIb7kdo

---

## api 设计

**Google API 设计指南 [https://cloud.google.com/apis/design?hl=zh-cn](https://cloud.google.com/apis/design?hl=zh-cn)**

yapi

API 文件夹同步到另一个仓库中，自动同步（p18）

定义路径：e.g. kratos/demo/v1/demo.proto

定义 pb 时，request 与 response 都定义对象，方便扩展

**使用 pb 来定义配置文件** ：pb 基础类型识别是否填写：将其包装为结构体，使用指针是否为空来判断；[https://github.com/protocolbuffers/protobuf/blob/master/src/google/protobuf/wrappers.proto](https://github.com/protocolbuffers/protobuf/blob/master/src/google/protobuf/wrappers.proto)

```protobuf
// Wrapper message for `float`.
//
// The JSON representation for `FloatValue` is JSON number.
message FloatValue {
  // The float value.
  float value = 1;
}
```

部分更新字段：使用 mask 

```protobuf
message UpdateBookRequest {
	Book book = 1;
	google.protobuf.FieldMask mask = 2;
}
```



## 处理错误

不建议使用全局错误码。需要将别人的错误码转换为自己的错误码

利用好 HTTP status code，使用一小组标准错误配合大量资源

错误信息设计：[https://github.com/go-kratos/kratos/blob/v2/errors/errors.go](https://github.com/go-kratos/kratos/blob/v2/errors/errors.go)

```go
// StatusError contains an error response from the server.
type StatusError struct {
	// Code is the gRPC response status code and will always be populated.
	Code int `json:"code"`
	// Message is the server response message and is only populated when
	// explicitly referenced by the JSON server response.
	Message string `json:"message"`
	// Details provide more context to an error.
	Details []interface{} `json:"details"`
}

// ErrorInfo is a detailed error code & message from the API frontend.
type ErrorInfo struct {
	// Reason is the typed error code. For example: "some_example".
	Reason string `json:"reason"`
	// Message is the human-readable description of the error.
	Message string `json:"message"`
}
```

service error --> grpc error --> service error (插件形式)

---

## 配置管理



Functional Option：

- [https://github.com/go-kratos/kratos/blob/v2/app.go](https://github.com/go-kratos/kratos/blob/v2/app.go)
- [https://github.com/go-kratos/kratos/blob/master/pkg/cache/redis/conn.go](https://github.com/go-kratos/kratos/blob/master/pkg/cache/redis/conn.go)

配置初始化 与 系统初始化 分离

```txt
+-------------------+
|  Config Web UI    +------------+
+-------------------+            |
                                 v
+-------------------+        +---+-----------+      +----------+
|  Config API       +------->+  Config Data  +----->+ System   |
+-------------------+        +---+-----------+      +----------+
                                 ^
+-------------------+            |
|  Config File      +------------+
+-------------------+
```

## 测试

https://blog.golang.org/subtests

table-driven --> subtest + gomock

docker-compose

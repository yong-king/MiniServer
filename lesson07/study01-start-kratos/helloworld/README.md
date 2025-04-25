# kratos

## 快速开始，项目初始化

1. 环境准备
    - [go](https://golang.org/dl/)
    - [protoc](https://github.com/protocolbuffers/protobuf)
    - [protoc-gen-go](https://github.com/protocolbuffers/protobuf-go)
    1. go env -w GO111MODULE=on
2. **kratos 命令工具**
    1. kratos 是与 Kratos 框架配套的脚手架工具，kratos 能够
        - 通过模板快速创建项目
        - 快速创建与生成 protoc 文件
        - 使用开发过程中常用的命令
        - 极大提高开发效率，减轻心智负担
    2. 安装**CLI工具**
        
        ```go
        go install github.com/go-kratos/kratos/cmd/kratos/v2@latest
        ```
        
3. 创建项目
    
    ```go
    # 使用默认模板创建项目
    kratos new helloworld
    ```
    

# 个人业务

## 编写proto配置文件

```go
// 添加 Proto 文件
kratos proto add api/xxx/v1/demo.proto
// 生成 Proto 代码
// 可以直接通过 make 命令生成
make api

// 或使用 kratos cli 进行生成
kratos proto client api/xxx/v1/demo.proto

// 生成 Service 代码
kratos proto server api/xxx/v1/demo.proto -t internal/service
```

## 添加service代码

```go
func (s *TodoService) CreateTodo(ctx context.Context, req *pb.CreateTodoRequest) (*pb.CreateTodoReply, error) {
	// 检验参数
	...
	// 调用业务逻辑
	...
	// 返回
	return &pb.CreateTodoReply{}, nil
}

// 调用业务逻辑需要用到biz的逻辑文件，因此
type TodoService struct {
	pb.UnimplementedTodoServer

	uc *biz.TodoUsecase
}

// 依赖提供 
var ProviderSet = wire.NewSet(NewTodoService)
```

## 添加biz逻辑代码

```go
// 根据你的业务逻辑修改biz代码文件中代码
// XXX is a XXX model. 】
你的信息的结构体

// XXXRepo is a XXX repo.
你的操作逻辑接口

// XXXUsecase is a XXX usecase.
封装接口的结构头

// NewXXXUsecase new a XXX usecase.
结构体构建方法

// CreateXXX creates a XXX, and returns the new XXX.
外部调用的结构体的创建逻辑方法

// 依赖提供
var ProviderSet = wire.NewSet(NewXXXUsecase)
```

## 添加data代码

```go
1. 实现biz中的接口的方法
2. 依赖提供
var ProviderSet = wire.NewSet(NewData, NewTodoRepo)
```

## 修改server代码

grpc和http，修改服务入参的代码

提供依赖

## 依赖注入

cmd/xxx文件

1. **安装工具**

```go
# 导入到项目中
go get -u github.com/google/wire

# 安装命令
go install github.com/google/wire/cmd/wire

cd cmd -XXX wire
```

控制反转（Inversion of Control，缩写为IoC），是**面向对象编程中的一种设计原则，可以用来减低计算机代码之间的耦合度**。其中最常见的方式叫做依赖注入（Dependency Injection，简称DI）。依赖注入是生成灵活和松散耦合代码的标准技术，通过明确地向组件提供它们所需要的所有依赖关系。

Go社区中有很多依赖注入框架。比如：Uber的[dig](https://github.com/uber-go/dig)和Facebook的[inject](https://github.com/facebookgo/inject)都使用反射来做运行时依赖注入。

[Wire](https://link.segmentfault.com/?enc=fD0W0IcWPmDhPSWFze1Pwg%3D%3D.JhEtO7p6UW6DvZshpJcIwaqeaKeXMIL%2FEVAB8gqQKHM%3D) 是一个的 Google 开源的依赖注入工具，通过自动生成代码的方式在**编译期**完成依赖注入。

`wire`中有两个核心概念：提供者（provider）和注入器（injector）。

Provider

`Wire`中的提供者就是一个可以产生值的普通函数。

提供者函数必须是可导出的（首字母大写）以便被其他包导入。提供者函数也是可以返回错误的。

提供者函数可以分组为提供者函数集（**provider set**）。使用`wire.NewSet` 函数可以将多个提供者函数添加到一个集合中。如果经常同时使用多个提供者函数，这非常有用。

还可以将其他提供者函数集添加到提供者函数集中。

**Injector**

应用程序中是用一个注入器来连接提供者，注入器就是一个按照依赖顺序调用提供者。

使用 `wire`时，你只需要编写注入器的函数签名，然后 `wire`会生成对应的函数体。

要声明一个注入器函数只需要在函数体中调用`wire.Build`。这个函数的返回值也无关紧要，只要它们的类型正确即可。这些值在生成的代码中将被忽略。假设上面的提供者函数是在一个名为 `wire_demo/demo` 的包中定义的，下面将声明一个注入器来得到一个`Z`。

```go
//go:build wireinject
// +build wireinject
```

与提供者一样，注入器也可以输入参数（然后将其发送给提供者），并且可以返回错误。

`wire.Build`的参数和`wire.NewSet`一样：都是提供者集合。这些就在该注入器的代码生成期间使用的提供者集。

将上面的代码保存到`wire.go`中，文件最上面的`//go:build wireinject` 是必须的（Go 1.18之前的版本使用`// +build wireinject`），它确保`wire.go`不会参与最终的项目编译。

```go
go install github.com/google/wire/cmd/wire@latest
wire
```

参考：[https://www.liwenzhou.com/posts/Go/wire/](https://www.liwenzhou.com/posts/Go/wire/)

## 自定义http响应以及错误响应

```go
opts = []http.ServerOption
-->
type ServerOption func(*Server)
--> 
func NewServer(opts ...ServerOption)
-->
enc:         DefaultResponseEncoder,
-->
// DefaultResponseEncoder encodes the object to the HTTP response.
func DefaultResponseEncoder(w http.ResponseWriter, r *http.Request, v interface{}) error {
	if v == nil {
		return nil
	}
	if rd, ok := v.(Redirector); ok {
		url, code := rd.Redirect()
		http.Redirect(w, r, url, code)
		return nil
	}
	codec, _ := CodecForRequest(r, "Accept")
	data, err := codec.Marshal(v)
	if err != nil {
		return err
	}
	w.Header().Set("Content-Type", httputil.ContentType(codec.Name()))
	_, err = w.Write(data)
	if err != nil {
		return err
	}
	return nil
}

// DefaultErrorEncoder encodes the error to the HTTP response.
func DefaultErrorEncoder(w http.ResponseWriter, r *http.Request, err error) {
	se := errors.FromError(err)
	codec, _ := CodecForRequest(r, "Accept")
	body, err := codec.Marshal(se)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", httputil.ContentType(codec.Name()))
	w.WriteHeader(int(se.Code))
	_, _ = w.Write(body)
}

```

### 自定义

```go
type httpResonse struct {
	Code int `json:"code"`
	Msg  string `json:"msg"`
	Data any    `json:"data"`
}

func responseEncoder(w http.ResponseWriter, r *http.Request, v interface{}) error {
	if v == nil {
		return nil
	}
	if rd, ok := v.(kratoshttp.Redirector); ok {
		url, code := rd.Redirect()
		http.Redirect(w, r, url, code)
		return nil
	}
	codec, _ := kratoshttp.CodecForRequest(r, "Accept")

	**// 构造自定义结构体
	resp := &httpResonse{
		Code: http.StatusOK,
		Msg: "success",
		Data: v,
	}**

	data, err := codec.Marshal(resp)
	if err != nil {
		return err
	}
	w.Header().Set("Content-Type", "application/" + codec.Name())
	_, err = w.Write(data)
	if err != nil {
		return err
	}
	return nil
}

	// 替换默认到http响应编码器
	opts = append(opts, http.ResponseEncoder(responseEncoder))
```

### 自定义错误响应

```go
// DefaultErrorEncoder encodes the error to the HTTP response.
func errorEncoder(w http.ResponseWriter, r *http.Request, err error) {
	if err == nil {
		return
	}
	// 判断err是否已经是一个err类型
	// se := errors.FromError(err)
	resp := new(httpResonse)
	// 检查 err 是否是 gRPC 错误，从错误中提取出 gRPC 错误信息
	if gs, ok := status.FromError(err); ok{
		resp = &httpResonse{
			// httpstatus.FromGRPCCode 将其转换为 HTTP 状态码
			Code: kratosstatus.FromGRPCCode(gs.Code()),
			Msg: gs.Message(),
			Data: nil,
		}
	}else {
		resp = &httpResonse{
			Code: http.StatusInternalServerError, //500
			Msg: "内部错误",
			Data: nil,
		}
	}
	codec, _ := kratoshttp.CodecForRequest(r, "Accept")
	body, err := codec.Marshal(resp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/" + codec.Name())
	w.WriteHeader(int(resp.Code))
	_, _ = w.Write(body)
}

	// 替换默认到错误编码响应
	opts = append(opts, http.ErrorEncoder(errorEncoder))
```

### api自定义错误，proto

```go
enum ErrorReason {
     // 设置缺省错误码
     option (errors.default_code) = 500;
  
     // 为某个枚举单独设置错误码
     TODO_NOT_FOUND = 0 [(errors.code) = 404];
   
     INVALID_PARAM = 1 [(errors.code) = 400];
}

// 调用
return &pb.GetTodoReply{}, pb.ErrorTodoNotFound("id:%v todo is not found", req.Id)
```

## 日志

```go
f, err := os.OpenFile("test.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return
	}
	// zap 日志库
	writeSyncer := zapcore.AddSync(f)

	encoder := zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
	core := zapcore.NewCore(encoder, writeSyncer, zapcore.DebugLevel)
	z := zap.New(core)
	// logger := log.With(log.NewStdLogger(os.Stdout),
	// 输出到日志文件中
	// logger := log.With(log.NewStdLogger(f),
	logger := log.With(kratoszap.NewLogger(z),
		"ts", log.DefaultTimestamp,
		"caller", log.DefaultCaller,
		"service.id", id,
		"service.name", Name,
		"service.version", Version,
		"trace.id", tracing.TraceID(),
		"span.id", tracing.SpanID(),
	)
	// Helper
	h := log.NewHelper(
		log.NewFilter(
			logger,
			// 日志过滤，按key过滤
			log.FilterKey("password"),
		),
	)
```

使用

```go
type TodoService struct {
	pb.UnimplementedTodoServer

	uc *biz.TodoUsecase
	log *log.Helper
}

func NewTodoService(uc *biz.TodoUsecase, logger log.Logger) *TodoService {
	return &TodoService{
		uc: uc,
		log: log.NewHelper(logger),
	}
}

s.log.WithContext(ctx).Errorw("uc.CreateTodo failed", err.Error())
```

## 中间件

一个请求进入时的处理顺序为 Middleware 注册的顺序，而响应返回的处理顺序为注册顺序的倒序，即先进后出(FILO)。

```go
// Handler defines the handler invoked by Middleware.
type Handler func(ctx context.Context, req any) (any, error)

// Middleware is HTTP/gRPC transport middleware.
type Middleware func(Handler) Handler
```

自定义中间件

```go
func Middleware1() middleware.Middleware{
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req any) (reply any, err error) {
			// 执行前
			fmt.Println("middle start")
			// 做token校验
			if tr, ok := transport.FromServerContext(ctx); ok{
				token := tr.RequestHeader().Get("token")
				fmt.Printf("toerkn:%v", token)
			}
			defer func() {
				fmt.Println("middle end")
			}()
			return handler(ctx, req)
		}
	}
}

var opts = []http.ServerOption{
		http.Middleware(
			recovery.Recovery(), // 全局中间件
			selector.Server(Middleware1(), // 特定path的中间件
		).
			Path("/api.bubble.v1.Todo/CreateTodo").Build(),
		),
	}
```
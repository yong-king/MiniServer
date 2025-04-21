# go-zero

# 常用

```go
// go api
$ goctl api go --help
Generate go files for provided api in api file

Usage:
  goctl api go [flags]

Flags:
      --api string      The api file
      --branch string   The branch of the remote repo, it does work with --remote
      --dir string      The target dir
  -h, --help            help for go
      --home string     The goctl home path of the template, --home and --remote cannot be set at the same time, if they are, --remote has higher priority
      --remote string   The remote git repo of the template, --home and --remote cannot be set at the same time, if they are, --remote has higher priority
                        The git repo directory must be consistent with the https://github.com/zeromicro/go-zero-template directory structure
      --style string    The file naming format, see [https://github.com/zeromicro/go-zero/blob/master/tools/goctl/config/readme.md] (default "gozero")
```

```go
// go rpc
$ goctl rpc protoc --help
Generate grpc code

Usage:
  goctl rpc protoc [flags]

      --branch string     The branch of the remote repo, it does work with --remote
  -c, --client            Whether to generate rpc client (default true)
  -h, --help              help for protoc
      --home string       The goctl home path of the template, --home and --remote cannot be set at the same time, if they are, --remote has higher
 priority
  -m, --multiple          Generated in multiple rpc service mode
      --remote string     The remote git repo of the template, --home and --remote cannot be set at the same time, if they are, --remote has higher
 priority
                          The git repo directory must be consistent with the https://github.com/zeromicro/go-zero-template directory structure     
      --style string      The file naming format, see [https://github.com/zeromicro/go-zero/blob/master/tools/goctl/config/readme.md]
  -v, --verbose           Enable log output
      --zrpc_out string   The zrpc output directory
      
# 单个 rpc 服务生成示例指令
$ goctl rpc protoc greet.proto --go_out=./pb --go-grpc_out=./pb --zrpc_out=. --client=true 
# 多个 rpc 服务生成示例指令
$ goctl rpc protoc greet.proto --go_out=./pb --go-grpc_out=./pb --zrpc_out=. --client=true -m
```

```go
// go model
$ goctl model mysql ddl --help
Generate mysql model from ddl

Usage:
  goctl model mysql ddl [flags]

Flags:
      --branch string     The branch of the remote repo, it does work with --remote
  -c, --cache             Generate code with cache [optional]
      --database string   The name of database [optional]
  -d, --dir string        The target dir
  -h, --help              help for ddl
      --home string       The goctl home path of the template, --home and --remote cannot be set at the same time, if they are, --remote has higher priority
      --idea              For idea plugin [optional]
      --remote string     The remote git repo of the template, --home and --remote cannot be set at the same time, if they are, --remote has higher priority
                          The git repo directory must be consistent with the https://github.com/zeromicro/go-zero-template directory structure
  -s, --src string        The path or path globbing patterns of the ddl
      --style string      The file naming format, see [https://github.com/zeromicro/go-zero/tree/master/tools/goctl/config/readme.md]

Global Flags:
  -i, --ignore-columns strings   Ignore columns while creating or updating rows (default [create_at,created_at,create_time,update_at,updated_at,update_time])
      --strict                   Generate model in strict mode
      -p, --prefix string            The cache prefix, effective when --cache is true (default "cache")
```

## API

```go
// 示例
syntax = "v1"

info (
    title:   "api 文件完整示例写法"
    desc:    "演示如何编写 api 文件"
    author:  "keson.an"
    date:    "2022 年 12 月 26 日"
    version: "v1"
)

type UpdateReq {
    Arg1 string `json:"arg1"`
}

type ListItem {
    Value1 string `json:"value1"`
}

type LoginReq {
    Username string `json:"username"`
    Password string `json:"password"`
}

type LoginResp {
    Name string `json:"name"`
}

type FormExampleReq {
    Name string `form:"name"`
}

type PathExampleReq {
    // path 标签修饰的 id 必须与请求路由中的片段对应，如
    // id 在 service 语法块的请求路径上一定会有 :id 对应，见下文。
    ID string `path:"id"`
}

type PathExampleResp {
    Name string `json:"name"`
}

@server (
    jwt:        Auth // 对当前 Foo 语法块下的所有路由，开启 jwt 认证，不需要则请删除此行
    prefix:     /v1 // 对当前 Foo 语法块下的所有路由，新增 /v1 路由前缀，不需要则请删除此行
    group:      g1 // 对当前 Foo 语法块下的所有路由，路由归并到 g1 目录下，不需要则请删除此行
    timeout:    3s // 对当前 Foo 语法块下的所有路由进行超时配置，不需要则请删除此行
    middleware: AuthInterceptor // 对当前 Foo 语法块下的所有路由添加中间件，不需要则请删除此行
    maxBytes:   1048576 // 对当前 Foo 语法块下的所有路由添加请求体大小控制，单位为 byte,goctl 版本 >= 1.5.0 才支持
)
service Foo {
    // 定义没有请求体和响应体的接口，如 ping
    @handler ping
    get /ping

    // 定义只有请求体的接口，如更新信息
    @handler update
    post /update (UpdateReq)

    // 定义只有响应体的结构，如获取全部信息列表
    @handler list
    get /list returns ([]ListItem)

    // 定义有结构体和响应体的接口，如登录
    @handler login
    post /login (LoginReq) returns (LoginResp)

    // 定义表单请求
    @handler formExample
    post /form/example (FormExampleReq)

    // 定义 path 参数
    @handler pathExample
    get /path/example/:id (PathExampleReq) returns (PathExampleResp)
}

goctl api go -api user.api -dir . -style=goZero
```

## 语法

单行注释以 `//` 开始，行尾结束。

多行注释（文档注释）以 `/*` 开始，以第一个 `*/` 结束。

**字符串**

原始字符串的字符序列在两个反引号之间，除反引号外，任何字符都可以出现，如 `foo`；

普通字符串的字符序列在两个双引号之间，除双引号外，任何字符都可以出现，如 "foo"。

在 api 语言中，双引号字符串不支持 `\"` 来实现字符串转义。

**syntax 语句**

syntax 语句用于标记 api 语言的版本，不同的版本可能语法结构有所不同，随着版本的提升会做不断的优化

```go
syntax = "v1"
```

**info 语句**

info 语句是 api 语言的 meta 信息，其仅用于对当前 api 文件进行描述，**暂**不参与代码生成，其和注释还是有一些区别，注释一般是依附某个 syntax 语句存在，而 info 语句是用于描述整个 api 信息的，当然，不排除在将来会参与到代码生成里面来

**import 语句**

`import` 语句是在 api 中引入其他 api 文件的语法块，其支持相对/绝对路径，**不支持** `package` 的设计

```go
// 单行 import
import "foo"
import "/path/to/file"

// import 组
import ()
import (
    "bar"
    "relative/to/file"
)
```

**数据类型**

api 中的数据类型基本沿用了 Golang 的数据类型，用于对 rest 服务的请求/响应体结构的描述，

```go
// 空结构体
type Foo {}

// 单个结构体
type Bar {
    Foo int               `json:"foo"`
    Bar bool              `json:"bar"`
    Baz []string          `json:"baz"`
    Qux map[string]string `json:"qux"`
}

type Baz {
    Bar    `json:"baz"`
    Array [3]int `json:"array"`
    // 结构体内嵌 goctl 1.6.8 版本支持
    Qux {
        Foo string `json:"foo"`
        Bar bool   `json:"bar"`
    } `json:"baz"`
}

// 空结构体组
type ()

// 结构体组
type (
    Int int
    Integer = int
    Bar {
        Foo int               `json:"foo"`
        Bar bool              `json:"bar"`
        Baz []string          `json:"baz"`
        Qux map[string]string `json:"qux"`
    }
)
// !!! 不支持 package 设计，如 time.Time。
```

**service 语句**

service 语句是对 HTTP 服务的直观描述，包含请求 handler，请求方法，请求路由，请求体，响应体，jwt 开关，中间件声明等定义。

**@server 语句**

@server 语句是对一个服务语句的 meta 信息描述，其对应特性包含但不限于：

- jwt 开关
- 中间件
- 路由分组
- 路由前缀

```go
// 空内容
@server()

// 有内容
@server (
    // jwt 声明
    // 如果 key 固定为 “jwt:”，则代表开启 jwt 鉴权声明
    // value 则为配置文件的结构体名称
    jwt: Auth

    // 路由前缀
    // 如果 key 固定为 “prefix:”
    // 则代表路由前缀声明，value 则为具体的路由前缀值，字符串中没让必须以 / 开头
    prefix: /v1

    // 路由分组
    // 如果 key 固定为 “group:”，则代表路由分组声明
    // value 则为具体分组名称，在 goctl生成代码后会根据此值进行文件夹分组
    group: Foo

    // 中间件
    // 如果 key 固定为 middleware:”，则代表中间件声明
    // value 则为具体中间件函数名称，在 goctl生成代码后会根据此值进生成对应的中间件函数
    middleware: AuthInterceptor

    // 超时控制
    // 如果 key 固定为  timeout:”，则代表超时配置
    // value 则为具体中duration，在 goctl生成代码后会根据此值进生成对应的超时配置
    timeout: 3s

    // 其他 key-value，除上述几个内置 key 外，其他 key-value
    // 也可以在作为 annotation 信息传递给 goctl 及其插件，但就
    // 目前来看，goctl 并未使用。
    foo: bar
)
```

## mysql配置及model操作

1. sql语句
2. 创建数据库表
    1. sql中的唯一索引会生成相对应的查询方法 ！！！
3. goctl model mysql datasource 指令用于从数据库连接生成 model 代码。

```go
$ goctl model mysql datasource --help
Generate model from datasource

Usage:
  goctl model mysql datasource [flags]

Flags:
      --branch string   The branch of the remote repo, it does work with --remote
  -c, --cache           Generate code with cache [optional]
  -d, --dir string      The target dir
  -h, --help            help for datasource
      --home string     The goctl home path of the template, --home and --remote cannot be set at the same time, if they are, --remote has higher priority
      --idea            For idea plugin [optional]
      --remote string   The remote git repo of the template, --home and --remote cannot be set at the same time, if they are, --remote has higher priority
                        The git repo directory must be consistent with the https://github.com/zeromicro/go-zero-template directory structure
      --style string    The file naming format, see [https://github.com/zeromicro/go-zero/tree/master/tools/goctl/config/readme.md]
  -t, --table strings   The table or table globbing patterns in the database
      --url string      The data source of database,like "root:password@tcp(127.0.0.1:3306)/database"

Global Flags:
  -i, --ignore-columns strings   Ignore columns while creating or updating rows (default [create_at,created_at,create_time,update_at,updated_at,update_time])
      --strict                   Generate model in strict mode
      -p, --prefix string            The cache prefix, effective when --cache is true (default "cache")
  
  goctl model mysql datasource -url="root:password@tcp(127.0.0.1:3306)/database" -table="*" -dir="./model"
```

生成

```go
.
├── usermodel.go
├── usermodel_gen.go
└── vars.go
```

1. 修改config 以及配置文件 

```go
type Config struct {
	rest.RestConf

	***MysqlDb struct{
		DbSource string `json:"DbSource"`
	}***
}

MysqlDB:
  DbSource: root:password@tcp(127.0.0.1:3306)/database
```

1. 添加调用信息

<aside>
💡

这里是我的理解

ctx := svc.NewServiceContext(c) // 上下文信息

mysql要在其中传递

1. model.UserModel 是一个接口，在其中的userModel接口需要实现增删改成的方法
2. 这些方法由defaultUserModel实现，即defaultUserModel实现了这个接口
3. newUserModel 是 defaultUserModel的构造方法

```go
customUserModel struct {
		*defaultUserModel
	}
// NewUserModel returns a model for the database table.
func NewUserModel(conn sqlx.SqlConn) UserModel {
	return &customUserModel{
		defaultUserModel: newUserModel(conn),
	}
}
```

1. customUserModel 定义了 *defaultUserModel， NewUserModel 是 customUserModel构造方法，因此做以下修改

```go
type ServiceContext struct {
	Config config.Config
	//UserModel: 类型为 model.UsersModel，表示与用户相关的数据库模型
	//用于处理与用户相关的数据操作（如用户的创建、读取、更新和删除等）
	UserModel model.UserModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,

		// UserModel 指针类型 --> userModel 指针类型
		// userModel 包含增删改查
		// defaultUserModel 实现了userModel接口
		// newUserModel 是 defaultUserModel 的构建方法
		//  NewUserModel(conn sqlx.SqlConn) UserModel 
		//通过调用 model.NewUsersModel 函数对UserModel 进行初始化
		//sqlx.NewMysql 是数据库连接,链接字符串为config中的MysqlDb.DbSource
		UserModel: model.NewUserModel(sqlx.NewMysql(c.MysqlDb.DbSource)),
	}
}

```

</aside>

1. 调用

```go
ctx := svc.NewServiceContext(c) // 上下文，cc.MysqlDb.DbSource
// ctx.UserModel

1. handler.RegisterHandlers(server, ctx)
2. Handler: SignupHandler(serverCtx)
3. l := logic.NewSignupLogic(r.Context(), svcCtx)
	resp, err := l.Signup(&req)
4. Singup是我们的业务逻辑部分
```

1. 业务逻辑

```go
_, err = l.svcCtx.UserModel.Insert(context.Background(), user)
```

### 配置cache

1. 生成model文件

```go
goctl model mysql datasource -url="root:password@tcp(127.0.0.1:3306)/database" -table="*" -dir="./model" -c
.
├── usermodel.go
├── usermodel_gen.go
└── vars.go

```

1. 修改配置文件

```go
// api-yaml
CacheRedis:
  - Host: 127.0.0.1:6379
    Pass: password // 可以省略，但前面不要加-
// config.go
CacheRedis cache.CacheConf
```

1. 修改model文件

```go
// model.go
1. 这里会变为
func newUserModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) *defaultUserModel {
	return &defaultUserModel{
		CachedConn: sqlc.NewConn(conn, c, opts...),
		table:      "`user`",
	}
}
-->
type CacheConf = ClusterConf
-->
type (
	// A ClusterConf is the config of a redis cluster that used as cache.
	ClusterConf []NodeConf

	// A NodeConf is the config of a redis node that used as cache.
	NodeConf struct {
		redis.RedisConf
		Weight int `json:",default=100"`
	}
)
-->
type (
	// A RedisConf is a redis config.
	RedisConf struct {
		Host     string
		Type     string `json:",default=node,options=node|cluster"`
		User     string `json:",optional"`
		Pass     string `json:",optional"`
		Tls      bool   `json:",optional"`
		NonBlock bool   `json:",default=true"`
		// PingTimeout is the timeout for ping redis.
		PingTimeout time.Duration `json:",default=1s"`
	}
2. mdeol.go 修改
// NewUserModel returns a model for the database table.
func NewUserModel(conn sqlx.SqlConn, **c cache.CacheConf**) UserModel {
	return &customUserModel{
		defaultUserModel: newUserModel(conn, **c)**,
	}
}

func (m *customUserModel) withSession(session sqlx.Session, **c cache.CacheConf**) UserModel {
	return NewUserModel(sqlx.NewSqlConnFromSession(session), **c)**
}
```

1. 修改serviceContext

```go
model.NewUserModel(sqlx.NewMysql(c.MysqlDb.DbSource), **c.CacheRedis**),
```

### 配置日志

1. 修改配置文件

```go
1. 找到配置文件中关于日志的信息
 // config.go 
 rest.RestConf
 -->
 service.ServiceConf
 -->
 Log        logx.LogConf
 --> 
 type LogConf struct {
	// ServiceName represents the service name.
	ServiceName string `json:",optional"`
	// Mode represents the logging mode, default is `console`.
	// console: log to console.
	// file: log to file.
	// volume: used in k8s, prepend the hostname to the log file name.
	Mode string `json:",default=console,options=[console,file,volume]"`
	// Encoding represents the encoding type, default is `json`.
	// json: json encoding.
	// plain: plain text encoding, typically used in development.
	Encoding string `json:",default=json,options=[json,plain]"`
	// TimeFormat represents the time format, default is `2006-01-02T15:04:05.000Z07:00`.
	TimeFormat string `json:",optional"`
	// Path represents the log file path, default is `logs`.
	Path string `json:",default=logs"`
	// Level represents the log level, default is `info`.
	Level string `json:",default=info,options=[debug,info,error,severe]"`
	// MaxContentLength represents the max content bytes, default is no limit.
	MaxContentLength uint32 `json:",optional"`
	// Compress represents whether to compress the log file, default is `false`.
	Compress bool `json:",optional"`
	// Stat represents whether to log statistics, default is `true`.
	Stat bool `json:",default=true"`
	// KeepDays represents how many days the log files will be kept. Default to keep all files.
	// Only take effect when Mode is `file` or `volume`, both work when Rotation is `daily` or `size`.
	KeepDays int `json:",optional"`
	// StackCooldownMillis represents the cooldown time for stack logging, default is 100ms.
	StackCooldownMillis int `json:",default=100"`
	// MaxBackups represents how many backup log files will be kept. 0 means all files will be kept forever.
	// Only take effect when RotationRuleType is `size`.
	// Even though `MaxBackups` sets 0, log files will still be removed
	// if the `KeepDays` limitation is reached.
	MaxBackups int `json:",default=0"`
	// MaxSize represents how much space the writing log file takes up. 0 means no limit. The unit is `MB`.
	// Only take effect when RotationRuleType is `size`
	MaxSize int `json:",default=0"`
	// Rotation represents the type of log rotation rule. Default is `daily`.
	// daily: daily rotation.
	// size: size limited rotation.
	Rotation string `json:",default=daily,options=[daily,size]"`
	// FileTimeFormat represents the time format for file name, default is `2006-01-02T15:04:05.000Z07:00`.
	FileTimeFormat string `json:",optional"`
}
2. 修改配置文件
// api.yaml
rest.RestConf ->  service.ServiceConf -> Log  logx.LogConf

Log:
  ServiceName: mall
  Mode: file
  Encoding: plain
  Path: logs
  Level: debug
  Stat: true
// 按自己的需求配置
```

1. 在业务中添加日志逻辑
2. 日志类型
    1. logc 是对 logx 的封装，可以带上 context 进行日志打印
        
        ```go
        logx.WithContext(ctx).Info("hello world")
        logc.Info(ctx, "hello world")
        // 代码是等效的
        ```
        
    
    ```go
    type Logger interface {
        Debug(...any)
        Debugf(string, ...any)
        Debugv(any)
        Debugw(string, ...LogField)
        Error(...any)
        Errorf(string, ...any)
        Errorv(any)
        Errorw(string, ...LogField)
        Info(...any)
        Infof(string, ...any)
        Infov(any)
        Infow(string, ...LogField)
        Slow(...any)
        Slowf(string, ...any)
        Slowv(any)
        Sloww(string, ...LogField)
        WithCallerSkip(skip int) Logger
        WithContext(ctx context.Context) Logger
        WithDuration(d time.Duration) Logger
        WithFields(fields ...LogField) Logger
    }
    // https://go-zero.dev/docs/components/logx
    ```
    

### JWT(Json Web Token)

1. 生成 JWT 方法

```go
// 生成JWT方法
// @secretKey: JWT 加解密密钥
// @iat: 时间戳
// @seconds: 过期时间，单位秒
// @payload: 数据载体
func (l *LoginLogic)getJwtToken(secretKey string, iat, seconds int64, userId int64) (string, error) {
	claims := make(jwt.MapClaims)
	claims["exp"] = iat + seconds
	claims["iat"] = iat
	claims["userId"] = userId
	claims["auth"] = "ysh"
	token := jwt.New(jwt.SigningMethodHS256)
	token.Claims = claims
	return token.SignedString([]byte(secretKey))
  }
```

1. 配置文件的修改

```go
// yaml
Auth:
  AccessSecret: dlrb&zrn&ysh
  AccessExpire: 60
// config.go
	Auth struct {// JWT 认证需要的密钥和过期时间配置
        AccessSecret string
        AccessExpire int64
    }
// api
@server (
	prefix: /v1
	jwt:    Auth // 开启 jwt 认证
)
```

1. 加入需要JWT的逻辑

```go
token, err := l.getJwtToken(l.svcCtx.Config.Auth.AccessSecret, now, expire, user.UserId)
```

1. 获取JWT信息

```go
value:=l.ctx.Value("custom-key")
```

### 自定义中间件

1. 方式一
    
    ```go
    1. 修改api文件
    @server (
    	prefix:     /v1
    	jwt:        Auth // 开启 jwt 认证
    	middleware: Cost // 添加中间件（路由中间件）
    )
    
    2. 修改serviceContext.go
    type ServiceContext struct {
    	Config config.Config
    	Cost rest.Middleware
    }
    
    func NewServiceContext(c config.Config) *ServiceContext {
    	return &ServiceContext{
    		Config: c,
    		Cost: middleware.NewCostMiddleware().Handle,
    	}
    }
    
    3. 中间件逻辑 costMiddleware.go
    func (m *CostMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
    	return func(w http.ResponseWriter, r *http.Request) {
    		// TODO generate middleware implement function, delete after code implementation
    		// 中间件逻辑
    		now := time.Now()
    		// Passthrough to next handler if need
    		next(w, r) // 实际执行的后续接口handler处理函数
    		logx.Infof("-->cost:%v\n", time.Since(now))
    	}
    }
    ```
    
2. 方案二
    
    ```go
    // rest.Middleware --> Middleware func(next http.HandlerFunc) http.HandlerFunc
    // type HandlerFunc func(ResponseWriter, *Request)
    1. 前提
    1.1 结构体
    	type bodyCopy struct{
    		http.ResponseWriter	// 结构体嵌入接口类型，默认实现了接口的所有方法
    		body *bytes.Buffer // 我们记录响应体的内容
    	}
    1.2 重写Write 
    func (bc bodyCopy) Write(b []byte) (int, error) {
    	// 1. 先记录到我们的这里
    	bc.body.Write(b)
    	// 2. 再往HTTP响应体写内容
    	return bc.ResponseWriter.Write(b)
    }
    1.3 结构体构造方法
    func NewbodyCopy(w http.ResponseWriter) *bodyCopy {
    	return &bodyCopy{
    		ResponseWriter: w,
    		body: bytes.NewBuffer([]byte{}),
    	}
    }
    2. 中间件逻辑
    func CopyResq(next http.HandlerFunc) http.HandlerFunc {
    	return func(w http.ResponseWriter, r *http.Request) {
    		// 初始化一个自定义的 ResponseWriter.Write
    		bc := NewbodyCopy(w)
    		// 实际执行完后，会执行bc.body.Write(b)， 然后再往HTTP响应体写内容
    		next(bc, r)
    		// 处理后的请求
    		logx.Infof("-->reqL%v resp:%v\n", r.URL, bc.body.String())
    	}
    }
    3. 使用中间件 user.go
    	server.Use(middleware.CopyResq)
    	
    // 使用其他中间，还有一种，基于闭包，修改其中内容
    func MiddlewareWithAnotherService(ok bool) rest.Middleware {
    	return func(next http.HandlerFunc) http.HandlerFunc {
    		return func(w http.ResponseWriter, r *http.Request) {
    			if ok {
    				fmt.Println("ok!")
    			}
    			next(w, r)
    		}
    	}
    }
    
    server.Use(middleware.MiddlewareWithAnotherService(true))
    ```
    

## GRPC

1. proto文件

```go
$ goctl rpc protoc --help
Generate grpc code

Usage:
  goctl rpc protoc [flags]

      --branch string     The branch of the remote repo, it does work with --remote
  -c, --client            Whether to generate rpc client (default true)
  -h, --help              help for protoc
      --home string       The goctl home path of the template, --home and --remote cannot be set at the same time, if they are, --remote has higher
 priority
  -m, --multiple          Generated in multiple rpc service mode
      --remote string     The remote git repo of the template, --home and --remote cannot be set at the same time, if they are, --remote has higher
 priority
                          The git repo directory must be consistent with the https://github.com/zeromicro/go-zero-template directory structure     
      --style string      The file naming format, see [https://github.com/zeromicro/go-zero/blob/master/tools/goctl/config/readme.md]
  -v, --verbose           Enable log output
      --zrpc_out string   The zrpc output directory
      
goctl rpc protoc greet.proto --go_out=./pb --go-grpc_out=./pb --zrpc_out=.

.
├── etc
│   └── user.yaml
├── internal
│   ├── config
│   │   └── config.go
│   ├── logic
│   │   └── getuserlogic.go
│   ├── server
│   │   └── userserver.go
│   └── svc
│       └── servicecontext.go
├── pb
│   └── user
│       ├── user.pb.go
│       └── user_grpc.pb.go
├── user.go
├── user.proto
└── userclient
    └── user.go

```

1. 修改配置文件

```go
// yaml
Name: user.rpc
ListenOn: 0.0.0.0:8080
Mode: dev
Etcd:
  Hosts:
  - 127.0.0.1:2379
  Key: user.rpc

MysqlDB:
  DbSource: root:password@tcp(127.0.0.1:3307)/database?parseTime=true

CacheRedis:
  - Host: 127.0.0.1:6379
    Pass: password
// config.go
type Config struct {
	zrpc.RpcServerConf

	// mysql
	MysqlDb struct{
		DbSource string `json:"DbSource"`
	}

	// redis
	CacheRedis cache.CacheConf
}

// Mode: dev 用于调试grpc
// rpc服务测试工具
// 一个测试grpc服务的ui工具 https://github.com/fullstorydev/grpcui
// 安装
go install github.com/fullstorydev/grpcui/cmd/grpcui@latest
grpcui -plaintext localhost:8080

// etcd  hub.docker.com
https://hub.docker.com/r/bitnami/etcd

// 拉取镜像
docker pull bitnami/etcd

1. docker network create app-tier --driver bridge
2. docker run -d --name Etcd-server \
    --network app-tier \
    --publish 2379:2379 \
    --publish 2380:2380 \
    --env ALLOW_NONE_AUTHENTICATION=yes \
    --env ETCD_ADVERTISE_CLIENT_URLS=http://etcd-server:2379 \
    bitnami/etcd:latest
3. docker run -it --rm \
    --network app-tier \
    --env ALLOW_NONE_AUTHENTICATION=yes \
    bitnami/etcd:latest etcdctl --endpoints http://etcd-server:2379 put /message Hello

```

1. 修改serviceContext

```go
type UserServer struct {
	svcCtx *svc.ServiceContext
	user.UnimplementedUserServer
}

func NewUserServer(svcCtx *svc.ServiceContext) *UserServer {
	return &UserServer{
		svcCtx: svcCtx,
	}
}

func (s *UserServer) GetUser(ctx context.Context, in *user.GetUserReq) (*user.GetUserResp, error) {
	l := logic.NewGetUserLogic(ctx, s.svcCtx)
	return l.GetUser(in)
}

```

1. 业务逻辑

```go
// logic....go
func (l *GetUserLogic) GetUser(in *user.GetUserReq) (*user.GetUserResp, error) {
```

### 调用rpc服务

需求，当一个服务需要用到另一个服务时，可以通过rpc调用另一个服务

1. 当前服务配置
    1. 数据库配置， 这里记住，你的sql语句中的唯一索引会生成相对应的方法。修改配置文件和配置服务文件。
    
    ```go
    goctl model mysql datasource -url="root:password@tcp(127.0.0.1:3306)/database" -table="*" -dir="./model" -c
    ```
    
    b. api文件
    
    ```go
    goctl api go -api user.api -dir . -style=goZero
    ```
    
    c. 修改配置文件，用于rpc调用
    
    ```go
    // config.go
    UserRPC zrpc.RpcClientConf	// 连接其他微服务的RPC客户端
    
    // yaml
    UserRPC:
      Etcd:
        Hosts: 
          - 127.0.0.1:2379
        Key: user.rpc
       
    // serviceContext.go
    type ServiceContext struct {
    	Config config.Config
    	UserRPC userclient.User
    }
    
    func NewServiceContext(c config.Config) *ServiceContext {
    	return &ServiceContext{
    		Config: c,
    		UserRPC: userclient.NewUser(zrpc.MustNewClient(c.UserRPC)),
    	}
    }
    ```
    
    d. 业务逻辑
    
    ```go
    // rpc调用
    l.svcCtx.UserRPC.GetUser
    ```
    
    e. 另一个rpc服务必须的跑起来！！
    

### 使用consul

服务注册

1. 修改配置文件
    
    ```go
    // config.go
    go get -u github.com/zeromicro/zero-contrib/zrpc/registry/consul
    type Config struct {
    	zrpc.RpcServerConf
    	Consul consul.Conf
    }
    
    // yaml 
    // 1.注释掉Etcd相关
    // 2.添加consul相关
    Consul:
      Host: 127.0.0.1:8500
      Key: consul-user.rpc
    ```
    
2. 启动服务注册到consul

```go
// api.go	
	// 注册consul
	_ = consul.RegisterService(c.ListenOn, c.Consul)
```

服务发现

1. 修改配置文件

```go
// yaml
UserRPC:
  Target: consul://127.0.0.1:8500/consul-user.rpc?wait=14s
// 注释掉etct相关
```

1. 启动导入

```go
// api.go
	_ "github.com/zeromicro/zero-contrib/zrpc/registry/consul"
```

### RPC拦截器和metadata

元数据（[metadata](https://pkg.go.dev/google.golang.org/grpc/metadata)）是指在处理RPC请求和响应过程中需要但又不属于具体业务（例如身份验证详细信息）的信息，采用键值对列表的形式，其中键是`string`类型，值通常是`[]string`类型，但也可以是二进制数据。gRPC中的 metadata 类似于我们在 HTTP headers中的键值对，元数据可以包含认证token、请求标识和监控标签等。

```go
md := metadata.New(map[string]string{"key1": "val1", "key2": "val2"})

md := metadata.Pairs(
    "key1", "val1",
    "key1", "val1-2", // "key1"的值将会是 []string{"val1", "val1-2"}
    "key2", "val2",
)

// 从请求上下文中获取元数据
metadata.FromIncomingContext(ctx)

// 发送metadata
// 创建带有metadata的context
md := metadata.Pairs("k1", "v1", "k1", "v2", "k2", "v3")
ctx := metadata.NewOutgoingContext(context.Background(), md)
```

**拦截器（中间件）**

gRPC 为在每个 ClientConn/Server 基础上实现和安装拦截器提供了一些简单的 API。 拦截器拦截每个 RPC 调用的执行。用户可以使用拦截器进行日志记录、身份验证/授权、指标收集以及许多其他可以跨 RPC 共享的功能。

**客户端端拦截器**

[UnaryClientInterceptor](https://godoc.org/google.golang.org/grpc#UnaryClientInterceptor) 是客户端一元拦截器的类型

```go
func(ctx context.Context, method string, req, reply interface{}, cc *ClientConn, invoker UnaryInvoker, opts ...CallOption) error

```

- 预处理：用户可以通过检查传入的参数(如 RPC 上下文、方法字符串、要发送的请求和 CallOptions 配置)来获得有关当前 RPC 调用的信息。
- RPC调用：预处理完成后，可以通过执行`invoker`执行 RPC 调用。
- 调用后：一旦调用者返回应答和错误，用户就可以对 RPC 调用进行后处理。通常，它是关于处理返回的响应和错误的。 若要在 `ClientConn` 上安装一元拦截器，请使用`DialOptionWithUnaryInterceptor`的`DialOption`配置 Dial 。

**server端拦截器**

```go
func(ctx context.Context, req interface{}, info *UnaryServerInfo, handler UnaryHandler) (resp interface{}, err error)

```

1. 客户端拦截器（知识）

```go
//serviceContecxt.go
UserRPC: userclient.NewUser(zrpc.MustNewClient(c.UserRPC))
--> 
func MustNewClient(c RpcClientConf, options ...ClientOption) Client {
	cli, err := NewClient(c, options...)
	logx.Must(err)
	return cli
}
```

1. 拦截器逻辑

```go
func Interceptor(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	fmt.Println("客户端拦截器 in")
	// RPC调用前
	// 编写客户端拦截器的逻辑
	adminID := ctx.Value(CtxKeyAdmindID).(string)
	md := metadata.Pairs(
		"token", "ysh&dlrb",
		"adminID", **adminID,** 
	)
	ctx = metadata.NewOutgoingContext(ctx, md)

	err := invoker(ctx, method, req, reply, cc, opts...) // 实际的RPC调用

	// RPC调用后
	fmt.Println("客户端拦截器 out")
	return err
}
```

<aside>
💡

```go
type CtxKey string
const(
	CtxKeyAdmindID CtxKey = "adminID"
)
// 用这种方法可以避免冲突
```

</aside>

1. 数据的传入（**adminI**）
    
    ```go
    l.ctx = context.WithValue(l.ctx, interceptor.CtxKeyAdmindID, "666") // 在调用rpc之前传入
    ```
    
2. 添加上下文服务配置

```go
func NewServiceContext(c config.Config) *ServiceContext {
	conn := sqlx.NewMysql(c.MysqlDb.DbSource)
	return &ServiceContext{
		Config: c,
		OrderModel: model.NewOrderModel(conn, c.CacheRedis),
		UserRPC: userclient.NewUser(
			zrpc.MustNewClient(
				c.UserRPC, 
				zrpc.WithUnaryClientInterceptor(interceptor.YshInterceptor),
			),
		),
	}
}
```

1. 服务端拦截器（服务启动之前注册）

```go
// 注册服务端拦截器
	s.AddUnaryInterceptors(myInterceptor)
```

1. 拦截器逻辑

```go
func yshInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
	// 调用前
	fmt.Println("服务端拦截器 in")
	// 拦截器业务逻辑
	// 获取元数据
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.InvalidArgument, "need metadata")
	}
	fmt.Println("metadata:%#v\n", md)

	// 根据metadata中的数据进行一些校验处理
	if md["token"][0] != "ysh&dlrb" {
		return nil, status.Errorf(codes.Unauthenticated, "invalid token")
	}

	m, err := handler(ctx, req) // 实际RPC方法

	// 调用后
	fmt.Println("服务端拦截器 out")
	return m, err
}
```

### 错误处理

1. 自定义错误格式

```go
const(
	DefaultErrorCode = 1001
	RpcErroCode = 1002
	SqlErrorCode = 1003
	QuerNoFoundErrorCode = 1004
	RedisErrorCode = 1005
)

// CodeError 自定义错误类型
type CodeError struct{
	Code int `json:"code"`
	Msg string `json:"msg"`
}

// CodeErrorResponse 自定义响应错误类型
type CodeErrorResponse struct{
	Code int `josn:"code"`
	Msg string `json:"msg"`
}

// NewCodeError 返回自定义错误
func NewCodeError(code int, msg string) error {
	return CodeError{
		Code: code,
		Msg: msg,
	}
}
// Error CodeError实现error接口
func (e CodeError) Error() string {
	return e.Msg
}

// NewDefaultCodeError 返回默认自定义错误
func NewDefaultCodeError(msg string) error {
	return CodeError{
		Code: DefaultErrorCode,
		Msg: msg,
	}
}

// Data 返回自定义类型的错误响应
func (e *CodeError) Data() *CodeErrorResponse {
	return &CodeErrorResponse{
		Code: e.Code,
		Msg: e.Msg,
	}
}
```

1. 业务中按需返回自定义的错误

```go
return nil, errorx.NewCodeError(errorx.SqlErrorCode, "内部错误")
```

1. 处理自定义错误

```go
// api.go	
	// 注册自定义错误处理方法
	httpx.SetErrorHandlerCtx(func(cte context.Context, err error)(int, any) {
		switch e := err.(type) {
		case errorx.CodeError: // 自定义错误类型
		return http.StatusOK, e.Data()
		default:
			return http.StatusInternalServerError, nil
		}
	})

```

### 定制模版

```go
//例如：
// 实现统一格式的 body 响应:
{
  "code": 0,
  "msg": "OK",
  "data": {}
  // ①
}
```

**准备工作**

提前在 `module` 为 `greet` 的工程下的 `response` 包中写一个 `Response` 方法

```go
package response

import (
    "net/http"

    "github.com/zeromicro/go-zero/rest/httpx"
)

type Body struct {
    Code int         `json:"code"`
    Msg  string      `json:"msg"`
    Data interface{} `json:"data,omitempty"`
}

func Response(w http.ResponseWriter, resp interface{}, err error) {
    var body Body
    if err != nil {
        body.Code = -1
        body.Msg = err.Error()
    } else {
        body.Msg = "OK"
        body.Data = resp
    }
    httpx.OkJson(w, body)
}
```

**修改 `handler` 模板**

```go
// 在goctl env环境变量下看版本，然后
$ vim ~/.goctl/${goctl版本号}/api/handler.tpl
// 如果本地没有~/.goctl/${goctl版本号}/api/handler.tpl文件，
// 可以通过模板初始化命令goctl template init进行初始化
```

```go
// ① 替换为你真实的response包名，仅供参考

// ② 自定义模板内容
package handler

import (
    "net/http"
    "greet/response"// ①. 导入的是上面的响应的包
    {{.ImportPackages}}
)

func {{.HandlerName}}(svcCtx *svc.ServiceContext) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        {{if .HasRequest}}var req types.{{.RequestType}}
        if err := httpx.Parse(r, &req); err != nil {
            httpx.Error(w, err)
            return
        }{{end}}

        l := {{.LogicName}}.New{{.LogicType}}(r.Context(), svcCtx)
        {{if .HasResp}}resp, {{end}}err := l.{{.Call}}({{if .HasRequest}}&req{{end}})
        **{{if .HasResp}}response.Response(w, resp, err){{else}}response.Response(w, nil, err){{end}}**//②

    }
}
```

```go
goctl api go -api xxx.api -dir . -sytle=gozero // 生成新的版本，要删除以前的返回响应的包

https://go-zero.dev/docs/tutorials/customization/template
// 官方文档

gotcl template clean // 清楚自定义模版
```
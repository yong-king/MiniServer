# 微服务 与 云原生

什么是RPC？

RPC（Remote Procedure Call），即远程过程调用。它允许像调用本地服务一样调用远程服务。
RPC是一种服务器-客户端（Client/Server）模式，经典实现是一个通过发送请求-接受回应进行信息交互的系统。
首先与RPC（远程过程调用）相对应的是本地调用。

# RPC原理与Go RPC

RPC就是为了解决类似远程、跨内存空间、的函数/方法调用的。要实现RPC就需要解决以下三个问题

1. 如何确定要执行的函数？ 在本地调用中，函数主体通过函数指针函数指定，然后调用 add 函数，编译器通过函数指针函数自动确定 add 函数在内存中的位置。但是在 RPC 中，调用不能通过函数指针完成，因为它们的内存地址可能完全不同。因此，调用方和被调用方都需要维护一个{ function <-> ID }映射表，以确保调用正确的函数。

2. 如何表达参数？ 本地过程调用中传递的参数是通过堆栈内存结构实现的，但 RPC 不能直接使用内存传递参数，因此参数或返回值需要在传输期间序列化并转换成字节流，反之亦然。
3. 如何进行网络传输？ 函数的调用方和被调用方通常是通过网络连接的，也就是说，function ID 和序列化字节流需要通过网络传输，因此，只要能够完成传输，调用方和被调用方就不受某个网络协议的限制。.例如，一些 RPC 框架使用 TCP 协议，一些使用 HTTP。

以往实现跨服务调用的时候，我们会采用RESTful API的方式，被调用方会对外提供一个HTTP接口，调用方按要求发起HTTP请求并接收API接口返回的响应数据。

<aside>
💡

使用 RPC 的目的是让我们调用远程方法像调用本地方法一样无差别。并且基于RESTful API通常是基于HTTP协议，传输数据采用JSON等文本协议，相较于RPC 直接使用TCP协议，传输数据多采用二进制协议来说，RPC通常相比RESTful API性能会更好。

RESTful API多用于前后端之间的数据传输，而目前微服务架构下各个微服务之间多采用RPC调用。

</aside>

## 基础的RPC

Go语言的 rpc 包提供对通过网络或其他 i/o 连接导出的对象方法的访问，服务器注册一个对象，并把它作为服务对外可见（服务名称就是类型名称）。注册后，对象的导出方法将支持远程访问。服务器可以注册不同类型的多个对象(服务) ，但是不支持注册同一类型的多个对象。

rpc 包默认使用的是 gob 协议对传输数据进行序列化/反序列化，比较有局限性。

![image.png](%E5%BE%AE%E6%9C%8D%E5%8A%A1%20%E4%B8%8E%20%E4%BA%91%E5%8E%9F%E7%94%9F%201cf3895181d28010a46afdc8337886b8/image.png)

① 服务调用方（client）以本地调用方式调用服务；

② client stub接收到调用后负责将方法、参数等组装成能够进行网络传输的消息体；

③ client stub找到服务地址，并将消息发送到服务端；

④ server 端接收到消息；

⑤ server stub收到消息后进行解码；

⑥ server stub根据解码结果调用本地的服务；

⑦ 本地服务执行并将结果返回给server stub；

⑧ server stub将返回结果打包成能够进行网络传输的消息体；

⑨ 按地址将消息发送至调用方；

⑩ client 端接收到消息；

⑪ client stub收到消息并进行解码；

⑫ 调用方得到最终结果。

使用RPC框架的目标是只需要关心第1步和最后1步，中间的其他步骤统统封装起来，让使用者无需关心。例如社区中各式RPC框架（grpc、thrift等）就是为了让RPC调用更方便。

# GRPC

GRPC是什么

`gRPC`是一种现代化开源的高性能RPC框架，能够运行于任意环境之中。最初由谷歌进行开发。它使用HTTP/2作为传输协议。

在gRPC里，客户端可以像调用本地方法一样直接调用其他机器上的服务端应用程序的方法，帮助你更容易创建分布式应用程序和服务。与许多RPC系统一样，gRPC是基于定义一个服务，指定一个可以远程调用的带有参数和返回类型的的方法。在服务端程序中实现这个接口并且运行gRPC服务处理客户端调用。在客户端，有一个stub提供和服务端相同的方法。

**为什么要用gRPC**

使用gRPC， 我们可以一次性的在一个`.proto`文件中定义服务并使用任何支持它的语言去实现客户端和服务端，反过来，它们可以应用在各种场景中，从Google的服务器到你自己的平板电脑—— gRPC帮你解决了不同语言及环境间通信的复杂性。使用`protocol buffers`还能获得其他好处，包括高效的序列化，简单的IDL以及容易进行接口更新。总之一句话，使用gRPC能让我们更容易编写跨语言的分布式代码。

<aside>
💡

**Protocol Buffers**（简称 **Protobuf**）是由 **Google** 开发的一种 **数据序列化格式**，主要用于高效的结构化数据传输。它用于在不同的系统之间交换数据，特别是在分布式系统和网络通信中，像 gRPC 就是使用 Protobuf 来定义接口和交换消息的。

**IDL**（Interface description language）是指接口描述语言，是用来描述软件组件接口的一种计算机语言，是跨平台开发的基础。IDL通过一种中立的方式来描述接口，使得在不同平台上运行的对象和用不同语言编写的程序可以相互通信交流；比如，一个组件用C++写成，另一个组件用Go写成。

</aside>

## 安装gRPC

**安装gRPC**

在你的项目目录下执行以下命令，获取 gRPC 作为项目依赖

```go
go get google.golang.org/grpc@latest
```

**安装Protocol Buffers v3‘**

安装用于生成gRPC服务代码的协议编译器，最简单的方法是从下面的链接：[https://github.com/google/protobuf/releases](https://github.com/google/protobuf/releases)下载适合你平台的预编译好的二进制文件（`protoc-<version>-<platform>.zip`）

- bin 目录下的 protoc 是可执行文件。
- include 目录下的是 google 定义的`.proto`文件，我们`import "google/protobuf/timestamp.proto"`就是从此处导入。

我们需要将下载得到的可执行文件`protoc`所在的 bin 目录加到我们电脑的环境变量中。

**安装插件**

使用Go语言做开发，接下来执行下面的命令安装`protoc`的Go插件：

安装go语言插件：

```go
go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
```

该插件会根据`.proto`文件生成一个后缀为`.pb.go`的文件，包含所有`.proto`文件中定义的类型及其序列化方法。

安装grpc插件

```go
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2
```

该插件会生成一个后缀为`_grpc.pb.go`的文件，其中包含：

- 一种接口类型(或存根) ，供客户端调用的服务方法。
- 服务器要实现的接口类型。

上述命令会默认将插件安装到`$GOPATH/bin`，为了`protoc`编译器能找到这些插件，请确保你的`$GOPATH/bin`在环境变量中。

**检查**

```go
确认 protoc 安装完成。
❯ protoc --version
libprotoc 3.20.1
确认 protoc-gen-go 安装完成。
❯ protoc-gen-go --version
protoc-gen-go v1.28.0
确认 protoc-gen-go-grpc 安装完成。
❯ protoc-gen-go-grpc --version
protoc-gen-go-grpc 1.2.0
```

## 开发

**编写`.proto`文件定义服务**

像许多 RPC 系统一样，gRPC 基于定义服务的思想，指定可以通过参数和返回类型远程调用的方法。默认情况下，gRPC 使用 [protocol buffers](https://developers.google.com/protocol-buffers)作为接口定义语言(IDL)来描述服务接口和有效负载消息的结构。可以根据需要使用其他的IDL代替。

在gRPC中你可以定义四种类型的服务方法。

1. 普通 rpc，客户端向服务器发送一个请求，然后得到一个响应，就像普通的函数调用一样。
2. 服务器流式 rpc，其中客户端向服务器发送请求，并获得一个流来读取一系列消息。客户端从返回的流中读取，直到没有更多的消息。gRPC 保证在单个 RPC 调用中的消息是有序的。
3. 客户端流式 rpc，其中客户端写入一系列消息并将其发送到服务器，同样使用提供的流。一旦客户端完成了消息的写入，它就等待服务器读取消息并返回响应。同样，gRPC 保证在单个 RPC 调用中对消息进行排序。
4. 双向流式 rpc，其中双方使用读写流发送一系列消息。这两个流独立运行，因此客户端和服务器可以按照自己喜欢的顺序读写: 例如，服务器可以等待接收所有客户端消息后再写响应，或者可以交替读取消息然后写入消息，或者其他读写组合。每个流中的消息是有序的。

**生成指定语言的代码**

在 `.proto` 文件中的定义好服务之后，gRPC 提供了生成客户端和服务器端代码的 protocol buffers 编译器插件。

我们使用这些插件可以根据需要生成`Java`、`Go`、`C++`、`Python`等语言的代码。我们通常会在客户端调用这些 API，并在服务器端实现相应的 API。

- 在服务器端，服务器实现服务声明的方法，并运行一个 gRPC 服务器来处理客户端发来的调用请求。gRPC 底层会对传入的请求进行解码，执行被调用的服务方法，并对服务响应进行编码。
- 在客户端，客户端有一个称为存根（stub）的本地对象，它实现了与服务相同的方法。然后，客户端可以在本地对象上调用这些方法，将调用的参数包装在适当的 protocol buffers 消息类型中—— gRPC 在向服务器发送请求并返回服务器的 protocol buffers 响应之后进行处理。

**编写业务逻辑代码**

在服务端编写业务代码实现具体的服务方法，在客户端按需调用这些方法。

## protocol buffers

为了生成 Go 代码，必须为每个 `.proto` 文件（包括那些被生成的 `.proto` 文件传递依赖的文件）提供 Go 包的导入路径。有两种方法可以指定 Go 导入路径：

- 通过在 `.proto` 文件中声明它。
- 通过在调用 `protoc` 时在命令行上声明它。

建议在 `.proto` 文件中声明它，以便 `.proto` 文件的 Go 包可以与 `.proto` 文件本身集中标识，并简化调用 `protoc`时传递的标志集。 如果给定 `.proto` 文件的 Go 导入路径由 `.proto` 文件本身和命令行提供，则后者优先于前者。

```go
option go_package = "example.com/project/protos/fizz";
```

```go
syntax = "proto3";

message SearchRequest {
  string query = 1;
  int32 page_number = 2;
  int32 result_per_page = 3;
}
```

• 文件的第一行指定使用 `proto3` 语法: 如果不这样写，protocol buffer编译器将假定你使用 `proto2`。这个声明必须是文件的第一个非空非注释行。

• `SearchRequest` 消息定义指定了三个字段(名称/值对) ，每个字段表示希望包含在此类消息中的每一段数据。每个字段都有一个名称和一个类型。

**指定字段类型**

**分配字段编号**

消息定义中的每个字段都有一个**唯一的编号**。这些字段编号用来在消息二进制格式中标识字段，在消息类型使用后就不能再更改。注意，范围1到15中的字段编号需要一个字节进行编码，包括字段编号和字段类型。范围16到2047的字段编号采用两个字节。因此，应该为经常使用的消息元素保留数字1到15的编号。切记为将来可能添加的经常使用的元素留出一些编号。

**指定字段规则**

- `singular`: 格式正确的消息可以有这个字段的零个或一个(但不能多于一个)。这是 proto3语法的默认字段规则。
- `repeated`: 该字段可以在格式正确的消息中重复任意次数(包括零次)。重复值的顺序将被保留。

**添加更多消息类型**

**保留字段**

如果你通过完全删除字段或将其注释掉来**更新**消息类型，那么未来的用户在对该类型进行自己的更新时可以重用字段号。如果其他人以后加载旧版本的相同`.proto`文件，这可能会导致严重的问题，包括数据损坏，隐私漏洞等等。确保这种情况不会发生的一种方法是指定已删除字段的字段编号(和/或名称，这也可能导致 JSON 序列化问题)是保留的（`reserved`）。如果将来有任何用户尝试使用这些字段标识符，protocol buffer编译器将发出提示。

```go
message Foo {
  reserved 2, 15, 9 to 11;
  reserved "foo", "bar";
}
// 注意，不能在同一个reserved语句中混合字段名和字段编号。
```

```go
double -- float64
float  -- float32
```

当解析消息时，如果编码消息不包含特定的 singular 元素，则解析对象中的相应字段将设置为该字段的默认值。

- 对于字符串，默认值为空字符串。
- 对于字节，默认值为空字节。
- 对于布尔值，默认值为 false。
- 对于数值类型，默认值为零。
- 对于枚举，默认值是第一个定义的枚举值，该值必须为0。

repeated 字段的默认值为空(通常是适当语言中的空列表)。

**枚举**

```go
message SearchRequest {
  string query = 1;
  int32 page_number = 2;
  int32 result_per_page = 3;
  enum Corpus {
    UNIVERSAL = 0;
    WEB = 1;
    IMAGES = 2;
    LOCAL = 3;
    NEWS = 4;
    PRODUCTS = 5;
    VIDEO = 6;
  }
  Corpus corpus = 4;
}
```

Corpus enum 的第一个常量映射为零: 每个 enum 定义必须包含一个常量，该常量映射为零作为它的第一个元素。这是因为:

1. 必须有一个零值，这样我们就可以使用0作为数值默认值。
2. 零值必须是第一个元素，以便与 proto2语义兼容，其中第一个枚举值总是默认值。

你可以通过将相同的值分配给不同的枚举常量来定义别名。为此，你需要将 `allow _ alias` 选项设置为 `true`，否则，当发现别名时，protocol 编译器将生成错误消息。

```go
message MyMessage1 {
  enum EnumAllowingAlias {
    option allow_alias = true;
    UNKNOWN = 0;
    STARTED = 1;
    RUNNING = 1;
  }
}
```

枚举的常数必须在32位整数的范围内。由于枚举值在传输时使用[变长编码](https://developers.google.com/protocol-buffers/docs/encoding)，因此负值效率低，因此不推荐使用

如果通过完全删除枚举条目或注释掉枚举类型来更新枚举类型，那么未来的用户在自己更新该类型时可以重用该数值。这可能会导致严重的问题，如果以后有人加载旧版本的相同`.proto`文件，包括数据损坏，隐私漏洞等等。确保不发生这种情况的一种方法是指定已删除条目的数值(和/或名称，这也可能导致 JSON 序列化问题)为 `reserved`。

**使用其他消息类型**

```go
message SearchResponse {
  repeated Result results = 1;
}

message Result {
  string url = 1;
  string title = 2;
  repeated string snippets = 3;
}
```

**oneof**

如果你有一条包含多个字段的消息，并且最多同时设置其中一个字段，那么你可以通过使用`oneof`来实现并节省内存。

`oneof`字段类似于常规字段，只不过`oneof`中的所有字段共享内存，而且最多可以同时设置一个字段。设置其中的任何成员都会自动清除所有其他成员。根据所选择的语言，可以使用特殊 `case()`或 `WhichOneof()` 方法检查 one of 中的哪个值被设置(如果有的话)。

```go
message SampleMessage {
  oneof test_oneof {
    string name = 4;
    SubMessage sub_message = 9;
  }
}
```

**Maps**

```go
map<key_type, value_type> map_field = N;
```

其中`key_type`可以是任何整型或字符串类型(因此，除了浮点类型和字节以外的任何标量类型) 。注意，枚举不是有效的`key_type`。`value_type`可以是除另一个映射以外的任何类型。

**定义服务**

## oneof

如果你有一条包含多个字段的消息，并且最多同时设置其中一个字段，那么你可以通过使用`oneof`来实现并节省内存。

`oneof`字段类似于常规字段，只不过`oneof`中的所有字段共享内存，而且最多可以同时设置一个字段。设置其中的任何成员都会自动清除所有其他成员。

可以在`oneof`中添加除了map字段和repeated字段外的任何类型的字段。

### WrapValues

当我们有如下消息定义时，我们拿到一个book消息，当`book.Price = 0`时我们没办法区分`book.Price`字段是未赋值还是被赋值为0。

```go
message Book {
    string title = 1;
    string author = 2;
    int64 price = 3;
}
```

这种场景推荐使用`google/protobuf/wrappers.proto`中定义的WrapValue，本质上就是使用自定义message代替基本类型。

```go
message Book {
    string title = 1;
    string author = 2;
    google.protobuf.Int64Value price = 3;
}
// client
wrapperspb.Int64Value{Value: 9900}
// server
book.GetPrice() == nil
```

Protobuf **v3.15.0** 版本之后又支持使用`optional`显式指定字段为可选

```go
message Book {
    string title = 1;
    string author = 2;
    //google.protobuf.Int64Value price = 3;
    optional int64 price = 3;  // 使用optional
}
// client
Price: proto.Int64(9900),
// server
**book.Price == nil**
```

### FiledMask

需要实现一个更新书籍信息接口，但是如果我们的`Book`中定义有很多很多字段时，我们不太可能每次请求都去全量更新`Book`的每个字段，因为通常每次操作只会更新1到2个字段。

google/protobuf/field_mask.proto,它能够记录在一次更新请求中涉及到的具体字段路径。

```go
message UpdateBookRequest {
    // 操作人 
    string op = 1;
    // 要更新的书籍信息
    Book book = 2;

    // 要更新的字段
    google.protobuf.FieldMask update_mask = 3;
}
// client
UpdateMask: &fieldmaskpb.FieldMask{Paths: paths},
// server
import fieldmask_utils "github.com/mennanov/fieldmask-utils"
mask, _ := fieldmask_utils.MaskFromProtoFieldMask(updateReq.UpdateMask, generator.CamelCase)
var bookDst = make(map[string]interface{})
// 将数据读取到map[string]interface{}
// fieldmask-utils支持读取到结构体等，更多用法可查看文档。
fieldmask_utils.StructToMap(mask, updateReq.Book, bookDst)
```

通过`paths`记录本次更新的字段路径，如果是嵌套的消息类型则通过`x.y`的方式标识。

在收到更新消息后，我们需要根据`UpdateMask`字段中记录的更新路径去读取更新数据。这里借助第三方库[github.com/mennanov/fieldmask-utils](https://github.com/mennanov/fieldmask-utils)实现。

**2022-11-20更新**：由于`github.com/golang/protobuf/protoc-gen-go/generator`包已弃用，而`MaskFromProtoFieldMask`函数（签名如下

```go
func MaskFromProtoFieldMask(fm *field_mask.FieldMask, naming func(string) string) (Mask, error)
```

接收的`naming`参数本质上是一个将字段掩码字段名映射到 Go 结构中使用的名称的函数，它必须根据你的实际需求实现。

例如在我们这个示例中，还可以使用`github.com/iancoleman/strcase`包提供的`ToCamel`方法：

```go
import "github.com/iancoleman/strcase"
import fieldmask_utils "github.com/mennanov/fieldmask-utils"

mask, _ := fieldmask_utils.MaskFromProtoFieldMask(updateReq.UpdateMask, strcase.ToCamel)
var bookDst = make(map[string]interface{})
// 将数据读取到map[string]interface{}
// fieldmask-utils支持读取到结构体等，更多用法可查看文档。
fieldmask_utils.StructToMap(mask, updateReq.Book, bookDst)
// do update with bookDst
fmt.Printf("bookDst:%#v\n", bookDst)

```

# 流式grpc

**服务端流式RPC**

客户端发出一个RPC请求，服务端与客户端之间建立一个单向的流，服务端可以向流中写入多个响应消息，最后主动关闭流；而客户端需要监听这个流，不断获取响应直到流关闭。应用场景举例：客户端向服务端发送一个股票代码，服务端就把该股票的实时数据源源不断的返回给客户端。

```go
// 服务端返回流式数据
rpc LotsOfReplies(HelloRequest) returns (stream HelloResponse);
// server
 stream.Send(data)
 // client
 stream.Recv()
```

**客户端流式RPC**

客户端传入多个请求对象，服务端返回一个响应结果。典型的应用场景举例：物联网终端向服务器上报数据、大数据流式计算等。

```go
// 客户端发送流式数据
rpc LotsOfGreetings(stream HelloRequest) returns (HelloResponse);
// serve 
stream.Recv()
stream.SendAndClose()
// client
stream.Send
stream.CloseAndRecv()
```

**双向流式RPC**

双向流式RPC即客户端和服务端均为流式的RPC，能发送多个请求对象也能接收到多个响应对象。典型应用示例：聊天应用等。

```go
// 双向流式数据
rpc BidiHello(stream HelloRequest) returns (stream HelloResponse);
// server
stream.Recv()
	// 处理请求数据
stream.Send
// client
go func() {stream.Recv() }()
stream.Send
stream.CloseSend()
```

### metadata

元数据（[metadata](https://pkg.go.dev/google.golang.org/grpc/metadata)）是指在处理RPC请求和响应过程中需要但又不属于具体业务（例如身份验证详细信息）的信息，采用键值对列表的形式，其中键是`string`类型，值通常是`[]string`类型，但也可以是二进制数据。gRPC中的 metadata 类似于我们在 HTTP headers中的键值对，元数据可以包含认证token、请求标识和监控标签等。

metadata中的键是大小写不敏感的，由字母、数字和特殊字符`-`、`_`、`.`组成并且不能以`grpc-`开头（gRPC保留自用），二进制值的键名必须以`-bin`结尾。

元数据对 gRPC 本身是不可见的，我们通常是在应用程序代码或中间件中处理元数据，我们不需要在`.proto`文件中指定元数据。

如何访问元数据取决于具体使用的编程语言。 在Go语言中我们是用[google.golang.org/grpc/metadata](https://pkg.go.dev/google.golang.org/grpc/metadata)这个库来操作metadata。

**创建新的metadata**

```go
md := metadata.New(map[string]string{"key1": "val1", "key2": "val2"})
md := metadata.Pairs(
    "key1", "val1",
    "key1", "val1-2", // "key1"的值将会是 []string{"val1", "val1-2"}
    "key2", "val2",
)
// 注意: 所有的键将自动转换为小写，因此“ kEy1”和“ Key1”将是相同的键，它们的值将合并到相同的列表中。这种情况适用于 New 和 Pair。
```

**元数据中存储二进制数据**

**元数据中存储二进制数据在元数据中，键始终是字符串。但是值可以是字符串或二进制数据。要在元数据中存储二进制数据值，只需在密钥中添加“-bin”后缀。在创建元数据时，将对带有“-bin”后缀键的值进行编码:**

```go
md := metadata.Pairs(
    "key", "string value",
    "key-bin", string([]byte{96, 102}), // 二进制数据在发送前会进行(base64) 编码
                                        // 收到后会进行解码
)
```

**从请求上下文中获取元数据**

可以使用 `FromIncomingContext` 可以从RPC请求的上下文中获取元数据:

```go
func (s *server) SomeRPC(ctx context.Context, in *pb.SomeRequest) (*pb.SomeResponse, err) {
    md, ok := metadata.FromIncomingContext(ctx)
    // do something with metadata
}
```

**发送和接收元数据-客户端**

**发送metadata**

有两种方法可以将元数据发送到服务端。推荐的方法是使用 `AppendToOutgoingContext` 将 kv 对附加到context。无论context中是否已经有元数据都可以使用这个方法。如果先前没有元数据，则添加元数据; 如果context中已经存在元数据，则将 kv 对合并进去。

```go
// 创建带有metadata的context
ctx := metadata.AppendToOutgoingContext(ctx, "k1", "v1", "k1", "v2", "k2", "v3")

// 添加一些 metadata 到 context (e.g. in an interceptor)
ctx := metadata.AppendToOutgoingContext(ctx, "k3", "v4")

// 发起普通RPC请求
response, err := client.SomeRPC(ctx, someRequest)

// 或者发起流式RPC请求
stream, err := client.SomeStreamingRPC(ctx)
```

或者，可以使用 `NewOutgoingContext` 将元数据附加到context。但是，这将替换context中的任何已有的元数据，因此必须注意保留现有元数据(如果需要的话)。这个方法比使用 `AppendToOutgoingContext` 要慢。这方面的一个例子如下:

```go
// 创建带有metadata的context
md := metadata.Pairs("k1", "v1", "k1", "v2", "k2", "v3")
ctx := metadata.NewOutgoingContext(context.Background(), md)

// 添加一些metadata到context (e.g. in an interceptor)
send, _ := metadata.FromOutgoingContext(ctx)
newMD := metadata.Pairs("k3", "v3")
ctx = metadata.NewOutgoingContext(ctx, metadata.Join(send, newMD))

// 发起普通RPC请求
response, err := client.SomeRPC(ctx, someRequest)

// 或者发起流式RPC请求
stream, err := client.SomeStreamingRPC(ctx)
```

**接收metadata**

客户端可以接收的元数据包括header和trailer。

<aside>
💡

trailer可以用于服务器希望在处理请求后给客户端发送任何内容，例如在流式RPC中只有等所有结果都流到客户端后才能计算出负载信息，这时候就不能使用headers（header在数据之前，trailer在数据之后）。

</aside>

**普通调用**

可以使用 [CallOption](https://godoc.org/google.golang.org/grpc#CallOption) 中的 [Header](https://godoc.org/google.golang.org/grpc#Header) 和 [Trailer](https://godoc.org/google.golang.org/grpc#Trailer) 函数来获取普通RPC调用发送的header和trailer:

```go
var header, trailer metadata.MD // 声明存储header和trailer的变量
r, err := client.SomeRPC(
    ctx,
    someRequest,
    grpc.Header(&header),    // 将会接收header
    grpc.Trailer(&trailer),  // 将会接收trailer
)

// do something with header and trailer
```

**流式调用**

```go
stream, err := client.SomeStreamingRPC(ctx)

// 接收 header
header, err := stream.Header()

// 接收 trailer
trailer := stream.Trailer()
```

**发送和接收元数据-服务器端**

**接收metadata**

**要读取客户端发送的元数据，服务器需要从 RPC 上下文检索它。如果是普通RPC调用，则可以使用 RPC 处理程序的上下文。对于流调用，服务器需要从流中获取上下文。**

**普通调用**

```go
func (s *server) SomeRPC(ctx context.Context, in *pb.someRequest) (*pb.someResponse, error) {
    md, ok := metadata.FromIncomingContext(ctx)
    // do something with metadata
}
```

**流式调用**

```go
func (s *server) SomeStreamingRPC(stream pb.Service_SomeStreamingRPCServer) error {
    md, ok := metadata.FromIncomingContext(stream.Context()) // get context from stream
    // do something with metadata
}
```

**发送metadata**

**普通调用**

在普通调用中，服务器可以调用 [grpc](https://godoc.org/google.golang.org/grpc) 模块中的 [SendHeader](https://godoc.org/google.golang.org/grpc#SendHeader) 和 [SetTrailer](https://godoc.org/google.golang.org/grpc#SetTrailer) 函数向客户端发送header和trailer。这两个函数将context作为第一个参数。它应该是 RPC 处理程序的上下文或从中派生的上下文：

```go
func (s *server) SomeRPC(ctx context.Context, in *pb.someRequest) (*pb.someResponse, error) {
    // 创建和发送 header
    header := metadata.Pairs("header-key", "val")
    grpc.SendHeader(ctx, header)
    // 创建和发送 trailer
    trailer := metadata.Pairs("trailer-key", "val")
    grpc.SetTrailer(ctx, trailer)
}
```

**流式调用**

对于流式调用，可以使用接口 [ServerStream](https://godoc.org/google.golang.org/grpc#ServerStream) 中的 `SendHeader` 和 `SetTrailer` 函数发送header和trailer:

```go
func (s *server) SomeStreamingRPC(stream pb.Service_SomeStreamingRPCServer) error {
    // 创建和发送 header
    header := metadata.Pairs("header-key", "val")
    stream.SendHeader(header)
    // 创建和发送 trailer
    trailer := metadata.Pairs("trailer-key", "val")
    stream.SetTrailer(trailer)
}
```

### 错误处理

**gRPC code**

```go
import "google.golang.org/grpc/codes"
```

**gRPC Status**

```go
import "google.golang.org/grpc/status"
```

**创建错误**

当遇到错误时，gRPC服务的方法函数应该创建一个 `status.Status`。通常我们会使用 `status.New`函数并传入适当的`status.Code`和错误描述来生成一个`status.Status`。调用`status.Err`方法便能将一个`status.Status`转为`error`类型。也存在一个简单的`status.Error`方法直接生成`error`。

```go
// 创建status.Status
st := status.New(codes.NotFound, "some description")
err := st.Err()  // 转为error类型

// vs.

err := status.Error(codes.NotFound, "some description")

```

**为错误添加其他详细信息**

在某些情况下，可能需要为服务器端的特定错误添加详细信息。`status.WithDetails`就是为此而存在的，它可以添加任意多个`proto.Message`，我们可以使用`google.golang.org/genproto/googleapis/rpc/errdetails`中的定义或自定义的错误详情。

```go
st := status.New(codes.ResourceExhausted, "Request limit exceeded.")
ds, _ := st.WithDetails(
	// proto.Message
)
return nil, ds.Err()
```

然后，客户端可以通过首先将普通`error`类型转换回`status.Status`，然后使用`status.Details`来读取这些详细信息。

```go
s := status.Convert(err)
for _, d := range s.Details() {
	// ...
}
```

### 加密或认证

**无加密认证**

```go
//client 
grpc.WithTransportCredentials(insecure.NewCredentials()
```

**使用服务器身份验证 SSL/TLS**

gRPC 内置支持 SSL/TLS，可以通过 SSL/TLS 证书建立安全连接，对传输的数据进行加密处理。

**生成证书**

**生成私钥**

```go
openssl ecparam -genkey -name secp384r1 -out server.key
// 这里生成的是ECC私钥
```

**生成自签名的证书**

```go
// server.cnf
[ req ]
default_bits       = 4096
default_md		= sha256
distinguished_name = req_distinguished_name
req_extensions     = req_ext

[ req_distinguished_name ]
countryName                 = Country Name (2 letter code)
countryName_default         = CN
stateOrProvinceName         = State or Province Name (full name)
stateOrProvinceName_default = GUANGZHOU
localityName                = Locality Name (eg, city)
localityName_default        = GUANGZHOU
organizationName            = Organization Name (eg, company)
organizationName_default    = DEV
commonName                  = Common Name (e.g. server FQDN or YOUR name)
commonName_max              = 64
commonName_default          = yuanshuhao.com

[ req_ext ]
subjectAltName = @alt_names

[alt_names]
DNS.1   = localhost
DNS.2   = yuanshuhao.com
IP      = 127.0.0.1

// 生成自签名证书
openssl req -nodes -new -x509 -sha256 -days 3650 -config server.cnf -extensions 'req_ext' -key server.key -out server.crt
```

**建立安全连接**

Server端使用`credentials.NewServerTLSFromFile`函数分别加载证书`server.cert`和秘钥`server.key`

```go
creds, _ := credentials.NewServerTLSFromFile(certFile, keyFile)
s := grpc.NewServer(grpc.Creds(creds))
```

client端使用上一步生成的证书文件——`server.cert`建立安全连接

```go
creds, _ := credentials.NewClientTLSFromFile(certFile, "")
conn, _ := grpc.NewClient("127.0.0.1:8972", grpc.WithTransportCredentials(creds))
```

### 拦截器（中间件）

gRPC 为在每个 ClientConn/Server 基础上实现和安装拦截器提供了一些简单的 API。 拦截器拦截每个 RPC 调用的执行。用户可以使用拦截器进行日志记录、身份验证/授权、指标收集以及许多其他可以跨 RPC 共享的功能。

在 gRPC 中，拦截器根据拦截的 RPC 调用类型可以分为两类。第一个是普通拦截器（一元拦截器），它拦截普通RPC 调用。另一个是流拦截器，它处理流式 RPC 调用。而客户端和服务端都有自己的普通拦截器和流拦截器类型。因此，在 gRPC 中总共有四种不同类型的拦截器。

**客户端端拦截器**

**普通拦截器/一元拦截器**

[UnaryClientInterceptor](https://godoc.org/google.golang.org/grpc#UnaryClientInterceptor) 是客户端一元拦截器的类型，它的函数前面如下：

```go
func(ctx context.Context, method string, req, reply interface{}, cc *ClientConn, invoker UnaryInvoker, opts ...CallOption) error
```

一元拦截器的实现通常可以分为三个部分: 调用 RPC 方法之前（预处理）、调用 RPC 方法（RPC调用）和调用 RPC 方法之后（调用后）。

- 预处理：用户可以通过检查传入的参数(如 RPC 上下文、方法字符串、要发送的请求和 CallOptions 配置)来获得有关当前 RPC 调用的信息。
- RPC调用：预处理完成后，可以通过执行`invoker`执行 RPC 调用。
- 调用后：一旦调用者返回应答和错误，用户就可以对 RPC 调用进行后处理。通常，它是关于处理返回的响应和错误的。 若要在 `ClientConn` 上安装一元拦截器，请使用`DialOptionWithUnaryInterceptor`的`DialOption`配置 Dial 。

**流拦截器**

[StreamClientInterceptor](https://godoc.org/google.golang.org/grpc#StreamClientInterceptor)是客户端流拦截器的类型。它的函数签名是

```go
func(ctx context.Context, desc *StreamDesc, cc *ClientConn, method string, streamer Streamer, opts ...CallOption) (ClientStream, error)
```

流拦截器的实现通常包括预处理和流操作拦截。

- 预处理：类似于上面的一元拦截器。
- 流操作拦截：流拦截器并没有事后进行 RPC 方法调用和后处理，而是拦截了用户在流上的操作。首先，拦截器调用传入的`streamer`以获取 `ClientStream`，然后包装 `ClientStream` 并用拦截逻辑重载其方法。最后，拦截器将包装好的 `ClientStream` 返回给用户进行操作。

若要为 `ClientConn` 安装流拦截器，请使用`WithStreamInterceptor`的 DialOption 配置 Dial。

**server端拦截器**

**普通拦截器/一元拦截器**

[UnaryServerInterceptor](https://godoc.org/google.golang.org/grpc#UnaryServerInterceptor)是服务端的一元拦截器类型，它的函数签名是

```go
func(ctx context.Context, req interface{}, info *UnaryServerInfo, handler UnaryHandler) (resp interface{}, err error)
```

服务端一元拦截器具体实现细节和客户端版本的类似。

若要为服务端安装一元拦截器，请使用 `UnaryInterceptor` 的`ServerOption`配置 `NewServer`。

**流拦截器**

[StreamServerInterceptor](https://godoc.org/google.golang.org/grpc#StreamServerInterceptor)是服务端流式拦截器的类型，它的签名如下：

```go
func(srv interface{}, ss ServerStream, info *StreamServerInfo, handler StreamHandler) error
```

若要为服务端安装流拦截器，请使用 `StreamInterceptor` 的`ServerOption`来配置 `NewServer`。

## grpc-gateway

[gRPC-Gateway](https://github.com/grpc-ecosystem/grpc-gateway) 是一个 protoc 插件。它读取 gRPC 服务定义并生成一个反向代理服务器，该服务器将 RESTful JSON API 转换为 gRPC。此服务器根据 gRPC 定义中的自定义选项生成。

![image.png](%E5%BE%AE%E6%9C%8D%E5%8A%A1%20%E4%B8%8E%20%E4%BA%91%E5%8E%9F%E7%94%9F%201cf3895181d28010a46afdc8337886b8/image%201.png)

添加 gRPC-Gateway 注释

这些注释定义了 gRPC 服务如何映射到 JSON 请求和响应。使用 protocol buffers时，每个 RPC 服务必须使用 `google.api.HTTP` 注释来定义 HTTP 方法和路径。

需要将 `google/api/http.proto` 导入到 `proto` 文件中。我们还需要添加所需的 HTTP-> gRPC 映射。

```go
// 导入google/api/annotations.proto
import "google/api/annotations.proto";
// 定义一个Greeter服务
service Greeter {
  // 打招呼方法
  rpc SayHello (HelloRequest) returns (HelloReply) {
    // 这里添加了google.api.http注释
    option (google.api.http) = {
      post: "/v1/example/echo" // 路径
      body: "*" // 所有的body都映射
    };
  }
}

go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@v2

$ protoc -I=proto \
   --go_out=proto --go_opt=paths=source_relative \
   --go-grpc_out=proto --go-grpc_opt=paths=source_relative \
   --grpc-gateway_out=proto --grpc-gateway_opt=paths=source_relative \
   helloworld/hello_world.proto

// server
// 创建一个连接到刚刚启动的grpc服务器的客户端连接
// gRPC-Gateway 就是通过它来代理请求（将HTTP请求转为RPC请求）
	conn, err := grpc.NewClient(
		"127.0.0.1:8080",
		// grpc.WithBlock()：确保连接是 同步 的，直到连接建立成功。
		grpc.WithBlock(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil{
		log.Fatalln("falied to dial server:", err)
	}

	// 创建 gRPC-Gateway Mux（路由器）：它会将 HTTP 请求映射到 gRPC 服务。
	gwmux := runtime.NewServeMux()

	// 注册 gRPC 服务与 HTTP 路由的映射
	// proto.RegisterGretterHandler 是自动生成的代码，用于将 gRPC 服务注册到 gRPC-Gateway 的HTTP 路由中。
	// 这一行的作用是将 gRPC 服务的接口（conn）与 HTTP 请求的路由 (gwmux) 关联起来。
	err = proto.RegisterGretterHandler(context.Background(), gwmux, conn)
	if err != nil {
		log.Fatalln("falied to register gateway:", err)
	}

	// 这一行代码创建了一个 HTTP 服务器，绑定在 8090 端口，
	// 并将之前创建的 gwmux（HTTP 路由器）作为处理请求的 handler。
	gwServer := &http.Server{Addr: ":8090", Handler: gwmux}

	log.Fatalln(gwServer.ListenAndServe())
	
	// !!!
		go func() {
		log.Fatalln(s.Serve(lis))
	}()
```

每个映射指定一个 URL 路径模板和一个 HTTP 方法。 路径模板可以引用 gRPC 请求消息中的一个或多个字段，只要每个字段是原始（非消息）类型的非重复字段即可。

映射到 URL 查询参数的字段必须具有原始类型或重复的原始类型或非重复的消息类型。 在重复类型的情况下，参数可以在URL中重复为`...?param=A&param=B`。在消息类型的情况下，消息的每个字段都映射到一个单独的参数

对于允许请求正文（request body）的 HTTP 方法，`body`字段指定映射关系

特殊名称 `*` 可用于主体映射来定义不受路径模板绑定的每个字段都应映射到请求正文。

在正文映射中使用 `*` 时，不可能有 HTTP 参数，因为所有不受路径绑定的字段都在正文中结束。这使得该选项在定义 REST API 时很少在实践中使用。 `*` 的常见用法是在根本不使用 URL 来传输数据的自定义方法中。

可以使用 `additional_bindings` 选项为一个 RPC 定义多个 HTTP 方法。

1. 路径变量**不得**引用任何repeated或mapped的字段，因为客户端库无法处理此类变量扩展。

1. 路径变量**不得**捕获前导“/”字符。 原因是最常见的用例“{var}”没有捕获前导“/”字符。 为了一致性，所有路径变量必须共享相同的行为。
2. 不能将重复的消息字段（repeated message）映射到 URL 查询参数，因为没有客户端库可以支持如此复杂的映射。
3. 如果 API 需要为请求或响应正文使用 JSON 数组，它可以将请求或响应正文映射到repeated字段。 但是，某些 gRPC 转码实现可能不支持此功能。

## 基于游标分页

**基于偏移量的分页**

```go
SELECT id, title FROM books ORDER BY id ASC LIMIT 10 OFFSET 10;
```

### 优势

1. 简单
2. 支持跳页访问

### 劣势

1. 基于偏移量的分页在数据量很大的场景下，查询效率会比较低。通常 OFFSET 越高，查询时间就越长。
2. 在并发场景下会出现元素重复（offset在第二页时有人在第一页新插入一个数据）或被跳过（offset在第二页时有人在第一页删掉了一个数据）。
3. 显式的page参数在支持跳页的同时也会被爬虫并发请求。

**基于游标的分页/基于令牌的分页**

基于游标的分页是指接口在返回响应数据的同时返回一个`cursor`——通常是一个不透明字符串。它表示的是这一页数据的最后那个元素（这就像是我们玩单机游戏的存档点，这一次我们从这里离开，下一次将从这里继续），通过这个`cursor` API 就能准确的返回下一页的数据。

用于游标的字段必须是唯一的、连续的列，数据集将基于该列进行排序。在处理实时数据时使用基于游标的分页。第一页请求不需要提供游标，但是后续的请求必须携带游标。

```go
SELECT id, title FROM books WHERE id > 10 ORDER BY id ASC LIMIT 10;
```

### 优势

1. 性能好
2. 并发安全
3. 防止被无脑批量爬取

### 劣势

1. 实现稍复杂
2. 不支持跳页（但现在流行无限滑动翻页）。
3. 不太适合多检索条件的场景。

我们在使用基于游标的分页时，通常并不会把具体的`cursor`数据显式拼接到API URL中，而是使用通常会被命名为`next`、`next_cursor`、`after`或`page_token`的不透明字符串。

**基于游标的分页实现方案**

```go
type Page struct {
	NextID        string `json:"next_id"`
	NextTimeAtUTC int64  `json:"next_time_at_utc"`
	PageSize      int64  `json:"page_size"`
}

// Encode 返回分页token
func (p Page) Encode() Token {
	b, err := json.Marshal(p)
	if err != nil {
		return Token("")
	}
	return Token(base64.StdEncoding.EncodeToString(b))
}

type Token string

// Decode 解析分页信息
func (t Token) Decode() Page {
	var result Page
	if len(t) == 0 {
		return result
	}

	bytes, err := base64.StdEncoding.DecodeString(string(t))
	if err != nil {
		return result
	}

	err = json.Unmarshal(bytes, &result)
	if err != nil {
		return result
	}

	return result
}
```

**同一个端口提供HTTP API和gRPC API**

```go
**gwmux := runtime.NewServeMux()**
dops := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
err = helloworldpb.RegisterGreeterHandlerFromEndpoint(context.Background(), gwmux, "127.0.0.1:8091", dops)

mux := http.NewServeMux()
	mux.Handle("/", gwmux)

	// 定义HTTP server配置
	gwServer := &http.Server{
		Addr:    "127.0.0.1:8091",
		Handler: grpcHandlerFunc(s, mux), // 请求的统一入口
	}
	
	// grpcHandlerFunc 将gRPC请求和HTTP请求分别调用不同的handler处理
func grpcHandlerFunc(grpcServer *grpc.Server, otherHandler http.Handler) http.Handler {
	return h2c.NewHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.ProtoMajor == 2 && strings.Contains(r.Header.Get("Content-Type"), "application/grpc") {
			grpcServer.ServeHTTP(w, r)
		} else {
			otherHandler.ServeHTTP(w, r)
		}
	}), &http2.Server{})
}
```

- **区分 HTTP 和 gRPC 请求**： 传统的 HTTP 请求和 gRPC 请求通常都使用 HTTP/2 协议，但它们的处理方式是不同的。gRPC 使用的是 `application/grpc` 作为 `Content-Type`，因此通过检查 `Content-Type` 头部，能够区分是普通的 HTTP 请求还是 gRPC 请求。
- **处理不同类型的请求**： 由于 gRPC 和传统 HTTP 请求的处理方式不同（gRPC 使用 Protobuf 进行数据传输，而 HTTP 请求通常使用 JSON 或 XML 等格式），你需要在同一个服务器上为这两种不同的请求提供不同的处理方式。这就是为什么需要通过检查请求头来区分它们。

## **gRPC中的名称解析和负载均衡**

**DNS解析器**

gRPC中默认使用的名称解析器是 DNS，即在gRPC客户端执行`grpc.Dial`时提供域名，默认会将DNS解析出对应的IP列表返回。

使用默认DNS解析器的名称语法为：`dns:[//authority/]host[:port]`

**consul resolver**

**自定义解析器**

```go
package main

import (
	"google.golang.org/grpc/resolver"
)

const (
	mySchme    = "ysh"
	myEndpoion = "resolver.ysh.com"
)

var addrs = []string{"127.0.0.1:8972", "127.0.0.1:8973", "127.0.0.1:8974"}

// 自定义name resolver，实现Resolver接口
type yshResolver struct {
	target      resolver.Target
	cc          resolver.ClientConn
	addresStore map[string][]string
}

func (r *yshResolver) ResolveNow(o resolver.ResolveNowOptions){
	addrsStrs := r.addresStore[r.target.Endpoint()]
	addrList := make([]resolver.Address, len(addrsStrs))
	for i, s := range addrsStrs {
		addrList[i] = resolver.Address{Addr: s}
	}
	r.cc.UpdateState(resolver.State{Addresses: addrList})
}

func (*yshResolver) Close() {}

// 需实现 Builder 接口
type yshResloverBuiler struct{}

func (*yshResloverBuiler) Build(target resolver.Target, cc resolver.ClientConn,  opts resolver.BuildOptions)  (resolver.Resolver, error) {
	r := &yshResolver{
		target: target,
		cc: cc,
		addresStore: map[string][]string{
			myEndpoion: addrs,
		},
	}
	r.ResolveNow(resolver.ResolveNowOptions{})
	return r, nil
}

func (*yshResloverBuiler) Scheme() string {return mySchme}

func init() {
	resolver.Register(&yshResloverBuiler{})
}

conn, err := grpc.NewClient("ysh:///resolver.ysh.com",
		grpc.WithTransportCredentials(insecure.NewCredentials()))
		
也可以在客户端建立连接时通过grpc.WithResolvers指定使用的名称解析器，使用这种方法就不需要事先注册名称解析器了。
conn, err := grpc.NewClient("ysh:///resolver.ysh.com",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithResolvers(&yshResolverBuilder{}),)
```

**负载均衡策略**

gRPC 中的负载均衡是基于每次调用的，而不是基于每个连接的。换句话说，即使所有请求来自单个客户端，我们仍然希望它们在所有服务器之间负载均衡（雨露均沾）。

gRPC-go 内置支持有 `pick_first` (默认值)和 `round_robin` 两种策略。

- `pick_first`是 gRPC 负载均衡的默认值，因此不需要设置。`pick_first` 会尝试连接取到的第一个服务端地址，如果连接成功，则将其用于所有 RPC，如果连接失败，则尝试下一个地址(并继续这样做，直到一个连接成功)。因此，所有的 RPC 将被发送到同一个后端。所有接收到的响应都显示相同的后端地址。
- `round_robin` 连接到它所看到的所有地址，并按顺序一次向每个server发送一个 RPC。例如，我们现在注册有两个server，第一个 RPC 将被发送到 server-1，第二个 RPC 将被发送到 server-2，第三个 RPC 将再次被发送到 server-1。

```go
	conn, err := grpc.NewClient("ysh:///resolver.ysh.com",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"round_robin"}`))
```

# **基于 consul 实现服务注册与发现**

[https://zhuanlan.zhihu.com/p/32052223](https://zhuanlan.zhihu.com/p/464258409)

[https://zhuanlan.zhihu.com/p/32052223](https://zhuanlan.zhihu.com/p/32052223)

[Consul](https://www.consul.io/)是一个分布式、高可用性和数据中心感知的解决方案，用于跨动态、分布式基础设施连接和配置应用程序。

Consul 提供了一个控制平面（control plane），使你能够注册、访问和安全的部署在你的网络上的服务。控制平面是网络基础设施的一部分，它维护一个中央注册中心来跟踪服务及其各自的 IP 地址。

当使用 Consul 的服务网格功能时，Consul 在请求路径中动态配置边车（sidecar）和网关代理，这使你能够授权服务到服务的连接，将请求路由到健康的服务实例，并在不修改服务代码的情况下强制 mTLS 加密。这可以确保通信的性能和可靠性。

![截屏2025-04-15 10.40.31.png](%E5%BE%AE%E6%9C%8D%E5%8A%A1%20%E4%B8%8E%20%E4%BA%91%E5%8E%9F%E7%94%9F%201cf3895181d28010a46afdc8337886b8/%E6%88%AA%E5%B1%8F2025-04-15_10.40.31.png)

**特性**

- 多数据中心——Consul 被构建为数据中心感知的，可以支持任意数量的区域，而不需要复杂的配置。
- 服务网格/服务细分——使用自动 TLS 加密和基于身份的授权实现安全的服务对服务通信。应用程序可以在服务网格配置中使用 sidecar 代理为入站和出站连接建立 TLS 连接，而完全不知道连接方。
- 服务发现——Consul 使服务注册自己和通过 DNS 或 HTTP 接口发现其他服务变得简单。也可以注册外部服务，如 SaaS 提供者。
- 健康检查——健康检查使得 Consul 能够就集群中的任何问题快速向操作员发出警报。与服务发现的集成可以防止将流量路由到不健康的主机，并启用服务级断路器。
- Key/Value存储——一个灵活的键/值存储允许存储动态配置、特性标记、协调、领导选举等。简单的 HTTP API 使得在任何地方都可以轻松使用。

**数据中心**

Consul 控制平面包含一个或多个数据中心。数据中心是执行基本 consul 操作的 consul 基础设施的最小单元。一个数据中心至少包含一个 consul 服务器代理（server agent），但是一个实际部署包含三个或五个服务器代理（server agent）和几个 consul 客户机代理（client agent）。你可以创建多个数据中心，并允许不同数据中心中的节点相互交互。

**集群**

Consul 代理之间相互联通的集合称为集群。数据中心和集群的概念通常可以互换使用。但是，在某些情况下，集群只引用 ConsulServer 代理，比如在 HCP Consull 中。在其他上下文中，例如 Consul Enterprise 包含的管理分区特性，集群可能引用客户机代理的集合。

**代理（Agents）**

你可以运行 Consul 二进制文件来启动 consul 代理，这些代理是实现 consul 控制平面功能的守护进程。代理可以作为服务器或客户端启动。

**server agent**

Consul 服务器代理存储所有状态信息，包括服务和节点 IP 地址、健康检查和配置。建议在一个集群中部署三到五台服务器。部署的服务器越多，出现故障时的弹性和可用性就越大。然而，更多的服务器会降低一致性，而一致性是一个关键的服务器功能，它使 Consul 能够高效地处理信息。

**共识协议**

Consul 集群通过一个称为“一致同意”的过程选出一个服务器作为 leader。leader 处理所有查询和事务，这可以防止包含多个服务器的集群中发生冲突的更新。

当前不充当群集 leader 的服务器称为 follower。follower 从客户机代理向集群 leader 转发请求。leader 将请求复制到集群中的所有其他服务器。复制能确保如果 leader 不可用，群集中的其他服务器可以选择另一个 leader 而不丢失任何数据。

Consul 服务器在端口8300上使用 Raft 算法建立共识。有关 Raft 算法介绍，可以参考[https://thesecretlivesofdata.com/raft](https://thesecretlivesofdata.com/raft/)。

**client agent**

Consul 客户向 consul 集群报告节点和服务状态。在典型的部署中，必须在数据中心的每个计算节点上运行客户端代理。客户端使用远程过程调用(RPC)与服务器交互。默认情况下，客户机向端口`8300`上的服务器发送 RPC 请求。

对于可以使用的客户机代理或服务的数量没有限制，但是生产部署应该跨多个 Consul 数据中心分发服务。使用多数据中心部署增强了基础设施的弹性，并限制了控制平面出现问题。建议每个数据中心最多部署5,000个客户机代理。一些大型组织已经在一个多数据中心部署中部署了成千上万的客户机代理和成千上万的服务实例。

**LAN gossip pool**

客户机和服务器代理参与 LAN gossip pool，以便分发和执行节点健康检查。池中的代理将健康检查信息传播到整个群集。代理在`8301`端口上使用 UDP 进行 gossip 通信。如果 UDP 不可用，代理的 gossip 就会退而求其次使用 TCP。

### go sdk

连接consul

```go
import "github.com/hashicorp/consul/api"
// 连接到consul
	cc, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		log.Fatalf("failed to conn consul")
	}
	// 获取本机的出口ip
	ipinfo, err := GetOutboundIP()
	// 将我们的gprc服务注册到consul
	// 1.定义服务
		srv := &api.AgentServiceRegistration{
		ID:      fmt.Sprintf("%s-%s-%d", serviceName, ipinfo.String(), 8976), // 服务唯一ID
		Name:    serviceName,
		Tags:    []string{"ysh"},
		Address: ipinfo.String(),
		Port:    8976,
	}
	// 2. 注册到consun
	cc.Agent().ServiceRegister(srv)
	
	// GetOutboundIP 获取本机的出口IP
func GetOutboundIP() (net.IP, error) {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP, nil
}

```

**健康检查**

```go
import "google.golang.org/grpc/health"
import healthpb "google.golang.org/grpc/health/grpc_health_v1"

// 给我们的gprc的服务增减增加注册健康检查
	healthpb.RegisterHealthServer(s, health.NewServer())// consul 发来健康检查的RPC请求，这个负责返回OK
	
// 配置健康检查
check := &api.AgentServiceCheck{
		GRPC:                           fmt.Sprintf("%s:%d", ipinfo.String(), 8976), // 外网地址
		Timeout:                        "5s",
		Interval:                       "5s",  // 间隔
		DeregisterCriticalServiceAfter: "10s", // 10秒钟后注销掉不健康的服务节点
	}
	
	srv := &api.AgentServiceRegistration{
		ID:      fmt.Sprintf("%s-%s-%d", serviceName, ipinfo.String(), 8976), // 服务唯一ID
		Name:    serviceName,
		Tags:    []string{"ysh"},
		Address: ipinfo.String(),
		Port:    8976,
		Check:   check,
	}
```

<aside>
💡

你的 Go 服务
|
| 1. 创建 Consul 客户端（127.0.0.1:8500）
v
cc := api.NewClient(...)
|
| 2. 定义服务 + 健康检查
v
srv := &api.AgentServiceRegistration{...}
|
| 3. 注册服务
v
cc.Agent().ServiceRegister(srv)
|
|==> 发起 HTTP 请求到 Consul 的 Agent API
|==> Consul 注册服务，开始定时健康检查

</aside>

**服务发现**

```go
import _ "github.com/mbobakov/grpc-consul-resolver"

// 1. 连接到consul
		cc, err := api.NewClient(api.DefaultConfig())
		if err != nil{
			fmt.Printf("conn consul failed:%v\n", err)
			return
		}

		// 2. 根据服务名称查询实例
		// cc.Agent().Services()  // 列出所有的
		serviceMap, err := cc.Agent().ServicesWithFilter("Service=`hello`")// 查询服务名称是hello的所有服务节点
		if err != nil {
			fmt.Printf("query `hello` service failed:%v\n",err)
			return
		}
		var addr string
		for k, v := range serviceMap {
			fmt.Printf("%s:%#v\n", k, v)
			addr = fmt.Sprintf("%s:%d", v.Address, v.Port)
			break
		}
		conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
		
//或者下面，下面的最好！！！
conn, err := grpc.NewClient("consul://localhost:8500/hello", // grpc中使用consul名称解析器，
		grpc.WithTransportCredentials(insecure.NewCredentials()))
```

**注销服务**

```go
// ctrl + c 退出程序
	quitCh := make(chan os.Signal, 1)
	signal.Notify(quitCh, syscall.SIGTERM, syscall.SIGINT)
	fmt.Println("wait quit signal...")
	<-quitCh // 没收到信号就阻塞

	// 程序退出，注销服务
	fmt.Println("service quit...")
	err = cc.Agent().ServiceDeregister(serviceID) // 注销服务
	if err != nil{
		fmt.Printf("service deregister failed: %v\n", err)
	}
```

**负载均衡**

```go
conn, err := grpc.NewClient("consul://localhost:8500/hello", // grpc中使用consul名称解析器，
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"round_robin"}`),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
```

## 函数选项模式

```go
package main

import "fmt"

type ServiceConfig struct {
	A string
	B string
	C int

	X struct{}
	Y Info
}

type Info struct {
	addr string
}

// NewServiceConfig 创建一个ServiceConfig的函数
func NewServiceConfig(a, b string, c int) *ServiceConfig {
	return &ServiceConfig{
		A: a,
		B: b,
		C: c,
	}
}

const defaultValueC = 1

// target: 想要A和B必须传入，C可以不传，不传就用默认值
func NewServiceConfig2(a, b string, c ...int) *ServiceConfig {
	valueC := defaultValueC
	if len(c) > 0 {
		valueC = c[0]
	}
	return &ServiceConfig{
		A: a,
		B: b,
		C: valueC,
	}
}

// Options模式
type FuncServiceConfigOption func(*ServiceConfig)

func NewServiceConfig3(a, b string, opst ...FuncServiceConfigOption) *ServiceConfig {
	sc := &ServiceConfig{
		A: a,
		B: b,
		C: defaultValueC,
	}
	// 针对可能传进来的FuncServiceConfigOption参数做处理
	for _, opt := range opst {
		opt(sc)
	}
	return sc
}

func WithC(c int) FuncServiceConfigOption {
	return func(sc *ServiceConfig) {
		sc.C = c
	}
}

func WithY(info Info) FuncServiceConfigOption {
	return func(sc *ServiceConfig) {
		sc.Y = info
	}
}

// 上面方法存在一些问题
// 我可以直接通过sc.C来修改，那我的WithC有什么用呢
// 不想被直接修改，只能通过我们提供的方法来进行修改

type config struct {
	name string
	age  int
}

const defaultName = "ysh"

func NewConfig(age int, opts ...ConfigOption) *config {
	cfg := &config{
		age:  age,
		name: defaultName,
	}
	for _, opt := range opts {
		opt.apply(cfg)
	}
	return cfg
}

type ConfigOption interface {
	apply(*config)
}

type funcConfigOption struct {
	f func(*config)
}

func (f funcConfigOption) apply(c *config) {
	f.f(c)
}

func NewfuncConfigOption(f func(*config)) funcConfigOption {
	return funcConfigOption{f: f}
}

func WithName(name string) ConfigOption {
	return NewfuncConfigOption(func(c *config) { c.name = name })
}

func main() {
	//info := Info{addr: "127.0.0.1"}
	//sc := NewServiceConfig3("ysh", "py", WithC(10), WithY(info))
	//fmt.Printf("sc:%#v\n", sc)

	cfg := NewConfig(18)
	fmt.Printf("cfg:%#v\n", cfg)
	cfg2 := NewConfig(18, WithName("张三"))
	fmt.Printf("cfg:%#v\n", cfg2)
}
```

<aside>
💡

## 一、👓 什么是接口（Go Interface）

在 Go 里，**接口是一组方法的集合**，只要一个类型实现了接口里所有方法，就隐式实现了这个接口。

✅ 接口的作用：

就是让你写通用代码，可以接收多种实现，只要它们都“符合某个行为”。

对比项 | 函数式写法（简单版） | 接口 + struct 写法（扩展版）
参数类型 | func(*config) | ConfigOption 接口
写法简洁 | ✅ 更简单 | ❌ 需要定义接口和结构体
调用时易读性 | 一般（传函数） | 很好（像自然语言）
支持组合、缓存、条件逻辑 | ❌ 难实现 | ✅ 支持封装逻辑，能做组合
支持多种 Option 类型（如 JSON、Env） | ❌ 几乎不行 | ✅ 很自然
对调用者隐藏内部实现细节 | ❌ 直接暴露函数 | ✅ 只暴露接口，解耦调用者
易测试、Mock、封装 | ❌ 比较难 | ✅ 容易扩展封装
适合小项目 | ✅ | ❌ 过度设计
适合中大型项目 / 框架 | ❌ | ✅ 标准实践

</aside>

> 👉 接口 + struct 写法的本质：
> 
> 
> 把“函数行为”变成“可组合、可扩展、可抽象的对象”。
> 

它适合那些需要**灵活配置、可扩展、构建框架或组件库**的地方。

# Go kit

使用 Go kit 构建的服务分为三层：

1. 传输层（Transport layer）
2. 端点层（Endpoint layer）
3. 服务层（Service layer）

请求在第1层进入服务，向下流到第3层，响应则相反。

**Transports**

```go
"httptransport "github.com/go-kit/kit/transport/http""
```

传输域绑定到具体的传输协议，如 HTTP 或 gRPC。在一个微服务可能支持一个或多个传输协议的世界中，这是非常强大的：你可以在单个微服务中支持原有的 HTTP API 和新增的 RPC 服务。

**Endpoints**

```go
"github.com/go-kit/kit/endpoint"
```

端点就像控制器上的动作/处理程序; 它是安全性和抗脆弱性逻辑的所在。如果实现两种传输(HTTP 和 gRPC) ，则可能有两种将请求发送到同一端点的方法。

**Services**

服务（指Go kit中的service层）是实现所有业务逻辑的地方。服务层通常将多个端点粘合在一起。在 Go kit 中，服务层通常被抽象为接口，这些接口的实现包含业务逻辑。Go kit 服务层应该努力遵守整洁架构或六边形架构。也就是说，业务逻辑不需要了解端点（尤其是传输域）概念：你的服务层不应该关心HTTP 头或 gRPC 错误代码。

**Middlewares**

Go kit 试图通过使用中间件（或装饰器）模式来执行严格的关注分离（separation of concerns）。中间件可以包装端点或服务以添加功能，比如日志记录、速率限制、负载平衡或分布式跟踪。围绕一个端点或服务链接多个中间件是很常见的。

Go kit 试图通过精心使用中间件（或装饰器）模式来强制执行严格的关注分离（**separation of concerns**）。
# 链路追踪

借助分布式链路追踪能够帮助我们在复杂的分布式系统中快速定位问题、排除故障。

把请求链路中进行的每个网络调用都会被捕获并表示为一个跨度。

分布式链路追踪工具将唯一的链路追踪上下文（trace ID）插入到每个请求的标头中，并借助各种实现工具确保链路追踪上下文在整个请求链路中传播。

# **OpenTelemetry**

对于 OpenTelemetry，不同的角色（Dev和Ops）侧重点也不尽相同。

如果你的角色是Dev，那么你可能更关注如何通过编写代码使程序获得可观测性。

如果你的角色是Ops，那么你可能更关注如何从多个服务中收集 traces、metrics 和 logs 数据，并将它们发送到可观测后台。

OpenTelemetry，也称为 OTel，是一个与供应商无关的开源可观测性框架，用于检测、生成、收集和导出遥测数据，如链路追

## **相关概念**

OpenTelemetry 的目的是收集、处理和输出信号。 信号是用于描述操作系统和平台上运行的应用程序基本活动的系统输出。信号可以是你想在特定时间点测量的东西，如温度或内存使用率，也可以是你想追踪的分布式系统组件中发生的事件。你可以将不同的信号组合在一起，从不同角度观察同一项技术的内部运作。

OpenTelemetry 目前支持 [traces](https://opentelemetry.io/docs/concepts/signals/traces/)、[metrics](https://opentelemetry.io/docs/concepts/signals/metrics/)、[logs](https://opentelemetry.io/docs/concepts/signals/logs/)和[baggage](https://opentelemetry.io/docs/concepts/signals/baggage/)。

- **trace**：在分布式应用程序中的完整请求链路信息。
- **指标**是在运行时捕获的服务指标，应用程序和请求指标是可用性和性能的重要指标。
- **日志**是系统或应用程序在特定时间点发生的事件的文本记录。
- **baggage**：是在信号之间传递的上下文信息。

**OTLP**

[开放遥测协议（OTLP）](https://opentelemetry.io/docs/specs/otlp/)规范描述了遥测源、中间节点（如采集器和遥测后端）之间的遥测数据编码、传输和交付机制。

OTLP 是在 OpenTelemetry 项目范围内设计的通用遥测数据传输协议。

**Collector**

OpenTelemetry Collector 提供了一种与供应商无关的接收、处理和导出遥测数据的方式。它无需运行、操作和维护多个代理/收集器。它具有更好的可扩展性，支持向一个或多个开源或商业后端发送开源可观测性数据格式（如 Jaeger、Prometheus、Fluent Bit 等）。本地收集器代理是仪器库导出遥测数据的默认位置。

**可观测后台**

[Jaeger](https://www.jaegertracing.io/) 和 [Zipkin](https://zipkin.io/) 是社区中比较流行的方案，他们都提供有可视化的WebUI方便查询。

# **OpenTelemetry Go**

1. 业务代码
2. 添加**OpenTelemetry instrumentation**
    
    ```go
    go get "go.opentelemetry.io/otel" \
      "go.opentelemetry.io/otel/exporters/stdout/stdoutmetric" \
      "go.opentelemetry.io/otel/exporters/stdout/stdouttrace" \
      "go.opentelemetry.io/otel/propagation" \
      "go.opentelemetry.io/otel/sdk/metric" \
      "go.opentelemetry.io/otel/sdk/resource" \
      "go.opentelemetry.io/otel/sdk/trace" \
      "go.opentelemetry.io/otel/semconv/v1.24.0" \
      "go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
    // 这里安装的是 OpenTelemety SDK 组件和 net/http 测量仪器。如果要对不同的库进行网络请求检测，则需要安装相应的仪器库。
    ```
    
3. **初始化OpenTelemetry SDK**
    
    ```go
    // otel.go
    
    package main
    
    import (
    	"context"
    	"errors"
    	"time"
    
    	"go.opentelemetry.io/otel"
    	"go.opentelemetry.io/otel/exporters/stdout/stdoutmetric"
    	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
    	"go.opentelemetry.io/otel/propagation"
    	"go.opentelemetry.io/otel/sdk/metric"
    	"go.opentelemetry.io/otel/sdk/trace"
    )
    
    // setupOTelSDK 引导 OpenTelemetry pipeline。
    // 如果没有返回错误，请确保调用 shutdown 进行适当清理。
    func setupOTelSDK(ctx context.Context) (shutdown func(context.Context) error, err error) {
    	var shutdownFuncs []func(context.Context) error
    
    	// shutdown 会调用通过 shutdownFuncs 注册的清理函数。
    	// 调用产生的错误会被合并。
    	// 每个注册的清理函数将被调用一次。
    	shutdown = func(ctx context.Context) error {
    		var err error
    		for _, fn := range shutdownFuncs {
    			err = errors.Join(err, fn(ctx))
    		}
    		shutdownFuncs = nil
    		return err
    	}
    
    	// handleErr 调用 shutdown 进行清理，并确保返回所有错误信息。
    	handleErr := func(inErr error) {
    		err = errors.Join(inErr, shutdown(ctx))
    	}
    
    	// 设置传播器
    	prop := newPropagator()
    	otel.SetTextMapPropagator(prop)
    
    	// 设置 trace provider.
    	tracerProvider, err := newTraceProvider()
    	if err != nil {
    		handleErr(err)
    		return
    	}
    	shutdownFuncs = append(shutdownFuncs, tracerProvider.Shutdown)
    	otel.SetTracerProvider(tracerProvider)
    
    	// 设置 meter provider.
    	meterProvider, err := newMeterProvider()
    	if err != nil {
    		handleErr(err)
    		return
    	}
    	shutdownFuncs = append(shutdownFuncs, meterProvider.Shutdown)
    	otel.SetMeterProvider(meterProvider)
    
    	return
    }
    
    func newPropagator() propagation.TextMapPropagator {
    	return propagation.NewCompositeTextMapPropagator(
    		propagation.TraceContext{},
    		propagation.Baggage{},
    	)
    }
    
    func newTraceProvider() (*trace.TracerProvider, error) {
    	traceExporter, err := stdouttrace.New(
    		stdouttrace.WithPrettyPrint())
    	if err != nil {
    		return nil, err
    	}
    
    	traceProvider := trace.NewTracerProvider(
    		trace.WithBatcher(traceExporter,
    			// 默认为 5s。为便于演示，设置为 1s。
    			trace.WithBatchTimeout(time.Second)),
    	)
    	return traceProvider, nil
    }
    
    func newMeterProvider() (*metric.MeterProvider, error) {
    	metricExporter, err := stdoutmetric.New()
    	if err != nil {
    		return nil, err
    	}
    
    	meterProvider := metric.NewMeterProvider(
    		metric.WithReader(metric.NewPeriodicReader(metricExporter,
    			// 默认为 1m。为便于演示，设置为 3s。
    			metric.WithInterval(3*time.Second))),
    	)
    	return meterProvider, nil
    }
    
    ```
    
4. **测量 HTTP server**
    
    ```go
    // main.go 
    func newHTTPHandler() http.Handler {
    	mux := http.NewServeMux()
    
    	// handleFunc 是 mux.HandleFunc 的替代品，。
    	// 它使用 http.route 模式丰富了 handler 的 HTTP 测量
    	handleFunc := func(pattern string, handlerFunc func(http.ResponseWriter, *http.Request)) {
    		// 为 HTTP 测量配置 "http.route"。
    		handler := otelhttp.WithRouteTag(pattern, http.HandlerFunc(handlerFunc))
    		mux.Handle(pattern, handler)
    	}
    
    	// Register handlers.
    	handleFunc("/roll", roll)
    
    	// 为整个服务器添加 HTTP 测量。
    	handler := otelhttp.NewHandler(mux, "/")
    	return handler
    }
    
    func main() {
    	if err := run(); err != nil {
    		log.Fatalln(err)
    	}
    }
    
    func run() (err error) {
    	// 平滑处理 SIGINT (CTRL+C) .
    	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
    	defer stop()
    
    	// 设置 OpenTelemetry.
    	otelShutdown, err := setupOTelSDK(ctx)
    	if err != nil {
    		return
    	}
    	// 妥善处理停机，确保无泄漏
    	defer func() {
    		err = errors.Join(err, otelShutdown(context.Background()))
    	}()
    
    	// 启动 HTTP server.
    	srv := &http.Server{
    		Addr:         ":8080",
    		BaseContext:  func(_ net.Listener) context.Context { return ctx },
    		ReadTimeout:  time.Second,
    		WriteTimeout: 10 * time.Second,
    		Handler:      **newHTTPHandler(),**
    	}
    	srvErr := make(chan error, 1)
    	go func() {
    		srvErr <- srv.ListenAndServe()
    	}()
    
    	// 等待中断.
    	select {
    	case err = <-srvErr:
    		// 启动 HTTP 服务器时出错.
    		return
    	case <-ctx.Done():
    		// 等待第一个 CTRL+C.
    		// 尽快停止接收信号通知.
    		stop()
    	}
    
    	// 调用 Shutdown 时，ListenAndServe 会立即返回 ErrServerClosed。
    	err = srv.Shutdown(context.Background())
    	return
    }
    ```
    
5. **添加自定义测量**
    
    ```go
    // 业务逻辑文件
    "go.opentelemetry.io/otel"
    	"go.opentelemetry.io/otel/attribute"
    	"go.opentelemetry.io/otel/metric"
    var (
    	tracer  = otel.Tracer("roll") )
    	
    	ctx, span := tracer.Start(r.Context(), "roll") // 开始 span
    	defer span.End()                               // 结束 span
    	
    	// 添加span属性
    	rollValueAttr := attribute.Int("roll.value", number)
    	span.SetAttributes(rollValueAttr) // span 添加属性
    ```
    

## **将链路追踪数据发送至 Jaeger**

Jaeger 官方提供的 **all-in-one** 是为快速本地测试而设计的可执行文件。它包括 **Jaeger UI**、**jaeger-collector**、**jaeger-query** 和 **jaeger-agent**，以及一个内存存储组件。

```go
docker run --rm --name jaeger \
  -e COLLECTOR_ZIPKIN_HOST_PORT=:9411 \
  -p 6831:6831/udp \
  -p 6832:6832/udp \
  -p 5778:5778 \
  -p 16686:16686 \
  -p 4317:4317 \
  -p 4318:4318 \
  -p 14250:14250 \
  -p 14268:14268 \
  -p 14269:14269 \
  -p 9411:9411 \
  jaegertracing/all-in-one:1.55

//  http://localhost:16686 
```

| Port | Protocol | Component | Function |
| --- | --- | --- | --- |
| 6831 | UDP | agent | accept jaeger.thrift over Thrift-compact protocol (used by most SDKs) |
| 6832 | UDP | agent | accept jaeger.thrift over Thrift-binary protocol (used by Node.js SDK) |
| 5775 | UDP | agent | (deprecated) accept zipkin.thrift over compact Thrift protocol (used by legacy clients only) |
| 5778 | HTTP | agent | serve configs (sampling, etc.) |
| 16686 | HTTP | query | serve frontend |
| 4317 | HTTP | collector | accept OpenTelemetry Protocol (OTLP) over gRPC |
| 4318 | HTTP | collector | accept OpenTelemetry Protocol (OTLP) over HTTP |
| 14268 | HTTP | collector | accept jaeger.thrift directly from clients |
| 14250 | HTTP | collector | accept model.proto |
| 9411 | HTTP | collector | Zipkin compatible endpoint (optional) |

 HTTP 协议的`4318` 端口上报链路追踪数据。

1. **上报至 Jaeger**
    
    ```go
    go get go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp
    
    // otel.go
    
    func newJaegerTraceProvider(ctx context.Context) (*trace.TracerProvider, error) {
    	// 创建一个使用 HTTP 协议连接本机Jaeger的 Exporter
    	traceExporter, err := otlptracehttp.New(ctx,
    		otlptracehttp.WithEndpoint("127.0.0.1:4318"),
    		otlptracehttp.WithInsecure())
    	if err != nil {
    		return nil, err
    	}
    	traceProvider := trace.NewTracerProvider(
    		trace.WithBatcher(traceExporter,
    			// 默认为 5s。为便于演示，设置为 1s。
    			trace.WithBatchTimeout(time.Second)),
    	)
    	return traceProvider, nil
    }
    
    // 修改设置 trace provider 部分
    // 设置 trace provider.
    //tracerProvider, err := newTraceProvider()
    tracerProvider, err := newJaegerTraceProvider(ctx)
    ```
    

# **Jaeger**

分布式追踪可观测平台（如 Jaeger）对于架构为微服务的现代软件应用程序至关重要。Jaeger 可以映射分布式系统中的请求流和数据流。这些请求可能会调用多个服务，而这些服务可能会带来各自的延迟或错误。Jaeger 将这些不同组件之间的点连接起来，帮助识别性能瓶颈、排除故障并提高整体应用程序的可靠Jaeger是100%开源、云原生、可无限扩展的。

[Jaeger](https://www.jaegertracing.io/) 是一个分布式追踪系统它可以用于监控基于微服务的分布式系统：

- 分布式上下文传递
- 分布式事务监听
- 根因分析
- 服务依赖性分析
- 性能/延迟优化

Jaeger 项目和 OpenTelemetry 项目有着不同的目标。OpenTelemetry 的目标是提供多种语言的应用程序接口（API）和 SDK，允许应用程序将各种遥测数据输出到任意数量的度量和跟踪后端。Jaeger 项目主要是跟踪后端，它接收跟踪遥测数据，并对数据进行处理、汇总、数据挖掘和可视化。

```go
docker run --rm --name jaeger \
  -e COLLECTOR_ZIPKIN_HOST_PORT=:9411 \
  -p 6831:6831/udp \
  -p 6832:6832/udp \
  -p 5778:5778 \
  -p 16686:16686 \
  -p 4317:4317 \
  -p 4318:4318 \
  -p 14250:14250 \
  -p 14268:14268 \
  -p 14269:14269 \
  -p 9411:9411 \
  jaegertracing/all-in-one:1.55
```

```go
// 常量定义
const (
	serviceName    = "Go-Jaeger-Demo"  // 服务名称（在Jaeger UI中显示）
	jaegerEndpoint = "127.0.0.1:4318"  // Jaeger OTLP HTTP 接收端地址（默认端口4318）
)

// setupTracer 初始化 OpenTelemetry TracerProvider
// 返回一个 shutdown 函数，用于优雅关闭 tracer

	// 创建 Jaeger Trace Provider
	tracerProvider, err := newJaegerTraceProvider(ctx)
	
	// 设置全局 TracerProvider（后续调用 otel.Tracer 会使用这个 provider）
	otel.SetTracerProvider(tracerProvider)
	
	// 返回 shutdown 函数，用于程序退出时关闭 tracer
	return tracerProvider.Shutdown, nil

// newJaegerTraceProvider 创建并配置 Jaeger Trace Provider
	// 创建 OTLP HTTP Exporter，用于将 trace 数据发送到 Jaeger
	// WithInsecure() 表示不使用 TLS（本地测试用）
	exp, err := otlptracehttp.New(ctx,
		otlptracehttp.WithEndpoint(jaegerEndpoint),  // Jaeger 地址
		otlptracehttp.WithInsecure())                // 禁用 TLS（生产环境不建议）
		
		// 创建 Resource（用于标识服务，Jaeger UI 会显示这个信息）
		res, err := resource.New(ctx, resource.WithAttributes(
		semconv.ServiceName(serviceName),  // 设置服务名
	))
	
	// 创建 TracerProvider 并配置：
	traceProvider := traceSDK.NewTracerProvider(
		traceSDK.WithResource(res),                   // 绑定 Resource
		traceSDK.WithSampler(traceSDK.AlwaysSample()), // 采样策略（这里设置为全部采样，适合测试）
		traceSDK.WithBatcher(exp, traceSDK.WithBatchTimeout(time.Second)), // 批量发送，间隔1秒
	)
	
// testTracer 演示如何创建 Span（追踪链路）
// 它会创建一个父 Span 和多个子 Span，形成调用链
		// 获取 Tracer 实例（"test-tracer" 是自定义名称）
		tracer := otel.Tracer("test-tracer")
		// 定义 Span 的 Attributes（附加信息，会在 Jaeger 中显示）
	baseAttrs := []attribute.KeyValue{
		attribute.String("domain", "yuanshuhao.com"),  // 字符串类型
		attribute.Bool("plagiarize", false),          // 布尔类型
		attribute.Int("code", 7),                     // 整数类型
	}
	
	// 创建父 Span（"parent-span" 是 Span 名称）
	// 返回的 ctx 会包含这个 Span 的上下文，用于创建子 Span
	ctx, span := tracer.Start(ctx, "parent-span", trace.WithAttributes(baseAttrs...))
	
	// 确保 Span 在函数结束时结束（否则数据不会发送）
	defer span.End()
	// 创建 10 个子 Span，模拟多个操作
	for i := range 10 { // Go 1.22+ 的循环语法
		// 创建子 Span（名称格式 "span-0" ~ "span-9"）
		_, iSpan := tracer.Start(ctx, fmt.Sprintf("span-%d", i))
		
		// 模拟耗时操作（随机 sleep 0~100ms）
		time.Sleep(time.Duration(rand.Int63n(100)) * time.Millisecond)
		
		// 结束子 Span
		iSpan.End()
	}
	fmt.Println("done!")
}

// 创建根 Context
	ctx := context.Background()
	
	// 初始化 Tracer
	shutdown, err := setupTracer(ctx)
	if err != nil {
		panic(err)
	}
	
	// 程序退出时关闭 Tracer（确保数据发送完成）
	defer func() {
		_ = shutdown(ctx)
	}()

	// 运行测试函数，生成 Trace 数据
	testTracer(ctx)	
	
	// http://localhost:16686
```

# **基于OTel的HTTP链路追踪**

```go
const (
	serviceName     = "httpclient-Demo"
	peerServiceName = "blog"
	jaegerEndpoint  = "127.0.0.1:4318"
	blogURL         = "https://yuanshuhao.com"
)

// newJaegerTraceProvider 创建一个 Jaeger Trace Provider
func newJaegerTraceProvider(ctx context.Context) (*sdktrace.TracerProvider, error) {
	// 创建一个使用 HTTP 协议连接本机Jaeger的 Exporter
	exp, err := otlptracehttp.New(ctx,
		otlptracehttp.WithEndpoint(jaegerEndpoint),
		otlptracehttp.WithInsecure())
	if err != nil {
		return nil, err
	}

	res, err := resource.New(ctx, resource.WithAttributes(semconv.ServiceName(serviceName)))
	if err != nil {
		return nil, err
	}
	// 创建 Provider
	traceProvider := sdktrace.NewTracerProvider(
		sdktrace.WithResource(res),
		sdktrace.WithSampler(sdktrace.AlwaysSample()), // 采样
		sdktrace.WithBatcher(exp, sdktrace.WithBatchTimeout(time.Second)),
	)
	return traceProvider, nil
}

// initTracer 初始化 Tracer
func initTracer(ctx context.Context) (*sdktrace.TracerProvider, error) {
	tp, err := newJaegerTraceProvider(ctx)
	if err != nil {
		return nil, err
	}

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(
		propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}),
	)
	return tp, nil
}
```

```go
tr := otel.Tracer("http-client")

	ctx, span := tr.Start(ctx, "GET BLOG", trace.WithAttributes(semconv.PeerService(peerServiceName)))
	defer span.End()

	// 创建一个 http client,带有链路追踪的配置
	client := http.Client{Transport: otelhttp.NewTransport(http.DefaultTransport)}

	// 构造请求
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, blogURL, nil)

	// 发起请求
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
```

深入`net/http`内部的追踪可以使用`net/http/httptrace`，会采集`dns`、`connect`、`tls`等环节

```go
// 创建 http client，配置trace
	client := http.Client{
		Transport: otelhttp.NewTransport(
			http.DefaultTransport,
			otelhttp.WithClientTrace(func(ctx context.Context) *httptrace.ClientTrace {
				return otelhttptrace.NewClientTrace(ctx)
			}),
		),
	}

```

**gin框架Jaeger示例**

```go
go get go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin

// 设置 otelgin 中间件
	r.Use(otelgin.Middleware(serviceName))

	// 在响应头记录 TRACE-ID
	r.Use(func(c *gin.Context) {
		c.Header("Trace-Id", trace.SpanFromContext(c.Request.Context()).SpanContext().TraceID().String())
	})

```

## **gRPC的链路追踪**

```go
go get go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc
```

1. server
    
    ```go
    s := grpc.NewServer(
      grpc.StatsHandler(otelgrpc.NewServerHandler()), // 设置 StatsHandler
    )
    ```
    
2. client
    
    ```go
    conn, err := grpc.NewClient(
      *addr,
      grpc.WithTransportCredentials(insecure.NewCredentials()),
      grpc.WithStatsHandler(otelgrpc.NewClientHandler()), // 设置 StatsHandler
    )
    ```
    

## **GORM配置链路追踪**

```go
go get gorm.io/plugin/opentelemetry/tracing
```

在初始化 gorm.DB 之后，通过安装插件的方式引入 tracing 和 metrics 。

```go
// 1.
if err := db.Use(tracing.NewPlugin()); err != nil {
		panic(err)
	}
// 2. 如果只想采集 tracing 数据，可以按如下方式配置插件。
if err := db.Use(tracing.NewPlugin(tracing.WithoutMetrics())); err != nil {
		panic(err)
	}
	
span.RecordError(err) // 记录error
span.SetStatus(codes.Error, err.Error())
```

## **go-redis配置链路追踪**

```go
go get github.com/redis/go-redis/extra/redisotel/v9
// 确保下载的 redisotel 版本与你使用的go-redis版本一致

// 启用 tracing
if err := redisotel.InstrumentTracing(rdb); err != nil {
	panic(err)
}
```

## zap**日志库配置链路追踪**

```go
go get github.com/uptrace/opentelemetry-go-extra/otelzap
```

```go
// 创建 logger
	logger := otelzap.New(
		zap.NewExample(),                    // zap实例，按需配置
		otelzap.WithMinLevel(zap.InfoLevel), // 指定日志级别
		**otelzap.WithTraceIDField(true),**      // 在日志中记录 traceID
	)
	defer logger.Sync()

	// 替换全局的logger
	undo := otelzap.ReplaceGlobals(logger)
	defer undo()

	otelzap.L().Info("replaced zap's global loggers")        // 记录日志
	otelzap.Ctx(context.TODO()).Info("... and with context") // 从ctx中获取traceID并记录
}

**otelzap.WithTraceIDField(true), 
// github.com/uptrace/opentelemetry-go-extra/otelzap v0.2.4 // indirect
才有，往上就没有了**
```

## kratos、zero链路追踪

[https://go-kratos.dev/docs/component/middleware/tracing](https://go-kratos.dev/docs/component/middleware/tracing)

[https://go-zero.dev/docs/tutorials/monitor/index?_highlight=链路#链路追踪](https://go-zero.dev/docs/tutorials/monitor/index?_highlight=%E9%93%BE%E8%B7%AF#%E9%93%BE%E8%B7%AF%E8%BF%BD%E8%B8%AA)
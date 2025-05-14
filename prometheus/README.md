# Prometheus

prometheus 是⽬前主流的⼀个开源监控系统和告警⼯具包，它可以与 Kubernetes 等现代基础设施平台配合，轻松集成到云原⽣环境中，提供对容器化应⽤、微服务架构等的全⾯监控。本

Prometheus 收集并存储其指标作为时间序列数据，即指标信息与其记录的时间戳⼀起存储，同时存储的还有可选的称为标签的键值对。

prometheus 的主要特性有：
多维数据模型，包含由 metric 名称和键值对标识的时间序列数据
PromQL，⼀种可以灵活利⽤上述维度数据的查询语⾔
不依赖于分布式存储; 单个服务器节点是⾃治的
通过基于 HTTP 的拉模式（pull）进⾏时间序列数据收集
可以通过⼀个中间⽹关（Pushgateway）以推模式上报时间序列数据
通过服务发现或静态配置发现监控⽬标
⽀持多种模式的图表和仪表板

**metric**

，metric 就是⽤数字来测量/度量。时间序列⼀词指的是记录⼀段时间内的变化。⽤户想要测量的内容会因应⽤⽽异。对于 Web 服务器，可以测量请求耗时；对于数据库，可以测量活动连接数或活动查询数等等。

Metrics 在理解应⽤程序以某种⽅式⼯作的原因⽅⾯发挥着重要作⽤。假设有⼀个 Web 应⽤程序正在运⾏，你发现它的运⾏速度很慢。这时候你需要⼀些信息来了解应⽤程序的运⾏情况。例如，当请求数量较多时，应⽤程序可能会变慢。如果掌握了请求数指标，就可以确定原因并增加服务器数量来处理负载。

Prometheus 的⽣态系统由多个组成部分组成，其中许多是可选的：
prometheus Server：⽤于抓取和存储时间序列数据
client libraries：⽤于检测应⽤程序代码的客户端库
push gateway：⽤于短时任务的推送接收器
exporter：⽤于 HAProxy、StatsD、Graphite 等服务的专⽤输出程序
altermanager：处理告警的警报管理器
各种⽀持⼯具
⼤多数 Prometheus 组件都是使⽤ Go 编写的，因此很容易构建和部署为静态⼆进制⽂件。

Prometheus 可直接或间接通过推送⽹关（Pushgateway）抓取监控指标（适⽤于短时任务）。

它在本地存储所有抓取到的样本数据，并在这些数据上执⾏⼀系列规则，以从现有数据中汇总并记录新的时间序列或⽣成告警。

可以使⽤ Grafana 或其他 API 消费者对收集到的数据进⾏可视化展示。

Prometheus 适⽤于记录⽂本格式的时间序列数据，它既适⽤于以机器为中⼼的监控，也适⽤于⾼度动态的⾯向服务架构的监控。在微服务的世界中，它天然⽀持对多维数据的收集和查询。Prometheus 是专为提⾼系统可靠性⽽设计的，它可以协助在故障期间快速诊断问题，每个 Prometheus Server 都是相互独⽴的，不依赖于⽹络存储或其他远程服务。当基础架构出现故障时，你可以通过 Prometheus 快速定位故障点，⽽且不会消耗⼤量的基础架构资源。

从根本上说，Prometheus 将所有数据都存储为时间序列：属于同⼀指标（metric）和同⼀组标注维度（label）的带时间戳的值流。除了存储的时间序列外，Prometheus 还可以根据查询结果⽣成临时派⽣时间序列。

在时间序列中的每⼀个点称为⼀个样本（sample），样本由以下三部分组成：
指标metric：metric name 和描述当前样本特征的 label sets
时间戳timestamp：⼀个精确到毫秒的时间戳
样本值value： ⼀个 float64 的浮点型数据表示当前样本的值

## **Metric name和 label**

每个时间序列都由其指标名称和称为标签的可选键值对唯⼀标识。

### **Metric name**

指定要测量的系统的⼀般功能（例如 http_requests_total -接收的http请求总数）。

指标名称可以包含ASCII字⺟、数字、下划线和冒号。它必须匹配正则表达式 [a-zA-Z_:][a-zA-Z0-
9_:]* 。

### **Metric labels**

使 Prometheus 的维度数据模型能够识别同⼀指标名称的任何给定标签组合。它标识了该度量的特定维度实例化（例如：所有发送 POST 到 /api/tracks 的HTTP请求）。Prometheus 查询语⾔允许基于这些维度进⾏筛选和聚合。

任何标签值的更改，包括添加或删除标签，都将创建⼀个新的时间序列。

label 可以包含 ASCII 字⺟、数字以及下划线。必须匹配 [a-zA-Z_][a-zA-Z0-9_]* 。

以 __ （两个“下划线”）开头的标签名称保留供内部使⽤。
标签值可以包含任何 Unicode 字符。
标签值为空的标签被视为等同于不存在的标签。

### **Sample（采样/样本）**

样本构成实际的时间序列数据。每个样本包括
⼀个 float64值
毫秒精度的时间戳

### **Notation（表达式）**

给定⼀个指标名称和⼀组标签，时间序列通常使⽤以下符号进⾏表示：

```go
<metric name>{<label name>=<label value>, ...}
```

例如，指标名称为 api_http_requests_total ，带有 method=“POST” 和 handler="/messages" label 的时间序列可以这样写：

```go
api_http_requests_total{method="POST", handler="/messages"}
```

### **Metric 类型**

Prometheus 客户端库提供四种核⼼指标类型。这些类型⽬前仅在客户端库（以便根据特定类型的使⽤情况定制应⽤程序接⼝）和传输协议中有所区别。⽬前 prometheus 服务端还没有使⽤类型信息，⽽是将所有数据平铺为⽆类型的时间序列。未来的版本中可能会有所改变。

### **Counter（计数器）**

Counter 是⼀种累积度量，表示单个单调递增的计数器（只增不减），其值只能在重新启动时增加或重置为零。
例如，可以使⽤ counter 来表示已服务的请求数、已完成的任务数或错误数。

不要使⽤ counter 记录可能减⼩的值。例如，不要使⽤ counter 记录当前正在运⾏的进程数，⽽是应该使⽤ gauge 类型来记录。

### **Gauge（仪表盘）**

Gauge 是⼀种度量标准，代表⼀个可以任意升降的单⼀数值。（可增可减）
Gauge 通常⽤于测量温度或当前内存使⽤量等值，但也⽤于上下变化的“计数”，如并发请求的数量。

### **Histogram（直⽅图）**

histogram 对观测结果进⾏采样(通常是请求耗时或响应体⼤⼩) ，并按可配置的桶进⾏计数。它还提供了所有观察值的总和。

基本度量名称为 <basename> 的 histogram 会在抓取过程中暴露多个时间序列：
观察桶的累积计数器，对外展示为 <basename>_bucket{le="<upper inclusive bound>"}
所有观测值的总和，对外展示为 <basename>_sum
已观察到的事件数，对外展示为 <basename>_count （与上⾯的 <basename>_bucket{le="+Inf"} 相同）

使⽤ histogram_quantile() 函数可以根据 histogram 甚⾄ histogram 的聚合计算分位数。histogram 也适⽤于计算 Apdex 得分。在对bucket进⾏操作时，请记住 histogram 是累积的。

### **Summary（摘要）**

与 histogram 类似， summary 对观察结果（通常是请求耗时和响应体⼤⼩）进⾏采样。虽然它还提供了观测的总计数和所有观测值的总和，但它在滑动时间窗⼝内计算可配置的分位数。

基本度量名称为 <basename> 的 summary 会在抓取过程中暴露多个时间序列：

观测事件的流式**φ-quantiles**(0 ≤ φ ≤ 1) 分位数，对外展示为 <basename>{quantile="<φ>"}

所有观测值的总和，对外展示为 <basename>_sum

已观察到的事件数，对外展示为 <basename>_count

关于 histogram 和 summary 的区别，可以简单概括为 histogram 分桶记录数据，后续可在服务端使⽤表达式函数进⾏各种计算；⽽ summary 在客户端上报时就按配置上报计算好的φ-分位数。

1. 如果需要多个实例的数据进⾏汇总，请选择 histogram 。
2. 除此以外，如果对将要观察的值的范围和分布有所了解，请选择 histogram 。⽆论值的范围和分布如何，如果需要准确的分位数，请选择 summary 。

## **job 和 instance**

⽤ Prometheus 的术语来说，⼀个可以抓取的端点被称为⼀个 instance ，通常对应于⼀个进程。具有相同功能的实例集合（例如，为提⾼可扩展性/可靠性⽽创建的副本进程）称为 job。

当 Prometheus 抓取⽬标时，它会⾃动在抓取的时间序列上附加下⾯的标签，⽤于区分不同的⽬标：
job ：⽬标所属的已配置作业名
instance ：被抓取的⽬标URL的 <host>:<port> 部分。

对于每⼀次抓取，prometheus 都会按照以下时间序列存储⼀个样本：

up{job="<job-name>", instance="<instance-id>"} :如果实例是健康的，即可访问的，就是 1 或者如
果抓取失败，则为 0 。

scrape_duration_seconds{job="<job-name>", instance="<instance-id>"}
scrape_samples_post_metric_relabeling{job="<job-name>", instance="<instance-id>"}
scrape_samples_scraped{job="<job-name>", instance="<instance-id>"}
scrape_series_added{job="<job-name>", instance="<instance-id>"}

# 使用

## 安装

Prometheus ⽀持预编译⼆进制⽂件安装、源码安装、Docker等⽅式，由于我们是学习 Prometheus 的基本使⽤，所以在本地使⽤ Docker 快速开启⼀个实例。

```go
docker run -d --name=prometheus -p 9090:9090 prom/prometheus
```

上⾯的命令将使⽤⼀个示例配置启动 Prometheus Server，启动完成后，可以通过 [http://localhost:9090](http://localhost:9090/) 访问Prometheus的UI界⾯。

如果你有⾃定义的 prometheus.yml 配置。

```go
# my global config
global:
 scrape_interval: 15s # Set the scrape interval to every 15 seconds. Default is every
1 minute.
 evaluation_interval: 15s # Evaluate rules every 15 seconds. The default is every 1
minute.
# scrape_timeout is set to the global default (10s).
# Alertmanager configuration
alerting:
 alertmanagers:
 - static_configs:
 - targets:
# - alertmanager:9093
# Load rules once and periodically evaluate them according to the global
'evaluation_interval'.
rule_files:
# - "first_rules.yml"
# - "second_rules.yml"
# A scrape configuration containing exactly one endpoint to scrape:
# Here it's Prometheus itself.
scrape_configs:
# The job name is added as a label `job=<job_name>` to any timeseries scraped from
this config.
 - job_name: "prometheus"
# metrics_path defaults to '/metrics'
# scheme defaults to 'http'.
 static_configs:
 - targets: ["localhost:9090"]
```

[https://prometheus.io/docs/prometheus/latest/configuration/configuration/](https://prometheus.io/docs/prometheus/latest/configuration/configuration/)

```go
docker run \
-d \
--name=prometheus \
-p 9090:9090 \
-v /path/to/prometheus.yml:/etc/prometheus/prometheus.yml \
 prom/prometheus
```

Prometheus 数据存储在容器内的 /prometheus ⽬录中，因此每次重新启动容器时都会清除数据。要保存数据，需要为容器设置持久存储（或绑定挂载）。

运⾏具有持久存储的 prometheus 容器：

```go
# Create persistent volume for your data
docker volume create prometheus-data
# Start Prometheus container
docker run \
-d \
--name=prometheus \
-p 9090:9090 \
-v /path/to/prometheus.yml:/etc/prometheus/prometheus.yml \
-v prometheus-data:/prometheus \
 prom/prometheus
```

## **采集指标**

Prometheus 通过在⽬标节点的 HTTP 端⼝上采集 metric 数据来监控⽬标节点。因为 Prometheus 也以相同的⽅式暴露⾃⼰的指标数据，可以通过http://127.0.0.1:9090/metrics查看。

## **可视化**

### **grafana可视化**

Grafana 是⼀款开源的数据可视化⼯具，⽀持多种数据源（如Graphite、InfluxDB、OpenTSDB、Prometheus、Elasticsearch等）并且具有快速灵活的客户端图表，提供了丰富的仪表盘插件和⾯板插件，⽀持多种展示⽅式，如折线图、柱状图、饼图、点状图等，满⾜⽤户不同的可视化需求。

**安装 grafana**

```go
docker run -d --name=grafana --add-host=host.docker.internal:host-gateway -p 3000:3000 grafana/grafana-oss
// 启动成功后，使⽤浏览器打开http://localhost:3000。默认的登录账号是 "admin" / "admin"。
```

### **配置 prometheus 数据源**

1. 点击左侧菜单栏⾥的 『Connections』 图标。
2. 在数据源列表⾥找到 『prometheus 图标』或者搜索框输⼊ "prometheus" 搜索。
3. 点击 『prometheus 图标』，进⼊数据源⻚⾯。
4. 点击⻚⾯右上⻆蓝⾊ 『Add new data source』 按钮，添加数据源。
5. 填写 Prometheus server URL (例如, [http://localhost:9090/](http://localhost:9090/) )。
6. 根据需要调整其他数据源设置(例如, 认证或请求⽅法等)。
7. 点击⻚⾯下⽅的 『Save & Test』保存并测试数据源连接

### **添加仪表板**

1. 点击左侧菜单栏中的 『Dashboards』 。
2. 点击⻚⾯中间的 『+ Create Dashboard』 按钮。
3. 在打开的⻚⾯点击『+ Add visualization』按钮。
4. 在打开的⻚⾯上选择上⼀节添加的 prometheus data source。
5. 在打开的⻚⾯输⼊查询表达式 prometheus_target_interval_length_seconds ，点击『Run queries』执
⾏查询即可看到图表。
6. 点击右上⻆的『Save』保存仪表板。

## **Exporter 采集数据**

在 Prometheus 的架构设计中，Prometheus Server 并不直接负责监控特定的⽬标，其主要任务负责数据的收集、存储以及对外提供数据查询⽀持。因此为了能够能够监控到某些指标，如主机的CPU使⽤率，我们需要使⽤到Exporter。Prometheus 周期性的从 Exporter暴露的HTTP服务地址（通常是/metrics）拉取监控样本数据。

⼴义上讲所有可以向 Prometheus 提供监控样本数据的程序都可以被称为⼀个 Exporter。⽽⼀个 Exporter 实例被称为 target，如下所示，Prometheus 通过轮询的⽅式定期从这些 target 中获取样本数据

**Exporter 有两种运⾏⽅式**

独⽴运⾏

被监控对象⽆法直接提供监控接⼝，可能的原因有：

1. 不能直接提供 HTTP 接⼝，如监控 Linux 系统状态指标。

2. 项⽬发布时间较早，不⽀持 Prometheus 监控接⼝，如 MySQL、Redis；

对于以上场景，可以选择使⽤独⽴运⾏的 Exporter。

集成到应⽤中

通过在应⽤程序内部使⽤ Prometheus 提供的 Client Library，将程序内部的运⾏状态主动暴露给
Prometheus，适⽤于需要较多⾃定义监控指标的项⽬。⽬前⼀些开源项⽬就增加了对 Prometheus 监控的原⽣⽀持，如 Kubernetes，ETCD 等。
我们也可以在业务代码中增加⾃定义指标数据上报⾄ Prometheus 。

### **社区提供的 exporter**

**范围**

**常⽤Exporter**

数据库

MySQL server exporter (**official**)、MSSQL server exporter、Elasticsearch

exporter、MongoDB exporter、Redis exporter 等

硬件

apcupsd exporter，Node/system metrics exporter (**official**)，NVIDIA GPU exporter,

Windows exporter等

问题跟踪和持续集成

Jenkins exporter，JIRA exporter

消息队列

Kafka exporter, RabbitMQ exporter, RocketMQ exporter, NSQ exporter等

存储

Ceph exporter, Hadoop HDFS FSImage exporter等

HTTP服务

Apache exporter, HAProxy exporter (**official**), Nginx metric library等

API服务

AWS ECS exporter，Azure Health exporter, Cloudflare exporter等

⽇志

Fluentd exporter ,Grok exporter等

监控系统

Alibaba Cloudmonitor exporter, AWS CloudWatch exporter (**official**), Azure Monitor

exporter, JMX exporter (**official**), TencentCloud monitor exporter等

其它

eBPF exporter，Kibana Exporter，SSH exporter,等

此外还有⼀些第三⽅软件默认就提供 Prometheus 格式的指标数据，因此不需要单独的 Exporter 。

Envoy

Etcd (**direct**)

Flink

Grafana

Kong

Kubernetes (**direct**)

RabbitMQ

⽐如 MySQL 数据库的exporter：推荐看⼀下!
[https://yunlzheng.gitbook.io/prometheus-book/part-ii-prome](https://yunlzheng.gitbook.io/prometheus-book/part-ii-prome)
theus-jin-jie/exporter/commonly-eporter-usage/use-promethues-monitor-mysql

## **使⽤ prometheus client 库实现**

### **Prometheus Go Client**

```go
go get github.com/gin-gonic/gin
go get github.com/prometheus/client_golang/prometheus
go get github.com/prometheus/client_golang/prometheus/promhttp
```

```go
// ⾃定义业务状态码 Counter 指标
var statusCounter = prometheus.NewCounterVec(
 prometheus.CounterOpts{
 Name: "api_response_status_count",
 },
 []string{"method", "path", "status"},
)
```

```go
func initRegistry() *prometheus.Registry {
 // 创建⼀个 registry
 reg := prometheus.NewRegistry()
 // 添加 Go 编译信息
 reg.MustRegister(collectors.NewBuildInfoCollector())
 // Go runtime metrics
 reg.MustRegister(collectors.NewGoCollector(
 collectors.WithGoCollectorRuntimeMetrics(
 collectors.GoRuntimeMetricsRule{Matcher: regexp.MustCompile("/.*")},
 ),
 ))
 // 注册⾃定义的业务指标
 reg.MustRegister(statusCounter)
 return reg
}
```

```go
// 记录
 statusCounter.WithLabelValues(
 c.Request.Method,
 c.Request.URL.Path,
 strconv.Itoa(status),
 ).Inc()
 c.JSON(200, gin.H{
 "status": status,
 "message": "pong",
 })
 })
 reg := initRegistry()
```

```go
// 对外提供 /metrics 接⼝，⽀持 prometheus 采集
 r.GET("/metrics", gin.WrapH(promhttp.HandlerFor(
 reg,
 promhttp.HandlerOpts{Registry: reg},
 )))
```

其上有需要修改的地方，一是起容器时，这里我的mac是这样，需要指定

```go
--add-host=host.docker.internal:host-gateway
```

其次，示例中是8083端口，需要修改prometheus的配置文件，`prometheus.yml`

```go
scrape_configs:
  - job_name: 'myapp'
    static_configs:
      - targets: ['host.docker.internal:8083']
```

Go-zero prometheus: [https://go-zero.dev/docs/tutorials/monitor/index#指标��%](https://go-zero.dev/docs/tutorials/monitor/index#%E6%8C%87%E6%A0%87%E7%9B%25)
91%E6%8E%A7
Kratos prometheus：[https://go-kratos.dev/docs/component/middleware/metrics/](https://go-kratos.dev/docs/component/middleware/metrics/)

## **官⽅⽂档原⽂**

[https://prometheus.io/docs/prometheus/latest/querying/basics/](https://prometheus.io/docs/prometheus/latest/querying/basics/)

[https://prometheus.io/docs/prometheus/latest/querying/examples/](https://prometheus.io/docs/prometheus/latest/querying/examples/)

# **编写⼀个完善的技术⽅案**

⼀个完善的技术⽅案⾄少要包含：

需求背景
要实现的功能
技术选型（对⽐）
流程图
接⼝时序图
安全相关

资源安全：防刷、防资损、机制漏洞
数据安全：⽔平越权等
系统安全：防注⼊等，传⻩图或者⾮法图⽚怎么办？

上线⽅案

依赖项
上线顺序

现在你负责的业务要加⼀个拍照搜索功能，技术⽅案应该怎么写？

拍照搜索怎么实现？

可以采买某云⼚商的商⽤服务
公司⾃⼰训练的模型。

1、如果使⽤商⽤服务（要花钱），怎么防⽌接⼝被刷？

按⽤户做限流、
如果是不登陆也能⽤的功能怎么限？

⽤滑动窗⼝限制（⼤家都⽤⼀个池⼦）
⽤机器标识限制

2、怎么样去评估资源⽤量？QPS=1 100块，QPS=100 要10000块。

经过实际的流量评估，⼤概QPS=1就够了，就购买QPS=1级别的资源。
要留个升级的机制！

上线之后效果好⽤的⼈多，QPS=1不够⽤了，要申请预算增加QPS。
要在代码⾥记录调⽤商业服务的返回值，⽐如正常返回code=0, 返回code=10024 表示query busy，QPS不够⽤了。
在代码⾥⾃定义指标，记录返回的 code ，计算 code=10024 的占⽐或增⻓率。

3、针对 公司⾃⼰训练的模型 的⽅式 也应该多想⼀步。

⾃⼰公司训练模型的话有没有可能需要去做调优？要记录这份数据

模型的识别率在业务侧⾃⼰统计⼀下
可以记录识别的详细⽇志：输⼊是什么，输出是什么
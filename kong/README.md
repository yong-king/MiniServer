# 网关

## **API⽹关介绍**

**API⽹关** 简单来说是⼀种主要⼯作在七层、专⻔⽤于 API 的管理和流量转发的基础设施，并拥有强⼤的扩展性。⽹关的⻆⾊是作为⼀个API架构，⽤来保护、增强和控制对于API服务的访问。它是⼀个处于应⽤程序或服务（提供RESTful API接⼝服务）之前的系统，⽤来管理授权、访问控制和流量限制等。这样RESTful API接⼝服务就被⽹关保护起来，对所有的调⽤者透明。因此，隐藏在API⽹关后⾯的业务系统就可以专注于创建和管理服务，⽆需关⼼这些策略性的请求。

API Gateway是⼀种服务，他是外部世界进⼊应⽤程序的⼊⼝点。它负责请求路由、API组合和身份验证等各项功能。

API Gateway 还可以作为客户端应⽤程序和后端微服务架构之间的反向代理。

## **常⽤API⽹关**

**Kong**

Kong 基于OpenResty（Nginx + Lua模块）编写的⾼可⽤、易扩展的，性能⾼效且稳定，⽀持多个可⽤插件（限流、鉴权）等，开箱即⽤。

相对于纯Nginx ，Kong具有以下优点：

（1）⾼性能：亚毫秒级处理延迟，可⽀持关键任务⽤例和⾼吞吐量。
（2）可扩展性：可插拔的体系结构，可通过Kong的Plugin SDK扩展 Kong。
（3）可移植性：Kong 可以部署在任何平台、或者云。

**APISIX**

它构建于 NGINX + ngx_lua 的技术基础之上，充分利⽤了 LuaJIT 所提供的强⼤性能。

APISIX 主要分为两个部分：

1. APISIX 核⼼：包括 Lua 插件、多语⾔插件运⾏时（Plugin Runner）、Wasm 插件运⾏时等；
2. 功能丰富的各种内置插件：包括可观测性、安全、流量控制等。

**APISIX的优势**

在单体服务时代，使⽤ Nginx 可以应对⼤多数的场景，⽽到了云原⽣时代，Nginx 因为其⾃身架构的原因则会出现
两个问题：
⾸先是 Nginx 不⽀持集群管理。⼏乎每家互联⽹⼚商都有⾃⼰的 Nginx 配置管理系统，系统虽然⼤同⼩异但是⼀直没有统⼀的⽅案。

其次是 Nginx 不⽀持配置的热加载。很多公司⼀旦修改了配置，重新加载 Nginx 的时间可能需要半个⼩时以上。并且在 Kubernetes 体系下，上游会经常发⽣变化，如果使⽤ Nginx 来处理就需要频繁重启服务，这对于企业是不可接受的。

⽽ Kong 的出现则解决了 Nginx 的痛点，但是⼜带来了新的问题：

Kong 需要依赖于 PostgreSQL 或 Cassandra 数据库，这使 Kong 的整个架构⾮常臃肿，并且会给企业带来⾼可⽤的问题。如果数据库故障了，那么整个 API ⽹关都会出现故障。
Kong 的路由使⽤的是遍历查找，当⽹关内有超过上千个路由时，它的性能就会出现⽐较急剧的降。

**Traefik**

Traefik 是⼀个使⽤Go语⾔开发的云原⽣的新型的 HTTP 反向代理、负载均衡软件。它负责接收系统的请求，然后使⽤合适的组件来对这些请求进⾏处理。它⽀持多种后台 (Docker, Swarm,Kubernetes, Marathon, Mesos,Consul, Etcd, Zookeeper, BoltDB, Rest API, file…) 来⾃动化、动态的应⽤它的配置⽂件设置。

# **Kong**

## **Kong介绍**

Kong 基于OpenResty（Nginx + Lua模块）编写的⾼可⽤、易扩展的，性能⾼效且稳定，⽀持多个可⽤插件（限流、鉴权）等，开箱即⽤。

相对于纯Nginx ，Kong具有以下优点：
（1）⾼性能：亚毫秒级处理延迟，可⽀持关键任务⽤例和⾼吞吐量。
（2）可扩展性：可插拔的体系结构，可通过Kong的Plugin SDK扩展 Kong。
（3）可移植性：Kong 可以部署在任何平台、或者云。

### **相关概念**

Kong⽹关管理员使⽤对象模型来定义他们想要的流量管理策略。该模型中的两个重要对象是服务和路由。服务和路由以协调的⽅式进⾏配置，以定义请求和响应在系统中的路由路径。

**service**

在 Kong Gateway 中，服务是现有上游应⽤程序的抽象。对应着后端的⼀个服务，或者是⼀个App。

服务可以存储诸如插件配置和策略之类的对象集合，并且它们可以与路由关联。

**route**

路由是到上游应⽤程序中的资源的路径。不同的route对应着service中的不同接⼝。

在 Kong Gateway 中，路由通常映射到通过 Kong Gateway 应⽤程序公开的接⼝。路由还可以定义将请求匹配到关联服务的规则。因此，⼀个路由可以引⽤多个端点。路由应该有⼀个名称、路径或路径，并引⽤现有的服务。

路由⽀持以下配置：
Protocols: ⽤于与上游应⽤程序通信的协议。
Hosts：匹配路由的域名列表。
Methods：匹配路由的 HTTP ⽅法
Headers：请求标头中预期的值列表
Redirect status codes：HTTPS 状态码
Tags：⽤于将路由分组的可选字符串集

**Upstream和Target**

Upstream：和nginx中的upstream差不多，都对应着⼀组服务节点。

Target：对应着⼀个api服务节点。

## **Kong管理后台**

**Konga**

[https://github.com/pantsel/konga](https://github.com/pantsel/konga) 之前社区版本的管理后台（⾃从官⽅推出Kong Manager之后就没更新了）

**Kong Manager**

Kong Manager是Kong官⽅推出的管理后台，也有开源版本 Kong Manager Open Source（也称OSS版本）。

### **Kong安装**

Kong⽀持多种平台安装。[https://konghq.com/install#kong-community](https://konghq.com/install#kong-community)

### **本地环境部署**

本地快速搭建kong环境，推荐使⽤官⽅提供的 docker-compose.yml

1. git clone 官⽅提供的代码库⾄本机，然后切换到该⽬录下。
    
    ```go
    $ git clone https://github.com/Kong/docker-kong
    $ cd docker-kong/compose/
    ```
    
2. 执⾏以下命令启动Kong⽹关相关实例。
    
    ```go
    $ KONG_DATABASE=postgres docker-compose --profile database up -d
    ```
    

Kong ⽹关现在可以在本地主机的下列端⼝上使⽤:
:8000 - Kong ⽹关的流量⼊⼝（外部流量通过8000端⼝流⼊，经过Kong转发⾄实际服务节点）
:8001 - Kong ⽹关的配置端⼝（可以通过Admin API 或decK对Kong进⾏配置）
:8002 - 可以通过localhost:8002访问 Kong 的管理 Web ⽤户界⾯（Kong Manager）

1. **添加Service**
    1. Name service名称，全局唯一
    2. 标签
    3. Full URL  upstream url 上游URL
2. 添加路由
    1. Name
    2. service 下拉选择
    3. paths 路由规则

http://localhost:8000/xxx

1. 负载均衡
    1. **设置upstreams 和 targets**
        1. service中修改Host 为 xxx_upstream
        2. upstream 添加
            1. Name 下拉框选择
            2. target：ip：port
2. 发布策略
    1. routes中添加header，使用modHeader扩展测试
    2. 调整权重
3. 插件
    1. **Rate Limiting（速率限制）**
    2. **Proxy Caching （缓存）**
    3. **Authentication with JWT （认证）**
        1. C端⽤户 --> 登录 --> 下发JWT
        2. C端⽤户携带JWT --> Kong⽹关（认证校验）--> 业务接⼝
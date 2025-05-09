# 评价项目

# 模块及实现

### 项目框架搭建

1. 创建项目
    
    ```go
    kratos new review-service
    ```
    
2. 添加proto
    
    ```go
    kratos proto add api/review/v1/review.proto
    ```
    
3. 生成客户端代码
    
    ```go
    kratos proto client api/review/v1/review.proto
    ```
    
4. 生成服务端代码
    
    ```go
    kratos proto server api/review/v1/review.proto -t internal/service
    ```
    

### 项目依赖准备

1. Mysql环境（docker）
2. Redis环境（docker）
3. 建立数据库表
4. 修改config配置文件 mysql，redis 
    
    ```go
    make config
    ```
    

## 开发接口流程

1. 定义api文件
    
    根据需求编写api文件
    
2. 生成客户端和服务端代码
    
    ```go
    make api 
    ```
    
3. 填充业务逻辑
    
    ./internal
    
    server → service → biz → data
    
4. 更新依赖注入

## 业务开发

### 评论服务

1. 创建评论
    1. 雪花算法生成ID
    2. validate参数校验
        1. 下载插件
        2. 在api中的pb文件编写 校验规则
        3. 生成代码
        4. 注册参数中间件
        
        ```go
        // 下载插件
        go install github.com/envoyproxy/protoc-gen-validate@latest
        
        // 导入
        import "validate/validate.proto";
        
        // 生成代码
        make validate
        .PHONY: validate
        # generate validate proto
        validate:
        	protoc --proto_path=. \
                   --proto_path=./third_party \
                   --go_out=paths=source_relative:. \
                   --validate_out=paths=source_relative,lang=go:. \
                   $(API_PROTO_FILES)
        
                   
        //中间件 server
        		http.Middleware(
        			recovery.Recovery(),
        			validate.Validator(),
        		),
        	}
        			grpc.Middleware(
        			recovery.Recovery(),
        			validate.Validator(),
        		),
        	}
        ```
        
2. 错误处理
    1. 定义proto文件
    2. 生成代码
    3. 业务中使用生成的错误代码返回
    
    ```go
    // 安装
    go install github.com/go-kratos/kratos/cmd/protoc-gen-go-errors/v2@latest
    
    protoc --proto_path=. \
             --proto_path=./third_party \
             --go_out=paths=source_relative:. \
             --go-errors_out=paths=source_relative:. \
             $(API_PROTO_FILES)
             
    // 或者
    make errors
    ```
    
3. 评价详情
4. 审核评价
5. 申诉评价
    
    ```go
    r.data.query.ReviewAppealInfo.
    		WithContext(ctx).
    		Clauses(clause.OnConflict{
    			Columns: []clause.Column{
    				{Name: "review_id"}, // ON DUPLICATE KEY
    			},
    			DoUpdates: clause.Assignments(map[string]interface{}{ // UPDATE
    				"status":     10,
    				"content":    appeal.Content,
    				"reason":     appeal.Reason,
    				"pic_info":   appeal.PicInfo,
    				"video_info": appeal.VideoInfo,
    			}),
    		}).
    		Create(appeal) // INSERT
    		当存在时就更新，当不存在是就创建
    		
    		INSERT INTO `table` *** ON DUPLICATE KEY UPDATE ***;
    ```
    
6. 回复评价
    1. review-service 新增grpc调用
        1. gorm-gen事务操作
            
            ```go
            r.data.query.Transaction(func(tx *query.Query) error {
            		// 回复一条插入数据
            		if err := tx.ReviewReplyInfo.WithContext(ctx).Save(reply); err != nil {
            			r.log.WithContext(ctx).Errorf("SaveReply create reply fail, err:%v", err)
            			return err
            		}
            		// 评价表更新hasReply字段\
            		if _, err := tx.ReviewInfo.WithContext(ctx).Where(tx.ReviewInfo.ReviewID.Eq(review.ReviewID)).Update(tx.ReviewInfo.HasReply, 1); err != nil {
            			r.log.WithContext(ctx).Errorf("SaveReply update reply fail, err:%v", err)
            			return err
            		}
            		return nil
            	})
            ```
            
        2. 防止水平越权
            
            ```go
            if review.StoreID != reply.StoreID {
            		return nil, errors.New("水平越权")
            	}
            ```
            
7. 评价列表

### 评价服务c端

1. 发表评价
2. 查看评价
3. 查看自己的评价

### 评价b端

1. 店铺评价列表
2. 店铺评价详情
3. 回复评价
4. 申诉评价

### 评价o端

1. 评价列表（筛选）
2. 评价详情
3. 评价审核
4. 评价申诉

**关键点**

雪花算法⽣成ID

validate参数校验

GORM事务操作

接⼝幂等

接⼝防⽌⽔平越权

## go submodule

项⽬中如何管理pb⽂件
proto⽂件要⽤⼀个
protoc要使⽤同⼀个版本
通常在公司中都是把 proto ⽂件和⽣成的不同语⾔的代码都放在⼀个单独的公⽤代码库中。
别的项⽬直接引⽤这个公⽤代码库。

语法：git⼦模块
项⽬中添加⼦模块。将 [git@github.com](mailto:git@github.com):Q1mi/reviewapis.git 作为当前项⽬的⼦⽬录，⽬录名为 api 

```go
git submodule add git@xxxxxxx /api
```

当前⽬录下会多⼀个 .gitmodules ⽂件

```go
# ⽤来初始化本地配置⽂件
git submodule init
# 从该项⽬中抓取所有数据并检出⽗项⽬中列出的合适的提交。
git submodule update
```

## 服务注册与服务发现

1. consul注册中心
    1. internal/conf/xx.proto
    2. configs/xx.yaml

```go
message Registry{
  message Consul{
    string address = 1;
    string scheme = 2;
  }
  Consul consul = 1;
}

consul:
  address: 127.0.0.1:8500
  scheme: http
  
 make config
```

1. service添加服务注册
    1. 注册的时机 --> internal/server层 --> 提供构造函数--> wire注⼊
    2. main函数传⼊conf.Registry配置
    3. 指定应⽤程序的name和version，在注册时使⽤

```go
// server
func NewRegistrar(cfg *conf.Registry) registry.Registrar {
	// new consul client
	c := api.DefaultConfig()
	c.Address = cfg.Consul.Address
	c.Scheme = cfg.Consul.Scheme
	client, err := api.NewClient(c)
	if err != nil {
		panic(err)
	}
	// new reg with consul client
	reg := consul.New(client, consul.WithHealthCheck(true))
	return reg
}

var ProviderSet = wire.NewSet(NewRegistrar, NewGRPCServer, NewHTTPServer)

// mian.go
func newApp(logger log.Logger, r registry.Registrar, gs *grpc.Server, hs *http.Server) *kratos.App {
	return kratos.New(
		kratos.ID(id),
		kratos.Name(Name),
		kratos.Version(Version),
		kratos.Metadata(map[string]string{}),
		kratos.Logger(logger),
		kratos.Server(
			gs,
			hs,
		),
		kratos.Registrar(r),
	)
}

var (
	// Name is the name of the compiled software.
	Name string =  "review.service"
	// Version is the version of the compiled software.
	Version string = "v0.1"
	// flagconf is the config flag.
	flagconf string

	id, _ = os.Hostname()
)

// 解析registry.yaml中的配置
	var rc conf.Registry
	if err := c.Scan(&rc); err != nil {
		panic(err)
	}

	app, cleanup, err := wireApp(bc.Server, &rc, bc.Data, logger)
	
// wire.go
func wireApp(*conf.Server, *conf.Registry, *conf.Data, log.Logger) (*kratos.App, func(), error) {
	panic(wire.Build(server.ProviderSet, data.ProviderSet, biz.ProviderSet, service.ProviderSet, newApp))
}

cd cmd/review -> wire
```

1. review-b添加服务发现流程\
    1. 服务发现的时机 --> internal/data 层 --> 提供构造函数 --> wire注⼊
    2. main函数传⼊conf.Registry配置

```go
registry:
  consul:
    address: 127.0.0.1:8500
    scheme: http
  
  message Bootstrap {
  Server server = 1;
  Data data = 2;
  Registry registry = 3;
}

message Registry {
  message Consul {
    string address = 1;
    string scheme = 2;
  }
  Consul consul = 1;
}

// data.go 
func NewDiscovever(conf *conf.Registry) registry.Discovery{
	// new consul client
	c := api.DefaultConfig()
	c.Address = conf.Consul.Address
	c.Scheme = conf.Consul.Scheme
	client, err := api.NewClient(c)
	if err != nil {
		panic(err)
	}
	// new dis with consul client
	dis := consul.New(client)
	return dis
}

func NewReviewServiceClient(d registry.Discovery) v1.ReviewClient {
	endpoint := "discovery:///review.service"
	conn, err := grpc.DialInsecure(context.Background(),
		// grpc.WithEndpoint("127.0.0.1:9001"),
		grpc.WithEndpoint(endpoint),
		grpc.WithDiscovery(d),
		grpc.WithMiddleware(
			recovery.Recovery(),
			validate.Validator(),
		))
	if err != nil {
		panic(err)
	}
	return v1.NewReviewClient(conn)
}

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewData, NewBusinessRepo, NewReviewServiceClient, NewDiscovever)

// main.go
app, cleanup, err := wireApp(bc.Server, bc.Registry, bc.Data, logger)
	if err != nil {
		panic(err)
	}
	defer cleanup()
	
// wire.go 
func wireApp(*conf.Server, *conf.Registry, *conf.Data, log.Logger) (*kratos.App, func(), error) {
	panic(wire.Build(server.ProviderSet, data.ProviderSet, biz.ProviderSet, service.ProviderSet, newApp))
}

cd cmd/service-b -> wire
```

## canal

Canal 是阿里开源的一款 MySQL 数据库增量日志解析工具，提供增量数据订阅和消费。使用Canal能够实现异步更新数据，配合MQ使用可在很多业务场景下发挥巨大作用。

1. **工作原理**
    
    **MySQL主备复制原理**
    
    - MySQL master 将数据变更写入二进制日志（binary log）, 日志中的记录叫做二进制日志事件（binary log events，可以通过 show binlog events 进行查看）
    - MySQL slave 将 master 的 binary log events 拷贝到它的中继日志(relay log)
    - MySQL slave 重放 relay log 中事件，将数据变更反映到它自己的数据
    
    **Canal 工作原理**
    
    - Canal 模拟 MySQL slave 的交互协议，伪装自己为 MySQL slave ，向 MySQL master 发送 dump 协议
    - MySQL master 收到 dump 请求，开始推送 binary log 给 slave (即 Canal )
    - Canal 解析 binary log 对象(原始为 byte 流)
2. **环境准备**
    1. MySQL环境
        1. **开启binlog**
        2. my.cnf,修改配置文件之后，重启MySQL。
            
            ```go
            /etc/mysql/my.cnf
            [mysqld]
            log-bin=mysql-bin # 开启 binlog
            binlog-format=ROW # 选择 ROW 模式
            server_id=1 # 配置 MySQL replaction 需要定义，不要和 canal 的 slaveId 重复
            
            //查看
            show variables like 'log_bin'; -> on
            show variables like 'binlog_format'; -> row
            
            ```
            
        3. **添加授权**
            
            ```go
            CREATE USER canan(name) IDENTIFIED BY 'password';  // 按需填写name和password
            GRANT SELECT, REPLICATION SLAVE, REPLICATION CLIENT ON *.* TO 'canal'@'%';
            -- GRANT ALL PRIVILEGES ON *.* TO 'canal'@'%' ;
            FLUSH PRIVILEGES;
            
            ```
            
    2. **安装Canal**
        
        docker:
        
        ```go
        docker pull canal/canal-server:latest
        
        // 启动容器
        docker run -d \
          --name canal-server \
          --add-host=host.docker.internal:host-gateway \
          canal/canal-server:latest
        
        // 进入容器
        docker exec -it canal-server /bin/bash
        
        // 修该配置
        vi canal-server/conf/example/instance.properties
        
        canal.instance.master.address=host.docker.internal:3306
        
        canal.instance.tsdb.dbUsername=canal
        canal.instance.tsdb.dbPassword=password // 上面的name和password
        ```
        

## kafka

Kafka是一种高吞吐量的分布式发布订阅消息系统

Apache Kafka是⼀个开源的分布式流系统，该项⽬旨在提供⼀个统⼀的、⾼吞吐量、低延迟的平台，⽤于处理实时数据流。它具有以下特点：
⽀持消息的发布和订阅，类似于 RabbtMQ、RocketMQ 等消息队列；
⽀持数据实时处理；
能保证消息的可靠性投递；
⽀持消息的持久化存储，并通过多副本分布式的存储⽅案来保证消息的容错；

⾼吞吐率，单 Broker 可以轻松处理数千个分区以及每秒百万级的消息量。

Kafka是⼀个数据流系统，允许开发⼈员在新事件发⽣时实时做出反应。Kafka体系结构由存储层和计算层组成。存
储层旨在⾼效地存储数据，是⼀个分布式系统，可以轻松地扩展系统以适应增⻓。
计算层由四个核⼼组件组成——⽣产者、消费者、流和连接器API，它们允许Kafka在分布式系统中扩展应⽤程序。

1. ⽣产者（Producer）
2. 消费者（Consumer）
3. 流处理（Streams）
4. 连接器（Connectors）APIs

**相关术语**

Messages And Batches：Kafka 的基本数据单元被称为 message(消息)，为减少⽹络开销，提⾼效率，多个

消息会被放⼊同⼀批次 (Batch) 中后再写⼊。

Topic：⽤来对消息进⾏分类，每个进⼊到Kafka的信息都会被放到⼀个Topic下

Broker：⽤来实现数据存储的主机服务器,kafka节点

Partition：每个Topic中的消息会被分为若⼲个Partition，以提⾼消息的处理效率

Producer：消息的⽣产者Consumer：消息的消费者

Consumer Group：消息的消费群组

Kafka 的消息通过 Topics(主题) 进⾏分类，⼀个主题可以被分为若⼲个 Partitions(分区)，⼀个分区就是⼀个提交⽇志 (commit log)。消息以追加的⽅式写⼊分区，然后以先⼊先出的顺序读取。Kafka 通过分区来实现数据的冗余和伸缩性，分区可以分布在不同的服务器上，这意味着⼀个 Topic 可以横跨多个服务器，以提供⽐单个服务器更强⼤的性能。由于⼀个 Topic 包含多个分区，因此⽆法在整个 Topic 范围内保证消息的顺序性，但可以保证消息在单个分区内的顺序性。

为了分散主题中事件的存储和处理，Kafka使⽤了分区的概念。⼀个主题由⼀个或多个分区组成，这些分区可以位于Kafka集群中的不同节点上。每个分区都是⼀个有序的，不可变的记录序列，不断附加到结构化的提交⽇志中。分区中的记录每都分配了⼀个称为偏移的顺序ID号，它唯⼀地标识分区中的每个记录。Kafka集群⽀持按配置持久化保存所有已发布的记录。例如，如果保留策略设置为两天，则在发布记录后的两天内，它可供消费，之后将被丢弃以释放空间。

**⽣产者**

⽣产者负责创建消息。⼀般情况下，⽣产者在把消息均衡地分布到在主题的所有分区上，⽽并不关⼼消息会被写到哪个分区。如果我们想要把消息写到指定的分区，可以通过⾃定义分区器来实现。

消费者是消费者群组的⼀部分，消费者负责消费消息。消费者可以订阅⼀个或者多个主题，并按照消息⽣成的顺序来读取它们。消费者通过检查消息的偏移量 (offset) 来区分读取过的消息。偏移量是⼀个不断递增的数值，在创建消息时，Kafka 会把它添加到其中，在给定的分区⾥，每个消息的偏移量都是唯⼀的。消费者把每个分区最后读取的偏移量保存在 Zookeeper 或 Kafka 上，如果消费者关闭或者重启，它还可以重新获取该偏移量，以保证读取状态不会丢失。

⼀个分区只能被同⼀个消费者群组⾥⾯的⼀个消费者读取，但可以被不同消费者群组中所组成的多个消费者共同读取。多个消费者群组中消费者共同读取同⼀个主题时，彼此之间互不影响。

⼀个独⽴的 Kafka 服务器被称为 Broker。Broker 接收来⾃⽣产者的消息，为消息设置偏移量，并提交消息到磁盘保存。Broker 为消费者提供服务，对读取分区的请求做出响应，返回已经提交到磁盘的消息。Broker 是集群 (Cluster) 的组成部分。每⼀个集群都会选举出⼀个 Broker 作为集群控制器 (Controller)，集群控制器负责管理⼯作，包括将分区分配给 Broker 和监控 Broker。
在集群中，⼀个分区 (Partition) 从属⼀个 Broker，该 Broker 被称为分区的⾸领 (Leader)。⼀个分区可以分配给多个 Brokers，这个时候会发⽣分区复制。这种复制机制为分区提供了消息冗余，如果有⼀个 Broker 失效，其他
Broker 可以接管领导权。

1. 环境准备
    
    ```go
    // docker-compose.yml
    version: '2.1'
    
    services:
      zoo1:
        image: confluentinc/cp-zookeeper:7.3.2
        hostname: zoo1
        container_name: zoo1
        ports:
          - "2181:2181"
        environment:
          ZOOKEEPER_CLIENT_PORT: 2181
          ZOOKEEPER_SERVER_ID: 1
          ZOOKEEPER_SERVERS: zoo1:2888:3888
    
      kafka1:
        image: confluentinc/cp-kafka:7.3.2
        hostname: kafka1
        container_name: kafka1
        ports:
          - "9092:9092"
          - "29092:29092"
          - "9999:9999"
        environment:
          KAFKA_ADVERTISED_LISTENERS: INTERNAL://kafka1:19092,EXTERNAL://${DOCKER_HOST_IP:-127.0.0.1}:9092,DOCKER://host.docker.internal:29092
          KAFKA_LISTENER_SECURITY_PROTOCOL_MAP: INTERNAL:PLAINTEXT,EXTERNAL:PLAINTEXT,DOCKER:PLAINTEXT
          KAFKA_INTER_BROKER_LISTENER_NAME: INTERNAL
          KAFKA_ZOOKEEPER_CONNECT: "zoo1:2181"
          KAFKA_BROKER_ID: 1
          KAFKA_LOG4J_LOGGERS: "kafka.controller=INFO,kafka.producer.async.DefaultEventHandler=INFO,state.change.logger=INFO"
          KAFKA_OFFSETS_TOPIC_REPLICATION_FACTOR: 1
          KAFKA_TRANSACTION_STATE_LOG_REPLICATION_FACTOR: 1
          KAFKA_TRANSACTION_STATE_LOG_MIN_ISR: 1
          KAFKA_JMX_PORT: 9999
          KAFKA_JMX_HOSTNAME: ${DOCKER_HOST_IP:-127.0.0.1}
          KAFKA_AUTHORIZER_CLASS_NAME: kafka.security.authorizer.AclAuthorizer
          KAFKA_ALLOW_EVERYONE_IF_NO_ACL_FOUND: "true"
        depends_on:
          - zoo1
      kafka-ui:
        container_name: kafka-ui
        image: provectuslabs/kafka-ui:latest
        extra_hosts:  # 👈 添加此配置
          - "host.docker.internal:host-gateway" // 避免访问不到,需要添加
        ports:
          - 8080:8080
        depends_on:
          - kafka1
        environment:
          DYNAMIC_CONFIG_ENABLED: "TRUE"
    
    docker-compose up -d
    
    // 127.0.0.1:8080 
    ```
    
2. 添加配置
    1. cluster name：都可以
    2. bootstrap servers ： host.docker.internal 29092
    3. validate→ok0>submit

### cancal&kafka&go

1. **安装kafka-go**
    
    ```go
    go get github.com/segmentio/kafka-go
    ```
    
2. **Canal Kafka**
    
    ```go
    vi canal-server/conf/example/instance.properties
    
    canal.mq.dynamicTopic=mytest,.*,mytest.user,mytest\\..*,.*\\..* // 按需要修改
    例如：
    topic3:commit\\..* //commit是数据库
    ```
    
    修改canal配置文件
    
    ```go
    vi /home/admin/canal-server/conf/canal.properties
    
    # 可选项: tcp(默认), kafka,RocketMQ,rabbitmq,pulsarmq
    canal.serverMode = kafka
    
    ##################################################
    #########                    Kafka                   #############
    ##################################################
    # 此处配置修改为你的Kafka环境地址
    kafka.bootstrap.servers = 127.0.0.1:9092
    ```
    

## Elasticsearch

Elasticsearch 是一个高度可扩展的开源实时搜索和分析引擎，它允许用户在近实时的时间内执行全文搜索、结构化搜索、聚合、过滤等功能。Elasticsearch 基于 Lucene 构建，提供了强大的全文搜索功能，并且具有广泛的应用领域，包括日志和实时分析、社交媒体、电子商务等。

Elasticsearch 为所有类型的数据提供近乎实时的搜索和分析。无论是结构化文本还是非结构化文本、数字数据或地理空间数据，Elasticsearch 都能够以支持快速搜索的方式有效地存储和索引它们。除了简单的数据检索和聚合信息之外，还可以用 Elasticsearch 发现数据中的趋势和模式。随着数据和查询量的增长，Elasticsearch 的分布式特性能够横向扩展至数以百计的服务器存储以及处理PB级的数据，同时可以在极短的时间内索引、搜索和分析大量的数据。

- 为APP或网站增加搜索功能
- 存储和分析日志、指标和安全事件数据
- 使用机器学习实时自动建模数据的行为
- 使用Elasticsearch作为存储引擎自动化业务工作流
- 使用Elasticsearch作为地理信息系统（GIS）管理、集成和分析空间信息
- 使用Elasticsearch作为生物信息学研究工具存储和处理遗传数据

Elasticsearch 架构主要由三个组件构成：索引、分片和节点。

- 索引是文档的逻辑分组，类似于数据库中的表；
- 分片是索引的物理分区，用于提高数据分布和查询性能；
- 节点是运行 Elasticsearch 的服务器实例。

Elasticsearch 通过以下步骤完成搜索和分析任务：

1. 接收用户查询请求：Elasticsearch 通过 RESTful API 或 JSON 请求接收用户的查询请求。
2. 路由请求：接收到查询请求后，Elasticsearch 根据请求中的索引和分片信息将请求路由到相应的节点。
3. 执行查询：节点执行查询请求，并在相应的索引中查找匹配的文档。
4. 返回结果：查询结果以 JSON 格式返回给用户，包括匹配的文档和相关字段信息。

索引（Index）

在Elasticsearch中，索引是存储相关数据的数据结构，可以理解为数据库中的表。索引是通过对数据源进行索引创建的，它是一种对数据进行结构化和半结构化处理的结果。每个索引都有自己的映射（mapping），用于定义每个字段的数据类型和其他属性。

在Elasticsearch中，索引的创建和定义通常是通过REST API或者相关Java API来实现的。在创建索引时，我们需要指定一些参数，比如分片数量和副本数量。分片是将索引数据水平切分为多个小块的过程，这样可以提高数据检索和处理的效率。副本则是将索引数据复制到一个或多个节点上，以提高数据的可靠性和查询的可用性。

索引的映射（mapping）是用于定义索引中每个字段的数据类型和其他属性。在创建索引时，需要定义每个字段的数据类型（如文本、数字、日期等）和其他属性（如是否需要分析、是否存储等）。此外，映射还可以定义其他高级功能，如聚合、排序和过滤等。

**文档（Document）**

文档是Elasticsearch中存储和检索的基本单位，它是序列化为JSON格式的数据结构。每个文档都有一个唯一的标识符，称为_id字段，用于唯一标识该文档。每个文档都存储在一个索引中，并且可以包含多个字段，这些字段可以是不同的数据类型，如文本、数字、日期等。

在Elasticsearch中，文档的属性包括_index、_type和_source等。_index表示文档所属的索引名称，_type表示文档所属的类型名称（在早期的Elasticsearch版本中，这是必需的，但在7.x版本之后已经不再需要），_source表示文档的原始JSON数据。

当我们在Elasticsearch中执行搜索查询时，实际上是在查询文档。我们可以使用简单的关键字搜索，也可以使用复杂的查询语句来搜索多个字段。在搜索时，Elasticsearch会使用反向索引来快速定位匹配的文档。反向索引是一个为每个字段建立的倒排索引，它允许Elasticsearch根据关键词在字段中快速查找包含该关键词的文档。

**集群（Cluster）**

一个Elasticsearch集群通常包含了多个节点（Node）和一个或多个索引（Index），并且这些节点和索引共同构成了整个Elasticsearch集群，在所有节点上提供联合索引和搜索功能。

每个Cluster都有一个唯一的名称，即cluster name，它用于标识和区分不同的Elasticsearch集群。

**节点（Node）**

在Elasticsearch集群中，Node是指运行Elasticsearch实例的服务器。每个Node都有自己的名称和标识符，并且都有自己的数据存储和索引存储。

一个Elasticsearch集群由一个或多个Node组成，这些Node通过它们的集群名称进行标识。在默认情况下，如果Elasticsearch已经开始运行，它会自动生成一个叫做“elasticsearch”的集群。我们也可以在配置文件（elasticsearch.yml）中定制我们的集群名字。

Node在Elasticsearch中扮演着不同的角色。根据节点的配置和功能，可以将Node分为以下几种类型：

- Master Node：负责整个Cluster的配置和管理任务，如创建、更新和删除索引，添加或删除Node等。一个Cluster中至少需要有一个Master Node。
- Data Node：主要负责数据的存储和处理，它们可以处理数据的CRUD操作、搜索操作、聚合操作等。一个Cluster中可以有多个Data Node。
- Ingest Node：主要负责对文档进行预处理，如解析、转换、过滤等操作，然后再将文档写入到Index中。每个Cluster中至少需要有一个Ingest Node。 除了上述的三种类型外，还可以有Tribe Node、Remote Cluster Client等特殊用途的Node。

Node之间是对等关系（去中心化），每个节点上面的集群状态数据都是实时同步的。如果Master节点出故障，按照预定的程序，其他一台Node机器会被选举成为新的Master。

需要注意的是，一个Node可以同时拥有一种或几种功能，如一个Node可以同时是Master Node和Data Node。

**分片（Shards）**

在Elasticsearch中，Shards是索引的分片，每个Shard都是一个基于Lucene的索引。当索引的数据量太大时，由于内存的限制、磁盘处理能力不足、无法足够快的响应客户端的请求等，一个节点可能不够用。这种情况下，数据可以被分为较小的分片，每个分片放到不同的服务器上。每个分片可以有零个或多个副本。这不仅能够提高查询效率，还能够提高系统的可靠性和可用性。如果某个节点或Shard发生故障，Elasticsearch可以从其他节点或Shard的副本中恢复数据，从而保证数据的可靠性和可用性。

每个Shard都存储在集群中的某个节点上，每个节点可以存储一个或多个Shard。当查询一个索引时，Elasticsearch会在所有的Shard上执行查询，并将结果合并返回给用户。

对于每个索引，在创建时需要指定主分片的数量，一旦索引创建后，主分片的数量就不能更改。

**副本（Replicas）**

在Elasticsearch中，Replicas是指索引的副本。它们的作用主要有两点：

- 提高系统的容错性。当某个节点发生故障，或者某个分片（Shard）损坏或丢失时，可以从副本中恢复数据。这意味着，即使一个节点或分片出现问题，也不会导致整个索引的数据丢失。这种机制可以增加系统的可靠性，并减少因节点或分片故障导致的宕机时间。
- 提高查询效率。Elasticsearch会自动对搜索请求进行负载均衡，可以将搜索请求分配到多个节点上，从而并行处理搜索请求，提高查询效率。这种负载均衡机制可以在节点之间分发查询请求，使得每个节点都可以处理一部分查询请求，从而避免了一个节点的瓶颈效应。

需要注意的是，在Elasticsearch中，每个索引可以有多个副本（Replicas），但是每个副本只能有一个主分片（Primary Shard）。可以增加或删除副本的数量。

| ES概念 | 关系型数据库 |
| --- | --- |
| Index（索引）支持全文检索 | Table（表） |
| Document（文档），不同文档可以有不同的字段集合 | Row（数据行） |
| Field（字段） | Column（数据列） |
| Mapping（映射） | Schema（模式） |
1. **搭建Elasticsearch环境**

使用docker compose 快速搭建一套Elasticsearch和Kibana环境。

Kibana 提供了一个好用的开发者控制台，非常适合用来练习Elasticsearch命令。

```go
services:
  elasticsearch:
    container_name: elasticsearch
    image: docker.elastic.co/elasticsearch/elasticsearch:8.9.1
    environment:
      - node.name=elasticsearch
      - ES_JAVA_OPTS=-Xms512m -Xmx512m
      - discovery.type=single-node
      - xpack.security.enabled=false
    ports:
      - 9200:9200
      - 9300:9300
    networks:
      - elastic
  kibana:
    image: docker.elastic.co/kibana/kibana:8.9.1
    container_name: kibana
    ports:
      - 5601:5601
    networks:
      - elastic
    depends_on:
      - elasticsearch

networks:
  elastic:

--> docker-compose up
--> http://127.0.0.1:5601
--> Explore on my own
Dev Tools
```

## go-kafka-elasticsearch

1. kafka 记载sql的变化
2. 从kafak中读取到es中并创建或更新

数据流向：mysql→kafka→elasticsearch

```go
type JobWork struct {
	kafkaReader *kafka.Reader
	esClient *EsClient
	log *log.Helper
}
```

1. 配置kafak和es

```go
func NewKafkaReader(conf *conf.Kafka) *kafka.Reader {
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers:   conf.Brokers,
		Topic:     conf.Topic,
		GroupID: conf.GroupId,
	})
}

type EsClient struct{
	*elasticsearch.TypedClient
	index string
}

func NewESClient(cfg *conf.Elasticsearch) (*EsClient, error) {
	// ES 配置
	c := elasticsearch.Config{
		Addresses: cfg.Addresses,
	}

	// 创建客户端连接
	client, err := elasticsearch.NewTypedClient(c)
	if err != nil {
		return nil, err
	}
	return &EsClient{
		TypedClient: client,
		index:       cfg.Index,
	}, nil
}
```

1. 要注册到服务

```go
return kratos.New->kratos.Server

type Server interface {
	Start(context.Context) error
	Stop(context.Context) error
}

// 实现
func (jw JobWork) Start(ctx context.Context) error {
	jw.log.Debug("JobWorker start....")
	// 1. 从kafka中获取MySQL中的数据变更消息
	// 接收消息
	for {
		m, err := jw.kafkaReader.ReadMessage(ctx)
		if errors.Is(err, context.Canceled){
			return nil
		}
		if err != nil {
			jw.log.Errorf("read message failed:%v\n", err)
			break 
		}
		jw.log.Debugf("message at offset %d: %s = %s\n", m.Offset, string(m.Key), string(m.Value))

		// 2.将完整数据写入ES
		msg := new(Msg)
		if err := json.Unmarshal(m.Value, msg); err != nil{
			jw.log.Errorf("unmarshal msg from kafka failed, err:%v", err)
			continue
		}
		if msg.Type == "INSERT" {
			// 往ES中新增文档
			for idx := range msg.Data {
				jw.indexDocument(msg.Data[idx])
			}
		} else {
			// 往ES中更新文档
			for idx := range msg.Data {
				jw.updateDocument(msg.Data[idx])
			}
		}
	}
	return nil
}

func (jw JobWork) Stop(ctx context.Context) error {
	jw.log.Debug("JobWorker stop....")
	// 程序退出前关闭Reader
	return jw.kafkaReader.Close()
}
```

1. 消息体

```go
type Msg struct{
	Type     string `json:"type"`
	Database string `json:"databse"`
	Table    string `json:"table"`
	IsDdl    bool   `json:"isDdl"`
	Data     []map[string]interface{}
}
```

1. es操作

```go
// indexDocument 索引文档
func (jw JobWork) indexDocument(d map[string]interface{}) {
	reviewID := d["review_id"].(string)
	// 添加文档
	resp, err := jw.esClient.Index(jw.esClient.index).
		Id(reviewID).
		Document(d).
		Do(context.Background())
	if err != nil {
		jw.log.Errorf("indexing document failed, err:%v\n", err)
		return
	}
	jw.log.Debugf("result:%#v\n", resp.Result)
}

// updateDocument 更新文档
func (jw JobWork) updateDocument(d map[string]interface{}) {
	reviewID := d["review_id"].(string)
	resp, err := jw.esClient.Update(jw.esClient.index, reviewID).
		Doc(d). // 使用结构体变量更新
		Do(context.Background())
	if err != nil {
		jw.log.Debugf("update document failed, err:%v\n", err)
		return
	}
	jw.log.Debugf("result:%v\n", resp.Result)
}
```

1. 依赖注入

```go
// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewESClient, NewJobWrok, NewKafkaReader)
-->
// wireApp init kratos application.
func wireApp(*conf.Server, *conf.Kafka, *conf.Elasticsearch,  *conf.Data, log.Logger) (*kratos.App, func(), error) {
	panic(wire.Build(server.ProviderSet, data.ProviderSet, biz.ProviderSet, service.ProviderSet, job.ProviderSet, newApp))
}

-->
func newApp(logger log.Logger, gs *grpc.Server, hs *http.Server, js *job.JobWork)
-->
app, cleanup, err := wireApp(bc.Server, bc.Kafka, bc.Elasticsearch, bc.Data, logger)
```

从es中读取数据

1. 配置文件，这里以通过store_id获取信息
    
    ```go
    // 根据商家id查询评价列表（分页）
    	rpc ListReviewByStoreID(ListReviewByStoreIDRequest) returns (ListReviewByStoreIDResponse){}
    	
    	message ListReviewByStoreIDRequest {
    	int64 storeID = 1 [(validate.rules).int64 = {gt: 0}];
    	int32 page = 2 [(validate.rules).int32 = {gt: 0}];
    	int32 size = 3 [(validate.rules).int32 = {gt: 0}];
    }
    
    message ListReviewByStoreIDResponse{
    	repeated ReviewInfo list = 1;
    }
    ```
    
2. 生成代码，业务逻辑
    1. service
        1. 调用biz获取信息
            
            ```go
            reviewList, err := s.uc.ListReviewByStoreID(ctx, req.StoreID, req.Page, req.Size)
            	if err != nil{
            		return nil, err
            	}
            ```
            
        2. 格式化
            
            ```go
            list := make([]*pb.ReviewInfo, 0, len(reviewList))
            	for _, r := range reviewList {
            		list = append(list, &pb.ReviewInfo{
            			ReviewID: r.ReviewID,
            			UserID: r.UserID,
            			OrderID: r.OrderID,
            			Score: r.Score,
            			ServiceScore: r.ServiceScore,
            			ExpressScore: r.ExpressScore,
            			Content: r.Content,
            			PicInfo: r.PicInfo,
            			VideoInfo: r.VideoInfo,
            			Status: r.Status,
            		})
            	}
            ```
            
    2. biz
        1. 对参数进行校验和初始化page, size, limit
            
            ```go
            if page <= 0{
            		page = 1
            	}
            	if size <= 0 || size > 50{
            		size = 10
            	}
            	offset := (page - 1) * size
            	limit := size
            ```
            
        2. 调用data层获取信息
            
            ```go
            uc.log.WithContext(ctx).Debugf("[biz] ListReviewByStoreID storeID:%v\n", storeID)
            ```
            
        3. !!!这里的创建时间会报错，格式不符合go的所以重新定义
            
            ```go
            type MyReviewinfo struct{
            	*model.ReviewInfo
            	CreateAt     MyTime `json:"create_at"` // 创建时间
            	UpdateAt     MyTime `json:"update_at"` // 创建时间
            	Anonymous    int32  `json:"anonymous,string"`
            	Score        int32  `json:"score,string"`
            	ServiceScore int32  `json:"service_score,string"`
            	ExpressScore int32  `json:"express_score,string"`
            	HasMedia     int32  `json:"has_media,string"`
            	Status       int32  `json:"status,string"`
            	IsDefault    int32  `json:"is_default,string"`
            	HasReply     int32  `json:"has_reply,string"`
            	ID           int64  `json:"id,string"`
            	Version      int32  `json:"version,string"`
            	ReviewID     int64  `json:"review_id,string"`
            	OrderID      int64  `json:"order_id,string"`
            	SkuID        int64  `json:"sku_id,string"`
            	SpuID        int64  `json:"spu_id,string"`
            	StoreID      int64  `json:"store_id,string"`
            	UserID       int64  `json:"user_id,string"`
            }
            
            type MyTime time.Time
            
            // UnmarshalJSON json.Unmarshal 的时候会自动调用这个方法
            func (t *MyTime) UnmarshalJSON(data []byte) error {
            	s := strings.Trim(string(data), `"`)
            	tmp, err := time.Parse(time.DateTime, s)
            	if err != nil{
            		return err
            	}
            	*t = MyTime(tmp)
            	return nil
            }
            ```
            
    3. data
        1. 去es查询
            1. 连接es
                
                ```go
                type Data struct {
                	query *query.Query
                	log *log.Helper
                	es *elasticsearch.TypedClient
                }
                
                func NewEsclient(cfg *conf.Elasticsearch) (*elasticsearch.TypedClient, error) {
                	// ES 配置
                	c := elasticsearch.Config{
                		Addresses: cfg.GetAddresses(),
                	}
                	// 创建客户端
                	return elasticsearch.NewTypedClient(c)
                }
                
                // ProviderSet is data providers.
                var ProviderSet = wire.NewSet(NewData, NewReviewRepo, NewDB, NewEsclient)
                ```
                
            2. 查询
                
                ```go
                resq, err := r.data.es.Search().
                		Index("review").
                		From(int(offset)).
                		Size(int(limit)).
                		Query(&types.Query{
                			Bool: &types.BoolQuery{
                				Filter: []types.Query{
                					{
                						Term: map[string]types.TermQuery{
                							"store_id": {Value: storeID},
                						},
                					},
                				},
                			},
                		}).
                		Do(ctx)
                	if err != nil{
                		return nil, err
                	}
                ```
                
            3. 反序列化
                
                ```go
                // 返序列华
                	list := make([]*biz.MyReviewinfo, 0, resq.Hits.Total.Value)
                
                	for _, hit := range resq.Hits.Hits{
                		tmp := &biz.MyReviewinfo{}
                		if err := json.Unmarshal(hit.Source_, tmp); err != nil{
                			r.log.Errorf("json.Unmarshal(hit.Source_, tmp) failed, err:%v", err)
                			continue
                		}
                		list = append(list, tmp)
                	}
                	return list, nil
                ```
                
    4. cmd 依赖注入
        
        ```go
        app, cleanup, err := wireApp(bc.Server, &rc, bc.Data, bc.Elasticsearch, logger)
        
        func wireApp(*conf.Server, *conf.Registry, *conf.Data, *conf.Elasticsearch, log.Logger) (*kratos.App, func(), error) {
        	panic(wire.Build(server.ProviderSet, data.ProviderSet, biz.ProviderSet, service.ProviderSet, newApp))
        }
        ```
        
    5. elasticsearch 配置
        
        ```go
        elasticsearch:
          addresses:
            - "http://127.0.0.1:9200"
            
        message Bootstrap {
          Server server = 1;
          Data data = 2;
          Snowflake snowflake = 3;
          Elasticsearch elasticsearch = 4;
        }
        
        message Elasticsearch{
          repeated string addresses = 1;
        }   
        ```
        

## 设置缓存

1. 配置文件
    
    ```go
    message Data {
      message Database {
        string driver = 1;
        string source = 2;
      }
      message Redis {
        string network = 1;
        string addr = 2;
        google.protobuf.Duration read_timeout = 3;
        google.protobuf.Duration write_timeout = 4;
      }
      Database database = 1;
      Redis redis = 2;
    }
    
      redis:
        addr: 127.0.0.1:6379
        read_timeout: 0.2s
        write_timeout: 0.2s
    ```
    
2. redis连接
    
    ```go
    // NewRedisClient redis连接
    func NewRedisClient(cfg *conf.Data) *redis.Client{
    	return redis.NewClient(&redis.Options{
    		Addr: cfg.Redis.Addr,
    		WriteTimeout: cfg.Redis.WriteTimeout.AsDuration(),
    		ReadTimeout: cfg.Redis.ReadTimeout.AsDuration(),
    	})
    }
    
    // ProviderSet is data providers.
    var ProviderSet = wire.NewSet(NewData, NewReviewRepo, NewDB, NewEsclient, NewRedisClient)
    
    // Data .
    type Data struct {
    	query *query.Query
    	log *log.Helper
    	es *elasticsearch.TypedClient
    	rdb *redis.Client
    }
    ```
    
    依赖注入
    
3. 业务逻辑
    1.  通过singleflight 合并短时间内大量的并发查询
        
        ```go
        var g singleflight.Group
        
        // getDataBySingleflight 合并短时间内大量的并发查询
        func (r *reviewRepo) getDataBySingleflight(ctx context.Context, key string)([]byte, error){
        	v, err, shared := g.Do(key, func() (interface{}, error){
        		// 查缓存
        		data, err := r.getDataFromCache(ctx, key)
        		if err == nil{
        			return data, nil
        		}
        
        		// 缓存中没有, 只有在缓存中没有这个key的错误时才查ES
        		if errors.Is(err, redis.Nil){
        			// 查ES
        			data, err := r.getDataFromEs(ctx, key)
        			if err == nil{
        				// 设置缓存
        				return data, r.setCache(ctx, key, data)
        			}
        			return nil, err
        		}
        
        		// 查缓存失败
        		return nil, err
        	})
        	r.log.Debugf("getDataBySingleflight ret: v:%v, err: %v shared:%v\n", v, err, shared)
        	if err != nil {
        		return nil, err
        	}
        	return v.([]byte), nil
        }
        ```
        
    2. 先查询Redis缓存
        
        ```go
        // getDataFromCache 读取缓存数据
        func  (r *reviewRepo) getDataFromCache(ctx context.Context, key string) ([]byte, error){
        	r.log.Debugf("getDataFromCache key:%v\n", key)
        	return r.data.rdb.Get(ctx, key).Bytes()
        }
        ```
        
    3. 存没有则查询ES
        
        ```go
        // getDataFromEs 从es读取数据
        func (r *reviewRepo) getDataFromEs(ctx context.Context, key string) ([]byte, error){
        	values := strings.Split(key, ":")
        	if len(values) < 4 {
        		return nil, errors.New("invalid key")
        	}
        	index, storeID, offsetStr, limitStr := values[0], values[1],  values[2],  values[3]
        
        	offset, err := strconv.Atoi(offsetStr)
        	if err != nil {
        		return nil, err
        	}
        	limit, err := strconv.Atoi(limitStr)
        	if err != nil {
        		return nil, err
        	}
        
        	resq, err := r.data.es.Search().
        		Index(index).
        		From(offset).
        		Size(limit).
        		Query(&types.Query{
        			Bool: &types.BoolQuery{
        				Filter: []types.Query{
        					{
        						Term: map[string]types.TermQuery{
        							"store_id": {Value: storeID},
        						},
        					},
        				},
        			},
        		}).
        		Do(ctx)
        	if err != nil{
        		return nil, err
        	}
        
        	return json.Marshal(resq.Hits)
        }
        ```
        
    4. 设置缓存
        
        ```go
        // setCache 设置缓存
        func (r *reviewRepo) setCache(ctx context.Context, key string,  data []byte) error {
        	return r.data.rdb.Set(ctx, key, data, time.Second*10).Err()
        }
        ```
        

## 生成api文档

kratos 框架⽀持⽣成 openapi.yaml ⽂件。

1. 安装⽣成openapi⽂件的 protoc 插件
    
    ```go
    go install github.com/google/gnostic/cmd/protoc-gen-openapi@latest
    ```
    
2. 项⽬根⽬录下执⾏以下命令根据API的 proto⽂件⽣成 openapi.yaml⽂件
    
    ```go
    make api
    ```
    
3. [swagger.io](http://swagger.io/) 提供了开源的 Swagger Editor ，直接导⼊项⽬⽬录下的 openapi.yaml ⽂件皆可。
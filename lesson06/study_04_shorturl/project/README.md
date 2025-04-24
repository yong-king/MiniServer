# 短链接项目

## 什么是短链接

通俗来说就是将⽐较⻓的⼀个URL⽹址，通过程序计算等⽅式，转换为简短的⽹址字符串。

ysh&dlrb — > https://youngking.com/dlrb

在许多短信里附带的连接就是短链接

**为什么要⽤短⽹址/短链接？**

公司内部有很多需要发送链接的场景，业务侧的链接通常会⽐较⻓，在发送短信、IM⼯具发送消息、push等场景
下⻓链接有以下劣势：

1. 短信内容超⻓，1条消息被拆分成多条短信发送，浪费钱。
2. 微博等平台有字数限制。
3. ⻜书、钉钉等IM⼯具对⻓链接（带特殊服务号的）识别有问题。
4. 短链接转成⼆维码更清晰。

短链接代码示例

```go
// api
syntax = "v1"

type Request {
	ShortURL string `path:"shortURL"`
}

type Response {
	LongURL string `json:"longURL"`
}

service shorurl-api {
	@handler ShorurlHandler
	get /:shortURL (Request) returns (Response)
}

// logic.shorturllogit.go
func (l *ShorurlLogic) Shorurl(req *types.Request) (resp *types.Response, err error) {
	// todo: add your logic here and delete this line
	// 将请求的端连接转换为长连接
	// ysh&dlrb --> https://baidu.com
	if req.ShortURL == "ysh&dlrb" {
		return &types.Response{LongURL: "https://baidu.com"}, nil
	}
	// 如果查询不到，就跳转为https://google.com
	return &types.Response{LongURL: "https://google.com"}, nil
}

// handler.shorturlhandler.go
func ShorurlHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.Request
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewShorurlLogic(r.Context(), svcCtx)
		resp, err := l.Shorurl(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			// httpx.OkJsonCtx(r.Context(), w, resp)

			// w.Header().Set("location", resp.LongURL) // 重定向
			// w.WriteHeader(http.StatusFound) // 状态码

			http.Redirect(w, r, resp.LongURL, http.StatusFound) // 重定向
		}
	}
}

httpx.OkJsonCtx -> func OkJsonCtx(ctx context.Context, w http.ResponseWriter, v any)
WriteJsonCtx(ctx, w, http.StatusOK, v) ->  doWriteJson(w, code, v)
doWriteJson(w, code, v) -> 
	w.Header().Set(ContentType, header.ContentTypeJson)
	w.WriteHeader(code)
```

# 项目

**需求背景**

**为什么要设计短链系统？**

公司内部业务需要发送⼤量的营销短信、通知类短信。

需要⼀个短链接服务满⾜各业务线的使⽤。

提供转链接⼝

后续⽀持提供点击的统计数据报表

**需求描述**

1. 输⼊⼀个⻓⽹址得到⼀个唯⼀的短⽹址。
2. ⽤户点击短⽹址能够正常跳转到对应的⽹址。
3. 为了保证业务的延续性，短⽹址⻓期有效。

**需求分析**

**产品定位**

1. 公司内部业务使⽤的短⽹址服务，只接收公司内部的⻓链转短链需求。（不对外提供短链功能。）

2. 基本在国内使⽤（点击链接的⽤户绝⼤多数为国内⽤户）

3. 后续可能会要求提供短链的访问数据报表

**规模**

1. ⼤致服务于公司内部x条业务线。

2. ⼤致服务的⽤户规模有x亿。

3. xx QPS

**技术指标**

1. 延时x ms内

2. 可靠性99.99%

3. 安全性

**需求拆解**

根据需求分析，可以将需求拆分为**转链模块**、**存储**和**访问链接模块**。

**转链模块**

1. 相同的⻓链要转为同⼀个短链

2. ⽣成的短链为尽量短的字符。

作为⼀个开发想得再多⼀点，引申出来的需求点或注意事项。

1. 需要避免某些不合适的词（例如 f**k 、 stupid ）
2. 避免⽣成的短链出现某些特殊含义的词 version 、 health 等
3. 避免循环转链（把已经是短链的再拿来转短链）

**存储**

1. 保存原始⻓链接与短链接的对应关系

2. 能够根据短链接查找到原始的⻓链接

**查看链接模块**

1. 根据短链查询到⻓链后返回重定向响应。

2. 后续数据报表需求可能需要采集并统计请求头数据。

**系统设计**

**总体设计⽅案**

通过分析可以得知，这是⼀个典型的**读多写少**的系统。

并且我们进⼀步分析这个短链系统区别于其他读多写少的业务场景，它的特点是数据写⼊后基本不会改变。（好处

是不需要考虑数据⼀致性的问题，可以放⼼⼤胆的使⽤缓存系统来提⾼读的效率）

**短链⽣成⽅式**

关于⽣成短链有以下⼏种⽅案，

**hash**

使⽤hash函数对⻓链接进⾏hash，得到hash值作为短链标识符。

优势：简单

缺点：数据量⼤之后，会出现哈希冲突

扩展：
MurmurHash 是⼀种⾮加密型哈希函数，和其它流⾏哈希函数相⽐，对于规律性较强的key随机分布特性表现更良
好，在很多开源的软件项⽬（Redis，Memcached，Cassandra，HBase，Lucene都⽤它）都有使⽤。有以下⼏
个特性：
随机分布特性表现好
算法速度快

**发号器/⾃增序列**

每收到⼀个转链请求，就使⽤发号器⽣成递增（1、2、3、4...以此递增）的序号，然后将该序号转为**62进制**，最后

拼接到短域名后即得到最终的短链。

**发号器⽅案的优劣如下**

优势

⽣成的id递增

理论上容量⾜够满⾜现实需求

缺点：

⾼并发下的发号器设计是难点。

**发号器实现⽅式**

常⻅的发号器实现⽅式有以下⼏种：

1. 基于uuid实现

优势：不会重复、性能好

劣势：数字太⼤了，32位16进制数

2. 基于redis实现发号器

优势：⾼性能

劣势：需搭建⾼可⽤架构并考虑持久化

3. 基于雪花算法的分布式ID⽣成器

优势：⾼性能、⾼可⽤

劣势：实现复杂，依赖时钟

4. 基于MySQL⾃增主键的发号器

优势：简单、可靠劣势：依赖MySQL，性能会成为瓶颈，但可通过分⽚扩展可⽤性

## 1.创建数据库

### sequence数据库：发号数据库

```sql
CREATE TABLE `sequence` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `stub` varchar(1) NOT NULL,
  `timestamp` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_uniq_stub` (`stub`)
) ENGINE=MyISAM DEFAULT CHARSET=utf8 COMMENT = '序号表';
```

### short_url_map数据库：短链接映射长链接数据库

```sql
CREATE TABLE `short_url_map` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键',
    `create_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `create_by` VARCHAR(64) NOT NULL DEFAULT '' COMMENT '创建者',
    `is_del` tinyint UNSIGNED NOT NULL DEFAULT '0' COMMENT '是否删除：0正常1删除',
    
    `lurl` varchar(2048) DEFAULT NULL COMMENT '长链接',
    `md5` char(32) DEFAULT NULL COMMENT '长链接MD5',
    `surl` varchar(11) DEFAULT NULL COMMENT '短链接',
    PRIMARY KEY (`id`),
    INDEX(`is_del`),
    UNIQUE(`md5`),
    UNIQUE(`surl`)
)ENGINE=INNODB DEFAULT CHARSET=utf8mb4 COMMENT = '长短链映射表';

// 实际应用中最好是建立连个数据库，防止一个崩掉导致整个都崩掉
```

**分⽚部署**

为了避免单点故障，我们将我们的ID⽣成器分成奇数和偶数两部分，分别部署在两个MySQL服务器。

两个数据表配置不同的 auto-increment-offset ， server1 ⽣成1、3、5、7、9...， server2 ⽣成2、4、6、8...。

## 搭建go-zero框架骨架

### api文件，使用goctl生成代码

```bash
syntax = "v1"

info(
    title: "shortener"
    desc: "shorturl to longurl"
    author: "ysh"
    email: "youngking98.com"
    version: "1.0"
)

type ConvertRequest{
    LongUrl string `json:"longUrl"`
}
type ConvertResponse{
    ShortUrl string `json:"shortUrl"`
}

type ShowResquest{
    ShortUrl string `path:"shortUrl"`
}
type ShowResponse{
    LongUrl string `json:"longUrl"`
}

@server(
    prefix: /v1
)
service shortener-api {
    @handler ConvertHandler
    post /convert(ConvertRequest) returns(ConvertResponse)
    
    @handler ShowHandler
    get /:shortUrl(ShowResquest) returns(ShowResponse)
}
```

根据api文件生成go代码

```bash
goctl api go -api shortener.api -dir . -style=goZero
```

根据数据表生成model代码

```bash
goctl model mysql datasource -url="root:password@tcp(127.0.0.1:3306)/database" -table="short_url_map" -dir="./model" -c  
goctl model mysql datasource -url="root:password@tcp(127.0.0.1:3306)/database" -table="sequence" -dir="./model" -c  
```

下载依赖

```bash
go mod tidy
```

修改配置文件

```bash
// yaml
type Config struct {
	rest.RestConf

	// sql
	ShortUrlMapDb ShortUrlMapDb
	SequenceDB struct{
		DNS string
	}
}

type ShortUrlMapDb struct{
	DNS string
}

// 实际应用中最好是建立连个数据库，防止一个崩掉导致整个都崩掉

// config.go
ShortUrlMapDb:
  DSN: root:password@tcp(127.0.0.1:3306)/db1?charset=utf8mb4&parseTime=true&loc=Asia%2FShanghai

SequenceDB:
  DSN: root:password@tcp(127.0.0.1:3306)/db1?charset=utf8mb4&parseTime=true&loc=Asia%2FShanghai
```

### 业务逻辑

1. 参数校验
    1. 长链接不能为空
        1. validate校验,在路由中处理，如果路由都不能通过，则不会进来
        2. 在api文件中加入
            
            ```bash
            validate:"required
            ```
            
        3. 使用
            
            ```bash
            // handler.go
            validator.New().StruceCtx(context, &req)
            ```
            
    2. 长链接必须是能连通的
        1. 基于http.Get()
        
        ```go
        var client = &http.Client{
        	Transport: &http.Transport{
        		DisableKeepAlives: true,
        	},
        	Timeout: 2 * time.Second,
        }
        
        resp, err := client.Get(url)
        ```
        
    3. 判断之前是否已经转链过
        1. 长链生成md5
        
        ```go
        h := md5.New()
        h.Write(data)
        return hex.EncodeToString(h.Sum(nil)) 
        ```
        
        1.  长链生成的md5是否在数据库中
        
        ```go
        l.svcCtx.ShortUrlDB.FindOneByMd5(l.ctx, sql.NullString{String: md5Value, Valid: true})
        ```
        
    4. 输入的不能是一个短链接（循环转链）
        1. 基于url.Path和path.Base
            
            ```go
            myUrl, err := url.Parse(targetUrl)
            path.Base(myUrl.Path), nil
            ```
            
        2. 到数据库中查询是否存在这个短链接
            
            ```go
            l.svcCtx.ShortUrlDB.FindOneBySurl(l.ctx, sql.NullString{String: basePath, Valid: true})
            ```
            
2. 取号

对于每一个转链请求，都在sequence中生成一条数据，并取出主键id为号码

有两种取号的方法，一是基于mysql取号，而是基于redis取号

1. mysql取号
    
    ```sql
    `REPLACE INTO sequence (stub) VALUES ('a')` 
    // 基于这个REPLACE INTO，并取出主键id
    // 这里的数据库表中永远只有一条数据
    ```
    
    1. 定义结构体
    
    ```sql
    type MySQL struct{
    	conn sqlx.SqlConn
    }
    
    func NewMysql(dns string) Sequence {
    	return &MySQL{
    		conn: sqlx.NewMysql(dns),
    	}
    }
    ```
    
    1. 预编译sql语句
    
    ```sql
    stmt, err = m.conn.Prepare(sqlReplaceIntoSub)
    defer stmt.Close()
    ```
    
    1. 执行语句
    
    ```sql
    rest, err = stmt.Exec()
    ```
    
    1. 取出刚刚插入的id
    
    ```sql
    lid, err = rest.LastInsertId()
    ```
    
    1. 返回
2. redis取号
    
    基于redis的Incr原子操作
    
    1. 定义结构体
        
        ```sql
        type Redis struct{
        	rds *redis.Redis
        }
        
        func NewRedis(conf redis.RedisConf) Sequence{
        	return &Redis{
        		rds: redis.MustNewRedis(conf),
        	}
        }
        ```
        
    2. 获取下一个字增的序列号
        
        ```sql
        val, err = r.rds.Incr(sequenceKey)
        ```
        
    3. 返回
    
    <aside>
    💡
    
    这里需要注意的是！！！
    
    在配置文件中加入redis的配置
    
    ```go
    // api
    SequenceRDB:
      Host: 127.0.0.1:6379
      Pass: youngking98
    // config
    SequenceRDB redis.RedisConf
    ```
    
    </aside>
    
3. 这里可以定义一个接口，这样就可以只用更改调用，而不用管我下面到底用到redis还是sql
    
    ```go
    type Sequence interface{
    	Next() (uint64, error)
    }
    ```
    
4. 修改服务上下文配置
    
    ```go
    type ServiceContext struct {
    	Config config.Config
    	SequenceDB sequence.Sequence
    }
    
    func NewServiceContext(c config.Config) *ServiceContext {
    	conn := sqlx.NewMysql(c.ShortUrlMapDb.DSN)
    	return &ServiceContext{
    		Config: c,
    		// SequenceDB: sequence.NewMysql(c.SequenceDB.DNS),
    		SequenceDB: sequence.NewRedis(c.SequenceRDB),
    	}
    }
    ```
    
5. 业务逻辑
    
    ```go
    seq, err := l.svcCtx.SequenceDB.Next()
    	if err != nil {
    		logx.Errorw("SequenceDB.Next failed", logx.LogField{Key: "err", Value: err.Error()})
    		return nil, err
    	}
    ```
    

1. 号码转链
    1. uint64 —> string
        
        ```go
        func Int2String(seq uint64) string{
        	if seq < 61 {
        		return string(base62String[seq])
        	}
        	bl := []byte{}
        	for seq > 0{
        		mod := seq % 62
        		div := seq / 62
        		bl = append(bl, base62String[mod])
        		seq = div
        	}
        	return string(reverse(bl))
        }
        
        base62String // 在配置文件中定义，在程序启动时加载
        ```
        
    2. 排除敏感词，建立敏感词黑名单
        
        ```go
        1. 在配置文件中定义黑名单
        yaml, config
        2. 在上下文配置文件中定义
        severContext.go 用map定义
        if _, ok := l.svcCtx.ShortUrlBlackList[short]; !ok {
        			break // 生成不在黑名单里的短链接就跳出for循环
        		}
        ```
        
    3. 如果出现敏感词就一直循环，直到正确

1. 查看短链
    1. 根据锻炼到数据库中查询原始长链
        
        ```go
        u, err := l.svcCtx.ShortUrlDB.FindOneBySurl(l.ctx, sql.NullString{String: req.ShortUrl, Valid: true})
        ```
        
    2. 返回查询到的长链，在调用handler层返回重定位响应
        
        ```go
        // 返回重定向的响应
        http.Redirect(w, r, resp.LongUrl, http.StatusFound)
        ```
        
    3. 这里如果大量的查询，会导致都去查询数据库，数据库承载不起，所以可以引入缓存Redis
        1. 方法一：使用go-zero生成的带存的方法，上面已经提到过，但是有一个问题就是存入了一些无用的信息
            
            ```go
            key: "cache:db1:shortUrlMap:surl:{D true}"
            ```
            
        2. 方法二，自己链接Redis传入
            
            ```go
            key: "shortlink:C"
            ```
            
            我们上面实现了sequence的接口，在实现一个LinkMapping接口
            
            ```go
            type LinkMapping interface {
                SetShortLink(shortLink, longLink string) error
                GetLongLink(shortLink string) (string, error)
            }
            // 用 Redis 实现
            // 实现 LinkMapping 接口的 SetShortLink 方法，存储短链接和长链接的映射
            func (r *Redis) SetShortLink(shortLink, longLink string) error {
                err := r.rds.Set("shortlink:"+shortLink, longLink)
                if err != nil {
                    logx.Errorw("rds.Set failed", logx.LogField{Key: "err", Value: err.Error()})
                    return err
                }
                return nil
            }
            
            // 实现 LinkMapping 接口的 GetLongLink 方法，查询短链接对应的长链接
            func (r *Redis) GetLongLink(shortLink string) (string, error) {
            	if shortLink == "" {
            		return "", errors.New("need shortlink")
            	}
                longLink, err := r.rds.Get("shortlink:" + shortLink)
                if err != nil {
                    logx.Errorw("rds.Get failed", logx.LogField{Key: "err", Value: err.Error()})
                    return "", err
                }
            
                if longLink == "" {
                    return "", fmt.Errorf("short link not found")
                }
            
                return longLink, nil
            }
            
            func NewRedis(conf redis.RedisConf) ***Redis**{
            	return &Redis{
            		rds: redis.MustNewRedis(conf),
            	}
            }
            ```
            
            修改serviceContext.go
            
            ```go
            LinkMappingDB sequence.LinkMapping //struct
            
            LinkMappingDB: sequence.NewRedis(c.SequenceRDB), // func
            
            // 这样外部只用调用LinkMappingDB接口，就可以使用NewRedis中的方法了，
            // 当一个结构体实现了一个接口后，该接口的实例就能够调用结构体中实现的相关方法
            ```
            
            业务逻辑
            
            ```go
            // 保存长短链接映射到redis
            l.svcCtx.LinkMappingDB.SetShortLink(short, req.LongUrl)
            
            // 查询长短链接
            ul, err := l.svcCtx.LinkMappingDB.GetLongLink(req.ShortUrl)
            ```
            

### 解决缓冲相关问题

**缓存相关问题**

使⽤Redis作为缓存，那么就需要考虑⼏个核⼼问题。

1. 缓存怎么设置，LRU
    
    1. Redis集群部署
    
    2. 根据数据量设置内存⼤⼩，内存淘汰策略LRU，移除最近最少使⽤的key。
    
2. 如果解决缓存击穿问题？引申：什么是缓存雪崩、缓存击穿、缓存穿透
    
    1. 过期时间设⼤
    
    2. 加锁
    
    3. 使⽤singleflight 合并请求
    
    singleflight: 提供了重复函数调用抑制机制，使用它可以避免同时进行相同的函数调用。第一个调用未完成时后续的重复调用会等待，当第一个调用完成时则会与它们分享结果，这样以来虽然只执行了一次函数调用但是所有调用都拿到了最终的调用结果。
    
    `singleflight`包中定义了一个名为`Group`的结构体类型，它表示一类工作，并形成一个命名空间，在这个命名空间中，可以使用重复抑制来执行工作单元。
    
    `Group`类型有以下三个方法。
    
    `Do` 执行并返回给定函数的结果，确保一次只有一个给定key在执行。如果进入重复调用，重复调用方将等待原始调用方完成并会收到相同的结果。返回值`shared`表示是否给多个调用方赋值 v。
    
    需要注意的是，使用`Do`方法时，如果第一次调用发生了阻塞，那么后续的调用也会发生阻塞。在极端场景下可能导致程序hang住。
    
    `singleflight`包提供了`DoChan`方法，支持我们异步获取调用结果。
    
    `DoChan` 类似于 `Do`，但不是直接返回结果而是返回一个通道，该通道将在结果准备就绪时接收结果。返回的通道将不会关闭。
    
    `Result` 保存 `Do` 的结果，因此它们可以在通道上传递。
    
    为了避免第一次调用阻塞所有调用的情况，我们可以结合使用select和`DoChan`为函数调用设置超时时间。
    
    如果在某些场景下允许第一个调用失败后再次尝试调用该函数，而不希望同一时间内的多次请求都因第一个调用返回失败而失败，那么可以通过调用`Forget`方法来忘记这个key。
    
    `Forget`告诉`singleflight`忘记一个key。将来对这个key的 `Do` 调用将调用该函数，而不是等待以前的调用完成。
    
    go-zero天生支持singleflight
    
    ```go
    l.svcCtx.ShortUrlDB.FindOneBySurl(l.ctx, sql.NullString{String: req.ShortUrl, Valid: true})
    -->
    FindOneBySurl(ctx context.Context, surl sql.NullString) (*ShortUrlMap, error)
    --> 转到实现
    m.QueryRowIndexCtx(ctx, &resp, db1ShortUrlMapSurlKey, m.formatPrimary, func(ctx context.Context, conn sqlx.SqlConn, v any) (i any, e error) {
    		query := fmt.Sprintf("select %s from %s where `surl` = ? limit 1", shortUrlMapRows, m.table)
    -->
    QueryRowIndexCtx -> cc.cache.TakeWithExpireCtx
    --> 
    TakeWithExpireCtx
    	|New -->  NewNode(redis.MustNewRedis(c[0].RedisConf), barrier, st, errNotFound, opts...)
    -->
    func NewNode(rds *redis.Redis, **barrier syncx.SingleFlight**, st *Stat,
    	errNotFound error, opts ...Option) 45
    ```
    
3. 如何解决缓存穿透问题？
    
    什么是缓存穿透？
    攻击者恶意请求短链服务，短时间⼤量请求并不存在的短链
    
    1. 布隆过滤器（简单，如果不在那⼀定不在）
        1. 为什么需要使⽤布隆过滤器？
        2. 节省空间。并不存储原始数据，只⽤来判断某个元素是否存在。
        3. 防⽌缓存穿透、推荐系统去重、⿊⽩名单、垃圾邮件过滤
        
        <aside>
        💡
        
        假设集合里面有3个元素{x, y, z}，哈希函数的个数为3。首先将位数组进行初始化，将里面每个位都设置位0。对于集合里面的每一个元素，将元素依次通过3个哈希函数进行映射，每次映射都会产生一个哈希值，这个值对应位数组上面的一个点，然后将位数组对应的位置标记为1。查询W元素是否存在集合中的时候，同样的方法将W通过哈希映射到位数组上的3个点。如果3个点的其中有一个点不为1，则可以判断该元素一定不存在集合中。反之，如果3个点都为1，则该元素可能存在集合中。注意：此处不能判断该元素是否一定存在集合中，可能存在一定的误判率。可以从图中可以看到：假设某个元素通过映射对应下标为4，5，6这3个点。虽然这3个点都为1，但是很明显这3个点是不同元素经过哈希得到的位置，因此这种情况说明元素虽然不在集合中，也可能对应的都是1，这是误判率存在的原因。
        
        </aside>
        
    2. 布⾕⻦过滤器（⽀持删除）
        
        使用两个哈希函数对一个`key`进行哈希，得到桶中的两个位置，此时
        
        - 如果两个位置都为为空则将`key`随机存入其中一个位置
        - 如果只有一个位置为空则存入为空的位置
        - 如果都不为空，则随机踢出一个元素，踢出的元素再重新计算哈希找到相应的位置
        
        当然假如存在绝对的空间不足，那老是踢出也不是办法，所以一般会设置一个**踢出阈值**，如果在某次插入行为过程中连续踢出超过阈值，则进行扩容。
        
        基本的布谷鸟过滤器也是由两个或者多个哈希函数构成，布谷鸟过滤器的布谷鸟哈希表的基本单位称为**条目（entry）**。 每个条目存储一个**指纹（fingerprint）**，指纹指的是使用一个哈希函数生成的n位比特位，n的具体大小由所能接受的误判率来设置，论文中的例子使用的是8bits的指纹大小。
        
        哈希表由一个桶数组组成，其中一个桶可以有多个条目（比如上述图c中有四个条目）。而每个桶中有四个指纹位置，意味着一次哈希计算后布谷鸟有四个“巢“可用，而且四个巢是连续位置，可以更好的利用cpu高速缓存。也就是说每个桶的大小是4*8bits。
        
        给定一个项x，算法首先根据上述插入公式，计算x的指纹和两个候选桶。然后读取这两个桶：如果两个桶中的任何现有指纹匹配，则布谷鸟过滤器返回true，否则过滤器返回false。此时，只要不发生桶溢出，就可以确保没有假阴性。
        
        标准布隆过滤器不能删除，因此删除单个项需要重建整个过滤器，而计数布隆过滤器需要更多的空间。布谷鸟过滤器就像计数布隆过滤器，可以通过从哈希表删除相应的指纹删除插入的项，其他具有类似删除过程的过滤器比布谷鸟过滤器更复杂。
        
        具体删除的过程也很简单，检查给定项的两个候选桶；如果任何桶中的指纹匹配，则从该桶中删除匹配指纹的一份副本。
        
    
    ```go
    //1.基于redis的，go-zero自带
    // serviceContext.go
    Filter *bloom.Filter
    
    // 初始化布隆过滤器
    	store := redis.New(c.CacheRedis[0].Host)
    	filter := bloom.New(store, "bloom_filter", 20*(1<<20))
    	Filter: filter,
    	
    // 业务逻辑
    	// 4.2将生成的短链接加入到布隆过滤器
    	err = l.svcCtx.Filter.Add([]byte(short))
    	
    	// 1.0 布隆过滤器
    	// 不存在的短链接直接返回404,不需要后续处理
    	exist, err := l.svcCtx.Filter.Exists([]byte(req.ShortUrl))
    	if err != nil {
    		logx.Errorw("Filter.Exists failed", logx.LogField{Key: "err", Value: err.Error()})
    		return nil, err
    	}
    	if !exist{
    		return nil, errors.New("404")
    	}
    	
    	// 2. 基于内存的，每次关机都得重新加载,这里不是基于go-zero的。
    	
    import (
    	"errors"
    
    	**"github.com/bits-and-blooms/bloom/v3"**
    	"github.com/zeromicro/go-zero/core/logx"
    	"github.com/zeromicro/go-zero/core/stores/sqlx"
    )
    
    var filter = bloom.NewWithEstimates(1<<20, 0.01)
    
    func loadDataToBloomFilter(conn sqlx.SqlConn, filter *bloom.BloomFilter) error {
    	if conn == nil || filter == nil{
    		return errors.New("loadDataToBloomFilter invalid param")
    	}
    
    	// 查总数
    	total := 0
    	if err := conn.QueryRow(&total, "select count(*) from short_url_map where is_del=0"); err != nil{
    		logx.Errorw(" conn.QueryRow failed", logx.LogField{Key: "err", Value: err.Error()})
    		return err
    	}
    	logx.Infow("total data", logx.LogField{Key: "total", Value: total})
    	if total == 0 {
    		logx.Info("no data need to load")
    		return nil
    	}
    
    	pageTotal := 0
    	pageSize := 20
    	if total%pageSize == 0{
    		pageTotal = total/pageSize
    	}else {
    		pageTotal = total/pageSize + 1
    	}
    	logx.Infow("pageTotal", logx.LogField{Key: "pageTotal", Value: pageTotal})
    
    	// 循环查询所有数据
    	for page := 1; page <= pageTotal; page++ {
    		offset := pageSize * (pageTotal - 1)
    		surls := []string{}
    		if err := conn.QueryRow(&surls, "select surl from short_url_map where is_del=0 limit ?,?", offset, pageSize); err != nil {
    			return err
    		}
    		for _, surl := range surls{
    			filter.AddString(surl)
    		}
    	}
    	logx.Info("load data to bloom success")
    	return nil
    }
    ```
    

### 编写测试单元

1. 方式一：右键点击函数，选择Go: Generate Unit Test Founction
2. 方式二：创建测试单元文件夹，go文件名+”_test “, 方法名为Test+要测试的方法名称
    1. **GoConvey**
    
    ```go
    go get github.com/smartystreets/goconvey
    
    import (
    	"testing"
    
    	c "github.com/smartystreets/goconvey/convey"  // 别名导入
    )
    
    func TestXXX(t *testing.T) {
    	c.Convey("基础用例", t, func() {
    		var (
    			// 你的方法中的参数
    		)
    		got := XXX(xxx, xxx, ...)
    		c.So(got, c.ShouldResemble, expect) // 断言
    	})
    }
    ```
    

### 部署

部署该项⽬的⼀种推荐⽅法是在通过 Nginx 代理，即将我们的短链服务部署在 Nginx 后。
通过这种⽅式，可以通过 Nginx 的访问⽇志（access.log）来统计访问数据。（例如通过EFK采集⽇志，统计报
表）

⻓链转短链：
单独部署为⼀个微服务（转链服务）
对其他服务提供转链服务，需要鉴权（接你们公司鉴权）。
通过RESTful API调⽤我们的转链接⼝
通过RPC⽅式调⽤我们的转链⽅法（⾃⼰实现⼀个RPC版本的转链）

查看短链接：
单独部署为⼀个服务（查看短链服务）
通过nginx转发查看请求， /[0-9a-zA-Z]* --> 转发到我们的查链服务
通过 access.log 收集（EFK）并统计访问数据

## 扩展

1.  如何⽀持**⾃定义短链**？
    
    维护⼀个已经使⽤的序号，后续⽣成序号时判断是否已经被分配。
    
2.  如何让短链⽀持过期时间？
    
    每个链接映射额外记录⼀个『过期时间』字段，到期后将该映射记录删除。关于删除的策略有以下⼏种：
    
    1. 延迟删除：每次请求时判断是否过期，如果过期则删除。
    
    实现简单，性能损失⼩
    
    存储空间的利⽤效率低，已经过期得数据可能永远不会被删除
    
    2. 定时删除：创建记录时根据过期时间设置定时器
    
    过期数据能被及时删除，存储空间的利⽤率⾼
    
    占⽤内存⼤，性能差
    
    3. 轮询删除：通过异步脚本在业务低峰期周期性扫表清理过期数据兼顾效率和磁盘利⽤率
    
3. 如何提⾼吞吐量？
整个系统分为『⽣成短链（写）』和『访问短链（读）』两部分
    1. ⽔平扩展多节点，根据序号分⽚
4. 延迟优化
整个系统分为『⽣成短链（写）』和『访问短链（读）』两部分
    1. 存储层
        1. 数据结构简单可以直接改⽤kv存储
        2. 对存储节点进⾏分⽚
    2. 缓存层
        1. 增加缓存层，本地缓存-->redis缓存
        2. 使⽤布隆过滤器判断⻓链接映射是否已存在，判断短链接是否有效
    3. ⽹络
        1. 基于地理位置就近访问数据节点
        

**短链接项⽬示例**

**项⽬介绍**：⼀个⽤于公司内部营销短信和App push的短链接服务，包含转链、存储、链接跳转功能，并可提供短链接点击数据报表功能。

**个⼈职责**：

负责项⽬的整体设计和开发, 负责实现转链和链接跳转模块逻辑。

基于MySQL主键实现了⾼可⽤的发号器组件。

在转链前进⾏特殊词过滤和防⽌循环转链的校验处理。

查看链接服务采⽤布隆过滤器防⽌缓存穿透，使⽤singleflight防⽌缓存击穿。

查看链接服务单独部署，Nginx转发请求通过EFK采集access⽇志⽅式统计短链接的点击数据。

**项⽬收获**

熟悉了常⽤的发号器设计⽅案，能够结合实际情况选择最适合的⽅案。

熟悉了go-zero框架的使⽤，对Go语⾔操作MySQL和Redis都更加熟悉。

锻炼了⾃⼰的⾃主学习能⼒，积累了项⽬设计和开发的经验。

1. 说清楚项⽬背景
    1. 内部业务对外营销需要有⼀个短链接功能，并且能回收点击数据。
2. 说清楚项⽬架构
    1. 能够画出实际的架构图，并能够清楚的说出每个组件的功能。
3. 项⽬是如何部署的？
    1. 转链单独作为⼀个微服务部署。
    2. 其他项⽬通过 API接⼝和RPC⽅式接⼊（创建短链接）
    3. 查看短链也是单独部署，接前⾯是nginx，通过nginx的access⽇志统计点击数据。
4. 说清楚项⽬实现过程中的重点和难点
    1. 为什么使⽤302跳转，⽽不使⽤301跳转？301与302的区别是什么？
        1. 需要记录访问数据，如果⽤301永久重定向跳转，下⼀次访问时浏览器有缓存就不再请求短链服务器了。这样会丢失访问数据。302是临时重定向，每次访问短链都会去请求短链服务器（除⾮响应中⽤ Cache-Control 或 Expired 暗示浏览器缓存）,虽然⽤ 302重定向会给 server 增加⼀点压⼒，但是能准确记录每⼀次短链请求数据。
        2. 防浏览器缓存，领导让你把这个短链接ban了，你短链服务端删除了数据，但是浏览器有缓存还是能访问。
    2. **HTTP 响应状态码**
        
        HTTP 响应状态码用来表明特定 [HTTP](https://developer.mozilla.org/zh-CN/docs/Web/HTTP) 请求是否成功完成。 响应被归为以下五大类：
        
        1. [信息响应](https://developer.mozilla.org/zh-CN/docs/Web/HTTP/Reference/Status#%E4%BF%A1%E6%81%AF%E5%93%8D%E5%BA%94) (`100`–`199`)
        2. [成功响应](https://developer.mozilla.org/zh-CN/docs/Web/HTTP/Reference/Status#%E6%88%90%E5%8A%9F%E5%93%8D%E5%BA%94) (`200`–`299`)
        3. [重定向消息](https://developer.mozilla.org/zh-CN/docs/Web/HTTP/Reference/Status#%E9%87%8D%E5%AE%9A%E5%90%91%E6%B6%88%E6%81%AF) (`300`–`399`)
        4. [客户端错误响应](https://developer.mozilla.org/zh-CN/docs/Web/HTTP/Reference/Status#%E5%AE%A2%E6%88%B7%E7%AB%AF%E9%94%99%E8%AF%AF%E5%93%8D%E5%BA%94) (`400`–`499`)
        5. [服务端错误响应](https://developer.mozilla.org/zh-CN/docs/Web/HTTP/Reference/Status#%E6%9C%8D%E5%8A%A1%E7%AB%AF%E9%94%99%E8%AF%AF%E5%93%8D%E5%BA%94) (`500`–`599`)
        
        [https://developer.mozilla.org/zh-CN/docs/Web/HTTP/Reference/Status](https://developer.mozilla.org/zh-CN/docs/Web/HTTP/Reference/Status)
        
        **http 状态码301、302、303、307、308 的区别**
        
        **301 Moved Permanently 永久重定向,默认情况下，永久重定向是会被浏览器缓存的。**
        
        **302 Found 临时重定向 与 301 状态码类似；但是，客户端应该使用 Location 首部给出的 URL 来临时定位资源。将来的请求仍应使用老的 URL。在浏览器的实现中，302默认以get重新发出请求。**
        
        **303 See Other 临时重定向** 
        
        虽然 [RFC 1945](https://links.jianshu.com/go?to=https%3A%2F%2Ftools.ietf.org%2Fhtml%2Frfc1945) 和 [RFC 2068](https://links.jianshu.com/go?to=https%3A%2F%2Ftools.ietf.org%2Fhtml%2Frfc2068) 规范不允许客户端在重定向时改变请求的方法，但是很多现存的浏览器在收到302响应时，直接使用GET方式访问在Location中规定的URI，而无视原先请求的方法。[[2]](https://links.jianshu.com/go?to=https%3A%2F%2Fzh.wikipedia.org%2Fwiki%2FHTTP_303%23cite_note-ruby-on-rails-ActionController-Redirecting-redirect_to-2)因此状态码303被添加了进来，用以明确服务器期待客户端进行何种反应。[[3]](https://links.jianshu.com/go?to=https%3A%2F%2Fzh.wikipedia.org%2Fwiki%2FHTTP_303%23cite_note-RFC7230-10-3)重定向到新地址时，客户端必须使用GET方法请求新地址。
        
        **307 Temporary Redirect**
        
        这个状态码和302相似，有一个唯一的区别是不允许将请求方法从post改为get。
        
        **308 Permanent Redirect 永久重定向**
        
        此状态码类似于301（永久移动），但不允许更改从POST到GET的请求方法。
        
        永久重定向有两个： 301和308。
        
        两者都默认缓存，
        
        但是308不允许将请求方法从POST修改到GET, 301允许。
        
        临时重定向三个：302，303，307303强制浏览器可以将请求方法从POST修改到GET307不允许浏览器修改请求方法。302一开始的标准是不允许修改POST方法，但是浏览器的实现不遵循标准，标准就向现实妥协而做了修改。
        
5. 发号器的设计和实现
    1. 为什么要使⽤发号器的⽅案。
    2. 常⻅的发号器实现⽅式有哪些。
    3. 更进⼀步：如何实现⾼可⽤的发号器（MySQL主备+分⽚）
6. 如何降低查看链接耗时？
    1. 加Redis缓存保存 短链接->⻓链接
    2. 再进⼀步：添加本地缓存构成多级缓存
7. 如何解决缓存击穿问题？
    1. singleflight 合并请求
    2. singleflight的实现原理
8. 如何解决缓存穿透问题？
    1. 使⽤布隆过滤器过滤掉不存在的短链请求。
    2. 布隆过滤器的实现原理
    3. 布隆过滤器的优点和缺点是什么？
    4. 怎么⽀持删除短链接？
        1. 使⽤布⾕⻦过滤器
    5. 布隆过滤器和布⾕⻦过滤器的区别
        1. 算法：布隆过滤器多个hash函数。布⾕⻦过滤器⽤布⾕⻦哈希算法。布隆过滤器的多个
        哈希函数之间没关系。布⾕⻦过滤器的两个哈希函数可互相推导，两者有关系，⽤到了
        异或操作。
        2. 能否删除：布隆过滤器⽆法删除元素。布⾕⻦过滤器可以删除元素，有误删可能。
        3. 空间是否2的指数：布隆过滤器不需要2的指数。布⾕⻦过滤器必须是2的指数。
        4. 空间利⽤率：相同误判下，布⾕⻦空间节省40%多。
        5. 查询性能：布隆过滤器查询性能弱，原因是使⽤了多个hash函数，内存跨度⼤，缓存⾏
        命中率低。布⾕⻦过滤器访问内存次数低，效率相对⾼。
        6. 重复插⼊相同元素：布隆过滤器天然⾃带重复过滤。布⾕⻦过滤器会发⽣挤兑循环问
        题。
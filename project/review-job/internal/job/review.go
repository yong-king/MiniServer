package job

import (
	"context"
	"encoding/json"
	"errors"
	"review-job/internal/conf"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/segmentio/kafka-go"
)

// 1. 从kafka获取信息
// 2. 往es中添加信息

// type Server interface {
// 	Start(context.Context) error
// 	Stop(context.Context) error
// }

// func (s *Server) Start(ctx context.Context) error

type JobWork struct {
	kafkaReader *kafka.Reader
	esClient *EsClient
	log *log.Helper
}

type EsClient struct{
	*elasticsearch.TypedClient
	index string
}

func NewJobWrok(kafkaReader *kafka.Reader, esClient *EsClient, logger log.Logger) *JobWork{
	return &JobWork{
		kafkaReader: kafkaReader,
		esClient: esClient,
		log: log.NewHelper(logger),
	}
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

func NewKafkaReader(conf *conf.Kafka) *kafka.Reader {
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers:   conf.Brokers,
		Topic:     conf.Topic,
		GroupID: conf.GroupId,
	})
}

type Msg struct{
	Type     string `json:"type"`
	Database string `json:"databse"`
	Table    string `json:"table"`
	IsDdl    bool   `json:"isDdl"`
	Data     []map[string]interface{}
}

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

func (jw JobWork) Stop(ctx context.Context) error {
	jw.log.Debug("JobWorker stop....")
	// 程序退出前关闭Reader
	return jw.kafkaReader.Close()
}


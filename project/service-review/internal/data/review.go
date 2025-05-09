package data

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"service-review/internal/biz"
	"service-review/internal/data/model"
	"service-review/internal/data/query"
	"strconv"
	"strings"
	"time"

	"golang.org/x/sync/singleflight"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type reviewRepo struct {
	data *Data
	log  *log.Helper
}

// NewReviewRepo .
func NewReviewRepo(data *Data, logger log.Logger) biz.ReviewRepo {
	return &reviewRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

// SaveReview 创建评论
func (r *reviewRepo) SaveReview(ctx context.Context, review *model.ReviewInfo) (*model.ReviewInfo, error) {
	err := r.data.query.ReviewInfo.WithContext(ctx).Create(review)
	return review, err
}

// GetReviewByOrderID 根据订单id获取评论信息
func (r *reviewRepo) GetReviewByOrderID(ctx context.Context, orderId int64) ([]*model.ReviewInfo, error) {
	return r.data.query.ReviewInfo.WithContext(ctx).Where(r.data.query.ReviewInfo.OrderID.Eq(orderId)).Find()
}

// GetReview 根据reviewId获取评价详情
func (r *reviewRepo) GetReview(ctx context.Context, reviewId int64) (*model.ReviewInfo, error) {
	return r.data.query.ReviewInfo.WithContext(ctx).Where(r.data.query.ReviewInfo.ReviewID.Eq(reviewId)).First()
}

// AuditRrview 审查评价
func (r *reviewRepo) AuditRrview(ctx context.Context, param *biz.AuditParma) error {
	_, err := r.data.query.ReviewInfo.WithContext(ctx).Where(r.data.query.ReviewInfo.ReviewID.Eq(param.ReviewID)).Updates(
		map[string]interface{}{
			"status":    param.Status,
			"op_user":   param.OpUser,
			"op_reason": param.OpReason,
			"op_remarks":  param.OpMarks,
		},
	)
	return err
}

// IsAppealing 评价申请状态
func (r *reviewRepo) IsAppealing(ctx context.Context, param *biz.AppealParam) (*model.ReviewInfo, error) {
	appeal, err := r.data.query.ReviewInfo.WithContext(ctx).Where(r.data.query.ReviewInfo.ReviewID.Eq(param.ReviewID),
		r.data.query.ReviewInfo.StoreID.Eq(param.StoreID)).First()
	if err == gorm.ErrRecordNotFound{
		return &model.ReviewInfo{}, nil
	}
	if err != nil {
		return nil, err
	}
	return appeal, nil
}

// UpdateAppeal 更新申诉
func (r *reviewRepo) UpdateAppeal(ctx context.Context, param *biz.AppealParam) error {
	_, err := r.data.query.ReviewAppealInfo.WithContext(ctx).Where(r.data.query.ReviewAppealInfo.ReviewID.Eq(param.ReviewID),
		r.data.query.ReviewAppealInfo.StoreID.Eq(param.StoreID)).Updates(map[string]interface{}{
		"status":    10,
		"review_id": param.ReviewID,
		"store_id":  param.StoreID,
		"content":   param.Content,
		"pic_info":   param.Picinfo,
		"video_Info": param.VideoInfo,
		"appeal_id": param.AppealID,
	})
	if err != nil {
		return err
	}
	return nil
}

// CreateAppealReview 创建申诉
// 这个方法，有就更新，没有就创建
func (r *reviewRepo) CreateAppealReview(ctx context.Context, appeal *model.ReviewAppealInfo) error {
	return r.data.query.ReviewAppealInfo.
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
}

func (r *reviewRepo) SaveReply(ctx context.Context, reply *model.ReviewReplyInfo) (*model.ReviewReplyInfo, error) {
	// 1. 参数校验
	// 1.1 已经回复过的就不能在回复了
	// 1.2 水平越权 （A商家只能回复自己的不能回复B商家的）

	review, err := r.data.query.WithContext(ctx).ReviewInfo.Where(r.data.query.ReviewInfo.ReviewID.Eq(reply.ReviewID)).First()
	if err != nil {
		return nil, err
	}
	if review.HasReply == 1 {
		return nil, errors.New("该评价已被回复")
	}
	if review.StoreID != reply.StoreID {
		return nil, errors.New("水平越权")
	}
	r.log.Debugf("--------->更新数据!")
	// 2. 更新数据库中的数据（评价回复表和评价表要同时更新，涉及到事务操作）
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
	return reply, err
}

func (r *reviewRepo) ListReviewByUserID(ctx context.Context, userID int64, offset int, limit int) ([]*model.ReviewInfo, error) {
	return r.data.query.ReviewInfo.WithContext(ctx).Where(r.data.query.ReviewInfo.UserID.Eq(userID)).
		Order(r.data.query.ReviewInfo.ID.Desc()).Limit(limit).Offset(offset).Find()
}


func (r *reviewRepo) AuditAppeal(ctx context.Context, param *biz.AuditAppealParam) error{
	err := r.data.query.Transaction(func(tx *query.Query) error {
		// 申诉表
		if _, err := tx.ReviewAppealInfo.WithContext(ctx).
		Where(r.data.query.ReviewAppealInfo.AppealID.Eq(param.AppealID)).Updates(map[string]interface{}{
			"status":  param.Status,
			"op_user": param.OpUser,
		}); err != nil {
			return err
		}
		// 评论表
		if param.Status == 20 { // 申诉通过则需要隐藏评价
			if _, err := tx.ReviewInfo.WithContext(ctx).Where(tx.ReviewInfo.ReviewID.Eq(param.ReviewID)).Update(tx.ReviewInfo.Status, 40); err != nil{
				return err
			}
		}
		return nil
	})
	return err
}

func (r *reviewRepo) ListReviewByStoreID(ctx context.Context, storeID int64, offset int32, limit int32) ([]*biz.MyReviewinfo, error){
	// 去es里查询评价
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
}

var g singleflight.Group

func (r *reviewRepo) getData2(ctx context.Context, storeID int64, offset int32, limit int32) ([]*biz.MyReviewinfo, error){
	// 1. 先查询Redis缓存
	// 2. 缓存没有则查询ES
	// 3. 通过singleflight 合并短时间内大量的并发查询
	key := fmt.Sprintf("review:%d:%d:%d", storeID, offset, limit)
	b, err := r.getDataBySingleflight(ctx, key)
	if err != nil {
		return nil, err
	}

	hm := new(types.HitsMetadata)
	if err := json.Unmarshal(b, hm); err != nil{
		return nil, err
	}

	// 反序列化
	// 反序列化数据
	// resp.Hits.Hits[0].Source_(json.RawMessage)  ==>  model.ReviewInfo
	list := make([]*biz.MyReviewinfo, 0, hm.Total.Value)

	for _, hit := range hm.Hits{
		tmp := &biz.MyReviewinfo{}
		if err := json.Unmarshal(hit.Source_, tmp); err != nil{
			r.log.Errorf("json.Unmarshal(hit.Source_, tmp) failed, err:%v", err)
			continue
		}
		list = append(list, tmp)
	}
	return list, nil

}

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

// getDataFromCache 读取缓存数据
func  (r *reviewRepo) getDataFromCache(ctx context.Context, key string) ([]byte, error){
	r.log.Debugf("getDataFromCache key:%v\n", key)
	return r.data.rdb.Get(ctx, key).Bytes()
}

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

// setCache 设置缓存
func (r *reviewRepo) setCache(ctx context.Context, key string,  data []byte) error {
	return r.data.rdb.Set(ctx, key, data, time.Second*10).Err()
}
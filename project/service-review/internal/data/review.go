package data

import (
	"context"
	"service-review/internal/biz"
	"service-review/internal/data/model"

	"github.com/go-kratos/kratos/v2/log"
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
func (r *reviewRepo) SaveReview(ctx context.Context, review *model.ReviewInfo) (*model.ReviewInfo, error){
	err := r.data.query.ReviewInfo.WithContext(ctx).Create(review)
	return review, err
}

// GetReviewByOrderID 根据订单id获取评论信息
func (r *reviewRepo) GetReviewByOrderID(ctx context.Context, orderId int64) ([]*model.ReviewInfo, error){
	return r.data.query.ReviewInfo.WithContext(ctx).Where(r.data.query.ReviewInfo.OrderID.Eq(orderId)).Find()
}
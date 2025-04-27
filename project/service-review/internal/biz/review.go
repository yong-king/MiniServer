package biz

import (
	"context"
	"fmt"
	v1 "service-review/api/review/v1"
	"service-review/internal/data/model"
	"service-review/pkg/snowflake"

	"github.com/go-kratos/kratos/v2/log"
)

// ReviewRepo is a Review repo.
type ReviewRepo interface {
	SaveReview(context.Context, *model.ReviewInfo) (*model.ReviewInfo, error)
	GetReviewByOrderID(context.Context, int64) ([]*model.ReviewInfo, error)
}


// ReviewUsecase is a Review usecase.
type ReviewUsecase struct {
	repo ReviewRepo
	log  *log.Helper
}

// NewReviewUsecase new a Review usecase.
func NewReviewUsecase(repo ReviewRepo, logger log.Logger) *ReviewUsecase{
	return &ReviewUsecase{repo: repo, log: log.NewHelper(logger)}
}

// CreateReview 创建评论
func (uc *ReviewUsecase) CreateReview(ctx context.Context, review *model.ReviewInfo) (*model.ReviewInfo, error) {
	uc.log.WithContext(ctx).Infof("CreateGreeter: %v", review)
	// 1. 参数校验
	// 1.1 参数基础校验，应该在上一层，即service层完成，采用validata方法
	// 1.2 参数业务校验：已经评价过的orderID就不让在评价了
	reviews , err := uc.repo.GetReviewByOrderID(ctx, review.OrderID)
	if err != nil {
		return nil, v1.ErrorDbFailed("数据库查询失败")
	}
	if len(reviews) > 0 {
		// 已经评价过
		fmt.Printf("订单已评价, len(reviews):%d\n", len(reviews))
		return nil, v1.ErrorOrderReviewed("订单%d已被评价", review.OrderID)
	}
	// 2. 生成评价id
	// 2.1 雪花算法生成id
	// 2.2 公司自己生成id的算法
	review.ReviewID = snowflake.GenID()

	// 查询订单和商品信息快照
	// 实际业务场景下就需要查询订单服务和商家服务（比如说通过RPC调用订单服务和商家服务）
	// 拼装数据入库
	return uc.repo.SaveReview(ctx, review)
}
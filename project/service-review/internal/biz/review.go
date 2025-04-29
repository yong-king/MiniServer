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
	GetReview(context.Context, int64) (*model.ReviewInfo, error)
	AuditRrview(context.Context, *AuditParma) error
	IsAppealing(context.Context, *AppealParam) (*model.ReviewInfo, error)
	UpdateAppeal(context.Context, *AppealParam) error
	CreateAppealReview(context.Context, *model.ReviewAppealInfo) error
	SaveReply(context.Context, *model.ReviewReplyInfo) (*model.ReviewReplyInfo, error) 
	ListReviewByUserID(context.Context, int64, int, int) ([]*model.ReviewInfo, error)
	AuditAppeal(context.Context, *AuditAppealParam) error
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

// GetReview 根据reviewId获取评价详情
func (uc *ReviewUsecase) GetReview(ctx context.Context, reviewId int64) (*model.ReviewInfo, error){
	uc.log.WithContext(ctx).Infof("GetReview: %d", reviewId)
	// 1. 参数校验
	// 1.1 参数基础校验，应该在上一层，即service层完成，采用validata方法
	// 1.2 业务参数校验，无
	// 2. 查询数据库
	return uc.repo.GetReview(ctx, reviewId)
}

// AuditRrview 评价审查
func (uc *ReviewUsecase) AuditRrview(ctx context.Context, param *AuditParma) error{
	uc.log.WithContext(ctx).Infof("AuditRrview: %v", param)
	return uc.repo.AuditRrview(ctx, param)
}

func (uc *ReviewUsecase) AppealReview(ctx context.Context, param *AppealParam) (int64, error){
	uc.log.WithContext(ctx).Infof("AppealReview: %v", param)
	// 1. 参数校验
	// 1.1 基础参数校验，validate
	// 1.2 业务参数校验
	// 1.2.1 如果正在处于申诉状态，就不让在申诉 即 status > 10
	appeal, err := uc.repo.IsAppealing(ctx, param)
	if err != nil {
		return -1, v1.ErrorDbFailed("数据库查询失败")
	}
	if  appeal.Status > 10 {
		// 正在申诉
		return -1, v1.ErrorAppealReview("评价%d正在申诉中,申诉号是：%d", appeal.ReviewID)
	}
	// 1.2.2 如果有申诉过，就更新申诉
	// 2. 生成id
	param.AppealID = snowflake.GenID()
	apeal := &model.ReviewAppealInfo{
		ReviewID: param.ReviewID,
		StoreID: param.StoreID,
		AppealID: param.ReviewID,
		Reason: param.Reason,
		Content: param.Content,
		PicInfo: param.Picinfo,
		VideoInfo: param.VideoInfo,
	}
	// if appeal != nil {
	// 	err := uc.repo.UpdateAppeal(ctx, param)
	// 	if err != nil{
	// 		return -1, err
	// 	}
	// 	return param.AppealID, nil
	// }
	// 1.2.3 如果没有申诉过，就创建参数
	err = uc.repo.CreateAppealReview(ctx, apeal)
	if err != nil {
		return -1, err
	}
	return param.AppealID, nil
}

func (uc *ReviewUsecase) CreateReply(ctx context.Context, param *ReplyParam) (*model.ReviewReplyInfo, error){
	// 1. 参数校验
	// 1.1 已经回复过的就不能在回复了
	// 1.2 水平越权 （A商家只能回复自己的不能回复B商家的）
	// 2. 更新数据库中的数据（评价回复表和评价表要同时更新，涉及到事务操作）
	uc.log.WithContext(ctx).Infof("CreateReply: %v", param)
	reply := &model.ReviewReplyInfo{
		ReplyID:   snowflake.GenID(),
		ReviewID:  param.ReviewID,
		StoreID:   param.StoreID,
		Content:   param.Content,
		PicInfo:   param.Picinfo,
		VideoInfo: param.VideoInfo,
	}
	uc.log.Debugf("------->[biz CreateReply] reply:%v\n", reply)
	return uc.repo.SaveReply(ctx, reply)
}


func (uc *ReviewUsecase) ListReviewByUserId (ctx context.Context, userID int64, page, size int) ([]*model.ReviewInfo, error) {
	// 参数校验
	uc.log.WithContext(ctx).Infof("ListReviewByUserId, userid:%d, page:%d, size:%d", userID, page, size)
	if page <= 0{
		page = 1
	}
	if size <= 0 || size > 50{
		size = 10
	}
	offset := (page - 1) * size
	limit := size
	return uc.repo.ListReviewByUserID(ctx, userID, offset, limit)
}

func (uc *ReviewUsecase) AuditAppeal (ctx context.Context, param *AuditAppealParam) error {
	uc.log.WithContext(ctx).Infof("AuditAppeal: %v", param)
	return uc.repo.AuditAppeal(ctx, param)
}
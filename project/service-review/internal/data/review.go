package data

import (
	"context"
	"errors"
	"service-review/internal/biz"
	"service-review/internal/data/model"
	"service-review/internal/data/query"

	"github.com/go-kratos/kratos/v2/log"
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
	if review.StoreID == reply.StoreID {
		return nil, errors.New("水平越权")
	}

	// 2. 更新数据库中的数据（评价回复表和评价表要同时更新，涉及到事务操作）
	r.data.query.Transaction(func(tx *query.Query) error {
		// 回复一条插入数据
		if err := tx.ReviewReplyInfo.WithContext(ctx).Save(reply); err != nil {
			r.log.WithContext(ctx).Errorf("SaveReply create reply fail, err:%v", err)
			return err
		}
		// 评价表更新hasReply字段\
		if _, err := tx.ReviewInfo.WithContext(ctx).Where(tx.ReviewInfo.ReviewID.Eq(reply.ReplyID)).Update(tx.ReviewInfo.HasReply, 1); err != nil {
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
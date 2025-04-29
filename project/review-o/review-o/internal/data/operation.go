package data

import (
	"context"
	v1 "review-o/api/review/v1"
	"review-o/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
)

type operationRepo struct {
	data *Data
	log  *log.Helper
}

// NewGreeterRepo .
func NewOperationRepo(data *Data, logger log.Logger) biz.OperationRepo {
	return &operationRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (r *operationRepo) AuditReview(ctx context.Context, param *biz.AuditReviewParam) error{
	ret, err := r.data.rc.AuditRrview(ctx, &v1.AuditRrviewRequest{
		ReviewID: param.ReviewID,
		Status: param.Status,
		OpUser: param.OpUser,
		OpReason: param.OpReason,
		OpMarks: param.OpRemarks,
	})
	r.log.Debugf("[data AuditReview] ret:%v\n, err:%v\n", ret, err)
	return err
}

func (r *operationRepo) AuditAppeal(ctx context.Context, param *biz.AuditAppealParam) error{
	ret, err := r.data.rc.AuditAppeal(ctx, &v1.AuditAppealRequest{
		AppealID: param.AppealID,
		ReviewID: param.ReviewID,
		Status: param.Status,
		OpUser: param.OpUser,
		OpRemarks: param.OpRemarks,
	})
	r.log.Debugf("[data AuditAppeal] ret:%v\n, err:%v\n", ret, err)
	return err
}
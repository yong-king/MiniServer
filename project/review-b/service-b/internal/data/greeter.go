package data

import (
	"context"

	v1 "service-b/api/review/v1"
	"service-b/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
)

type businessRepo struct {
	data *Data
	log  *log.Helper
}

// NewGreeterRepo .
func NewBusinessRepo(data *Data, logger log.Logger) biz.BusinessRepo {
	return &businessRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (r *businessRepo) Reply(ctx context.Context, param *biz.ReviewReplyParam) (int64, error) {
	// rpc调用
	reply, err := r.data.rc.ReplyReview(ctx, &v1.ReplyReviewRequest{
		ReviewID: param.ReviewID,
		StoreID: param.StoreID,
		Content: param.Content,
		Picinfo: param.PicInfo,
		VideoInfo: param.VideoInfo,
	})
	if err != nil {
		return 0, err
	}
	return reply.GetReplyID(), nil
}

func  (r *businessRepo) AppealReview(ctx context.Context, param *biz.AppealReviewParam) (int64, error){
	reply, err := r.data.rc.AppealReview(ctx, &v1.AppealReviewRequest{
		ReviewID: param.ReviewID,
		StoreID: param.StoreID,
		Reason: param.Reason,
		Content: param.Content,
		Picinfo: param.PicInfo,
		VideoInfo: param.VideoInfo,
	})
	if err != nil {
		return 0, err
	}
	return reply.GetAppleID(), nil
}
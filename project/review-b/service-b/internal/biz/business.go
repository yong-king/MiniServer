package biz

import (
	"context"


	"github.com/go-kratos/kratos/v2/log"
)


// BusinessRepo is a Greater repo.
type BusinessRepo interface {
	Reply(context.Context, *ReviewReplyParam) (int64, error)
	AppealReview(context.Context, *AppealReviewParam) (int64, error)
}

// GreeterUsecase is a Greeter usecase.
type BusinessUsecase struct {
	repo BusinessRepo
	log  *log.Helper
}

// NewBusinessUsecase new a Greeter usecase.
func NewBusinessUsecase(repo BusinessRepo, logger log.Logger) *BusinessUsecase {
	return &BusinessUsecase{repo: repo, log: log.NewHelper(logger)}
}



func (uc *BusinessUsecase) CreateReply (ctx context.Context, param *ReviewReplyParam) (int64, error){
	uc.log.WithContext(ctx).Infof("[biz] CreateReply param:%v\n", param)
	return uc.repo.Reply(ctx, param)
}

func (uc *BusinessUsecase) AppealReview (ctx context.Context, param *AppealReviewParam) (int64, error){
	return uc.repo.AppealReview(ctx, param)
}
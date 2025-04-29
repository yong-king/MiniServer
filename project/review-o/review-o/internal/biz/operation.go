package biz

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
)


// OperationRepo is a Greater repo.
type OperationRepo interface {
	AuditReview(context.Context, *AuditReviewParam) error
	AuditAppeal(context.Context, *AuditAppealParam) error
}

// OperationUsecase is a Greeter usecase.
type OperationUsecase struct {
	repo OperationRepo
	log  *log.Helper
}

// NewGreeterUsecase new a Greeter usecase.
func NewOperationUsecase(repo OperationRepo, logger log.Logger) *OperationUsecase {
	return &OperationUsecase{repo: repo, log: log.NewHelper(logger)}
}

// CreateGreeter creates a Greeter, and returns the new Greeter.
func (uc *OperationUsecase) CreateOperation(ctx context.Context) (error) {
	uc.log.WithContext(ctx).Infof("CreateOperation:")
	return nil
}

func (uc *OperationUsecase) AuditReview (ctx context.Context, param *AuditReviewParam) error{
	return uc.repo.AuditReview(ctx, param)
}	

func (uc *OperationUsecase)  AuditAppeal(ctx context.Context, param *AuditAppealParam) error {
	return uc.repo.AuditAppeal(ctx, param)
}
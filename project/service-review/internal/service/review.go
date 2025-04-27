package service

import (
	"context"

	pb "service-review/api/review/v1"
	"service-review/internal/biz"
	"service-review/internal/data/model"
)

type ReviewService struct {
	pb.UnimplementedReviewServer

	uc *biz.ReviewUsecase
}

func NewReviewService(uc *biz.ReviewUsecase) *ReviewService {
	return &ReviewService{
		uc: uc,
	}
}

func (s *ReviewService) CreateReview(ctx context.Context, req *pb.CreateReviewRequest) (*pb.CreateReviewReply, error) {
	// 调用biz层
	var anonymous int32
	if req.Annoymous{
		anonymous = 1
	}
	review, err := s.uc.CreateReview(ctx, &model.ReviewInfo{
		UserID: req.UserId,
		OrderID: req.OrderId,
		Score: req.Score,
		ServiceScore: req.ServiceScore,
		ExpressScore: req.ExpreeScore,
		Content: req.Content,
		PicInfo: req.Picinfo,
		VideoInfo: req.VideoInfo,
		Anonymous: anonymous,
		Status: 0,
	})
	if err != nil {
		return nil, err
	}
	return &pb.CreateReviewReply{ReviewId: review.ReviewID}, nil
}
func (s *ReviewService) UpdateReview(ctx context.Context, req *pb.UpdateReviewRequest) (*pb.UpdateReviewReply, error) {
	return &pb.UpdateReviewReply{}, nil
}
func (s *ReviewService) DeleteReview(ctx context.Context, req *pb.DeleteReviewRequest) (*pb.DeleteReviewReply, error) {
	return &pb.DeleteReviewReply{}, nil
}
func (s *ReviewService) GetReview(ctx context.Context, req *pb.GetReviewRequest) (*pb.GetReviewReply, error) {
	return &pb.GetReviewReply{}, nil
}
func (s *ReviewService) ListReview(ctx context.Context, req *pb.ListReviewRequest) (*pb.ListReviewReply, error) {
	return &pb.ListReviewReply{}, nil
}

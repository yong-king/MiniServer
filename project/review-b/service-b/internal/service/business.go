package service

import (
	"context"
	"fmt"

	pb "service-b/api/business/v1"
	"service-b/internal/biz"
)

type BusinessService struct {
	pb.UnimplementedBusinessServer

	uc *biz.BusinessUsecase
}

func NewBusinessService(uc *biz.BusinessUsecase) *BusinessService {
	return &BusinessService{uc: uc}
}

func (s *BusinessService) ReplyReview(ctx context.Context, req *pb.ReplyReviewRequest) (*pb.ReplyReviewReply, error) {
	fmt.Printf("[service] ReplyReview req:%v\n", req)
	param := &biz.ReviewReplyParam{
		ReviewID: req.GetReviewID(),
		StoreID: req.GetStoreID(),
		Content: req.GetContent(),
		PicInfo: req.GetPicInfo(),
		VideoInfo: req.GetVideoInfo(),
	}
	replyId , err := s.uc.CreateReply(ctx, param)
	if err != nil {
		return nil, err
	}
	return &pb.ReplyReviewReply{ReplyID: replyId}, nil
}

func (s *BusinessService) AppealReview(ctx context.Context, req *pb.AppealReviewRequest) (*pb.AppealReviewReply, error){
	param := &biz.AppealReviewParam{
		ReviewID: req.GetReviewID(),
		StoreID: req.GetStoreID(),
		Reason: req.GetReason(),
		Content: req.GetContent(),
		PicInfo: req.GetPicInfo(),
		VideoInfo: req.GetVideoInfo(),
	}
	appealId, err := s.uc.AppealReview(ctx, param)
	if err != nil{
		return nil, err
	}
	return &pb.AppealReviewReply{AppealID: appealId}, nil
}
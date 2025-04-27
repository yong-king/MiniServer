package service

import (
	"context"
	"fmt"

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

// CreateReview 创建评论
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

// GetReview 根据评价id获取评价详情
func (s *ReviewService) GetReview(ctx context.Context, req *pb.GetReviewRequest) (*pb.GetReviewReply, error) {
	fmt.Printf("[service GetReview req:%#v]\n", req)
	// 调用biz层
	review, err := s.uc.GetReview(ctx, req.ReviewID)
	if err != nil {
		return nil, err
	}
	return &pb.GetReviewReply{
		Data: &pb.ReviewInfo{
			ReviewID: review.ReviewID,
			UserID: review.UserID,
			OrderID: review.OrderID,
			Score: review.Score,
			ExpressScore: review.ExpressScore,
			ServiceScore: review.ServiceScore,
			Content: review.Content,
			PicInfo: review.PicInfo,
			VideoInfo: review.VideoInfo,
			Status: review.Status,
		},
	}, nil
}

// AuditRrview o端申诉评价
func (s *ReviewService) AuditRrview(ctx context.Context, req *pb.AuditRrviewRequest) (*pb.AuditRrviewReply, error){
	fmt.Printf("[serivce AuditRrview req:%#v]\n", req)
	// 调用biz层
	err := s.uc.AuditRrview(ctx, &biz.AuditParma{
		ReviewID: req.ReviewID,
		Status: req.Status,
		OpUser: req.OpUser,
		OpReason: req.OpReason,
		OpMarks: req.GetOpMarks(),
	})
	if err != nil {
		return nil, err
	}
	return &pb.AuditRrviewReply{ReviewID: req.ReviewID, Status: req.Status}, nil
}

// AppealReview 评价申述
func (s *ReviewService) AppealReview(ctx context.Context, req *pb.AppealReviewRequest) (*pb.AppealReviewReply, error){
	fmt.Printf("[serivce AuditRrview req:%#v]\n", req)
	// 调用biz层
	appealID, err := s.uc.AppealReview(ctx, &biz.AppealParam{
		ReviewID: req.ReviewID,
		StoreID: req.StoreID,
		Reason: req.Reason,
		Content: req.Content,
		Picinfo: req.Picinfo,
		VideoInfo: req.VideoInfo,
	})
	if err != nil {
		return nil, err
	}
	return &pb.AppealReviewReply{AppleID: appealID}, nil
}

// ReplyReview 回复评价
func (s *ReviewService) ReplyReview(ctx context.Context, req *pb.ReplyReviewRequest) (*pb.ReplyReviewReply, error) {
	fmt.Printf("[serivce ReplyReview req:%#v]\n", req)
	// 调用biz层
	reply, err := s.uc.CreateReply(ctx, &biz.ReplyParam{
		ReviewID: req.ReviewID,
		StoreID: req.StoreID,
		Content: req.Content,
		Picinfo: req.Picinfo,
		VideoInfo: req.VideoInfo,
	})
	if err != nil{
		return nil, err
	}
	return &pb.ReplyReviewReply{ReplyID: reply.ReplyID}, nil
}

// ListReviewByUserId 根据用户id返回评价列表
func (s *ReviewService) ListReviewByUserId(ctx context.Context, req *pb.ListReviewByUserIdRequest) (*pb.ListReviewByUserIdReply, error) {
	fmt.Printf("[serivce ListReviewByUserId req:%#v]\n", req)
	// 调用biz层
	dataList, err := s.uc.ListReviewByUserId(ctx, req.GetUserID(), int(req.GetPage()) ,int(req.GetSize()))
	if err != nil {
		return nil, err
	}
	list := make([]*pb.ReviewInfo, 0, len(dataList))
	for _, review := range dataList{
		list = append(list, &pb.ReviewInfo{
			ReviewID:     review.ReviewID,
			UserID:       review.UserID,
			OrderID:      review.OrderID,
			Score:        review.Score,
			ServiceScore: review.ServiceScore,
			ExpressScore: review.ExpressScore,
			Content:      review.Content,
			PicInfo:      review.PicInfo,
			VideoInfo:    review.VideoInfo,
			Status:       review.Status,
		})
	}
	return &pb.ListReviewByUserIdReply{List: list}, nil
}

// AuditAppeal o端评价申诉审核
func (s *ReviewService) AuditAppeal(ctx context.Context, req *pb.AuditAppealRequest) (*pb.AuditAppealReply, error){
	fmt.Printf("[serivce AuditAppeal req:%#v]\n", req)
	err := s.uc.AuditAppeal(ctx, &biz.AuditAppealParam{
		AppealID: req.AppealID,
		ReviewID: req.ReviewID,
		Status: req.Status,
		OpUser: req.OpUser,
	})
	if err != nil{
		return nil, err
	}
	return &pb.AuditAppealReply{}, nil
}
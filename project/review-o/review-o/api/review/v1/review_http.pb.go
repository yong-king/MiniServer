// Code generated by protoc-gen-go-http. DO NOT EDIT.
// versions:
// - protoc-gen-go-http v2.8.4
// - protoc             v5.29.3
// source: review/v1/review.proto

package v1

import (
	context "context"
	http "github.com/go-kratos/kratos/v2/transport/http"
	binding "github.com/go-kratos/kratos/v2/transport/http/binding"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the kratos package it is being compiled against.
var _ = new(context.Context)
var _ = binding.EncodeURL

const _ = http.SupportPackageIsVersion1

const OperationReviewAppealReview = "/api.review.v1.Review/AppealReview"
const OperationReviewAuditAppeal = "/api.review.v1.Review/AuditAppeal"
const OperationReviewAuditRrview = "/api.review.v1.Review/AuditRrview"
const OperationReviewCreateReview = "/api.review.v1.Review/CreateReview"
const OperationReviewGetReview = "/api.review.v1.Review/GetReview"
const OperationReviewListReviewByUserId = "/api.review.v1.Review/ListReviewByUserId"
const OperationReviewReplyReview = "/api.review.v1.Review/ReplyReview"

type ReviewHTTPServer interface {
	// AppealReview b端申述评价
	AppealReview(context.Context, *AppealReviewRequest) (*AppealReviewReply, error)
	// AuditAppeal o 端评价申诉审核
	AuditAppeal(context.Context, *AuditAppealRequest) (*AuditAppealReply, error)
	// AuditRrview o端审核评价
	AuditRrview(context.Context, *AuditRrviewRequest) (*AuditRrviewReply, error)
	// CreateReview c端创建评价
	CreateReview(context.Context, *CreateReviewRequest) (*CreateReviewReply, error)
	// GetReview c端获取评价详情
	GetReview(context.Context, *GetReviewRequest) (*GetReviewReply, error)
	// ListReviewByUserId C端查看userID下所有评价
	ListReviewByUserId(context.Context, *ListReviewByUserIdRequest) (*ListReviewByUserIdReply, error)
	// ReplyReview b端回复评价
	ReplyReview(context.Context, *ReplyReviewRequest) (*ReplyReviewReply, error)
}

func RegisterReviewHTTPServer(s *http.Server, srv ReviewHTTPServer) {
	r := s.Route("/")
	r.POST("/v1/review", _Review_CreateReview0_HTTP_Handler(srv))
	r.GET("/v1/review/{reviewID}", _Review_GetReview0_HTTP_Handler(srv))
	r.POST("v1/review/audit", _Review_AuditRrview0_HTTP_Handler(srv))
	r.POST("/v1/review/appeal", _Review_AppealReview0_HTTP_Handler(srv))
	r.POST("/v1/review/reply", _Review_ReplyReview1_HTTP_Handler(srv))
	r.POST("/v1/appeal/audit", _Review_AuditAppeal1_HTTP_Handler(srv))
	r.GET("/v1/{userID}/reviews", _Review_ListReviewByUserId0_HTTP_Handler(srv))
}

func _Review_CreateReview0_HTTP_Handler(srv ReviewHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in CreateReviewRequest
		if err := ctx.Bind(&in); err != nil {
			return err
		}
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationReviewCreateReview)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.CreateReview(ctx, req.(*CreateReviewRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*CreateReviewReply)
		return ctx.Result(200, reply)
	}
}

func _Review_GetReview0_HTTP_Handler(srv ReviewHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in GetReviewRequest
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		if err := ctx.BindVars(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationReviewGetReview)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.GetReview(ctx, req.(*GetReviewRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*GetReviewReply)
		return ctx.Result(200, reply)
	}
}

func _Review_AuditRrview0_HTTP_Handler(srv ReviewHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in AuditRrviewRequest
		if err := ctx.Bind(&in); err != nil {
			return err
		}
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationReviewAuditRrview)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.AuditRrview(ctx, req.(*AuditRrviewRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*AuditRrviewReply)
		return ctx.Result(200, reply)
	}
}

func _Review_AppealReview0_HTTP_Handler(srv ReviewHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in AppealReviewRequest
		if err := ctx.Bind(&in); err != nil {
			return err
		}
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationReviewAppealReview)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.AppealReview(ctx, req.(*AppealReviewRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*AppealReviewReply)
		return ctx.Result(200, reply)
	}
}

func _Review_ReplyReview1_HTTP_Handler(srv ReviewHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in ReplyReviewRequest
		if err := ctx.Bind(&in); err != nil {
			return err
		}
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationReviewReplyReview)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.ReplyReview(ctx, req.(*ReplyReviewRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*ReplyReviewReply)
		return ctx.Result(200, reply)
	}
}

func _Review_AuditAppeal1_HTTP_Handler(srv ReviewHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in AuditAppealRequest
		if err := ctx.Bind(&in); err != nil {
			return err
		}
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationReviewAuditAppeal)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.AuditAppeal(ctx, req.(*AuditAppealRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*AuditAppealReply)
		return ctx.Result(200, reply)
	}
}

func _Review_ListReviewByUserId0_HTTP_Handler(srv ReviewHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in ListReviewByUserIdRequest
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		if err := ctx.BindVars(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationReviewListReviewByUserId)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.ListReviewByUserId(ctx, req.(*ListReviewByUserIdRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*ListReviewByUserIdReply)
		return ctx.Result(200, reply)
	}
}

type ReviewHTTPClient interface {
	AppealReview(ctx context.Context, req *AppealReviewRequest, opts ...http.CallOption) (rsp *AppealReviewReply, err error)
	AuditAppeal(ctx context.Context, req *AuditAppealRequest, opts ...http.CallOption) (rsp *AuditAppealReply, err error)
	AuditRrview(ctx context.Context, req *AuditRrviewRequest, opts ...http.CallOption) (rsp *AuditRrviewReply, err error)
	CreateReview(ctx context.Context, req *CreateReviewRequest, opts ...http.CallOption) (rsp *CreateReviewReply, err error)
	GetReview(ctx context.Context, req *GetReviewRequest, opts ...http.CallOption) (rsp *GetReviewReply, err error)
	ListReviewByUserId(ctx context.Context, req *ListReviewByUserIdRequest, opts ...http.CallOption) (rsp *ListReviewByUserIdReply, err error)
	ReplyReview(ctx context.Context, req *ReplyReviewRequest, opts ...http.CallOption) (rsp *ReplyReviewReply, err error)
}

type ReviewHTTPClientImpl struct {
	cc *http.Client
}

func NewReviewHTTPClient(client *http.Client) ReviewHTTPClient {
	return &ReviewHTTPClientImpl{client}
}

func (c *ReviewHTTPClientImpl) AppealReview(ctx context.Context, in *AppealReviewRequest, opts ...http.CallOption) (*AppealReviewReply, error) {
	var out AppealReviewReply
	pattern := "/v1/review/appeal"
	path := binding.EncodeURL(pattern, in, false)
	opts = append(opts, http.Operation(OperationReviewAppealReview))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "POST", path, in, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *ReviewHTTPClientImpl) AuditAppeal(ctx context.Context, in *AuditAppealRequest, opts ...http.CallOption) (*AuditAppealReply, error) {
	var out AuditAppealReply
	pattern := "/v1/appeal/audit"
	path := binding.EncodeURL(pattern, in, false)
	opts = append(opts, http.Operation(OperationReviewAuditAppeal))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "POST", path, in, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *ReviewHTTPClientImpl) AuditRrview(ctx context.Context, in *AuditRrviewRequest, opts ...http.CallOption) (*AuditRrviewReply, error) {
	var out AuditRrviewReply
	pattern := "v1/review/audit"
	path := binding.EncodeURL(pattern, in, false)
	opts = append(opts, http.Operation(OperationReviewAuditRrview))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "POST", path, in, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *ReviewHTTPClientImpl) CreateReview(ctx context.Context, in *CreateReviewRequest, opts ...http.CallOption) (*CreateReviewReply, error) {
	var out CreateReviewReply
	pattern := "/v1/review"
	path := binding.EncodeURL(pattern, in, false)
	opts = append(opts, http.Operation(OperationReviewCreateReview))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "POST", path, in, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *ReviewHTTPClientImpl) GetReview(ctx context.Context, in *GetReviewRequest, opts ...http.CallOption) (*GetReviewReply, error) {
	var out GetReviewReply
	pattern := "/v1/review/{reviewID}"
	path := binding.EncodeURL(pattern, in, true)
	opts = append(opts, http.Operation(OperationReviewGetReview))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "GET", path, nil, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *ReviewHTTPClientImpl) ListReviewByUserId(ctx context.Context, in *ListReviewByUserIdRequest, opts ...http.CallOption) (*ListReviewByUserIdReply, error) {
	var out ListReviewByUserIdReply
	pattern := "/v1/{userID}/reviews"
	path := binding.EncodeURL(pattern, in, true)
	opts = append(opts, http.Operation(OperationReviewListReviewByUserId))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "GET", path, nil, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *ReviewHTTPClientImpl) ReplyReview(ctx context.Context, in *ReplyReviewRequest, opts ...http.CallOption) (*ReplyReviewReply, error) {
	var out ReplyReviewReply
	pattern := "/v1/review/reply"
	path := binding.EncodeURL(pattern, in, false)
	opts = append(opts, http.Operation(OperationReviewReplyReview))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "POST", path, in, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

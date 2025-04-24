package service

import (
	"context"
	"errors"

	pb "helloworld/api/bubble/v1"
	"helloworld/internal/biz"
)

type TodoService struct {
	pb.UnimplementedTodoServer

	uc *biz.TodoUsecase
}

func NewTodoService(uc *biz.TodoUsecase) *TodoService {
	return &TodoService{
		uc: uc,
	}
}

// CreateTodo 添加信息
func (s *TodoService) CreateTodo(ctx context.Context, req *pb.CreateTodoRequest) (*pb.CreateTodoReply, error) {
	// 检验参数
	if len(req.Title) == 0 {
		return &pb.CreateTodoReply{}, errors.New("need title")
	}
	// 调用业务逻辑
	td, err := s.uc.CreateTodo(ctx, &biz.Todo{Title: req.Title})
	if err != nil {
		return &pb.CreateTodoReply{}, err
	}
	// 返回
	return &pb.CreateTodoReply{Id: td.ID, Title: td.Title, Status: td.Status}, nil
}

// UpdateTodo 更新信息
func (s *TodoService) UpdateTodo(ctx context.Context, req *pb.UpdateTodoRequest) (*pb.UpdateTodoReply, error) {
	// 参数校验
	if req.Id <= 0 || req.Title == ""{
		return nil, errors.New("invalid param")
	}
	// 调用biz逻辑
	err := s.uc.Update(ctx, &biz.Todo{ID: req.Id, Title: req.Title, Status: req.Status})
	if err != nil {
		return &pb.UpdateTodoReply{}, err
	}
	// 返回
	return &pb.UpdateTodoReply{}, nil
}

// DeleteTodo 删除信息
func (s *TodoService) DeleteTodo(ctx context.Context, req *pb.DeleteTodoRequest) (*pb.DeleteTodoReply, error) {
	// 参数校验
	if req.Id <= 0 {
		return nil, errors.New("invalid id")
	}
	// 调用diz逻辑
	err := s.uc.Delete(ctx, req.Id)
	if err != nil {
		return &pb.DeleteTodoReply{}, err
	}
	return &pb.DeleteTodoReply{}, nil
}

// GetTodo 根据id获取信息
func (s *TodoService) GetTodo(ctx context.Context, req *pb.GetTodoRequest) (*pb.GetTodoReply, error) {
	// 参数校验
	if req.Id <= 0 {
		return nil, errors.New("id error")
	}
	// 调用biz业务逻辑
	ret, err := s.uc.Get(ctx, req.Id)
	if err != nil {
		return &pb.GetTodoReply{}, pb.ErrorTodoNotFound("id:%v todo is not found", req.Id)
	}
	// 返回
	return &pb.GetTodoReply{Todo: &pb.Todo{Id: ret.ID, Status: ret.Status, Title: ret.Title}}, nil
}

// ListTodo 获取清单列表
func (s *TodoService) ListTodo(ctx context.Context, req *pb.ListTodoRequest) (*pb.ListTodoReply, error) {
	// 调用diz逻辑
	dataList, err := s.uc.ListTodo(ctx)
	if err != nil {
		return nil, err
	}

	reply := &pb.ListTodoReply{}
	for _, data := range dataList{
		reply.Todo = append(reply.Todo, &pb.Todo{
			Id: data.ID,
			Title: data.Title,
			Status: data.Status,
		})
	}
	return reply, nil
}

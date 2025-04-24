package biz

import (
	"context"

	v1 "helloworld/api/helloworld/v1"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
)

var (
	// ErrUserNotFound is user not found.
	ErrUserNotFound = errors.NotFound(v1.ErrorReason_USER_NOT_FOUND.String(), "user not found")
)

// Todo is a Todo model.
type Todo struct {
	ID int64
	Title string
	Status bool
}

// TodoRepo is a Todo repo.
type TodoRepo interface {
	Save(context.Context, *Todo) (*Todo, error)
	Update(context.Context, *Todo) error
	Delete(context.Context, int64) error
	FindByID(context.Context, int64) (*Todo, error)
	ListAll(context.Context) ([]*Todo, error)
}

// TodoUsecase is a Todo usecase.
type TodoUsecase struct {
	repo TodoRepo
	log  *log.Helper
}

// NewTodoUsecase new a Todo usecase.
func NewTodoUsecase(repo TodoRepo, logger log.Logger) *TodoUsecase {
	return &TodoUsecase{repo: repo, log: log.NewHelper(logger)}
}

// CreateTodo creates a Todo, and returns the new Todo.
func (uc *TodoUsecase) CreateTodo (ctx context.Context, t *Todo) (*Todo, error) {
	uc.log.WithContext(ctx).Infof("CreateTodo: %v", t)
	return uc.repo.Save(ctx, t)
}

// Get 根据id获取信息
func (uc *TodoUsecase) Get(ctx context.Context, id int64) (*Todo, error) {
	uc.log.WithContext(ctx).Infof("GetTodo: %v", id)
	return uc.repo.FindByID(ctx, id)
}

// Update 更新信息
func (uc *TodoUsecase) Update(ctx context.Context, t *Todo) (error) {
	uc.log.WithContext(ctx).Infof("UpdateTodo: %v", t)
	return uc.repo.Update(ctx, t)
}

// Delete 删除指定信息
func (uc *TodoUsecase) Delete(ctx context.Context, id int64) (error) {
	uc.log.WithContext(ctx).Infof("DeleteTodo: %v", id)
	return uc.repo.Delete(ctx, id)
}

// 
func (uc *TodoUsecase) ListTodo (ctx context.Context) ([]*Todo, error) {
	return uc.repo.ListAll(ctx)
}
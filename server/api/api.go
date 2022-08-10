package api

import (
	"context"

	ent "github.com/gregwebs/go-grpc-openapi-ent/ent"
	apiv1 "github.com/gregwebs/go-grpc-openapi-ent/gen/proto/go/todo/v1"
	dal "github.com/gregwebs/go-grpc-openapi-ent/server/dal"
)

func NewtodoService(db *ent.Client) *todoService {
	return &todoService{db: db}
}

type todoService struct {
	db *ent.Client
}

func (s *todoService) GetUser(ctx context.Context, req *apiv1.GetUserRequest) (*apiv1.GetUserResponse, error) {
	emp, err := dal.UserByEmail(ctx, s.db, req.GetEmail())
	if err != nil {
		return nil, err
	}
	return &apiv1.GetUserResponse{
		User: &apiv1.User{
			Id:    uint32(emp.ID),
			Name:  emp.Name,
			Email: emp.Email,
		},
	}, nil
}

func (s *todoService) ListTaskUsers(ctx context.Context, req *apiv1.ListTaskUsersRequest) (*apiv1.ListTaskUsersResponse, error) {
	task, err := dal.TaskUsers(ctx, s.db, int(req.GetId()))
	if err != nil {
		return nil, err
	}

	users := make([]*apiv1.User, len(task.Edges.Users))
	for i, user := range task.Edges.Users {
		users[i] = &apiv1.User{
			Id:    uint32(user.ID),
			Name:  user.Name,
			Email: user.Email,
		}
	}
	return &apiv1.ListTaskUsersResponse{
		Task: &apiv1.Task{
			Id:          uint32(task.ID),
			Description: task.Description,
		},
		Users: users,
	}, nil
}

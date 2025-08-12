package clients

import (
	"context"
	"time"

	"google.golang.org/grpc"
	userpb "github.com/huntercenter1/backend-test/proto"
)

type UserClient interface {
	Validate(ctx context.Context, userID string) (bool, error)
}

type userClient struct{ cc userpb.UserServiceClient }

func NewUserClient(addr string) (UserClient, func() error, error) {
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil { return nil, nil, err }
	return &userClient{cc: userpb.NewUserServiceClient(conn)}, conn.Close, nil
}

func (c *userClient) Validate(ctx context.Context, userID string) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second); defer cancel()
	resp, err := c.cc.ValidateUser(ctx, &userpb.ValidateUserRequest{UserId: userID})
	if err != nil { return false, err }
	return resp.GetValid(), nil
}

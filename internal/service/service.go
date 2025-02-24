package service

import (
	"context"
	"twitter/internal/common"
	"twitter/internal/db"
	"twitter/internal/dto"
)

type PostService interface {
	CreatePost(context.Context, *dto.CreatePostRequest) (*dto.CreatePostResponse, error)
	GetPostsByAuthorIdsAndTimeWindow(context.Context, *dto.GetPostsByAuthorIdsAndTimeWindowRequest) (*dto.GetPostsByAuthorIdsAndTimeWindowResponse, error)
}
type AfterPostCreate interface {
	AfterPostCreateHook(context.Context, *dto.Post)
}

type UserService interface {
	CreateUser(context.Context, *dto.CreateUserRequest) (*dto.CreateUserResponse, error)
	FollowUser(context.Context, *dto.FollowUserRequest) (*dto.FollowUserResponse, error)
	GetFollowingUsers(context.Context, *dto.GetFollowingUsersRequest) (*dto.GetFollowingUsersResponse, error)

	GetFollowerUsers(context.Context, *dto.GetFollowerUsersRequest) (*dto.GetFollowerUsersResponse, error)
	GetAllUsers(context.Context, *dto.GetAllUsersRequest) (*dto.GetAllUsersResponse, error)
	GetUser(context.Context, *dto.GetUserRequest) (*dto.GetUserResponse, error)
	GetUsers(context.Context, []int) (*dto.GetAllUsersResponse, error)
}

type TimelineService interface {
	GetTimelineForUser(context.Context, *dto.GetTimelineForUserRequest) (*dto.GetTimelineForUserResponse, error)
}

type ServiceDependencies struct {
	UserService     UserService
	PostService     PostService
	TimelineService TimelineService
}

func NewServiceDependencies(db *db.Db, lg common.Logger) *ServiceDependencies {
	sdp := &ServiceDependencies{}
	sdp.UserService = NewUserService(sdp, db, lg)
	sdp.TimelineService = NewTimelineService(sdp, db, lg)
	sdp.PostService = NewPostService(sdp, db, lg)
	return sdp
}

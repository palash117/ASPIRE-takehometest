package service

import (
	"context"
	"twitter/internal/common"
	"twitter/internal/dao"
	"twitter/internal/db"
	"twitter/internal/dto"
	"twitter/internal/model"
)

func NewUserService(dependencies *ServiceDependencies, db *db.Db, lg common.Logger) UserService {
	userDao := dao.NewUserDao(db)
	dao.AutoMigrate(userDao.(dao.Migratable))
	return &userService{dependencies: dependencies, userDao: userDao, lg: lg}
}

type userService struct {
	dependencies *ServiceDependencies
	userDao      dao.UserDao
	lg           common.Logger
}

func (u *userService) CreateUser(ctx context.Context, req *dto.CreateUserRequest) (*dto.CreateUserResponse, error) {
	mdl := model.FromUserDto(&dto.User{
		Email:    req.Email,
		UserName: req.UserName,
		//Password: req.Pa
	})
	mdl, err := u.userDao.Create(ctx, mdl)
	if err != nil {
		return nil, err
	}
	return &dto.CreateUserResponse{User: *(mdl.ToDto())}, nil
}

func (u *userService) FollowUser(ctx context.Context, req *dto.FollowUserRequest) (*dto.FollowUserResponse, error) {
	follower, err := u.userDao.IsFollowing(ctx, int(req.FollowerUserId), int(req.FollowingUserId))
	if follower == nil {
		follower, err = u.userDao.Follow(ctx, int(req.FollowerUserId), int(req.FollowingUserId))
		if err != nil {
			return nil, err
		}
	}
	resp := &dto.FollowUserResponse{
		Follower: *(follower.ToDto()),
	}
	return resp, nil

}

func (u *userService) GetFollowingUsers(ctx context.Context, req *dto.GetFollowingUsersRequest) (*dto.GetFollowingUsersResponse, error) {
	followees, err := u.userDao.GetFollowees(ctx, (int)(req.PrimaryUserId))
	if err != nil {
		return nil, err
	}
	resp := &dto.GetFollowingUsersResponse{
		FollowingUsers: make([]*dto.Follower, len(followees)),
	}
	for i, flwr := range followees {
		resp.FollowingUsers[i] = flwr.ToDto()
	}
	return resp, nil
}

func (u *userService) GetFollowerUsers(ctx context.Context, req *dto.GetFollowerUsersRequest) (*dto.GetFollowerUsersResponse, error) {
	followers, err := u.userDao.GetFollowers(ctx, int(req.PrimaryUserId))
	if err != nil {
		return nil, err
	}
	resp := &dto.GetFollowerUsersResponse{
		FollowerUsers: make([]*dto.Follower, len(followers)),
	}
	for i, flwr := range followers {
		resp.FollowerUsers[i] = flwr.ToDto()
	}
	return resp, nil
}
func (u *userService) GetAllUsers(ctx context.Context, req *dto.GetAllUsersRequest) (*dto.GetAllUsersResponse, error) {
	users, err := u.userDao.GetAll(ctx)
	if err != nil {
		return nil, err

	}
	resp := &dto.GetAllUsersResponse{
		Users: make([]*dto.User, len(users)),
	}
	for i, usr := range users {
		resp.Users[i] = usr.ToDto()
	}
	return resp, err

}
func (u *userService) GetUser(ctx context.Context, req *dto.GetUserRequest) (*dto.GetUserResponse, error) {

	user, err := u.userDao.GetById(ctx, (int)(req.UserId))
	if err != nil {
		return nil, err

	}
	return &dto.GetUserResponse{
		User: *(user.ToDto()),
	}, nil

}
func (u *userService) GetUsers(ctx context.Context, ids []int) (*dto.GetAllUsersResponse, error) {

	users, err := u.userDao.GetByIds(ctx, ids)
	if err != nil {
		return nil, err

	}
	resp := &dto.GetAllUsersResponse{
		Users: make([]*dto.User, len(users)),
	}
	for i, usr := range users {
		resp.Users[i] = usr.ToDto()
	}
	return resp, err
}

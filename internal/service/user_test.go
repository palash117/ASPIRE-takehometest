package service

import (
	"context"
	"testing"
	"twitter/internal/dto"

	"github.com/stretchr/testify/assert"
)

func TestUsers(t *testing.T) {
	//t.Skip()
	sdp := SetupTestServiceDependencies(t)
	ctx := context.Background()
	t.Run(" create user & get users green", func(t *testing.T) {

		usr1, err := sdp.UserService.CreateUser(ctx, &dto.CreateUserRequest{
			Email:    "tst1@mail.com",
			UserName: "tst1",
		})
		assert.NoError(t, err)
		assert.NotNil(t, usr1)
		assert.NotEqual(t, 0, int(usr1.Id))
		usr2, err := sdp.UserService.CreateUser(ctx, &dto.CreateUserRequest{
			Email:    "tst2@mail.com",
			UserName: "tst2",
		})
		assert.NoError(t, err)
		assert.NotNil(t, usr2)
		assert.NotEqual(t, 0, int(usr2.Id))
		gusr, err := sdp.UserService.GetUser(ctx, &dto.GetUserRequest{
			UserId: usr1.Id,
		})
		assert.NoError(t, err)
		assert.NotNil(t, gusr)
		assert.Equal(t, usr1.Email, gusr.Email)

		usrs, err := sdp.UserService.GetAllUsers(ctx, &dto.GetAllUsersRequest{})
		assert.NoError(t, err)
		assert.NotNil(t, usrs)
		assert.NotEqual(t, 0, len(usrs.Users))

		usrs, err = sdp.UserService.GetUsers(ctx, []int{int(usr1.Id), int(usr2.Id)})
		assert.NoError(t, err)
		assert.NotNil(t, usrs)
		assert.Equal(t, 2, len(usrs.Users))

	})
	t.Run("create and follow usrs", func(t *testing.T) {

		usr1, err := sdp.UserService.CreateUser(ctx, &dto.CreateUserRequest{
			Email:    "tst1@mail.com",
			UserName: "tst1",
		})
		assert.NoError(t, err)
		assert.NotNil(t, usr1)
		assert.NotEqual(t, 0, int(usr1.Id))
		usr2, err := sdp.UserService.CreateUser(ctx, &dto.CreateUserRequest{
			Email:    "tst2@mail.com",
			UserName: "tst2",
		})
		assert.NoError(t, err)
		assert.NotNil(t, usr2)
		assert.NotEqual(t, 0, int(usr2.Id))

		flw, err := sdp.UserService.FollowUser(ctx, &dto.FollowUserRequest{
			FollowerUserId:  usr1.Id,
			FollowingUserId: usr2.Id,
		})
		assert.NoError(t, err)
		assert.NotNil(t, flw)

		flwrs, err := sdp.UserService.GetFollowerUsers(ctx, &dto.GetFollowerUsersRequest{
			PrimaryUserId: usr2.Id,
		})
		assert.NoError(t, err)
		assert.NotNil(t, flwrs)
		assert.Equal(t, 1, len(flwrs.FollowerUsers))
		assert.Equal(t, flw.Follower.FollowId, flwrs.FollowerUsers[0].FollowId)

		flwng, err := sdp.UserService.GetFollowingUsers(ctx, &dto.GetFollowingUsersRequest{
			PrimaryUserId: usr1.Id,
		})
		assert.NoError(t, err)
		assert.NotNil(t, flwng)
		assert.Equal(t, 1, len(flwng.FollowingUsers))
		assert.Equal(t, flw.Follower.FollowId, flwng.FollowingUsers[0].FollowId)

	})
}

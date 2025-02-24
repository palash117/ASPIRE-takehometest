package dao

import (
	"context"
	"testing"
	"time"
	"twitter/internal/db"
	"twitter/internal/model"

	"github.com/stretchr/testify/assert"
)

func TestUserDao(t *testing.T) {
	//t.Skip()
	_, usrDao, _ := SetupDaos(t)
	ctx := context.Background()
	t.Run("test Create green", func(t *testing.T) {
		rsp, err := usrDao.Create(ctx, &model.User{
			Email:         "ab@mail.com",
			UserName:      "asdf",
			FollowerCount: 1,
			CreatedAt:     time.Now(),
		})
		assert.NoError(t, err)
		assert.NotNil(t, rsp)
		assert.NotEqual(t, 0, rsp.Id)
		rsp, err = usrDao.Create(ctx, &model.User{
			Email:         "ab@mail.com",
			UserName:      "asdf",
			FollowerCount: 1,
			CreatedAt:     time.Now(),
		})
		assert.NoError(t, err)
		assert.NotNil(t, rsp)
		assert.NotEqual(t, 0, rsp.Id)

	})
	t.Run("test GetById green", func(t *testing.T) {
		rsp, err := usrDao.GetById(ctx, 1)
		assert.NoError(t, err)
		assert.NotNil(t, rsp)
		assert.Equal(t, 1, rsp.Id)

	})
	t.Run("test GetById red id missing in db", func(t *testing.T) {
		rsp, err := usrDao.GetById(ctx, 3)
		assert.Nil(t, rsp)
		assert.NotNil(t, err)
		assert.Equal(t, err.Error(), "User entity not found with id 3")

	})
	t.Run("test GetByIds green", func(t *testing.T) {
		rsp, err := usrDao.GetByIds(ctx, []int{1, 2})
		assert.NoError(t, err)
		assert.NotNil(t, rsp)
		assert.Equal(t, 2, len(rsp))

	})
	t.Run("test GetByIds red invalid id", func(t *testing.T) {
		rsp, err := usrDao.GetByIds(ctx, []int{1, 2, 5})
		assert.NoError(t, err)
		assert.NotNil(t, rsp)
		assert.Equal(t, 2, len(rsp))

	})
	t.Run("test GetAll green", func(t *testing.T) {
		rsp, err := usrDao.GetAll(ctx)
		assert.NoError(t, err)
		assert.NotNil(t, rsp)
		assert.Equal(t, 2, len(rsp))

	})
	t.Run("test Follow green", func(t *testing.T) {
		rsp, err := usrDao.Follow(ctx, 1, 2)
		assert.NoError(t, err)
		assert.NotNil(t, rsp)
		assert.Equal(t, 1, rsp.FollowerId)
		assert.Equal(t, 2, rsp.FolloweeId)

	})
	t.Run("test IsFollowing green", func(t *testing.T) {
		rsp, err := usrDao.IsFollowing(ctx, 1, 2)
		assert.NoError(t, err)
		assert.NotNil(t, rsp)
		assert.Equal(t, 1, rsp.FollowerId)
		assert.Equal(t, 2, rsp.FolloweeId)

	})
	t.Run("test IsFollowing red", func(t *testing.T) {
		rsp, err := usrDao.IsFollowing(ctx, 2, 1)
		assert.Equal(t, err.Error(), db.NoEntriesFound{}.Error())
		var nilRf *model.Follower
		assert.Equal(t, nilRf, rsp)

	})
	t.Run("test GetFollowees green", func(t *testing.T) {
		rsp, err := usrDao.GetFollowees(ctx, 1)
		assert.NoError(t, err)
		assert.NotNil(t, rsp)
		assert.Equal(t, 1, len(rsp))
		assert.Equal(t, 1, rsp[0].FollowerId)
		assert.Equal(t, 2, rsp[0].FolloweeId)

	})
	t.Run("test GetFollowees  with 0 followees", func(t *testing.T) {
		rsp, err := usrDao.GetFollowees(ctx, 2)
		assert.NoError(t, err)
		assert.NotNil(t, rsp)
		assert.Equal(t, 0, len(rsp))

	})
	t.Run("test GetFollowing green", func(t *testing.T) {
		rsp, err := usrDao.GetFollowers(ctx, 2)
		assert.NoError(t, err)
		assert.NotNil(t, rsp)
		assert.Equal(t, 1, len(rsp))
		assert.Equal(t, 1, rsp[0].FollowerId)
		assert.Equal(t, 2, rsp[0].FolloweeId)

	})
	t.Run("test GetFollowing red invalid id", func(t *testing.T) {
		rsp, err := usrDao.GetFollowers(ctx, 3)
		assert.NoError(t, err)
		assert.NotNil(t, rsp)
		assert.Equal(t, 0, len(rsp))

	})
	t.Run("test GetFollowing 0 following", func(t *testing.T) {
		rsp, err := usrDao.GetFollowers(ctx, 1)
		assert.NoError(t, err)
		assert.NotNil(t, rsp)
		assert.Equal(t, 0, len(rsp))

	})
}

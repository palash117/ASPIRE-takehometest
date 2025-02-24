package dao

import (
	"context"
	"testing"
	"time"
	"twitter/internal/db"
	"twitter/internal/model"

	"github.com/stretchr/testify/assert"
)

func SetupDaos(t *testing.T) (PostDao, UserDao, TimelineDao) {
	t.Helper()
	gdb, err := db.SetupTestDb(t)
	assert.NoError(t, err)
	ddb := &db.Db{
		DB: gdb,
	}
	pstDao := NewPostDao(ddb)
	pstDao.(Migratable).AutoMigrate(false)

	usrDao := NewUserDao(ddb)
	usrDao.(Migratable).AutoMigrate(false)

	tmlnDao := NewTimelineDao(ddb)
	tmlnDao.(Migratable).AutoMigrate(false)

	return pstDao, usrDao, tmlnDao
}

func TestPostDao(t *testing.T) {
	//t.Skip()
	postDao, _, _ := SetupDaos(t)
	authorId := 11
	t.Run("test create posts green", func(t *testing.T) {
		ctx := context.Background()
		mdl, err := postDao.Create(ctx, &model.Post{
			AuthorId: authorId,
			Contents: "asdf",
			Subject:  "asdf",
		})
		assert.NoError(t, err)
		assert.NotNil(t, mdl)
		assert.Equal(t, 1, mdl.Id)

	})
	t.Run("test getPostsByAuthorIdsAndTimeWindow green", func(t *testing.T) {
		ctx := context.Background()
		psts, err := postDao.GetPostsByAuthorIdsAndTimeWindow(ctx, []int{authorId}, time.Now().Add(-1*time.Minute), time.Now())

		assert.NoError(t, err)
		assert.NotNil(t, psts)
		assert.Equal(t, 1, len(psts))
		assert.Equal(t, authorId, psts[0].AuthorId)
	})
	t.Run("test getPostsByAuthorIdsAndTimeWindow outside timewindow", func(t *testing.T) {
		ctx := context.Background()
		psts, err := postDao.GetPostsByAuthorIdsAndTimeWindow(ctx, []int{authorId}, time.Now().Add(-2*time.Hour), time.Now().Add(-1*time.Hour))

		assert.NoError(t, err)
		assert.NotNil(t, psts)
		assert.Equal(t, 0, len(psts))
	})
}

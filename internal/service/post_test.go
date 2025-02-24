package service

import (
	"context"
	"testing"
	"time"
	"twitter/internal/common"
	"twitter/internal/db"
	"twitter/internal/dto"

	"github.com/stretchr/testify/assert"
)

func SetupTestServiceDependencies(t *testing.T) *ServiceDependencies {
	//t.Helper()
	gdb, err := db.SetupTestDb(t)
	if err != nil {
		panic("unable to setup test db")
	}
	db := &db.Db{
		DB: gdb,
	}
	lg := common.NewLogger()
	sdp := NewServiceDependencies(db, lg)
	return sdp
}

func TestPosts(t *testing.T) {
	sdp := SetupTestServiceDependencies(t)
	ctx := context.Background()
	t.Run("create post", func(t *testing.T) {
		pst, err := sdp.PostService.CreatePost(ctx, &dto.CreatePostRequest{
			AuthorId: dto.Id(10),
			Contents: "asdf",
			Subject:  "asdfasdf",
		})
		assert.NoError(t, err)
		assert.NotNil(t, pst)
		assert.Equal(t, 1, int(pst.Id))

	})

	t.Run(" test GetPostsByAuthorIdsAndTimeWindow", func(t *testing.T) {
		psts, err := sdp.PostService.GetPostsByAuthorIdsAndTimeWindow(ctx, &dto.GetPostsByAuthorIdsAndTimeWindowRequest{
			AuthorIds: []int{10},
			From:      time.Now().Add(-1 * time.Hour),
			Till:      time.Now(),
		})
		assert.NoError(t, err)
		assert.Equal(t, 1, len(psts.Posts))

	})

}

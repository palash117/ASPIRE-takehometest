package dao

import (
	"context"
	"testing"
	"time"
	"twitter/internal/model"

	"github.com/stretchr/testify/assert"
)

func TestTimeline(t *testing.T) {
	//t.Skip()
	_, _, tmlnDao := SetupDaos(t)
	ctx := context.Background()
	t.Run(" test Create green", func(t *testing.T) {
		rsp, err := tmlnDao.Create(ctx, &model.Timeline{
			UserId:    10,
			PostId:    11,
			AuthorId:  12,
			Subject:   "asdf",
			Contents:  "asdfgh",
			CreatedAt: time.Now(),
		})
		assert.NoError(t, err)
		assert.NotNil(t, rsp)
		assert.Equal(t, 10, rsp.UserId)
		assert.Equal(t, 11, rsp.PostId)
		assert.Equal(t, 12, rsp.AuthorId)
	})
	t.Run(" test CreateInBatches green", func(t *testing.T) {
		rsp, err := tmlnDao.CreateInBatches(ctx, []*model.Timeline{
			&model.Timeline{
				UserId:    10,
				PostId:    11,
				AuthorId:  12,
				Subject:   "asdf",
				Contents:  "asdfgh",
				CreatedAt: time.Now(),
			},
			&model.Timeline{
				UserId:    10,
				PostId:    11,
				AuthorId:  12,
				Subject:   "asdf",
				Contents:  "asdfgh",
				CreatedAt: time.Now(),
			},
			&model.Timeline{
				UserId:    10,
				PostId:    11,
				AuthorId:  13,
				Subject:   "asdf",
				Contents:  "asdfgh",
				CreatedAt: time.Now(),
			},
			&model.Timeline{
				UserId:    10,
				PostId:    11,
				AuthorId:  12,
				Subject:   "asdf",
				Contents:  "asdfgh",
				CreatedAt: time.Now(),
			},
			&model.Timeline{
				UserId:    10,
				PostId:    11,
				AuthorId:  12,
				Subject:   "asdf",
				Contents:  "asdfgh",
				CreatedAt: time.Now(),
			},
			&model.Timeline{
				UserId:    10,
				PostId:    11,
				AuthorId:  13,
				Subject:   "asdf",
				Contents:  "asdfgh",
				CreatedAt: time.Now(),
			},
			&model.Timeline{
				UserId:    10,
				PostId:    11,
				AuthorId:  12,
				Subject:   "asdf",
				Contents:  "asdfgh",
				CreatedAt: time.Now(),
			},
			&model.Timeline{
				UserId:    10,
				PostId:    11,
				AuthorId:  12,
				Subject:   "asdf",
				Contents:  "asdfgh",
				CreatedAt: time.Now(),
			},
			&model.Timeline{
				UserId:    10,
				PostId:    11,
				AuthorId:  13,
				Subject:   "asdf",
				Contents:  "asdfgh",
				CreatedAt: time.Now(),
			},
			&model.Timeline{
				UserId:    10,
				PostId:    11,
				AuthorId:  12,
				Subject:   "asdf",
				Contents:  "asdfgh",
				CreatedAt: time.Now(),
			},
			&model.Timeline{
				UserId:    10,
				PostId:    11,
				AuthorId:  12,
				Subject:   "asdf",
				Contents:  "asdfgh",
				CreatedAt: time.Now(),
			},
			&model.Timeline{
				UserId:    10,
				PostId:    11,
				AuthorId:  13,
				Subject:   "asdf",
				Contents:  "asdfgh",
				CreatedAt: time.Now(),
			},
		})
		assert.NoError(t, err)
		assert.NotNil(t, rsp)
		assert.Equal(t, 12, len(rsp))
	})

	t.Run(" test GetTimelineForUser green", func(t *testing.T) {
		rsp, err := tmlnDao.GetTimelineForUser(ctx, 10, 0, 10)
		assert.NoError(t, err)
		assert.NotNil(t, rsp)
		assert.Equal(t, 10, len(rsp), "checking timeline length")
		rsp, err = tmlnDao.GetTimelineForUser(ctx, 10, 1, 10)
		assert.NoError(t, err)
		assert.NotNil(t, rsp)
		assert.Equal(t, 3, len(rsp), "checking timeline length")

	})

}

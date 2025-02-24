package service

import (
	"context"
	"testing"
	"time"
	"twitter/internal/dto"

	"github.com/stretchr/testify/assert"
)

func TestTimelineService(t *testing.T) {
	//t.Skip()

	sdp := SetupTestServiceDependencies(t)
	ctx := context.Background()
	t.Run("test timeline update ", func(t *testing.T) {
		usr1, err := sdp.UserService.CreateUser(ctx, &dto.CreateUserRequest{
			Email:    "tst1@mail.com",
			UserName: "tst1",
		})
		assert.NoError(t, err)
		usr2, err := sdp.UserService.CreateUser(ctx, &dto.CreateUserRequest{
			Email:    "tst2@mail.com",
			UserName: "tst2",
		})
		assert.NoError(t, err)
		_, err = sdp.UserService.FollowUser(ctx, &dto.FollowUserRequest{
			FollowerUserId:  usr1.Id,
			FollowingUserId: usr2.Id,
		})
		assert.NoError(t, err)
		_, err = sdp.PostService.CreatePost(ctx, &dto.CreatePostRequest{
			AuthorId: usr2.Id,
			Subject:  "asdf",
			Contents: "asdfasdf",
		})
		time.Sleep(100 * time.Millisecond)

		tmln, err := sdp.TimelineService.GetTimelineForUser(ctx, &dto.GetTimelineForUserRequest{
			UserId: usr1.Id,
		})

		assert.NoError(t, err)
		assert.NotNil(t, tmln)
		assert.Equal(t, 1, len(tmln.Timeline.Posts))
		assert.Equal(t, "asdfasdf", tmln.Timeline.Posts[0].Contents)
		assert.Equal(t, "asdf", tmln.Timeline.Posts[0].Subject)

	})
}

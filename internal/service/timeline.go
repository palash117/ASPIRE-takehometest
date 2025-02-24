package service

import (
	"context"
	"fmt"
	"twitter/internal/common"
	"twitter/internal/dao"
	"twitter/internal/db"
	"twitter/internal/dto"
	"twitter/internal/model"
)

const (
	POPPULAR_FOLLOW_COUNT = 3
	PAGE_SIZE             = 10
)

func NewTimelineService(dependencies *ServiceDependencies, db *db.Db, lg common.Logger) TimelineService {
	timelineDao := dao.NewTimelineDao(db)
	dao.AutoMigrate(timelineDao.(dao.Migratable))

	timelineServ := &timelineService{
		dependencies:    dependencies,
		timelineDao:     timelineDao,
		postCreatedChan: make(chan *dto.Post, 10),
		lg:              lg,
	}
	go PostCreatedConsumer(timelineServ)
	return timelineServ
}

func PostCreatedConsumer(timelineService *timelineService) {
	for pst := range timelineService.postCreatedChan {
		err := timelineService.UpdateTimelines(context.Background(), pst)
		if err != nil {
			timelineService.lg.Println(fmt.Sprintf("Error updating timeline %s", err))
		}

	}

}

type timelineService struct {
	dependencies    *ServiceDependencies
	timelineDao     dao.TimelineDao
	postCreatedChan chan *dto.Post
	lg              common.Logger
}

func (t *timelineService) UpdateTimelines(ctx context.Context, pst *dto.Post) error {
	author, err := t.dependencies.UserService.GetUser(ctx, &dto.GetUserRequest{UserId: pst.AuthorId})

	if err != nil {

		return err
	}

	if author.FollowerCount >= POPPULAR_FOLLOW_COUNT {
		return nil
	}

	followers, err := t.dependencies.UserService.GetFollowerUsers(ctx, &dto.GetFollowerUsersRequest{PrimaryUserId: pst.AuthorId})
	if err != nil {
		return err
	}
	timelineEntries := make([]*model.Timeline, 0, len(followers.FollowerUsers))
	for _, follower := range followers.FollowerUsers {
		tmln := &model.Timeline{
			UserId:    (int)(follower.FollowerId),
			PostId:    (int)(pst.Id),
			AuthorId:  (int)(pst.AuthorId),
			Subject:   pst.Subject,
			Contents:  pst.Contents,
			CreatedAt: pst.CreatedAt,
		}
		timelineEntries = append(timelineEntries, tmln)
	}
	_, err = t.timelineDao.CreateInBatches(ctx, timelineEntries)
	return err

}

func (t *timelineService) GetTimelineForUser(ctx context.Context, req *dto.GetTimelineForUserRequest) (*dto.GetTimelineForUserResponse, error) {
	resp := &dto.GetTimelineForUserResponse{
		Timeline: dto.Timeline{
			Posts: []*dto.Post{},
		},
		FamousPosts: []*dto.Post{},
		TimelinePaginationResponse: dto.TimelinePaginationResponse{
			PageNo: req.PageNo,
		},
	}
	tmln, err := t.timelineDao.GetTimelineForUser(ctx, (int)(req.UserId), req.PageNo, PAGE_SIZE)
	if err != nil {
		return nil, err
	}
	resp.Timeline = model.ToTimelineDto(tmln)

	if len(tmln) > 0 {
		till := tmln[0].CreatedAt
		from := tmln[len(tmln)-1].CreatedAt

		followings, err := t.dependencies.UserService.GetFollowingUsers(ctx, &dto.GetFollowingUsersRequest{PrimaryUserId: req.UserId})
		usrIds := make([]int, len(followings.FollowingUsers))
		for _, usr := range followings.FollowingUsers {
			usrIds = append(usrIds, (int)(usr.FolloweeId))
		}
		users, err := t.dependencies.UserService.GetUsers(ctx, usrIds)
		if err != nil {
			return nil, err
		}

		famUsrIds := make([]int, 0, len(users.Users))
		for _, usr := range users.Users {
			if usr.FollowerCount >= POPPULAR_FOLLOW_COUNT {
				famUsrIds = append(famUsrIds, (int)(usr.Id))
			}
		}
		popularPosts, err := t.dependencies.PostService.GetPostsByAuthorIdsAndTimeWindow(ctx, &dto.GetPostsByAuthorIdsAndTimeWindowRequest{
			AuthorIds: famUsrIds,
			From:      from,
			Till:      till,
		})
		if err != nil {
			return nil, err
		}
		if popularPosts != nil {
			resp.FamousPosts = popularPosts.Posts
		}

	}
	return resp, nil

}

func (t *timelineService) AfterPostCreateHook(ctx context.Context, createdPost *dto.Post) {
	fmt.Println("AfterPostCreateHook called")
	t.postCreatedChan <- createdPost
}

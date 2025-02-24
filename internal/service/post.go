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

func NewPostService(dependencies *ServiceDependencies, db *db.Db, lg common.Logger) PostService {
	postDao := dao.NewPostDao(db)
	dao.AutoMigrate(postDao.(dao.Migratable))
	return &postService{
		dependencies:   dependencies,
		postDao:        postDao,
		postCreateHook: (dependencies.TimelineService).(AfterPostCreate),
		lg:             lg,
	}
}

type postService struct {
	dependencies   *ServiceDependencies
	postDao        dao.PostDao
	postCreateHook AfterPostCreate
	lg             common.Logger
}

func (p *postService) CreatePost(ctx context.Context, req *dto.CreatePostRequest) (*dto.CreatePostResponse, error) {
	mdl := &model.Post{
		AuthorId: (int)(req.AuthorId),
		Contents: req.Contents,
		Subject:  req.Subject,
	}
	pst, err := p.postDao.Create(ctx, mdl)

	if err != nil {
		p.lg.Println(fmt.Sprintf("Post create err %s", err))

	}
	resp := pst.ToDto()
	p.postCreateHook.AfterPostCreateHook(ctx, resp)
	return &dto.CreatePostResponse{Post: *resp}, err
}

func (p *postService) GetPostsByAuthorIdsAndTimeWindow(ctx context.Context, req *dto.GetPostsByAuthorIdsAndTimeWindowRequest) (*dto.GetPostsByAuthorIdsAndTimeWindowResponse, error) {
	psts, err := p.postDao.GetPostsByAuthorIdsAndTimeWindow(ctx, req.AuthorIds, req.From, req.Till)
	if err != nil {
		return nil, err
	}
	resp := &dto.GetPostsByAuthorIdsAndTimeWindowResponse{Posts: make([]*dto.Post, len(psts))}
	for i, pst := range psts {
		resp.Posts[i] = pst.ToDto()
	}
	return resp, nil
}

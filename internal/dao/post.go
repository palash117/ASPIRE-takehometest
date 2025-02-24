package dao

import (
	"context"
	"time"
	"twitter/internal/db"
	"twitter/internal/model"

	"gorm.io/gorm"
)

type PostDao interface {
	GetPostsByAuthorIdsAndTimeWindow(ctx context.Context, authorIds []int, from, till time.Time) ([]*model.Post, error)
	Create(ctx context.Context, mdl *model.Post) (*model.Post, error)
}
type postDao struct {
	db *gorm.DB
}

func NewPostDao(db *db.Db) PostDao {
	return &postDao{db: db.DB}
}

func (p *postDao) AutoMigrate(dropTable bool) {
	if dropTable {
		p.db.Migrator().DropTable(&model.Timeline{})
	}
	p.db.AutoMigrate(&model.Post{})
}
func (p *postDao) Create(ctx context.Context, mdl *model.Post) (*model.Post, error) {
	rslt := p.db.Create(mdl)
	return mdl, rslt.Error
}

func (p *postDao) GetPostsByAuthorIdsAndTimeWindow(ctx context.Context, authorIds []int, from, till time.Time) ([]*model.Post, error) {
	var posts []*model.Post
	rslt := p.db.Where("author_id in ? and created_at>= ? and created_at <= ? ", authorIds, from, till).Order("created_at desc").Find(&posts)
	return posts, rslt.Error
}

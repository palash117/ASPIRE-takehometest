package dao

import (
	"context"
	"twitter/internal/db"
	"twitter/internal/model"

	"gorm.io/gorm"
)

func AutoMigrate(migratable Migratable) {
	migratable.AutoMigrate(false)
}

type Migratable interface {
	AutoMigrate(bool)
}

type TimelineDao interface {
	Create(ctx context.Context, mdl *model.Timeline) (*model.Timeline, error)
	GetTimelineForUser(ctx context.Context, userId int, pageNo int, pageSize int) ([]*model.Timeline, error)
	CreateInBatches(ctx context.Context, mdl []*model.Timeline) ([]*model.Timeline, error)
}

type timelineDao struct {
	db *gorm.DB
}

func NewTimelineDao(db *db.Db) TimelineDao {
	return &timelineDao{db: db.DB}
}

func (t *timelineDao) AutoMigrate(dropTable bool) {
	if dropTable {
		t.db.Migrator().DropTable(&model.Timeline{})
	}
	t.db.AutoMigrate(&model.Timeline{})
}

func (t *timelineDao) GetTimelineForUser(ctx context.Context, userId int, pageNo int, pageSize int) ([]*model.Timeline, error) {
	var tmln []*model.Timeline
	rslt := t.db.Where("user_id = ?", userId).Limit(pageSize).Offset((pageNo) * pageSize).Order("created_at desc").Find(&tmln)
	return tmln, rslt.Error
}

func (t *timelineDao) Create(ctx context.Context, mdl *model.Timeline) (*model.Timeline, error) {
	rslt := t.db.Create(mdl)
	return mdl, rslt.Error
}

func (t *timelineDao) CreateInBatches(ctx context.Context, mdl []*model.Timeline) ([]*model.Timeline, error) {
	rslt := t.db.CreateInBatches(&mdl, len(mdl))
	return mdl, rslt.Error
}

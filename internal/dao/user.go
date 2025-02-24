package dao

import (
	"context"
	"fmt"
	"time"
	"twitter/internal/db"
	"twitter/internal/model"

	"gorm.io/gorm"
)

type UserDao interface {
	Create(ctx context.Context, mdl *model.User) (*model.User, error)
	GetById(ctx context.Context, userId int) (*model.User, error)
	GetByIds(ctx context.Context, userIds []int) ([]*model.User, error)
	GetAll(ctx context.Context) ([]*model.User, error)
	IsFollowing(ctx context.Context, followerId int, followeeId int) (*model.Follower, error)
	Follow(ctx context.Context, followerId int, followeeId int) (*model.Follower, error)
	GetFollowees(ctx context.Context, userId int) ([]*model.Follower, error)
	GetFollowers(ctx context.Context, userId int) ([]*model.Follower, error)
}

type userDao struct {
	db *gorm.DB
}

func NewUserDao(db *db.Db) UserDao {
	return &userDao{db: db.DB}
}

func (u *userDao) AutoMigrate(dropTable bool) {
	if dropTable {
		u.db.Migrator().DropTable(&model.User{}, &model.Follower{})
	}

	u.db.AutoMigrate(&model.User{})
	u.db.AutoMigrate(&model.Follower{})
}

func (u *userDao) Create(ctx context.Context, mdl *model.User) (*model.User, error) {
	rslt := u.db.Create(mdl)
	return mdl, rslt.Error
}

func (u *userDao) GetById(ctx context.Context, userId int) (*model.User, error) {
	var (
		user model.User
	)
	rslt := u.db.First(&user, userId)
	if user.Id == 0 {
		return nil, EntityNotFound(&user, userId)
	}

	return &user, rslt.Error
}
func (u *userDao) GetByIds(ctx context.Context, userIds []int) ([]*model.User, error) {
	var (
		users []*model.User
	)
	rslt := u.db.Where("id in ?", userIds).Find(&users)
	return users, rslt.Error
}

func (u *userDao) GetAll(ctx context.Context) ([]*model.User, error) {
	var (
		users []*model.User
	)
	u.db.Find(&users)
	return users, nil
}

func (u *userDao) IsFollowing(ctx context.Context, followerId int, followeeId int) (*model.Follower, error) {
	var followers []*model.Follower
	rslt := u.db.Where("follower_id = ? and followee_id = ?", followerId, followeeId).Find(&followers)
	if len(followers) == 0 {
		return nil, db.NoEntriesFound{}
	}
	return followers[0], rslt.Error
}

func (u *userDao) Follow(ctx context.Context, followerId int, followeeId int) (*model.Follower, error) {
	var follower model.Follower
	follower.FollowerId = followerId
	follower.FolloweeId = followeeId
	follower.CreatedAt = time.Now()

	rslt := u.db.Create(&follower)
	if rslt.Error == nil {
		u.db.Model(&model.User{}).Where("id = ?", followeeId).Update("follower_count", gorm.Expr("follower_count + ?", 1))
	}

	return &follower, rslt.Error

}

func (u *userDao) GetFollowers(ctx context.Context, userId int) ([]*model.Follower, error) {
	var followers []*model.Follower
	rslt := u.db.Where("followee_id = ?", userId).Find(&followers)
	return followers, rslt.Error
}

func (u *userDao) GetFollowees(ctx context.Context, userId int) ([]*model.Follower, error) {
	var followers []*model.Follower
	rslt := u.db.Where("follower_id = ?", userId).Find(&followers)
	return followers, rslt.Error
}

func EntityNotFound(ent model.Entity, id int) error {

	return fmt.Errorf(fmt.Sprintf("%s entity not found with id %d", ent.GetName(), id))
}

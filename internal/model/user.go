package model

import (
	"time"
	"twitter/internal/dto"
)

type User struct {
	Id            int            `gorm:"primaryKey;column:id"`
	Email         string         `gorm:"column:email"`
	UserName      string         `gorm:"column:user_name"`
	UserStatus    dto.UserStatus `gorm:"column:user_status"`
	FollowerCount int            `gorm:"column:follower_count"`
	CreatedAt     time.Time      `gorm:"column:created_at"`
	PasswordHash  string         `gorm:"column:password_hash"`
}

func (p *User) GetName() string {
	return "User"
}

func (u *User) ToDto() *dto.User {
	return &dto.User{
		Id:            (dto.Id)(u.Id),
		Email:         u.Email,
		UserName:      u.UserName,
		UserStatus:    u.UserStatus,
		FollowerCount: u.FollowerCount,
		CreatedAt:     u.CreatedAt,
	}
}

func FromUserDto(u *dto.User) *User {
	return &User{
		Id:            (int)(u.Id),
		Email:         u.Email,
		UserName:      u.UserName,
		UserStatus:    u.UserStatus,
		FollowerCount: u.FollowerCount,
		CreatedAt:     u.CreatedAt,
	}
}

type Follower struct {
	Id         int       `gorm:"primaryKey;column:id"`
	FollowerId int       `gorm:"column:follower_id"`
	FolloweeId int       `gorm:"column:followee_id"`
	CreatedAt  time.Time `gorm:"column:created_at"`
}

func (f *Follower) ToDto() *dto.Follower {
	return &dto.Follower{

		FollowId:   dto.Id(f.Id),
		FollowerId: dto.Id(f.FollowerId),
		FolloweeId: dto.Id(f.FolloweeId),
		CreatedAt:  f.CreatedAt,
	}
}

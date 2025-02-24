package model

import (
	"time"
	"twitter/internal/dto"
)

type Timeline struct {
	Id        int       `gorm:"primaryKey;column:id"`
	UserId    int       `gorm:"column:user_id"`
	PostId    int       `gorm:"column:post_id"`
	AuthorId  int       `gorm:"column:author_id"`
	Subject   string    `gorm:"column:subject"`
	Contents  string    `gorm:"column:contents"`
	CreatedAt time.Time `gorm:"column:created_at"`
}

func ToTimelineDto(ts []*Timeline) dto.Timeline {
	dtoTmln := dto.Timeline{
		Posts: make([]*dto.Post, len(ts)),
	}
	for i, t := range ts {
		dtoTmln.Posts[i] = &dto.Post{
			Id:        dto.Id(t.PostId),
			AuthorId:  dto.Id(t.AuthorId),
			Subject:   t.Subject,
			Contents:  t.Contents,
			CreatedAt: t.CreatedAt,
		}

	}
	return dtoTmln
}

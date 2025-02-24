package model

import (
	"time"
	"twitter/internal/dto"
)

type Entity interface {
	GetName() string
}

type Post struct {
	Id        int       `gorm:"primaryKey;column:id"`
	AuthorId  int       `gorm:"primaryKey;column:author_id"`
	Subject   string    `gorm:"primaryKey;column:subject"`
	Contents  string    `gorm:"primaryKey;column:contents"`
	CreatedAt time.Time `gorm:"primaryKey;column:created_at"`
}

func (p *Post) ToDto() *dto.Post {
	return &dto.Post{
		Id:        dto.Id(p.Id),
		AuthorId:  dto.Id(p.AuthorId),
		Subject:   p.Subject,
		Contents:  p.Contents,
		CreatedAt: p.CreatedAt,
	}

}

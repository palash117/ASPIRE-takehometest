package dto

import (
	"time"
)

type Id int

type Timeline struct {
	Posts []*Post
}
type TimelinePaginationRequest struct {
	PageNo int `json:"pageNo"`
}
type TimelinePaginationResponse struct {
	PageNo int `json:"pageNo"`
}

type GetTimelineForUserRequest struct {
	UserId Id
	TimelinePaginationRequest
}
type GetTimelineForUserResponse struct {
	Timeline    `json:"timeline"`
	FamousPosts []*Post `json:"famousPosts"`
	TimelinePaginationResponse
}

type UserStatus int32

const (
	UserStatus_Unspecified = iota
	UserStatus_Activated
	UserStatus_Deactivated
)

type User struct {
	Id            Id         `json:"id"`
	Email         string     `json:"email"`
	UserName      string     `json:"userName"`
	UserStatus    UserStatus `json:"userStatus"`
	FollowerCount int        `json:"followerCount"`
	CreatedAt     time.Time  `json:"createdAt"`
}

type Follower struct {
	FollowId   Id        `json:"id"`
	FollowerId Id        `json:"followerId"`
	FolloweeId Id        `json:"followeeId"`
	CreatedAt  time.Time `json:"createdAt"`
}

type CreateUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	UserName string `json:"userName"`
}
type CreateUserResponse struct {
	User
}
type FollowUserRequest struct {
	FollowerUserId  Id `json:"followerUserId"`
	FollowingUserId Id `json:"followUserId"`
}
type FollowUserResponse struct {
	Follower
}

type GetFollowingUsersRequest struct {
	PrimaryUserId Id `json:"userId"`
}
type GetFollowingUsersResponse struct {
	FollowingUsers []*Follower `json:"followers"`
}
type GetFollowerUsersRequest struct {
	PrimaryUserId Id `json:"userId"`
}
type GetFollowerUsersResponse struct {
	FollowerUsers []*Follower `json:"followers"`
}
type GetAllUsersRequest struct {
	PageSize   int
	PageNumber int
}
type GetAllUsersResponse struct {
	Users []*User `json:"users"`
}
type GetUserRequest struct {
	UserId Id `json:"userId"`
}
type GetUserResponse struct {
	User
}

type Post struct {
	Id        Id
	AuthorId  Id
	Subject   string
	Contents  string
	CreatedAt time.Time
}

type CreatePostRequest struct {
	AuthorId Id
	Subject  string `json:"subject"`
	Contents string `json:"contents"`
}
type CreatePostResponse struct {
	Post
}
type GetPostsByAuthorIdsAndTimeWindowRequest struct {
	AuthorIds []int
	From      time.Time
	Till      time.Time
}

type GetPostsByAuthorIdsAndTimeWindowResponse struct {
	Posts []*Post
}

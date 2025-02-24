package controller

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
	"twitter/internal/dto"
	"twitter/internal/service"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type PostMock struct {
	mock.Mock
}

func (p *PostMock) CreatePost(ctx context.Context, req *dto.CreatePostRequest) (*dto.CreatePostResponse, error) {
	args := p.Called(ctx, req)
	resp := args[0].(*dto.CreatePostResponse)
	var err error
	if args[1] != nil {
		err = args[1].(error)

	}
	return resp, err
}
func (p *PostMock) GetPostsByAuthorIdsAndTimeWindow(ctx context.Context, req *dto.GetPostsByAuthorIdsAndTimeWindowRequest) (*dto.GetPostsByAuthorIdsAndTimeWindowResponse, error) {

	args := p.Called(req)
	return args[0].(*dto.GetPostsByAuthorIdsAndTimeWindowResponse), args[1].(error)
}

type UserMock struct {
	mock.Mock
}

func (u *UserMock) CreateUser(ctx context.Context, req *dto.CreateUserRequest) (*dto.CreateUserResponse, error) {
	args := u.Called(ctx, req)
	resp := args[0].(*dto.CreateUserResponse)
	var err error
	if args[1] != nil {
		err = args[1].(error)

	}
	return resp, err
}
func (u *UserMock) FollowUser(ctx context.Context, req *dto.FollowUserRequest) (*dto.FollowUserResponse, error) {
	args := u.Called(ctx, req)
	resp := args[0].(*dto.FollowUserResponse)
	var err error
	if args[1] != nil {
		err = args[1].(error)

	}
	return resp, err

}
func (u *UserMock) GetFollowingUsers(ctx context.Context, req *dto.GetFollowingUsersRequest) (*dto.GetFollowingUsersResponse, error) {
	args := u.Called(ctx, req)
	resp := args[0].(*dto.GetFollowingUsersResponse)
	var err error
	if args[1] != nil {
		err = args[1].(error)

	}
	return resp, err

}

func (u *UserMock) GetFollowerUsers(ctx context.Context, req *dto.GetFollowerUsersRequest) (*dto.GetFollowerUsersResponse, error) {
	args := u.Called(ctx, req)
	resp := args[0].(*dto.GetFollowerUsersResponse)
	var err error
	if args[1] != nil {
		err = args[1].(error)

	}
	return resp, err

}
func (u *UserMock) GetAllUsers(ctx context.Context, req *dto.GetAllUsersRequest) (*dto.GetAllUsersResponse, error) {
	args := u.Called(ctx, req)
	resp := args[0].(*dto.GetAllUsersResponse)
	var err error
	if args[1] != nil {
		err = args[1].(error)

	}
	return resp, err
}
func (u *UserMock) GetUser(ctx context.Context, req *dto.GetUserRequest) (*dto.GetUserResponse, error) {
	args := u.Called(ctx, req)
	return args[0].(*dto.GetUserResponse), args[1].(error)
}
func (u *UserMock) GetUsers(ctx context.Context, req []int) (*dto.GetAllUsersResponse, error) {
	args := u.Called(ctx, req)
	return args[0].(*dto.GetAllUsersResponse), args[1].(error)
}

type TimelineMock struct {
	mock.Mock
}

func (t *TimelineMock) GetTimelineForUser(ctx context.Context, req *dto.GetTimelineForUserRequest) (*dto.GetTimelineForUserResponse, error) {
	args := t.Called(ctx, req)
	resp := args[0].(*dto.GetTimelineForUserResponse)
	var err error
	if args[1] != nil {
		err = args[1].(error)

	}
	return resp, err
}

func setupServiceDependencies(t *testing.T) *service.ServiceDependencies {
	t.Helper()
	sdp := &service.ServiceDependencies{
		UserService:     &UserMock{},
		PostService:     &PostMock{},
		TimelineService: &TimelineMock{},
	}
	return sdp
}

type logRecorder struct {
	prints []string
}

func NewLogRecorder() *logRecorder {
	return &logRecorder{
		prints: []string{},
	}

}

func (l *logRecorder) Println(str string) {
	l.prints = append(l.prints, str)

}

func TestHandlers(t *testing.T) {
	userId := 12
	var nilErr error
	t.Run("", func(t *testing.T) {

	})

	t.Run("test create post green", func(t *testing.T) {
		//t.Skip()
		lg := NewLogRecorder()
		ctx := context.Background()
		sdp := setupServiceDependencies(t)
		robj := &dto.CreatePostRequest{
			AuthorId: dto.Id(userId),
			Subject:  "asdf",
			Contents: "asdfasdf",
		}
		mockedSrvResp := &dto.CreatePostResponse{
			Post: dto.Post{
				AuthorId: dto.Id(userId),
				Subject:  "asdf",
				Contents: "asdfasdf",
			},
		}
		postsMock := sdp.PostService.(*PostMock)
		postsMock.On("CreatePost", mock.Anything, robj).Return(mockedSrvResp, nilErr)

		cntrlr := NewController(sdp, lg)

		req, err := http.NewRequest("POST", "/twitter/12/post", stringify(t, robj))
		req = req.WithContext(ctx)
		assert.NoError(t, err)
		rr := httptest.NewRecorder()
		router := mux.NewRouter()
		router.HandleFunc("/twitter/{user_id}/post", cntrlr.CreatePostHandler)

		router.ServeHTTP(rr, req)
		got := dto.CreatePostResponse{}
		assert.NoError(t, parse(t, rr, &got), "error parsing response stream")

		assert.Equal(t, int(got.AuthorId), userId)
		assert.Equal(t, got.Subject, "asdf")
		assert.Equal(t, got.Contents, "asdfasdf")

	})

	t.Run("test create post red oversized contents", func(t *testing.T) {
		//t.Skip()
		lg := NewLogRecorder()
		ctx := context.Background()
		sdp := setupServiceDependencies(t)
		robj := &dto.CreatePostRequest{
			AuthorId: dto.Id(userId),
			Subject:  "asdf",
			Contents: "aasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfsdfasdf",
		}
		mockedSrvResp := &dto.CreatePostResponse{
			Post: dto.Post{
				AuthorId: dto.Id(userId),
				Subject:  "asdf",
				Contents: "asdfasdf",
			},
		}
		postsMock := sdp.PostService.(*PostMock)
		postsMock.On("CreatePost", mock.Anything, robj).Return(mockedSrvResp, nilErr)

		cntrlr := NewController(sdp, lg)

		req, err := http.NewRequest("POST", "/twitter/12/post", stringify(t, robj))
		req = req.WithContext(ctx)
		assert.NoError(t, err)
		rr := httptest.NewRecorder()
		router := mux.NewRouter()
		router.HandleFunc("/twitter/{user_id}/post", cntrlr.CreatePostHandler)

		router.ServeHTTP(rr, req)
		//got := dto.CreatePostResponse{}
		//assert.NoError(t, parse(t, rr, &got), "error parsing response stream")

		assert.Equal(t, rr.Body.String(), "{\"field\":\"contents\",\"message\":\"Larger than limit 100\"}\n")

	})
	t.Run("test create post red oversized subject", func(t *testing.T) {
		//t.Skip()
		lg := NewLogRecorder()
		ctx := context.Background()
		sdp := setupServiceDependencies(t)

		robj := &dto.CreatePostRequest{
			AuthorId: dto.Id(userId),
			Subject:  "aasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdf",
			Contents: "asdf asdf asdf ",
		}
		mockedSrvResp := &dto.CreatePostResponse{
			Post: dto.Post{
				AuthorId: dto.Id(userId),
				Subject:  "asdf",
				Contents: "asdfasdf",
			},
		}
		postsMock := sdp.PostService.(*PostMock)
		postsMock.On("CreatePost", mock.Anything, robj).Return(mockedSrvResp, nilErr)

		cntrlr := NewController(sdp, lg)

		req, err := http.NewRequest("POST", "/twitter/12/post", stringify(t, robj))
		req = req.WithContext(ctx)
		assert.NoError(t, err)
		rr := httptest.NewRecorder()
		router := mux.NewRouter()
		router.HandleFunc("/twitter/{user_id}/post", cntrlr.CreatePostHandler)

		router.ServeHTTP(rr, req)
		//got := dto.CreatePostResponse{}
		//assert.NoError(t, parse(t, rr, &got), "error parsing response stream")

		assert.Equal(t, "{\"field\":\"subject\",\"message\":\"Larger than limit 50\"}\n", rr.Body.String())

	})

	t.Run("test create user", func(t *testing.T) {
		//t.Skip()
		lg := NewLogRecorder()
		ctx := context.Background()
		sdp := setupServiceDependencies(t)
		robj := &dto.CreateUserRequest{
			UserName: "asfd",
			Email:    "asdf@asfd.com",
		}
		mockedSrvResp := &dto.CreateUserResponse{
			User: dto.User{
				Id:       1,
				UserName: "asdf",
				Email:    "asdf@asdf.com",
			},
		}
		userMock := sdp.UserService.(*UserMock)
		userMock.On("CreateUser", mock.Anything, robj).Return(mockedSrvResp, nilErr)

		cntrlr := NewController(sdp, lg)

		req, err := http.NewRequest("POST", "/twitter/user", stringify(t, robj))
		req = req.WithContext(ctx)
		assert.NoError(t, err)
		rr := httptest.NewRecorder()
		router := mux.NewRouter()
		router.HandleFunc("/twitter/user", cntrlr.CreateUserHandler)

		router.ServeHTTP(rr, req)
		got := dto.CreateUserResponse{}
		assert.NoError(t, parse(t, rr, &got), "error parsing response stream")

		assert.NotEqual(t, 0, int(got.Id))
		assert.Equal(t, "asdf", got.UserName)
		assert.Equal(t, "asdf@asdf.com", got.Email)

	})

	t.Run("test get all users", func(t *testing.T) {
		//t.Skip()
		lg := NewLogRecorder()
		ctx := context.Background()
		sdp := setupServiceDependencies(t)
		robj := &dto.GetAllUsersRequest{}
		mockedSrvResp := &dto.GetAllUsersResponse{
			Users: []*dto.User{
				{
					Id:       1,
					UserName: "asdf",
					Email:    "asdf@asdf.com",
				},
				{
					Id:       1,
					UserName: "asdf",
					Email:    "asdf@asdf.com",
				},
			},
		}
		userMock := sdp.UserService.(*UserMock)
		userMock.On("GetAllUsers", mock.Anything, robj).Return(mockedSrvResp, nilErr)

		cntrlr := NewController(sdp, lg)

		req, err := http.NewRequest("POST", "/twitter/user/all", stringify(t, robj))
		req = req.WithContext(ctx)
		assert.NoError(t, err)
		rr := httptest.NewRecorder()
		router := mux.NewRouter()
		router.HandleFunc("/twitter/user/all", cntrlr.GetAllUsersHandler)

		router.ServeHTTP(rr, req)
		got := dto.GetAllUsersResponse{}
		assert.NoError(t, parse(t, rr, &got), "error parsing response stream")

		assert.Equal(t, 2, len(got.Users))

	})

	t.Run("test get user timeline ", func(t *testing.T) {
		lg := NewLogRecorder()
		ctx := context.Background()
		sdp := setupServiceDependencies(t)
		robj := &dto.GetTimelineForUserRequest{
			UserId: 1,
		}
		mockedSrvResp := &dto.GetTimelineForUserResponse{
			Timeline: dto.Timeline{
				Posts: []*dto.Post{
					{
						Id:        1,
						AuthorId:  2,
						Subject:   "asdf",
						Contents:  "asdfasdf",
						CreatedAt: getSampleDateObj(),
					},
					{
						Id:        1,
						AuthorId:  2,
						Subject:   "asdf",
						Contents:  "asdfasdf",
						CreatedAt: getSampleDateObj(),
					},
				},
			},
			FamousPosts: []*dto.Post{
				{
					Id:        1,
					AuthorId:  2,
					Subject:   "asdf",
					Contents:  "asdfasdf",
					CreatedAt: getSampleDateObj(),
				},
			},
		}
		timelineMock := sdp.TimelineService.(*TimelineMock)
		timelineMock.On("GetTimelineForUser", mock.Anything, robj).Return(mockedSrvResp, nilErr)

		cntrlr := NewController(sdp, lg)

		req, err := http.NewRequest("POST", "/twitter/1/timeline", stringify(t, robj))
		req = req.WithContext(ctx)
		assert.NoError(t, err)
		rr := httptest.NewRecorder()
		router := mux.NewRouter()
		router.HandleFunc("/twitter/{user_id}/timeline", cntrlr.GetTimelineHandler)

		router.ServeHTTP(rr, req)
		got := dto.GetTimelineForUserResponse{}
		assert.NoError(t, parse(t, rr, &got), "error parsing response stream")

		assert.Equal(t, 2, len(got.Timeline.Posts))
		assert.Equal(t, 1, len(got.FamousPosts))

	})
	t.Run("test get user timeline ", func(t *testing.T) {
		lg := NewLogRecorder()
		ctx := context.Background()
		sdp := setupServiceDependencies(t)
		robj := &dto.GetTimelineForUserRequest{
			UserId: 1,
		}
		mockedSrvResp := &dto.GetTimelineForUserResponse{
			Timeline: dto.Timeline{
				Posts: []*dto.Post{
					{
						Id:        1,
						AuthorId:  2,
						Subject:   "asdf",
						Contents:  "asdfasdf",
						CreatedAt: getSampleDateObj(),
					},
					{
						Id:        1,
						AuthorId:  2,
						Subject:   "asdf",
						Contents:  "asdfasdf",
						CreatedAt: getSampleDateObj(),
					},
				},
			},
			FamousPosts: []*dto.Post{
				{
					Id:        1,
					AuthorId:  2,
					Subject:   "asdf",
					Contents:  "asdfasdf",
					CreatedAt: getSampleDateObj(),
				},
			},
		}
		timelineMock := sdp.TimelineService.(*TimelineMock)
		timelineMock.On("GetTimelineForUser", mock.Anything, robj).Return(mockedSrvResp, nilErr)

		cntrlr := NewController(sdp, lg)

		req, err := http.NewRequest("POST", "/twitter/1/timeline", stringify(t, robj))
		req = req.WithContext(ctx)
		assert.NoError(t, err)
		rr := httptest.NewRecorder()
		router := mux.NewRouter()
		router.HandleFunc("/twitter/{user_id}/timeline", cntrlr.GetTimelineHandler)

		router.ServeHTTP(rr, req)
		got := dto.GetTimelineForUserResponse{}
		assert.NoError(t, parse(t, rr, &got), "error parsing response stream")

		assert.Equal(t, 2, len(got.Timeline.Posts))
		assert.Equal(t, 1, len(got.FamousPosts))

	})

	t.Run("test follow user", func(t *testing.T) {
		//t.Skip()
		lg := NewLogRecorder()
		ctx := context.Background()
		sdp := setupServiceDependencies(t)
		robj := &dto.FollowUserRequest{
			FollowerUserId:  dto.Id(1),
			FollowingUserId: dto.Id(2),
		}
		mockedSrvResp := &dto.FollowUserResponse{
			Follower: dto.Follower{
				FollowId:   1,
				FollowerId: 12,
				FolloweeId: 2,
				CreatedAt:  getSampleDateObj(),
			},
		}
		userMock := sdp.UserService.(*UserMock)
		userMock.On("FollowUser", mock.Anything, robj).Return(mockedSrvResp, nilErr)

		cntrlr := NewController(sdp, lg)

		req, err := http.NewRequest("POST", "/twitter/1/follow", stringify(t, robj))
		req = req.WithContext(ctx)
		assert.NoError(t, err)
		rr := httptest.NewRecorder()
		router := mux.NewRouter()
		router.HandleFunc("/twitter/{user_id}/follow", cntrlr.FollowUserHandler)

		router.ServeHTTP(rr, req)
		got := dto.FollowUserResponse{}
		assert.NoError(t, parse(t, rr, &got), "error parsing response stream")

		assert.Equal(t, dto.Id(12), got.Follower.FollowerId)
		assert.Equal(t, dto.Id(2), got.Follower.FolloweeId)

	})
	t.Run("test get followers", func(t *testing.T) {
		lg := NewLogRecorder()
		ctx := context.Background()
		sdp := setupServiceDependencies(t)
		robj := &dto.GetFollowerUsersRequest{
			PrimaryUserId: dto.Id(1),
		}
		mockedSrvResp := &dto.GetFollowerUsersResponse{
			FollowerUsers: []*dto.Follower{
				{
					FollowId:   dto.Id(1),
					FollowerId: dto.Id(2),
					FolloweeId: dto.Id(1),
					CreatedAt:  getSampleDateObj(),
				},
			},
		}
		userMock := sdp.UserService.(*UserMock)
		userMock.On("GetFollowerUsers", mock.Anything, robj).Return(mockedSrvResp, nilErr)

		cntrlr := NewController(sdp, lg)

		req, err := http.NewRequest("POST", "/twitter/1/followers", stringify(t, robj))
		req = req.WithContext(ctx)
		assert.NoError(t, err)
		rr := httptest.NewRecorder()
		router := mux.NewRouter()
		router.HandleFunc("/twitter/{user_id}/followers", cntrlr.GetFollowersHandler)

		router.ServeHTTP(rr, req)
		got := dto.GetFollowerUsersResponse{}
		assert.NoError(t, parse(t, rr, &got), "error parsing response stream")

		assert.Equal(t, 1, len(got.FollowerUsers))

	})
	t.Run("test get followings", func(t *testing.T) {
		lg := NewLogRecorder()
		ctx := context.Background()
		sdp := setupServiceDependencies(t)
		robj := &dto.GetFollowingUsersRequest{
			PrimaryUserId: dto.Id(1),
		}
		mockedSrvResp := &dto.GetFollowingUsersResponse{
			FollowingUsers: []*dto.Follower{
				{
					FollowId:   dto.Id(1),
					FollowerId: dto.Id(2),
					FolloweeId: dto.Id(1),
					CreatedAt:  getSampleDateObj(),
				},
			},
		}
		userMock := sdp.UserService.(*UserMock)
		userMock.On("GetFollowingUsers", mock.Anything, robj).Return(mockedSrvResp, nilErr)

		cntrlr := NewController(sdp, lg)

		req, err := http.NewRequest("POST", "/twitter/1/following", stringify(t, robj))
		req = req.WithContext(ctx)
		assert.NoError(t, err)
		rr := httptest.NewRecorder()
		router := mux.NewRouter()
		router.HandleFunc("/twitter/{user_id}/following", cntrlr.GetFolloweesHandler)

		router.ServeHTTP(rr, req)
		got := dto.GetFollowingUsersResponse{}
		assert.NoError(t, parse(t, rr, &got), "error parsing response stream")

		assert.Equal(t, 1, len(got.FollowingUsers))

	})

}

func parse(t *testing.T, rr *httptest.ResponseRecorder, b any) error {
	t.Helper()
	err := json.Unmarshal([]byte(rr.Body.String()), &b)
	return err
}
func stringify(t *testing.T, obj any) io.Reader {
	byts, _ := json.Marshal(obj)
	return bytes.NewReader(byts)
}
func getSampleDateObj() time.Time {
	layout := "2006-01-02 15:04:05"
	value := "2025-02-24 00:00:00"
	date, _ := time.Parse(layout, value)
	return date
}

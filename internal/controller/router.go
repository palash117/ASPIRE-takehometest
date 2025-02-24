package controller

import (
	"net/http"
	"twitter/internal/common"
	"twitter/internal/service"

	"github.com/gorilla/mux"
)

type Controller interface {
	RegistreRoutes(*mux.Router)
	CreatePostHandler(w http.ResponseWriter, r *http.Request)
	GetTimelineHandler(w http.ResponseWriter, r *http.Request)
	CreateUserHandler(w http.ResponseWriter, r *http.Request)
	GetFollowersHandler(w http.ResponseWriter, r *http.Request)
	FollowUserHandler(w http.ResponseWriter, r *http.Request)
	GetFolloweesHandler(w http.ResponseWriter, r *http.Request)
	GetAllUsersHandler(w http.ResponseWriter, r *http.Request)
}

func NewController(dependencies *service.ServiceDependencies, lg common.Logger) Controller {
	return &cntrlr{dependencies: dependencies, lg: lg}
}

type cntrlr struct {
	dependencies *service.ServiceDependencies
	lg           common.Logger
}

func (c *cntrlr) RegistreRoutes(r *mux.Router) {

	twitterRouter := r.PathPrefix("/twitter").Subrouter()
	twitterRouter.HandleFunc("/user", c.CreateUserHandler).Methods("Post")
	twitterRouter.HandleFunc("/user/all", c.GetAllUsersHandler).Methods("Get")
	twitterRouter.HandleFunc("/{user_id}/post", c.CreatePostHandler).Methods("Post")
	twitterRouter.HandleFunc("/{user_id}/timeline", c.GetTimelineHandler).Methods("Get")
	twitterRouter.HandleFunc("/user/{user_id}/followers", c.GetFollowersHandler).Methods("Get")
	twitterRouter.HandleFunc("/user/{user_id}/following", c.GetFolloweesHandler).Methods("Get")
	twitterRouter.HandleFunc("/user/{user_id}/follow", c.FollowUserHandler).Methods("Post")

	twitterRouter.Use(mux.CORSMethodMiddleware(r))
	twitterRouter.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Route not found", http.StatusNotFound)
	})

}

package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"twitter/internal/dto"

	"github.com/gorilla/mux"
)

func (c *cntrlr) CreatePostHandler(w http.ResponseWriter, r *http.Request) {

	var (
		createPostReq  dto.CreatePostRequest
		createPostResp *dto.CreatePostResponse
		err            error
	)
	defer r.Body.Close()
	err = json.NewDecoder(r.Body).Decode(&createPostReq)

	if err != nil {
		c.handleHttpErr(fmt.Sprintf("err parsing req: %s", err), http.StatusBadRequest, w)
		return
	}
	createPostReq.AuthorId, err = getUserId(r)

	if err != nil {
		c.handleHttpErr(fmt.Sprintf("err parsing req: %s", err), http.StatusBadRequest, w)
		return
	}
	httpErr := validatePostReq(&createPostReq)
	if httpErr != nil {
		c.handleHttpErrObj(httpErr, http.StatusBadRequest, w)
		return
	}
	createPostResp, err = c.dependencies.PostService.CreatePost(r.Context(), &createPostReq)
	if err != nil {

		c.handleHttpErr(fmt.Sprintf("internal err %s", err), http.StatusInternalServerError, w)
		return
	}
	c.writeToResponseStream(createPostResp, http.StatusOK, w)

}
func (c *cntrlr) GetTimelineHandler(w http.ResponseWriter, r *http.Request) {

	var (
		getTimelineReq  dto.GetTimelineForUserRequest
		getTimelineResp *dto.GetTimelineForUserResponse
		err             error
	)
	getTimelineReq.UserId, err = getUserId(r)
	if err != nil {
		c.handleHttpErr(fmt.Sprintf("err parsing req: %s", err), http.StatusBadRequest, w)
		return
	}
	getTimelineReq.PageNo = getPageNo(r)
	getTimelineResp, err = c.dependencies.TimelineService.GetTimelineForUser(r.Context(), &getTimelineReq)
	if err != nil {

		c.handleHttpErr(fmt.Sprintf("internal err %s", err), http.StatusInternalServerError, w)
		return
	}
	c.writeToResponseStream(getTimelineResp, http.StatusOK, w)

}

// managment apis
func (c *cntrlr) CreateUserHandler(w http.ResponseWriter, r *http.Request) {

	var (
		createUserReq  dto.CreateUserRequest
		createUserResp *dto.CreateUserResponse
		err            error
	)
	defer r.Body.Close()
	err = json.NewDecoder(r.Body).Decode(&createUserReq)

	if err != nil {
		c.handleHttpErr(fmt.Sprintf("err parsing req: %s", err), http.StatusBadRequest, w)
		return
	}
	createUserResp, err = c.dependencies.UserService.CreateUser(r.Context(), &createUserReq)
	if err != nil {

		c.handleHttpErr(fmt.Sprintf("internal err %s", err), http.StatusInternalServerError, w)
		return
	}
	c.writeToResponseStream(createUserResp, http.StatusOK, w)

}

func (c *cntrlr) GetFollowersHandler(w http.ResponseWriter, r *http.Request) {

	var (
		req  = dto.GetFollowerUsersRequest{}
		resp *dto.GetFollowerUsersResponse
		err  error
	)

	req.PrimaryUserId, err = getUserId(r)
	if err != nil {
		c.handleHttpErr(fmt.Sprintf("err parsing req: %s", err), http.StatusBadRequest, w)
		return
	}
	resp, err = c.dependencies.UserService.GetFollowerUsers(r.Context(), &req)
	if err != nil {

		c.handleHttpErr(fmt.Sprintf("internal err %s", err), http.StatusInternalServerError, w)
		return
	}
	c.writeToResponseStream(resp, http.StatusOK, w)

}

func (c *cntrlr) FollowUserHandler(w http.ResponseWriter, r *http.Request) {
	c.lg.Println("GetFollowees called")

	var (
		req  = dto.FollowUserRequest{}
		resp *dto.FollowUserResponse
		err  error
	)
	defer r.Body.Close()
	err = json.NewDecoder(r.Body).Decode(&req)

	req.FollowerUserId, err = getUserId(r)
	if err != nil {
		c.handleHttpErr(fmt.Sprintf("err parsing req: %s", err), http.StatusBadRequest, w)
		return
	}
	resp, err = c.dependencies.UserService.FollowUser(r.Context(), &req)
	if err != nil {

		c.handleHttpErr(fmt.Sprintf("internal err %s", err), http.StatusInternalServerError, w)
		return
	}
	c.writeToResponseStream(resp, http.StatusOK, w)

}

func (c *cntrlr) GetFolloweesHandler(w http.ResponseWriter, r *http.Request) {
	c.lg.Println("GetFollowees called")

	var (
		req  = dto.GetFollowingUsersRequest{}
		resp *dto.GetFollowingUsersResponse
		err  error
	)
	req.PrimaryUserId, err = getUserId(r)
	if err != nil {
		c.handleHttpErr(fmt.Sprintf("err parsing req: %s", err), http.StatusBadRequest, w)
		return
	}
	resp, err = c.dependencies.UserService.GetFollowingUsers(r.Context(), &req)
	if err != nil {

		c.handleHttpErr(fmt.Sprintf("internal err %s", err), http.StatusInternalServerError, w)
		return
	}
	c.writeToResponseStream(resp, http.StatusOK, w)

}

func (c *cntrlr) GetAllUsersHandler(w http.ResponseWriter, r *http.Request) {
	c.lg.Println("calling GetAllUsers")

	var (
		req  dto.GetAllUsersRequest
		resp *dto.GetAllUsersResponse
		err  error
	)

	if err != nil {
		c.handleHttpErr(fmt.Sprintf("err parsing req: %s", err), http.StatusBadRequest, w)
		return
	}
	resp, err = c.dependencies.UserService.GetAllUsers(r.Context(), &req)
	if err != nil {

		c.handleHttpErr(fmt.Sprintf("internal err %s", err), http.StatusInternalServerError, w)
		return
	}
	c.writeToResponseStream(resp, http.StatusOK, w)

}

func (c *cntrlr) writeToResponseStream(resp any, statusCode int, w http.ResponseWriter) {
	//respJson, _ := json.Marshal(resp)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
func getUserId(r *http.Request) (dto.Id, error) {
	pathParams := mux.Vars(r)
	id, err := strconv.Atoi(pathParams["user_id"])
	if err != nil {
		return 0, fmt.Errorf("invalid user_id")
	}
	return dto.Id(id), nil
}

// func getQueryParam(r *http.Request, query string)
type HttpError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func (h *HttpError) Error() string {
	if h.Field != "" {
		return fmt.Sprintf("field: %s, messge: %s", h.Field, h.Message)
	}
	return fmt.Sprintf("message: %s", h.Message)
}

func getPageNo(r *http.Request) int {
	pageNoStr := r.URL.Query().Get("pageNo")
	pageNo, _ := strconv.Atoi(pageNoStr)
	if pageNo < 0 {
		pageNo = 0
	}
	return pageNo
}

func (c *cntrlr) handleHttpErr(errString string, statusCode int, w http.ResponseWriter) {
	c.lg.Println(errString)
	herr := &HttpError{Message: errString}
	c.handleHttpErrObj(herr, statusCode, w)
	return
}
func (c *cntrlr) handleHttpErrObj(err *HttpError, statusCode int, w http.ResponseWriter) {
	errString, _ := json.Marshal(err)
	http.Error(w, string(errString), statusCode)
	return
}
func validatePostReq(req *dto.CreatePostRequest) *HttpError {
	if len(req.Subject) > 50 {
		return &HttpError{Field: "subject", Message: "Larger than limit 50"}
	}
	if len(req.Contents) > 100 {
		return &HttpError{Field: "contents", Message: "Larger than limit 100"}
	}
	if len(req.Subject) == 0 {
		return &HttpError{Field: "subject", Message: "Empty not allowerd"}
	}
	if len(req.Contents) == 100 {
		return &HttpError{Field: "contents", Message: "Empty not allowed"}
	}

	return nil
}

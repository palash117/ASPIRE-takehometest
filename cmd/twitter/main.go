package main

import (
	"fmt"
	"net/http"
	"twitter/internal"
	"twitter/internal/common"
	"twitter/internal/controller"
	"twitter/internal/db"
	"twitter/internal/service"

	"github.com/gorilla/mux"
)

var (
	dependencies *internal.Dependencies
	cntrlr       controller.Controller
)

func main() {
	initialize()
	startServer()
}

func initialize() {

	lg := common.NewLogger()
	db := db.SetupDb()
	dependencies := service.NewServiceDependencies(db, lg)
	cntrlr = controller.NewController(dependencies, lg)
}
func startServer() {
	r := mux.NewRouter()
	cntrlr.RegistreRoutes(r)
	err := http.ListenAndServe(":8080", r)
	if err != nil {
		dependencies.Logger.Println(fmt.Sprintf("err is %s", err))

	}
}

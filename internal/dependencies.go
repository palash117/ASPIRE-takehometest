package internal

import (
	"twitter/internal/common"
	"twitter/internal/db"
	"twitter/internal/service"
)

type Dependencies struct {
	ServiceDependencies *service.ServiceDependencies
	Logger              common.Logger
	Db                  *db.Db
}

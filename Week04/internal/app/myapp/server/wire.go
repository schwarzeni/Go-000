// +build wireinject

package server

import (
	"github.com/google/wire"
	"github.com/schwarzeni/Go-000/Week04/internal/app/myapp/dao"
	"github.com/schwarzeni/Go-000/Week04/internal/app/myapp/service"
	"github.com/schwarzeni/Go-000/Week04/pkg/db"
)

func NewServerWithWire(svcConfig *HTTPServerConfig, dbConfig *db.DBConfig) (*HTTPServer, error) {
	wire.Build(NewHTTPServer, NewHandler, service.NewService, dao.NewDao, db.NewGorm)
	return &HTTPServer{}, nil
}

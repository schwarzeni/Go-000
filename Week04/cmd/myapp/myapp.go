package main

import (
	"log"
	"syscall"

	"github.com/schwarzeni/Go-000/Week04/internal/app/myapp/server"
	"github.com/schwarzeni/Go-000/Week04/pkg/appgroup"
	"github.com/schwarzeni/Go-000/Week04/pkg/db"
)

func main() {
	httpServer, err := server.NewServerWithWire(
		&server.HTTPServerConfig{Addr: ":8080"},
		&db.DBConfig{
			URL:     "127.0.0.1:3306",
			User:    "root",
			PassWD:  "root",
			CharSet: "utf8mb4",
			DBName:  "learningdb",
		})
	if err != nil {
		log.Fatal("failed to start server", err)
	}

	app := appgroup.NewAppGroup(appgroup.WithGracefullyQuit(syscall.SIGINT, syscall.SIGTERM))
	app.Add(httpServer)

	app.Start()
	if err := app.Wait(); err != nil {
		log.Fatal(err)
	}
}

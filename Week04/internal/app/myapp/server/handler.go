package server

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/schwarzeni/Go-000/Week04/internal/app/myapp/service"
)

func NewHandler(svc *service.Service) http.Handler {
	router := gin.Default()
	router.GET("/v1/api/articles", func(ctx *gin.Context) {
		articles, err := svc.GetArticles()
		if err != nil {
			log.Printf("%+v", err)
			ctx.Status(http.StatusInternalServerError)
			return
		}
		ctx.JSON(http.StatusOK, articles)
	})
	return router
}

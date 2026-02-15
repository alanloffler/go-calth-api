package business

import (
	"github.com/alanloffler/go-calth-api/internal/database/sqlc"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine, q *sqlc.Queries) {
	var repo *BusinessRepository = NewBusinessRepository(q)
	var handler *BusinessHandler = NewBusinessHandler(repo)
	var businesses *gin.RouterGroup = router.Group("/businesses")

	businesses.POST("", handler.Create)
	businesses.GET("", handler.GetAll)
	businesses.GET("/:id", handler.GetOneByID)
	businesses.PATCH("/:id", handler.Update)
}

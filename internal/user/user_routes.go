package user

import (
	"github.com/alanloffler/go-calth-api/internal/database/sqlc"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine, q *sqlc.Queries) {
	var repo *UserRepository = NewUserRepository(q)
	var handler *UserHandler = NewUserHandler(repo)
	var users *gin.RouterGroup = router.Group("/users")

	users.POST("", handler.Create)
	users.GET("", handler.GetAll)
	users.GET("/:id", handler.GetByID)
	users.PATCH("/:id", handler.Update)
	users.PATCH("/:id/restore", handler.Restore)
	users.DELETE("/:id/soft", handler.SoftDelete)
}

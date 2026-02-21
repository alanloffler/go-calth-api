package user

import (
	"github.com/alanloffler/go-calth-api/internal/database/sqlc"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.RouterGroup, q *sqlc.Queries) {
	var repo *UserRepository = NewUserRepository(q)
	var handler *UserHandler = NewUserHandler(repo)
	var users *gin.RouterGroup = router.Group("/users")

	users.POST("", handler.Create)
	users.GET("", handler.GetAll)
	users.GET("/soft", handler.GetAllWithSoftDeleted)
	users.GET("/role/:role/soft", handler.GetAllByRoleWithSoftDeleted)
	users.GET("/:id", handler.GetByID)
	users.GET("/:id/soft", handler.GetByIDWithSoftDeleted)
	users.PATCH("/:id", handler.Update)
	users.PATCH("/:id/restore", handler.Restore)
	users.DELETE("/:id", handler.Delete)
	users.DELETE("/:id/soft", handler.SoftDelete)
}

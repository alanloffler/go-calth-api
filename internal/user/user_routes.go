package user

import (
	"github.com/alanloffler/go-calth-api/internal/database/sqlc"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func RegisterRoutes(router *gin.RouterGroup, q *sqlc.Queries, pool *pgxpool.Pool) {
	var repo *UserRepository = NewUserRepository(q)
	var handler *UserHandler = NewUserHandler(repo, pool)
	var users *gin.RouterGroup = router.Group("/users")

	users.POST("", handler.Create)
	users.GET("", handler.GetAll)
	users.GET("/soft", handler.GetAllWithSoftDeleted)
	users.GET("/role/:role", handler.GetAllByRole)
	users.GET("/role/:role/soft", handler.GetAllByRoleWithSoftDeleted)
	users.GET("/:id", handler.GetByID)
	users.GET("/:id/soft", handler.GetByIDWithSoftDeleted)
	users.PATCH("/:id", handler.Update)
	users.PATCH("/:id/restore", handler.Restore)
	users.DELETE("/:id", handler.Delete)
	users.DELETE("/:id/soft", handler.SoftDelete)
	// Checks
	users.GET("/check/email/:email", handler.CheckEmailAvailability)
	users.GET("/check/ic/:ic", handler.CheckIcAvailability)
	users.GET("/check/username/:userName", handler.CheckUsernameAvailability)
}

package role

import (
	"github.com/alanloffler/go-calth-api/internal/database/sqlc"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine, q *sqlc.Queries) {
	var repo *RoleRepository = NewRoleRepository(q)
	var handler *RoleHandler = NewRoleHandler(repo)
	var roles *gin.RouterGroup = router.Group("/roles")

	roles.POST("", handler.Create)
	roles.GET("", handler.GetAll)
	roles.GET("/soft", handler.GetAllWithSoftDeleted)
	roles.GET("/:id", handler.GetOneByID)
	roles.GET("/:id/soft", handler.GetOneByIDWithSoftDeleted)
	roles.PATCH("/:id", handler.Update)
	roles.PATCH("/:id/restore", handler.Restore)
	roles.DELETE("/:id", handler.Delete)
	roles.DELETE("/:id/soft", handler.SoftDelete)
}

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
	roles.DELETE("/:id", handler.Delete)
}

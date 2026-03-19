package role

import (
	"github.com/alanloffler/go-calth-api/internal/database/sqlc"
	"github.com/alanloffler/go-calth-api/internal/middleware"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func RegisterRoutes(router *gin.RouterGroup, q *sqlc.Queries, pool *pgxpool.Pool) {
	var repo *RoleRepository = NewRoleRepository(q)
	var handler *RoleHandler = NewRoleHandler(repo, pool)
	var roles *gin.RouterGroup = router.Group("/roles")

	roles.POST("", middleware.PermissionMiddleware(q, "roles-create"), handler.Create)

	roles.GET("", middleware.PermissionMiddleware(q, "roles-view"), handler.GetAll)
	roles.GET("/soft", middleware.PermissionMiddleware(q, "roles-view"), handler.GetAllWithSoftDeleted)
	roles.GET("/:id", middleware.PermissionMiddleware(q, "roles-view"), handler.GetOneByID)
	roles.GET("/:id/soft", middleware.PermissionMiddleware(q, "roles-view"), handler.GetOneByIDWithSoftDeleted)

	roles.PATCH("/:id", middleware.PermissionMiddleware(q, "roles-update"), handler.Update)
	roles.PATCH("/:id/restore", middleware.PermissionMiddleware(q, "roles-restore"), handler.Restore)

	roles.DELETE("/:id", middleware.PermissionMiddleware(q, "roles-delete-hard"), handler.Delete)
	roles.DELETE("/:id/soft", middleware.PermissionMiddleware(q, "roles-delete"), handler.SoftDelete)
}

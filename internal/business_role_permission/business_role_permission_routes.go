package business_role_permission

import (
	"github.com/alanloffler/go-calth-api/internal/database/sqlc"
	"github.com/alanloffler/go-calth-api/internal/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.RouterGroup, q *sqlc.Queries) {
	repo := NewBusinessRolePermissionRepository(q)
	handler := NewBusinessRolePermissionHandler(repo)
	overrides := router.Group("/role-overrides")

	overrides.GET("/:roleId", middleware.PermissionMiddleware(q, "roles-view"), handler.ListEffective)
	overrides.GET("/:roleId/overrides", middleware.PermissionMiddleware(q, "roles-view"), handler.ListOverrides)

	overrides.PUT("/:roleId/permissions/:permissionId", middleware.PermissionMiddleware(q, "roles-update"), handler.Upsert)

	overrides.DELETE("/:roleId/permissions/:permissionId", middleware.PermissionMiddleware(q, "roles-update"), handler.ResetOne)
	overrides.DELETE("/:roleId", middleware.PermissionMiddleware(q, "roles-update"), handler.ResetAll)
}

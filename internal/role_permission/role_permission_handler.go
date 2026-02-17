package role_permission

import (
	"net/http"

	"github.com/alanloffler/go-calth-api/internal/common/response"
	"github.com/alanloffler/go-calth-api/internal/database/sqlc"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type RolePermissionHandler struct {
	repo *RolePermissionRepository
}

func NewRolePermissionHandler(repo *RolePermissionRepository) *RolePermissionHandler {
	return &RolePermissionHandler{repo: repo}
}

type CreateRolePermissionRequest struct {
	RoleID       string `json:"roleId" binding:"required,uuid"`
	PermissionID string `json:"permissionId" binding:"required,uuid"`
}

func (h *RolePermissionHandler) Create(c *gin.Context) {
	var req CreateRolePermissionRequest
	if err := c.ShouldBindBodyWithJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, "Error al crear permiso de rol", err))
		return
	}

	roleID, err := uuid.Parse(req.RoleID)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, "Formato de UUID inválido", err))
		return
	}

	permissionID, err := uuid.Parse(req.PermissionID)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, "Formato de UUID inválido", err))
		return
	}

	params := sqlc.CreateRolePermissionParams{
		RoleID:       pgtype.UUID{Bytes: roleID, Valid: true},
		PermissionID: pgtype.UUID{Bytes: permissionID, Valid: true},
	}

	result, err := h.repo.Create(c.Request.Context(), params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(http.StatusInternalServerError, "Error al crear permiso de rol", err))
		return
	}

	c.JSON(http.StatusOK, response.Success("Permiso de rol creado", &result))
}

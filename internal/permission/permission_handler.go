package permission

import (
	"net/http"

	"github.com/alanloffler/go-calth-api/internal/common/response"
	"github.com/alanloffler/go-calth-api/internal/database/sqlc"
	"github.com/gin-gonic/gin"
)

type PermissionHandler struct {
	repo *PermissionRepository
}

func NewPermissionHandler(repo *PermissionRepository) *PermissionHandler {
	return &PermissionHandler{repo: repo}
}

type CreatePermissionRequest struct {
	Name        string `json:"name" binding:"required,min=3,max=100"`
	Category    string `json:"category" binding:"required,min=3,max=100"`
	ActionKey   string `json:"actionKey" binding:"required,min=3,max=100"`
	Description string `json:"description" binding:"required,min=3,max=255"`
}

func (h *PermissionHandler) Create(c *gin.Context) {
	var req CreatePermissionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, "Error al crear permiso", err))
		return
	}

	permission, err := h.repo.Create(c.Request.Context(), sqlc.CreatePermissionParams{
		Name:        req.Name,
		Category:    req.Category,
		ActionKey:   req.ActionKey,
		Description: req.Description,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(http.StatusInternalServerError, "Error al crear permiso", err))
		return
	}

	c.JSON(http.StatusOK, response.Created("Permiso creado", &permission))
}

func (h *PermissionHandler) GetAll(c *gin.Context) {
	permissions, err := h.repo.GetAll(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusNotFound, response.Error(http.StatusNotFound, "Permisos no encontrados", err))
		return
	}

	c.JSON(http.StatusOK, response.Success("Permisos encontrados", &permissions))
}

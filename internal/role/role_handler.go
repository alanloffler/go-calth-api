package role

import (
	"net/http"

	"github.com/alanloffler/go-calth-api/internal/common/response"
	"github.com/alanloffler/go-calth-api/internal/database/sqlc"
	"github.com/gin-gonic/gin"
)

type RoleHandler struct {
	repo *RoleRepository
}

func NewRoleHandler(repo *RoleRepository) *RoleHandler {
	return &RoleHandler{repo: repo}
}

type CreateRoleRequest struct {
	Name        string `json:"name" binding:"required,min=3,max=100"`
	Value       string `json:"value" binding:"required,min=3,max=100"`
	Description string `json:"description" binding:"required,min=3,max=100"`
}

func (h *RoleHandler) Create(c *gin.Context) {
	var req CreateRoleRequest
	if err := c.ShouldBindBodyWithJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, "Error de validación de datos", err))
		return
	}

	role, err := h.repo.Create(c.Request.Context(), sqlc.CreateRoleParams{
		Name:        req.Name,
		Value:       req.Value,
		Description: req.Description,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(http.StatusInternalServerError, "Error al crear rol", err))
		return
	}

	c.JSON(http.StatusOK, response.Success("Rol creado", &role))
}

package role

import (
	"net/http"

	"github.com/alanloffler/go-calth-api/internal/common/response"
	"github.com/alanloffler/go-calth-api/internal/database/sqlc"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
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

func (h *RoleHandler) GetAll(c *gin.Context) {
	permissions, err := h.repo.GetAll(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusNotFound, response.Error(http.StatusNotFound, "Roles no encontrados", err))
		return
	}

	c.JSON(http.StatusOK, response.Success("Roles encontrados", &permissions))
}

func (h *RoleHandler) Delete(c *gin.Context) {
	var id pgtype.UUID
	if err := id.Scan(c.Param("id")); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, "Formato de ID inválido", err))
		return
	}

	rows, err := h.repo.Delete(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(http.StatusInternalServerError, "Error al eliminar rol", err))
		return
	}

	if rows == 0 {
		c.JSON(http.StatusNotFound, response.Error(http.StatusNotFound, "Rol no encontrado"))
		return
	}

	c.JSON(http.StatusOK, response.Success[any]("Rol eliminado", nil))
}

func (h *RoleHandler) SoftDelete(c *gin.Context) {
	var id pgtype.UUID
	if err := id.Scan(c.Param("id")); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, "Formato de ID inválido", err))
		return
	}

	role, err := h.repo.SoftDelete(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, "Error al eliminar rol"))
		return
	}

	c.JSON(http.StatusOK, response.Success("Rol eliminado", &role))
}

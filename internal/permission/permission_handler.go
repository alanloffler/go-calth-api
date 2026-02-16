package permission

import (
	"net/http"

	"github.com/alanloffler/go-calth-api/internal/common/response"
	"github.com/alanloffler/go-calth-api/internal/common/utils"
	"github.com/alanloffler/go-calth-api/internal/database/sqlc"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
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

type UpdatePermissionRequest struct {
	Name        *string `json:"name" binding:"omitempty,min=3,max=50"`
	Category    *string `json:"category" binding:"omitempty,min=3,max=100"`
	ActionKey   *string `json:"actionKey" binding:"omitempty,min=3,max=100"`
	Description *string `json:"description" binding:"omitempty,min=3,max=255"`
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

func (h *PermissionHandler) GetAllWithSoftDeleted(c *gin.Context) {
	permissions, err := h.repo.GetAllWithSoftDeleted(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusNotFound, response.Error(http.StatusNotFound, "Permisos no encontrados", err))
		return
	}

	c.JSON(http.StatusOK, response.Success("Permisos encontrados", &permissions))
}

func (h *PermissionHandler) GetOneByID(c *gin.Context) {
	var id pgtype.UUID
	if err := id.Scan(c.Param("id")); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, "Formato de ID inválido", err))
		return
	}

	permission, err := h.repo.GetOneByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, response.Error(http.StatusNotFound, "Permiso no encontrado", err))
		return
	}

	c.JSON(http.StatusOK, response.Success("Permiso encontrado", &permission))
}

func (h *PermissionHandler) Update(c *gin.Context) {
	var id pgtype.UUID
	if err := id.Scan(c.Param("id")); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, "Formato de ID inválido", err))
		return
	}

	var req UpdatePermissionRequest
	if err := c.ShouldBindBodyWithJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, "Error al actualizar permiso", err))
		return
	}

	permission, err := h.repo.Update(c.Request.Context(), sqlc.UpdatePermissionParams{
		ID:          id,
		Name:        utils.ToPgText(req.Name),
		Category:    utils.ToPgText(req.Category),
		ActionKey:   utils.ToPgText(req.ActionKey),
		Description: utils.ToPgText(req.Description),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(http.StatusInternalServerError, "Error al actualizar permiso", err))
		return
	}

	c.JSON(http.StatusOK, response.Success("Permiso actualizado", &permission))
}

func (h *PermissionHandler) Delete(c *gin.Context) {
	var id pgtype.UUID
	if err := id.Scan(c.Param("id")); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, "Formato de ID inválido", err))
		return
	}

	err := h.repo.Delete(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(http.StatusInternalServerError, "Error al eliminar permiso", err))
		return
	}

	c.JSON(http.StatusOK, response.Success[any]("Permiso eliminado", nil))
}

func (h *PermissionHandler) SoftDelete(c *gin.Context) {
	var id pgtype.UUID
	if err := id.Scan(c.Param("id")); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, "Formato de ID inválido", err))
		return
	}

	permission, err := h.repo.SoftDelete(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(http.StatusInternalServerError, "Error al eliminar permiso", err))
		return
	}

	c.JSON(http.StatusOK, response.Success("Permiso eliminado", &permission))
}

func (h *PermissionHandler) Restore(c *gin.Context) {
	var id pgtype.UUID
	if err := id.Scan(c.Param("id")); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, "Formato de ID inválido", err))
		return
	}

	permission, err := h.repo.Restore(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(http.StatusInternalServerError, "Error al restaurar permiso", err))
		return
	}

	c.JSON(http.StatusOK, response.Success("Permiso restaurado", &permission))
}

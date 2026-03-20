package setting

import (
	"net/http"

	"github.com/alanloffler/go-calth-api/internal/common/response"
	"github.com/alanloffler/go-calth-api/internal/common/utils"
	"github.com/alanloffler/go-calth-api/internal/database/sqlc"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
)

type SettingHandler struct {
	repo *SettingRepository
}

func NewSettingHandler(repo *SettingRepository) *SettingHandler {
	return &SettingHandler{repo: repo}
}

type UpdateSettingRequest struct {
	Module    *string `json:"module" binding:"omitempty,min=2,max=50"`
	Submodule *string `json:"submodule" binding:"omitempty,min=2,max=50"`
	Key       *string `json:"key" binding:"omitempty,min=2,max=50"`
	Value     *string `json:"value" binding:"omitempty,min=2,max=255"`
	Title     *string `json:"title" binding:"omitempty,min=2,max=100"`
}

func (h *SettingHandler) GetAll(c *gin.Context) {
	settings, err := h.repo.GetAll(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(http.StatusInternalServerError, "Error al obtener configuraciones", err))
		return
	}

	c.JSON(http.StatusOK, response.Success("Configuraciones encontradas", &settings))
}

func (h *SettingHandler) GetByModule(c *gin.Context) {
	module := c.Param("module")
	if module == "" {
		c.JSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, "Módulo requerido"))
		return
	}

	result, err := h.repo.GetByModule(c.Request.Context(), module)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(http.StatusInternalServerError, "Error al obtener configuraciones", err))
		return
	}

	c.JSON(http.StatusOK, response.Success("Configuraciones encontradas", &result))
}

func (h *SettingHandler) Update(c *gin.Context) {
	var id pgtype.UUID
	if err := id.Scan(c.Param("id")); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, "Formato de ID inválido", err))
		return
	}

	var req UpdateSettingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, "Error de validación de datos", err))
		return
	}

	setting, err := h.repo.Update(c.Request.Context(), sqlc.UpdateSettingParams{
		ID:        id,
		Module:    utils.ToPgText(req.Module),
		Submodule: utils.ToPgText(req.Submodule),
		Key:       utils.ToPgText(req.Key),
		Value:     utils.ToPgText(req.Value),
		Title:     utils.ToPgText(req.Title),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(http.StatusInternalServerError, "Error al actualizar configuraciones", err))
		return
	}

	c.JSON(http.StatusOK, response.Success("Configuraciones actualizadas", &setting))
}

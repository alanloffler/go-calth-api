package blocked_day

import (
	"net/http"
	"time"

	"github.com/alanloffler/go-calth-api/internal/common/ctxkeys"
	"github.com/alanloffler/go-calth-api/internal/common/response"
	"github.com/alanloffler/go-calth-api/internal/database/sqlc"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
)

type BlockedDayHandler struct {
	repo *BlockedDayRepository
}

func NewBlockedDayHandler(repo *BlockedDayRepository) *BlockedDayHandler {
	return &BlockedDayHandler{repo: repo}
}

type CreateBlockedDayRequest struct {
	Date           string `json:"date" binding:"required,datetime=2006-01-02T15:04:05Z07:00"`
	Reason         string `json:"reason" binding:"required,min=3,max=50"`
	ProfessionalID string `json:"professionalId" binding:"required,uuid"`
}

func (h *BlockedDayHandler) Create(c *gin.Context) {
	businessID, ok := ctxkeys.BusinessID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, response.Error(http.StatusUnauthorized, "Usuario no autenticado"))
		return
	}

	var req CreateBlockedDayRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, "Error de validación de datos", err))
		return
	}

	date, err := time.Parse(time.RFC3339, req.Date)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, "Formato de fecha inválido", err))
		return
	}

	var professionalID pgtype.UUID
	if err := professionalID.Scan(req.ProfessionalID); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, "Formato de ID del profesional inválido", err))
		return
	}

	bd, err := h.repo.Create(c.Request.Context(), sqlc.CreateBlockedDayParams{
		Date:           pgtype.Timestamptz{Time: date, Valid: true},
		Reason:         req.Reason,
		BusinessID:     businessID,
		ProfessionalID: professionalID,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(http.StatusInternalServerError, "Error al crear el día bloqueado", err))
		return
	}

	c.JSON(http.StatusOK, response.Created("Día bloqueado creado", &bd))
}

func (h *BlockedDayHandler) GetByProfessionalID(c *gin.Context) {
	businessID, ok := ctxkeys.BusinessID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, response.Error(http.StatusUnauthorized, "Usuario no autenticado"))
		return
	}

	var professionalID pgtype.UUID
	if err := professionalID.Scan(c.Param("professionalId")); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, "Formato de ID del profesional inválido", err))
		return
	}

	bd, err := h.repo.GetByProfessionalID(c.Request.Context(), sqlc.GetBlockedDaysProfessionalIDParams{
		BusinessID:     businessID,
		ProfessionalID: professionalID,
	})
	if err != nil {
		c.JSON(http.StatusNotFound, response.Error(http.StatusNotFound, "Dias bloqueados no encontrados", err))
		return
	}

	c.JSON(http.StatusOK, response.Success("Dias bloqueados encontrados", &bd))
}

func (h *BlockedDayHandler) Delete(c *gin.Context) {
	businessID, ok := ctxkeys.BusinessID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, response.Error(http.StatusUnauthorized, "Usuario no autenticado"))
		return
	}

	var id pgtype.UUID
	if err := id.Scan(c.Param("id")); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, "Formato de ID inválido", err))
		return
	}

	affected, err := h.repo.Delete(c.Request.Context(), sqlc.DeleteBlockedDayParams{
		BusinessID: businessID,
		ID:         id,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(http.StatusInternalServerError, "Error eliminando día bloqueado", err))
		return
	}
	if affected == 0 {
		c.JSON(http.StatusNotFound, response.Error(http.StatusNotFound, "Día bloqueado no encontrado"))
		return
	}

	c.JSON(http.StatusOK, response.Success[any]("Día bloqueado eliminado", nil))
}

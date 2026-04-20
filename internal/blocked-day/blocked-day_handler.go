package blocked_day

import (
	"errors"
	"net/http"
	"time"

	"github.com/alanloffler/go-calth-api/internal/common/ctxkeys"
	"github.com/alanloffler/go-calth-api/internal/common/response"
	"github.com/alanloffler/go-calth-api/internal/common/utils"
	"github.com/alanloffler/go-calth-api/internal/database/sqlc"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgconn"
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
	Recurrent      bool   `json:"recurrent"`
}

type UpdateBlockedDayRequest struct {
	Date   string `json:"date" binding:"omitempty,datetime=2006-01-02T15:04:05Z07:00"`
	Reason string `json:"reason" binding:"omitempty,min=3,max=50"`
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
		Recurrent:      pgtype.Bool{Bool: req.Recurrent, Valid: true},
	})
	if err != nil {
		if pgErr, ok := errors.AsType[*pgconn.PgError](err); ok {
			if pgErr.Code == "23505" {
				c.JSON(http.StatusConflict, response.Error(http.StatusConflict, "El día bloqueado ya existe"))
				return
			}
		}

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
		c.JSON(http.StatusNotFound, response.Error(http.StatusNotFound, "Días bloqueados no encontrados", err))
		return
	}

	c.JSON(http.StatusOK, response.Success("Días bloqueados encontrados", &bd))
}

func (h *BlockedDayHandler) Update(c *gin.Context) {
	businessID, ok := ctxkeys.BusinessID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, response.Error(http.StatusUnauthorized, "Usuario no autorizado"))
		return
	}

	var id pgtype.UUID
	if err := id.Scan(c.Param("id")); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, "Formato de ID inválido", err))
		return
	}

	var req UpdateBlockedDayRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, "Error de validación de datos", err))
		return
	}

	params := sqlc.UpdateBlockedDayParams{
		BusinessID: businessID,
		ID:         id,
		Reason:     utils.ToPgText(&req.Reason),
	}

	if req.Date != "" {
		date, err := time.Parse(time.RFC3339, req.Date)
		if err != nil {
			c.JSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, "Formato de fecha inválido", err))
			return
		}
		params.Date = pgtype.Timestamptz{Time: date, Valid: true}
	}

	affected, err := h.repo.Update(c.Request.Context(), params)
	if err != nil {
		if pgErr, ok := errors.AsType[*pgconn.PgError](err); ok {
			if pgErr.Code == "23505" {
				c.JSON(http.StatusConflict, response.Error(http.StatusConflict, "El día bloqueado ya existe, elige otra fecha"))
				return
			}
		}

		c.JSON(http.StatusInternalServerError, response.Error(http.StatusInternalServerError, "Error al actualizar el día bloqueado", err))
		return
	}
	if affected == 0 {
		c.JSON(http.StatusNotFound, response.Error(http.StatusNotFound, "Día bloqueado no encontrado"))
		return
	}

	c.JSON(http.StatusOK, response.Success[any]("Día bloqueado actualizado", nil))
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

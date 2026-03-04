package medical_history

import (
	"net/http"
	"time"

	"github.com/alanloffler/go-calth-api/internal/common/ctxkeys"
	"github.com/alanloffler/go-calth-api/internal/common/response"
	"github.com/alanloffler/go-calth-api/internal/database/sqlc"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
)

type MedicalHistoryHandler struct {
	repo *MedicalHistoryRepository
}

type CreateMedicalHistoryRequest struct {
	UserID         string `json:"userId" binding:"required,uuid"`
	ProfessionalID string `json:"professionalId" binding:"required,uuid"`
	EventID        string `json:"eventId" binding:"omitempty,uuid"`
	Date           string `json:"date" binding:"required,datetime=2006-01-02T15:04:05Z07:00"`
	Reason         string `json:"reason" binding:"required,min=3,max=100"`
	Recipe         bool   `json:"recipe"`
	Comments       string `json:"comments" binding:"required,min=3"`
}

func NewMedicalHistoryHandler(repo *MedicalHistoryRepository) *MedicalHistoryHandler {
	return &MedicalHistoryHandler{repo: repo}
}

func (h *MedicalHistoryHandler) CreateMedicalHistory(c *gin.Context) {
	businessID, ok := ctxkeys.BusinessID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, response.Error(http.StatusUnauthorized, "Usuario no autenticado"))
		return
	}

	var req CreateMedicalHistoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, "Error de validación de datos", err))
		return
	}

	date, err := time.Parse(time.RFC3339, req.Date)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, "Formato de fecha inválido", err))
		return
	}

	var userID pgtype.UUID
	if err := userID.Scan(req.UserID); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, "Formato de ID del usuario inválido", err))
		return
	}

	var professionalID pgtype.UUID
	if err := professionalID.Scan(req.ProfessionalID); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, "Formato de ID del profesional inválido", err))
		return
	}

	var eventID pgtype.UUID
	if req.EventID != "" {
		if err := eventID.Scan(req.EventID); err != nil {
			c.JSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, "Formato de ID del evento inválido", err))
			return
		}
	}

	mh, err := h.repo.CreateMedicaHistory(c.Request.Context(), sqlc.CreateMedicalHistoryParams{
		BusinessID:     businessID,
		UserID:         userID,
		ProfessionalID: professionalID,
		EventID:        eventID,
		Date:           pgtype.Timestamptz{Time: date, Valid: true},
		Reason:         req.Reason,
		Recipe:         req.Recipe,
		Comments:       req.Comments,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(http.StatusInternalServerError, "Error al crear el historial médico", err))
		return
	}

	c.JSON(http.StatusOK, response.Created("Historia médica creada", &mh))
}

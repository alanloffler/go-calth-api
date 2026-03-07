package medical_history

import (
	"net/http"
	"time"

	"github.com/alanloffler/go-calth-api/internal/common/ctxkeys"
	"github.com/alanloffler/go-calth-api/internal/common/response"
	"github.com/alanloffler/go-calth-api/internal/database/sqlc"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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

type MedicalHistoryResponse struct {
	ID             string               `json:"id"`
	BusinessID     string               `json:"businessId"`
	UserID         string               `json:"userId"`
	ProfessionalID string               `json:"professionalId"`
	EventID        *string              `json:"eventId"`
	Date           string               `json:"date"`
	Reason         string               `json:"reason"`
	Recipe         bool                 `json:"recipe"`
	Comments       string               `json:"comments"`
	CreatedAt      string               `json:"createdAt"`
	UpdatedAt      string               `json:"updatedAt"`
	DeletedAt      *string              `json:"deletedAt"`
	User           UserResponse         `json:"user"`
	Professional   ProfessionalResponse `json:"professional"`
}

type UserResponse struct {
	IC        string `json:"ic"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

type ProfessionalResponse struct {
	FirstName string                      `json:"firstName"`
	LastName  string                      `json:"lastName"`
	Profile   ProfessionalProfileResponse `json:"professionalProfile"`
}

type ProfessionalProfileResponse struct {
	ProfessionalPrefix string `json:"professionalPrefix"`
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

func (h *MedicalHistoryHandler) GetAllByPatientIDWithSoftDeleted(c *gin.Context) {
	businessID, ok := ctxkeys.BusinessID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, response.Error(http.StatusUnauthorized, "Usuario no autenticado"))
		return
	}

	var id pgtype.UUID
	if err := id.Scan(c.Param("id")); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, "Formato de ID del usuario inválido", err))
		return
	}

	mhs, err := h.repo.GetAllByPatientIDWithSoftDeleted(c.Request.Context(), sqlc.GetMedicalHistoriesByPatientIDWithSoftDeletedParams{BusinessID: businessID, UserID: id})
	if err != nil {
		c.JSON(http.StatusNotFound, response.Error(http.StatusNotFound, "Historias médicas no encontradas", err))
		return
	}

	result := make([]MedicalHistoryResponse, len(mhs))
	for i, mh := range mhs {
		var eventID *string
		if mh.EventID.Valid {
			s := uuid.UUID(mh.EventID.Bytes).String()
			eventID = &s
		}

		var deletedAt *string
		if mh.DeletedAt.Valid {
			s := mh.DeletedAt.Time.Format(time.RFC3339)
			deletedAt = &s
		}

		result[i] = MedicalHistoryResponse{
			ID:             uuid.UUID(mh.ID.Bytes).String(),
			BusinessID:     uuid.UUID(mh.BusinessID.Bytes).String(),
			UserID:         uuid.UUID(mh.UserID.Bytes).String(),
			ProfessionalID: uuid.UUID(mh.ProfessionalID.Bytes).String(),
			EventID:        eventID,
			Date:           mh.Date.Time.Format(time.RFC3339),
			Reason:         mh.Reason,
			Recipe:         mh.Recipe,
			Comments:       mh.Comments,
			User: UserResponse{
				IC:        mh.Ic.String,
				FirstName: mh.FirstName.String,
				LastName:  mh.LastName.String,
			},
			Professional: ProfessionalResponse{
				FirstName: mh.FirstName_2.String,
				LastName:  mh.LastName_2.String,
				Profile: ProfessionalProfileResponse{
					ProfessionalPrefix: mh.ProfessionalPrefix.String,
				},
			},
			CreatedAt: mh.CreatedAt.Time.Format(time.RFC3339),
			UpdatedAt: mh.UpdatedAt.Time.Format(time.RFC3339),
			DeletedAt: deletedAt,
		}
	}

	c.JSON(http.StatusOK, response.Success("Historias médicas encontradas", &result))
}

func (h *MedicalHistoryHandler) SoftDelete(c *gin.Context) {
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

	rows, err := h.repo.SoftDelete(c.Request.Context(), sqlc.SoftDeleteMedicalHistoryParams{
		BusinessID: businessID,
		ID:         id,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(http.StatusInternalServerError, "Error al eliminar historia médica", err))
		return
	}
	if rows == 0 {
		c.JSON(http.StatusNotFound, response.Error(http.StatusNotFound, "Historia médica no encontrada"))
		return
	}

	c.JSON(http.StatusOK, response.Success[any]("Historia médica eliminada", nil))
}

func (h *MedicalHistoryHandler) Restore(c *gin.Context) {
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

	rows, err := h.repo.Restore(c.Request.Context(), sqlc.RestoreMedicalHistoryParams{
		BusinessID: businessID,
		ID:         id,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(http.StatusInternalServerError, "Error al restaurar la historia médica", err))
		return
	}
	if rows == 0 {
		c.JSON(http.StatusNotFound, response.Error(http.StatusNotFound, "Historia médica no encontrada"))
		return
	}

	c.JSON(http.StatusOK, response.Success[any]("Historia médica restaurada", nil))
}

func (h *MedicalHistoryHandler) Delete(c *gin.Context) {
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

	err := h.repo.Delete(c.Request.Context(), sqlc.DeleteMedicalHistoryParams{
		BusinessID: businessID,
		ID:         id,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(http.StatusInternalServerError, "Error eliminando historia médica", err))
		return
	}

	c.JSON(http.StatusOK, response.Success[any]("Historia médica eliminada", nil))
}

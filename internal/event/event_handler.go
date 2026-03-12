package event

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/alanloffler/go-calth-api/internal/common/ctxkeys"
	"github.com/alanloffler/go-calth-api/internal/common/response"
	"github.com/alanloffler/go-calth-api/internal/database/sqlc"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
)

type EventHandler struct {
	repo *EventRepository
}

type CreateEventRequest struct {
	Title          string `json:"title" binding:"required,min=3,max=255"`
	StartDate      string `json:"startDate" binding:"required,datetime=2006-01-02T15:04:05Z07:00"`
	EndDate        string `json:"endDate" binding:"required,datetime=2006-01-02T15:04:05Z07:00"`
	ProfessionalID string `json:"professionalId" binding:"required,uuid"`
	UserID         string `json:"userId" binding:"required,uuid"`
}

type UpdateEventRequest struct {
	Title          *string `json:"title" binding:"omitempty,min=3,max=255"`
	StartDate      *string `json:"startDate" binding:"omitempty,datetime=2006-01-02T15:04:05Z07:00"`
	EndDate        *string `json:"endDate" binding:"omitempty,datetime=2006-01-02T15:04:05Z07:00"`
	ProfessionalID *string `json:"professionalId" binding:"omitempty,uuid"`
	UserID         *string `json:"userId" binding:"omitempty,uuid"`
	Status         *string `json:"status" binding:"omitempty"`
}

type UpdateEventStatusRequest struct {
	Status string `json:"status" binding:"required"`
}

func NewEventHandler(repo *EventRepository) *EventHandler {
	return &EventHandler{repo: repo}
}

func (h *EventHandler) Create(c *gin.Context) {
	businessID, ok := ctxkeys.BusinessID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, response.Error(http.StatusUnauthorized, "Usuario no autenticado"))
		return
	}

	var req CreateEventRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, "Error de validación de datos", err))
		return
	}

	startTime, err := time.Parse(time.RFC3339, req.StartDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, "Formato de fecha de inicio inválido", err))
		return
	}

	endTime, err := time.Parse(time.RFC3339, req.EndDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, "Formato de fecha de finalización inválido", err))
		return
	}

	startDate := pgtype.Timestamptz{Time: startTime, Valid: true}
	endDate := pgtype.Timestamptz{Time: endTime, Valid: true}

	var professionalID pgtype.UUID
	if err := professionalID.Scan(req.ProfessionalID); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, "Formato de ID del profesional inválido", err))
		return
	}

	var userID pgtype.UUID
	if err := userID.Scan(req.UserID); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, "Formato de ID del usuario inválido", err))
		return
	}

	event, err := h.repo.Create(c.Request.Context(), sqlc.CreateEventParams{
		Title:          req.Title,
		StartDate:      startDate,
		EndDate:        endDate,
		BusinessID:     businessID,
		ProfessionalID: professionalID,
		UserID:         userID,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(http.StatusInternalServerError, "Error al crear evento", err))
		return
	}

	c.JSON(http.StatusOK, response.Created("Evento creado", &event))
}

func (h *EventHandler) GetByBusinessID(c *gin.Context) {
	businessID, ok := ctxkeys.BusinessID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, response.Error(http.StatusUnauthorized, "Usuario no autenticado"))
		return
	}

	limit, err := strconv.ParseInt(c.Query("limit"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, "Límite inválido", err))
		return
	}

	rawEvents, err := h.repo.GetEventsByBusinessID(c.Request.Context(), sqlc.GetEventsByBusinessIDParams{BusinessID: businessID, Limit: int32(limit)})
	if err != nil {
		c.JSON(http.StatusNotFound, response.Error(http.StatusNotFound, "Eventos no encontrados", err))
		return
	}

	events := make([]json.RawMessage, len(rawEvents))
	for i, e := range rawEvents {
		events[i] = json.RawMessage(e)
	}

	c.JSON(http.StatusOK, response.Success("Eventos encontrados", &events))
}

func (h *EventHandler) GetByProfessionalID(c *gin.Context) {
	businessID, ok := ctxkeys.BusinessID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, response.Error(http.StatusUnauthorized, "Usuario no autenticado"))
		return
	}

	var id pgtype.UUID
	if err := id.Scan(c.Param("id")); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, "Formato de ID de professional inválido", err))
		return
	}

	rawEvents, err := h.repo.GetEventsByProfessionalID(c.Request.Context(), sqlc.GetEventsByProfessionalIDParams{
		BusinessID:     businessID,
		ProfessionalID: id,
	})
	if err != nil {
		c.JSON(http.StatusNotFound, response.Error(http.StatusNotFound, "Eventos no encontrados", err))
		return
	}

	events := make([]json.RawMessage, len(rawEvents))
	for i, e := range rawEvents {
		events[i] = json.RawMessage(e)
	}

	c.JSON(http.StatusOK, response.Success("Eventos encontrados", &events))
}

func (h *EventHandler) GetEventsByBusinessProfessionalPatient(c *gin.Context) {
	businessID, ok := ctxkeys.BusinessID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, response.Error(http.StatusUnauthorized, "Usuario no autenticado"))
		return
	}

	var professionalID pgtype.UUID
	if err := professionalID.Scan(c.Query("professional")); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, "Formato de ID de profesional inválido", err))
		return
	}

	var userID pgtype.UUID
	if err := userID.Scan(c.Param("patient_id")); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, "Formato de ID de paciente inválido", err))
		return
	}

	log.Println(businessID)
	log.Println(professionalID)
	log.Println(userID)

	rawEvents, err := h.repo.GetEventsByBusinessProfessionalPatient(c.Request.Context(), sqlc.GetEventsByBusinessProfessionalPatientParams{
		BusinessID:     businessID,
		ProfessionalID: professionalID,
		UserID:         userID,
	})
	if err != nil {
		c.JSON(http.StatusNotFound, response.Error(http.StatusNotFound, "Eventos no encontrados", err))
		return
	}

	events := make([]json.RawMessage, len(rawEvents))
	for i, e := range rawEvents {
		events[i] = json.RawMessage(e)
	}

	c.JSON(http.StatusOK, response.Success("Eventos encontrados", &events))
}

func (h *EventHandler) GetProfessionalEventsByDay(c *gin.Context) {
	businessID, ok := ctxkeys.BusinessID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, response.Error(http.StatusUnauthorized, "Usuario no autenticado"))
		return
	}

	var id pgtype.UUID
	if err := id.Scan(c.Param("id")); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, "Formato de ID profesional inválido", err))
		return
	}

	loc, err := time.LoadLocation("America/Argentina/Buenos_Aires")
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(http.StatusInternalServerError, "Error de zona horaria", err))
		return
	}

	dayStr := c.Param("day")
	dayTime, err := time.ParseInLocation("2006-01-02", dayStr, loc)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, "Formato de día inválido", err))
		return
	}

	startOfDay := pgtype.Timestamptz{Time: dayTime, Valid: true}
	endOfDay := pgtype.Timestamptz{Time: dayTime.Add(24*time.Hour - time.Second), Valid: true}

	rawEvents, err := h.repo.GetProfessionalEventsByDay(c.Request.Context(), sqlc.GetProfessionalEventsByDayParams{
		BusinessID:     businessID,
		ProfessionalID: id,
		StartDate:      startOfDay,
		StartDate_2:    endOfDay,
	})
	if err != nil {
		c.JSON(http.StatusNotFound, response.Error(http.StatusNotFound, "Eventos no encontrados", err))
		return
	}

	events := make([]json.RawMessage, len(rawEvents))
	for i, e := range rawEvents {
		events[i] = json.RawMessage(e)
	}

	c.JSON(http.StatusOK, response.Success("Eventos encontrados", &events))
}

func (h *EventHandler) GetProfessionalEventsByDayArray(c *gin.Context) {
	businessID, ok := ctxkeys.BusinessID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, response.Error(http.StatusUnauthorized, "Usuario no autenticado"))
		return
	}

	var id pgtype.UUID
	if err := id.Scan(c.Param("id")); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, "Formato de ID de profesional inválido", err))
		return
	}

	loc, err := time.LoadLocation("America/Argentina/Buenos_Aires")
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(http.StatusInternalServerError, "Error de zona horaria", err))
		return
	}

	dayStr := c.Param("day")
	dayTime, err := time.ParseInLocation("2006-01-02", dayStr, loc)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, "Formato de día inválido", err))
		return
	}

	startOfDay := pgtype.Timestamptz{Time: dayTime, Valid: true}
	endOfDay := pgtype.Timestamptz{Time: dayTime.Add(24*time.Hour - time.Second), Valid: true}

	slotTimes, err := h.repo.GetProfessionalEventsByDayArray(c.Request.Context(), sqlc.GetProfessionalEventsByDayArrayParams{
		BusinessID:     businessID,
		ProfessionalID: id,
		StartDate:      startOfDay,
		StartDate_2:    endOfDay,
	})
	if err != nil {
		c.JSON(http.StatusNotFound, response.Error(http.StatusNotFound, "Eventos no encontrados", err))
		return
	}

	dates := make([]string, len(slotTimes))
	for i, st := range slotTimes {
		dates[i] = st.Time.In(loc).Format("15:04")
	}

	c.JSON(http.StatusOK, response.Success("Fechas encontradas", &dates))
}

func (h *EventHandler) GetEventsFiltered(c *gin.Context) {
	businessID, ok := ctxkeys.BusinessID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, response.Error(http.StatusUnauthorized, "Usuario no autenticado"))
		return
	}

	params := sqlc.GetEventsFilteredParams{
		BusinessID: businessID,
	}

	limit := int32(10)
	if limitStr := c.Query("limit"); limitStr != "" {
		parsedLimit, err := strconv.ParseInt(limitStr, 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, "Límite inválido", err))
			return
		}
		limit = int32(parsedLimit)
	}
	params.QueryLimit = limit

	pageIndex := int32(1)
	if pageStr := c.Query("page"); pageStr != "" {
		parsedPage, err := strconv.ParseInt(pageStr, 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, "Página inválida", err))
			return
		}
		pageIndex = int32(parsedPage)
	}
	params.QueryOffset = (pageIndex - 1) * limit

	sortByMapping := map[string]string{
		"startDate":              "start_date",
		"status":                 "status",
		"title":                  "title",
		"user_firstName":         "user.first_name",
		"professional_firstName": "professional.first_name",
	}

	if sortBy := c.Query("sortBy"); sortBy != "" {
		if mapped, ok := sortByMapping[sortBy]; ok {
			params.SortBy = pgtype.Text{String: mapped, Valid: true}
		} else {
			c.JSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, "sortBy inválido"))
			return
		}
	}

	if sortOrder := c.Query("sortOrder"); sortOrder != "" {
		sortOrder = strings.ToLower(sortOrder)
		if sortOrder != "asc" && sortOrder != "desc" {
			c.JSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, "sortOrder debe ser ASC o DESC"))
			return
		}
		params.SortOrder = pgtype.Text{String: sortOrder, Valid: true}
	}

	if professionalIDStr := c.Query("professionalId"); professionalIDStr != "" {
		var professionalID pgtype.UUID
		if err := professionalID.Scan(professionalIDStr); err != nil {
			c.JSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, "Formato de ID de profesional inválido", err))
			return
		}
		params.ProfessionalID = professionalID
	}

	if patientIDStr := c.Query("patientId"); patientIDStr != "" {
		var patientID pgtype.UUID
		if err := patientID.Scan(patientIDStr); err != nil {
			c.JSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, "Formato de ID de paciente inválido", err))
			return
		}
		params.PatientID = patientID
	}

	if statusStr := c.Query("status"); statusStr != "" {
		params.Status = pgtype.Text{String: statusStr, Valid: true}
	}

	if dateStr := c.Query("date"); dateStr != "" {
		loc, err := time.LoadLocation("America/Argentina/Buenos_Aires")
		if err != nil {
			c.JSON(http.StatusInternalServerError, response.Error(http.StatusInternalServerError, "Error de zona horaria", err))
			return
		}
		date, err := time.ParseInLocation("2006-01-02", dateStr, loc)
		if err != nil {
			c.JSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, "Formato de fecha inválido", err))
			return
		}
		params.StartOfDay = pgtype.Timestamp{Time: date, Valid: true}
		params.EndOfDay = pgtype.Timestamp{Time: date.Add(24*time.Hour - time.Second), Valid: true}
	}

	countParams := sqlc.GetEventsFilteredCountParams{
		BusinessID:     params.BusinessID,
		StartOfDay:     params.StartOfDay,
		EndOfDay:       params.EndOfDay,
		PatientID:      params.PatientID,
		ProfessionalID: params.ProfessionalID,
		Status:         params.Status,
	}

	total, err := h.repo.GetEventsFilteredCount(c.Request.Context(), countParams)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(http.StatusInternalServerError, "Error al contar eventos", err))
		return
	}

	rawEvents, err := h.repo.GetEventsFiltered(c.Request.Context(), params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(http.StatusInternalServerError, "Error al obtener eventos", err))
		return
	}

	events := make([]json.RawMessage, len(rawEvents))
	for i, e := range rawEvents {
		events[i] = json.RawMessage(e)
	}

	result := response.PaginatedData[json.RawMessage]{
		Result: events,
		Total:  total,
	}
	c.JSON(http.StatusOK, response.Success("Eventos encontrados", &result))
}

func (h *EventHandler) GetByID(c *gin.Context) {
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

	rawEvent, err := h.repo.GetByID(c.Request.Context(), sqlc.GetEventByIDParams{
		BusinessID: businessID,
		ID:         id,
	})
	if err != nil {
		c.JSON(http.StatusNotFound, response.Error(http.StatusNotFound, "Evento no encontrado", err))
		return
	}

	event := json.RawMessage(rawEvent)

	c.JSON(http.StatusOK, response.Success("Evento encontrado", &event))
}

func (h *EventHandler) UpdateEvent(c *gin.Context) {
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

	var req UpdateEventRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, "Error de validación de datos", err))
		return
	}

	params := sqlc.UpdateEventParams{
		BusinessID: businessID,
		ID:         id,
	}

	if req.Title != nil {
		params.Title = pgtype.Text{String: *req.Title, Valid: true}
	}

	if req.StartDate != nil {
		t, err := time.Parse(time.RFC3339, *req.StartDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, "Formato de fecha de inicio inválido", err))
			return
		}
		params.StartDate = pgtype.Timestamptz{Time: t, Valid: true}
	}

	if req.EndDate != nil {
		t, err := time.Parse(time.RFC3339, *req.EndDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, "Formato de fecha de finalización inválido", err))
			return
		}
		params.EndDate = pgtype.Timestamptz{Time: t, Valid: true}
	}

	if req.ProfessionalID != nil {
		var professionalID pgtype.UUID
		if err := professionalID.Scan(*req.ProfessionalID); err != nil {
			c.JSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, "Formato de ID del profesional inválido", err))
			return
		}
		params.ProfessionalID = professionalID
	}

	if req.UserID != nil {
		var userID pgtype.UUID
		if err := userID.Scan(*req.UserID); err != nil {
			c.JSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, "Formato de ID del paciente inválido", err))
			return
		}
		params.UserID = userID
	}

	if req.Status != nil {
		status := sqlc.EventStatus(*req.Status)
		switch status {
		case sqlc.EventStatusAbsent, sqlc.EventStatusAttended, sqlc.EventStatusCancelled, sqlc.EventStatusInProgress, sqlc.EventStatusPending:
		default:
			c.JSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, "Estado inválido"))
			return
		}
		params.Status = sqlc.NullEventStatus{EventStatus: status, Valid: true}
	}

	event, err := h.repo.UpdateEvent(c.Request.Context(), params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(http.StatusInternalServerError, "Error al actualizar evento", err))
		return
	}

	c.JSON(http.StatusOK, response.Success("Evento actualizado", &event))
}

func (h *EventHandler) UpdateEventStatus(c *gin.Context) {
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

	var req UpdateEventStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, "Error de validación de datos", err))
		return
	}

	status := sqlc.EventStatus(req.Status)
	switch status {
	case sqlc.EventStatusAbsent, sqlc.EventStatusAttended, sqlc.EventStatusCancelled, sqlc.EventStatusInProgress, sqlc.EventStatusPending:
	default:
		c.JSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, "Estado inválido"))
		return
	}

	event, err := h.repo.UpdateEventStatus(c.Request.Context(), sqlc.UpdateEventStatusParams{
		BusinessID: businessID,
		ID:         id,
		Status:     status,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(http.StatusInternalServerError, "Error al actualizar estado del evento", err))
		return
	}

	c.JSON(http.StatusOK, response.Success("Estado del evento actualizado", &event))
}

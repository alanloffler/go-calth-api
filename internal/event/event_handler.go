package event

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/alanloffler/go-calth-api/internal/common/ctxkeys"
	"github.com/alanloffler/go-calth-api/internal/common/response"
	"github.com/alanloffler/go-calth-api/internal/database/sqlc"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type EventHandler struct {
	repo *EventRepository
	pool *pgxpool.Pool
}

type CreateEventRequest struct {
	Title          string   `json:"title" binding:"required,min=3,max=255"`
	StartDate      string   `json:"startDate" binding:"required,datetime=2006-01-02T15:04:05Z07:00"`
	EndDate        string   `json:"endDate" binding:"required,datetime=2006-01-02T15:04:05Z07:00"`
	ProfessionalID string   `json:"professionalId" binding:"required,uuid"`
	UserID         string   `json:"userId" binding:"required,uuid"`
	RecurringDates []string `json:"recurringDates" binding:"omitempty,dive,datetime=2006-01-02T15:04:05Z07:00"`
}

type UpdateEventRequest struct {
	Title          *string `json:"title" binding:"omitempty,min=3,max=255"`
	StartDate      *string `json:"startDate" binding:"omitempty,datetime=2006-01-02T15:04:05Z07:00"`
	EndDate        *string `json:"endDate" binding:"omitempty,datetime=2006-01-02T15:04:05Z07:00"`
	ProfessionalID *string `json:"professionalId" binding:"omitempty,uuid"`
	UserID         *string `json:"userId" binding:"omitempty,uuid"`
	Status         *string `json:"status" binding:"omitempty"`
	RecurrentID    *string `json:"recurrentId" binding:"omitempty,uuid"`
}

type UpdateEventStatusRequest struct {
	Status string `json:"status" binding:"required"`
}

func NewEventHandler(repo *EventRepository, pool *pgxpool.Pool) *EventHandler {
	return &EventHandler{repo: repo, pool: pool}
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

	if len(req.RecurringDates) > 0 {
		h.createRecurring(c, req, startTime, endTime, businessID, professionalID, userID)
		return
	}

	event, err := h.repo.Create(c.Request.Context(), sqlc.CreateEventParams{
		Title:          req.Title,
		StartDate:      pgtype.Timestamptz{Time: startTime, Valid: true},
		EndDate:        pgtype.Timestamptz{Time: endTime, Valid: true},
		BusinessID:     businessID,
		ProfessionalID: professionalID,
		UserID:         userID,
	})
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			c.JSON(http.StatusConflict, response.Error(http.StatusConflict, "El horario ya fue ocupado por otro usuario"))
			return
		}
		c.JSON(http.StatusInternalServerError, response.Error(http.StatusInternalServerError, "Error al crear evento", err))
		return
	}

	c.JSON(http.StatusOK, response.Created("Evento creado", &event))
}

func (h *EventHandler) createRecurring(c *gin.Context, req CreateEventRequest, startTime, endTime time.Time, businessID, professionalID, userID pgtype.UUID) {
	ctx := c.Request.Context()
	duration := endTime.Sub(startTime)

	var recurrentID pgtype.UUID
	if err := recurrentID.Scan(uuid.New().String()); err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(http.StatusInternalServerError, "Error al generar ID de recurrente", err))
		return
	}

	tx, err := h.pool.Begin(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(http.StatusInternalServerError, "Error al iniciar transacción", err))
		return
	}
	defer tx.Rollback(ctx)

	qtx := sqlc.New(tx)
	events := make([]sqlc.Event, 0, len(req.RecurringDates))

	for _, dateStr := range req.RecurringDates {
		recurringStart, err := time.Parse(time.RFC3339, dateStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, "Formato de fecha recurrente inválido", err))
			return
		}
		recurringEnd := recurringStart.Add(duration)

		event, err := qtx.CreateEvent(ctx, sqlc.CreateEventParams{
			Title:          req.Title,
			StartDate:      pgtype.Timestamptz{Time: recurringStart, Valid: true},
			EndDate:        pgtype.Timestamptz{Time: recurringEnd, Valid: true},
			BusinessID:     businessID,
			ProfessionalID: professionalID,
			UserID:         userID,
			RecurrentID:    recurrentID,
		})
		if err != nil {
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) && pgErr.Code == "23505" {
				c.JSON(http.StatusConflict, response.Error(http.StatusConflict, "Uno o más horarios ya están ocupados"))
				return
			}
			c.JSON(http.StatusInternalServerError, response.Error(http.StatusInternalServerError, "Error al crear turnos recurrentes", err))
			return
		}
		events = append(events, event)
	}

	if err := tx.Commit(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(http.StatusInternalServerError, "Error al confirmar transacción", err))
		return
	}

	c.JSON(http.StatusCreated, response.Created("Turnos recurrentes creados", &events))
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

	rawEvents, err := h.repo.GetByBusinessID(c.Request.Context(), sqlc.GetByBusinessIDParams{BusinessID: businessID, Limit: int32(limit)})
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

	var startDate, endDate pgtype.Timestamptz

	if startDateStr := c.Query("startDate"); startDateStr != "" {
		loc, err := time.LoadLocation("America/Argentina/Buenos_Aires")
		if err != nil {
			c.JSON(http.StatusInternalServerError, response.Error(http.StatusInternalServerError, "Error de zona horaria", err))
			return
		}
		parsedStartDate, err := time.ParseInLocation(time.RFC3339, startDateStr, loc)
		if err != nil {
			c.JSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, "Formato de fecha de inicio inválido", err))
			return
		}
		startDate = pgtype.Timestamptz{Time: parsedStartDate, Valid: true}
	}

	if endDateStr := c.Query("endDate"); endDateStr != "" {
		loc, err := time.LoadLocation("America/Argentina/Buenos_Aires")
		if err != nil {
			c.JSON(http.StatusInternalServerError, response.Error(http.StatusInternalServerError, "Error de zona horaria", err))
			return
		}
		parsedEndDate, err := time.ParseInLocation(time.RFC3339, endDateStr, loc)
		if err != nil {
			c.JSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, "Formato de fecha de fin inválido", err))
			return
		}
		endDate = pgtype.Timestamptz{Time: parsedEndDate, Valid: true}
	}

	rawEvents, err := h.repo.GetByProfessionalID(c.Request.Context(), sqlc.GetEventsByProfessionalIDParams{
		BusinessID:     businessID,
		ProfessionalID: id,
		StartDate:      startDate,
		EndDate:        endDate,
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

func (h *EventHandler) GetByBusinessProfessionalPatient(c *gin.Context) {
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

	rawEvents, err := h.repo.GetByBusinessProfessionalPatient(c.Request.Context(), sqlc.GetByBusinessProfessionalPatientParams{
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

func (h *EventHandler) GetByProfessionalDay(c *gin.Context) {
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

	rawEvents, err := h.repo.GetByProfessionalDay(c.Request.Context(), sqlc.GetByProfessionalDayParams{
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

func (h *EventHandler) GetByProfessionalDayArray(c *gin.Context) {
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

	slotTimes, err := h.repo.GetByProfessionalDayArray(c.Request.Context(), sqlc.GetByProfessionalDayArrayParams{
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

func (h *EventHandler) GetFiltered(c *gin.Context) {
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

	if recurrent := c.Query("recurrent"); recurrent != "" {
		params.Recurrent = pgtype.Text{String: recurrent, Valid: true}
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

	countParams := sqlc.GetFilteredCountParams{
		BusinessID:     params.BusinessID,
		StartOfDay:     params.StartOfDay,
		EndOfDay:       params.EndOfDay,
		PatientID:      params.PatientID,
		ProfessionalID: params.ProfessionalID,
		Recurrent:      params.Recurrent,
		Status:         params.Status,
	}

	total, err := h.repo.GetFilteredCount(c.Request.Context(), countParams)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(http.StatusInternalServerError, "Error al contar eventos", err))
		return
	}

	rawEvents, err := h.repo.GetFiltered(c.Request.Context(), params)
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

func (h *EventHandler) GetDaysWithEvents(c *gin.Context) {
	businessID, ok := ctxkeys.BusinessID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, response.Error(http.StatusUnauthorized, "Usuario no autenticado"))
		return
	}

	var professionalID pgtype.UUID
	if err := professionalID.Scan(c.Param("id")); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, "Formato de ID profesional inválido", err))
		return
	}

	fromDateStr := c.Query("fromDate")
	toDateStr := c.Query("toDate")

	if fromDateStr == "" || toDateStr == "" {
		c.JSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, "fromDate y toDate son requeridos"))
		return
	}

	loc, err := time.LoadLocation("America/Argentina/Buenos_Aires")
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(http.StatusInternalServerError, "Error de zona horaria", err))
		return
	}

	fromDate, err := time.ParseInLocation("2006-01-02", fromDateStr, loc)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, "Formato de fecha inválido", err))
		return
	}

	toDate, err := time.ParseInLocation("2006-01-02", toDateStr, loc)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, "Formato de fecha inválido", err))
		return
	}

	days, err := h.repo.GetDaysWithEvents(c.Request.Context(), sqlc.GetDaysWithEventsParams{
		BusinessID:     businessID,
		ProfessionalID: professionalID,
		StartDate:      pgtype.Timestamptz{Time: fromDate, Valid: true},
		StartDate_2:    pgtype.Timestamptz{Time: toDate.Add(24*time.Hour - time.Second), Valid: true},
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(http.StatusInternalServerError, "Error al obtener días con eventos", err))
		return
	}

	result := make(map[string]bool)
	for _, d := range days {
		result[d.Time.Format("2")] = true
	}

	c.JSON(http.StatusOK, response.Success("Días ocupados", &result))
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

	rawEvent, err := h.repo.GetByID(c.Request.Context(), sqlc.GetByIDParams{
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

func (h *EventHandler) CheckRecurring(c *gin.Context) {
	businessID, ok := ctxkeys.BusinessID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, response.Error(http.StatusUnauthorized, "Usuario no autenticado"))
		return
	}

	var professionalID pgtype.UUID
	if err := professionalID.Scan(c.Query("professionalId")); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, "Formato de ID del profesional inválido", err))
		return
	}

	startDateStr := c.Query("startDate")
	parsedTime, err := time.Parse(time.RFC3339, startDateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, "Formato de fecha inválido", err))
		return
	}

	occurrences, err := strconv.Atoi(c.Query("days"))
	if err != nil || occurrences <= 0 {
		c.JSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, "Parámetro 'days' inválido"))
		return
	}

	recurringDates := generateRecurringDates(parsedTime, int32(occurrences))

	existingEvents, err := h.repo.CheckRecurring(c, sqlc.CheckRecurringEventsParams{
		BusinessID:     businessID,
		ProfessionalID: professionalID,
		StartDate:      pgtype.Timestamptz{Time: parsedTime, Valid: true},
		Column4:        pgtype.Text{String: strconv.Itoa(occurrences * 7), Valid: true},
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(http.StatusInternalServerError, "Error al verificar eventos recurrentes", err))
		return
	}

	busySlots := make(map[string]bool, len(existingEvents))
	for _, e := range existingEvents {
		busySlots[e.StartDate.Time.UTC().Format("2006-01-02T15:04")] = true
	}

	type recurringResult struct {
		Date      time.Time `json:"date"`
		Available bool      `json:"available"`
	}

	results := make([]recurringResult, len(recurringDates))
	for i, d := range recurringDates {
		results[i] = recurringResult{
			Date:      d,
			Available: !busySlots[d.UTC().Format("2006-01-02T15:04")],
		}
	}

	c.JSON(http.StatusOK, response.Success("Recurrencia verificada", &results))
}

func (h *EventHandler) Update(c *gin.Context) {
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
		case sqlc.EventStatusAbsent, sqlc.EventStatusPresent, sqlc.EventStatusCancelled, sqlc.EventStatusInProgress, sqlc.EventStatusPending:
		default:
			c.JSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, "Estado inválido"))
			return
		}
		params.Status = sqlc.NullEventStatus{EventStatus: status, Valid: true}
	}

	if req.RecurrentID != nil {
		var recurrentID pgtype.UUID
		if err := recurrentID.Scan(*req.RecurrentID); err != nil {
			c.JSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, "Formato de ID de recurrente inválido", err))
			return
		}
		params.RecurrentID = recurrentID
	}

	event, err := h.repo.Update(c.Request.Context(), params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(http.StatusInternalServerError, "Error al actualizar evento", err))
		return
	}

	c.JSON(http.StatusOK, response.Success("Evento actualizado", &event))
}

func (h *EventHandler) UpdateStatus(c *gin.Context) {
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
	case sqlc.EventStatusAbsent, sqlc.EventStatusPresent, sqlc.EventStatusCancelled, sqlc.EventStatusInProgress, sqlc.EventStatusPending:
	default:
		c.JSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, "Estado inválido"))
		return
	}

	affected, err := h.repo.UpdateStatus(c.Request.Context(), sqlc.UpdateStatusParams{
		BusinessID: businessID,
		ID:         id,
		Status:     status,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(http.StatusInternalServerError, "Error al actualizar estado del evento", err))
		return
	}
	if affected == 0 {
		c.JSON(http.StatusNotFound, response.Error(http.StatusNotFound, "Evento no encontrado"))
		return
	}

	c.JSON(http.StatusOK, response.Success[any]("Estado del evento actualizado", nil))
}

func (h *EventHandler) Delete(c *gin.Context) {
	ctx := c.Request.Context()

	businessID, ok := ctxkeys.BusinessID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, response.Error(http.StatusUnauthorized, "Usuario no autenticado"))
		return
	}

	var id pgtype.UUID
	if err := id.Scan(c.Param("id")); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, "Formato de ID de turno inválido", err))
		return
	}

	// PLAN:
	// 1. Get the event's recurrent_id
	recurrentID, err := h.repo.GetEventRecurrentID(ctx, sqlc.GetEventRecurrentIDParams{
		BusinessID: businessID,
		ID:         id,
	})
	if err != nil {
		c.JSON(http.StatusNotFound, response.Error(http.StatusNotFound, "Turno no encontrado", err))
		return
	}

	// 2. Not recurring: simple delete
	if !recurrentID.Valid {
		affected, err := h.repo.Delete(ctx, sqlc.DeleteEventParams{
			BusinessID: businessID,
			ID:         id,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, response.Error(http.StatusInternalServerError, "Error al eliminar el turno", err))
			return
		}
		if affected == 0 {
			c.JSON(http.StatusNotFound, response.Error(http.StatusNotFound, "Turno no encontrado"))
			return
		}
		c.JSON(http.StatusOK, response.Success[any]("Turno eliminado", nil))
		return
	}

	// 3. Recurring: count siblings (including self)
	siblingIDs, err := h.repo.GetIDsByRecurrentID(ctx, sqlc.GetIDsByRecurrentIDParams{
		RecurrentID: recurrentID,
		BusinessID:  businessID,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(http.StatusInternalServerError, "Error al obtener turnos recurrentes", err))
		return
	}

	// 4. Transactional delete
	tx, err := h.pool.Begin(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(http.StatusInternalServerError, "Error al iniciar transacción", err))
		return
	}
	defer tx.Rollback(ctx)

	qtx := sqlc.New(tx)

	affected, err := qtx.DeleteEvent(ctx, sqlc.DeleteEventParams{
		BusinessID: businessID,
		ID:         id,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(http.StatusInternalServerError, "Error al eliminar el turno", err))
		return
	}
	if affected == 0 {
		c.JSON(http.StatusNotFound, response.Error(http.StatusInternalServerError, "Turno no encontrado"))
		return
	}

	// 5. Only 1 sibling remains after deletion: clear its recurrent_id
	if len(siblingIDs) == 2 {
		for _, siblindID := range siblingIDs {
			if siblindID == id {
				continue
			}
			if _, err := qtx.ClearRecurrentID(ctx, sqlc.ClearRecurrentIDParams{
				BusinessID: businessID,
				ID:         siblindID,
			}); err != nil {
				c.JSON(http.StatusInternalServerError, response.Error(http.StatusInternalServerError, "Error al actualizar turno huérfano", err))
				return
			}
		}
	}

	if err := tx.Commit(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(http.StatusInternalServerError, "Error al confirmar transacción", err))
		return
	}

	c.JSON(http.StatusOK, response.Success[any]("Turno eliminado", nil))
}

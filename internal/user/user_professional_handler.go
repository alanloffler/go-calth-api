package user

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/alanloffler/go-calth-api/internal/common/ctxkeys"
	"github.com/alanloffler/go-calth-api/internal/common/response"
	"github.com/alanloffler/go-calth-api/internal/common/utils"
	"github.com/alanloffler/go-calth-api/internal/database/sqlc"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type CreateProfessionalRequest struct {
	User    CreateUserData                `json:"user" binding:"required"`
	Profile CreateProfessionalProfileData `json:"profile" binding:"required"`
}

type CreateProfessionalProfileData struct {
	LicenseID           string  `json:"license_id" binding:"required"`
	ProfessionalPrefix  string  `json:"professional_prefix" binding:"required"`
	Specialty           string  `json:"specialty" binding:"required"`
	WorkingDays         []int   `json:"working_days" binding:"required"`
	StartHour           string  `json:"start_hour" binding:"required"`
	EndHour             string  `json:"end_hour" binding:"required"`
	SlotDuration        string  `json:"slot_duration" binding:"required"`
	DailyExceptionStart *string `json:"daily_exception_start"`
	DailyExceptionEnd   *string `json:"daily_exception_end"`
}

func (h *UserHandler) CreateProfessional(c *gin.Context) {
	businessID, ok := ctxkeys.BusinessID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, response.Error(http.StatusUnauthorized, "Usuario no autenticado"))
		return
	}

	var req CreateProfessionalRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, "Error al crear profesional", err))
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.User.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(http.StatusInternalServerError, "Error al procesar contraseña", err))
		return
	}

	ctx := c.Request.Context()

	role, err := h.repo.q.GetRoleByValue(ctx, "professional")
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, "Rol de profesional no encontrado", err))
		return
	}

	tx, err := h.pool.Begin(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(http.StatusInternalServerError, "Error al iniciar transacción", err))
		return
	}
	defer tx.Rollback(ctx)

	qtx := sqlc.New(tx)

	user, err := qtx.CreateUser(ctx, sqlc.CreateUserParams{
		Ic:          req.User.Ic,
		UserName:    req.User.UserName,
		FirstName:   req.User.FirstName,
		LastName:    req.User.LastName,
		Email:       req.User.Email,
		Password:    string(hashedPassword),
		PhoneNumber: req.User.PhoneNumber,
		RoleID:      role.ID,
		BusinessID:  businessID,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(http.StatusInternalServerError, "Error al crear usuario", err))
		return
	}

	days := make([]string, len(req.Profile.WorkingDays))
	for i, d := range req.Profile.WorkingDays {
		days[i] = strconv.Itoa(d)
	}
	workingDays := strings.Join(days, ",")

	_, err = qtx.CreateProfessionalProfile(ctx, sqlc.CreateProfessionalProfileParams{
		BusinessID:          businessID,
		UserID:              user.ID,
		LicenseID:           req.Profile.LicenseID,
		ProfessionalPrefix:  req.Profile.ProfessionalPrefix,
		Specialty:           req.Profile.Specialty,
		WorkingDays:         workingDays,
		StartHour:           req.Profile.StartHour,
		EndHour:             req.Profile.EndHour,
		SlotDuration:        req.Profile.SlotDuration,
		DailyExceptionStart: utils.ToPgText(req.Profile.DailyExceptionStart),
		DailyExceptionEnd:   utils.ToPgText(req.Profile.DailyExceptionEnd),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(http.StatusInternalServerError, "Error al crear perfil de profesional", err))
		return
	}

	if err := tx.Commit(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(http.StatusInternalServerError, "Error al confirmar transacción", err))
		return
	}

	c.JSON(http.StatusCreated, response.Created("Profesional creado", &user))
}

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
	"github.com/jackc/pgx/v5/pgtype"
	"golang.org/x/crypto/bcrypt"
)

type CreateProfessionalRequest struct {
	User    CreateUserData                `json:"user" binding:"required"`
	Profile CreateProfessionalProfileData `json:"profile" binding:"required"`
}

type CreateProfessionalProfileData struct {
	LicenseID           string  `json:"licenseId" binding:"required"`
	ProfessionalPrefix  string  `json:"professionalPrefix" binding:"required"`
	Specialty           string  `json:"specialty" binding:"required"`
	WorkingDays         []int   `json:"workingDays" binding:"required"`
	StartHour           string  `json:"startHour" binding:"required"`
	EndHour             string  `json:"endHour" binding:"required"`
	SlotDuration        string  `json:"slotDuration" binding:"required"`
	DailyExceptionStart *string `json:"dailyExceptionStart"`
	DailyExceptionEnd   *string `json:"dailyExceptionEnd"`
}

type userWithProfessionalProfile struct {
	ID                  pgtype.UUID                 `json:"id"`
	Ic                  string                      `json:"ic"`
	UserName            string                      `json:"userName"`
	FirstName           string                      `json:"firstName"`
	LastName            string                      `json:"lastName"`
	Email               string                      `json:"email"`
	PhoneNumber         string                      `json:"phoneNumber"`
	Role                *userRole                   `json:"role"`
	BusinessID          pgtype.UUID                 `json:"businessID"`
	CreatedAt           pgtype.Timestamptz          `json:"createdAt"`
	UpdatedAt           pgtype.Timestamptz          `json:"updatedAt"`
	DeletedAt           pgtype.Timestamptz          `json:"deletedAt"`
	ProfessionalProfile professionalProfileResponse `json:"professionalProfile"`
}

type professionalProfileResponse struct {
	ID                  pgtype.UUID        `json:"id"`
	LicenseID           string             `json:"licenseId"`
	ProfessionalPrefix  string             `json:"professionalPrefix"`
	Specialty           string             `json:"specialty"`
	WorkingDays         []string           `json:"workingDays"`
	StartHour           string             `json:"startHour"`
	EndHour             string             `json:"endHour"`
	SlotDuration        string             `json:"slotDuration"`
	DailyExceptionStart pgtype.Text        `json:"dailyExceptionStart"`
	DailyExceptionEnd   pgtype.Text        `json:"dailyExceptionEnd"`
	CreatedAt           pgtype.Timestamptz `json:"createdAt"`
	UpdatedAt           pgtype.Timestamptz `json:"updatedAt"`
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
		c.JSON(http.StatusInternalServerError, response.Error(http.StatusInternalServerError, "Error al crear profesional", err))
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

func (h *UserHandler) GetProfessionalByID(c *gin.Context) {
	h.getProfessionalByID(c, false)
}

func (h *UserHandler) GetProfessionalByIDWithSoftDeleted(c *gin.Context) {
	h.getProfessionalByID(c, true)
}

func (h *UserHandler) getProfessionalByID(c *gin.Context, withSoftDeleted bool) {
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

	ctx := c.Request.Context()

	var user userWithProfessionalProfile
	if withSoftDeleted {
		row, err := h.repo.GetByIDWithSoftDeleted(ctx, sqlc.GetUserByIDWithSoftDeletedParams{BusinessID: businessID, ID: id})
		if err != nil {
			c.JSON(http.StatusNotFound, response.Error(http.StatusNotFound, "Profesional no encontrado", err))
			return
		}

		user = userWithProfessionalProfile{
			ID:          row.ID,
			Ic:          row.Ic,
			UserName:    row.UserName,
			FirstName:   row.FirstName,
			LastName:    row.LastName,
			Email:       row.Email,
			PhoneNumber: row.PhoneNumber,
			BusinessID:  row.BusinessID,
			CreatedAt:   row.CreatedAt,
			UpdatedAt:   row.UpdatedAt,
			DeletedAt:   row.DeletedAt,
		}
		if row.RoleID.Valid {
			user.Role = &userRole{
				ID:          row.RoleID,
				Name:        row.RoleName.String,
				Value:       row.RoleValue.String,
				Description: row.RoleDescription.String,
			}
		}
	} else {
		row, err := h.repo.GetByID(ctx, sqlc.GetUserByIDParams{BusinessID: businessID, ID: id})
		if err != nil {
			c.JSON(http.StatusNotFound, response.Error(http.StatusNotFound, "Profesional no encontrado", err))
			return
		}

		user = userWithProfessionalProfile{
			ID:          row.ID,
			Ic:          row.Ic,
			UserName:    row.UserName,
			FirstName:   row.FirstName,
			LastName:    row.LastName,
			Email:       row.Email,
			PhoneNumber: row.PhoneNumber,
			BusinessID:  row.BusinessID,
			CreatedAt:   row.CreatedAt,
			UpdatedAt:   row.UpdatedAt,
			DeletedAt:   row.DeletedAt,
		}
		if row.RoleID.Valid {
			user.Role = &userRole{
				ID:          row.RoleID,
				Name:        row.RoleName.String,
				Value:       row.RoleValue.String,
				Description: row.RoleDescription.String,
			}
		}
	}

	profile, err := h.professionalProfileRepo.GetProfessionalProfileByUserID(c.Request.Context(), sqlc.GetProfessionalProfileByUserIDParams{BusinessID: businessID, UserID: id})
	if err != nil {
		c.JSON(http.StatusNotFound, response.Error(http.StatusNotFound, "Perfil profesional no encontrado", err))
		return
	}

	workingDays := strings.Split(profile.WorkingDays, ",")

	profResponse := professionalProfileResponse{
		ID:                  profile.ID,
		LicenseID:           profile.LicenseID,
		ProfessionalPrefix:  profile.ProfessionalPrefix,
		Specialty:           profile.Specialty,
		WorkingDays:         workingDays,
		StartHour:           profile.StartHour,
		EndHour:             profile.EndHour,
		SlotDuration:        profile.SlotDuration,
		DailyExceptionStart: profile.DailyExceptionStart,
		DailyExceptionEnd:   profile.DailyExceptionEnd,
		CreatedAt:           profile.CreatedAt,
		UpdatedAt:           profile.UpdatedAt,
	}

	user.ProfessionalProfile = profResponse

	c.JSON(http.StatusOK, response.Success("Profesional encontrado", &user))
}

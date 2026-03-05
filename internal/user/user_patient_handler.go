package user

import (
	"fmt"
	"log"
	"math"
	"math/big"
	"net/http"
	"time"

	"github.com/alanloffler/go-calth-api/internal/common/ctxkeys"
	"github.com/alanloffler/go-calth-api/internal/common/response"
	"github.com/alanloffler/go-calth-api/internal/common/utils"
	"github.com/alanloffler/go-calth-api/internal/database/sqlc"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
	"golang.org/x/crypto/bcrypt"
)

type CreatePatientRequest struct {
	User    CreateUserData           `json:"user" binding:"required"`
	Profile CreatePatientProfileData `json:"profile" binding:"required"`
}

type UpdatePatientRequest struct {
	User    UpdateUserData           `json:"user" binding:"required"`
	Profile UpdatePatientProfileData `json:"profile" binding:"required"`
}

type UpdatePatientResponse struct {
	User    sqlc.User           `json:"user"`
	Profile sqlc.PatientProfile `json:"profile"`
}

type CreatePatientProfileData struct {
	Gender                string  `json:"gender" binding:"required"`
	BirthDay              string  `json:"birthDay" binding:"required"`
	BloodType             string  `json:"bloodType" binding:"required"`
	Weight                float64 `json:"weight" binding:"required,gt=0,lt=999.99"`
	Height                float64 `json:"height" binding:"required,gt=0,lt=300"`
	EmergencyContactName  string  `json:"emergencyContactName" binding:"required"`
	EmergencyContactPhone string  `json:"emergencyContactPhone" binding:"required,len=10,numeric"`
}

type UpdatePatientProfileData struct {
	Gender                *string  `json:"gender" binding:"omitempty"`
	BirthDay              *string  `json:"birthDay" binding:"omitempty"`
	BloodType             *string  `json:"bloodType" binding:"omitempty"`
	Weight                *float64 `json:"weight" binding:"omitempty,gt=0,lt=999.99"`
	Height                *float64 `json:"height" binding:"omitempty,gt=0,lt=300"`
	EmergencyContactName  *string  `json:"emergencyContactName" binding:"omitempty"`
	EmergencyContactPhone *string  `json:"emergencyContactPhone" binding:"omitempty,len=10,numeric"`
}

type userWithPatientProfile struct {
	ID             pgtype.UUID            `json:"id"`
	Ic             string                 `json:"ic"`
	UserName       string                 `json:"userName"`
	FirstName      string                 `json:"firstName"`
	LastName       string                 `json:"lastName"`
	Email          string                 `json:"email"`
	PhoneNumber    string                 `json:"phoneNumber"`
	Role           *userRole              `json:"role"`
	BusinessID     pgtype.UUID            `json:"businessId"`
	CreatedAt      pgtype.Timestamptz     `json:"createdAt"`
	UpdatedAt      pgtype.Timestamptz     `json:"updatedAt"`
	DeletedAt      pgtype.Timestamptz     `json:"deletedAt"`
	PatientProfile patientProfileResponse `json:"patientProfile"`
}

type patientProfileResponse struct {
	Gender                string  `json:"gender"`
	BirthDay              string  `json:"birthDay"`
	BloodType             string  `json:"bloodType"`
	Weight                float64 `json:"weight"`
	Height                float64 `json:"height"`
	EmergencyContactName  string  `json:"emergencyContactName"`
	EmergencyContactPhone string  `json:"emergencyContactPhone"`
}

func (h *UserHandler) CreatePatient(c *gin.Context) {
	businessID, ok := ctxkeys.BusinessID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, response.Error(http.StatusUnauthorized, "Usuario no autenticado"))
		return
	}

	var req CreatePatientRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, "Error al crear usuario", err))
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.User.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(http.StatusInternalServerError, "Error al procesar contraseña", err))
		return
	}

	ctx := c.Request.Context()

	role, err := h.repo.q.GetRoleByValue(ctx, "patient")
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, "Rol de paciente no encontrado", err))
		return
	}

	birthDay, err := time.Parse("2006-01-02", req.Profile.BirthDay)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, "Formato de fecha inválido", err))
		return
	}

	var pgBirthDay pgtype.Date
	if err := pgBirthDay.Scan(birthDay); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, "Formato de fecha inválido", err))
		return
	}

	var weight pgtype.Numeric
	if err := weight.Scan(fmt.Sprintf("%g", req.Profile.Weight)); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, "Formato de peso inválido", err))
		return
	}

	var height pgtype.Numeric
	if err := height.Scan(fmt.Sprintf("%g", req.Profile.Height)); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, "Formato de altura inválido", err))
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

	_, err = qtx.CreatePatientProfile(ctx, sqlc.CreatePatientProfileParams{
		BusinessID:            businessID,
		UserID:                user.ID,
		Gender:                req.Profile.Gender,
		BirthDay:              pgBirthDay,
		BloodType:             req.Profile.BloodType,
		Weight:                weight,
		Height:                height,
		EmergencyContactName:  req.Profile.EmergencyContactName,
		EmergencyContactPhone: req.Profile.EmergencyContactPhone,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(http.StatusInternalServerError, "Error al crear perfil de paciente", err))
		return
	}

	if err := tx.Commit(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(http.StatusInternalServerError, "Error al confirmar transacción", err))
		return
	}

	c.JSON(http.StatusCreated, response.Created("Usuario creado", &user))
}

func (h *UserHandler) GetPatientByID(c *gin.Context) {
	h.getPatientByID(c, false)
}

func (h *UserHandler) GetPatientByIDWithSoftDeleted(c *gin.Context) {
	h.getPatientByID(c, true)
}

func (h *UserHandler) getPatientByID(c *gin.Context, withSoftDeleted bool) {
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
	log.Println(id)

	ctx := c.Request.Context()

	var user userWithPatientProfile
	if withSoftDeleted {
		row, err := h.repo.GetByIDWithSoftDeleted(ctx, sqlc.GetUserByIDWithSoftDeletedParams{BusinessID: businessID, ID: id})
		if err != nil {
			c.JSON(http.StatusNotFound, response.Error(http.StatusNotFound, "Paciente no encontrado", err))
			return
		}

		user = userWithPatientProfile{
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
			c.JSON(http.StatusNotFound, response.Error(http.StatusNotFound, "Paciente no encontrado", err))
			return
		}

		user = userWithPatientProfile{
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

	profile, err := h.patientProfileRepo.GetPatientProfileByUserID(c.Request.Context(), sqlc.GetPatientProfileByUserIDParams{
		BusinessID: businessID,
		UserID:     id,
	})
	if err != nil {
		c.JSON(http.StatusNotFound, response.Error(http.StatusNotFound, "Perfil de paciente no encontrado", err))
		return
	}

	profResponse := patientProfileResponse{
		Gender:    profile.Gender,
		BirthDay:  profile.BirthDay.Time.Format("2006/01/02"),
		BloodType: profile.BloodType,
		Weight: func() float64 {
			f, _ := new(big.Float).SetInt(profile.Weight.Int).Float64()
			return f * math.Pow10(int(profile.Weight.Exp))
		}(),
		Height: func() float64 {
			f, _ := new(big.Float).SetInt(profile.Height.Int).Float64()
			return f * math.Pow10(int(profile.Height.Exp))
		}(),
		EmergencyContactName:  profile.EmergencyContactName,
		EmergencyContactPhone: profile.EmergencyContactPhone,
	}

	user.PatientProfile = profResponse
	log.Print(user)

	c.JSON(http.StatusOK, response.Success("Paciente encontrado", &user))
}

func (h *UserHandler) UpdatePatient(c *gin.Context) {
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

	var req UpdatePatientRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, "Error al actualizar paciente"))
		return
	}

	ctx := c.Request.Context()

	var passwordHash pgtype.Text
	if req.User.Password != nil {
		hashed, err := bcrypt.GenerateFromPassword([]byte(*req.User.Password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, response.Error(http.StatusInternalServerError, "Error al procesar contraseña", err))
			return
		}
		passwordHash = pgtype.Text{String: string(hashed), Valid: true}
	}

	var pgBirthDay pgtype.Date
	if req.Profile.BirthDay != nil {
		birthDay, err := time.Parse("2006-01-02", *req.Profile.BirthDay)
		if err != nil {
			c.JSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, "Formato de fecha inválido", err))
			return
		}
		if err := pgBirthDay.Scan(birthDay); err != nil {
			c.JSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, "Formato de fecha inválido", err))
			return
		}
	}

	var weight pgtype.Numeric
	if req.Profile.Weight != nil {
		if err := weight.Scan(fmt.Sprintf("%g", *req.Profile.Weight)); err != nil {
			c.JSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, "Formato de peso inválido", err))
			return
		}
	}

	var height pgtype.Numeric
	if req.Profile.Height != nil {
		if err := height.Scan(fmt.Sprintf("%g", *req.Profile.Height)); err != nil {
			c.JSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, "Formato de altura inválido", err))
			return
		}
	}

	tx, err := h.pool.Begin(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(http.StatusInternalServerError, "Error al iniciar transacción", err))
		return
	}
	defer tx.Rollback(ctx)

	qtx := sqlc.New(tx)

	updatedUser, err := qtx.UpdateUser(ctx, sqlc.UpdateUserParams{
		ID:          id,
		Ic:          utils.ToPgText(req.User.Ic),
		UserName:    utils.ToPgText(req.User.UserName),
		FirstName:   utils.ToPgText(req.User.FirstName),
		LastName:    utils.ToPgText(req.User.LastName),
		Email:       utils.ToPgText(req.User.Email),
		Password:    passwordHash,
		PhoneNumber: utils.ToPgText(req.User.PhoneNumber),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(http.StatusInternalServerError, "Error al actualizar usuario", err))
		return
	}

	updatedProfile, err := qtx.UpdatePatientProfile(ctx, sqlc.UpdatePatientProfileParams{
		BusinessID:            businessID,
		UserID:                id,
		Gender:                utils.ToPgText(req.Profile.Gender),
		BirthDay:              pgBirthDay,
		BloodType:             utils.ToPgText(req.Profile.BloodType),
		Weight:                weight,
		Height:                height,
		EmergencyContactName:  utils.ToPgText(req.Profile.EmergencyContactName),
		EmergencyContactPhone: utils.ToPgText(req.Profile.EmergencyContactPhone),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(http.StatusInternalServerError, "Error al actualizar perfil", err))
		return
	}

	if err := tx.Commit(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(http.StatusInternalServerError, "Error al confirmar transacción", err))
		return
	}

	c.JSON(http.StatusOK, response.Success("Paciente actualizado", &UpdatePatientResponse{
		User:    updatedUser,
		Profile: updatedProfile,
	}))
}

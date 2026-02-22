package user

import (
	"fmt"
	"net/http"
	"time"

	"github.com/alanloffler/go-calth-api/internal/common/ctxkeys"
	"github.com/alanloffler/go-calth-api/internal/common/response"
	"github.com/alanloffler/go-calth-api/internal/database/sqlc"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

type UserHandler struct {
	repo *UserRepository
	pool *pgxpool.Pool
}

func NewUserHandler(repo *UserRepository, pool *pgxpool.Pool) *UserHandler {
	return &UserHandler{repo: repo, pool: pool}
}

type CreateUserRequest struct {
	Ic             string                       `json:"ic" binding:"required,len=8"`
	UserName       string                       `json:"userName" binding:"required,min=3,max=100"`
	FirstName      string                       `json:"firstName" binding:"required,min=3,max=100"`
	LastName       string                       `json:"lastName" binding:"required,min=3,max=100"`
	Email          string                       `json:"email" binding:"required,email,max=100"`
	Password       string                       `json:"password" binding:"required,min=8,max=100"`
	PhoneNumber    string                       `json:"phoneNumber" binding:"required,len=10,numeric"`
	RoleID         string                       `json:"roleId" binding:"required,uuid"`
	BusinessID     string                       `json:"businessId" binding:"required,uuid"`
	PatientProfile *CreatePatientProfileRequest `json:"patientProfile"`
}

type UpdateUserRequest struct {
	Ic          string  `json:"ic" binding:"required,len=8"`
	UserName    string  `json:"userName" binding:"required,min=3,max=100"`
	FirstName   string  `json:"firstName" binding:"required,min=3,max=100"`
	LastName    string  `json:"lastName" binding:"required,min=3,max=100"`
	Email       string  `json:"email" binding:"required,email,max=100"`
	Password    *string `json:"password" binding:"omitempty,min=8,max=100"`
	PhoneNumber string  `json:"phoneNumber" binding:"required,len=10,numeric"`
	RoleID      string  `json:"roleId" binding:"required,uuid"`
}

type CreatePatientProfileRequest struct {
	Gender                string  `json:"gender" binding:"required"`
	BirthDay              string  `json:"birthDay" binding:"required"`
	BloodType             string  `json:"bloodType" binding:"required"`
	Weight                float64 `json:"weight" binding:"required,gt=0,lt=999.99"`
	Height                float64 `json:"height" binding:"required,gt=0,lt=300"`
	EmergencyContactName  string  `json:"emergencyContactName" binding:"required"`
	EmergencyContactPhone string  `json:"emergencyContactPhone" binding:"required,len=10,numeric"`
}

type userRole struct {
	ID          pgtype.UUID `json:"id"`
	Name        string      `json:"name"`
	Value       string      `json:"value"`
	Description string      `json:"description"`
}

type userByRoleResponse struct {
	ID          pgtype.UUID        `json:"id"`
	Ic          string             `json:"ic"`
	UserName    string             `json:"userName"`
	FirstName   string             `json:"firstName"`
	LastName    string             `json:"lastName"`
	Email       string             `json:"email"`
	PhoneNumber string             `json:"phoneNumber"`
	Role        *userRole          `json:"role"`
	BusinessID  pgtype.UUID        `json:"businessID"`
	CreatedAt   pgtype.Timestamptz `json:"createdAt"`
	UpdatedAt   pgtype.Timestamptz `json:"updatedAt"`
	DeletedAt   pgtype.Timestamptz `json:"deletedAt"`
}

func (h *UserHandler) Create(c *gin.Context) {
	var req CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, "Error al crear usuario", err))
		return
	}

	var roleID pgtype.UUID
	if err := roleID.Scan(req.RoleID); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, "Formato de ID inválido", err))
		return
	}

	var businessID pgtype.UUID
	if err := businessID.Scan(req.BusinessID); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, "Formato de ID inválido", err))
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(http.StatusInternalServerError, "Error al procesar contraseña", err))
		return
	}

	ctx := c.Request.Context()

	role, err := h.repo.q.GetRoleByID(ctx, roleID)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, "Rol no encontrado", err))
		return
	}

	userArg := sqlc.CreateUserParams{
		Ic:          req.Ic,
		UserName:    req.UserName,
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		Email:       req.Email,
		Password:    string(hashedPassword),
		PhoneNumber: req.PhoneNumber,
		RoleID:      roleID,
		BusinessID:  businessID,
	}

	if role.Value != "patient" {
		user, err := h.repo.Create(ctx, userArg)
		if err != nil {
			c.JSON(http.StatusInternalServerError, response.Error(http.StatusInternalServerError, "Error al crear usuario", err))
			return
		}
		c.JSON(http.StatusCreated, response.Created("Usuario creado", &user))
	}

	if req.PatientProfile == nil {
		c.JSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, "Perfil de paciente requerido"))
		return
	}

	birthDay, err := time.Parse("2006-01-02", req.PatientProfile.BirthDay)
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
	if err := weight.Scan(fmt.Sprintf("%g", req.PatientProfile.Weight)); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, "Formato de peso inválido", err))
		return
	}

	var height pgtype.Numeric
	if err := height.Scan(fmt.Sprintf("%g", req.PatientProfile.Height)); err != nil {
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

	user, err := qtx.CreateUser(ctx, userArg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(http.StatusInternalServerError, "Error al crear usuario", err))
		return
	}

	_, err = qtx.CreatePatientProfile(ctx, sqlc.CreatePatientProfileParams{
		BusinessID:            businessID,
		UserID:                user.ID,
		Gender:                req.PatientProfile.Gender,
		BirthDay:              pgBirthDay,
		BloodType:             req.PatientProfile.BloodType,
		Weight:                weight,
		Height:                height,
		EmergencyContactName:  req.PatientProfile.EmergencyContactName,
		EmergencyContactPhone: req.PatientProfile.EmergencyContactPhone,
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

func (h *UserHandler) GetAll(c *gin.Context) {
	users, err := h.repo.GetAll(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusNotFound, response.Error(http.StatusNotFound, "Usuarios no encontrados", err))
		return
	}

	c.JSON(http.StatusOK, response.Success("Usuarios encontrados", &users))
}

func (h *UserHandler) GetAllWithSoftDeleted(c *gin.Context) {
	users, err := h.repo.GetAllWithSoftDeleted(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusNotFound, response.Error(http.StatusNotFound, "Usuarios no encontrados", err))
		return
	}

	c.JSON(http.StatusOK, response.Success("Usuarios encontrados", &users))
}

func (h *UserHandler) GetAllByRole(c *gin.Context) {
	businessID, ok := ctxkeys.BusinessID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, response.Error(http.StatusUnauthorized, "Usuario no autenticado"))
		return
	}

	role := c.Param("role")

	rows, err := h.repo.GetAllByRole(c.Request.Context(), sqlc.GetUsersByRoleParams{
		BusinessID: businessID,
		Value:      role,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(http.StatusInternalServerError, "Usuarios no encontrados", err))
		return
	}

	users := make([]userByRoleResponse, len(rows))
	for i, row := range rows {
		users[i] = userByRoleResponse{
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
			users[i].Role = &userRole{
				ID:          row.RoleID,
				Name:        row.RoleName.String,
				Value:       row.RoleValue.String,
				Description: row.RoleDescription.String,
			}
		}
	}

	c.JSON(http.StatusOK, response.Success("Usuarios encontrados", &users))
}

func (h *UserHandler) GetAllByRoleWithSoftDeleted(c *gin.Context) {
	businessID, ok := ctxkeys.BusinessID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, response.Error(http.StatusUnauthorized, "Usuario no autenticado"))
		return
	}

	role := c.Param("role")

	rows, err := h.repo.GetAllByRoleWithSoftDeleted(c.Request.Context(), sqlc.GetUsersByRoleWithSoftDeletedParams{
		BusinessID: businessID,
		Value:      role,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(http.StatusInternalServerError, "Usuarios no encontrados", err))
		return
	}

	users := make([]userByRoleResponse, len(rows))
	for i, row := range rows {
		users[i] = userByRoleResponse{
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
			users[i].Role = &userRole{
				ID:          row.RoleID,
				Name:        row.RoleName.String,
				Value:       row.RoleValue.String,
				Description: row.RoleDescription.String,
			}
		}
	}

	c.JSON(http.StatusOK, response.Success("Usuarios encontrados", &users))
}

func (h *UserHandler) GetByID(c *gin.Context) {
	var id pgtype.UUID
	if err := id.Scan(c.Param("id")); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, "Formato de ID inválido", err))
		return
	}

	user, err := h.repo.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, response.Error(http.StatusNotFound, "Usuario no encontrado", err))
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h *UserHandler) GetByIDWithSoftDeleted(c *gin.Context) {
	businessID, ok := ctxkeys.BusinessID(c)
	if !ok {
		c.JSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, "ID de negocio inválido"))
		return
	}

	var id pgtype.UUID
	if err := id.Scan(c.Param("id")); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, "Formato de ID inválido", err))
		return
	}

	row, err := h.repo.GetByIDWithSoftDeleted(c.Request.Context(), sqlc.GetUserByIDWithSoftDeletedParams{ID: id, BusinessID: businessID})
	if err != nil {
		c.JSON(http.StatusNotFound, response.Error(http.StatusNotFound, "Usuario no encontrado", err))
		return
	}

	user := userByRoleResponse{
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

	c.JSON(http.StatusOK, response.Success("Usuario encontrado", &user))
}

func (h *UserHandler) Update(c *gin.Context) {
	var id pgtype.UUID
	if err := id.Scan(c.Param("id")); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, "Formato de ID inválido", err))
		return
	}

	var req UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, "Error de validación de datos", err))
		return
	}

	var roleID pgtype.UUID
	if err := roleID.Scan(req.RoleID); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, "Formato de ID inválido", err))
		return
	}

	current, err := h.repo.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, response.Error(http.StatusNotFound, "Usuario no encontrado", err))
		return
	}

	passwordHash := current.Password
	if req.Password != nil {
		hashed, err := bcrypt.GenerateFromPassword([]byte(*req.Password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, response.Error(http.StatusInternalServerError, "Error al procesar contraseña", err))
			return
		}
		passwordHash = string(hashed)
	}

	user, err := h.repo.Update(c.Request.Context(), sqlc.UpdateUserParams{
		ID:          id,
		Ic:          req.Ic,
		UserName:    req.UserName,
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		Email:       req.Email,
		Password:    passwordHash,
		PhoneNumber: req.PhoneNumber,
		RoleID:      roleID,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(http.StatusInternalServerError, "Error al actualizar usuario", err))
		return
	}

	c.JSON(http.StatusOK, response.Success("Usuario actualizado", &user))
}

func (h *UserHandler) Delete(c *gin.Context) {
	var id pgtype.UUID
	if err := id.Scan(c.Param("id")); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, "Formato de ID inválido", err))
		return
	}

	err := h.repo.Delete(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, "Error al eliminar usuario"))
		return
	}

	c.JSON(http.StatusOK, response.Success[any]("Usuario eliminado", nil))
}

func (h *UserHandler) SoftDelete(c *gin.Context) {
	var id pgtype.UUID
	if err := id.Scan(c.Param("id")); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, "Formato de ID inválido", err))
		return
	}

	rows, err := h.repo.SoftDelete(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, "Error al eliminar usuario"))
		return
	}
	if rows == 0 {
		c.JSON(http.StatusNotFound, response.Error(http.StatusNotFound, "Usuario no encontrado"))
		return
	}

	c.JSON(http.StatusOK, response.Success[any]("Usuario eliminado", nil))
}

func (h *UserHandler) Restore(c *gin.Context) {
	var id pgtype.UUID
	if err := id.Scan(c.Param("id")); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, "Formato de ID inválido", err))
		return
	}

	rows, err := h.repo.Restore(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, "Error al restaurar usuario"))
		return
	}
	if rows == 0 {
		c.JSON(http.StatusNotFound, response.Error(http.StatusNotFound, "Usuario no encontrado"))
		return
	}

	c.JSON(http.StatusOK, response.Success[any]("Usuario restaurado", nil))
}

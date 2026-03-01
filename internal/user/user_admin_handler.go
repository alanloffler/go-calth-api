package user

import (
	"net/http"

	"github.com/alanloffler/go-calth-api/internal/common/ctxkeys"
	"github.com/alanloffler/go-calth-api/internal/common/response"
	"github.com/alanloffler/go-calth-api/internal/database/sqlc"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
	"golang.org/x/crypto/bcrypt"
)

type CreateAdminRequest struct {
	Ic          string `json:"ic" binding:"required,len=8"`
	UserName    string `json:"userName" binding:"required,min=3,max=100"`
	FirstName   string `json:"firstName" binding:"required,min=3,max=100"`
	LastName    string `json:"lastName" binding:"required,min=3,max=100"`
	Email       string `json:"email" binding:"required,email,max=100"`
	Password    string `json:"password" binding:"required,min=8,max=100"`
	PhoneNumber string `json:"phoneNumber" binding:"required,len=10,numeric"`
}

type UpdateAdminRequest struct {
	User UpdateUserData `json:"user" binding:"required"`
}

func (h *UserHandler) CreateAdmin(c *gin.Context) {
	businessID, ok := ctxkeys.BusinessID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, response.Error(http.StatusUnauthorized, "Usuario no autenticado"))
		return
	}

	var req CreateAdminRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, "Error de validación de datos", err))
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(http.StatusInternalServerError, "Error al procesar contraseña", err))
		return
	}

	ctx := c.Request.Context()

	role, err := h.repo.q.GetRoleByValue(ctx, "admin")
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(http.StatusInternalServerError, "Rol de administrador no encontrado", err))
		return
	}

	user, err := h.repo.Create(ctx, sqlc.CreateUserParams{
		Ic:          req.Ic,
		UserName:    req.UserName,
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		Email:       req.Email,
		Password:    string(hashedPassword),
		PhoneNumber: req.PhoneNumber,
		RoleID:      role.ID,
		BusinessID:  businessID,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(http.StatusInternalServerError, "Error al crear administrador", err))
		return
	}

	c.JSON(http.StatusOK, response.Created("Administrador creado", &user))
}

func (h *UserHandler) UpdateAdmin(c *gin.Context) {
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

	var req UpdateAdminRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, "Error de validación de datos", err))
		return
	}

	current, err := h.repo.GetByID(c.Request.Context(), sqlc.GetUserByIDParams{BusinessID: businessID, ID: id})
	if err != nil {
		c.JSON(http.StatusNotFound, response.Error(http.StatusNotFound, "Administrador no encontrado", err))
		return
	}

	ic := current.Ic
	if req.User.Ic != nil {
		ic = *req.User.Ic
	}
	userName := current.UserName
	if req.User.UserName != nil {
		userName = *req.User.UserName
	}
	firstName := current.FirstName
	if req.User.FirstName != nil {
		firstName = *req.User.FirstName
	}
	lastName := current.LastName
	if req.User.LastName != nil {
		lastName = *req.User.LastName
	}
	email := current.Email
	if req.User.Email != nil {
		email = *req.User.Email
	}
	phoneNumber := current.PhoneNumber
	if req.User.PhoneNumber != nil {
		phoneNumber = *req.User.PhoneNumber
	}
	passwordHash := current.Password
	if req.User.Password != nil {
		hashed, err := bcrypt.GenerateFromPassword([]byte(*req.User.Password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, response.Error(http.StatusInternalServerError, "Error al procesar contraseña", err))
			return
		}
		passwordHash = string(hashed)
	}

	user, err := h.repo.Update(c.Request.Context(), sqlc.UpdateUserParams{
		ID:          id,
		Ic:          ic,
		UserName:    userName,
		FirstName:   firstName,
		LastName:    lastName,
		Email:       email,
		Password:    passwordHash,
		PhoneNumber: phoneNumber,
		RoleID:      current.RoleID,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(http.StatusInternalServerError, "Error al actualizar el administrador", err))
		return
	}

	c.JSON(http.StatusOK, response.Success("Administrador actualizado", &user))
}

package user

import (
	"net/http"

	"github.com/alanloffler/go-calth-api/internal/common/ctxkeys"
	"github.com/alanloffler/go-calth-api/internal/common/response"
	"github.com/alanloffler/go-calth-api/internal/database/sqlc"
	"github.com/gin-gonic/gin"
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
	h.Update(c)
}

package user

import (
	"net/http"

	"github.com/alanloffler/go-calth-api/internal/common/ctxkeys"
	"github.com/alanloffler/go-calth-api/internal/common/response"
	"github.com/alanloffler/go-calth-api/internal/common/utils"
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
	Ic          *string `json:"ic" binding:"omitempty,len=8"`
	UserName    *string `json:"userName" binding:"omitempty,min=3,max=100"`
	FirstName   *string `json:"firstName" binding:"omitempty,min=3,max=100"`
	LastName    *string `json:"lastName" binding:"omitempty,min=3,max=100"`
	Email       *string `json:"email" binding:"omitempty,email,max=100"`
	Password    *string `json:"password" binding:"omitempty,min=8,max=100"`
	PhoneNumber *string `json:"phoneNumber" binding:"omitempty,len=10,numeric"`
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

	role, err := h.repo.GetQueries().GetRoleByValue(ctx, "admin")
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

	var passwordHash pgtype.Text
	if req.Password != nil {
		hashed, err := bcrypt.GenerateFromPassword([]byte(*req.Password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, response.Error(http.StatusInternalServerError, "Error al procesar contraseña", err))
			return
		}
		passwordHash = pgtype.Text{String: string(hashed), Valid: true}
	}

	affected, err := h.repo.Update(c.Request.Context(), sqlc.UpdateUserParams{
		BusinessID:  businessID,
		ID:          id,
		Ic:          utils.ToPgText(req.Ic),
		UserName:    utils.ToPgText(req.UserName),
		FirstName:   utils.ToPgText(req.FirstName),
		LastName:    utils.ToPgText(req.LastName),
		Email:       utils.ToPgText(req.Email),
		Password:    passwordHash,
		PhoneNumber: utils.ToPgText(req.PhoneNumber),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(http.StatusInternalServerError, "Error al actualizar usuario", err))
		return
	}
	if affected == 0 {
		c.JSON(http.StatusNotFound, response.Error(http.StatusNotFound, "Usuario no encontrado"))
		return
	}

	c.JSON(http.StatusOK, response.Success[any]("Usuario actualizado", nil))
}

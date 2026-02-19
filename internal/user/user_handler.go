package user

import (
	"net/http"

	"github.com/alanloffler/go-calth-api/internal/common/response"
	"github.com/alanloffler/go-calth-api/internal/database/sqlc"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
	"golang.org/x/crypto/bcrypt"
)

type UserHandler struct {
	repo *UserRepository
}

func NewUserHandler(repo *UserRepository) *UserHandler {
	return &UserHandler{repo: repo}
}

type CreateUserRequest struct {
	Ic          string `json:"ic" binding:"required,len=8"`
	UserName    string `json:"userName" binding:"required,min=3,max=100"`
	FirstName   string `json:"firstName" binding:"required,min=3,max=100"`
	LastName    string `json:"lastName" binding:"required,min=3,max=100"`
	Email       string `json:"email" binding:"required,email,max=100"`
	Password    string `json:"password" binding:"required,min=8,max=100"`
	PhoneNumber string `json:"phoneNumber" binding:"required,len=10,numeric"`
	RoleID      string `json:"roleId" binding:"required,uuid"`
	BusinessID  string `json:"businessId" binding:"required,uuid"`
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

	user, err := h.repo.Create(c.Request.Context(), sqlc.CreateUserParams{
		Ic:          req.Ic,
		UserName:    req.UserName,
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		Email:       req.Email,
		Password:    string(hashedPassword),
		PhoneNumber: req.PhoneNumber,
		RoleID:      roleID,
		BusinessID:  businessID,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(http.StatusInternalServerError, "Error al crear usuario", err))
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
	var id pgtype.UUID
	if err := id.Scan(c.Param("id")); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, "Formato de ID inválido", err))
		return
	}

	user, err := h.repo.GetByIDWithSoftDeleted(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, response.Error(http.StatusNotFound, "Usuario no encontrado", err))
		return
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

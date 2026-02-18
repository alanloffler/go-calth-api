package auth

import (
	"errors"
	"net/http"

	"github.com/alanloffler/go-calth-api/internal/common/response"
	"github.com/alanloffler/go-calth-api/internal/database/sqlc"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	service *AuthService
	queries *sqlc.Queries
}

func NewAuthHandler(service *AuthService, queries *sqlc.Queries) *AuthHandler {
	return &AuthHandler{service: service, queries: queries}
}

type LoginRequest struct {
	BusinessID string `json:"businessID" binding:"required,uuid"`
	Email      string `json:"email" binding:"required,email"`
	Password   string `json:"password" binding:"required"`
}

type LogoutRequest struct {
	RefreshToken string `json:"refreshToken" binding:"required"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refreshToken" binding:"required"`
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindBodyWithJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, "Error de validación de datos", err))
		return
	}

	var businessID pgtype.UUID
	if err := businessID.Scan(req.BusinessID); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, "Formato de ID inválido", err))
		return
	}

	user, err := h.queries.GetUserByEmail(c.Request.Context(), sqlc.GetUserByEmailParams{
		BusinessID: businessID,
		Email:      req.Email,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.JSON(http.StatusUnauthorized, response.Error(http.StatusUnauthorized, "Credenciales inválidas", err))
			return
		}
		c.JSON(http.StatusInternalServerError, response.Error(http.StatusInternalServerError, "Error al buscar usuario", err))
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, response.Error(http.StatusUnauthorized, "Credenciales inválidas"))
		return
	}

	tokenPair, err := h.service.GenerateTokenPair(user.ID.String(), user.BusinessID.String(), user.RoleID.String())
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(http.StatusInternalServerError, "Error al generar tokens", err))
		return
	}

	_, err = h.queries.UpdateRefreshToken(c.Request.Context(), sqlc.UpdateRefreshTokenParams{
		ID:           user.ID,
		RefreshToken: pgtype.Text{String: tokenPair.RefreshToken, Valid: true},
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(http.StatusInternalServerError, "Error al guardar token de refresco", err))
		return
	}

	c.JSON(http.StatusOK, response.Success("Inicio de sesión exitoso", tokenPair))
}

func (h *AuthHandler) Logout(c *gin.Context) {
	var req LogoutRequest
	if err := c.ShouldBindBodyWithJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, "Error de validación de datos", err))
		return
	}

	claims, err := h.service.ValidateRefreshToken(req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, response.Error(http.StatusUnauthorized, "Token inválido o expirado", err))
		return
	}

	var userID pgtype.UUID
	if err := userID.Scan(claims.UserID); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, "ID de usuario inválido"))
		return
	}

	_, err = h.queries.ClearRefreshToken(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(http.StatusInternalServerError, "Error al cerrar sesión", err))
		return
	}

	c.JSON(http.StatusOK, response.Success[any]("Sesión cerrada exitosamente", nil))
}

func (h *AuthHandler) Refresh(c *gin.Context) {
	var req RefreshRequest
	if err := c.ShouldBindBodyWithJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, "Error de validación de datos", err))
		return
	}

	claims, err := h.service.ValidateRefreshToken(req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, response.Error(http.StatusUnauthorized, "Token inválido o expirado", err))
		return
	}

	var userID pgtype.UUID
	if err := userID.Scan(claims.UserID); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, "ID de usuario inválido"))
		return
	}

	user, err := h.queries.GetUserByID(c.Request.Context(), userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.JSON(http.StatusUnauthorized, response.Error(http.StatusUnauthorized, "Usuario no encontrado", err))
			return
		}
		c.JSON(http.StatusInternalServerError, response.Error(http.StatusInternalServerError, "Error al buscar usuario", err))
		return
	}

	if !user.RefreshToken.Valid || user.RefreshToken.String != req.RefreshToken {
		c.JSON(http.StatusUnauthorized, response.Error(http.StatusUnauthorized, "Token de refresco inválido", err))
		return
	}

	tokenPair, err := h.service.GenerateTokenPair(userID.String(), user.BusinessID.String(), user.RoleID.String())
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(http.StatusInternalServerError, "Error al generar tokens", err))
		return
	}

	_, err = h.queries.UpdateRefreshToken(c.Request.Context(), sqlc.UpdateRefreshTokenParams{
		ID:           user.ID,
		RefreshToken: pgtype.Text{String: tokenPair.RefreshToken, Valid: true},
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(http.StatusInternalServerError, "Error al guardar token", err))
		return
	}

	c.JSON(http.StatusOK, response.Success("Token refrescado exitosamente", tokenPair))
}

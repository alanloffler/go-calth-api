package auth

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/alanloffler/go-calth-api/internal/common/ctxkeys"
	"github.com/alanloffler/go-calth-api/internal/common/response"
	"github.com/alanloffler/go-calth-api/internal/config"
	"github.com/alanloffler/go-calth-api/internal/database/sqlc"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"golang.org/x/crypto/bcrypt"
)

type getMePermission struct {
	ID        pgtype.UUID `json:"id"`
	ActionKey string      `json:"actionKey"`
}

type getMeRolePermission struct {
	RoleID       pgtype.UUID     `json:"roleId"`
	PermissionID pgtype.UUID     `json:"permissionId"`
	Permission   getMePermission `json:"permission"`
}

type getMeRole struct {
	ID              pgtype.UUID           `json:"id"`
	Name            string                `json:"name"`
	Value           string                `json:"value"`
	RolePermissions []getMeRolePermission `json:"rolePermissions"`
}

type getMeResponse struct {
	ID          pgtype.UUID        `json:"id"`
	Ic          string             `json:"ic"`
	UserName    string             `json:"userName"`
	FirstName   string             `json:"firstName"`
	LastName    string             `json:"lastName"`
	Email       string             `json:"email"`
	PhoneNumber string             `json:"phoneNumber"`
	RoleID      pgtype.UUID        `json:"roleId"`
	BusinessID  pgtype.UUID        `json:"businessId"`
	CreatedAt   pgtype.Timestamptz `json:"createdAt"`
	UpdatedAt   pgtype.Timestamptz `json:"updatedAt"`
	Role        *getMeRole         `json:"role"`
}

type AuthHandler struct {
	cfg     *config.Config
	repo    *AuthRepository
	service *AuthService
}

func NewAuthHandler(cfg *config.Config, repo *AuthRepository, service *AuthService) *AuthHandler {
	return &AuthHandler{cfg: cfg, repo: repo, service: service}
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindBodyWithJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, "Error de validación de datos", err))
		return
	}

	origin := c.GetHeader("Origin")
	if origin == "" {
		c.JSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, "Origen requerido"))
		return
	}

	slug, err := extractSubdomain(origin)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, "Subdominio inválido", err))
		return
	}

	business, err := h.repo.GetBusinessBySlug(c.Request.Context(), slug)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.JSON(http.StatusUnauthorized, response.Error(http.StatusUnauthorized, "Negocio no encontrado"))
			return
		}
		c.JSON(http.StatusInternalServerError, response.Error(http.StatusInternalServerError, "Error al buscar negocio", err))
		return
	}

	user, err := h.repo.GetUserByEmail(c.Request.Context(), sqlc.GetUserByEmailParams{
		BusinessID: business.ID,
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

	_, err = h.repo.UpdateRefreshToken(c.Request.Context(), sqlc.UpdateRefreshTokenParams{
		ID:           user.ID,
		RefreshToken: pgtype.Text{String: tokenPair.RefreshToken, Valid: true},
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(http.StatusInternalServerError, "Error al guardar token de refresco", err))
		return
	}

	accessMaxAge := parseDurationToSeconds(h.cfg.JwtAccessExpiry)
	refreshMaxAge := parseDurationToSeconds(h.cfg.JwtRefreshExpiry)

	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("access_token", tokenPair.AccessToken, accessMaxAge, "/", h.cfg.CookieDomain, h.cfg.CookieSecure, true)
	c.SetCookie("refresh_token", tokenPair.RefreshToken, refreshMaxAge, "/", h.cfg.CookieDomain, h.cfg.CookieSecure, true)

	c.JSON(http.StatusOK, response.Success[any]("Inicio de sesión exitoso", nil))
}

func (h *AuthHandler) Logout(c *gin.Context) {
	refreshToken, err := c.Cookie("refresh_token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, response.Error(http.StatusUnauthorized, "Token requerido"))
		return
	}

	claims, err := h.service.ValidateRefreshToken(refreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, response.Error(http.StatusUnauthorized, "Token inválido o expirado", err))
		return
	}

	var userID pgtype.UUID
	if err := userID.Scan(claims.UserID); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, "ID de usuario inválido"))
		return
	}

	_, err = h.repo.ClearRefreshToken(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(http.StatusInternalServerError, "Error al cerrar sesión", err))
		return
	}

	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("access_token", "", -1, "/", h.cfg.CookieDomain, h.cfg.CookieSecure, true)
	c.SetCookie("refresh_token", "", -1, "/", h.cfg.CookieDomain, h.cfg.CookieSecure, true)

	c.JSON(http.StatusOK, response.Success[any]("Sesión cerrada exitosamente", nil))
}

func (h *AuthHandler) Refresh(c *gin.Context) {
	refreshToken, err := c.Cookie("refresh_token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, response.Error(http.StatusUnauthorized, "Token requerido"))
		return
	}

	claims, err := h.service.ValidateRefreshToken(refreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, response.Error(http.StatusUnauthorized, "Token inválido o expirado", err))
		return
	}

	var userID pgtype.UUID
	if err := userID.Scan(claims.UserID); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, "ID de usuario inválido"))
		return
	}

	user, err := h.repo.GetUserByID(c.Request.Context(), userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.JSON(http.StatusUnauthorized, response.Error(http.StatusUnauthorized, "Usuario no encontrado", err))
			return
		}
		c.JSON(http.StatusInternalServerError, response.Error(http.StatusInternalServerError, "Error al buscar usuario", err))
		return
	}

	if !user.RefreshToken.Valid || user.RefreshToken.String != refreshToken {
		c.JSON(http.StatusUnauthorized, response.Error(http.StatusUnauthorized, "Token de refresco inválido"))
		return
	}

	tokenPair, err := h.service.GenerateTokenPair(userID.String(), user.BusinessID.String(), user.RoleID.String())
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(http.StatusInternalServerError, "Error al generar tokens", err))
		return
	}

	_, err = h.repo.UpdateRefreshToken(c.Request.Context(), sqlc.UpdateRefreshTokenParams{
		ID:           userID,
		RefreshToken: pgtype.Text{String: tokenPair.RefreshToken, Valid: true},
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(http.StatusInternalServerError, "Error al guardar token", err))
		return
	}

	accessMaxAge := parseDurationToSeconds(h.cfg.JwtAccessExpiry)
	refreshMaxAge := parseDurationToSeconds(h.cfg.JwtRefreshExpiry)

	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("access_token", tokenPair.AccessToken, accessMaxAge, "/", h.cfg.CookieDomain, h.cfg.CookieSecure, true)
	c.SetCookie("refresh_token", tokenPair.RefreshToken, refreshMaxAge, "/", h.cfg.CookieDomain, h.cfg.CookieSecure, true)

	c.JSON(http.StatusOK, response.Success[any]("Token refrescado exitosamente", nil))
}

func (h *AuthHandler) GetMe(c *gin.Context) {
	businessID, ok := ctxkeys.BusinessID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, response.Error(http.StatusUnauthorized, "Usuario no autenticado"))
		return
	}

	userID, ok := ctxkeys.UserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, response.Error(http.StatusUnauthorized, "Usuario no autenticado"))
		return
	}

	rows, err := h.repo.GetMe(c.Request.Context(), sqlc.GetMeParams{
		ID:         userID,
		BusinessID: businessID,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(http.StatusInternalServerError, "Error al buscar usuario", err))
		return
	}
	if len(rows) == 0 {
		c.JSON(http.StatusNotFound, response.Error(http.StatusNotFound, "Usuario no encontrado"))
		return
	}

	first := rows[0]
	result := getMeResponse{
		ID:          first.ID,
		Ic:          first.Ic,
		UserName:    first.UserName,
		FirstName:   first.FirstName,
		LastName:    first.LastName,
		Email:       first.Email,
		PhoneNumber: first.PhoneNumber,
		RoleID:      first.RoleID,
		BusinessID:  first.BusinessID,
		CreatedAt:   first.CreatedAt,
		UpdatedAt:   first.UpdatedAt,
	}

	if first.RoleID_2.Valid {
		role := getMeRole{
			ID:              first.RoleID_2,
			Name:            first.RoleName.String,
			Value:           first.RoleValue.String,
			RolePermissions: []getMeRolePermission{},
		}
		for _, row := range rows {
			if !row.RpPermissionID.Valid {
				continue
			}
			role.RolePermissions = append(role.RolePermissions, getMeRolePermission{
				RoleID:       row.RpRoleID,
				PermissionID: row.RpPermissionID,
				Permission: getMePermission{
					ID:        row.PermissionID,
					ActionKey: row.PermissionActionKey.String,
				},
			})
		}
		result.Role = &role
	}

	c.JSON(http.StatusOK, response.Success("Usuario encontrado", &result))
}

// Helpers
func extractSubdomain(origin string) (string, error) {
	host := origin
	if _, after, ok := strings.Cut(origin, "://"); ok {
		host = after
	}

	if idx := strings.LastIndex(host, ":"); idx != -1 {
		host = host[:idx]
	}

	parts := strings.Split(host, ".")
	if len(parts) < 2 {
		return "", errors.New("Subdominio no encontrado")
	}

	return parts[0], nil
}

func parseDurationToSeconds(expiry string) int {
	d, err := time.ParseDuration(expiry)
	if err != nil {
		return 0
	}
	return int(d.Seconds())
}

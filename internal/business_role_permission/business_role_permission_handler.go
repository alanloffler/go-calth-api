package business_role_permission

import (
	"net/http"

	"github.com/alanloffler/go-calth-api/internal/common/ctxkeys"
	"github.com/alanloffler/go-calth-api/internal/common/response"
	"github.com/alanloffler/go-calth-api/internal/database/sqlc"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
)

type BusinessRolePermissionHandler struct {
	repo *BusinessRolePermissionRepository
}

type UpsertOverrideRequest struct {
	Effect string `json:"effect" binding:"required,oneof=grant deny"`
}

func NewBusinessRolePermissionHandler(repo *BusinessRolePermissionRepository) *BusinessRolePermissionHandler {
	return &BusinessRolePermissionHandler{repo: repo}
}

func (h *BusinessRolePermissionHandler) ListEffective(c *gin.Context) {
	businessID, ok := ctxkeys.BusinessID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, response.Error(http.StatusUnauthorized, "Usuario no autenticado"))
		return
	}

	var roleID pgtype.UUID
	if err := roleID.Scan(c.Param("roleId")); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, "Formato de ID de rol inválido", err))
		return
	}

	rows, err := h.repo.ListEffective(c.Request.Context(), sqlc.ListEffectivePermissionsParams{
		BusinessID: businessID,
		RoleID:     roleID,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(http.StatusInternalServerError, "Error al obtener permisos efectivos", err))
		return
	}

	c.JSON(http.StatusOK, response.Success("Permisos efectivos obtenidos", &rows))
}

func (h *BusinessRolePermissionHandler) ListOverrides(c *gin.Context) {
	businessID, ok := ctxkeys.BusinessID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, response.Error(http.StatusUnauthorized, "Usuario no autenticado"))
		return
	}

	var roleID pgtype.UUID
	if err := roleID.Scan(c.Param("roleId")); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, "Formato de ID de rol inválido", err))
		return
	}

	rows, err := h.repo.ListOverrides(c.Request.Context(), sqlc.GetBusinessRoleOverridesParams{
		BusinessID: businessID,
		RoleID:     roleID,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(http.StatusInternalServerError, "Error al obtener overrides", err))
		return
	}

	c.JSON(http.StatusOK, response.Success("Overrides encontrados", &rows))
}

func (h *BusinessRolePermissionHandler) Upsert(c *gin.Context) {
	businessID, ok := ctxkeys.BusinessID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, response.Error(http.StatusUnauthorized, "Usuario no autenticado"))
		return
	}

	var roleID pgtype.UUID
	if err := roleID.Scan(c.Param("roleId")); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, "Formato de ID de rol inválido", err))
		return
	}

	var permissionID pgtype.UUID
	if err := permissionID.Scan(c.Param("permissionId")); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, "Formato de ID de permiso inválido", err))
		return
	}

	var req UpsertOverrideRequest
	if err := c.ShouldBindBodyWithJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, "Error de validación de datos", err))
		return
	}

	row, err := h.repo.Upsert(c.Request.Context(), sqlc.UpsertBusinessRolePermissionParams{
		BusinessID:   businessID,
		RoleID:       roleID,
		PermissionID: permissionID,
		Effect:       req.Effect,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(http.StatusInternalServerError, "Error al guardar override", err))
		return
	}

	c.JSON(http.StatusOK, response.Success("Override guardado", &row))
}

func (h *BusinessRolePermissionHandler) ResetOne(c *gin.Context) {
	businessID, ok := ctxkeys.BusinessID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, response.Error(http.StatusUnauthorized, "Usuario no autenticado"))
		return
	}

	var roleID pgtype.UUID
	if err := roleID.Scan(c.Param("roleId")); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, "Formato de ID de rol inválido", err))
		return
	}

	var permissionID pgtype.UUID
	if err := permissionID.Scan(c.Param("permissionId")); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, "Formato de ID de permiso inválido", err))
		return
	}

	affected, err := h.repo.DeleteOne(c.Request.Context(), sqlc.DeleteBusinessRolePermissionParams{
		BusinessID:   businessID,
		RoleID:       roleID,
		PermissionID: permissionID,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(http.StatusInternalServerError, "Error al eliminar override", err))
		return
	}
	if affected == 0 {
		c.JSON(http.StatusNotFound, response.Error(http.StatusNotFound, "Override no encontrado"))
		return
	}

	c.JSON(http.StatusOK, response.Success[any]("Override eliminado", nil))
}

func (h *BusinessRolePermissionHandler) ResetAll(c *gin.Context) {
	businessID, ok := ctxkeys.BusinessID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, response.Error(http.StatusUnauthorized, "Usuario no autenticado"))
		return
	}

	var roleID pgtype.UUID
	if err := roleID.Scan(c.Param("roleId")); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, "Formato de ID de rol inválido", err))
		return
	}

	_, err := h.repo.DeleteAll(c.Request.Context(), sqlc.DeleteBusinessRoleOverridesParams{
		BusinessID: businessID,
		RoleID:     roleID,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(http.StatusInternalServerError, "Error al eliminar overrides", err))
		return
	}

	c.JSON(http.StatusOK, response.Success[any]("Overrides eliminados", nil))
}

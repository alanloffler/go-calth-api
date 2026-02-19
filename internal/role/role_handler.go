package role

import (
	"net/http"

	"github.com/alanloffler/go-calth-api/internal/common/response"
	"github.com/alanloffler/go-calth-api/internal/common/utils"
	"github.com/alanloffler/go-calth-api/internal/database/sqlc"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type RoleHandler struct {
	repo *RoleRepository
	pool *pgxpool.Pool
}

func NewRoleHandler(repo *RoleRepository, pool *pgxpool.Pool) *RoleHandler {
	return &RoleHandler{repo: repo, pool: pool}
}

type PermissionAction struct {
	ID    string `json:"id" binding:"required,uuid"`
	Value bool   `json:"value"`
}

type PermissionGroup struct {
	Actions []PermissionAction `json:"actions" binding:"required,dive"`
}

type CreateRoleRequest struct {
	Name        string            `json:"name" binding:"required,min=3,max=100"`
	Value       string            `json:"value" binding:"required,min=3,max=100"`
	Description string            `json:"description" binding:"required,min=3,max=100"`
	Permissions []PermissionGroup `json:"permissions"`
}

type UpdateRoleRequest struct {
	Name        *string           `json:"name" binding:"omitempty,min=3,max=100"`
	Value       *string           `json:"value" binding:"omitempty,min=3,max=100"`
	Description *string           `json:"description" binding:"omitempty,min=3,max=100"`
	Permissions []PermissionGroup `json:"permissions"`
}

type RoleWithPermissions struct {
	ID              pgtype.UUID            `json:"id"`
	Name            string                 `json:"name"`
	Value           string                 `json:"value"`
	Description     string                 `json:"description"`
	RolePermissions []RolePermissionDetail `json:"rolePermissions"`
	CreatedAt       pgtype.Timestamptz     `json:"createdAt"`
	UpdatedAt       pgtype.Timestamptz     `json:"updatedAt"`
	DeletedAt       pgtype.Timestamptz     `json:"deletedAt"`
}

type RolePermissionDetail struct {
	RoleID       pgtype.UUID        `json:"roleId"`
	PermissionID pgtype.UUID        `json:"permissionId"`
	CreatedAt    pgtype.Timestamptz `json:"createdAt"`
	UpdatedAt    pgtype.Timestamptz `json:"updatedAt"`
	Permission   PermissionDetail   `json:"permission"`
}

type PermissionDetail struct {
	ID          pgtype.UUID        `json:"id"`
	Name        string             `json:"name"`
	Category    string             `json:"category"`
	ActionKey   string             `json:"actionKey"`
	Description string             `json:"description"`
	CreatedAt   pgtype.Timestamptz `json:"createdAt"`
	UpdatedAt   pgtype.Timestamptz `json:"updatedAt"`
	DeletedAt   pgtype.Timestamptz `json:"deletedAt"`
}

func (h *RoleHandler) Create(c *gin.Context) {
	var req CreateRoleRequest
	if err := c.ShouldBindBodyWithJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, "Error de validación de datos", err))
		return
	}

	ctx := c.Request.Context()

	_, err := h.repo.q.GetRoleByValue(ctx, req.Value)
	if err == nil {
		c.JSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, "El rol ya existe"))
		return
	}

	// Enabled permission IDs not duplicated
	seen := make(map[string]bool)
	var permissionIDs []pgtype.UUID
	for _, group := range req.Permissions {
		for _, action := range group.Actions {
			if action.Value && !seen[action.ID] {
				seen[action.ID] = true
				parsed, err := uuid.Parse(action.ID)
				if err != nil {
					c.JSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, "ID de permiso inválido", err))
					return
				}

				permissionIDs = append(permissionIDs, pgtype.UUID{Bytes: parsed, Valid: true})
			}
		}
	}

	// No permissions, create without transaction
	if len(permissionIDs) == 0 {
		role, err := h.repo.Create(ctx, sqlc.CreateRoleParams{
			Name:        req.Name,
			Value:       req.Value,
			Description: req.Description,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, response.Error(http.StatusInternalServerError, "Error al crear rol", err))
			return
		}

		c.JSON(http.StatusOK, response.Success("Rol creado", &role))
		return
	}

	// With permissions, use transaction
	tx, err := h.pool.Begin(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(http.StatusInternalServerError, "Error al iniciar transacción", err))
		return
	}
	defer tx.Rollback(ctx)

	qtx := sqlc.New(tx)

	role, err := qtx.CreateRole(ctx, sqlc.CreateRoleParams{
		Name:        req.Name,
		Value:       req.Value,
		Description: req.Description,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(http.StatusInternalServerError, "Error al crear rol", err))
		return
	}

	for _, permID := range permissionIDs {
		_, err := qtx.CreateRolePermission(ctx, sqlc.CreateRolePermissionParams{
			RoleID:       role.ID,
			PermissionID: permID,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, response.Error(http.StatusInternalServerError, "Error al asignar permiso", err))
			return
		}
	}

	if err := tx.Commit(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(http.StatusInternalServerError, "Error al confirmar transacción", err))
		return
	}

	c.JSON(http.StatusCreated, response.Created("Rol creado", &role))
}

func (h *RoleHandler) GetAll(c *gin.Context) {
	permissions, err := h.repo.GetAll(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusNotFound, response.Error(http.StatusNotFound, "Roles no encontrados", err))
		return
	}

	c.JSON(http.StatusOK, response.Success("Roles encontrados", &permissions))
}

func (h *RoleHandler) GetAllWithSoftDeleted(c *gin.Context) {
	permissions, err := h.repo.GetAllWithSoftDeleted(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusNotFound, response.Error(http.StatusNotFound, "Roles no encontrados", err))
		return
	}

	c.JSON(http.StatusOK, response.Success("Roles encontrados", &permissions))
}

func (h *RoleHandler) GetOneByID(c *gin.Context) {
	var id pgtype.UUID
	if err := id.Scan(c.Param("id")); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, "Formato de ID inválido", err))
		return
	}

	rows, err := h.repo.GetOneByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, response.Error(http.StatusNotFound, "Rol no encontrado", err))
		return
	}

	if len(rows) == 0 {
		c.JSON(http.StatusNotFound, response.Error(http.StatusNotFound, "Rol no encontrado", err))
		return
	}

	role := RoleWithPermissions{
		ID:              rows[0].ID,
		Name:            rows[0].Name,
		Value:           rows[0].Value,
		Description:     rows[0].Description,
		CreatedAt:       rows[0].CreatedAt,
		UpdatedAt:       rows[0].UpdatedAt,
		DeletedAt:       rows[0].DeletedAt,
		RolePermissions: []RolePermissionDetail{},
	}

	for _, row := range rows {
		if row.RoleID.Valid {
			role.RolePermissions = append(role.RolePermissions, RolePermissionDetail{
				RoleID:       row.RoleID,
				PermissionID: row.PermissionID,
				CreatedAt:    row.RpCreatedAt,
				UpdatedAt:    row.RpUpdatedAt,
				Permission: PermissionDetail{
					ID:          row.PID,
					Name:        row.PName.String,
					Category:    row.PCategory.String,
					ActionKey:   row.PActionKey.String,
					Description: row.PDescription.String,
					CreatedAt:   row.PCreatedAt,
					UpdatedAt:   row.PUpdatedAt,
					DeletedAt:   row.PDeletedAt,
				},
			})
		}
	}

	c.JSON(http.StatusOK, response.Success("Rol encontrado", &role))
}

func (h *RoleHandler) GetOneByIDWithSoftDeleted(c *gin.Context) {
	var id pgtype.UUID
	if err := id.Scan(c.Param("id")); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, "Formato de ID inválido", err))
		return
	}

	rows, err := h.repo.GetOneByIDWithSoftDeleted(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, response.Error(http.StatusNotFound, "Rol no encontrado", err))
		return
	}

	if len(rows) == 0 {
		c.JSON(http.StatusNotFound, response.Error(http.StatusNotFound, "Rol no encontrado", err))
		return
	}

	role := RoleWithPermissions{
		ID:              rows[0].ID,
		Name:            rows[0].Name,
		Value:           rows[0].Value,
		Description:     rows[0].Description,
		CreatedAt:       rows[0].CreatedAt,
		UpdatedAt:       rows[0].UpdatedAt,
		DeletedAt:       rows[0].DeletedAt,
		RolePermissions: []RolePermissionDetail{},
	}

	for _, row := range rows {
		if row.RoleID.Valid {
			role.RolePermissions = append(role.RolePermissions, RolePermissionDetail{
				RoleID:       row.RoleID,
				PermissionID: row.PermissionID,
				CreatedAt:    row.RpCreatedAt,
				UpdatedAt:    row.RpUpdatedAt,
				Permission: PermissionDetail{
					ID:          row.PID,
					Name:        row.PName.String,
					Category:    row.PCategory.String,
					ActionKey:   row.PActionKey.String,
					Description: row.PDescription.String,
					CreatedAt:   row.PCreatedAt,
					UpdatedAt:   row.PUpdatedAt,
					DeletedAt:   row.PDeletedAt,
				},
			})
		}
	}

	c.JSON(http.StatusOK, response.Success("Rol encontrado", &role))
}

func (h *RoleHandler) Update(c *gin.Context) {
	var id pgtype.UUID
	if err := id.Scan(c.Param("id")); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, "Formato de ID inválido", err))
		return
	}

	var req UpdateRoleRequest
	if err := c.ShouldBindBodyWithJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, "Error al actualizar rol", err))
		return
	}

	ctx := c.Request.Context()

	if req.Value != nil {
		existing, err := h.repo.GetOneByValue(ctx, *req.Value)
		if err == nil && existing.ID != id {
			c.JSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, "El rol ya existe"))
			return
		}
	}

	seen := make(map[string]bool)
	var permissionIDs []pgtype.UUID
	for _, group := range req.Permissions {
		for _, action := range group.Actions {
			if action.Value && !seen[action.ID] {
				seen[action.ID] = true
				parsed, err := uuid.Parse(action.ID)
				if err != nil {
					c.JSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, "ID de permiso inválido", err))
					return
				}

				permissionIDs = append(permissionIDs, pgtype.UUID{Bytes: parsed, Valid: true})
			}
		}
	}

	tx, err := h.pool.Begin(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(http.StatusInternalServerError, "Error al iniciar transacción", err))
		return
	}
	defer tx.Rollback(ctx)

	qtx := sqlc.New(tx)

	_, err = qtx.UpdateRole(ctx, sqlc.UpdateRoleParams{
		ID:          id,
		Name:        utils.ToPgText(req.Name),
		Value:       utils.ToPgText(req.Value),
		Description: utils.ToPgText(req.Description),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(http.StatusInternalServerError, "Error al actualizar rol", err))
		return
	}

	err = qtx.DeleteRolePermissionsByRoleID(ctx, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(http.StatusInternalServerError, "Error al eliminar permisos del rol", err))
		return
	}

	for _, permID := range permissionIDs {
		_, err := qtx.CreateRolePermission(ctx, sqlc.CreateRolePermissionParams{
			RoleID:       id,
			PermissionID: permID,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, response.Error(http.StatusInternalServerError, "Error al asignar permiso", err))
			return
		}
	}

	if err := tx.Commit(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(http.StatusInternalServerError, "Error al confirmar la transacción", err))
		return
	}

	rows, err := h.repo.GetOneByID(ctx, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(http.StatusInternalServerError, "Error al obtener rol actualizado", err))
		return
	}

	role := RoleWithPermissions{
		ID:              rows[0].ID,
		Name:            rows[0].Name,
		Value:           rows[0].Value,
		Description:     rows[0].Description,
		CreatedAt:       rows[0].CreatedAt,
		UpdatedAt:       rows[0].UpdatedAt,
		DeletedAt:       rows[0].DeletedAt,
		RolePermissions: []RolePermissionDetail{},
	}

	for _, row := range rows {
		if row.RoleID.Valid {
			role.RolePermissions = append(role.RolePermissions, RolePermissionDetail{
				RoleID:       row.RoleID,
				PermissionID: row.PermissionID,
				CreatedAt:    row.RpCreatedAt,
				UpdatedAt:    row.RpUpdatedAt,
				Permission: PermissionDetail{
					ID:          row.PID,
					Name:        row.PName.String,
					Category:    row.PCategory.String,
					ActionKey:   row.PActionKey.String,
					Description: row.PDescription.String,
					CreatedAt:   row.PCreatedAt,
					UpdatedAt:   row.PUpdatedAt,
					DeletedAt:   row.PDeletedAt,
				},
			})
		}
	}

	c.JSON(http.StatusOK, response.Success("Rol actualizado", &role))
}

func (h *RoleHandler) Delete(c *gin.Context) {
	var id pgtype.UUID
	if err := id.Scan(c.Param("id")); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, "Formato de ID inválido", err))
		return
	}

	rows, err := h.repo.Delete(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(http.StatusInternalServerError, "Error al eliminar rol", err))
		return
	}

	if rows == 0 {
		c.JSON(http.StatusNotFound, response.Error(http.StatusNotFound, "Rol no encontrado"))
		return
	}

	c.JSON(http.StatusOK, response.Success[any]("Rol eliminado", nil))
}

func (h *RoleHandler) SoftDelete(c *gin.Context) {
	var id pgtype.UUID
	if err := id.Scan(c.Param("id")); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, "Formato de ID inválido", err))
		return
	}

	rows, err := h.repo.SoftDelete(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, "Error al eliminar rol"))
		return
	}
	if rows == 0 {
		c.JSON(http.StatusNotFound, response.Error(http.StatusNotFound, "Role no encontrado"))
		return
	}

	c.JSON(http.StatusOK, response.Success[any]("Rol eliminado", nil))
}

func (h *RoleHandler) Restore(c *gin.Context) {
	var id pgtype.UUID
	if err := id.Scan(c.Param("id")); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, "Formato de ID inválido", err))
		return
	}

	rows, err := h.repo.Restore(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(http.StatusInternalServerError, "Error al restaurar rol", err))
		return
	}
	if rows == 0 {
		c.JSON(http.StatusNotFound, response.Error(http.StatusNotFound, "Rol no encontrado"))
		return
	}

	c.JSON(http.StatusOK, response.Success[any]("Rol restaurado", nil))
}

package user

import (
	"net/http"
	"net/mail"
	"strings"

	"github.com/alanloffler/go-calth-api/internal/common/ctxkeys"
	"github.com/alanloffler/go-calth-api/internal/common/response"
	"github.com/alanloffler/go-calth-api/internal/common/utils"
	"github.com/alanloffler/go-calth-api/internal/database/sqlc"
	"github.com/alanloffler/go-calth-api/internal/patient_profile"
	"github.com/alanloffler/go-calth-api/internal/professional_profile"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

type UserHandler struct {
	repo                    UserRepositoryInterface
	pool                    *pgxpool.Pool
	patientProfileRepo      *patient_profile.PatientProfileRepository
	professionalProfileRepo *professional_profile.ProfessionalProfileRepository
}

func NewUserHandler(
	repo UserRepositoryInterface,
	pool *pgxpool.Pool,
	patientProfileRepo *patient_profile.PatientProfileRepository,
	professionalProfileRepo *professional_profile.ProfessionalProfileRepository,
) *UserHandler {
	return &UserHandler{repo: repo, pool: pool, patientProfileRepo: patientProfileRepo, professionalProfileRepo: professionalProfileRepo}
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

type UpdateUserData struct {
	Ic          *string `json:"ic" binding:"omitempty,len=8"`
	UserName    *string `json:"userName" binding:"omitempty,min=3,max=100"`
	FirstName   *string `json:"firstName" binding:"omitempty,min=3,max=100"`
	LastName    *string `json:"lastName" binding:"omitempty,min=3,max=100"`
	Email       *string `json:"email" binding:"omitempty,email,max=100"`
	Password    *string `json:"password" binding:"omitempty,min=8,max=100"`
	PhoneNumber *string `json:"phoneNumber" binding:"omitempty,len=10,numeric"`
}

type CreateUserData struct {
	Ic          string `json:"ic" binding:"required,len=8"`
	UserName    string `json:"userName" binding:"required,min=3,max=100"`
	FirstName   string `json:"firstName" binding:"required,min=3,max=100"`
	LastName    string `json:"lastName" binding:"required,min=3,max=100"`
	Email       string `json:"email" binding:"required,email,max=100"`
	Password    string `json:"password" binding:"required,min=8,max=100"`
	PhoneNumber string `json:"phoneNumber" binding:"required,len=10,numeric"`
}

type UpdateRequest struct {
	User UpdateUserData `json:"user" binding:"required"`
}

type userRole struct {
	ID          pgtype.UUID `json:"id"`
	Name        string      `json:"name"`
	Value       string      `json:"value"`
	Description string      `json:"description"`
}

type userByRoleResponse struct {
	ID                 pgtype.UUID                `json:"id"`
	Ic                 string                     `json:"ic"`
	UserName           string                     `json:"userName"`
	FirstName          string                     `json:"firstName"`
	LastName           string                     `json:"lastName"`
	Email              string                     `json:"email"`
	PhoneNumber        string                     `json:"phoneNumber"`
	Role               *userRole                  `json:"role"`
	ProfessionaProfile *professionalProfileInList `json:"professionalProfile"`
	BusinessID         pgtype.UUID                `json:"businessId"`
	CreatedAt          pgtype.Timestamptz         `json:"createdAt"`
	UpdatedAt          pgtype.Timestamptz         `json:"updatedAt"`
	DeletedAt          pgtype.Timestamptz         `json:"deletedAt"`
}

type professionalProfileInList struct {
	ProfessionalPrefix string `json:"professionalPrefix"`
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
		if row.ProfessionalPrefix.Valid {
			users[i].ProfessionaProfile = &professionalProfileInList{
				ProfessionalPrefix: row.ProfessionalPrefix.String,
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

	user, err := h.repo.GetByID(c.Request.Context(), sqlc.GetUserByIDParams{BusinessID: businessID, ID: id})
	if err != nil {
		c.JSON(http.StatusNotFound, response.Error(http.StatusNotFound, "Usuario no encontrado", err))
		return
	}

	c.JSON(http.StatusOK, response.Success("Usuario encontrado", &user))
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

	var req UpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, "Error de validación de datos", err))
		return
	}

	var passwordHash pgtype.Text
	if req.User.Password != nil {
		hashed, err := bcrypt.GenerateFromPassword([]byte(*req.User.Password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, response.Error(http.StatusInternalServerError, "Error al procesar contraseña", err))
			return
		}
		passwordHash = pgtype.Text{String: string(hashed), Valid: true}
	}

	affected, err := h.repo.Update(c.Request.Context(), sqlc.UpdateUserParams{
		BusinessID:  businessID,
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
	if affected == 0 {
		c.JSON(http.StatusNotFound, response.Error(http.StatusNotFound, "Usuario no encontrado"))
		return
	}

	c.JSON(http.StatusOK, response.Success[any]("Usuario actualizado", nil))
}

func (h *UserHandler) Delete(c *gin.Context) {
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

	affected, err := h.repo.Delete(c.Request.Context(), sqlc.DeleteUserParams{
		BusinessID: businessID,
		ID:         id,
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, "Error al eliminar usuario"))
		return
	}
	if affected == 0 {
		c.JSON(http.StatusNotFound, response.Error(http.StatusNotFound, "Usuario no encontrado"))
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

// Checks
func (h *UserHandler) CheckIcAvailability(c *gin.Context) {
	businessID, ok := ctxkeys.BusinessID(c)
	if !ok {
		c.JSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, "Formato de ID de negocio inválido"))
		return
	}

	ic := c.Param("ic")
	if ic == "" || len(ic) < 7 || len(ic) > 8 {
		c.JSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, "Formato de IC inválido"))
		return
	}

	available, err := h.repo.CheckIcAvailability(c.Request.Context(), sqlc.CheckIcAvailabilityParams{BusinessID: businessID, Ic: ic})
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(http.StatusInternalServerError, "Error al verificar disponibilidad de IC", err))
		return
	}

	available = !available

	c.JSON(http.StatusOK, response.Success("Disponibilidad de IC", &available))
}

func (h *UserHandler) CheckEmailAvailability(c *gin.Context) {
	businessID, ok := ctxkeys.BusinessID(c)
	if !ok {
		c.JSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, "Formato de ID de negocio inválido"))
		return
	}

	email := c.Param("email")
	if email == "" {
		c.JSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, "Formato de email inválido"))
		return
	}

	if _, err := mail.ParseAddress(email); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, "Formato de email inválido"))
		return
	}

	available, err := h.repo.CheckEmailAvailability(c.Request.Context(), sqlc.CheckEmailAvailabilityParams{BusinessID: businessID, Email: email})
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(http.StatusInternalServerError, "Error al verificar disponibilidad de email", err))
		return
	}

	available = !available

	c.JSON(http.StatusOK, response.Success("Disponibilidad de email", &available))
}

func (h *UserHandler) CheckUsernameAvailability(c *gin.Context) {
	businessID, ok := ctxkeys.BusinessID(c)
	if !ok {
		c.JSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, "Formato de ID de negocio inválido"))
		return
	}

	userName := c.Param("userName")
	if !strings.HasPrefix(userName, "@") || len(userName) < 4 {
		c.JSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, "Formato de nombre de usuario inválido"))
		return
	}

	available, err := h.repo.CheckUsernameAvailability(c.Request.Context(), sqlc.CheckUsernameAvailabilityParams{BusinessID: businessID, UserName: userName})
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(http.StatusInternalServerError, "Error al verificar disponibilidad de nombre de usuario", err))
		return
	}

	available = !available

	c.JSON(http.StatusOK, response.Success("Disponibilidad de nombre de usuario", &available))
}

package business

import (
	"log"
	"net/http"

	"github.com/alanloffler/go-calth-api/internal/common/response"
	"github.com/alanloffler/go-calth-api/internal/common/utils"
	"github.com/alanloffler/go-calth-api/internal/database/sqlc"
	"github.com/alanloffler/go-calth-api/internal/queue"
	"github.com/alanloffler/go-calth-api/internal/user"
	"github.com/gin-gonic/gin"
	"github.com/hibiken/asynq"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

type BusinessHandler struct {
	repo        *BusinessRepository
	userRepo    *user.UserRepository
	pool        *pgxpool.Pool
	queueClient *asynq.Client
	appDomain   string
}

func NewBusinessHandler(repo *BusinessRepository, userRepo *user.UserRepository, pool *pgxpool.Pool, queueClient *asynq.Client, appDomain string) *BusinessHandler {
	return &BusinessHandler{repo: repo, userRepo: userRepo, pool: pool, queueClient: queueClient, appDomain: appDomain}
}

type createBusinessData struct {
	Slug        string `json:"slug" binding:"required,min=3,max=50"`
	TaxId       string `json:"taxId" binding:"required,len=11,numeric"`
	CompanyName string `json:"companyName" binding:"required,min=3,max=100"`
	TradeName   string `json:"tradeName" binding:"required,min=3,max=100"`
	Description string `json:"description" binding:"required,min=3,max=100"`
	Street      string `json:"street" binding:"required,min=3,max=50"`
	City        string `json:"city" binding:"required,min=3,max=50"`
	Province    string `json:"province" binding:"required,min=3,max=50"`
	Country     string `json:"country" binding:"required,min=2,max=50"`
	ZipCode     string `json:"zipCode" binding:"required,min=4,max=10"`
	Timezone    string `json:"timezone" binding:"required,min=2,max=100"`
}

type createContactData struct {
	Email          string  `json:"email" binding:"required,email"`
	PhoneNumber    string  `json:"phoneNumber" binding:"required,len=10,numeric"`
	WhatsAppNumber *string `json:"whatsAppNumber" binding:"omitempty,len=10,numeric"`
	Website        *string `json:"website" binding:"omitempty,min=6"`
}

type createAdminData struct {
	Ic          string `json:"ic" binding:"required,len=8"`
	UserName    string `json:"userName" binding:"required,min=3,max=100"`
	FirstName   string `json:"firstName" binding:"required,min=3,max=100"`
	LastName    string `json:"lastName" binding:"required,min=3,max=100"`
	Email       string `json:"email" binding:"required,email,max=100"`
	Password    string `json:"password" binding:"required,min=8,max=100"`
	PhoneNumber string `json:"phoneNumber" binding:"required,len=10,numeric"`
	RoleId      string `json:"roleId" binding:"required,uuid"`
}

type CreateBusinessWithAdminRequest struct {
	Business createBusinessData `json:"business" binding:"required"`
	Contact  createContactData  `json:"contact" binding:"required"`
	Admin    createAdminData    `json:"admin" binding:"required"`
}

type UpdateBusinessRequest struct {
	Slug           *string `json:"slug" binding:"omitempty,min=3,max=50"`
	TaxId          *string `json:"taxId" binding:"omitempty,len=11,numeric"`
	CompanyName    *string `json:"companyName" binding:"omitempty,min=3,max=100"`
	TradeName      *string `json:"tradeName" binding:"omitempty,min=3,max=100"`
	Description    *string `json:"description" binding:"omitempty,min=3,max=100"`
	Street         *string `json:"street" binding:"omitempty,min=3,max=50"`
	City           *string `json:"city" binding:"omitempty,min=3,max=50"`
	Province       *string `json:"province" binding:"omitempty,min=3,max=50"`
	Country        *string `json:"country" binding:"omitempty,min=2,max=50"`
	ZipCode        *string `json:"zipCode" binding:"omitempty,min=4,max=10"`
	Timezone       *string `json:"timezone" binding:"omitempty,min=2,max=100"`
	Email          *string `json:"email" binding:"omitempty,email"`
	PhoneNumber    *string `json:"phoneNumber" binding:"omitempty,len=10,numeric"`
	WhatsappNumber *string `json:"whatsappNumber" binding:"omitempty,len=10,numeric"`
	Website        *string `json:"website" binding:"omitempty,min=6"`
}

type BusinessWithUsersResponse struct {
	sqlc.Business
	Users []businessUserResponse `json:"users"`
}

type businessUserResponse struct {
	ID          pgtype.UUID        `json:"id"`
	Ic          string             `json:"ic"`
	UserName    string             `json:"userName"`
	FirstName   string             `json:"firstName"`
	LastName    string             `json:"lastName"`
	Email       string             `json:"email"`
	PhoneNumber string             `json:"phoneNumber"`
	Role        *businessUserRole  `json:"role"`
	BusinessID  pgtype.UUID        `json:"businessId"`
	CreatedAt   pgtype.Timestamptz `json:"createdAt"`
	UpdatedAt   pgtype.Timestamptz `json:"updatedAt"`
	DeletedAt   pgtype.Timestamptz `json:"deletedAt"`
}

type businessUserRole struct {
	ID          pgtype.UUID `json:"id"`
	Name        string      `json:"name"`
	Value       string      `json:"value"`
	Description string      `json:"description"`
}

func (h *BusinessHandler) Create(c *gin.Context) {
	var req CreateBusinessWithAdminRequest
	if err := c.ShouldBindBodyWithJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, "Error al crear negocio", err))
		return
	}

	ctx := c.Request.Context()

	taxIdExists, err := h.repo.CheckTaxIDAvailability(ctx, req.Business.TaxId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(http.StatusInternalServerError, "Error al verificar CUIT", err))
		return
	}
	if taxIdExists {
		c.JSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, "CUIT no disponible, debes elegir un CUIT diferente"))
		return
	}

	slugExists, err := h.repo.CheckSlugAvailability(ctx, req.Business.Slug)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(http.StatusInternalServerError, "Error al verificar subdominio", err))
		return
	}
	if slugExists {
		c.JSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, "Subdominio no disponible, debes elegir un subdominio diferente"))
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Admin.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(http.StatusInternalServerError, "Error al procesar contraseña", err))
		return
	}

	tx, err := h.pool.Begin(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(http.StatusInternalServerError, "Error al iniciar transacción", err))
		return
	}
	defer tx.Rollback(ctx)

	qtx := sqlc.New(tx)

	var whatsappNumber pgtype.Text
	if req.Contact.WhatsAppNumber != nil {
		whatsappNumber = pgtype.Text{String: *req.Contact.WhatsAppNumber, Valid: true}
	}

	var website pgtype.Text
	if req.Contact.Website != nil {
		website = pgtype.Text{String: *req.Contact.Website, Valid: true}
	}

	business, err := qtx.CreateBusiness(ctx, sqlc.CreateBusinessParams{
		Slug:           req.Business.Slug,
		TaxID:          req.Business.TaxId,
		CompanyName:    req.Business.CompanyName,
		TradeName:      req.Business.TradeName,
		Description:    req.Business.Description,
		Street:         req.Business.Street,
		City:           req.Business.City,
		Province:       req.Business.Province,
		Country:        req.Business.Country,
		ZipCode:        req.Business.ZipCode,
		Timezone:       req.Business.Timezone,
		Email:          req.Contact.Email,
		PhoneNumber:    req.Contact.PhoneNumber,
		WhatsappNumber: whatsappNumber,
		Website:        website,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(http.StatusInternalServerError, "Error al crear negocio", err))
		return
	}

	var roleID pgtype.UUID
	if err := roleID.Scan(req.Admin.RoleId); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, "Formato de roleId inválido", err))
		return
	}

	_, err = qtx.CreateUser(ctx, sqlc.CreateUserParams{
		Ic:          req.Admin.Ic,
		UserName:    req.Admin.UserName,
		FirstName:   req.Admin.FirstName,
		LastName:    req.Admin.LastName,
		Email:       req.Admin.Email,
		Password:    string(hashedPassword),
		PhoneNumber: req.Admin.PhoneNumber,
		RoleID:      roleID,
		BusinessID:  business.ID,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(http.StatusInternalServerError, "Error al crear administrador", err))
		return
	}

	if err := tx.Commit(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(http.StatusInternalServerError, "Error al confirmar transacción", err))
		return
	}

	// queue send email
	if err := queue.EnqueueBusinessCreated(h.queueClient, queue.BusinessCreatedPayload{
		Email:        req.Contact.Email,
		BusinessName: req.Business.TradeName,
		BusinessLink: "https://" + req.Business.Slug + "." + h.appDomain,
	}); err != nil {
		log.Printf("failed to enqueue business_created email: %v", err)
	}

	c.JSON(http.StatusCreated, response.Created("Negocio creado", &business))
}

func (h *BusinessHandler) GetAll(c *gin.Context) {
	businesses, err := h.repo.GetAll(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusNotFound, response.Error(http.StatusNotFound, "Negocios no encontrados", err))
		return
	}

	c.JSON(http.StatusOK, response.Success("Negocios encontrados", &businesses))
}

func (h *BusinessHandler) GetOneByID(c *gin.Context) {
	var id pgtype.UUID
	if err := id.Scan(c.Param("id")); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, "Formato de Id inválido", err))
		return
	}

	business, err := h.repo.GetOneByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, response.Error(http.StatusNotFound, "Negocio no encontrado", err))
		return
	}

	users, err := h.userRepo.GetByBusinessID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(http.StatusInternalServerError, "Error al obtener usuarios", err))
		return
	}

	mapped := make([]businessUserResponse, len(users))
	for i, u := range users {
		mapped[i] = businessUserResponse{
			ID:          u.ID,
			Ic:          u.Ic,
			UserName:    u.UserName,
			FirstName:   u.FirstName,
			LastName:    u.LastName,
			Email:       u.Email,
			PhoneNumber: u.PhoneNumber,
			BusinessID:  business.ID,
			CreatedAt:   u.CreatedAt,
			UpdatedAt:   u.UpdatedAt,
			DeletedAt:   u.DeletedAt,
		}
		if u.RoleID.Valid {
			mapped[i].Role = &businessUserRole{
				ID:          u.RoleID,
				Name:        u.RoleName,
				Value:       u.RoleValue,
				Description: u.RoleDescription,
			}
		}
	}

	c.JSON(http.StatusOK, response.Success("Negocio encontrado", &BusinessWithUsersResponse{
		Business: business,
		Users:    mapped,
	}))
}

func (h *BusinessHandler) Update(c *gin.Context) {
	var id pgtype.UUID
	if err := id.Scan(c.Param("id")); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, "Formato de ID inválido", err))
		return
	}

	var req UpdateBusinessRequest
	if err := c.ShouldBindBodyWithJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, "Error al actualizar negocio", err))
		return
	}

	affected, err := h.repo.Update(c.Request.Context(), sqlc.UpdateBusinessParams{
		ID:             id,
		Slug:           utils.ToPgText(req.Slug),
		TaxID:          utils.ToPgText(req.TaxId),
		CompanyName:    utils.ToPgText(req.CompanyName),
		TradeName:      utils.ToPgText(req.TradeName),
		Description:    utils.ToPgText(req.Description),
		Street:         utils.ToPgText(req.Street),
		City:           utils.ToPgText(req.City),
		Province:       utils.ToPgText(req.Province),
		Country:        utils.ToPgText(req.Country),
		ZipCode:        utils.ToPgText(req.ZipCode),
		Timezone:       utils.ToPgText(req.Timezone),
		Email:          utils.ToPgText(req.Email),
		PhoneNumber:    utils.ToPgText(req.PhoneNumber),
		WhatsappNumber: utils.ToPgText(req.WhatsappNumber),
		Website:        utils.ToPgText(req.Website),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(http.StatusInternalServerError, "Error al actualizar negocio", err))
		return
	}
	if affected == 0 {
		c.JSON(http.StatusNotFound, response.Error(http.StatusNotFound, "Negocio no encontrado"))
		return
	}

	c.JSON(http.StatusOK, response.Success[any]("Negocio actualizado", nil))
}

func (h *BusinessHandler) Delete(c *gin.Context) {
	var id pgtype.UUID
	if err := id.Scan(c.Param("id")); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, "Formato de ID inválido", err))
		return
	}

	// TODO:
	// Plan:
	// log.Println("Init transaction")
	// log.Println("Remove users, on cascade other related content")
	// log.Println("Remove business")
	//
	// err := h.repo.Delete(c.Request.Context(), id)
	// if err != nil {
	// 	c.JSON(http.StatusInternalServerError, response.Error(http.StatusInternalServerError, "Error al eliminar negocio", err))
	// 	return
	// }
	// if rows == 0 {
	// 	c.JSON(http.StatusNotFound, response.Error(http.StatusNotFound, "Negocio no encontrado", nil))
	// 	return
	// }

	c.JSON(http.StatusOK, response.Success[any]("Negocio eliminado", nil))
}

func (h *BusinessHandler) CheckTaxIDAvailability(c *gin.Context) {
	ctx := c.Request.Context()

	taxId := c.Param("taxId")
	if taxId == "" || len(taxId) != 11 {
		c.JSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, "Formato de CUIT inválido"))
		return
	}

	available, err := h.repo.CheckTaxIDAvailability(ctx, taxId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(http.StatusInternalServerError, "Error al verificar el CUIT", err))
		return
	}

	available = !available

	var message string
	if available {
		message = "CUIT disponible"
	} else {
		message = "CUIT no disponible"
	}

	c.JSON(http.StatusOK, response.Success(message, &available))
}

func (h *BusinessHandler) CheckSlugAvailability(c *gin.Context) {
	ctx := c.Request.Context()

	slug := c.Param("slug")
	if slug == "" || len(slug) < 3 {
		c.JSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, "Formato de subdominio inválido"))
		return
	}

	available, err := h.repo.CheckSlugAvailability(ctx, slug)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(http.StatusInternalServerError, "Error al verificar subdominio", err))
		return
	}

	available = !available

	var message string
	if available {
		message = "Subdominio disponible"
	} else {
		message = "Subdominio no disponible"
	}

	c.JSON(http.StatusOK, response.Success(message, &available))
}

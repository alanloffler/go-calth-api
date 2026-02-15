package business

import (
	"net/http"

	"github.com/alanloffler/go-calth-api/internal/database/sqlc"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
)

type BusinessHandler struct {
	repo *BusinessRepository
}

func NewBusinessHandler(repo *BusinessRepository) *BusinessHandler {
	return &BusinessHandler{repo: repo}
}

type CreateBusinessRequest struct {
	Slug           string  `json:"slug" binding:"required,min=3,max=50"`
	TaxId          string  `json:"taxId" binding:"required,len=11,numeric"`
	CompanyName    string  `json:"companyName" binding:"required,min=3,max=100"`
	TradeName      string  `json:"tradeName" binding:"required,min=3,max=100"`
	Description    string  `json:"description" binding:"required,min=3,max=100"`
	Street         string  `json:"street" binding:"required,min=3,max=50"`
	City           string  `json:"city" binding:"required,min=3,max=50"`
	Province       string  `json:"province" binding:"required,min=3,max=50"`
	Country        string  `json:"country" binding:"required,min=3,max=50"`
	ZipCode        string  `json:"zipCode" binding:"required,min=4,max=10"`
	Email          string  `json:"email" binding:"required,email"`
	PhoneNumber    string  `json:"phoneNumber" binding:"required,len=10,numeric"`
	WhatsappNumber *string `json:"whatsappNumber" binding:"omitempty,len=10,numeric"`
	Website        *string `json:"website" binding:"omitempty,min=6"`
}

func (h *BusinessHandler) Create(c *gin.Context) {
	var req CreateBusinessRequest
	if err := c.ShouldBindBodyWithJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	business, err := h.repo.Create(c.Request.Context(), sqlc.CreateBusinessParams{
		Slug:           req.Slug,
		TaxID:          req.TaxId,
		CompanyName:    req.CompanyName,
		TradeName:      req.TradeName,
		Description:    req.Description,
		Street:         req.Street,
		City:           req.City,
		Province:       req.Province,
		Country:        req.Country,
		ZipCode:        req.ZipCode,
		Email:          req.Email,
		PhoneNumber:    req.PhoneNumber,
		WhatsappNumber: pgtype.Text{String: "", Valid: req.WhatsappNumber != nil},
		Website:        pgtype.Text{String: "", Valid: req.Website != nil},
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, business)
}

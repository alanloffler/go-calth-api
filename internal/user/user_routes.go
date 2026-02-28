package user

import (
	"github.com/alanloffler/go-calth-api/internal/database/sqlc"
	"github.com/alanloffler/go-calth-api/internal/patient_profile"
	"github.com/alanloffler/go-calth-api/internal/professional_profile"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func RegisterRoutes(router *gin.RouterGroup, q *sqlc.Queries, pool *pgxpool.Pool) {
	var repo *UserRepository = NewUserRepository(q)
	var ppRepo *patient_profile.PatientProfileRepository = patient_profile.NewPatientProfileRepository(q)
	var prpRepo *professional_profile.ProfessionalProfileRepository = professional_profile.NewProfessionalProfileRepository(q)
	var handler *UserHandler = NewUserHandler(repo, pool, ppRepo, prpRepo)
	var users *gin.RouterGroup = router.Group("/users")

	users.POST("/admin", handler.CreateAdmin)
	users.POST("/patient", handler.CreatePatient)
	users.POST("/professional", handler.CreateProfessional)

	users.GET("", handler.GetAll)
	users.GET("/soft", handler.GetAllWithSoftDeleted)
	users.GET("/role/:role", handler.GetAllByRole)
	users.GET("/role/:role/soft", handler.GetAllByRoleWithSoftDeleted)
	users.GET("/:id", handler.GetByID)
	users.GET("/:id/soft", handler.GetByIDWithSoftDeleted)
	users.GET("/:id/admin/profile/soft", handler.GetByIDWithSoftDeleted)
	users.GET("/:id/patient/profile/soft", handler.GetPatientByIDWithSoftDeleted)
	users.GET("/:id/professional/profile/soft", handler.GetProfessionalByIDWithSoftDeleted)
	users.GET("/:id/admin/profile", handler.GetByID)
	users.GET("/:id/patient/profile", handler.GetPatientByID)
	users.GET("/:id/professional/profile", handler.GetProfessionalByID)

	users.PATCH("/:id", handler.Update)
	users.PATCH("/:id/patient", handler.UpdatePatient)
	users.PATCH("/:id/professional", handler.UpdateProfessional)
	users.PATCH("/:id/restore", handler.Restore)
	users.PATCH("/:id/admin/restore", handler.Restore)
	users.PATCH("/:id/patient/restore", handler.Restore)
	users.PATCH("/:id/professional/restore", handler.Restore)

	users.DELETE("/:id", handler.Delete)
	users.DELETE("/:id/soft", handler.SoftDelete)
	users.DELETE("/:id/admin/soft", handler.SoftDelete)
	users.DELETE("/:id/patient/soft", handler.SoftDelete)
	users.DELETE("/:id/professional/soft", handler.SoftDelete)
	users.DELETE("/:id/admin", handler.Delete)
	users.DELETE("/:id/patient", handler.Delete)
	users.DELETE("/:id/professional", handler.Delete)
	// Checks
	users.GET("/check/email/:email", handler.CheckEmailAvailability)
	users.GET("/check/ic/:ic", handler.CheckIcAvailability)
	users.GET("/check/username/:userName", handler.CheckUsernameAvailability)
}

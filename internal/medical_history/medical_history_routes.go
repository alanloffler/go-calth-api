package medical_history

import (
	"github.com/alanloffler/go-calth-api/internal/database/sqlc"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.RouterGroup, q *sqlc.Queries) {
	var repo *MedicalHistoryRepository = NewMedicalHistoryRepository(q)
	var handler *MedicalHistoryHandler = NewMedicalHistoryHandler(repo)
	var medical_histories *gin.RouterGroup = router.Group("/medical-history")
}

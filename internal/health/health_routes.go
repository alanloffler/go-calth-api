package health

import (
	"net/http"

	"github.com/alanloffler/go-calth-api/internal/common/response"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func RegisterRoutes(router *gin.Engine, pool *pgxpool.Pool) {
	var health *gin.RouterGroup = router.Group("/health")

	health.GET("", func(c *gin.Context) {
		c.JSON(http.StatusOK, response.Success[any]("Calth API running", nil))
	})

	health.GET("/db", func(c *gin.Context) {
		if err := pool.Ping(c.Request.Context()); err != nil {
			c.JSON(http.StatusInternalServerError, response.Error(http.StatusInternalServerError, "Database is not connected", err))
			return
		}

		c.JSON(http.StatusOK, response.Success[any]("Database connected successfully", nil))
	})
}

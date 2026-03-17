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
		c.JSON(http.StatusOK, response.Success[any]("Calth API en ejecución", nil))
	})

	health.GET("/db", func(c *gin.Context) {
		if err := pool.Ping(c.Request.Context()); err != nil {
			c.JSON(http.StatusInternalServerError, response.Error(http.StatusInternalServerError, "La base de datos no está conectada", err))
			return
		}

		c.JSON(http.StatusOK, response.Success[any]("Base de datos conectada exitosamente", nil))
	})
}

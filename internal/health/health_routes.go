package health

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func RegisterRoutes(router *gin.Engine, pool *pgxpool.Pool) {
	var health *gin.RouterGroup = router.Group("/health")

	health.GET("", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Calth API running", "status": "success"})
	})

	health.GET("/db", func(c *gin.Context) {
		if err := pool.Ping(c.Request.Context()); err != nil {
			c.JSON(500, gin.H{"message": err.Error(), "status": "error"})
			return
		}
		c.JSON(200, gin.H{
			"message": "Database connected successfully", "status": "connected",
		})
	})
}

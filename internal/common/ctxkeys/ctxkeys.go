package ctxkeys

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
)

func BusinessID(c *gin.Context) (pgtype.UUID, bool) {
	return scanUUID(c, "businessID")
}

func scanUUID(c *gin.Context, key string) (pgtype.UUID, bool) {
	val, exists := c.Get(key)
	if !exists {
		return pgtype.UUID{}, false
	}

	var id pgtype.UUID
	if err := id.Scan(val.(string)); err != nil {
		return pgtype.UUID{}, false
	}

	return id, true
}

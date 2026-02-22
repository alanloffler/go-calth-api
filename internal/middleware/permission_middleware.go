package middleware

import (
	"net/http"

	"github.com/alanloffler/go-calth-api/internal/common/ctxkeys"
	"github.com/alanloffler/go-calth-api/internal/common/response"
	"github.com/alanloffler/go-calth-api/internal/database/sqlc"
	"github.com/gin-gonic/gin"
)

type PermissionMode string

const (
	PermissionSome  PermissionMode = "some"
	PermissionEvery PermissionMode = "every"
)

func PermissionMiddleware(q *sqlc.Queries, permissions any, mode ...PermissionMode) gin.HandlerFunc {
	var keys []string
	switch v := permissions.(type) {
	case string:
		keys = []string{v}
	case []string:
		keys = v
	}

	m := PermissionEvery
	if len(mode) > 0 {
		m = mode[0]
	}

	return func(c *gin.Context) {
		roleID, ok := ctxkeys.RoleID(c)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, response.Error(http.StatusUnauthorized, "Usuario no autenticado"))
			return
		}

		var matched int
		for _, key := range keys {
			has, err := q.HasPermission(c.Request.Context(), sqlc.HasPermissionParams{
				RoleID:    roleID,
				ActionKey: key,
			})
			if err != nil {
				c.AbortWithStatusJSON(http.StatusInternalServerError, response.Error(http.StatusInternalServerError, "Error al verificar permisos", err))
				return
			}
			if has {
				matched++
			}
		}

		allowed := false
		switch m {
		case PermissionSome:
			allowed = matched > 0
		case PermissionEvery:
			allowed = matched == len(keys)
		}

		if !allowed {
			c.AbortWithStatusJSON(http.StatusForbidden, response.Error(http.StatusForbidden, "Permisos insuficientes"))
			return
		}

		c.Next()
	}
}

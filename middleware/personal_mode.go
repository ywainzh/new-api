package middleware

import (
	"net/http"

	"github.com/QuantumNous/new-api/setting/operation_setting"
	"github.com/gin-gonic/gin"
)

func DisableInPersonalMode() gin.HandlerFunc {
	return func(c *gin.Context) {
		if operation_setting.IsPersonalModeEnabled() {
			c.JSON(http.StatusForbidden, gin.H{
				"success": false,
				"message": "Personal mode is enabled; this feature is disabled.",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}

package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/QuantumNous/new-api/setting/operation_setting"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func withPersonalMode(t *testing.T, enabled bool) {
	t.Helper()

	original := operation_setting.PersonalModeForced
	operation_setting.PersonalModeForced = enabled
	t.Cleanup(func() {
		operation_setting.PersonalModeForced = original
	})
}

func performPersonalModeRequest(t *testing.T) *httptest.ResponseRecorder {
	t.Helper()

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/api/test", DisableInPersonalMode(), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"success": true})
	})

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/api/test", nil)
	router.ServeHTTP(recorder, request)
	return recorder
}

func TestDisableInPersonalModeAllowsNormalMode(t *testing.T) {
	withPersonalMode(t, false)

	recorder := performPersonalModeRequest(t)

	require.Equal(t, http.StatusOK, recorder.Code)
	assert.JSONEq(t, `{"success":true}`, recorder.Body.String())
}

func TestDisableInPersonalModeBlocksPersonalMode(t *testing.T) {
	withPersonalMode(t, true)

	recorder := performPersonalModeRequest(t)

	require.Equal(t, http.StatusForbidden, recorder.Code)
	assert.JSONEq(t, `{"success":false,"message":"Personal mode is enabled; this feature is disabled."}`, recorder.Body.String())
}

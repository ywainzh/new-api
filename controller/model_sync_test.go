package controller

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/QuantumNous/new-api/common"
	"github.com/QuantumNous/new-api/constant"
	"github.com/QuantumNous/new-api/model"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func insertSyncTestChannel(t *testing.T, models string) {
	t.Helper()

	channel := model.Channel{
		Type:   constant.ChannelTypeOpenAI,
		Key:    "sk-test",
		Status: common.ChannelStatusEnabled,
		Name:   "sync test channel",
		Models: models,
		Group:  "default",
	}
	require.NoError(t, model.BatchInsertChannels([]model.Channel{channel}))
}

func resetSyncMetadataCache(t *testing.T) {
	t.Helper()

	cacheMutex.Lock()
	etagCache = make(map[string]string)
	bodyCache = make(map[string][]byte)
	cacheMutex.Unlock()
}

func TestSyncChannelModelsFromAbilitiesCreatesOfficialAndBasicModels(t *testing.T) {
	setupModelListControllerTestDB(t)
	resetSyncMetadataCache(t)
	insertSyncTestChannel(t, "official-model,custom-official-missing")

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/i18n/zh-CN/newapi/models.json":
			_, _ = w.Write([]byte(`{
				"success": true,
				"data": [{
					"model_name": "official-model",
					"description": "Official model description",
					"icon": "OpenAI",
					"tags": "chat,vision",
					"vendor_name": "Official Vendor",
					"status": 1,
					"name_rule": 2
				}]
			}`))
		case "/api/i18n/zh-CN/newapi/vendors.json":
			_, _ = w.Write([]byte(`{
				"success": true,
				"data": [{
					"name": "Official Vendor",
					"description": "Official vendor description",
					"icon": "VendorIcon",
					"status": 1
				}]
			}`))
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()
	t.Setenv("SYNC_UPSTREAM_BASE", server.URL)
	t.Setenv("SYNC_HTTP_RETRY", "1")

	result, err := SyncChannelModelsFromAbilities(context.Background(), "zh")
	require.NoError(t, err)
	assert.Equal(t, 2, result.CreatedModels)
	assert.Equal(t, 1, result.CreatedBasicModels)
	assert.Equal(t, 1, result.CreatedVendors)
	assert.ElementsMatch(t, []string{"official-model", "custom-official-missing"}, result.CreatedList)
	assert.Empty(t, result.SkippedModels)

	var official model.Model
	require.NoError(t, model.DB.Where("model_name = ?", "official-model").First(&official).Error)
	assert.Equal(t, "Official model description", official.Description)
	assert.Equal(t, "OpenAI", official.Icon)
	assert.Equal(t, "chat,vision", official.Tags)
	assert.Equal(t, model.NameRuleContains, official.NameRule)
	assert.Equal(t, 1, official.Status)
	assert.Equal(t, 1, official.SyncOfficial)
	assert.NotZero(t, official.VendorID)

	var vendor model.Vendor
	require.NoError(t, model.DB.First(&vendor, official.VendorID).Error)
	assert.Equal(t, "Official Vendor", vendor.Name)
	assert.Equal(t, "Official vendor description", vendor.Description)

	var basic model.Model
	require.NoError(t, model.DB.Where("model_name = ?", "custom-official-missing").First(&basic).Error)
	assert.Empty(t, basic.Description)
	assert.Empty(t, basic.Icon)
	assert.Empty(t, basic.Tags)
	assert.Zero(t, basic.VendorID)
	assert.Equal(t, 1, basic.Status)
	assert.Equal(t, 1, basic.SyncOfficial)

	repeat, err := SyncChannelModelsFromAbilities(context.Background(), "zh")
	require.NoError(t, err)
	assert.Zero(t, repeat.CreatedModels)
	assert.Zero(t, repeat.CreatedBasicModels)

	var count int64
	require.NoError(t, model.DB.Model(&model.Model{}).Where("model_name IN ?", []string{"official-model", "custom-official-missing"}).Count(&count).Error)
	assert.EqualValues(t, 2, count)
}

func TestSyncChannelModelsFromAbilitiesFallsBackToBasicWhenUpstreamFails(t *testing.T) {
	setupModelListControllerTestDB(t)
	resetSyncMetadataCache(t)
	insertSyncTestChannel(t, "offline-only-model")

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "offline", http.StatusInternalServerError)
	}))
	defer server.Close()
	t.Setenv("SYNC_UPSTREAM_BASE", server.URL)
	t.Setenv("SYNC_HTTP_RETRY", "1")

	result, err := SyncChannelModelsFromAbilities(context.Background(), "zh")
	require.NoError(t, err)
	assert.Equal(t, 1, result.CreatedModels)
	assert.Equal(t, 1, result.CreatedBasicModels)
	assert.NotEmpty(t, result.UpstreamError)

	var basic model.Model
	require.NoError(t, model.DB.Where("model_name = ?", "offline-only-model").First(&basic).Error)
	assert.Equal(t, "offline-only-model", basic.ModelName)
	assert.Equal(t, 1, basic.Status)
	assert.Equal(t, 1, basic.SyncOfficial)
}

func TestAddChannelAutoSyncsChannelModels(t *testing.T) {
	db := setupModelListControllerTestDB(t)
	resetSyncMetadataCache(t)
	require.NoError(t, db.AutoMigrate(&model.Log{}))

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "offline", http.StatusInternalServerError)
	}))
	defer server.Close()
	t.Setenv("SYNC_UPSTREAM_BASE", server.URL)
	t.Setenv("SYNC_HTTP_RETRY", "1")

	body, err := common.Marshal(gin.H{
		"mode": "single",
		"channel": gin.H{
			"type":   constant.ChannelTypeOpenAI,
			"key":    "sk-test",
			"status": common.ChannelStatusEnabled,
			"name":   "auto sync channel",
			"models": "auto-sync-model",
			"group":  "default",
		},
	})
	require.NoError(t, err)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/api/channel/", strings.NewReader(string(body)))
	c.Request.Header.Set("Content-Type", "application/json")

	AddChannel(c)

	require.Equal(t, http.StatusOK, w.Code)
	var response struct {
		Success   bool `json:"success"`
		ModelSync struct {
			Success bool `json:"success"`
			Data    struct {
				CreatedModels      int    `json:"created_models"`
				CreatedBasicModels int    `json:"created_basic_models"`
				UpstreamError      string `json:"upstream_error"`
			} `json:"data"`
		} `json:"model_sync"`
	}
	require.NoError(t, common.Unmarshal(w.Body.Bytes(), &response))
	assert.True(t, response.Success)
	assert.True(t, response.ModelSync.Success)
	assert.Equal(t, 1, response.ModelSync.Data.CreatedModels)
	assert.Equal(t, 1, response.ModelSync.Data.CreatedBasicModels)
	assert.NotEmpty(t, response.ModelSync.Data.UpstreamError)

	var synced model.Model
	require.NoError(t, model.DB.Where("model_name = ?", "auto-sync-model").First(&synced).Error)
	assert.Equal(t, "auto-sync-model", synced.ModelName)
	assert.Equal(t, 1, synced.Status)
	assert.Equal(t, 1, synced.SyncOfficial)
}

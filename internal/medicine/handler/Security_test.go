package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func setupMedicineTestEngine(role string) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(func(c *gin.Context) {
		c.Set("role", role)
		c.Set("user_id", uint(1))
		c.Next()
	})

	apiV1 := r.Group("/api/v1")
	StartMedicineHandler(&HandlerContract{
		Logger: logrus.New(),
		Router: apiV1,
	}, &MedicineHandlerProps{
		MedicineService:    nil,
		TransactionService: nil,
	})

	return r
}

func TestMedicineHandlerSecurity(t *testing.T) {
	t.Run("POST /api/v1/medicines by kasir -> 403", func(t *testing.T) {
		r := setupMedicineTestEngine("kasir")
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/medicines", nil)
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusForbidden, w.Code)
	})

	t.Run("POST /api/v1/medicines by owner -> 403", func(t *testing.T) {
		r := setupMedicineTestEngine("owner")
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/medicines", nil)
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusForbidden, w.Code)
	})

	t.Run("POST /api/v1/medicines by admin -> bukan 403", func(t *testing.T) {
		r := setupMedicineTestEngine("admin")
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/medicines", nil)
		r.ServeHTTP(w, req)
		assert.NotEqual(t, http.StatusForbidden, w.Code)
	})

	t.Run("PUT /api/v1/medicines/1 by kasir -> 403", func(t *testing.T) {
		r := setupMedicineTestEngine("kasir")
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("PUT", "/api/v1/medicines/1", nil)
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusForbidden, w.Code)
	})

	t.Run("PUT /api/v1/medicines/1 by owner -> 403", func(t *testing.T) {
		r := setupMedicineTestEngine("owner")
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("PUT", "/api/v1/medicines/1", nil)
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusForbidden, w.Code)
	})

	t.Run("PUT /api/v1/medicines/1 by admin -> bukan 403", func(t *testing.T) {
		r := setupMedicineTestEngine("admin")
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("PUT", "/api/v1/medicines/1", nil)
		r.ServeHTTP(w, req)
		assert.NotEqual(t, http.StatusForbidden, w.Code)
	})

	t.Run("PATCH /api/v1/medicines/1/status by kasir -> 403", func(t *testing.T) {
		r := setupMedicineTestEngine("kasir")
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("PATCH", "/api/v1/medicines/1/status", nil)
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusForbidden, w.Code)
	})

	t.Run("PATCH /api/v1/medicines/1/status by owner -> 403", func(t *testing.T) {
		r := setupMedicineTestEngine("owner")
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("PATCH", "/api/v1/medicines/1/status", nil)
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusForbidden, w.Code)
	})

	t.Run("PATCH /api/v1/medicines/1/status by admin -> bukan 403", func(t *testing.T) {
		r := setupMedicineTestEngine("admin")
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("PATCH", "/api/v1/medicines/1/status", nil)
		r.ServeHTTP(w, req)
		assert.NotEqual(t, http.StatusForbidden, w.Code)
	})

	t.Run("DELETE /api/v1/medicines/1 by kasir -> 403", func(t *testing.T) {
		r := setupMedicineTestEngine("kasir")
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", "/api/v1/medicines/1", nil)
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusForbidden, w.Code)
	})

	t.Run("DELETE /api/v1/medicines/1 by owner -> 403", func(t *testing.T) {
		r := setupMedicineTestEngine("owner")
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", "/api/v1/medicines/1", nil)
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusForbidden, w.Code)
	})

	t.Run("DELETE /api/v1/medicines/1 by admin -> bukan 403", func(t *testing.T) {
		r := setupMedicineTestEngine("admin")
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", "/api/v1/medicines/1", nil)
		r.ServeHTTP(w, req)
		assert.NotEqual(t, http.StatusForbidden, w.Code)
	})

	t.Run("POST /api/v1/medicines/scan by owner -> 403", func(t *testing.T) {
		r := setupMedicineTestEngine("owner")
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/medicines/scan", nil)
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusForbidden, w.Code)
	})

	t.Run("POST /api/v1/medicines/scan by kasir -> bukan 403", func(t *testing.T) {
		r := setupMedicineTestEngine("kasir")
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/medicines/scan", nil)
		r.ServeHTTP(w, req)
		assert.NotEqual(t, http.StatusForbidden, w.Code)
	})

	t.Run("POST /api/v1/medicines/scan by admin -> bukan 403", func(t *testing.T) {
		r := setupMedicineTestEngine("admin")
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/medicines/scan", nil)
		r.ServeHTTP(w, req)
		assert.NotEqual(t, http.StatusForbidden, w.Code)
	})

	t.Run("GET /api/v1/medicines by admin -> bukan 403", func(t *testing.T) {
		r := setupMedicineTestEngine("admin")
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/v1/medicines", nil)
		r.ServeHTTP(w, req)
		assert.NotEqual(t, http.StatusForbidden, w.Code)
	})

	t.Run("GET /api/v1/medicines by kasir -> bukan 403", func(t *testing.T) {
		r := setupMedicineTestEngine("kasir")
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/v1/medicines", nil)
		r.ServeHTTP(w, req)
		assert.NotEqual(t, http.StatusForbidden, w.Code)
	})

	t.Run("GET /api/v1/medicines by owner -> bukan 403", func(t *testing.T) {
		r := setupMedicineTestEngine("owner")
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/v1/medicines", nil)
		r.ServeHTTP(w, req)
		assert.NotEqual(t, http.StatusForbidden, w.Code)
	})
}

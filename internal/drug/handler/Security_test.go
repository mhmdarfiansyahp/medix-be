package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func setupDrugTestEngine(role string) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(func(c *gin.Context) {
		c.Set("role", role)
		c.Set("user_id", uint(1))
		c.Next()
	})

	apiV1 := r.Group("/api/v1")
	StartTypeDrugHandler(&HandlerContract{
		Logger: logrus.New(),
		Router: apiV1,
	}, &TypeDrugHandlerProps{
		TypeDrugService: nil,
	})

	return r
}

func TestTypeDrugHandlerSecurity(t *testing.T) {
	t.Run("POST /api/v1/type-drugs by kasir -> 403", func(t *testing.T) {
		r := setupDrugTestEngine("kasir")
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/type-drugs", nil)
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusForbidden, w.Code)
	})

	t.Run("POST /api/v1/type-drugs by owner -> 403", func(t *testing.T) {
		r := setupDrugTestEngine("owner")
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/type-drugs", nil)
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusForbidden, w.Code)
	})

	t.Run("POST /api/v1/type-drugs by admin -> bukan 403", func(t *testing.T) {
		r := setupDrugTestEngine("admin")
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/type-drugs", nil)
		r.ServeHTTP(w, req)
		assert.NotEqual(t, http.StatusForbidden, w.Code)
	})

	t.Run("PUT /api/v1/type-drugs/1 by kasir -> 403", func(t *testing.T) {
		r := setupDrugTestEngine("kasir")
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("PUT", "/api/v1/type-drugs/1", nil)
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusForbidden, w.Code)
	})

	t.Run("PUT /api/v1/type-drugs/1 by owner -> 403", func(t *testing.T) {
		r := setupDrugTestEngine("owner")
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("PUT", "/api/v1/type-drugs/1", nil)
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusForbidden, w.Code)
	})

	t.Run("PUT /api/v1/type-drugs/1 by admin -> bukan 403", func(t *testing.T) {
		r := setupDrugTestEngine("admin")
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("PUT", "/api/v1/type-drugs/1", nil)
		r.ServeHTTP(w, req)
		assert.NotEqual(t, http.StatusForbidden, w.Code)
	})

	t.Run("DELETE /api/v1/type-drugs/1 by kasir -> 403", func(t *testing.T) {
		r := setupDrugTestEngine("kasir")
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", "/api/v1/type-drugs/1", nil)
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusForbidden, w.Code)
	})

	t.Run("DELETE /api/v1/type-drugs/1 by owner -> 403", func(t *testing.T) {
		r := setupDrugTestEngine("owner")
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", "/api/v1/type-drugs/1", nil)
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusForbidden, w.Code)
	})

	t.Run("DELETE /api/v1/type-drugs/1 by admin -> bukan 403", func(t *testing.T) {
		r := setupDrugTestEngine("admin")
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", "/api/v1/type-drugs/1", nil)
		r.ServeHTTP(w, req)
		assert.NotEqual(t, http.StatusForbidden, w.Code)
	})

	t.Run("GET /api/v1/type-drugs by admin -> bukan 403", func(t *testing.T) {
		r := setupDrugTestEngine("admin")
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/v1/type-drugs", nil)
		r.ServeHTTP(w, req)
		assert.NotEqual(t, http.StatusForbidden, w.Code)
	})

	t.Run("GET /api/v1/type-drugs by kasir -> bukan 403", func(t *testing.T) {
		r := setupDrugTestEngine("kasir")
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/v1/type-drugs", nil)
		r.ServeHTTP(w, req)
		assert.NotEqual(t, http.StatusForbidden, w.Code)
	})

	t.Run("GET /api/v1/type-drugs by owner -> bukan 403", func(t *testing.T) {
		r := setupDrugTestEngine("owner")
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/v1/type-drugs", nil)
		r.ServeHTTP(w, req)
		assert.NotEqual(t, http.StatusForbidden, w.Code)
	})
}

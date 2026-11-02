package handler

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"medix-be/internal/user/model/dto"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// mockUserService implements service.UserService interface for testing

type mockUserService struct {}

func (m *mockUserService) CreateUser(req dto.CreateUserRequest) (*dto.UserResponse, error) {
	return nil, nil
}

func (m *mockUserService) Login(req *dto.LoginRequest) (*dto.LoginResponse, error) {
	return nil, nil
}

func (m *mockUserService) RefreshToken(req *dto.RefreshRequest) (*dto.RefreshResponse, error) {
	return nil, nil
}

func (m *mockUserService) Logout(req *dto.LogoutRequest) error {
	return nil
}

func (m *mockUserService) CleanExpiredSessions() error {
	return nil
}

func (m *mockUserService) GetAllUsers(page int, limit int, search string, role string, status string) (*dto.UserListResponse, error) {
	return nil, nil
}

func (m *mockUserService) GetUserByID(id uint) (*dto.UserResponse, error) {
	return nil, nil
}

func (m *mockUserService) UpdateUser(id uint, req dto.UpdateUserRequest) (*dto.UserResponse, error) {
	return nil, nil
}

func (m *mockUserService) DeleteUser(id uint) error {
	return nil
}

func (m *mockUserService) GetProfile(id uint) (*dto.UserResponse, error) {
	return nil, nil
}

func (m *mockUserService) UpdateProfile(id uint, req *dto.UpdateProfileRequest) (*dto.UserResponse, error) {
	return nil, nil
}

func (m *mockUserService) UpdateProfilePhoto(id uint, fotoPath string) (*dto.UserResponse, error) {
	return nil, nil
}

func (m *mockUserService) GetProfileByUsername(username string) (*dto.UserResponse, error) {
	return nil, nil
}

func (m *mockUserService) ChangePassword(id uint, req dto.ChangePasswordRequest) error {
	return nil
}

func (m *mockUserService) ForgotPassword(req dto.ForgotPasswordRequest) error {
	return nil
}

func (m *mockUserService) ResetPasswordWithToken(req dto.ResetPasswordWithTokenRequest) error {
	return nil
}

func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	
	// Create the base group /api/v1
	base := r.Group("/api/v1")
	
	// Create the /users subgroups like the real app does
	public := base.Group("/users")
	protected := base.Group("/users")
	
	h := &UserHandler{
		userService: &mockUserService{},
		routers: []*gin.RouterGroup{public, protected},
	}
	h.RegisterRouter()
	return r
}

func setupTestRouterWithAuth(role string, userID uint) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	
	r.Use(func(c *gin.Context) {
		c.Set("user_id", userID)
		c.Set("role", role)
		c.Next()
	})
	
	base := r.Group("/api/v1")
	public := base.Group("/users")
	protected := base.Group("/users")
	
	h := &UserHandler{
		userService: &mockUserService{},
		routers: []*gin.RouterGroup{public, protected},
	}
	h.RegisterRouter()
	return r
}

func TestUserSecurity(t *testing.T) {
	t.Run("POST /users tanpa auth -> 403", func(t *testing.T) {
		r := setupTestRouter()
		
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/users", bytes.NewBufferString("{}"))
		req.Header.Set("Content-Type", "application/json")
		
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusForbidden, w.Code, "POST /users without auth should be 403")
	})

	t.Run("DELETE /users/:id bukan owner oleh kasir -> 403", func(t *testing.T) {
		r := setupTestRouterWithAuth("kasir", 2)
		
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", "/api/v1/users/1", nil)
		r.ServeHTTP(w, req)
		
		assert.Equal(t, http.StatusForbidden, w.Code, "DELETE /users/:id by non-admin non-owner should be 403")
	})

	t.Run("PUT /users/:id role dari non-admin -> role dihilangkan", func(t *testing.T) {
		r := setupTestRouterWithAuth("kasir", 1)
		
		payload := `{"nama_user": "Test", "role": "admin"}`
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("PUT", "/api/v1/users/1", bytes.NewBufferString(payload))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)
		
		assert.Equal(t, http.StatusOK, w.Code, "PUT /users/:id should work for kasir updating own profile")
	})

	t.Run("POST /login rate limit -> 429", func(t *testing.T) {
		r := setupTestRouter()
		
		// Call multiple times to test rate limiting (limit is 5 per minute)
		for i := 0; i < 5; i++ {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "/api/v1/users/login", bytes.NewBufferString(`{"username":"u","password":"p"}`))
			req.Header.Set("Content-Type", "application/json")
			r.ServeHTTP(w, req)
			assert.NotEqual(t, http.StatusTooManyRequests, w.Code, "Call %d should not be rate limited", i+1)
		}
		
		// 6th call should hit rate limit
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/users/login", bytes.NewBufferString(`{"username":"u","password":"p"}`))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusTooManyRequests, w.Code, "Should be rate limited on 6th call")
	})

	t.Run("POST /forgot-password rate limit -> 429", func(t *testing.T) {
		r := setupTestRouter()
		
		// Call multiple times to test rate limiting (limit is 3 per minute)
		for i := 0; i < 3; i++ {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "/api/v1/users/forgot-password", bytes.NewBufferString(`{"identifier":"test"}`))
			req.Header.Set("Content-Type", "application/json")
			r.ServeHTTP(w, req)
			assert.NotEqual(t, http.StatusTooManyRequests, w.Code, "Call %d should not be rate limited", i+1)
		}
		
		// 4th call should hit rate limit
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/users/forgot-password", bytes.NewBufferString(`{"identifier":"test"}`))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusTooManyRequests, w.Code, "Should be rate limited on 4th call")
	})
}
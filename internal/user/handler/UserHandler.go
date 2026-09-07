package handler

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"medix-be/internal/common/response"
	"medix-be/internal/user/model/dto"
	"medix-be/internal/user/service"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type UserHandler struct {
	backgroundContext context.Context
	logger            *logrus.Logger
	routers           []*gin.RouterGroup
	userService       service.UserService
}

type UserHandlerProps struct {
	UserService service.UserService
}

type HandlerContract struct {
	BackgroundContext context.Context
	Logger            *logrus.Logger
	Router            *gin.RouterGroup
}

func StartUserHandler(publicContract, authContract *HandlerContract, props *UserHandlerProps) *UserHandler {
	handler := &UserHandler{
		backgroundContext: publicContract.BackgroundContext,
		logger:            publicContract.Logger,
		routers: []*gin.RouterGroup{
			publicContract.Router.Group("/users"),
			authContract.Router.Group("/users"),
		},
		userService: props.UserService,
	}
	handler.RegisterRouter()
	return handler
}

func (h *UserHandler) RegisterRouter() {
	public := h.routers[0]
	protected := h.routers[1]

	public.POST("/login", h.Login())
	public.POST("", h.CreateUser())

	protected.GET("", h.GetAllUsers())
	protected.GET("/:id", h.GetUserByID())
	protected.PUT("/:id", h.UpdateUser())
	protected.DELETE("/:id", h.DeleteUser())
	protected.GET("/profile", h.GetProfile())
	protected.PUT("/profile", h.UpdateProfile())
	protected.POST("/profile/photo", h.UploadProfilePhoto())
}

func (h *UserHandler) CreateUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		var payload dto.CreateUserRequest

		if err := c.ShouldBindJSON(&payload); err != nil {
			response.Error(c, http.StatusBadRequest, err.Error())
			return
		}

		res, err := h.userService.CreateUser(payload)
		if err != nil {
			response.Error(c, http.StatusInternalServerError, err.Error())
			return
		}

		response.Success(c, http.StatusCreated, "User created successfully", res)
	}
}

func (h *UserHandler) GetAllUsers() gin.HandlerFunc {
	return func(c *gin.Context) {

		page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
		limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

		if page < 1 {
			page = 1
		}

		if limit < 1 {
			limit = 10
		}

		search := c.Query("search")
		role := c.Query("role")
		status := c.Query("status")

		res, err := h.userService.GetAllUsers(
			page,
			limit,
			search,
			role,
			status,
		)

		if err != nil {
			response.Error(c, http.StatusInternalServerError, err.Error())
			return
		}

		response.Success(c, http.StatusOK, "User list retrieved successfully", res)
	}
}

func (h *UserHandler) GetUserByID() gin.HandlerFunc {
	return func(c *gin.Context) {
		idParam := c.Param("id")
		id, err := strconv.ParseUint(idParam, 10, 32)
		if err != nil {
			response.Error(c, http.StatusBadRequest, "ID user tidak valid")
			return
		}

		res, err := h.userService.GetUserByID(uint(id))
		if err != nil {
			response.Error(c, http.StatusNotFound, err.Error())
			return
		}

		response.Success(c, http.StatusOK, "User retrieved successfully", res)
	}
}

func (h *UserHandler) UpdateUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		idParam := c.Param("id")
		id, err := strconv.ParseUint(idParam, 10, 32)
		if err != nil {
			response.Error(c, http.StatusBadRequest, "invalid user ID")
			return
		}

		var payload dto.UpdateUserRequest
		if err := c.ShouldBindJSON(&payload); err != nil {
			response.Error(c, http.StatusBadRequest, err.Error())
			return
		}

		res, err := h.userService.UpdateUser(uint(id), payload)
		if err != nil {
			response.Error(c, http.StatusBadRequest, err.Error())
			return
		}

		response.Success(c, http.StatusOK, "User updated successfully", res)
	}
}

func (h *UserHandler) DeleteUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		idParam := c.Param("id")
		id, err := strconv.ParseUint(idParam, 10, 32)
		if err != nil {
			response.Error(c, http.StatusBadRequest, "ID user tidak valid")
			return
		}

		if err := h.userService.DeleteUser(uint(id)); err != nil {
			response.Error(c, http.StatusNotFound, err.Error())
			return
		}

		response.Success(c, http.StatusOK, "User deleted successfully", nil)
	}
}

func (h *UserHandler) Login() gin.HandlerFunc {

	return func(c *gin.Context) {

		var payload dto.LoginRequest

		if err := c.ShouldBindJSON(&payload); err != nil {
			response.Error(c, http.StatusBadRequest, err.Error())
			return
		}

		res, err := h.userService.Login(&payload)
		if err != nil {
			response.Error(c, http.StatusUnauthorized, err.Error())
			return
		}

		response.Success(c, http.StatusOK, "Login successful", res)
	}
}

func (h *UserHandler) GetProfile() gin.HandlerFunc {

	return func(c *gin.Context) {

		value, exists := c.Get("user_id")
		if !exists {
			response.Error(c, http.StatusUnauthorized, "user not authenticated")
			return
		}

		userID, ok := value.(uint)
		if !ok {
			response.Error(c, http.StatusUnauthorized, "invalid user ID")
			return
		}

		res, err := h.userService.GetProfile(userID)
		if err != nil {
			response.Error(c, http.StatusNotFound, err.Error())
			return
		}

		response.Success(c, http.StatusOK, "Profile retrieved successfully", res)
	}
}

func (h *UserHandler) UpdateProfile() gin.HandlerFunc {

	return func(c *gin.Context) {

		value, exists := c.Get("user_id")
		if !exists {
			response.Error(c, http.StatusUnauthorized, "user not authenticated")
			return
		}

		userID, ok := value.(uint)
		if !ok {
			response.Error(c, http.StatusUnauthorized, "invalid user ID")
			return
		}

		var payload dto.UpdateProfileRequest

		if err := c.ShouldBind(&payload); err != nil {
			response.Error(c, http.StatusBadRequest, err.Error())
			return
		}

		if !payload.ValidateNoTelp() {
			response.Error(c, http.StatusBadRequest, "invalid phone number format (example: 08123456789 or +628123456789)")
			return
		}

		res, err := h.userService.UpdateProfile(
			userID,
			&payload,
		)

		if err != nil {
			response.Error(c, http.StatusBadRequest, err.Error())
			return
		}

		response.Success(c, http.StatusOK, "Profile updated successfully", res)
	}
}

func (h *UserHandler) UploadProfilePhoto() gin.HandlerFunc {
	return func(c *gin.Context) {
		value, exists := c.Get("user_id")
		if !exists {
			response.Error(c, http.StatusUnauthorized, "user not authenticated")
			return
		}
		userID, ok := value.(uint)
		if !ok {
			response.Error(c, http.StatusUnauthorized, "invalid user ID")
			return
		}

		file, err := c.FormFile("foto")
		if err != nil {
			response.Error(c, http.StatusBadRequest, "photo must be uploaded")
			return
		}

		if file.Size > 2<<20 {
			response.Error(c, http.StatusBadRequest, "photo size max 2MB")
			return
		}

		ext := filepath.Ext(file.Filename)
		allowed := map[string]bool{".jpg": true, ".jpeg": true, ".png": true, ".webp": true}
		if !allowed[ext] {
			response.Error(c, http.StatusBadRequest, "photo format must be jpg, jpeg, png, or webp")
			return
		}

		os.MkdirAll("uploads/profiles", 0755)
		filename := fmt.Sprintf("user_%d_%d%s", userID, time.Now().UnixNano(), ext)
		savePath := filepath.Join("uploads/profiles", filename)

		if err := c.SaveUploadedFile(file, savePath); err != nil {
			response.Error(c, http.StatusInternalServerError, "failed to save photo")
			return
		}

		res, err := h.userService.UpdateProfilePhoto(userID, savePath)
		if err != nil {
			os.Remove(savePath)
			response.Error(c, http.StatusBadRequest, err.Error())
			return
		}

		response.Success(c, http.StatusOK, "Profile photo updated successfully", res)
	}
}

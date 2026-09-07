package handler

import (
	"context"
	"net/http"
	"strconv"

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

		response.Success(c, http.StatusCreated, "User berhasil dibuat", res)
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

		response.Success(c, http.StatusOK, "Daftar user berhasil diambil", res)
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

		response.Success(c, http.StatusOK, "User berhasil diambil", res)
	}
}

func (h *UserHandler) UpdateUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		idParam := c.Param("id")
		id, err := strconv.ParseUint(idParam, 10, 32)
		if err != nil {
			response.Error(c, http.StatusBadRequest, "ID user tidak valid")
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

		response.Success(c, http.StatusOK, "User berhasil diperbarui", res)
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

		response.Success(c, http.StatusOK, "User berhasil dihapus", nil)
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

		response.Success(c, http.StatusOK, "Login berhasil", res)
	}
}

func (h *UserHandler) GetProfile() gin.HandlerFunc {

	return func(c *gin.Context) {

		value, exists := c.Get("user_id")
		if !exists {
			response.Error(c, http.StatusUnauthorized, "User tidak terautentikasi")
			return
		}

		userID, ok := value.(uint)
		if !ok {
			response.Error(c, http.StatusUnauthorized, "User ID tidak valid")
			return
		}

		res, err := h.userService.GetProfile(userID)
		if err != nil {
			response.Error(c, http.StatusNotFound, err.Error())
			return
		}

		response.Success(c, http.StatusOK, "Profile berhasil diambil", res)
	}
}

func (h *UserHandler) UpdateProfile() gin.HandlerFunc {

	return func(c *gin.Context) {

		value, exists := c.Get("user_id")
		if !exists {
			response.Error(c, http.StatusUnauthorized, "User tidak terautentikasi")
			return
		}

		userID, ok := value.(uint)
		if !ok {
			response.Error(c, http.StatusUnauthorized, "User ID tidak valid")
			return
		}

		var payload dto.UpdateUserRequest

		if err := c.ShouldBindJSON(&payload); err != nil {
			response.Error(c, http.StatusBadRequest, err.Error())
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

		response.Success(c, http.StatusOK, "Profile berhasil diperbarui", res)
	}
}

package handler

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"medix-be/internal/user/model/dto"
	"medix-be/internal/user/service"
	"net/http"
	"strconv"
)

type UserHandler struct {
	backgroundContext context.Context
	logger            *logrus.Logger
	router            *gin.RouterGroup
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

func StartUserHandler(contract *HandlerContract, props *UserHandlerProps) *UserHandler {
	handler := &UserHandler{
		backgroundContext: contract.BackgroundContext,
		logger:            contract.Logger,
		router:            contract.Router.Group("/users"),
		userService:       props.UserService,
	}
	handler.RegisterRouter()
	return handler
}

func (h *UserHandler) RegisterRouter() {
	h.router.POST("/login", h.Login())
	h.router.POST("", h.CreateUser())
	h.router.GET("", h.GetAllUsers())
	h.router.GET("/:id", h.GetUserByID())
	h.router.PUT("/:id", h.UpdateUser())
	h.router.DELETE("/:id", h.DeleteUser())
	h.router.GET("/profile", h.GetProfile())
	h.router.PUT("/profile", h.UpdateProfile())
}

func (h *UserHandler) CreateUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		var payload dto.CreateUserRequest

		if err := c.ShouldBindJSON(&payload); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		res, err := h.userService.CreateUser(payload)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"message": "User created successfully",
			"data":    res,
		})
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
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message":    "Users retrieved successfully",
			"data":       res.Data,
			"pagination": res.Pagination,
		})
	}
}

func (h *UserHandler) GetUserByID() gin.HandlerFunc {
	return func(c *gin.Context) {
		idParam := c.Param("id")
		id, err := strconv.ParseUint(idParam, 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
			return
		}

		res, err := h.userService.GetUserByID(uint(id))
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "User retrieved successfully",
			"data":    res,
		})
	}
}

func (h *UserHandler) UpdateUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		idParam := c.Param("id")
		id, err := strconv.ParseUint(idParam, 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
			return
		}

		var payload dto.UpdateUserRequest
		if err := c.ShouldBindJSON(&payload); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		res, err := h.userService.UpdateUser(uint(id), payload)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "User updated successfully",
			"data":    res,
		})
	}
}

func (h *UserHandler) DeleteUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		idParam := c.Param("id")
		id, err := strconv.ParseUint(idParam, 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
			return
		}

		if err := h.userService.DeleteUser(uint(id)); err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "User deleted successfully",
		})
	}
}

func (h *UserHandler) Login() gin.HandlerFunc {

	return func(c *gin.Context) {

		var payload dto.LoginRequest

		if err := c.ShouldBindJSON(&payload); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		res, err := h.userService.Login(&payload)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Login berhasil",
			"data":    res,
		})
	}
}

func (h *UserHandler) GetProfile() gin.HandlerFunc {

	return func(c *gin.Context) {

		value, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "User tidak terautentikasi",
			})
			return
		}

		userID, ok := value.(uint)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "User ID tidak valid",
			})
			return
		}

		res, err := h.userService.GetProfile(userID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Profile retrieved successfully",
			"data":    res,
		})
	}
}

func (h *UserHandler) UpdateProfile() gin.HandlerFunc {

	return func(c *gin.Context) {

		value, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "User tidak terautentikasi",
			})
			return
		}

		userID, ok := value.(uint)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "User ID tidak valid",
			})
			return
		}

		var payload dto.UpdateUserRequest

		if err := c.ShouldBindJSON(&payload); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		res, err := h.userService.UpdateProfile(
			userID,
			&payload,
		)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Profile updated successfully",
			"data":    res,
		})
	}
}

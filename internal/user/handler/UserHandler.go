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
	h.router.POST("", h.CreateUser())
	h.router.GET("", h.GetAllUsers())
	h.router.GET("/:id", h.GetUserByID())
	h.router.PUT("/:id", h.UpdateUser())
	h.router.DELETE("/:id", h.DeleteUser())
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
		res, err := h.userService.GetAllUsers()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Users retrieved successfully",
			"data":    res,
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

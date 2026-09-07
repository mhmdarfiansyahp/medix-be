package handler

import (
	"context"
	"net/http"
	"strconv"

	"medix-be/internal/common/response"
	"medix-be/internal/transaction/model/dto"
	"medix-be/internal/transaction/service"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type TransactionHandler struct {
	backgroundContext  context.Context
	logger             *logrus.Logger
	router             *gin.RouterGroup
	transactionService service.TransactionService
}

type TransactionHandlerProps struct {
	TransactionService service.TransactionService
}

type HandlerContract struct {
	BackgroundContext context.Context
	Logger            *logrus.Logger
	Router            *gin.RouterGroup
}

func StartTransactionHandler(contract *HandlerContract, props *TransactionHandlerProps) *TransactionHandler {
	handler := &TransactionHandler{
		backgroundContext:  contract.BackgroundContext,
		logger:             contract.Logger,
		router:             contract.Router.Group("/transactions"),
		transactionService: props.TransactionService,
	}

	handler.RegisterRouter()
	return handler
}

func (h *TransactionHandler) RegisterRouter() {
	h.router.POST("", h.Create())
	h.router.GET("", h.GetAll())
	h.router.GET("/today", h.GetToday())
	h.router.GET("/:id", h.GetByID())
	h.router.PATCH("/:id/cancel", h.Cancel())
	h.router.GET("/:id/receipt", h.GetReceipt())
}

func (h *TransactionHandler) Create() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req dto.CreateTransactionRequest

		if err := c.ShouldBindJSON(&req); err != nil {
			response.Error(c, http.StatusBadRequest, err.Error())
			return
		}

		userIDValue, exists := c.Get("user_id")
		if !exists {
			response.Error(c, http.StatusUnauthorized, "user not authenticated")
			return
		}

		userID, ok := userIDValue.(uint)
		if !ok {
			response.Error(c, http.StatusUnauthorized, "invalid user ID")
			return
		}

		res, err := h.transactionService.CreateTransaction(
			userID,
			req,
		)

		if err != nil {
			response.Error(c, http.StatusBadRequest, err.Error())
			return
		}

		response.Success(c, http.StatusCreated, "Transaction saved successfully", res)
	}
}

func (h *TransactionHandler) GetAll() gin.HandlerFunc {
	return func(c *gin.Context) {
		res, err := h.transactionService.GetAllTransactions()
		if err != nil {
			response.Error(c, http.StatusInternalServerError, err.Error())
			return
		}

		response.Success(c, http.StatusOK, "Transaction list retrieved successfully", res)
	}
}

func (h *TransactionHandler) GetByID() gin.HandlerFunc {
	return func(c *gin.Context) {

		id, err := strconv.ParseUint(c.Param("id"), 10, 64)

		if err != nil {
			response.Error(c, http.StatusBadRequest, "ID transaksi tidak valid")
			return
		}

		res, err := h.transactionService.GetTransactionByID(uint(id))

		if err != nil {
			response.Error(c, http.StatusNotFound, err.Error())
			return
		}

		response.Success(c, http.StatusOK, "Transaction retrieved successfully", res)
	}
}

func (h *TransactionHandler) Cancel() gin.HandlerFunc {
	return func(c *gin.Context) {

		id, err := strconv.ParseUint(c.Param("id"), 10, 64)

		if err != nil {
			response.Error(c, http.StatusBadRequest, "ID transaksi tidak valid")
			return
		}

		userIDValue, exists := c.Get("user_id")

		if !exists {
			response.Error(c, http.StatusUnauthorized, "user not authenticated")
			return
		}

		userID, ok := userIDValue.(uint)

		if !ok {
			response.Error(c, http.StatusUnauthorized, "invalid user ID")
			return
		}

		err = h.transactionService.CancelTransaction(
			uint(id),
			userID,
		)

		if err != nil {
			response.Error(c, http.StatusBadRequest, err.Error())
			return
		}

		response.Success(c, http.StatusOK, "Transaction cancelled successfully", nil)
	}
}

func (h *TransactionHandler) GetToday() gin.HandlerFunc {
	return func(c *gin.Context) {
		userIDValue, exists := c.Get("user_id")

		if !exists {
			response.Error(c, http.StatusUnauthorized, "user not authenticated")
			return
		}

		userID, ok := userIDValue.(uint)
		if !ok {
			response.Error(c, http.StatusUnauthorized, "invalid user ID")
			return
		}

		res, err := h.transactionService.
			GetTodayTransactions(userID)

		if err != nil {
			response.Error(c, http.StatusInternalServerError, err.Error())
			return
		}

		response.Success(c, http.StatusOK, "Today's transaction history retrieved successfully", res)
	}
}

func (h *TransactionHandler) GetReceipt() gin.HandlerFunc {
	return func(c *gin.Context) {
		idParam := c.Param("id")

		id, err := strconv.ParseUint(idParam, 10, 64)
		if err != nil {
			response.Error(c, http.StatusBadRequest, "ID transaksi tidak valid")
			return
		}

		res, err := h.transactionService.GetReceipt(uint(id))

		if err != nil {
			response.Error(c, http.StatusNotFound, err.Error())
			return
		}

		response.Success(c, http.StatusOK, "Transaction receipt retrieved successfully", res)
	}
}

package handler

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"medix-be/internal/transaction/model/dto"
	"medix-be/internal/transaction/service"
	"net/http"
	"strconv"
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
	h.router.GET("/:id", h.GetByID())
}

func (h *TransactionHandler) Create() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req dto.CreateTransactionRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		res, err := h.transactionService.CreateTransaction(req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"message": "Transaksi berhasil disimpan",
			"data":    res,
		})
	}
}

func (h *TransactionHandler) GetAll() gin.HandlerFunc {
	return func(c *gin.Context) {
		res, err := h.transactionService.GetAllTransactions()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Daftar transaksi berhasil diambil",
			"data":    res,
		})
	}
}

func (h *TransactionHandler) GetByID() gin.HandlerFunc {
	return func(c *gin.Context) {
		idParam := c.Param("id")
		id, err := strconv.Atoi(idParam)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "ID transaksi tidak valid"})
			return
		}

		res, err := h.transactionService.GetTransactionByID(uint(id))
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Transaksi berhasil diambil",
			"data":    res,
		})
	}
}

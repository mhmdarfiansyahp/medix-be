package handler

import (
	"context"
	"net/http"
	"strconv"

	"medix-be/internal/common/response"
	"medix-be/internal/transaction/model/dto"
	"medix-be/internal/transaction/service"
	"medix-be/internal/middleware"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type TransactionHandler struct {
	backgroundContext  context.Context
	logger             *logrus.Logger
	router             *gin.RouterGroup
	transactionService *service.TransactionService
}

type TransactionHandlerProps struct {
	TransactionService *service.TransactionService
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
	h.router.POST("/add-to-cart", h.AddToCart())
	h.router.GET("", h.GetAll())
	h.router.GET("/today", h.GetToday())
	h.router.GET("/kasir-report", h.GetTodayKasirReport())
	h.router.GET("/:id", h.GetByID())
	h.router.GET("/:id/price-snapshot", h.GetTransactionByIDWithPriceSnapshot())
	h.router.GET("/:id/price-history", h.GetTransactionPriceHistory())
	h.router.PATCH("/:id/cancel", h.Cancel())
	h.router.GET("/:id/receipt", h.GetReceipt())
	h.router.POST("/:id/payment", h.ProcessPayment())
	h.router.POST("/returns", h.CreateReturn())
	h.router.GET("/returns/:id/view", h.ViewReturnDetails())
	h.router.POST("/returns/:id/reject", middleware.RequireRoles("admin", "owner"), h.RejectReturn())
	h.router.PUT("/returns/threshold", middleware.RequireRoles("admin", "owner"), h.ConfigureApprovalThreshold())
	h.router.GET("/returns/threshold", h.GetApprovalThreshold())
	h.router.POST("/returns/:id/approve", middleware.RequireRoles("admin", "owner"), h.ApproveReturn())
	h.router.GET("/:id/receipt/generate", h.GenerateReceipt())
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

func (h *TransactionHandler) AddToCart() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req dto.AddToCartRequest

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

		res, err := h.transactionService.AddToCart(userID, req)

		if err != nil {
			response.Error(c, http.StatusBadRequest, err.Error())
			return
		}

		response.Success(c, http.StatusOK, "Medicine added to cart successfully", res)
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

		roleValue, exists := c.Get("role")
		if !exists {
			response.Error(c, http.StatusUnauthorized, "user not authenticated")
			return
		}
		userRole, ok := roleValue.(string)
		if !ok {
			response.Error(c, http.StatusUnauthorized, "invalid role")
			return
		}

		err = h.transactionService.CancelTransaction(
			uint(id),
			userID,
			userRole,
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

		res, err := h.transactionService.GetTodayTransactions(userID)

		if err != nil {
			response.Error(c, http.StatusInternalServerError, err.Error())
			return
		}

		response.Success(c, http.StatusOK, "Today's transaction history retrieved successfully", res)
	}
}

func (h *TransactionHandler) GetTodayKasirReport() gin.HandlerFunc {
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

		kasirReport, err := h.transactionService.GetTodayKasirReport(userID)
		if err != nil {
			response.Error(c, http.StatusInternalServerError, err.Error())
			return
		}

		response.Success(c, http.StatusOK, "Kasir daily report retrieved successfully", kasirReport)
	}
}

func (h *TransactionHandler) GetTransactionByIDWithPriceSnapshot() gin.HandlerFunc {
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

		idParam := c.Param("id")
		id, err := strconv.ParseUint(idParam, 10, 64)
		if err != nil {
			response.Error(c, http.StatusBadRequest, "ID transaksi tidak valid")
			return
		}

		res, err := h.transactionService.GetTransactionByID(uint(id))
		if err != nil {
			response.Error(c, http.StatusNotFound, err.Error())
			return
		}

		if res == nil {
			response.Error(c, http.StatusNotFound, "transaction not found")
			return
		}

		if res.IDUser != userID {
			response.Error(c, http.StatusForbidden, "you do not have access to this transaction")
			return
		}

		response.Success(c, http.StatusOK, "Transaction with price snapshot retrieved successfully", res)
	}
}

func (h *TransactionHandler) GetTransactionPriceHistory() gin.HandlerFunc {
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

		idParam := c.Param("id")
		id, err := strconv.ParseUint(idParam, 10, 64)
		if err != nil {
			response.Error(c, http.StatusBadRequest, "ID transaksi tidak valid")
			return
		}

		transaction, err := h.transactionService.GetTransactionByID(uint(id))
		if err != nil {
			response.Error(c, http.StatusNotFound, err.Error())
			return
		}

		if transaction == nil {
			response.Error(c, http.StatusNotFound, "transaction not found")
			return
		}

		if transaction.IDUser != userID {
			response.Error(c, http.StatusForbidden, "you do not have access to this transaction")
			return
		}

		response.Success(c, http.StatusOK, "Transaction price history retrieved successfully", transaction)
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

func (h *TransactionHandler) ProcessPayment() gin.HandlerFunc {
	return func(c *gin.Context) {
		idParam := c.Param("id")

		id, err := strconv.ParseUint(idParam, 10, 64)
		if err != nil {
			response.Error(c, http.StatusBadRequest, "ID transaksi tidak valid")
			return
		}

		var req dto.PaymentRequest
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

		res, err := h.transactionService.ProcessPayment(
			uint(id),
			userID,
			req,
		)

		if err != nil {
			response.Error(c, http.StatusBadRequest, err.Error())
			return
		}

		response.Success(c, http.StatusOK, "Payment processed successfully", res)
	}
}

func (h *TransactionHandler) CreateReturn() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req dto.CreateReturnRequest
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

		res, err := h.transactionService.CreateReturn(userID, req)
		if err != nil {
			response.Error(c, http.StatusBadRequest, err.Error())
			return
		}

		response.Success(c, http.StatusOK, "Return created successfully", res)
	}
}

func (h *TransactionHandler) ViewReturnDetails() gin.HandlerFunc {
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

		returnID, err := strconv.ParseUint(c.Param("id"), 10, 64)
		if err != nil {
			response.Error(c, http.StatusBadRequest, "invalid return ID")
			return
		}

		returnData, err := h.transactionService.GetReturnByID(uint(returnID))
		if err != nil {
			response.Error(c, http.StatusNotFound, err.Error())
			return
		}

		role, _ := c.Get("role")
		roleStr, _ := role.(string)
		isAdmin := roleStr == "admin" || roleStr == "owner"

		if returnData.DiajukanOleh != userID && !isAdmin {
			response.Error(c, http.StatusForbidden, "Anda hanya dapat melihat retur milik sendiri")
			return
		}

		response.Success(c, http.StatusOK, "Return details retrieved successfully", returnData)
	}
}

func (h *TransactionHandler) RejectReturn() gin.HandlerFunc {
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

		returnID, err := strconv.ParseUint(c.Param("id"), 10, 64)
		if err != nil {
			response.Error(c, http.StatusBadRequest, "invalid return ID")
			return
		}

		var req dto.RejectReturnRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			response.Error(c, http.StatusBadRequest, err.Error())
			return
		}

		adminID := userID

		res, err := h.transactionService.RejectReturn(adminID, uint(returnID), req.Alasan)
		if err != nil {
			response.Error(c, http.StatusBadRequest, err.Error())
			return
		}

		response.Success(c, http.StatusOK, "Return rejected successfully", res)
	}
}

func (h *TransactionHandler) ConfigureApprovalThreshold() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req dto.ApprovalThresholdRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			response.Error(c, http.StatusBadRequest, err.Error())
			return
		}

		threshold, err := h.transactionService.SetApprovalThreshold(req.Threshold)
		if err != nil {
			response.Error(c, http.StatusInternalServerError, err.Error())
			return
		}

		response.Success(c, http.StatusOK, "Batas approval retur berhasil dikonfigurasi", threshold)
	}
}

func (h *TransactionHandler) GetApprovalThreshold() gin.HandlerFunc {
	return func(c *gin.Context) {
		threshold, err := h.transactionService.GetApprovalThreshold()
		if err != nil {
			response.Error(c, http.StatusInternalServerError, err.Error())
			return
		}

		response.Success(c, http.StatusOK, "Batas approval retur berhasil diambil", threshold)
	}
}

func (h *TransactionHandler) ApproveReturn() gin.HandlerFunc {
	return func(c *gin.Context) {
		returnID, err := strconv.ParseUint(c.Param("id"), 10, 64)
		if err != nil {
			response.Error(c, http.StatusBadRequest, "invalid return ID")
			return
		}

		value, exists := c.Get("user_id")
		if !exists {
			response.Error(c, http.StatusUnauthorized, "user not authenticated")
			return
		}

		adminID, ok := value.(uint)
		if !ok {
			response.Error(c, http.StatusUnauthorized, "invalid user ID")
			return
		}

		res, err := h.transactionService.ApproveReturn(adminID, uint(returnID))
		if err != nil {
			response.Error(c, http.StatusBadRequest, err.Error())
			return
		}

		response.Success(c, http.StatusOK, "Return approved successfully", res)
	}
}

func (h *TransactionHandler) GenerateReceipt() gin.HandlerFunc {
	return func(c *gin.Context) {
		idParam := c.Param("id")
		id, err := strconv.ParseUint(idParam, 10, 64)
		if err != nil {
			response.Error(c, http.StatusBadRequest, "ID transaksi tidak valid")
			return
		}

		format := c.Query("format")
		if format == "" {
			format = "whatsapp"
		}

		res, err := h.transactionService.GenerateReceipt(uint(id), format)
		if err != nil {
			response.Error(c, http.StatusBadRequest, err.Error())
			return
		}

		response.Success(c, http.StatusOK, "Receipt generated successfully", res)
	}
}


package handler

import (
	"context"
	"net/http"
	"strconv"
	"strings"

	"medix-be/internal/common/response"
	"medix-be/internal/medicine/model/dto"
	"medix-be/internal/medicine/service"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type MedicineHandler struct {
	backgroundContext context.Context
	logger            *logrus.Logger
	router            *gin.RouterGroup
	medicineService   service.MedicineService
}

type MedicineHandlerProps struct {
	MedicineService service.MedicineService
}

type HandlerContract struct {
	BackgroundContext context.Context
	Logger            *logrus.Logger
	Router            *gin.RouterGroup
}

func StartMedicineHandler(contract *HandlerContract, props *MedicineHandlerProps) *MedicineHandler {
	handler := &MedicineHandler{
		backgroundContext: contract.BackgroundContext,
		logger:            contract.Logger,
		router:            contract.Router.Group("/medicines"),
		medicineService:   props.MedicineService,
	}

	handler.RegisterRouter()
	return handler
}

func (h *MedicineHandler) RegisterRouter() {
	h.router.GET("/alerts/low-stock", h.GetLowStock())
	h.router.GET("/alerts/expiring", h.GetExpiring())
	h.router.GET("/alerts/summary", h.GetNotificationSummary())

	h.router.POST("", h.CreateMedicine())
	h.router.GET("", h.GetAllMedicines())
	h.router.GET("/:id", h.GetMedicineByID())
	h.router.GET("/barcode/:barcode", h.GetMedicineByBarcode())
	h.router.PUT("/:id", h.UpdateMedicine())
	h.router.PATCH("/:id/status", h.ToggleActiveStatus())
	h.router.DELETE("/:id", h.DeleteMedicine())
}

func (h *MedicineHandler) CreateMedicine() gin.HandlerFunc {
	return func(c *gin.Context) {
		var payload dto.CreateMedicineRequest

		if err := c.ShouldBindJSON(&payload); err != nil {
			response.Error(c, http.StatusBadRequest, err.Error())
			return
		}

		res, err := h.medicineService.CreateMedicine(payload)
		if err != nil {
			if strings.Contains(err.Error(), "barcode") &&
				(strings.Contains(err.Error(), "already used") || strings.Contains(err.Error(), "23505")) {
				response.Error(c, http.StatusConflict, err.Error())
				return
			}
			response.Error(c, http.StatusInternalServerError, err.Error())
			return
		}

		response.Success(c, http.StatusCreated, "Medicine added successfully", res)
	}
}

func (h *MedicineHandler) GetAllMedicines() gin.HandlerFunc {
	return func(c *gin.Context) {
		var filter dto.MedicineFilterParams

		if err := c.ShouldBindQuery(&filter); err != nil {
			response.Error(c, http.StatusBadRequest, err.Error())
			return
		}

		res, err := h.medicineService.GetAllMedicines(c.Request.Context(), filter)
		if err != nil {
			response.Error(c, http.StatusInternalServerError, err.Error())
			return
		}

		response.Success(c, http.StatusOK, "Medicine retrieved successfully", res)
	}
}

func (h *MedicineHandler) GetMedicineByID() gin.HandlerFunc {
	return func(c *gin.Context) {
		idParam := c.Param("id")
		id, err := strconv.ParseUint(idParam, 10, 32)
		if err != nil {
			response.Error(c, http.StatusBadRequest, "invalid medicine ID")
			return
		}

		res, err := h.medicineService.GetMedicineByID(uint(id))
		if err != nil {
			response.Error(c, http.StatusNotFound, err.Error())
			return
		}

		response.Success(c, http.StatusOK, "Medicine retrieved successfully", res)
	}
}

func (h *MedicineHandler) GetMedicineByBarcode() gin.HandlerFunc {
	return func(c *gin.Context) {
		barcode := c.Param("barcode")
		if barcode == "" {
			response.Error(c, http.StatusBadRequest, "barcode cannot be empty")
			return
		}

		res, err := h.medicineService.GetMedicineByBarcode(c.Request.Context(), barcode)
		if err != nil {
			response.Error(c, http.StatusNotFound, err.Error())
			return
		}

		response.Success(c, http.StatusOK, "Medicine retrieved successfully", res)
	}
}

func (h *MedicineHandler) UpdateMedicine() gin.HandlerFunc {
	return func(c *gin.Context) {
		idParam := c.Param("id")
		id, err := strconv.ParseUint(idParam, 10, 32)
		if err != nil {
			response.Error(c, http.StatusBadRequest, "invalid medicine ID")
			return
		}

		var payload dto.UpdateMedicineRequest
		if err := c.ShouldBindJSON(&payload); err != nil {
			response.Error(c, http.StatusBadRequest, err.Error())
			return
		}

		res, err := h.medicineService.UpdateMedicine(uint(id), payload)
		if err != nil {
			if strings.Contains(err.Error(), "barcode") &&
				(strings.Contains(err.Error(), "already used") || strings.Contains(err.Error(), "23505")) {
				response.Error(c, http.StatusConflict, err.Error())
				return
			}
			response.Error(c, http.StatusBadRequest, err.Error())
			return
		}

		response.Success(c, http.StatusOK, "Medicine updated successfully", res)
	}
}

func (h *MedicineHandler) ToggleActiveStatus() gin.HandlerFunc {
	return func(c *gin.Context) {
		idParam := c.Param("id")
		id, err := strconv.ParseUint(idParam, 10, 32)
		if err != nil {
			response.Error(c, http.StatusBadRequest, "invalid medicine ID")
			return
		}

		var payload struct {
			IsActive bool `json:"is_active"`
		}

		if err := c.ShouldBindJSON(&payload); err != nil {
			response.Error(c, http.StatusBadRequest, err.Error())
			return
		}

		err = h.medicineService.ToggleActiveStatus(c.Request.Context(), uint(id), payload.IsActive)
		if err != nil {
			response.Error(c, http.StatusInternalServerError, err.Error())
			return
		}

		response.Success(c, http.StatusOK, "Medicine status updated successfully", nil)
	}
}

func (h *MedicineHandler) DeleteMedicine() gin.HandlerFunc {
	return func(c *gin.Context) {
		idParam := c.Param("id")
		id, err := strconv.ParseUint(idParam, 10, 32)
		if err != nil {
			response.Error(c, http.StatusBadRequest, "invalid medicine ID")
			return
		}

		if err := h.medicineService.DeleteMedicine(uint(id)); err != nil {
			response.Error(c, http.StatusNotFound, err.Error())
			return
		}

		response.Success(c, http.StatusOK, "Medicine deleted successfully", nil)
	}
}

// GET /medicines/alerts/low-stock (US-13)
func (h *MedicineHandler) GetLowStock() gin.HandlerFunc {
	return func(c *gin.Context) {
		res, err := h.medicineService.GetLowStockDrugs(c.Request.Context())
		if err != nil {
			response.Error(c, http.StatusInternalServerError, err.Error())
			return
		}

		response.Success(c, http.StatusOK, "Successfully retrieved list of low stock medicines", res)
	}
}

// GET /medicines/alerts/expiring?days=30 (US-14)
func (h *MedicineHandler) GetExpiring() gin.HandlerFunc {
	return func(c *gin.Context) {
		daysParam := c.DefaultQuery("days", "30")
		days, err := strconv.Atoi(daysParam)
		if err != nil {
			days = 30
		}

		res, err := h.medicineService.GetExpiringDrugs(c.Request.Context(), days)
		if err != nil {
			response.Error(c, http.StatusInternalServerError, err.Error())
			return
		}

		response.Success(c, http.StatusOK, "Successfully retrieved list of expiring medicines", res)
	}
}

// GET /medicines/alerts/summary
func (h *MedicineHandler) GetNotificationSummary() gin.HandlerFunc {
	return func(c *gin.Context) {
		res, err := h.medicineService.GetNotificationSummary(c.Request.Context())
		if err != nil {
			response.Error(c, http.StatusInternalServerError, err.Error())
			return
		}

		response.Success(c, http.StatusOK, "Successfully retrieved stock notification summary", res)
	}
}

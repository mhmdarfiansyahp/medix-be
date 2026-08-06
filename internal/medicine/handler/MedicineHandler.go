package handler

import (
	"context"
	"net/http"
	"strconv"

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
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		res, err := h.medicineService.CreateMedicine(payload)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"message": "Obat berhasil ditambahkan",
			"data":    res,
		})
	}
}

func (h *MedicineHandler) GetAllMedicines() gin.HandlerFunc {
	return func(c *gin.Context) {
		var filter dto.MedicineFilterParams

		if err := c.ShouldBindQuery(&filter); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		res, err := h.medicineService.GetAllMedicines(c.Request.Context(), filter)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"data": res})
	}
}

func (h *MedicineHandler) GetMedicineByID() gin.HandlerFunc {
	return func(c *gin.Context) {
		idParam := c.Param("id")
		id, err := strconv.ParseUint(idParam, 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "ID obat tidak valid"})
			return
		}

		res, err := h.medicineService.GetMedicineByID(uint(id))
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"data": res})
	}
}

func (h *MedicineHandler) GetMedicineByBarcode() gin.HandlerFunc {
	return func(c *gin.Context) {
		barcode := c.Param("barcode")
		if barcode == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Barcode tidak boleh kosong"})
			return
		}

		res, err := h.medicineService.GetMedicineByBarcode(c.Request.Context(), barcode)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"data": res})
	}
}

func (h *MedicineHandler) UpdateMedicine() gin.HandlerFunc {
	return func(c *gin.Context) {
		idParam := c.Param("id")
		id, err := strconv.ParseUint(idParam, 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "ID obat tidak valid"})
			return
		}

		var payload dto.UpdateMedicineRequest
		if err := c.ShouldBindJSON(&payload); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		res, err := h.medicineService.UpdateMedicine(uint(id), payload)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Obat berhasil diperbarui",
			"data":    res,
		})
	}
}

func (h *MedicineHandler) ToggleActiveStatus() gin.HandlerFunc {
	return func(c *gin.Context) {
		idParam := c.Param("id")
		id, err := strconv.ParseUint(idParam, 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "ID obat tidak valid"})
			return
		}

		var payload struct {
			IsActive bool `json:"is_active"`
		}

		if err := c.ShouldBindJSON(&payload); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		err = h.medicineService.ToggleActiveStatus(c.Request.Context(), uint(id), payload.IsActive)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Status obat berhasil diperbarui"})
	}
}

func (h *MedicineHandler) DeleteMedicine() gin.HandlerFunc {
	return func(c *gin.Context) {
		idParam := c.Param("id")
		id, err := strconv.ParseUint(idParam, 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "ID obat tidak valid"})
			return
		}

		if err := h.medicineService.DeleteMedicine(uint(id)); err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Obat berhasil dihapus"})
	}
}

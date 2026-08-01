package handler

import (
	"context"
	"net/http"
	"strconv"

	"medix-be/internal/drug/model/dto"
	"medix-be/internal/drug/service"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type TypeDrugHandler struct {
	backgroundContext context.Context
	logger            *logrus.Logger
	router            *gin.RouterGroup
	typeDrugService   service.TypeDrugService
}

type TypeDrugHandlerProps struct {
	TypeDrugService service.TypeDrugService
}

type HandlerContract struct {
	BackgroundContext context.Context
	Logger            *logrus.Logger
	Router            *gin.RouterGroup
}

func StartTypeDrugHandler(contract *HandlerContract, props *TypeDrugHandlerProps) *TypeDrugHandler {
	handler := &TypeDrugHandler{
		backgroundContext: contract.BackgroundContext,
		logger:            contract.Logger,
		router:            contract.Router.Group("/type-drugs"),
		typeDrugService:   props.TypeDrugService,
	}

	handler.RegisterRouter()
	return handler
}

func (h *TypeDrugHandler) RegisterRouter() {
	h.router.POST("", h.CreateTypeDrug())
	h.router.GET("", h.GetAllTypeDrugs())
	h.router.GET("/:id", h.GetTypeDrugByID())
	h.router.PUT("/:id", h.UpdateTypeDrug())
	h.router.DELETE("/:id", h.DeleteTypeDrug())
}

func (h *TypeDrugHandler) CreateTypeDrug() gin.HandlerFunc {
	return func(c *gin.Context) {
		var payload dto.CreateTypeDrugRequest

		if err := c.ShouldBindJSON(&payload); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		res, err := h.typeDrugService.CreateTypeDrug(payload)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"message": "Jenis obat berhasil ditambahkan",
			"data":    res,
		})
	}
}

func (h *TypeDrugHandler) GetAllTypeDrugs() gin.HandlerFunc {
	return func(c *gin.Context) {
		res, err := h.typeDrugService.GetAllTypeDrugs()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"data": res})
	}
}

func (h *TypeDrugHandler) GetTypeDrugByID() gin.HandlerFunc {
	return func(c *gin.Context) {
		idParam := c.Param("id")
		id, err := strconv.ParseUint(idParam, 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "ID jenis obat tidak valid"})
			return
		}

		res, err := h.typeDrugService.GetTypeDrugByID(uint(id))
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"data": res})
	}
}

func (h *TypeDrugHandler) UpdateTypeDrug() gin.HandlerFunc {
	return func(c *gin.Context) {
		idParam := c.Param("id")
		id, err := strconv.ParseUint(idParam, 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "ID jenis obat tidak valid"})
			return
		}

		var payload dto.UpdateTypeDrugRequest
		if err := c.ShouldBindJSON(&payload); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		res, err := h.typeDrugService.UpdateTypeDrug(uint(id), payload)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Jenis obat berhasil diperbarui",
			"data":    res,
		})
	}
}

func (h *TypeDrugHandler) DeleteTypeDrug() gin.HandlerFunc {
	return func(c *gin.Context) {
		idParam := c.Param("id")
		id, err := strconv.ParseUint(idParam, 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "ID jenis obat tidak valid"})
			return
		}

		if err := h.typeDrugService.DeleteTypeDrug(uint(id)); err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Jenis obat berhasil dihapus"})
	}
}

package handler

import (
	"context"
	"net/http"
	"strconv"

	"medix-be/internal/common/response"
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
			response.Error(c, http.StatusBadRequest, err.Error())
			return
		}

		res, err := h.typeDrugService.CreateTypeDrug(payload)
		if err != nil {
			response.Error(c, http.StatusInternalServerError, err.Error())
			return
		}

		response.Success(c, http.StatusCreated, "Drug type added successfully", res)
	}
}

func (h *TypeDrugHandler) GetAllTypeDrugs() gin.HandlerFunc {
	return func(c *gin.Context) {
		res, err := h.typeDrugService.GetAllTypeDrugs()
		if err != nil {
			response.Error(c, http.StatusInternalServerError, err.Error())
			return
		}

		response.Success(c, http.StatusOK, "Drug type retrieved successfully", res)
	}
}

func (h *TypeDrugHandler) GetTypeDrugByID() gin.HandlerFunc {
	return func(c *gin.Context) {
		idParam := c.Param("id")
		id, err := strconv.ParseUint(idParam, 10, 32)
		if err != nil {
			response.Error(c, http.StatusBadRequest, "invalid drug type ID")
			return
		}

		res, err := h.typeDrugService.GetTypeDrugByID(uint(id))
		if err != nil {
			response.Error(c, http.StatusNotFound, err.Error())
			return
		}

		response.Success(c, http.StatusOK, "Drug type retrieved successfully", res)
	}
}

func (h *TypeDrugHandler) UpdateTypeDrug() gin.HandlerFunc {
	return func(c *gin.Context) {
		idParam := c.Param("id")
		id, err := strconv.ParseUint(idParam, 10, 32)
		if err != nil {
			response.Error(c, http.StatusBadRequest, "invalid drug type ID")
			return
		}

		var payload dto.UpdateTypeDrugRequest
		if err := c.ShouldBindJSON(&payload); err != nil {
			response.Error(c, http.StatusBadRequest, err.Error())
			return
		}

		res, err := h.typeDrugService.UpdateTypeDrug(uint(id), payload)
		if err != nil {
			response.Error(c, http.StatusBadRequest, err.Error())
			return
		}

		response.Success(c, http.StatusOK, "Drug type updated successfully", res)
	}
}

func (h *TypeDrugHandler) DeleteTypeDrug() gin.HandlerFunc {
	return func(c *gin.Context) {
		idParam := c.Param("id")
		id, err := strconv.ParseUint(idParam, 10, 32)
		if err != nil {
			response.Error(c, http.StatusBadRequest, "invalid drug type ID")
			return
		}

		if err := h.typeDrugService.DeleteTypeDrug(uint(id)); err != nil {
			response.Error(c, http.StatusNotFound, err.Error())
			return
		}

		response.Success(c, http.StatusOK, "Drug type deleted successfully", nil)
	}
}

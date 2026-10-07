package handler

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"medix-be/internal/common/response"
	"medix-be/internal/report/model/dto"
	"medix-be/internal/report/service"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type ReportHandler struct {
	backgroundContext context.Context
	logger            *logrus.Logger
	router            *gin.RouterGroup
	reportService     service.ReportService
}

type ReportHandlerProps struct {
	ReportService service.ReportService
}

type HandlerContract struct {
	BackgroundContext context.Context
	Logger            *logrus.Logger
	Router            *gin.RouterGroup
}

func StartReportHandler(contract *HandlerContract, props *ReportHandlerProps) *ReportHandler {
	handler := &ReportHandler{
		backgroundContext: contract.BackgroundContext,
		logger:            contract.Logger,
		router:            contract.Router.Group("/reports"),
		reportService:     props.ReportService,
	}

	handler.RegisterRouter()
	return handler
}

func (h *ReportHandler) RegisterRouter() {
		h.router.GET("/sales-summary", h.GetSalesSummary()) // US-15
		h.router.GET("/drug-ranking", h.GetDrugRanking())   // US-16
		h.router.GET("/export/excel", h.ExportExcel())       // US-17 sync
		h.router.POST("/export/excel", h.ExportExcelAsync()) // US-17 async
		h.router.POST("/export/pdf", h.ExportPDFAsync())     // US-17 async
		h.router.GET("/export/status/:exportID", h.GetExportStatus()) // US-17 status
		h.router.GET("/export/download/:exportID", h.DownloadExport()) // US-17 download
	}

// GET /reports/sales-summary?start_date=2026-01-01&end_date=2026-01-31&group_by=daily
func (h *ReportHandler) GetSalesSummary() gin.HandlerFunc {
	return func(c *gin.Context) {
		var params dto.ReportFilterParams
		if err := c.ShouldBindQuery(&params); err != nil {
			response.Error(c, http.StatusBadRequest, err.Error())
			return
		}

		res, err := h.reportService.GetSalesSummary(c.Request.Context(), params)
		if err != nil {
			response.Error(c, http.StatusInternalServerError, err.Error())
			return
		}

		response.Success(c, http.StatusOK, "Sales summary retrieved successfully", res)
	}
}

// GET /reports/drug-ranking?start_date=2026-01-01&end_date=2026-01-31
func (h *ReportHandler) GetDrugRanking() gin.HandlerFunc {
	return func(c *gin.Context) {
		var params dto.ReportFilterParams
		if err := c.ShouldBindQuery(&params); err != nil {
			response.Error(c, http.StatusBadRequest, err.Error())
			return
		}

		res, err := h.reportService.GetDrugRanking(c.Request.Context(), params)
		if err != nil {
			response.Error(c, http.StatusInternalServerError, err.Error())
			return
		}

		response.Success(c, http.StatusOK, "Drug ranking retrieved successfully", res)
	}
}

// GET /reports/export/excel?start_date=2026-01-01&end_date=2026-01-31
func (h *ReportHandler) ExportExcel() gin.HandlerFunc {
	return func(c *gin.Context) {
		var params dto.ReportFilterParams
		if err := c.ShouldBindQuery(&params); err != nil {
			response.Error(c, http.StatusBadRequest, err.Error())
			return
		}

		buf, err := h.reportService.ExportToExcel(c.Request.Context(), params)
		if err != nil {
			response.Error(c, http.StatusInternalServerError, err.Error())
			return
		}

		fileName := fmt.Sprintf("Sales_Report_%s.xlsx", time.Now().Format("20060102_150405"))
		c.Header("Content-Description", "File Transfer")
		c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", fileName))
		c.Data(http.StatusOK, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", buf.Bytes())
	}
}

// POST /reports/export/excel - Start async Excel export
func (h *ReportHandler) ExportExcelAsync() gin.HandlerFunc {
	return func(c *gin.Context) {
		var params dto.ReportFilterParams
		if err := c.ShouldBindQuery(&params); err != nil {
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

		exportID, err := h.reportService.ExportToExcelAsync(c.Request.Context(), params, userID)
		if err != nil {
			response.Error(c, http.StatusInternalServerError, err.Error())
			return
		}

		response.Success(c, http.StatusAccepted, "Export job started", map[string]string{
			"export_id": exportID,
			"status_url": "/api/v1/reports/export/status/" + exportID,
		})
	}
}

// POST /reports/export/pdf - Start async PDF export
func (h *ReportHandler) ExportPDFAsync() gin.HandlerFunc {
	return func(c *gin.Context) {
		var params dto.ReportFilterParams
		if err := c.ShouldBindQuery(&params); err != nil {
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

		exportID, err := h.reportService.ExportToPDFAsync(c.Request.Context(), params, userID)
		if err != nil {
			response.Error(c, http.StatusInternalServerError, err.Error())
			return
		}

		response.Success(c, http.StatusAccepted, "Export job started", map[string]string{
			"export_id": exportID,
			"status_url": "/api/v1/reports/export/status/" + exportID,
		})
	}
}

// GET /reports/export/status/:exportID - Check export job status
func (h *ReportHandler) GetExportStatus() gin.HandlerFunc {
	return func(c *gin.Context) {
		exportID := c.Param("exportID")

		status, err := h.reportService.GetExportStatus(exportID)
		if err != nil {
			response.Error(c, http.StatusNotFound, "Export job not found")
			return
		}

		response.Success(c, http.StatusOK, "Export status", status)
	}
}

// GET /reports/export/download/:exportID - Download completed export file
func (h *ReportHandler) DownloadExport() gin.HandlerFunc {
	return func(c *gin.Context) {
		exportID := c.Param("exportID")

		status, err := h.reportService.GetExportStatus(exportID)
		if err != nil {
			response.Error(c, http.StatusNotFound, "Export job not found")
			return
		}

		if status.Status != "completed" {
			response.Error(c, http.StatusBadRequest, "Export not completed yet")
			return
		}

		filePath := status.FileURL
		if filePath == "" {
			response.Error(c, http.StatusInternalServerError, "File path not available")
			return
		}

		content, err := os.ReadFile(filePath)
		if err != nil {
			response.Error(c, http.StatusInternalServerError, "Failed to read export file")
			return
		}

		c.Data(http.StatusOK, "application/octet-stream", content)
	}
}

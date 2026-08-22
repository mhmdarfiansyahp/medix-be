package handler

import (
	"context"
	"fmt"
	"net/http"
	"time"

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
	h.router.GET("/export/excel", h.ExportExcel())      // US-17
}

// GET /reports/sales-summary?start_date=2026-01-01&end_date=2026-01-31&group_by=daily
func (h *ReportHandler) GetSalesSummary() gin.HandlerFunc {
	return func(c *gin.Context) {
		var params dto.ReportFilterParams
		c.ShouldBindQuery(&params)

		res, err := h.reportService.GetSalesSummary(c.Request.Context(), params)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"data": res})
	}
}

// GET /reports/drug-ranking?start_date=2026-01-01&end_date=2026-01-31
func (h *ReportHandler) GetDrugRanking() gin.HandlerFunc {
	return func(c *gin.Context) {
		var params dto.ReportFilterParams
		c.ShouldBindQuery(&params)

		res, err := h.reportService.GetDrugRanking(c.Request.Context(), params)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"data": res})
	}
}

// GET /reports/export/excel?start_date=2026-01-01&end_date=2026-01-31
func (h *ReportHandler) ExportExcel() gin.HandlerFunc {
	return func(c *gin.Context) {
		var params dto.ReportFilterParams
		c.ShouldBindQuery(&params)

		buf, err := h.reportService.ExportToExcel(c.Request.Context(), params)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		fileName := fmt.Sprintf("Laporan_Penjualan_%s.xlsx", time.Now().Format("20060102_150405"))
		c.Header("Content-Description", "File Transfer")
		c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", fileName))
		c.Data(http.StatusOK, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", buf.Bytes())
	}
}

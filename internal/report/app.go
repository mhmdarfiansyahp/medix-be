package report

import (
	"medix-be/internal/report/handler"
	"medix-be/internal/report/repository"
	"medix-be/internal/report/service"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type ReportHandler struct {
	DB     *gorm.DB
	Logger *logrus.Logger
	Router *gin.RouterGroup
}

func StartApp(cfg *ReportHandler) {
	cfg.Logger.Info("Report module starting")

	reportRepo := repository.NewReportRepository(cfg.DB)
	reportService := service.NewReportService(reportRepo)
	handlerContract := &handler.HandlerContract{
		Logger: cfg.Logger,
		Router: cfg.Router,
	}

	handler.StartReportHandler(handlerContract, &handler.ReportHandlerProps{
		ReportService: reportService,
	})
}

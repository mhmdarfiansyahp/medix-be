package drug

import (
	"medix-be/internal/drug/repository"
	"medix-be/internal/drug/service"
	"medix-be/internal/drug/handler"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type TypeDrugHandler struct {
	DB *gorm.DB
	Logger *logrus.Logger
	Router *gin.RouterGroup
}

func StartApp(cfg *TypeDrugHandler) {
	cfg.Logger.Info("TypeDrug module starting...")

	drugRepo := repository.NewTypeDrugRepository(cfg.DB)
	drugSvc := service.NewTypeDrugService(drugRepo)
	handlerContract := &handler.HandlerContract{
		Logger: cfg.Logger,
		Router: cfg.Router,
	}

	handler.StartTypeDrugHandler(handlerContract, &handler.TypeDrugHandlerProps{
		TypeDrugService: drugSvc,
	})
}
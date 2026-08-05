package transaction

import (
	"medix-be/internal/transaction/handler"
	"medix-be/internal/transaction/repository"
	"medix-be/internal/transaction/service"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type TransactionHandler struct {
	DB     *gorm.DB
	Logger *logrus.Logger
	Router *gin.RouterGroup
}

func StartApp(cfg *TransactionHandler) {
	cfg.Logger.Info("Transaction module starting...")

	transactionRepo := repository.NewTransactionRepository(cfg.DB)
	transactionSvc := service.NewTransactionService(transactionRepo)
	handlerContract := &handler.HandlerContract{
		Logger: cfg.Logger,
		Router: cfg.Router,
	}
	
	handler.StartTransactionHandler(handlerContract, &handler.TransactionHandlerProps{
		TransactionService: transactionSvc,
	})

}
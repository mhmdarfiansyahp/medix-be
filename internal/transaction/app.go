package transaction

import (
	"medix-be/internal/transaction/handler"
	"medix-be/internal/transaction/repository"
	"medix-be/internal/transaction/service"
	"medix-be/internal/transaction/service/core"
	"medix-be/internal/transaction/service/returns"
	"medix-be/internal/transaction/service/receipt"
	"medix-be/internal/transaction/service/report"

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
	
	// Create each service
	coreSvc := core.NewCoreService(transactionRepo)
	returnsSvc := returns.NewReturnService(transactionRepo)
	receiptSvc := receipt.NewReceiptService(transactionRepo)
	reportSvc := report.NewReportService(transactionRepo)
	
	// Create main orchestrator service
	transactionSvc := service.NewTransactionService(
		cfg.Logger,
		coreSvc,
		returnsSvc,
		receiptSvc,
		reportSvc,
	)
	
	handlerContract := &handler.HandlerContract{
		Logger: cfg.Logger,
		Router: cfg.Router,
	}
	
	handler.StartTransactionHandler(handlerContract, &handler.TransactionHandlerProps{
		TransactionService: transactionSvc,
	})

}
package medicine

import (
	"fmt"

	"medix-be/internal/medicine/handler"
	"medix-be/internal/medicine/repository"
	"medix-be/internal/medicine/service"
	transRepo "medix-be/internal/transaction/repository"
	transService "medix-be/internal/transaction/service"
	transCore "medix-be/internal/transaction/service/core"
	transReturn "medix-be/internal/transaction/service/returns"
	transReceipt "medix-be/internal/transaction/service/receipt"
	transReport "medix-be/internal/transaction/service/report"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type ModuleConfig struct {
	DB     *gorm.DB
	Logger *logrus.Logger
	Router *gin.RouterGroup
}

func StartApp(cfg *ModuleConfig) {
	cfg.Logger.Info("Medicine module starting...")

	medicineRepo := repository.NewMedicineRepository(cfg.DB)
	medicineSvc := service.NewMedicineService(medicineRepo)

	// Create transaction services (for cart operations)
	transactionRepo := transRepo.NewTransactionRepository(cfg.DB)
	coreSvc := transCore.NewCoreService(transactionRepo)
	returnsSvc := transReturn.NewReturnService(transactionRepo)
	receiptSvc := transReceipt.NewReceiptService(transactionRepo)
	reportSvc := transReport.NewReportService(transactionRepo)
	transactionSvc := transService.NewTransactionService(
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

	handler.StartMedicineHandler(handlerContract, &handler.MedicineHandlerProps{
		MedicineService:    medicineSvc,
		TransactionService: transactionSvc,
	})

	fmt.Println("Medicine module successfully initialized!")
}

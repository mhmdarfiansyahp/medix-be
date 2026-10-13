package service

import (
	"github.com/sirupsen/logrus"

	"medix-be/internal/transaction/model/dto"
	"medix-be/internal/transaction/service/core"
	"medix-be/internal/transaction/service/receipt"
	"medix-be/internal/transaction/service/report"
)

// Main orchestrator that coordinates all transaction services.
// This is the service that handlers will use.
// Internally delegates to specialized services for each concern.

type TransactionService struct {
	logger         *logrus.Logger
	coreService    core.CoreService
	receiptService receipt.ReceiptService
	reportService  report.ReportService
}

func NewTransactionService(
	logger *logrus.Logger,
	coreService core.CoreService,
	receiptService receipt.ReceiptService,
	reportService report.ReportService,
) *TransactionService {
	return &TransactionService{
		logger:         logger,
		coreService:    coreService,
		receiptService: receiptService,
		reportService:  reportService,
	}
}

// Handler methods that delegate to appropriate services.
// This maintains backward compatibility while enabling modularity.

func (s *TransactionService) CreateTransaction(userID uint, req dto.CreateTransactionRequest) (*dto.TransactionResponse, error) {
	return s.coreService.CreateTransaction(userID, req)
}

func (s *TransactionService) AddToCart(userID uint, req dto.AddToCartRequest) (*dto.TransactionResponse, error) {
	return s.coreService.AddToCart(userID, req)
}

func (s *TransactionService) GetAllTransactions() ([]dto.TransactionResponse, error) {
	return s.coreService.GetAllTransactions()
}

func (s *TransactionService) GetTransactionByID(id uint) (*dto.TransactionResponse, error) {
	return s.coreService.GetTransactionByID(id)
}

func (s *TransactionService) CancelTransaction(id uint, userID uint, userRole string) error {
	return s.coreService.CancelTransaction(id, userID, userRole)
}

func (s *TransactionService) ProcessPayment(id uint, userID uint, req dto.PaymentRequest) (*dto.PaymentResponse, error) {
	return s.coreService.ProcessPayment(id, userID, req)
}

func (s *TransactionService) GetTodayTransactions(userID uint) (*dto.TodayTransactionResponse, error) {
	return s.coreService.GetTodayTransactions(userID)
}

// Receipt operations
func (s *TransactionService) GetReceipt(id uint) (*dto.ReceiptResponse, error) {
	return s.receiptService.GetReceipt(id)
}

func (s *TransactionService) GenerateReceipt(id uint, format string) (*dto.GenerateReceiptResponse, error) {
	return s.receiptService.GenerateReceipt(id, format)
}

// Report operations
func (s *TransactionService) GetTodayKasirReport(userID uint) (*dto.TodayKasirReportResponse, error) {
	return s.reportService.GetTodayKasirReport(userID)
}

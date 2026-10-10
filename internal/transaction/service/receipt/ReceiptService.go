package receipt

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"medix-be/internal/transaction/model/entities"
	"medix-be/internal/transaction/model/dto"
	"medix-be/internal/transaction/repository"

	"github.com/jung-kurt/gofpdf"
)

type ReceiptService interface {
	GetReceipt(id uint) (*dto.ReceiptResponse, error)
	GenerateReceipt(id uint, format string) (*dto.GenerateReceiptResponse, error)
}

type receiptService struct {
	repo repository.TransactionRepository
}

func NewReceiptService(repo repository.TransactionRepository) ReceiptService {
	return &receiptService{repo: repo}
}

func (s *receiptService) GetReceipt(id uint) (*dto.ReceiptResponse, error) {
	transaction, err := s.repo.FindByID(id)
	if err != nil {
		return nil, errors.New("transaction not found")
	}

	if transaction.Status == entities.StatusTransaksiDibatalkan {
		return nil, errors.New("transaksi sudah dibatalkan")
	}

	details := make([]dto.DetailItemResponse, 0)
	for _, detail := range transaction.Details {
		details = append(details, dto.DetailItemResponse{
			IDDetail:    detail.IDDetail,
			IDObat:      detail.IDObat,
			Jumlah:      detail.Jumlah,
			HargaSatuan: detail.HargaSatuan,
			Subtotal:    detail.Subtotal,
		})
	}

	return &dto.ReceiptResponse{
		IDTransaksi:  transaction.IDTransaksi,
		TglTransaksi: transaction.TglTransaksi,
		IDUser:       transaction.IDUser,
		Details:      details,
		TotalHarga:   transaction.TotalHarga,
		MetodeBayar:  transaction.MetodeBayar,
		UangDiterima: transaction.UangDiterima,
		Kembalian:    transaction.Kembalian,
	}, nil
}

func (s *receiptService) GenerateReceipt(id uint, format string) (*dto.GenerateReceiptResponse, error) {
	transaction, err := s.repo.FindByID(id)
	if err != nil {
		return nil, errors.New("transaksi tidak ditemukan")
	}

	if transaction.Status == entities.StatusTransaksiDibatalkan {
		return nil, errors.New("transaksi sudah dibatalkan")
	}

	if format == "pdf" {
		return s.generatePDFReceipt(transaction)
	}

	return s.generateWhatsAppText(transaction)
}

func (s *receiptService) generatePDFReceipt(transaction *entities.Transaksi) (*dto.GenerateReceiptResponse, error) {
	names, err := s.obatNames(transaction)
	if err != nil {
		return nil, err
	}

	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Helvetica", "B", 14)
	pdf.Cell(0, 10, "MEDIX - STRUK PEMBELIAN")
	pdf.Ln(8)

	pdf.SetFont("Helvetica", "", 9)
	metode := transaction.MetodeBayar
	if metode == "" {
		metode = "-"
	}
	pdf.Cell(0, 6, fmt.Sprintf("No Transaksi : %d", transaction.IDTransaksi))
	pdf.Ln(5)
	pdf.Cell(0, 6, fmt.Sprintf("Tanggal      : %s", transaction.TglTransaksi.Format("2006-01-02 15:04:05")))
	pdf.Ln(5)
	pdf.Cell(0, 6, fmt.Sprintf("Metode Bayar : %s", metode))
	pdf.Ln(8)

	colW := []float64{10, 70, 20, 30, 30}
	headers := []string{"No", "Obat", "Jml", "Harga", "Subtotal"}

	pdf.SetFont("Helvetica", "B", 9)
	for i, h := range headers {
		pdf.CellFormat(colW[i], 7, h, "1", 0, "C", false, 0, "")
	}
	pdf.Ln(7)

	pdf.SetFont("Helvetica", "", 8)
	for i, detail := range transaction.Details {
		nama := names[detail.IDObat]
		if len(nama) > 35 {
			nama = nama[:32] + "..."
		}
		pdf.CellFormat(colW[0], 7, fmt.Sprintf("%d", i+1), "1", 0, "C", false, 0, "")
		pdf.CellFormat(colW[1], 7, nama, "1", 0, "L", false, 0, "")
		pdf.CellFormat(colW[2], 7, fmt.Sprintf("%d", detail.Jumlah), "1", 0, "C", false, 0, "")
		pdf.CellFormat(colW[3], 7, fmt.Sprintf("%.0f", detail.HargaSatuan), "1", 0, "R", false, 0, "")
		pdf.CellFormat(colW[4], 7, fmt.Sprintf("%.0f", detail.Subtotal), "1", 0, "R", false, 0, "")
		pdf.Ln(7)
	}

	pdf.SetFont("Helvetica", "B", 9)
	pdf.CellFormat(30, 7, "TOTAL", "1", 0, "R", false, 0, "")
	pdf.CellFormat(0, 7, fmt.Sprintf("%.0f", transaction.TotalHarga), "1", 0, "R", false, 0, "")
	pdf.Ln(8)

	if transaction.MetodeBayar == "tunai" {
		pdf.SetFont("Helvetica", "", 9)
		pdf.Cell(0, 6, fmt.Sprintf("Uang Diterima : Rp %.0f", transaction.UangDiterima))
		pdf.Ln(5)
		pdf.Cell(0, 6, fmt.Sprintf("Kembalian     : Rp %.0f", transaction.Kembalian))
		pdf.Ln(8)
	}

	pdf.SetFont("Helvetica", "I", 8)
	pdf.Cell(0, 6, "Terima kasih atas pembelian Anda.")

	os.MkdirAll("uploads/receipts", 0755)
	filePath := filepath.Join("uploads", "receipts", fmt.Sprintf("receipt_%d.pdf", transaction.IDTransaksi))

	var buf bytes.Buffer
	if err := pdf.Output(&buf); err != nil {
		return nil, errors.New("failed to generate receipt PDF")
	}

	if err := os.WriteFile(filePath, buf.Bytes(), 0644); err != nil {
		return nil, errors.New("failed to save receipt PDF")
	}

	return &dto.GenerateReceiptResponse{
		ReceiptID: transaction.IDTransaksi,
		Message:   filePath,
	}, nil
}

func (s *receiptService) generateWhatsAppText(transaction *entities.Transaksi) (*dto.GenerateReceiptResponse, error) {
	names, err := s.obatNames(transaction)
	if err != nil {
		return nil, err
	}

	var msg string
	msg += "*MEDIX - STRUK PEMBELIAN*\n"
	msg += fmt.Sprintf("No Transaksi : %d\n", transaction.IDTransaksi)
	msg += fmt.Sprintf("Tanggal      : %s\n", transaction.TglTransaksi.Format("2006-01-02 15:04:05"))
	msg += fmt.Sprintf("Metode Bayar : %s\n", transaction.MetodeBayar)
	msg += "------------------------\n"
	for _, detail := range transaction.Details {
		msg += fmt.Sprintf("%s\n%d x Rp %.0f = Rp %.0f\n", names[detail.IDObat], detail.Jumlah, detail.HargaSatuan, detail.Subtotal)
	}
	msg += "------------------------\n"
	msg += fmt.Sprintf("TOTAL        : Rp %.0f\n", transaction.TotalHarga)
	if transaction.MetodeBayar == "tunai" {
		msg += fmt.Sprintf("Uang Diterima: Rp %.0f\n", transaction.UangDiterima)
		msg += fmt.Sprintf("Kembalian    : Rp %.0f\n", transaction.Kembalian)
	}
	msg += "_Terima kasih atas pembelian Anda._"

	return &dto.GenerateReceiptResponse{
		ReceiptID: transaction.IDTransaksi,
		Message:   msg,
	}, nil
}

func (s *receiptService) obatNames(transaction *entities.Transaksi) (map[uint]string, error) {
	names := make(map[uint]string, len(transaction.Details))
	for _, detail := range transaction.Details {
		if _, ok := names[detail.IDObat]; ok {
			continue
		}
		obat, err := s.repo.FindObatByID(s.repo.GetDB(), detail.IDObat)
		if err != nil {
			return nil, errors.New("obat tidak ditemukan")
		}
		names[detail.IDObat] = obat.NamaObat
	}
	return names, nil
}
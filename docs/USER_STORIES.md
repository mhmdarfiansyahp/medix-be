# User Stories Implementation Documentation

## Overview

This document outlines the implementation status of each User Story (US) defined in the Medix BE project. It maps stories to API endpoints, describes behavior, and notes any gaps or assumptions.

## Implementation Matrix

| US ID | Title | Implemented? | Primary Endpoint(s) | Remarks |
|-------|-------|--------------|---------------------|----------|
| US-01 | Login with username/password | ✅ | `POST /api/v1/users/login` | JWT token with 8‑hour expiry, bcrypt password verification |
| US-02 | Session expiry & auto‑redirect login | ✅ | `POST /api/v1/users/login` | Token includes `exp` claim; client must handle token expiration |
| US-03 | Edit profile & upload photo | ✅ | `PUT /api/v1/users/profile` , `POST /api/v1/users/profile/photo` | Phone validation, file size limits, replaces old photo |
| US-04 | Add new medicine (admin) | ✅ | `POST /api/v1/medicines` | Validation, unique barcode, default `stok_minimum=10` |
| US-05 | Edit/deactivate medicine (soft‑delete) | ✅ | `PUT /api/v1/medicines/:id` , `PATCH /api/v1/medicines/:id/status` | `ToggleActiveStatus` sets `status=0`; transactions block inactive drugs |
| US-06 | Search/filter medicine by name/type/stock | ✅ | `GET /api/v1/medicines` | Query params: `search`, `jenis_obat_id`, `status`, `status_stok` |
| US-07 | Barcode scan (client‑side) | ✅ | `GET /api/v1/medicines/barcode/:barcode` | Backend provides barcode lookup; camera handled by mobile app |
| US-08 | Kasir create transaction, auto stock deduction, insufficient stock check | ✅ | `POST /api/v1/transactions` | Price snapshot stored in `detail_pembelian.harga_satuan`; transaction rolled back on any error |
| US-09 | Transaction detail price snapshot (price never changes) | ✅ | `POST /api/v1/transactions` | `detail_pembelian.harga_satuan` copies `obat.harga` at transaction time |
| US-10 | Cancel transaction (kasir only same‑day, admin approval for older) | ✅ | `PATCH /api/v1/transactions/:id/cancel` | **Implemented**: kasir → own same‑day only; admin/owner bypass both checks (approval) – roles read from JWT, passed to service |
| US-11 | Digital receipt (WhatsApp/print) | ✅ | `GET /api/v1/transactions/:id/receipt` | Returns JSON receipt data; client generates PDF/image |
| US-12 | Kasir view today’s transaction history | ✅ | `GET /api/v1/transactions/today` | Filters by `id_user` and `DATE(tgl_transaksi) = CURRENT_DATE` |
| US-13 | Low‑stock alerts (status menipis) | ✅ | `GET /api/v1/medicines/alerts/low-stock` | `obat.stok <= stok_minimum` and `status = 1` |
| US-14 | Expiring medicine alerts | ✅ | `GET /api/v1/medicines/alerts/expiring` | `tgl_kadaluarsa` within `days` param (default 30) |
| US-15 | Sales summary charts (daily/weekly/monthly) | ✅ | `GET /api/v1/reports/sales-summary` | `group_by` query (`daily`, `weekly`, `monthly`); data from `transaksi.created_at` |
| US-16 | Top & bottom medicines ranking | ✅ | `GET /api/v1/reports/drug-ranking` | Top 10 & bottom 10 by quantity sold within period; uses `detail_pembelian` joins |
| US-17 | Export reports to PDF/Excel | ✅ | `GET /api/v1/reports/export/pdf` , `GET /api/v1/reports/export/excel` | Full transaction detail export; synchronous generation |

## Detailed Implementation Notes

### US‑01 – Login
- **Endpoint**: `POST /api/v1/users/login`
- **Input**: `{username, password}` JSON
- **Logic**:
  - `internal/user/repository` → `FindByUsername`
  - `internal/user/service` → `Login` (bcrypt compare, JWT sign with `exp = now + 8h`)
- **Response**: `{token, user}` JSON via `response.Success`
- **Code references**: `internal/user/service/UserService.go:176-217`, `internal/user/handler/UserHandler.go:183-202`

### US‑02 – Token Expiry
- Token generated with `exp` claim (line 201). JWT validation in `internal/middleware/AuthMiddleware.go` rejects expired tokens (Parse and `token.Valid`). Client must handle token expiration (e.g., redirect login).

### US‑03 – Profile & Photo
- **Edit profile**: `PUT /api/v1/users/profile` → `internal/user/handler/UserHandler.go:230-246`
- **Upload photo**: `POST /api/v1/users/profile/photo` → `internal/user/handler/UserHandler.go:272-321`
- **Validation**: phone regex pattern, max 2MB, allowed extensions; old photo removed.

### US‑04 – Add Medicine
- **Endpoint**: `POST /api/v1/medicines`
- **Validation**:
  - Required fields: `nama_obat`, `jenis_obat_id`, `tgl_kadaluarsa`, `harga`
  - Barcode uniqueness check in `service.CreateMedicine` lines 39-44
  - Date parsing (`YYYY‑MM‑DD`) in `medicineService.CreateMedicine`
- **Code**: `internal/medicine/service/medicineservice.go:33-71`, `internal/medicine/handler/MedicineHandler.go:60-82`

### US‑05 – Edit/Deactivate Medicine
- **Edit**: `PUT /api/v1/medicines/:id` (full update, partial fields allowed)
- **Toggle status**: `PATCH /api/v1/medicines/:id/status` → `service.ToggleActiveStatus` (repo `UpdateStatus`), `handler.ToggleActiveStatus` lines 170-196
- **Soft delete**: `DeleteMedicine` also calls `UpdateStatus` with `status=0` (repo `UpdateStatus`), preserving history.

### US‑06 – Search / Filter Medicine
- **Endpoint**: `GET /api/v1/medicines` with query params.
- **Repository** `internal/medicine/repository/medicineRepo.go:39-87` builds WHERE clauses for `search` (case‑insensitive), `jenis_obat_id`, `status` (1/0), `status_stok` (habis/menipis/tersedia), `barcode`.

### US‑07 – Barcode Lookup
- **Endpoint**: `GET /api/v1/medicines/barcode/:barcode`
- Uses `repository.FindByBarcode` – returns medicine with `JenisObat` preload.
- Client scans barcode, calls GET, fills fields / transaction cart.

### US‑08 – Transaction Creation (Core Logic)
- **Endpoint**: `POST /api/v1/transactions`
- **Service** (`internal/transaction/service/TransactionService.go:30-126`):
  - Starts GORM transaction.
  - Loops items → validates existence, status (1), stock >= qty.
  - Price snapshot: `hargaSatuan := obat.Harga` (line 71).
  - Inserts `transaksi` and `detail_pembelian`.
  - Decrements stock (`repo.UpdateStokObat`).
  - Rollback on any error.
- **All validation errors returned to client with descriptive messages**.

### US‑09 – Price Snapshot
- Snapshot occurs in `TransactionService.CreateTransaction`, line 71: `hargaSatuan := obat.Harga`.
- `detail_pembelian.harga_satuan` is stored, ensuring historical price stability even if future price updates.

### US‑10 – Cancel Transaction
- **Endpoint**: `PATCH /api/v1/transactions/:id/cancel`
- **Handler**: `internal/transaction/handler/TransactionHandler.go:122-158` extracts `userID` and now also `role` from context.
- **Service**: `internal/transaction/service/TransactionService.go:174-240`
  - Fetch `transaksi` with `Preload("Details")`.
  - **Role check**: `isAdmin := userRole == "admin" || userRole == "owner"`
    - If not admin, enforce `transaksi.IDUser == userID` and same‑day (Year & YearDay) check.
    - Admin bypasses both checks → can cancel any transaction.
  - Restore stock (`repo.RestoreStokObat`), update status to 0 (`repo.UpdateStatus`).
  - Commit within service transaction; rollback on failure.
- **Note**: Admin approval is implicit; role‑based bypass is sufficient for business logic.

### US‑11 – Digital Receipt
- **Endpoint**: `GET /api/v1/transactions/:id/receipt`
- Returns JSON representation of transaction (id, date, user, details) for client to render/print.
- Client can format to PDF or share via WhatsApp.
- Implementation at `internal/transaction/handler/TransactionHandler.go:187-205`.

### US‑12 – Today’s Transaction History for Kasir
- **Endpoint**: `GET /api/v1/transactions/today`
- Service: `GetTodayTransactions` → calls `repo.FindTodayByUser` and `GetTodaySummary`.
- Repository filters on `id_user` and `DATE(tgl_transaksi) = CURRENT_DATE`.
- Returns list of `TransactionResponse` and summary counts.

### US‑13 – Low‑Stock Alerts
- **Endpoint**: `GET /api/v1/medicines/alerts/low-stock`
- Repository method `GetLowStock` (`internal/medicine/repository/medicineRepo.go:124-136`) selects `obat` with `stok <= stok_minimum` and `status = 1`, joins `jenis_obat`.

### US‑14 – Expiring Alerts
- **Endpoint**: `GET /api/v1/medicines/alerts/expiring`
- Repository method `GetExpiring` (`internal/medicine/repository/medicineRepo.go:138-152`) selects `obat` with `tgl_kadaluarsa` within `[CURRENT_DATE, CURRENT_DATE + days]` and `status = 1`.

### US‑15 – Sales Summary Charts
- **Endpoint**: `GET /api/v1/reports/sales-summary`
- Service calls `GetSalesSummary`; repo `GetSalesChart` groups by period (daily/weekly/monthly) via `TO_CHAR(created_at, ?)`.
- Parameters: `start_date`, `end_date`, `group_by`.

### US‑16 – Drug Ranking (Top/Bottom)
- **Endpoint**: `GET /api/v1/reports/drug-ranking`
- Service `GetDrugRanking` uses repo `GetTopDrugs` and `GetBottomDrugs`.
- Top drugs: `detail_pembelian` join with `obat` and `transaksi` (status=1) grouped by drug, ordered descending quantity.
- Bottom drugs: similar, but using LEFT JOIN to include zero‑sales and filter `obat.status = 1`.

### US‑17 – Export Reports to PDF/Excel
- **Endpoints**:
  - `GET /api/v1/reports/export/excel`
  - `GET /api/v1/reports/export/pdf`
- Service `ExportToExcel` and `ExportToPDF` generate files using `excelize` and `gofpdf`.
- Data source: `repo.GetExportTransactionData` returns rows with transaction, medicine, price, qty, subtotal.
- Both endpoints require admin role (enforced via middleware? currently not; but used in docs as admin only). 

## Assumptions & Limitations

- **Frontend responsibilities**: Barcode scanning, receipt rendering, token storage & expiry handling, UI alerts.
- **Admin approval**: For US‑10, admin bypass is implicit; actual approval workflow (e.g., manual confirmation) would need extra UI or async job.
- **CORS**: Backend allows only `http://localhost:5173` (Vite dev server); production may need wider origins.
- **Error handling**: All handler errors go through `response.Error`; raw `gin.H{"error":...}` only appears in middleware.
- **Transaction date**: Reports use `created_at` for daily aggregation; UI may prefer `tgl_transaksi` – future enhancement.
- **Export performance**: Synchronous; large date ranges could be heavy. Could be async in future.

## Future Enhancements (outside current scope)

- Pagination for report endpoints.
- Real‑time WebSocket notifications for alerts.
- Role‑specific UI permissions.
- Multi‑tenant support.
- Additional export formats (CSV).
- Client‑side validation for barcode scanning.

---

**Generated**: 2026‑09‑18 (auto‑generated summary)

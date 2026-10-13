# Graph Report - medix-be  (2026-10-13)

## Corpus Check
- Corpus is ~27,492 words - fits in a single context window. You may not need a graph.

## Summary
- 571 nodes · 1327 edges · 67 communities (23 shown, 44 thin omitted)
- Extraction: 97% EXTRACTED · 3% INFERRED · 0% AMBIGUOUS · INFERRED: 37 edges (avg confidence: 0.9)
- Token cost: 0 input · 0 output

## Community Hubs (Navigation)
- Community 0: ScanMedicineRequest, ScanMedicineResponse
- Community 1: github.com/gin-gonic/gin.Context, github.com/gin-gonic/gin.HandlerFunc
- Community 2: main.go, main()
- Community 3: TypeDrugHandler, TypeDrug
- Community 4: .TableName(), .TableName()
- Community 5: .AddToCart(), .CreateTransaction()
- Community 6: CreateMedicineRequest Schema, CreateTransactionRequest Schema
- Community 7: ReturnItemRequest, ReturnItemResponse
- Community 8: Session, .TableName()
- Community 9: UserHandlerProps, UpdateUserRequest
- Community 10: DetailItemResponse, TransactionHandlerProps
- Community 11: PaginationResponse, .ValidateNoTelp()
- Community 12: go_pkg_crypto_rand, go_pkg_encoding_hex
- Community 13: LoginRequest, LoginResponse
- Community 14: .CancelTransaction(), .GetAllTransactions()
- Community 15: .GetTodayTransactions(), ReturnSummaryResponse
- Community 16: US-01 Login Kasir Mobile, Bcrypt Password Verification
- Community 17: model.User, UserRepository
- Community 18: Application Configuration, PostgreSQL Docker Service
- Community 19: US-07 Scan Barcode, US-08 Input Transaksi Mobile
- Community 20: go_pkg_encoding_json, go_pkg_testing
- Community 21: US-05 Edit & Nonaktifkan Obat, US-18 Kelola Akun User
- Community 22: US-04 Tambah Obat, Medicine Domain Module
- Community 23: US-15b Dashboard Ringkas Mobile, US-15 Grafik Penjualan
- Community 24: BearerAuth Security Scheme, JWT Authentication
- Community 25: US-02 Sesi Login Kasir Mobile, US-02 Sesi Login Admin/Owner Web
- Community 26: US-06 Cari Obat Kasir Mobile, US-06 Cari & Filter Obat Web
- Community 27: US-09 Harga Snapshot Mobile, US-09w Harga Snapshot Web
- Community 28: US-10 Retur Obat, US-20 Approval & Lihat Detail Retur
- Community 29: US-11 Struk Digital, US-11w Cetak/Kirim Struk Web
- Community 30: US-12 Riwayat Transaksi Kasir, US-12w Riwayat Transaksi Web
- Community 31: US-13b Push Notifikasi Stok Menipis, US-13 Alert Stok Menipis Web
- Community 33: Drug Type Domain Module, API Test: Create Jenis Obat
- Community 34: Medicines API Endpoints, US-13 Low Stock Alerts
- Community 35: Price Snapshot Design, US-09 Transaction Price Snapshot
- Community 36: Reports API Endpoints, US-17 Export Reports PDF/Excel
- Community 37: Transactions API Endpoints, US-08 Transaction Creation with Stock Deduction
- Community 38: Mobile Kasir Autentikasi & Profil
- Community 39: US-03 Edit Profil Kasir
- Community 40: Web Admin Autentikasi
- Community 41: Mobile Kasir Lihat & Cari Obat
- Community 42: US-06b Input Massal Bulk
- Community 43: Web Admin Manajemen Obat
- Community 44: Mobile Kasir Transaksi
- Community 45: US-19 Reset Password User
- Community 46: Web Admin Manajemen User
- Community 47: Mobile Notifikasi & Laporan Ringkas
- Community 48: US-14 Alert Mendekati Kadaluarsa
- Community 49: Web Admin Notifikasi & Stok
- Community 50: US-16 Obat Terlaris/Terendah
- Community 51: US-17 Ekspor Laporan
- Community 52: Web Admin Laporan & Dashboard
- Community 53: Web Admin Approval Retur
- Community 54: US-10w Batalkan Transaksi Web
- Community 55: Web Admin Transaksi Fallback
- Community 57: AGENTS.md Project Configuration
- Community 58: Medix BE Documentation
- Community 59: Registered API Endpoints
- Community 60: Jenis Obat API Collection Folder
- Community 61: API Test: Get All Jenis Obat
- Community 62: API Test: Update Jenis Obat
- Community 63: API Test: Get All Obat
- Community 64: medix-be
- Community 66: User Stories Implementation Matrix

## God Nodes (most connected - your core abstractions)
1. `Success()` - 57 edges
2. `TransactionRepository` - 34 edges
3. `TransactionHandler` - 27 edges
4. `TransactionService` - 26 edges
5. `UserService` - 22 edges
6. `TransactionResponse` - 21 edges
7. `MedicineHandler` - 20 edges
8. `UserHandler` - 20 edges
9. `MedicineRepository` - 16 edges
10. `MedicineService` - 15 edges

## Surprising Connections (you probably didn't know these)
- `Get Medicines API Request` --semantically_similar_to--> `CreateMedicineRequest Schema`  [INFERRED] [semantically similar]
  Medix/Obat/Get data.yml → docs/openapi.yaml
- `Create Transaction API Request` --semantically_similar_to--> `CreateTransactionRequest Schema`  [INFERRED] [semantically similar]
  Medix/Transaction/Create data.yml → docs/openapi.yaml
- `Create User API Request` --semantically_similar_to--> `LoginRequest Schema`  [INFERRED] [semantically similar]
  Medix/User/Create data.yml → docs/openapi.yaml
- `main()` --calls--> `ConnectDatabase()`  [EXTRACTED]
  cmd/api/main.go → config/config.go
- `main()` --calls--> `LoadEnv()`  [EXTRACTED]
  cmd/api/main.go → config/config.go

## Import Cycles
- None detected.

## Hyperedges (group relationships)
- **Authentication Flow (Mobile + Web)** — 01_auth__1__us01_login_kasir, 01_auth_us01_login_admin, 01_auth_jwt_token, 01_auth_bcrypt_hash, agents_md_auth_middleware, documentation_md_auth_flow [EXTRACTED 0.95]
- **Transaction Flow (Mobile + Web Fallback)** — 03_transaksi_us08_input_transaksi, 03_transaksi_us08b_proses_pembayaran, 03_transaksi_us09_harga_snapshot, 06_transaksi_us08w_input_transaksi_web, 06_transaksi_us09w_harga_snapshot_web, domain_transaction [EXTRACTED 0.95]
- **Retur & Approval Flow** — 03_transaksi_us10_retur_obat, 06_approval_transaksi_us20_approval_retur, 06_transaksi_us10w_batalkan_transaksi [EXTRACTED 0.95]
- **Bruno API Collection Structure** — medix_opencollection_medix_api_collection, medix_obat_folder_obat_collection, medix_transaction_folder_transaction_collection, medix_user_folder_user_collection [EXTRACTED 1.00]
- **Transaction-Medicine Stock Deduction Flow** — user_stories_us08_transaction_creation, user_stories_us09_price_snapshot, user_stories_us10_cancel_transaction, readme_price_snapshot_rationale [INFERRED 0.95]
- **Infrastructure Configuration Stack** — config_config_application_config, docker_compose_postgres_service, readme_medix_be_overview [INFERRED 0.85]

## Communities (67 total, 44 thin omitted)

### Community 0 - "Community 0: ScanMedicineRequest, ScanMedicineResponse"
Cohesion: 0.05
Nodes (31): ScanMedicineRequest, ScanMedicineResponse, Obat, CreateMedicineRequest, MedicineFilterParams, MedicineResponse, UpdateMedicineRequest, ExpiringDrugResponse (+23 more)

### Community 1 - "Community 1: github.com/gin-gonic/gin.Context, github.com/gin-gonic/gin.HandlerFunc"
Cohesion: 0.10
Nodes (11): MedicineHandler, MedicineHandlerProps, TransactionHandler, TypeDrugHandler, UserHandler, Error(), Success(), StartMedicineHandler() (+3 more)

### Community 2 - "Community 2: main.go, main()"
Cohesion: 0.10
Nodes (11): main(), ConnectDatabase(), LoadEnv(), ApprovalThresholdRequest, CancelTransactionResponse, RejectReturnRequest, Pengaturan, ReportHandlerProps (+3 more)

### Community 3 - "Community 3: TypeDrugHandler, TypeDrug"
Cohesion: 0.07
Nodes (24): TypeDrugHandler, TypeDrug, ReportHandler, TypeDrugHandlerProps, StartApp(), HandlerContract, StartTypeDrugHandler(), CreateTypeDrugRequest (+16 more)

### Community 4 - "Community 4: .TableName(), .TableName()"
Cohesion: 0.08
Nodes (16): StartApp(), StartApp(), GenerateReceiptResponse, Return, ReturnItem, DetailPembelian, Transaksi, TransactionRepository (+8 more)

### Community 5 - "Community 5: .AddToCart(), .CreateTransaction()"
Cohesion: 0.18
Nodes (8): transactionCoreService, AddToCartResponse, DetailItemRequest, TransactionService, AddToCartRequest, CreateTransactionRequest, TransactionResponse, toTransactionResponse()

### Community 6 - "Community 6: CreateMedicineRequest Schema, CreateTransactionRequest Schema"
Cohesion: 0.11
Nodes (19): CreateMedicineRequest Schema, CreateTransactionRequest Schema, LoginRequest Schema, Obat API Collection Folder, Get Medicines API Request, Update Medicine API Request, Update Medicine Status API Request, Medix Bruno API Collection (+11 more)

### Community 7 - "Community 7: ReturnItemRequest, ReturnItemResponse"
Cohesion: 0.17
Nodes (7): ReturnItemRequest, ReturnItemResponse, CreateReturnRequest, CreateReturnResponse, RejectReturnResponse, getThreshold(), ReturnService

### Community 8 - "Community 8: Session, .TableName()"
Cohesion: 0.16
Nodes (6): Session, StartApp(), SessionRepository, NewSessionRepository(), NewUserRepository(), NewUserService()

### Community 9 - "Community 9: UserHandlerProps, UpdateUserRequest"
Cohesion: 0.23
Nodes (7): UserHandlerProps, UpdateUserRequest, UserResponse, UserService, normalizeStatusFilter(), statusToInt(), toUserResponse()

### Community 10 - "Community 10: DetailItemResponse, TransactionHandlerProps"
Cohesion: 0.13
Nodes (6): DetailItemResponse, TransactionHandlerProps, ApproveReturnResponse, ApprovalThresholdResponse, ReceiptResponse, TransactionService

### Community 11 - "Community 11: PaginationResponse, .ValidateNoTelp()"
Cohesion: 0.16
Nodes (8): PaginationResponse, UserFilterRequest, UserStatus, ChangePasswordRequest, CreateUserRequest, LogoutRequest, UpdateProfileRequest, UserListResponse

### Community 12 - "Community 12: go_pkg_crypto_rand, go_pkg_encoding_hex"
Cohesion: 0.20
Nodes (3): ResetPasswordResponse, generateTempPassword(), intToStatus()

### Community 13 - "Community 13: LoginRequest, LoginResponse"
Cohesion: 0.22
Nodes (6): LoginRequest, LoginResponse, RefreshRequest, RefreshResponse, User, generateRefreshToken()

### Community 14 - "Community 14: .CancelTransaction(), .GetAllTransactions()"
Cohesion: 0.28
Nodes (3): PaymentRequest, PaymentResponse, CoreService

### Community 15 - "Community 15: .GetTodayTransactions(), ReturnSummaryResponse"
Cohesion: 0.25
Nodes (5): ReturnSummaryResponse, TransactionSummaryResponse, TodayKasirReportResponse, TodayTransactionResponse, toTransactionResponse()

### Community 16 - "Community 16: US-01 Login Kasir Mobile, Bcrypt Password Verification"
Cohesion: 0.32
Nodes (6): US-01 Login Kasir Mobile, Bcrypt Password Verification, JWT Token Authentication, US-01 Login Admin/Owner Web, Role-Based Access Control, API Test: Login Request

### Community 18 - "Community 18: Application Configuration, PostgreSQL Docker Service"
Cohesion: 0.33
Nodes (7): Application Configuration, PostgreSQL Docker Service, Diagram Viewing Instructions, OpenAPI 3.0 Specification, SuccessResponse Envelope Schema, Documentation Index, Medix BE README

### Community 19 - "Community 19: US-07 Scan Barcode, US-08 Input Transaksi Mobile"
Cohesion: 0.33
Nodes (6): US-07 Scan Barcode, US-08 Input Transaksi Mobile, US-08b Proses Pembayaran, US-08w Input Transaksi Web, Transaction Domain Module, API Test: Get Obat by Barcode

### Community 21 - "Community 21: US-05 Edit & Nonaktifkan Obat, US-18 Kelola Akun User"
Cohesion: 0.50
Nodes (3): US-05 Edit & Nonaktifkan Obat, US-18 Kelola Akun User, User Domain Module

### Community 22 - "Community 22: US-04 Tambah Obat, Medicine Domain Module"
Cohesion: 0.67
Nodes (3): US-04 Tambah Obat, Medicine Domain Module, API Test: Create Obat

### Community 23 - "Community 23: US-15b Dashboard Ringkas Mobile, US-15 Grafik Penjualan"
Cohesion: 0.67
Nodes (3): US-15b Dashboard Ringkas Mobile, US-15 Grafik Penjualan, Report Domain Module

### Community 24 - "Community 24: BearerAuth Security Scheme, JWT Authentication"
Cohesion: 0.67
Nodes (3): BearerAuth Security Scheme, JWT Authentication, US-01 Login with Username/Password

## Knowledge Gaps
- **77 isolated node(s):** `medix-be`, `Mobile Kasir Autentikasi & Profil`, `US-02 Sesi Login Kasir Mobile`, `US-03 Edit Profil Kasir`, `Web Admin Autentikasi` (+72 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 137 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **44 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `TransactionService` connect `Community 10: DetailItemResponse, TransactionHandlerProps` to `Community 1: github.com/gin-gonic/gin.Context, github.com/gin-gonic/gin.HandlerFunc`, `Community 2: main.go, main()`, `Community 3: TypeDrugHandler, TypeDrug`, `Community 4: .TableName(), .TableName()`, `Community 5: .AddToCart(), .CreateTransaction()`, `Community 7: ReturnItemRequest, ReturnItemResponse`, `Community 14: .CancelTransaction(), .GetAllTransactions()`, `Community 15: .GetTodayTransactions(), ReturnSummaryResponse`?**
  _High betweenness centrality (0.057) - this node is a cross-community bridge._
- **What connects `medix-be`, `Mobile Kasir Autentikasi & Profil`, `US-02 Sesi Login Kasir Mobile` to the rest of the system?**
  _77 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `Community 0: ScanMedicineRequest, ScanMedicineResponse` be split into smaller, more focused modules?**
  _Cohesion score 0.054340396445659606 - nodes in this community are weakly interconnected._
- **Why does `UserHandler` connect `Community 1: github.com/gin-gonic/gin.Context, github.com/gin-gonic/gin.HandlerFunc` to `Community 0: ScanMedicineRequest, ScanMedicineResponse`, `Community 9: UserHandlerProps, UpdateUserRequest`, `Community 2: main.go, main()`, `Community 3: TypeDrugHandler, TypeDrug`?**
  _High betweenness centrality (0.053) - this node is a cross-community bridge._
- **Should `Community 1: github.com/gin-gonic/gin.Context, github.com/gin-gonic/gin.HandlerFunc` be split into smaller, more focused modules?**
  _Cohesion score 0.09657594381035997 - nodes in this community are weakly interconnected._
- **Why does `UserService` connect `Community 9: UserHandlerProps, UpdateUserRequest` to `Community 1: github.com/gin-gonic/gin.Context, github.com/gin-gonic/gin.HandlerFunc`, `Community 8: Session, .TableName()`, `Community 11: PaginationResponse, .ValidateNoTelp()`, `Community 12: go_pkg_crypto_rand, go_pkg_encoding_hex`, `Community 13: LoginRequest, LoginResponse`, `Community 17: model.User, UserRepository`?**
  _High betweenness centrality (0.051) - this node is a cross-community bridge._
- **Should `Community 2: main.go, main()` be split into smaller, more focused modules?**
  _Cohesion score 0.09557109557109557 - nodes in this community are weakly interconnected._
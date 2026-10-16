# Product Requirements Document (PRD) — Medix Backend

## 1. Tujuan & Latar Belakang
Medix Backend (medix-be) adalah sistem RESTful API untuk manajemen apotek atau klinik (data obat, transaksi, pengguna, dan laporan). Dibangun dengan Go, Gin, GORM, dan PostgreSQL.

## 2. Fitur Utama
1. **Autentikasi & Manajemen Pengguna (`user`)**:
   - Login JWT (`POST /api/v1/users/login`).
   - Manajemen akun pengguna (Admin, Kasir).
2. **Manajemen Obat & Golongan (`medicine`, `drug`)**:
   - CRUD data obat (nama, barcode, stok, harga, status).
   - CRUD golongan/kategori obat (`type-drugs`).
3. **Transaksi Penjualan (`transaction`)**:
   - Pencatatan transaksi penjualan obat.
   - Perhitungan total dan pengelolaan stok otomatis.
4. **Laporan (`report`)**:
   - Rekapitulasi transaksi dan data operasional.

## 3. Spesifikasi Teknis
- **Bahasa**: Go 1.26.5
- **Framework**: Gin
- **ORM**: GORM + PostgreSQL
- **Migrasi**: golang-migrate (`migrations/data`)
- **Konfigurasi**: `config/config.yaml` (Viper)

## 4. Keamanan & Konfigurasi
- File sensitif `config/config.yaml` diabaikan oleh `.gitignore`.
- Menggunakan template `config/config.yaml.example` untuk deployment/setup awal.
- Dockerized via `docker-compose.yml` untuk kemudahan local deployment.

## 5. Deployment & Jalankan
- Lokal: `go run cmd/api/main.go`
- Docker: `docker compose up --build -d`

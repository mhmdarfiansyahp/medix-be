# medix-be

[![Go](https://img.shields.io/badge/Go-1.26.5-blue)]()
[![Gin](https://img.shields.io/badge/Gin-%7E1.12-green)]()
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-16-alpine)]()
[![License](https://img.shields.io/badge/License-MIT-blue)]()

REST API untuk sistem layanan medis **Medix**.

---

## ⚡ Quick Start

```bash
# Jalankan PostgreSQL via Docker
docker compose up -d

# Jalankan server (migrasi + seed otomatis)
go run cmd/api/main.go
```

Server berjalan di **http://localhost:8080**.

---

## 📦 Prasyarat

- Go 1.26+
- PostgreSQL 16+ (atau Docker)
- Docker (opsional)

---

## 🛠️ Konfigurasi

Konfigurasi utama di `config/config.yaml`. Salin `.env.example` ke `.env` jika perlu override:

```bash
cp .env.example .env
```

| Key | Default |
|-----|---------|
| `app.port` | `8080` |
| `database.name` | `medix` |
| `jwt.secret` | `change-this-secret-in-production` |
| `migration_path` | `file://migrations/data` |

---

## 📚 User Stories

### Autentikasi & Profil

| US | Deskripsi | Endpoint |
|----|-----------|----------|
| US-01 | Login (JWT 8 jam, bcrypt) | `POST /api/v1/users/login` |
| US-03 | Update profil + foto | `PUT /api/v1/users/profile`, `POST /api/v1/users/profile/photo` |

### Manajemen Obat

| US | Deskripsi | Endpoint |
|----|-----------|----------|
| US-04 | Tambah obat | `POST /api/v1/medicines` |
| US-05 | Edit / nonaktifkan obat | `PUT /api/v1/medicines/:id`, `PATCH /api/v1/medicines/:id/status` |
| US-06 | Search & filter | `GET /api/v1/medicines?search=&jenis_obat_id=&status=` |
| US-07 | Cari by barcode | `GET /api/v1/medicines/barcode/:barcode` |
| US-13 | Stok menipis | `GET /api/v1/medicines/alerts/low-stock` |
| US-14 | Dekat kadaluarsa | `GET /api/v1/medicines/alerts/expiring?days=30` |

### Transaksi / Kasir

| US | Deskripsi | Endpoint |
|----|-----------|----------|
| US-08 | Buat transaksi multi item | `POST /api/v1/transactions` |
| US-10 | Batalkan transaksi | `PATCH /api/v1/transactions/:id/cancel` |
| US-11 | Struk digital | `GET /api/v1/transactions/:id/receipt` |
| US-12 | Riwayat hari ini | `GET /api/v1/transactions/today` |

### Laporan & Dashboard

| US | Deskripsi | Endpoint |
|----|-----------|----------|
| US-15 | Ringkasan penjualan | `GET /api/v1/reports/sales-summary?group_by=daily\|weekly\|monthly` |
| US-16 | Obat terlaris | `GET /api/v1/reports/drug-ranking` |
| US-17 | Export Excel / PDF | `GET /api/v1/reports/export/excel`, `GET /api/v1/reports/export/pdf` |

Semua endpoint kecuali login dan create user memerlukan:

```
Authorization: Bearer <token>
```

---

## 🔑 Akun Dummy (Seed)

| Username | Password | Role |
|----------|----------|------|
| `admin` | `admin123` | admin |
| `kasir1` | `kasir123` | kasir |
| `kasir2` | `kasir123` | kasir (nonaktif) |
| `kasir3` | `kasir123` | kasir |
| `owner` | `owner123` | owner |

---

## 🏗️ Struktur Direktori

```
cmd/api/       → Entry point
config/        → Konfigurasi YAML + env
internal/      → Business logic (handler, service, repository)
migrations/    → Skrip migrasi + seed
Medix/         → Koleksi Postman
```

---

## 🧪 Build

```bash
go build ./...
```

---

## 📝 Catatan

- Migrasi **forward-only** (tidak ada `.down.sql`).
- CORS hardcoded ke `http://localhost:5173`.
- Status user disimpan sebagai `SMALLINT` (1 = aktif, 0 = nonaktif).

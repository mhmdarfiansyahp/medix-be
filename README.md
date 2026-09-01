# medix-be

REST API untuk sistem layanan medis Medix yang dibangun menggunakan **Go (Golang)**.

---

## 🛠️ Tech Stack & Dependencies

* **Language:** Go (Golang)
* **Framework:** Gin **
* **Database:** PostgreSQL **
* **Environment Configuration:** `.env`

---

## 📁 Struktur Direktori

* `cmd/api/` — Entry point utama aplikasi (`main.go`).
* `config/` — Konfigurasi database dan environment.
* `internal/` — Core business logic, handlers, services, dan repository.
* `migrations/` — File skrip migrasi database.
* `Medix/` — Modul / paket tambahan proyek.

---

## 🚀 Cara Menjalankan Aplikasi

### 1. Prasyarat
Pastikan kamu sudah menginstal:
* [Go](https://golang.org/doc/install) (versi 1.18+)
* Database pilihan (PostgreSQL / MySQL)

### 2. Clone Repository
```bash
git clone [https://github.com/mhmdarfiansyahp/medix-be.git](https://github.com/mhmdarfiansyahp/medix-be.git)
cd medix-be

# Run the application
make run

# Build binary
make build

# Run database migrations
make migrate-up

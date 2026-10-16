# Arsitektur Sistem Medix Backend (medix-be)

Dokumen ini mendeskripsikan arsitektur backend aplikasi Medix berdasarkan pemetaan knowledge graph sistem (menggunakan `graphify`), mencakup struktur modular, aliran data, lapisan layanan, serta pengelolaan basis data.

---

## 1. Ikhtisar Arsitektur

Backend Medix dibangun menggunakan bahasa **Go (v1.26.5)** dengan kerangka kerja **Gin** untuk routing HTTP dan **GORM** sebagai ORM untuk berinteraksi dengan basis data **PostgreSQL**.

### Prinsip Desain Modular
Setiap domain bisnis diisolasi ke dalam modul mandiri di bawah direktori `internal/`. Setiap modul mengimplementasikan pola standar:
1. **Model / DTO**: Struktur data request/response dan entitas basis data.
2. **Repository**: Lapisan akses data (GORM) yang berkomunikasi langsung dengan PostgreSQL.
3. **Service**: Lapisan logika bisnis (business logic).
4. **Handler**: Lapisan pengontrol HTTP (Gin router/handlers) yang menggunakan response envelope standar.
5. **App Wiring (`app.go`)**: Fungsi `StartApp(cfg)` yang merakit ketergantungan (dependency injection) repo → service → handler.

---

## 2. Modul Utama Domain

Berdasarkan analisis kluster sistem, arsitektur terbagi menjadi beberapa modul utama:

- **Medicine Module (`internal/medicine`)**: Mengelola data obat, pencarian berdasarkan barcode, stok menipis (*low stock*), tanggal kedaluwarsa (*expiring*), dan status aktif/nonaktif.
- **Drug Type Module (`internal/drug`)**: Mengelola kategori atau jenis obat (*type-drugs*).
- **Transaction Module (`internal/transaction`)**: Modul terbesar yang menangani keranjang belanja (*cart*), pemrosesan pembayaran (*payment*), riwayat transaksi, cetak struk PDF/WhatsApp, serta alur pengembalian barang (*returns*) dengan threshold persetujuan.
- **User Module (`internal/user`)**: Mengelola autentikasi (JWT), manajemen profil, peran (*role*), dan pembaruan sandi.
- **Report Module (`internal/report`)**: Menghasilkan laporan penjualan harian, pemeringkatan obat (*drug ranking*), serta ekspor laporan ke format Excel dan PDF.

---

## 3. Alur Komunikasi & Titik Masuk (Entrypoint)

- **Entrypoint Utama**: `cmd/api/main.go`
- **Inisialisasi Saat Startup**:
  1. Memuat konfigurasi dari `config/config.yaml` menggunakan **Viper**.
  2. Menjalankan migrasi basis data secara otomatis via **golang-migrate** dari direktori `migrations/data`.
  3. Menghubungkan ke PostgreSQL dan memastikan ekstensi `uuid-ossp` aktif.
  4. Menyiapkan CORS dan mendaftarkan router grup API untuk masing-masing modul (`/api/v1/...`).

---

## 4. Inti Abstraksi & God Nodes (Komponen Sentral)

Berdasarkan metrik konektivitas graph sistem, komponen dengan dependensi tertinggi (*God Nodes*) yang menjadi pilar utama aplikasi meliputi:

1. **`Success()` (`internal/common/response`)**: Helper pembungkus response sukses global yang digunakan secara masif di seluruh handler.
2. **`TransactionRepository` / `TransactionService` / `TransactionHandler`**: Tulang punggung operasional transaksi kasir dan keranjang.
3. **`UserService` / `UserHandler`**: Pusat autentikasi, manajemen sesi, dan keamanan pengguna.
4. **`MedicineRepository` / `MedicineService` / `MedicineHandler`**: Pusat inventaris dan pengelolaan data obat.

---

## 5. Keamanan & Autentikasi

- **JWT (JSON Web Token)**: Autentikasi rute terlindungi menggunakan `middleware.Auth()`. Middleware ini mendekode token dan menyuntikkan `user_id`, `username`, dan `role` ke dalam context Gin (`c.Set`).
- **Role-Based Access Control**: Menggunakan `middleware.RequireRoles(...)` untuk membatasi akses endpoint sensitif berdasarkan peran pengguna (misal: `admin`, `kasir`).

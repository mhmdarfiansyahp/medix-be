# Security Documentation & Vulnerability Assessment

## Overview

Medix Backend implements security best practices across authentication, authorization, data validation, and asset management. This document outlines security architecture and addresses key vulnerability considerations.

---

## 🔒 Tiga Celah Keamanan & Mitigasi (Top 3 Security Considerations)

### 1. JWT Secret Strength & Environment Exposure
- **Risiko / Celah:** Jika `JWT_SECRET` menggunakan string yang lemah atau hardcoded di konfigurasi default (`config/config.yaml`), token dapat di-forge (dipalsukan) oleh pihak luar untuk mendapatkan akses administrator.
- **Mitigasi di Sistem:**
  - Middleware autentikasi (`internal/middleware/AuthMiddleware.go`) wajib membaca `JWT_SECRET` dari environment variable (`os.Getenv("JWT_SECRET")`) atau konfigurasi aman.
  - Sesi token memiliki masa kedaluwarsa (expiration time) yang ketat dan mendukung mekanisme token refresh serta logout revokasi via tabel `sessions`.
  - **Rekomendasi Produksi:** Gunakan secret key acak minimal 32 karakter dan simpan di secret manager / environment variable sistem produksi (bukan di commit git).

### 2. Validasi & Penyimpanan File Upload (Path Traversal / Arbitrary File Upload)
- **Risiko / Celah:** Endpoint upload foto profil (`POST /users/profile/photo`) menerima file multipart. Jika ekstensi dan nama file tidak divalidasi ketat, penyerang dapat mengunggah file berbahaya (misal: script berbahaya) atau melakukan path traversal.
- **Mitigasi di Sistem:**
  - Validasi ukuran file maksimal 2MB (`file.Size > 2<<20`).
  - Whitelist ekstensi file secara ketat: `.jpg`, `.jpeg`, `.png`, `.webp`.
  - Penamaan file di-generate ulang secara aman di server menggunakan format `user_{id}_{timestamp}{ext}` (mengabaikan nama asli file dari user).
  - File disimpan di direktori terisolasi (`uploads/profiles`).

### 3. Insecure Direct Object Reference (IDOR) & Role Authorization
- **Risiko / Celah:** Pengguna (kasir/pelanggan) mencoba mengakses atau memanipulasi data transaksi atau profil milik pengguna lain dengan menebak ID (`/transactions/{id}`).
- **Mitigasi di Sistem:**
  - Implementasi RBAC (`RequireRoles`) pada route sensitif (admin, kasir, owner).
  - Validasi kepemilikan resource pada handler transaksi (`internal/transaction/handler/TransactionHandler.go`), di mana system memeriksa apakah `transaction.IDUser == userID` sebelum mengembalikan data sensitif atau riwayat harga.
  - Pemisahan route publik dan protected untuk mencegah akses tanpa token yang valid.

---

## 🛡️ Best Practices & Standard Security Controls

- **Password Hashing:** Menggunakan pustaka **Bcrypt** dengan cost factor yang aman untuk menyimpan hash password di database.
- **SQL Injection Prevention:** Menggunakan ORM **GORM** dengan parameterisasi query otomatis untuk mencegah celah SQL Injection.
- **CORS Policy:** Dibatasi secara eksplisit (default: `http://localhost:5173` untuk frontend Vite).
- **Database Constraints:** Validasi integritas data di level database menggunakan `CHECK` constraints, foreign keys, dan unique indexes.

## Reporting Security Issues
Jika Anda menemukan kerentanan keamanan (vulnerability) pada proyek ini, harap laporkan kepada tim pengembang internal Medix.

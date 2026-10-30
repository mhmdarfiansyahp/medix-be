# Database Documentation

## Entity Relationship Diagram

Refer to `docs/diagrams/er-diagram.mmd` for the visual representation.

## Tables

### 1. `users`
Stores system users (Admin, Kasir, Owner).

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| `id_user` | `serial` | `PK` | Unique user identifier |
| `nama_user` | `varchar(100)` | `NOT NULL` | Full name |
| `no_telp` | `varchar(20)` | | Phone number |
| `role` | `varchar(20)` | `NOT NULL` | `admin`, `kasir`, `owner` |
| `username` | `varchar(50)` | `UNIQUE`, `NOT NULL` | Login username |
| `password` | `varchar(255)` | `NOT NULL` | Bcrypt hashed password |
| `email` | `varchar(100)` | `UNIQUE` | User email |
| `status` | `smallint` | `DEFAULT 1`, `CHECK (1,0)` | `1`: aktif, `0`: nonaktif |
| `foto` | `varchar(255)` | | Path to profile image |
| `must_change_password` | `boolean` | `DEFAULT false` | Force password reset |
| `created_at` | `timestamptz` | `DEFAULT now()` | Creation timestamp |
| `updated_at` | `timestamptz` | `DEFAULT now()` | Last update timestamp |

### 2. `jenis_obat` (Golongan Obat)
Drug categorization.

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| `id_jenis` | `serial` | `PK` | Unique category identifier |
| `nama_jenis` | `varchar(50)` | `UNIQUE`, `NOT NULL` | Category name |

### 3. `obat`
Inventory data for medicines.

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| `id_obat` | `serial` | `PK` | Unique medicine identifier |
| `nama_obat` | `varchar(150)` | `NOT NULL` | Medicine name |
| `merk_obat` | `varchar(100)` | | Brand name |
| `jenis_obat_id` | `int` | `FK (jenis_obat.id_jenis)` | Reference to category |
| `barcode` | `varchar(50)` | `UNIQUE` (partial) | Scan code (nullable, unique when present) |
| `tgl_kadaluarsa` | `date` | `NOT NULL` | Expiry date |
| `harga` | `numeric(12,2)` | `NOT NULL`, `CHECK (>=0)` | Selling price |
| `stok` | `int` | `NOT NULL`, `DEFAULT 0`, `CHECK (>=0)` | Current stock |
| `stok_minimum` | `int` | `DEFAULT 10` | Low stock threshold |
| `keterangan` | `text` | | Notes/Usage |
| `status` | `smallint` | `DEFAULT 1`, `CHECK (1,0)` | `1`: aktif, `0`: nonaktif |
| `gambar` | `varchar(255)` | | Product image path |
| `created_at` | `timestamptz` | `DEFAULT now()` | Creation timestamp |
| `updated_at` | `timestamptz` | `DEFAULT now()` | Last update timestamp |

**Indexes:**
- `idx_obat_status` on `status`
- `idx_obat_kadaluarsa` on `tgl_kadaluarsa`
- `idx_obat_nama_lower` on `LOWER(nama_obat)`
- `idx_obat_jenis` on `jenis_obat_id`
- `idx_obat_barcode_unique` unique on `barcode` WHERE `barcode IS NOT NULL AND barcode <> ''`

### 4. `transaksi`
Sales transaction records.

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| `id_transaksi` | `serial` | `PK` | Unique transaction identifier |
| `id_user` | `int` | `FK (users.id_user)`, `NOT NULL` | Reference to creator (kasir/admin) |
| `tgl_transaksi` | `timestamptz` | `DEFAULT now()` | Date of transaction |
| `total_harga` | `numeric(14,2)` | `NOT NULL`, `CHECK (>=0)` | Grand total |
| `metode_bayar` | `varchar(20)` | | Payment method |
| `uang_diterima` | `numeric(14,2)` | `NOT NULL`, `DEFAULT 0`, `CHECK (>=0)` | Cash received |
| `kembalian` | `numeric(14,2)` | `NOT NULL`, `DEFAULT 0`, `CHECK (>=0)` | Change amount |
| `status` | `smallint` | `DEFAULT 1`, `CHECK (1,0,2)` | `1`: normal, `0`: cancelled/deleted, `2`: completed |
| `status_pembayaran` | `varchar(20)` | `DEFAULT 'pending'` | `pending`, `paid`, `cancelled` |
| `created_at` | `timestamptz` | `DEFAULT now()` | Log creation time |

**Indexes:**
- `idx_transaksi_user` on `id_user`
- `idx_transaksi_tgl` on `tgl_transaksi`

### 5. `detail_pembelian`
Line items for each transaction.

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| `id_detail` | `serial` | `PK` | Unique detail identifier |
| `id_transaksi` | `int` | `FK (transaksi.id_transaksi)`, `ON DELETE CASCADE` | Parent transaction |
| `id_obat` | `int` | `FK (obat.id_obat)`, `NOT NULL` | Sold medicine |
| `jumlah` | `int` | `NOT NULL`, `CHECK (>0)` | Quantity sold |
| `harga_satuan` | `numeric(12,2)` | `NOT NULL` | Price at time of sale |
| `subtotal` | `numeric(14,2)` | `NOT NULL` | `jumlah * harga_satuan` |

**Indexes:**
- `idx_detail_transaksi` on `id_transaksi`

### 6. `transaksi_retur`
Return transaction (for cancellations/refunds).

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| `id_return` | `serial` | `PK` | Unique return identifier |
| `id_transaksi` | `int` | `FK (transaksi.id_transaksi)`, `NOT NULL` | Original transaction |
| `alasan` | `text` | `NOT NULL` | Return reason |
| `tanggal_retur` | `timestamptz` | `NOT NULL`, `DEFAULT now()` | Return date |
| `diajukan_oleh` | `int` | `FK (users.id_user)`, `NOT NULL` | User who initiated return |
| `status` | `varchar(20)` | `NOT NULL`, `DEFAULT 'pending'`, `CHECK` | `pending`, `approved`, `rejected` |
| `alasan_penolakan` | `text` | | Rejection reason |

**Indexes:**
- `idx_retur_transaksi` on `id_transaksi`

### 7. `transaksi_retur_item`
Items for return transaction.

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| `id_item` | `serial` | `PK` | Unique item identifier |
| `id_return` | `int` | `FK (transaksi_retur.id_return)`, `NOT NULL` | Parent return |
| `id_obat` | `int` | `FK (obat.id_obat)`, `NOT NULL` | Returned medicine |
| `jumlah` | `int` | `NOT NULL`, `CHECK (>0)` | Quantity returned |
| `alasan_item` | `text` | | Item return reason |
| `kondisi_layak` | `boolean` | `NOT NULL`, `DEFAULT false` | Is condition good? |
| `stok_kembali` | `int` | `NOT NULL`, `DEFAULT 0` | Returned stock |

**Indexes:**
- `idx_retur_item_retur` on `id_return`

### 8. `sessions`
Web/Mobile user sessions.

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| `id_session` | `uuid` | `PK`, `DEFAULT uuid_generate_v4()` | Session ID |
| `id_user` | `int` | `FK (users.id_user)`, `ON DELETE CASCADE` | Owner |
| `refresh_token` | `varchar(255)` | `NOT NULL`, `UNIQUE` | JWT refresh token |
| `expires_at` | `timestamptz` | `NOT NULL` | Expiry time |
| `revoked` | `boolean` | `NOT NULL`, `DEFAULT false` | Is revoked? |
| `created_at` | `timestamptz` | `DEFAULT now()` | Creation timestamp |

**Indexes:**
- `idx_sessions_refresh_token` on `refresh_token`
- `idx_sessions_user_id` on `id_user`

### 9. `password_reset_tokens`
Password reset tokens for user authentication.

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| `id` | `serial` | `PK` | Unique identifier |
| `id_user` | `int` | `FK (users.id_user)`, `ON DELETE CASCADE` | User owner |
| `token` | `varchar(255)` | `NOT NULL`, `UNIQUE` | Reset token |
| `expires_at` | `timestamp` | `NOT NULL` | Expiry time |
| `used` | `boolean` | `NOT NULL`, `DEFAULT false` | Is token used? |
| `created_at` | `timestamp` | `NOT NULL`, `DEFAULT now()` | Creation timestamp |

### 10. `pengaturan`
System settings.

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| `key` | `varchar(100)` | `PK` | Setting key |
| `value` | `text` | `NOT NULL` | Setting value |

**Sample setting:**
- `key`: `retur_approval_threshold`, `value`: `100000` (return amount threshold for approval)

## Database Engine
- **PostgreSQL 14+**
- Extensions: `uuid-ossp` (for sessions)

## Migration Tool
- `golang-migrate/migrate`
- Files: `migrations/data/*.sql`

## Migration History
| File | Description |
|------|-------------|
| 000001 | Create users table |
| 000002 | Create jenis_obat table |
| 000003 | Create obat table (+ indexes) |
| 000004 | Create transaksi table (+ indexes) |
| 000005 | Create detail_pembelian table (+ index) |
| 000006 | Seed jenis_obat |
| 000007 | Seed users |
| 000008 | Seed obat |
| 000009 | Seed transaksi |
| 000010 | Seed detail_pembelian |
| 000011 | Create sessions table (+ indexes) |
| 000012 | Add payment columns to transaksi; create transaksi_retur & transaksi_retur_item |
| 000013 | Add must_change_password to users |
| 000014 | Add alasan_penolakan to transaksi_retur; create pengaturan table |
| 000015 | Update transaksi.status CHECK constraint (add value 2) |
| 000016 | Create password_reset_tokens table |
| 000017 | Add email column to users |
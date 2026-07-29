CREATE TABLE obat (
    id_obat        SERIAL PRIMARY KEY,
    nama_obat      VARCHAR(150) NOT NULL,
    merk_obat      VARCHAR(100),
    jenis_obat_id  INT REFERENCES jenis_obat(id_jenis),
    barcode        VARCHAR(50) UNIQUE,
    tgl_kadaluarsa DATE NOT NULL,
    harga          NUMERIC(12,2) NOT NULL CHECK (harga >= 0),
    stok           INT NOT NULL DEFAULT 0 CHECK (stok >= 0),
    stok_minimum   INT DEFAULT 10,
    keterangan     TEXT,
    status         VARCHAR(20) NOT NULL DEFAULT 'aktif' CHECK (status IN ('aktif','nonaktif')),
    gambar         VARCHAR(255),
    created_at     TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at     TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_obat_status ON obat(status);
CREATE INDEX idx_obat_kadaluarsa ON obat(tgl_kadaluarsa);
CREATE TABLE obat (
    id_obat        SERIAL PRIMARY KEY,
    nama_obat      VARCHAR(150) NOT NULL,
    merk_obat      VARCHAR(100),
    jenis_obat_id  INT REFERENCES jenis_obat(id_jenis),
    barcode        VARCHAR(50),
    tgl_kadaluarsa DATE NOT NULL,
    harga          NUMERIC(12,2) NOT NULL CHECK (harga >= 0),
    stok           INT NOT NULL DEFAULT 0 CHECK (stok >= 0),
    stok_minimum   INT DEFAULT 10,
    keterangan     TEXT,
    status         SMALLINT NOT NULL DEFAULT 1 CHECK (status IN (1, 0)),
    gambar         VARCHAR(255),
    created_at     TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at     TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_obat_status ON obat(status);

CREATE INDEX idx_obat_kadaluarsa ON obat(tgl_kadaluarsa);

CREATE INDEX idx_obat_nama_lower ON obat (LOWER(nama_obat));

CREATE INDEX idx_obat_jenis ON obat (jenis_obat_id);

CREATE UNIQUE INDEX idx_obat_barcode_unique ON obat (barcode) 
WHERE barcode IS NOT NULL AND barcode <> '';
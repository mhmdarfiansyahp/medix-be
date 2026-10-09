ALTER TABLE transaksi
    ADD COLUMN metode_bayar  VARCHAR(20),
    ADD COLUMN uang_diterima NUMERIC(14,2) NOT NULL DEFAULT 0 CHECK (uang_diterima >= 0),
    ADD COLUMN kembalian     NUMERIC(14,2) NOT NULL DEFAULT 0 CHECK (kembalian >= 0);

CREATE TABLE transaksi_retur (
    id_return      SERIAL PRIMARY KEY,
    id_transaksi   INT NOT NULL REFERENCES transaksi(id_transaksi),
    alasan         TEXT NOT NULL,
    tanggal_retur  TIMESTAMPTZ NOT NULL DEFAULT now(),
    diajukan_oleh  INT NOT NULL REFERENCES users(id_user),
    status         VARCHAR(20) NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'approved', 'rejected'))
);

CREATE INDEX idx_retur_transaksi ON transaksi_retur(id_transaksi);

CREATE TABLE transaksi_retur_item (
    id_item       SERIAL PRIMARY KEY,
    id_return     INT NOT NULL REFERENCES transaksi_retur(id_return),
    id_obat       INT NOT NULL REFERENCES obat(id_obat),
    jumlah        INT NOT NULL CHECK (jumlah > 0),
    alasan_item   TEXT,
    kondisi_layak BOOLEAN NOT NULL DEFAULT false,
    stok_kembali  INT NOT NULL DEFAULT 0
);

CREATE INDEX idx_retur_item_retur ON transaksi_retur_item(id_return);

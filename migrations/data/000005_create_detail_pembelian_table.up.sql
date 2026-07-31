CREATE TABLE detail_pembelian (
    id_detail    SERIAL PRIMARY KEY,
    id_transaksi INT NOT NULL REFERENCES transaksi(id_transaksi) ON DELETE CASCADE,
    id_obat      INT NOT NULL REFERENCES obat(id_obat),
    jumlah       INT NOT NULL CHECK (jumlah > 0),
    harga_satuan NUMERIC(12,2) NOT NULL,
    subtotal     NUMERIC(14,2) NOT NULL
);

CREATE INDEX idx_detail_transaksi ON detail_pembelian(id_transaksi);
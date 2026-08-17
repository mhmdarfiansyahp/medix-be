CREATE TABLE transaksi (
    id_transaksi  SERIAL PRIMARY KEY,
    id_user       INT NOT NULL REFERENCES users(id_user),
    tgl_transaksi TIMESTAMPTZ NOT NULL DEFAULT now(),
    total_harga   NUMERIC(14,2) NOT NULL CHECK (total_harga >= 0),
    status         SMALLINT NOT NULL DEFAULT 1 CHECK (status IN (1, 0)),
    created_at    TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_transaksi_user ON transaksi(id_user);
CREATE INDEX idx_transaksi_tgl ON transaksi(tgl_transaksi);
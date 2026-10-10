ALTER TABLE transaksi_retur
    ADD COLUMN alasan_penolakan TEXT;

CREATE TABLE pengaturan (
    key   VARCHAR(100) PRIMARY KEY,
    value TEXT NOT NULL
);

INSERT INTO pengaturan (key, value) VALUES ('retur_approval_threshold', '100000');

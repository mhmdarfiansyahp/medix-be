ALTER TABLE transaksi DROP CONSTRAINT transaksi_status_check;
ALTER TABLE transaksi ADD CONSTRAINT transaksi_status_check CHECK (status IN (1, 0, 2));

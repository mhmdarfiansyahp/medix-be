CREATE TABLE jenis_obat (
    id_jenis   SERIAL PRIMARY KEY,
    nama_jenis VARCHAR(50) NOT NULL UNIQUE
);
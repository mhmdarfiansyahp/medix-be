CREATE TABLE users (
    id_user    SERIAL PRIMARY KEY,
    nama_user  VARCHAR(100) NOT NULL,
    no_telp    VARCHAR(20),
    role       VARCHAR(20) NOT NULL CHECK (role IN ('admin','kasir','owner')),
    username   VARCHAR(50) NOT NULL UNIQUE,
    password   VARCHAR(255) NOT NULL,
    status     VARCHAR(20) NOT NULL DEFAULT 'aktif' CHECK (status IN ('aktif','nonaktif')),
    foto       VARCHAR(255),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
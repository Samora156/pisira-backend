-- ============================================================
--  PISIRA - Database Service Laptop (PostgreSQL)
-- ============================================================

-- Buat database (jalankan terpisah jika perlu)
-- CREATE DATABASE pisira_db;

-- ─── Tabel: users ────────────────────────────────────────────
CREATE TABLE IF NOT EXISTS users (
  id          SERIAL        PRIMARY KEY,
  nama        VARCHAR(100)  NOT NULL,
  email       VARCHAR(100)  NOT NULL UNIQUE,
  password    VARCHAR(255)  NOT NULL,
  role        VARCHAR(20)   NOT NULL DEFAULT 'teknisi'
                            CHECK (role IN ('admin', 'teknisi')),
  is_active   BOOLEAN       NOT NULL DEFAULT TRUE,
  created_at  TIMESTAMPTZ   NOT NULL DEFAULT NOW(),
  updated_at  TIMESTAMPTZ   NOT NULL DEFAULT NOW()
);

-- ─── Tabel: customers ────────────────────────────────────────
CREATE TABLE IF NOT EXISTS customers (
  id          SERIAL        PRIMARY KEY,
  nama        VARCHAR(100)  NOT NULL,
  no_hp       VARCHAR(20)   NOT NULL,
  email       VARCHAR(100),
  alamat      TEXT,
  created_at  TIMESTAMPTZ   NOT NULL DEFAULT NOW(),
  updated_at  TIMESTAMPTZ   NOT NULL DEFAULT NOW()
);

-- ─── Tabel: service_orders ───────────────────────────────────
CREATE TABLE IF NOT EXISTS service_orders (
  id                SERIAL        PRIMARY KEY,
  customer_id       INT           NOT NULL REFERENCES customers(id),
  no_order          VARCHAR(30)   NOT NULL UNIQUE,
  merk_laptop       VARCHAR(50)   NOT NULL,
  model_laptop      VARCHAR(100)  NOT NULL,
  sn_laptop         VARCHAR(100),
  keluhan           TEXT          NOT NULL,
  diagnosa          TEXT,
  status            VARCHAR(30)   NOT NULL DEFAULT 'menunggu'
                                  CHECK (status IN (
                                    'menunggu','diagnosa','menunggu_persetujuan',
                                    'proses','selesai','diambil','batal'
                                  )),
  tanggal_masuk     DATE          NOT NULL DEFAULT CURRENT_DATE,
  tanggal_estimasi  TIMESTAMPTZ,
  tanggal_selesai   TIMESTAMPTZ,
  tanggal_ambil     TIMESTAMPTZ,
  teknisi_id        INT           REFERENCES users(id),
  catatan_teknisi   TEXT,
  created_at        TIMESTAMPTZ   NOT NULL DEFAULT NOW(),
  updated_at        TIMESTAMPTZ   NOT NULL DEFAULT NOW()
);

-- ─── Tabel: estimasi ─────────────────────────────────────────
CREATE TABLE IF NOT EXISTS estimasi (
  id                   SERIAL          PRIMARY KEY,
  order_id             INT             NOT NULL UNIQUE REFERENCES service_orders(id),
  deskripsi_pekerjaan  TEXT            NOT NULL,
  biaya_jasa           NUMERIC(12, 2)  NOT NULL DEFAULT 0,
  biaya_sparepart      NUMERIC(12, 2)  NOT NULL DEFAULT 0,
  total                NUMERIC(12, 2)  NOT NULL DEFAULT 0,
  status_persetujuan   VARCHAR(20)     NOT NULL DEFAULT 'menunggu'
                                       CHECK (status_persetujuan IN ('menunggu','disetujui','ditolak')),
  catatan              TEXT,
  created_at           TIMESTAMPTZ     NOT NULL DEFAULT NOW(),
  updated_at           TIMESTAMPTZ     NOT NULL DEFAULT NOW()
);

-- ─── Tabel: invoice ──────────────────────────────────────────
CREATE TABLE IF NOT EXISTS invoice (
  id               SERIAL          PRIMARY KEY,
  order_id         INT             NOT NULL UNIQUE REFERENCES service_orders(id),
  no_invoice       VARCHAR(30)     NOT NULL UNIQUE,
  subtotal         NUMERIC(12, 2)  NOT NULL DEFAULT 0,
  diskon           NUMERIC(12, 2)  NOT NULL DEFAULT 0,
  total_bayar      NUMERIC(12, 2)  NOT NULL DEFAULT 0,
  metode_bayar     VARCHAR(20)     NOT NULL DEFAULT 'tunai'
                                   CHECK (metode_bayar IN ('tunai','transfer','qris')),
  status_bayar     VARCHAR(20)     NOT NULL DEFAULT 'belum_lunas'
                                   CHECK (status_bayar IN ('belum_lunas','lunas')),
  tanggal_invoice  TIMESTAMPTZ     NOT NULL DEFAULT NOW(),
  created_at       TIMESTAMPTZ     NOT NULL DEFAULT NOW(),
  updated_at       TIMESTAMPTZ     NOT NULL DEFAULT NOW()
);

-- ─── Index ───────────────────────────────────────────────────
CREATE INDEX IF NOT EXISTS idx_customers_no_hp       ON customers      (no_hp);
CREATE INDEX IF NOT EXISTS idx_orders_status         ON service_orders (status);
CREATE INDEX IF NOT EXISTS idx_orders_tanggal_masuk  ON service_orders (tanggal_masuk);
CREATE INDEX IF NOT EXISTS idx_invoice_status_bayar  ON invoice        (status_bayar);

-- ─── Trigger auto-update updated_at ──────────────────────────
CREATE OR REPLACE FUNCTION set_updated_at()
RETURNS TRIGGER AS $$
BEGIN
  NEW.updated_at = NOW();
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE TRIGGER trg_users_updated_at
  BEFORE UPDATE ON users
  FOR EACH ROW EXECUTE FUNCTION set_updated_at();

CREATE OR REPLACE TRIGGER trg_customers_updated_at
  BEFORE UPDATE ON customers
  FOR EACH ROW EXECUTE FUNCTION set_updated_at();

CREATE OR REPLACE TRIGGER trg_orders_updated_at
  BEFORE UPDATE ON service_orders
  FOR EACH ROW EXECUTE FUNCTION set_updated_at();

CREATE OR REPLACE TRIGGER trg_estimasi_updated_at
  BEFORE UPDATE ON estimasi
  FOR EACH ROW EXECUTE FUNCTION set_updated_at();

CREATE OR REPLACE TRIGGER trg_invoice_updated_at
  BEFORE UPDATE ON invoice
  FOR EACH ROW EXECUTE FUNCTION set_updated_at();

-- ─── Data awal: akun admin ───────────────────────────────────
-- Ganti hash di bawah dengan hasil bcrypt dari password Anda
-- Gunakan: https://bcrypt-generator.com (cost factor 10)
INSERT INTO users (nama, email, password, role)
VALUES ('Admin PISIRA', 'admin@pisira.com', '$2b$10$placeholder_ganti_dengan_hash_bcrypt', 'admin')
ON CONFLICT (email) DO NOTHING;

-- ─── Contoh data customer & order (untuk testing) ────────────
INSERT INTO customers (nama, no_hp, email, alamat) VALUES
  ('Budi Santoso', '081234567890', 'budi@email.com', 'Jl. Merdeka No. 10, Malang'),
  ('Siti Rahayu',  '082345678901', 'siti@email.com', 'Jl. Ijen No. 5, Malang'),
  ('Andi Wijaya',  '083456789012', NULL,              'Jl. Kawi No. 3, Malang')
ON CONFLICT DO NOTHING;

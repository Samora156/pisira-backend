# PISIRA Backend — Go + Gin + MySQL

Backend API untuk aplikasi manajemen servis laptop PISIRA.

## Teknologi
- **Go 1.22** — bahasa pemrograman utama
- **Gin** — web framework HTTP yang ringan dan cepat
- **sqlx** — helper query MySQL
- **JWT** — autentikasi token
- **bcrypt** — enkripsi password

---

## Cara Setup

### 1. Install Go
Download dari https://go.dev/dl/ (pilih versi 1.22+)

### 2. Clone / buat folder project
```bash
mkdir pisira-backend && cd pisira-backend
```

### 3. Install dependencies
```bash
go mod tidy
```

### 4. Siapkan database
Jalankan file SQL yang sudah dibuat sebelumnya:
```bash
mysql -u root -p < pisira_database.sql
```

### 5. Buat file .env
```bash
cp .env.example .env
```
Lalu edit `.env` sesuai konfigurasi database Anda.

### 6. Buat password admin (jalankan sekali)
```bash
# Contoh membuat hash bcrypt di Go
go run scripts/hash_password.go
```
Atau gunakan website: https://bcrypt-generator.com/
Salin hasilnya dan update tabel users di database:
```sql
UPDATE users SET password = 'HASH_BCRYPT_ANDA' WHERE email = 'admin@pisira.com';
```

### 7. Jalankan server
```bash
go run cmd/api/main.go
```
Server berjalan di: http://localhost:8080

---

## Daftar Endpoint API

### Auth
| Method | Endpoint | Keterangan |
|--------|----------|------------|
| POST | /api/auth/login | Login, mendapatkan token JWT |

### Customer
| Method | Endpoint | Keterangan |
|--------|----------|------------|
| GET | /api/customers | Daftar semua customer |
| GET | /api/customers/:id | Detail customer |
| POST | /api/customers | Tambah customer baru |
| PUT | /api/customers/:id | Update data customer |

### Service Order
| Method | Endpoint | Keterangan |
|--------|----------|------------|
| GET | /api/orders | Daftar semua order (filter: ?status=proses) |
| GET | /api/orders/:id | Detail order |
| POST | /api/orders | Buat order baru |
| PATCH | /api/orders/:id/status | Update status order |

### Estimasi
| Method | Endpoint | Keterangan |
|--------|----------|------------|
| GET | /api/orders/:id/estimasi | Lihat estimasi sebuah order |
| POST | /api/estimasi | Buat estimasi baru |
| PATCH | /api/orders/:id/estimasi/persetujuan | Update persetujuan customer |

### Invoice
| Method | Endpoint | Keterangan |
|--------|----------|------------|
| GET | /api/invoices | Daftar invoice |
| POST | /api/invoices | Buat invoice baru |
| PATCH | /api/invoices/:order_id/lunas | Tandai invoice lunas |

### Laporan (Admin only)
| Method | Endpoint | Keterangan |
|--------|----------|------------|
| GET | /api/laporan/bulanan?tahun=2024 | Laporan bulanan per tahun |

---

## Contoh Request

### Login
```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@pisira.com","password":"password123"}'
```

### Tambah Order (dengan token)
```bash
curl -X POST http://localhost:8080/api/orders \
  -H "Authorization: Bearer TOKEN_ANDA" \
  -H "Content-Type: application/json" \
  -d '{
    "customer_id": 1,
    "merk_laptop": "ASUS",
    "model_laptop": "VivoBook 14",
    "keluhan": "Layar retak dan tidak bisa menyala"
  }'
```

---

## Struktur Folder
```
pisira-backend/
├── cmd/api/main.go          ← entry point, setup server & routing
├── internal/
│   ├── config/config.go     ← load konfigurasi .env
│   ├── handler/handler.go   ← terima request HTTP, kirim response
│   ├── service/service.go   ← logika bisnis
│   ├── repository/          ← query ke database
│   ├── model/model.go       ← struct data
│   └── middleware/auth.go   ← cek JWT token
├── pkg/response/            ← format response standar
├── .env.example
└── go.mod
```

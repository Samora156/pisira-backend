package model

import "time"

// ─── User ────────────────────────────────────────────────────

type User struct {
	ID        int       `db:"id"         json:"id"`
	Nama      string    `db:"nama"       json:"nama"`
	Email     string    `db:"email"      json:"email"`
	Password  string    `db:"password"   json:"-"`
	Role      string    `db:"role"       json:"role"`
	IsActive  bool      `db:"is_active"  json:"is_active"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

type LoginRequest struct {
	Email    string `json:"email"    binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type LoginResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}

// ─── Customer ────────────────────────────────────────────────

type Customer struct {
	ID        int       `db:"id"         json:"id"`
	Nama      string    `db:"nama"       json:"nama"`
	NoHP      string    `db:"no_hp"      json:"no_hp"`
	Email     *string   `db:"email"      json:"email"`
	Alamat    *string   `db:"alamat"     json:"alamat"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

type CreateCustomerRequest struct {
	Nama   string  `json:"nama"   binding:"required"`
	NoHP   string  `json:"no_hp"  binding:"required"`
	Email  *string `json:"email"`
	Alamat *string `json:"alamat"`
}

// ─── Service Order ───────────────────────────────────────────

type ServiceOrder struct {
	ID              int        `db:"id"                json:"id"`
	CustomerID      int        `db:"customer_id"       json:"customer_id"`
	NoOrder         string     `db:"no_order"          json:"no_order"`
	MerkLaptop      string     `db:"merk_laptop"       json:"merk_laptop"`
	ModelLaptop     string     `db:"model_laptop"      json:"model_laptop"`
	SNLaptop        *string    `db:"sn_laptop"         json:"sn_laptop"`
	Keluhan         string     `db:"keluhan"           json:"keluhan"`
	Diagnosa        *string    `db:"diagnosa"          json:"diagnosa"`
	Status          string     `db:"status"            json:"status"`
	TanggalMasuk    time.Time  `db:"tanggal_masuk"     json:"tanggal_masuk"`
	TanggalEstimasi *time.Time `db:"tanggal_estimasi"  json:"tanggal_estimasi"`
	TanggalSelesai  *time.Time `db:"tanggal_selesai"   json:"tanggal_selesai"`
	TanggalAmbil    *time.Time `db:"tanggal_ambil"     json:"tanggal_ambil"`
	TeknisiID       *int       `db:"teknisi_id"        json:"teknisi_id"`
	CatatanTeknisi  *string    `db:"catatan_teknisi"   json:"catatan_teknisi"`
	CreatedAt       time.Time  `db:"created_at"        json:"created_at"`
	UpdatedAt       time.Time  `db:"updated_at"        json:"updated_at"`
	// field join
	CustomerNama    string     `db:"customer_nama"     json:"customer_nama"`
	CustomerNoHP    string     `db:"customer_no_hp"    json:"customer_no_hp"`
}

type CreateOrderRequest struct {
	CustomerID  int     `json:"customer_id"  binding:"required"`
	MerkLaptop  string  `json:"merk_laptop"  binding:"required"`
	ModelLaptop string  `json:"model_laptop" binding:"required"`
	SNLaptop    *string `json:"sn_laptop"`
	Keluhan     string  `json:"keluhan"      binding:"required"`
}

type UpdateOrderStatusRequest struct {
	Status         string  `json:"status"          binding:"required"`
	Diagnosa       *string `json:"diagnosa"`
	CatatanTeknisi *string `json:"catatan_teknisi"`
}

// ─── Estimasi ─────────────────────────────────────────────────

type Estimasi struct {
	ID                 int       `db:"id"                  json:"id"`
	OrderID            int       `db:"order_id"            json:"order_id"`
	DeskripsIPekerjaan string    `db:"deskripsi_pekerjaan" json:"deskripsi_pekerjaan"`
	BiayaJasa          float64   `db:"biaya_jasa"          json:"biaya_jasa"`
	BiayaSparepart     float64   `db:"biaya_sparepart"     json:"biaya_sparepart"`
	Total              float64   `db:"total"               json:"total"`
	StatusPersetujuan  string    `db:"status_persetujuan"  json:"status_persetujuan"`
	Catatan            *string   `db:"catatan"             json:"catatan"`
	CreatedAt          time.Time `db:"created_at"          json:"created_at"`
	UpdatedAt          time.Time `db:"updated_at"          json:"updated_at"`
}

type CreateEstimasiRequest struct {
	OrderID            int     `json:"order_id"            binding:"required"`
	DeskripsIPekerjaan string  `json:"deskripsi_pekerjaan" binding:"required"`
	BiayaJasa          float64 `json:"biaya_jasa"          binding:"required"`
	BiayaSparepart     float64 `json:"biaya_sparepart"`
	Catatan            *string `json:"catatan"`
}

// ─── Invoice ──────────────────────────────────────────────────

type Invoice struct {
	ID             int       `db:"id"              json:"id"`
	OrderID        int       `db:"order_id"        json:"order_id"`
	NoInvoice      string    `db:"no_invoice"      json:"no_invoice"`
	Subtotal       float64   `db:"subtotal"        json:"subtotal"`
	Diskon         float64   `db:"diskon"          json:"diskon"`
	TotalBayar     float64   `db:"total_bayar"     json:"total_bayar"`
	MetodeBayar    string    `db:"metode_bayar"    json:"metode_bayar"`
	StatusBayar    string    `db:"status_bayar"    json:"status_bayar"`
	TanggalInvoice time.Time `db:"tanggal_invoice" json:"tanggal_invoice"`
	CreatedAt      time.Time `db:"created_at"      json:"created_at"`
	UpdatedAt      time.Time `db:"updated_at"      json:"updated_at"`
	// field join
	CustomerNama   string    `db:"customer_nama"   json:"customer_nama"`
	NoOrder        string    `db:"no_order"        json:"no_order"`
}

type CreateInvoiceRequest struct {
	OrderID     int     `json:"order_id"     binding:"required"`
	Diskon      float64 `json:"diskon"`
	MetodeBayar string  `json:"metode_bayar" binding:"required"`
}

// ─── Laporan ──────────────────────────────────────────────────

type LaporanBulanan struct {
	Bulan         string  `db:"bulan"          json:"bulan"`
	TotalOrder    int     `db:"total_order"    json:"total_order"`
	TotalSelesai  int     `db:"total_selesai"  json:"total_selesai"`
	TotalPendapat float64 `db:"total_pendapat" json:"total_pendapat"`
}

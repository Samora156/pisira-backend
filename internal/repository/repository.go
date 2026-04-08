package repository

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/pisira/backend/internal/model"
)

type Repository struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) *Repository {
	return &Repository{db: db}
}

// ─── User ─────────────────────────────────────────────────────

func (r *Repository) GetUserByEmail(email string) (*model.User, error) {
	var user model.User
	err := r.db.Get(&user, "SELECT * FROM users WHERE email = $1 AND is_active = true", email)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *Repository) GetUserByID(id int) (*model.User, error) {
	var user model.User
	err := r.db.Get(&user, "SELECT * FROM users WHERE id = $1", id)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// ─── Customer ─────────────────────────────────────────────────

func (r *Repository) GetAllCustomers(search string) ([]model.Customer, error) {
	var customers []model.Customer
	query := "SELECT * FROM customers WHERE 1=1"
	args := []interface{}{}
	idx := 1

	if search != "" {
		query += fmt.Sprintf(" AND (nama ILIKE $%d OR no_hp ILIKE $%d)", idx, idx+1)
		s := "%" + search + "%"
		args = append(args, s, s)
	}
	query += " ORDER BY created_at DESC"

	err := r.db.Select(&customers, query, args...)
	return customers, err
}

func (r *Repository) GetCustomerByID(id int) (*model.Customer, error) {
	var c model.Customer
	err := r.db.Get(&c, "SELECT * FROM customers WHERE id = $1", id)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *Repository) CreateCustomer(req model.CreateCustomerRequest) (int, error) {
	var id int
	err := r.db.QueryRow(
		"INSERT INTO customers (nama, no_hp, email, alamat) VALUES ($1, $2, $3, $4) RETURNING id",
		req.Nama, req.NoHP, req.Email, req.Alamat,
	).Scan(&id)
	return id, err
}

func (r *Repository) UpdateCustomer(id int, req model.CreateCustomerRequest) error {
	_, err := r.db.Exec(
		"UPDATE customers SET nama=$1, no_hp=$2, email=$3, alamat=$4 WHERE id=$5",
		req.Nama, req.NoHP, req.Email, req.Alamat, id,
	)
	return err
}

// ─── Service Order ────────────────────────────────────────────

func (r *Repository) GetAllOrders(status, search string) ([]model.ServiceOrder, error) {
	var orders []model.ServiceOrder
	query := `
		SELECT so.*, c.nama AS customer_nama, c.no_hp AS customer_no_hp
		FROM service_orders so
		JOIN customers c ON c.id = so.customer_id
		WHERE 1=1`
	args := []interface{}{}
	idx := 1

	if status != "" {
		query += fmt.Sprintf(" AND so.status = $%d", idx)
		args = append(args, status)
		idx++
	}
	if search != "" {
		query += fmt.Sprintf(" AND (c.nama ILIKE $%d OR so.no_order ILIKE $%d OR so.merk_laptop ILIKE $%d)", idx, idx+1, idx+2)
		s := "%" + search + "%"
		args = append(args, s, s, s)
	}
	query += " ORDER BY so.created_at DESC"

	err := r.db.Select(&orders, query, args...)
	return orders, err
}

func (r *Repository) GetOrderByID(id int) (*model.ServiceOrder, error) {
	var o model.ServiceOrder
	err := r.db.Get(&o, `
		SELECT so.*, c.nama AS customer_nama, c.no_hp AS customer_no_hp
		FROM service_orders so
		JOIN customers c ON c.id = so.customer_id
		WHERE so.id = $1`, id)
	if err != nil {
		return nil, err
	}
	return &o, nil
}

func (r *Repository) CreateOrder(req model.CreateOrderRequest, teknisiID int) (int, error) {
	// Generate nomor order: ORD-YYYYMMDD-XXX
	today := time.Now().Format("20060102")
	var count int
	r.db.Get(&count, "SELECT COUNT(*) FROM service_orders WHERE tanggal_masuk = CURRENT_DATE")
	noOrder := fmt.Sprintf("ORD-%s-%03d", today, count+1)

	var id int
	err := r.db.QueryRow(`
		INSERT INTO service_orders
		  (customer_id, no_order, merk_laptop, model_laptop, sn_laptop, keluhan, tanggal_masuk, teknisi_id)
		VALUES ($1, $2, $3, $4, $5, $6, CURRENT_DATE, $7)
		RETURNING id`,
		req.CustomerID, noOrder, req.MerkLaptop, req.ModelLaptop,
		req.SNLaptop, req.Keluhan, teknisiID,
	).Scan(&id)
	return id, err
}

func (r *Repository) UpdateOrderStatus(id int, req model.UpdateOrderStatusRequest) error {
	_, err := r.db.Exec(`
		UPDATE service_orders
		SET status         = $1,
		    diagnosa       = $2,
		    catatan_teknisi = $3,
		    tanggal_selesai = CASE WHEN $1 = 'selesai' THEN NOW() ELSE tanggal_selesai END,
		    tanggal_ambil   = CASE WHEN $1 = 'diambil' THEN NOW() ELSE tanggal_ambil   END
		WHERE id = $4`,
		req.Status, req.Diagnosa, req.CatatanTeknisi, id,
	)
	return err
}

// ─── Estimasi ─────────────────────────────────────────────────

func (r *Repository) GetEstimasiByOrderID(orderID int) (*model.Estimasi, error) {
	var e model.Estimasi
	err := r.db.Get(&e, "SELECT * FROM estimasi WHERE order_id = $1", orderID)
	if err != nil {
		return nil, err
	}
	return &e, nil
}

func (r *Repository) CreateEstimasi(req model.CreateEstimasiRequest) (int, error) {
	total := req.BiayaJasa + req.BiayaSparepart
	var id int
	err := r.db.QueryRow(`
		INSERT INTO estimasi (order_id, deskripsi_pekerjaan, biaya_jasa, biaya_sparepart, total, catatan)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING order_id`,
		req.OrderID, req.DeskripsIPekerjaan, req.BiayaJasa, req.BiayaSparepart, total, req.Catatan,
	).Scan(&id)
	if err != nil {
		return 0, err
	}
	r.db.Exec(`UPDATE service_orders
		SET status = 'menunggu_persetujuan', tanggal_estimasi = NOW()
		WHERE id = $1`, req.OrderID)
	return id, nil
}

func (r *Repository) UpdatePersetujuanEstimasi(orderID int, status string) error {
	_, err := r.db.Exec(
		"UPDATE estimasi SET status_persetujuan = $1 WHERE order_id = $2", status, orderID,
	)
	if err != nil {
		return err
	}
	if status == "disetujui" {
		r.db.Exec("UPDATE service_orders SET status = 'proses' WHERE id = $1", orderID)
	} else if status == "ditolak" {
		r.db.Exec("UPDATE service_orders SET status = 'batal' WHERE id = $1", orderID)
	}
	return nil
}

// ─── Invoice ──────────────────────────────────────────────────

func (r *Repository) GetAllInvoices(statusBayar string) ([]model.Invoice, error) {
	var invoices []model.Invoice
	query := `
		SELECT inv.*, c.nama AS customer_nama, so.no_order
		FROM invoice inv
		JOIN service_orders so ON so.id = inv.order_id
		JOIN customers c ON c.id = so.customer_id
		WHERE 1=1`
	args := []interface{}{}

	if statusBayar != "" {
		query += " AND inv.status_bayar = $1"
		args = append(args, statusBayar)
	}
	query += " ORDER BY inv.created_at DESC"

	err := r.db.Select(&invoices, query, args...)
	return invoices, err
}

func (r *Repository) CreateInvoice(req model.CreateInvoiceRequest) (int, error) {
	var estimasi model.Estimasi
	if err := r.db.Get(&estimasi, "SELECT * FROM estimasi WHERE order_id = $1", req.OrderID); err != nil {
		return 0, fmt.Errorf("estimasi tidak ditemukan untuk order ini")
	}

	today := time.Now().Format("20060102")
	var count int
	r.db.Get(&count, "SELECT COUNT(*) FROM invoice WHERE created_at::date = CURRENT_DATE")
	noInvoice := fmt.Sprintf("INV-%s-%03d", today, count+1)

	subtotal := estimasi.Total
	totalBayar := subtotal - req.Diskon

	var id int
	err := r.db.QueryRow(`
		INSERT INTO invoice (order_id, no_invoice, subtotal, diskon, total_bayar, metode_bayar)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id`,
		req.OrderID, noInvoice, subtotal, req.Diskon, totalBayar, req.MetodeBayar,
	).Scan(&id)
	return id, err
}

func (r *Repository) LunaskanInvoice(orderID int) error {
	_, err := r.db.Exec(
		"UPDATE invoice SET status_bayar = 'lunas' WHERE order_id = $1", orderID,
	)
	if err != nil {
		return err
	}
	r.db.Exec("UPDATE service_orders SET status = 'diambil', tanggal_ambil = NOW() WHERE id = $1", orderID)
	return nil
}

// ─── Laporan ──────────────────────────────────────────────────

func (r *Repository) GetLaporanBulanan(tahun string) ([]model.LaporanBulanan, error) {
	var laporan []model.LaporanBulanan
	err := r.db.Select(&laporan, `
		SELECT
		  TO_CHAR(so.tanggal_masuk, 'YYYY-MM')          AS bulan,
		  COUNT(so.id)                                   AS total_order,
		  COUNT(*) FILTER (WHERE so.status = 'diambil') AS total_selesai,
		  COALESCE(SUM(inv.total_bayar), 0)              AS total_pendapat
		FROM service_orders so
		LEFT JOIN invoice inv
		       ON inv.order_id = so.id
		      AND inv.status_bayar = 'lunas'
		WHERE EXTRACT(YEAR FROM so.tanggal_masuk) = $1
		GROUP BY bulan
		ORDER BY bulan ASC`, tahun,
	)
	return laporan, err
}

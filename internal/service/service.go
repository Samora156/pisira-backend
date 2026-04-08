package service

import (
	"fmt"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/pisira/backend/internal/model"
	"github.com/pisira/backend/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	repo      *repository.Repository
	jwtSecret string
	jwtExpire int
}

func New(repo *repository.Repository, jwtSecret string, jwtExpire int) *Service {
	return &Service{repo: repo, jwtSecret: jwtSecret, jwtExpire: jwtExpire}
}

// ─── Auth ─────────────────────────────────────────────────────

func (s *Service) Login(req model.LoginRequest) (*model.LoginResponse, error) {
	user, err := s.repo.GetUserByEmail(req.Email)
	if err != nil {
		log.Printf("[LOGIN DEBUG] User tidak ditemukan untuk email: %s, error: %v", req.Email, err)
		return nil, fmt.Errorf("email atau password salah")
	}

	log.Printf("[LOGIN DEBUG] User ditemukan: %s", user.Email)
	log.Printf("[LOGIN DEBUG] Panjang hash di DB: %d", len(user.Password))
	log.Printf("[LOGIN DEBUG] 7 karakter pertama hash: %s", user.Password[:7])
	log.Printf("[LOGIN DEBUG] Panjang password input: %d", len(req.Password))

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		log.Printf("[LOGIN DEBUG] bcrypt error: %v", err)
		return nil, fmt.Errorf("email atau password salah")
	}

	token, err := s.generateJWT(user)
	if err != nil {
		return nil, err
	}

	log.Printf("[LOGIN DEBUG] Login berhasil untuk: %s", user.Email)
	return &model.LoginResponse{Token: token, User: *user}, nil
}

func (s *Service) generateJWT(user *model.User) (string, error) {
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"role":    user.Role,
		"exp":     time.Now().Add(time.Duration(s.jwtExpire) * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.jwtSecret))
}

// ─── Customer ─────────────────────────────────────────────────

func (s *Service) GetAllCustomers(search string) ([]model.Customer, error) {
	return s.repo.GetAllCustomers(search)
}

func (s *Service) GetCustomerByID(id int) (*model.Customer, error) {
	return s.repo.GetCustomerByID(id)
}

func (s *Service) CreateCustomer(req model.CreateCustomerRequest) (*model.Customer, error) {
	id, err := s.repo.CreateCustomer(req)
	if err != nil {
		return nil, err
	}
	return s.repo.GetCustomerByID(int(id))
}

func (s *Service) UpdateCustomer(id int, req model.CreateCustomerRequest) (*model.Customer, error) {
	if err := s.repo.UpdateCustomer(id, req); err != nil {
		return nil, err
	}
	return s.repo.GetCustomerByID(id)
}

// ─── Service Order ────────────────────────────────────────────

func (s *Service) GetAllOrders(status, search string) ([]model.ServiceOrder, error) {
	return s.repo.GetAllOrders(status, search)
}

func (s *Service) GetOrderByID(id int) (*model.ServiceOrder, error) {
	return s.repo.GetOrderByID(id)
}

func (s *Service) CreateOrder(req model.CreateOrderRequest, teknisiID int) (*model.ServiceOrder, error) {
	if _, err := s.repo.GetCustomerByID(req.CustomerID); err != nil {
		return nil, fmt.Errorf("customer dengan ID %d tidak ditemukan", req.CustomerID)
	}
	id, err := s.repo.CreateOrder(req, teknisiID)
	if err != nil {
		return nil, err
	}
	return s.repo.GetOrderByID(int(id))
}

func (s *Service) UpdateOrderStatus(id int, req model.UpdateOrderStatusRequest) (*model.ServiceOrder, error) {
	validStatus := map[string]bool{
		"menunggu": true, "diagnosa": true, "menunggu_persetujuan": true,
		"proses": true, "selesai": true, "diambil": true, "batal": true,
	}
	if !validStatus[req.Status] {
		return nil, fmt.Errorf("status '%s' tidak valid", req.Status)
	}
	if err := s.repo.UpdateOrderStatus(id, req); err != nil {
		return nil, err
	}
	return s.repo.GetOrderByID(id)
}

// ─── Estimasi ─────────────────────────────────────────────────

func (s *Service) GetEstimasiByOrderID(orderID int) (*model.Estimasi, error) {
	return s.repo.GetEstimasiByOrderID(orderID)
}

func (s *Service) CreateEstimasi(req model.CreateEstimasiRequest) (*model.Estimasi, error) {
	if req.BiayaJasa < 0 || req.BiayaSparepart < 0 {
		return nil, fmt.Errorf("biaya tidak boleh negatif")
	}
	id, err := s.repo.CreateEstimasi(req)
	if err != nil {
		return nil, err
	}
	return s.repo.GetEstimasiByOrderID(int(id))
}

func (s *Service) UpdatePersetujuan(orderID int, status string) error {
	if status != "disetujui" && status != "ditolak" {
		return fmt.Errorf("status persetujuan harus 'disetujui' atau 'ditolak'")
	}
	return s.repo.UpdatePersetujuanEstimasi(orderID, status)
}

// ─── Invoice ──────────────────────────────────────────────────

func (s *Service) GetAllInvoices(statusBayar string) ([]model.Invoice, error) {
	return s.repo.GetAllInvoices(statusBayar)
}

func (s *Service) CreateInvoice(req model.CreateInvoiceRequest) (*model.Invoice, error) {
	if req.Diskon < 0 {
		return nil, fmt.Errorf("diskon tidak boleh negatif")
	}
	id, err := s.repo.CreateInvoice(req)
	if err != nil {
		return nil, err
	}
	var inv model.Invoice
	invoices, _ := s.repo.GetAllInvoices("")
	for _, v := range invoices {
		if v.ID == int(id) {
			inv = v
			break
		}
	}
	return &inv, nil
}

func (s *Service) LunaskanInvoice(orderID int) error {
	return s.repo.LunaskanInvoice(orderID)
}

// ─── Laporan ──────────────────────────────────────────────────

func (s *Service) GetLaporanBulanan(tahun string) ([]model.LaporanBulanan, error) {
	if tahun == "" {
		tahun = time.Now().Format("2006")
	}
	return s.repo.GetLaporanBulanan(tahun)
}

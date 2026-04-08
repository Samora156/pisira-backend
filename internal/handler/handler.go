package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/pisira/backend/internal/model"
	"github.com/pisira/backend/internal/service"
	"github.com/pisira/backend/pkg/response"
)

type Handler struct {
	svc *service.Service
}

func New(svc *service.Service) *Handler {
	return &Handler{svc: svc}
}

// ─── Auth ─────────────────────────────────────────────────────

// POST /api/auth/login
func (h *Handler) Login(c *gin.Context) {
	var req model.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Data tidak valid: "+err.Error())
		return
	}
	result, err := h.svc.Login(req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.OK(c, "Login berhasil", result)
}

// ─── Customer ─────────────────────────────────────────────────

// GET /api/customers?search=budi
func (h *Handler) GetCustomers(c *gin.Context) {
	search := c.Query("search")
	customers, err := h.svc.GetAllCustomers(search)
	if err != nil {
		response.ServerError(c, err)
		return
	}
	response.OK(c, "Data customer berhasil diambil", customers)
}

// GET /api/customers/:id
func (h *Handler) GetCustomerByID(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	customer, err := h.svc.GetCustomerByID(id)
	if err != nil {
		response.NotFound(c, "Customer tidak ditemukan")
		return
	}
	response.OK(c, "Data customer ditemukan", customer)
}

// POST /api/customers
func (h *Handler) CreateCustomer(c *gin.Context) {
	var req model.CreateCustomerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	customer, err := h.svc.CreateCustomer(req)
	if err != nil {
		response.ServerError(c, err)
		return
	}
	response.Created(c, "Customer berhasil ditambahkan", customer)
}

// PUT /api/customers/:id
func (h *Handler) UpdateCustomer(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var req model.CreateCustomerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	customer, err := h.svc.UpdateCustomer(id, req)
	if err != nil {
		response.ServerError(c, err)
		return
	}
	response.OK(c, "Customer berhasil diupdate", customer)
}

// ─── Service Order ────────────────────────────────────────────

// GET /api/orders?status=proses&search=asus
func (h *Handler) GetOrders(c *gin.Context) {
	status := c.Query("status")
	search := c.Query("search")
	orders, err := h.svc.GetAllOrders(status, search)
	if err != nil {
		response.ServerError(c, err)
		return
	}
	response.OK(c, "Data order berhasil diambil", orders)
}

// GET /api/orders/:id
func (h *Handler) GetOrderByID(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	order, err := h.svc.GetOrderByID(id)
	if err != nil {
		response.NotFound(c, "Order tidak ditemukan")
		return
	}
	response.OK(c, "Data order ditemukan", order)
}

// POST /api/orders
func (h *Handler) CreateOrder(c *gin.Context) {
	var req model.CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	teknisiID, _ := c.Get("user_id")
	order, err := h.svc.CreateOrder(req, teknisiID.(int))
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Created(c, "Order servis berhasil dibuat", order)
}

// PATCH /api/orders/:id/status
func (h *Handler) UpdateOrderStatus(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var req model.UpdateOrderStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	order, err := h.svc.UpdateOrderStatus(id, req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.OK(c, "Status order berhasil diupdate", order)
}

// ─── Estimasi ─────────────────────────────────────────────────

// GET /api/orders/:id/estimasi
func (h *Handler) GetEstimasi(c *gin.Context) {
	orderID, _ := strconv.Atoi(c.Param("id"))
	estimasi, err := h.svc.GetEstimasiByOrderID(orderID)
	if err != nil {
		response.NotFound(c, "Estimasi belum dibuat untuk order ini")
		return
	}
	response.OK(c, "Data estimasi ditemukan", estimasi)
}

// POST /api/estimasi
func (h *Handler) CreateEstimasi(c *gin.Context) {
	var req model.CreateEstimasiRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	estimasi, err := h.svc.CreateEstimasi(req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Created(c, "Estimasi berhasil dibuat", estimasi)
}

// PATCH /api/orders/:id/estimasi/persetujuan
func (h *Handler) UpdatePersetujuan(c *gin.Context) {
	orderID, _ := strconv.Atoi(c.Param("id"))
	var body struct {
		Status string `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	if err := h.svc.UpdatePersetujuan(orderID, body.Status); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.OK(c, "Status persetujuan estimasi berhasil diupdate", nil)
}

// ─── Invoice ──────────────────────────────────────────────────

// GET /api/invoices?status_bayar=belum_lunas
func (h *Handler) GetInvoices(c *gin.Context) {
	statusBayar := c.Query("status_bayar")
	invoices, err := h.svc.GetAllInvoices(statusBayar)
	if err != nil {
		response.ServerError(c, err)
		return
	}
	response.OK(c, "Data invoice berhasil diambil", invoices)
}

// POST /api/invoices
func (h *Handler) CreateInvoice(c *gin.Context) {
	var req model.CreateInvoiceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	invoice, err := h.svc.CreateInvoice(req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Created(c, "Invoice berhasil dibuat", invoice)
}

// PATCH /api/invoices/:order_id/lunas
func (h *Handler) LunaskanInvoice(c *gin.Context) {
	orderID, _ := strconv.Atoi(c.Param("order_id"))
	if err := h.svc.LunaskanInvoice(orderID); err != nil {
		response.ServerError(c, err)
		return
	}
	response.OK(c, "Invoice berhasil dilunaskan", nil)
}

// ─── Laporan ──────────────────────────────────────────────────

// GET /api/laporan/bulanan?tahun=2024
func (h *Handler) GetLaporanBulanan(c *gin.Context) {
	tahun := c.Query("tahun")
	laporan, err := h.svc.GetLaporanBulanan(tahun)
	if err != nil {
		response.ServerError(c, err)
		return
	}
	response.OK(c, "Laporan bulanan berhasil diambil", laporan)
}

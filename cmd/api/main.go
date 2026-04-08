package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"github.com/pisira/backend/internal/config"
	"github.com/pisira/backend/internal/handler"
	"github.com/pisira/backend/internal/middleware"
	"github.com/pisira/backend/internal/repository"
	"github.com/pisira/backend/internal/service"
)

func main() {
	// 1. Load konfigurasi dari .env
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Gagal load config: %v", err)
	}

	// 2. Koneksi ke PostgreSQL
	db, err := sqlx.Connect("postgres", cfg.DSN())
	if err != nil {
		log.Fatalf("Gagal koneksi database: %v", err)
	}
	defer db.Close()
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	log.Println("PostgreSQL terkoneksi")

	// 3. Inisialisasi layer (Repository → Service → Handler)
	repo := repository.New(db)
	svc := service.New(repo, cfg.JWTSecret, cfg.JWTExpireHours)
	h := handler.New(svc)
	authMw := middleware.Auth(cfg.JWTSecret)

	// 4. Setup router Gin
	// SetMode HARUS dipanggil sebelum gin.Default() agar tidak ada warning
	if cfg.AppEnv == "production" {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.Default()

	// Fix warning "trusted all proxies": batasi ke lokal saja
	// Jika pakai reverse proxy (Nginx), ganti dengan IP server Nginx Anda
	r.SetTrustedProxies([]string{"127.0.0.1"})

	// CORS — izinkan request dari frontend
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// 5. Daftarkan semua route
	api := r.Group("/api")
	{
		api.POST("/auth/login", h.Login)

		protected := api.Group("", authMw)
		{
			protected.GET("/customers", h.GetCustomers)
			protected.GET("/customers/:id", h.GetCustomerByID)
			protected.POST("/customers", h.CreateCustomer)
			protected.PUT("/customers/:id", h.UpdateCustomer)

			protected.GET("/orders", h.GetOrders)
			protected.GET("/orders/:id", h.GetOrderByID)
			protected.POST("/orders", h.CreateOrder)
			protected.PATCH("/orders/:id/status", h.UpdateOrderStatus)

			protected.GET("/orders/:id/estimasi", h.GetEstimasi)
			protected.POST("/estimasi", h.CreateEstimasi)
			protected.PATCH("/orders/:id/estimasi/persetujuan", h.UpdatePersetujuan)

			protected.GET("/invoices", h.GetInvoices)
			protected.POST("/invoices", h.CreateInvoice)
			protected.PATCH("/invoices/:order_id/lunas", h.LunaskanInvoice)

			admin := protected.Group("", middleware.AdminOnly())
			{
				admin.GET("/laporan/bulanan", h.GetLaporanBulanan)
			}
		}
	}

	// 6. Jalankan server
	addr := fmt.Sprintf(":%s", cfg.AppPort)
	log.Printf("PISIRA Backend berjalan di http://localhost%s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("Server gagal berjalan: %v", err)
	}
}

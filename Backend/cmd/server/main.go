package main

import (
	"backend/internal/handlers"
	"backend/internal/middleware"
	"backend/internal/repositories"
	"backend/internal/services"
	"backend/pkg/config"
	"backend/pkg/database"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Initialize database
	db := database.ConnectDB(cfg)

	// Initialize repositories
	userRepo := repositories.NewUserRepository(db)
	courtRepo := repositories.NewCourtRepository(db)
	reservationRepo := repositories.NewReservationRepository(db)
	paymentRepo := repositories.NewPaymentRepository(db)

	// Initialize services
	authService := services.NewAuthService(userRepo)
	courtService := services.NewCourtService(courtRepo)
	reservationService := services.NewReservationService(reservationRepo, courtRepo)

	// ⚠️ MidtransService TIDAK menerima client eksternal
	midtransService := services.NewMidtransService(paymentRepo, reservationRepo)

	// ❌ Tidak ada PaymentService
	// paymentService := services.NewPaymentService(...)  ← HAPUS

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authService)
	courtHandler := handlers.NewCourtHandler(courtService)
	reservationHandler := handlers.NewReservationHandler(reservationService)

	// PaymentHandler menerima 4 parameter:
	// (midtransService, reservationRepo, userRepo, paymentRepo)
	paymentHandler := handlers.NewPaymentHandler(
		midtransService,
		reservationRepo,
		userRepo,
		paymentRepo,
	)

	// Setup router
	router := setupRouter(
		authHandler,
		courtHandler,
		reservationHandler,
		paymentHandler,
	)

	// Start server
	log.Printf("🚀 Server starting on port %s", cfg.Port)
	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

func setupRouter(
	authHandler *handlers.AuthHandler,
	courtHandler *handlers.CourtHandler,
	reservationHandler *handlers.ReservationHandler,
	paymentHandler *handlers.PaymentHandler,
) *gin.Engine {

	router := gin.Default()

	// Middleware
	router.Use(middleware.CORSMiddleware())
	router.Use(middleware.LoggerMiddleware())

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "OK",
			"message": "Badminton Reservation API is running",
		})
	})

	api := router.Group("/api/v1")

	// Public auth routes
	public := api.Group("/auth")
	{
		public.POST("/register", authHandler.Register)
		public.POST("/login", authHandler.Login)
	}

	// Public courts
	courtRoutes := api.Group("/courts")
	{
		courtRoutes.GET("", courtHandler.GetAllCourts)
		courtRoutes.GET("/available", courtHandler.GetAvailableCourts)
		courtRoutes.POST("/check-availability", courtHandler.CheckAvailability)
		courtRoutes.GET("/:id", courtHandler.GetCourtByID)
	}

	// Protected routes
	protected := api.Group("")
	protected.Use(middleware.AuthMiddleware())
	{
		protected.GET("/auth/profile", authHandler.GetProfile)

		reservationRoutes := protected.Group("/reservations")
		{
			reservationRoutes.POST("", reservationHandler.CreateReservation)
			reservationRoutes.GET("", reservationHandler.GetUserReservations)
			reservationRoutes.GET("/:id", reservationHandler.GetReservationByID)
			reservationRoutes.PUT("/:id/cancel", reservationHandler.CancelReservation)
		}

		paymentRoutes := protected.Group("/payments")
		{
			paymentRoutes.POST("", paymentHandler.CreatePayment)
			paymentRoutes.GET("", paymentHandler.GetUserPayments)
			paymentRoutes.GET("/:id", paymentHandler.GetPaymentByID)
		}
	}

	// Midtrans webhook (public)
	api.POST("/payments/notification", paymentHandler.HandlePaymentNotification)

	return router
}

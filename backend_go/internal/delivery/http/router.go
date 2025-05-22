package http

import (
	"github.com/gin-gonic/gin"
	"github.com/RenanGroh/SaaS_microbusiness_project/backend_go/internal/config"
	"github.com/RenanGroh/SaaS_microbusiness_project/backend_go/internal/delivery/http/middleware"
)

// Modifique a assinatura para incluir AppointmentHandler
func SetupRoutes(
	router *gin.Engine,
	cfg *config.Config,
	userHandler *UserHandler,
	appointmentHandler *AppointmentHandler, // <<< ADICIONADO
) {
	authMW := middleware.AuthMiddleware(cfg)

	apiV1 := router.Group("/api/v1")
	{
		// Rotas de Autenticação
		authRoutes := apiV1.Group("/auth")
		{
			authRoutes.POST("/login", userHandler.Login)
			// authRoutes.POST("/register", userHandler.CreateUser) // Opcional
		}

		// Rotas de Usuário
		userRoutes := apiV1.Group("/users")
		{
			userRoutes.POST("", userHandler.CreateUser) // Criação de usuário geralmente não precisa de auth
			userRoutes.GET("/by-email", userHandler.GetUserByEmail) // Pode ser pública ou protegida
			userRoutes.GET("/me", authMW, userHandler.GetUserProfile)
			userRoutes.GET("/:id", authMW, userHandler.GetUserByID) // Rota para buscar usuário por ID
		}

		// Rotas de Agendamento (todas protegidas)
		appointmentRoutes := apiV1.Group("/appointments")
		appointmentRoutes.Use(authMW) // Aplica o middleware de autenticação a todas as rotas de agendamento
		{
			appointmentRoutes.POST("", appointmentHandler.CreateAppointment)
			appointmentRoutes.GET("", appointmentHandler.ListUserAppointments)
			appointmentRoutes.GET("/:id", appointmentHandler.GetAppointmentByID)
			appointmentRoutes.PUT("/:id", appointmentHandler.UpdateAppointment)
			appointmentRoutes.PATCH("/:id/cancel", appointmentHandler.CancelAppointment) // Usando PATCH para mudança de status
			// appointmentRoutes.DELETE("/:id", appointmentHandler.DeleteAppointment) // Se for implementar delete
		}
	}

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "UP"})
	})
}
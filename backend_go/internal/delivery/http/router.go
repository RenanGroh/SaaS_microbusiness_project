package http

import (
	"github.com/gin-gonic/gin"
	"github.com/RenanGroh/SaaS_microbusiness_project/backend_go/internal/config" // Importar config
	"github.com/RenanGroh/SaaS_microbusiness_project/backend_go/internal/delivery/http/middleware" // Importar middleware
)

func SetupRoutes(router *gin.Engine, cfg *config.Config, userHandler *UserHandler) {

	authMW := middleware.AuthMiddleware(cfg)
	
	apiV1 := router.Group("/api/v1")
	{
		// Rotas de Autenticação
		authRoutes := apiV1.Group("/auth")
		{
			authRoutes.POST("/login", userHandler.Login)
		}

		// Rotas de Usuário
		userRoutes := apiV1.Group("/users")
		{
			userRoutes.POST("", userHandler.CreateUser)
			userRoutes.GET("/by-email", userHandler.GetUserByEmail) // Pode ou não ser protegida
			userRoutes.GET("/:id", authMW, userHandler.GetUserByID)
			// Exemplo de rota protegida: buscar o perfil do usuário logado
			userRoutes.GET("/me", authMW, userHandler.GetUserProfile) // <<< ROTA PROTEGIDA
		}
	}

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "UP"})
	})
}
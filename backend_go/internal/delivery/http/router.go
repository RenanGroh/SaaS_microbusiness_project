// Em internal/delivery/http/router.go
package http

import (
	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine, userHandler *UserHandler) {

	apiV1 := router.Group("/api/v1")
	{
		// Rotas de Autenticação
		authRoutes := apiV1.Group("/auth")
		{
			authRoutes.POST("/login", userHandler.Login) // <<< NOVA ROTA
			// authRoutes.POST("/register", userHandler.CreateUser) // Poderia mover o registro para cá também
		}

		// Rotas de Usuário
		userRoutes := apiV1.Group("/users")
		{
			userRoutes.POST("", userHandler.CreateUser) // Ou use /auth/register
			userRoutes.GET("/by-email", userHandler.GetUserByEmail)
			// userRoutes.GET("/:id", authMiddleware.RequireAuth, userHandler.GetUserByID) // Comentado
		}
	}

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "UP"})
	})
}
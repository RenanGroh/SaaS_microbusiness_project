package main

import (
	"log"

	"github.com/RenanGroh/SaaS_microbusiness_project/backend_go/internal/config"
	httpDelivery "github.com/RenanGroh/SaaS_microbusiness_project/backend_go/internal/delivery/http"
	gormPersistence "github.com/RenanGroh/SaaS_microbusiness_project/backend_go/internal/infra/persistence/gorm"
	"github.com/RenanGroh/SaaS_microbusiness_project/backend_go/internal/usecase"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.LoadConfig()

	db, err := gormPersistence.NewGormDB(cfg.DBDriver, cfg.DBSource)
	if err != nil {
		log.Fatalf("Falha ao conectar ao banco de dados: %v", err)
	}
	err = db.AutoMigrate(&gormPersistence.UserGormModel{})
	if err != nil {
		log.Fatalf("Falha ao rodar AutoMigrate: %v", err)
	}

	userGormRepo := gormPersistence.NewGormUserRepository(db)

	// Agora passamos jwtSecret e jwtExpirationHours para o UserUseCase
	userUC := usecase.NewUserUseCase(
		userGormRepo,
		cfg.JWTSecret,          // <<< ADICIONADO
		cfg.JWTExpirationHours, // <<< ADICIONADO
	)

	userHandler := httpDelivery.NewUserHandler(userUC)

	router := gin.Default()
	httpDelivery.SetupRoutes(router, cfg, userHandler)

	log.Printf("Servidor Bizly iniciando na porta %s", cfg.ServerPort)
	if err := router.Run(":" + cfg.ServerPort); err != nil {
		log.Fatalf("Falha ao iniciar o servidor: %v", err)
	}
}
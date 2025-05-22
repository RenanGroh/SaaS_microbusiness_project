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

	// Adicionar AppointmentGormModel ao AutoMigrate
	err = db.AutoMigrate(&gormPersistence.UserGormModel{}, &gormPersistence.AppointmentGormModel{}) // <<< ADICIONADO
	if err != nil {
		log.Fatalf("Falha ao rodar AutoMigrate: %v", err)
	}

	// Repositórios
	userGormRepo := gormPersistence.NewGormUserRepository(db)
	appointmentGormRepo := gormPersistence.NewGormAppointmentRepository(db) // <<< ADICIONADO

	// Casos de Uso
	userUC := usecase.NewUserUseCase(userGormRepo, cfg.JWTSecret, cfg.JWTExpirationHours)
	appointmentUC := usecase.NewAppointmentUseCase(appointmentGormRepo, userGormRepo) // <<< ADICIONADO (passa userRepo também)

	// Handlers HTTP
	userHandler := httpDelivery.NewUserHandler(userUC)
	appointmentHandler := httpDelivery.NewAppointmentHandler(appointmentUC) // <<< ADICIONADO

	// Roteador Gin
	router := gin.Default()
	// Passar appointmentHandler para SetupRoutes
	httpDelivery.SetupRoutes(router, cfg, userHandler, appointmentHandler) // <<< ADICIONADO

	log.Printf("Servidor Bizly iniciando na porta %s", cfg.ServerPort)
	if err := router.Run(":" + cfg.ServerPort); err != nil {
		log.Fatalf("Falha ao iniciar o servidor: %v", err)
	}
}
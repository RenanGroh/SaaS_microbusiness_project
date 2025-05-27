package main

import (
	"log"

	"github.com/RenanGroh/SaaS_microbusiness_project/backend_go/internal/config"
	httpDelivery "github.com/RenanGroh/SaaS_microbusiness_project/backend_go/internal/delivery/http"
	gormPersistence "github.com/RenanGroh/SaaS_microbusiness_project/backend_go/internal/infra/persistence/gorm"
	"github.com/RenanGroh/SaaS_microbusiness_project/backend_go/internal/usecase"

	"github.com/gin-contrib/cors" // <<< ADICIONE ESTE IMPORT
	"github.com/gin-gonic/gin"
	"time" // <<< ADICIONE ESTE IMPORT (se usar cors.Config com MaxAge)
)

func main() {
	cfg := config.LoadConfig()

	db, err := gormPersistence.NewGormDB(cfg.DBDriver, cfg.DBSource)
	if err != nil {
		log.Fatalf("Falha ao conectar ao banco de dados: %v", err)
	}
	err = db.AutoMigrate(&gormPersistence.UserGormModel{}, &gormPersistence.AppointmentGormModel{})
	if err != nil {
		log.Fatalf("Falha ao rodar AutoMigrate: %v", err)
	}

	userGormRepo := gormPersistence.NewGormUserRepository(db)
	appointmentGormRepo := gormPersistence.NewGormAppointmentRepository(db)

	userUC := usecase.NewUserUseCase(userGormRepo, cfg.JWTSecret, cfg.JWTExpirationHours)
	appointmentUC := usecase.NewAppointmentUseCase(appointmentGormRepo, userGormRepo)

	userHandler := httpDelivery.NewUserHandler(userUC)
	appointmentHandler := httpDelivery.NewAppointmentHandler(appointmentUC)

	// gin.SetMode(gin.ReleaseMode) // Descomente para produção
	router := gin.Default() // gin.Default() já inclui logger e recovery

	// --- CONFIGURAÇÃO DO CORS ---
	// Use cors.New com uma configuração customizada
	corsConfig := cors.Config{
		// AllowOrigins especifica quais origens são permitidas.
		// Para desenvolvimento, você pode usar a porta específica do seu Flutter Web.
		// Em produção, você colocaria o domínio real do seu frontend.
		AllowOrigins: []string{"http://localhost:59202", "http://127.0.0.1:59202"}, // Adicione a porta que o Flutter Web está usando
		// Você pode adicionar mais origens se necessário, ex: "http://localhost:3000" para um React app
		// Para permitir todas as origens (NÃO RECOMENDADO PARA PRODUÇÃO):
		// AllowAllOrigins: true,

		AllowMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders: []string{"Origin", "Content-Type", "Accept", "Authorization"},
		// ExposeHeaders permite que o cliente acesse certos cabeçalhos da resposta
		// ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true, // Se você precisar enviar cookies ou cabeçalhos de autenticação
		MaxAge:           12 * time.Hour,
	}
	router.Use(cors.New(corsConfig))
	// --- FIM DA CONFIGURAÇÃO DO CORS ---

	httpDelivery.SetupRoutes(router, cfg, userHandler, appointmentHandler)

	log.Printf("Servidor Bizly iniciando na porta %s", cfg.ServerPort)
	if err := router.Run(":" + cfg.ServerPort); err != nil {
		log.Fatalf("Falha ao iniciar o servidor: %v", err)
	}
}
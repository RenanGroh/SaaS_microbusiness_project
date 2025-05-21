package main

import (
	"log"
	// Removido "os" se não for usado diretamente aqui (agora está em config)

	// Ajuste os paths de import para corresponder ao seu nome de módulo no go.mod
	// Exemplo: se go.mod diz "module github.com/seu-usuario/bizly"
	"github.com/RenanGroh/SaaS_microbusiness_project/backend_go/internal/config"
	httpDelivery "github.com/RenanGroh/SaaS_microbusiness_project/backend_go/internal/delivery/http" // Alias para evitar conflito com pacote http padrão
	gormPersistence "github.com/RenanGroh/SaaS_microbusiness_project/backend_go/internal/infra/persistence/gorm"
	"github.com/RenanGroh/SaaS_microbusiness_project/backend_go/internal/usecase"

	"github.com/gin-gonic/gin"
	// "github.com/joho/godotenv" // Movido para dentro do pacote config
)

func main() {
	// godotenv.Load() foi movido para dentro de config.LoadConfig()

	cfg := config.LoadConfig() // Carrega configurações (DB DSN, Port, JWT Secret, etc.)

	log.Printf("Rodando na porta: %s", cfg.ServerPort)
	log.Printf("Usando segredo JWT: %s (Apenas para debug, não logue em produção!)", cfg.JWTSecret)

	// Inicializar Banco de Dados (GORM)
	db, err := gormPersistence.NewGormDB(cfg.DBDriver, cfg.DBSource)
	if err != nil {
		log.Fatalf("Falha ao conectar ao banco de dados: %v", err)
	}

	// AutoMigrate (ou rodar migrations dedicadas)
	// Idealmente, AutoMigrate só para UserGormModel e outras GORM models.
	// O pacote gormPersistence agora exporta os modelos GORM se necessário aqui,
	// ou você pode passar o *gorm.DB para uma função de migração dentro do pacote gormPersistence.
	err = db.AutoMigrate(&gormPersistence.UserGormModel{} /*, &gormPersistence.AppointmentGormModel{} */)
	if err != nil {
		log.Fatalf("Falha ao rodar AutoMigrate: %v", err)
	}

	// Inicializar Repositórios
	// A interface é definida em internal/repository, a implementação em internal/infra/...
	userGormRepo := gormPersistence.NewGormUserRepository(db)
	// appointmentGormRepo := gormPersistence.NewGormAppointmentRepository(db) // Quando tiver

	// Inicializar Casos de Uso
	// passwordHasher := security.NewBcryptPasswordHasher() // Se for usar a interface
	// tokenGenerator := security.NewJWTTokenGenerator(cfg.JWTSecret) // Se for usar a interface
	userUC := usecase.NewUserUseCase(userGormRepo /*, passwordHasher, tokenGenerator */)
	// appointmentUC := usecase.NewAppointmentUseCase(appointmentGormRepo) // Quando tiver

	// Inicializar Handlers HTTP
	userHandler := httpDelivery.NewUserHandler(userUC)
	// appointmentHandler := httpDelivery.NewAppointmentHandler(appointmentUC) // Quando tiver
	// authMiddleware := middleware.NewAuthMiddleware(tokenGenerator) // Quando tiver JWT

	// Configurar Roteador Gin
	// gin.SetMode(gin.ReleaseMode) // Para produção
	router := gin.Default() // gin.Default() já inclui logger e recovery

	// Setup das rotas usando a função do pacote httpDelivery
	httpDelivery.SetupRoutes(router, userHandler /*, appointmentHandler, authMiddleware */)

	log.Printf("Servidor Bizly iniciando na porta %s", cfg.ServerPort)
	if err := router.Run(":" + cfg.ServerPort); err != nil {
		log.Fatalf("Falha ao iniciar o servidor: %v", err)
	}
}
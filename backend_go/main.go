package main

import (
	"fmt"
	"log"
	"net/http"
	"os" // Para variáveis de ambiente (melhor prática)

	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger" // Para configurar o logger do GORM
)

// Defina suas structs de modelo aqui (ex: User, Appointment)
type User struct {
    gorm.Model // Inclui ID, CreatedAt, UpdatedAt, DeletedAt
    Name       string `gorm:"size:100;not null"`
    Email      string `gorm:"size:100;uniqueIndex;not null"`
    Password   string `gorm:"not null"` // Lembre-se de hashear!
    // Outros campos
}

// Variável global para o DB (ou injete via contexto/structs de handler)
var DB *gorm.DB

func connectDatabase() {
    dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=America/Sao_Paulo",
        os.Getenv("DB_HOST"),       // Ex: "localhost"
        os.Getenv("DB_USER"),       // Ex: "postgres"
        os.Getenv("DB_PASSWORD"),   // Ex: "sua_senha"
        os.Getenv("DB_NAME"),       // Ex: "meu_saas_db"
        os.Getenv("DB_PORT"),       // Ex: "5432"
    )

    // Configurar logger do GORM
    newLogger := logger.New(
        log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
        logger.Config{
            SlowThreshold:             200 * time.Millisecond, // Limite para SQL lento
            LogLevel:                  logger.Info,                 // Nível de Log (Info, Warn, Error)
            IgnoreRecordNotFoundError: true,                        // Ignorar erros de 'record not found'
            Colorful:                  true,                        // Saída colorida
        },
    )

    database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
        Logger: newLogger,
    })
    if err != nil {
        log.Fatal("Falha ao conectar ao banco de dados!", err)
    }

    // Auto Migrate (cria/altera tabelas baseado nas structs)
    // CUIDADO: Use com cautela em produção. Prefira migrations dedicadas.
    err = database.AutoMigrate(&User{}) // Adicione outros modelos aqui: &Appointment{}, etc.
    if err != nil {
        log.Fatal("Falha ao realizar auto-migrate!", err)
    }

    DB = database
    log.Println("Conexão com o banco de dados estabelecida e migrações (auto) completas.")
}

func main() {
    // Carregar variáveis de ambiente (ex: usando godotenv para dev)
    errEnv := godotenv.Load()
    if errEnv != nil {
		log.Println("Aviso: Não foi possível carregar o arquivo .env. Usando variáveis de ambiente do sistema se disponíveis.")
	}
    connectDatabase() // Conecta ao banco

    router := gin.Default()

    // Rota de teste para criar um usuário (SIMPLIFICADO, sem validação, sem hash de senha ainda)
    router.POST("/users", func(c *gin.Context) {
        var user User
        if err := c.ShouldBindJSON(&user); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }

        // AQUI você adicionaria hashing de senha antes de salvar
        // result := DB.Create(&user)
        // if result.Error != nil {
        //     c.JSON(http.StatusInternalServerError, gin.H{"error": "Falha ao criar usuário"})
        //     return
        // }

        // Placeholder para o exemplo
        log.Printf("Usuário recebido: %+v\n", user)
        user.Password = "SENHA_HASHEADA_AQUI" // SIMULAÇÃO
        result := DB.Create(&user) // Cria o usuário no banco
        if result.Error != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
            return
        }

        c.JSON(http.StatusCreated, user)
    })

    router.GET("/users", func(c *gin.Context){
        var users []User
        result := DB.Find(&users)
        if result.Error != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Falha ao buscar usuários"})
            return
        }
        c.JSON(http.StatusOK, users)
    })

    router.GET("/api/test/hello", func(c *gin.Context) {
        c.JSON(http.StatusOK, gin.H{
            "message": "Olá do Backend Go com Gin e GORM!",
        })
    })

    log.Println("Servidor iniciando na porta 8080...")
    err := router.Run(":8080")
    if err != nil {
        log.Fatal("Falha ao iniciar o servidor Gin: ", err)
    }
}
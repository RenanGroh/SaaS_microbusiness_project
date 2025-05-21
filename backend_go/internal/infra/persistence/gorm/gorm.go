package gorm

import (
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// NewGormDB inicializa e retorna uma conexão com o banco de dados GORM.
func NewGormDB(dbDriver, dbSource string) (*gorm.DB, error) {
	var dialector gorm.Dialector
	if dbDriver == "postgres" {
		dialector = postgres.Open(dbSource)
	} else {
		return nil, fmt.Errorf("driver de banco de dados não suportado: %s", dbDriver)
	}

	gormLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             200 * time.Millisecond,
			LogLevel:                  logger.Info, // Mude para logger.Warn em produção se quiser menos logs
			IgnoreRecordNotFoundError: true,
			Colorful:                  true,
		},
	)

	db, err := gorm.Open(dialector, &gorm.Config{
		Logger: gormLogger,
		PrepareStmt: true, // Habilita prepared statements para melhor performance
	})

	if err != nil {
		return nil, fmt.Errorf("falha ao conectar ao banco de dados: %w", err)
	}

	// Opcional: Configurar pool de conexões
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("falha ao obter o objeto sql.DB subjacente: %w", err)
	}
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	log.Println("Conexão com o banco de dados GORM estabelecida.")
	return db, nil
}
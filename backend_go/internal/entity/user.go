package entity

import (
	"time"
	//"github.com/google/uuid" // Se for usar UUID como ID
)

// User representa a entidade de usuário no domínio da aplicação.
// Não contém tags de banco de dados ou JSON.
type User struct {
	ID        uint      // GORM Model usa uint por padrão para ID
	Name      string
	Email     string
	Password  string // Este será o HASH da senha
	CreatedAt time.Time
	UpdatedAt time.Time
	// DeletedAt gorm.DeletedAt `gorm:"index"` // Se você precisar de soft delete, mas gorm.Model já inclui
}

// NewUser é um construtor para a entidade User (opcional, mas bom para validações iniciais)
 func NewUser(name, email, hashedPassword string) (*User, error) {
//  // Adicionar validações aqui se necessário antes de criar a struct
 	return &User{
 		Name:     name,
 		Email:    email,
 		Password: hashedPassword,
	}, nil
}
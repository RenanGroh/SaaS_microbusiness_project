package entity

import (
	"time"

	"github.com/google/uuid" // Se for usar UUID como ID
)

// User representa a entidade de usuário no domínio da aplicação.
// Não contém tags de banco de dados ou JSON.
type User struct {
	ID        uuid.UUID `json:"id"` // Mude para uuid.UUID
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Password  string    // Não exponha no JSON
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// Construtor opcional, para gerar o UUID na criação da entidade
func NewUser(name, email, hashedPassword string) (*User, error) {
    // Validar entradas...
    return &User{
        ID:        uuid.New(), // Gera um UUID v4
        Name:      name,
        Email:     email,
        Password:  hashedPassword,
        // CreatedAt e UpdatedAt podem ser definidos aqui ou pelo GORM/repo
    }, nil
}
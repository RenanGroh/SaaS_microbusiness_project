package repository

import "github.com/RenanGroh/SaaS_microbusiness_project/backend_go/internal/entity" // <<< AJUSTE O PATH DO MÓDULO AQUI

// UserRepository define a interface para interações com o armazenamento de dados de usuários.
// Os casos de uso dependem desta interface, não da implementação concreta.
type UserRepository interface {
	Create(user *entity.User) error
	FindByEmail(email string) (*entity.User, error)
	FindByID(id uint) (*entity.User, error)
	// GetAll() ([]*entity.User, error) // Exemplo de outro método
}
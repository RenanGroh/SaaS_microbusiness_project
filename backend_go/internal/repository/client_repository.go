package repository

import (
	"github.com/google/uuid"
	"github.com/RenanGroh/SaaS_microbusiness_project/backend_go/internal/entity"
)

// ClientRepository define a interface para interações com o armazenamento de dados de clientes.
type ClientRepository interface {
	Create(client *entity.Client) error
	FindByID(id uuid.UUID) (*entity.Client, error)
	FindByUserID(userID uuid.UUID) ([]*entity.Client, error)
	Update(client *entity.Client) error
	Delete(id uuid.UUID) error
}

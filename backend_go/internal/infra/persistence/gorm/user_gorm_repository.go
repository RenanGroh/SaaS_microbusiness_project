package gorm

import (
	"errors"
	"github.com/RenanGroh/SaaS_microbusiness_project/backend_go/internal/entity"     // <<< AJUSTE O PATH DO MÓDULO AQUI
	"github.com/RenanGroh/SaaS_microbusiness_project/backend_go/internal/repository" // <<< AJUSTE O PATH DO MÓDULO AQUI
	"gorm.io/gorm"
)

type gormUserRepository struct {
	db *gorm.DB
}

// NewGormUserRepository cria uma nova instância da implementação do repositório de usuário com GORM.
// Ele retorna a interface repository.UserRepository.
func NewGormUserRepository(db *gorm.DB) repository.UserRepository {
	return &gormUserRepository{db: db}
}

func (r *gormUserRepository) Create(userEntity *entity.User) error {
	userGorm := FromEntity(userEntity) // Converte entidade para modelo GORM
	result := r.db.Create(userGorm)
	if result.Error != nil {
		// Poderia verificar aqui se o erro é de "duplicate key" para email
		// e retornar um erro mais específico se necessário, ex: repository.ErrEmailExists
		return result.Error
	}
	// Atualizar a entidade original com ID e Timestamps gerados pelo GORM
	userEntity.ID = userGorm.ID
	userEntity.CreatedAt = userGorm.CreatedAt
	userEntity.UpdatedAt = userGorm.UpdatedAt
	return nil
}

func (r *gormUserRepository) FindByEmail(email string) (*entity.User, error) {
	var userGorm UserGormModel
	result := r.db.Where("email = ?", email).First(&userGorm)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil // Indica que não foi encontrado, sem erro de sistema
		}
		return nil, result.Error // Outro erro do GORM
	}
	return userGorm.ToEntity(), nil
}

func (r *gormUserRepository) FindByID(id uint) (*entity.User, error) {
	var userGorm UserGormModel
	// GORM busca por chave primária automaticamente se o segundo argumento for o valor da PK
	result := r.db.First(&userGorm, id)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil // Indica que não foi encontrado
		}
		return nil, result.Error
	}
	return userGorm.ToEntity(), nil
}
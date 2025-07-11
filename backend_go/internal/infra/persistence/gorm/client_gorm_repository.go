package gorm

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/RenanGroh/SaaS_microbusiness_project/backend_go/internal/entity"
	"github.com/RenanGroh/SaaS_microbusiness_project/backend_go/internal/repository"
	"gorm.io/gorm"
)

// ClientGormModel representa o modelo de cliente para o GORM.
type ClientGormModel struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	UserID    uuid.UUID `gorm:"type:uuid;not null"`
	Name      string    `gorm:"type:varchar(255);not null"`
	Email     string    `gorm:"type:varchar(255);uniqueIndex"`
	Phone     string    `gorm:"type:varchar(50)"`
	Notes     string    `gorm:"type:text"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"index"` // Para soft delete

	// Relacionamento com User (opcional, para carregar o usuário dono do cliente)
	User UserGormModel `gorm:"foreignKey:UserID"`
}

// TableName define o nome da tabela no banco de dados.
func (ClientGormModel) TableName() string {
	return "clients"
}

// ToEntity converte um ClientGormModel para uma entidade Client.
func (c *ClientGormModel) ToEntity() *entity.Client {
	return &entity.Client{
		ID:        c.ID,
		UserID:    c.UserID,
		Name:      c.Name,
		Email:     c.Email,
		Phone:     c.Phone,
		Notes:     c.Notes,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
	}
}

// ClientFromEntity converte uma entidade Client para ClientGormModel.
func ClientFromEntity(e *entity.Client) *ClientGormModel {
	return &ClientGormModel{
		ID:        e.ID,
		UserID:    e.UserID,
		Name:      e.Name,
		Email:     e.Email,
		Phone:     e.Phone,
		Notes:     e.Notes,
		CreatedAt: e.CreatedAt,
		UpdatedAt: e.UpdatedAt,
	}
}

// gormClientRepository implementa a interface ClientRepository usando GORM.
type gormClientRepository struct {
	db *gorm.DB
}

// NewGormClientRepository cria uma nova instância de GormClientRepository.
func NewGormClientRepository(db *gorm.DB) repository.ClientRepository {
	return &gormClientRepository{db: db}
}

// Create cria um novo cliente no banco de dados.
func (r *gormClientRepository) Create(clientEntity *entity.Client) error {
	clientGorm := ClientFromEntity(clientEntity)
	result := r.db.Create(clientGorm)
	if result.Error != nil {
		return result.Error
	}
	// Atualizar a entidade original com ID e Timestamps gerados pelo GORM
	clientEntity.ID = clientGorm.ID
	clientEntity.CreatedAt = clientGorm.CreatedAt
	clientEntity.UpdatedAt = clientGorm.UpdatedAt
	return nil
}

// FindByID busca um cliente pelo seu ID.
func (r *gormClientRepository) FindByID(id uuid.UUID) (*entity.Client, error) {
	var clientGorm ClientGormModel
	result := r.db.First(&clientGorm, "id = ?", id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil // Indica não encontrado
		}
		return nil, result.Error
	}
	return clientGorm.ToEntity(), nil
}

// FindByUserID busca todos os clientes associados a um UserID específico.
func (r *gormClientRepository) FindByUserID(userID uuid.UUID) ([]*entity.Client, error) {
	var clientsGorm []ClientGormModel
	result := r.db.Where("user_id = ?", userID).Find(&clientsGorm)
	if result.Error != nil {
		return nil, result.Error
	}

	var clientEntities []*entity.Client
	for _, cg := range clientsGorm {
		clientEntities = append(clientEntities, cg.ToEntity())
	}
	return clientEntities, nil
}

// Update atualiza um cliente existente no banco de dados.
func (r *gormClientRepository) Update(clientEntity *entity.Client) error {
	if clientEntity.ID == uuid.Nil {
		return errors.New("ID do cliente não pode ser nulo para atualização")
	}
	clientGorm := ClientFromEntity(clientEntity)
	result := r.db.Model(&ClientGormModel{}).Where("id = ?", clientGorm.ID).Updates(clientGorm)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("cliente não encontrado para atualização ou nenhum dado alterado")
	}
	return nil
}

// Delete exclui um cliente do banco de dados (soft delete).
func (r *gormClientRepository) Delete(id uuid.UUID) error {
	if id == uuid.Nil {
		return errors.New("ID do cliente não pode ser nulo para deleção")
	}
	result := r.db.Delete(&ClientGormModel{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("cliente não encontrado para deleção")
	}
	return nil
}

package gorm

import (
	"time"

	"github.com/RenanGroh/SaaS_microbusiness_project/backend_go/internal/entity" // <<< AJUSTE O PATH DO MÓDULO AQUI
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// UserGormModel representa a estrutura do usuário para o GORM, com tags de banco.
type UserGormModel struct {
	// gorm.Model // NÃO USE gorm.Model se for usar UUID como PK, pois gorm.Model usa ID uint
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key"` // Define como PK
	Name      string    `gorm:"size:100;not null"`
	Email     string    `gorm:"size:100;uniqueIndex;not null"`
	Password  string    `gorm:"not null"`
	CreatedAt time.Time // GORM vai popular automaticamente
	UpdatedAt time.Time // GORM vai popular automaticamente
    DeletedAt gorm.DeletedAt `gorm:"index"` // Para soft delete, se precisar
}

// ToEntity converte um UserGormModel para uma entity.User.
func (m *UserGormModel) ToEntity() *entity.User {
	return &entity.User{
		ID:        m.ID,
		Name:      m.Name,
		Email:     m.Email,
		Password:  m.Password,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}

// FromEntity converte uma entity.User para um UserGormModel para persistência.

func FromEntity(e *entity.User) *UserGormModel {
    // Se o ID da entidade for uuid.Nil, o GORM (com default:uuid_generate_v4()) irá gerar.
    // Se já tiver um ID, ele será usado.
	return &UserGormModel{
		ID:        e.ID,
		Name:      e.Name,
		Email:     e.Email,
		Password:  e.Password,
		CreatedAt: e.CreatedAt, // GORM pode sobrescrever se for valor zero
		UpdatedAt: e.UpdatedAt, // GORM pode sobrescrever se for valor zero
	}
}

// AppointmentGormModel (Exemplo para quando você adicionar agendamentos)
 type AppointmentGormModel struct {
 	gorm.Model
 	UserID    uint   `gorm:"not null"`
	User      UserGormModel `gorm:"foreignKey:UserID"` // Relacionamento
	ClientName string
	StartTime time.Time
	EndTime   time.Time
	Status    string
}

//func (m *AppointmentGormModel) ToEntity() *entity.Appointment { /* ... */ }
//func AppointmentFromEntity(e *entity.Appointment) *AppointmentGormModel { /* ... */ }
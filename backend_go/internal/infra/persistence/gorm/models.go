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

// -----------------------------------------------------------------------------
// AppointmentGormModel
// -----------------------------------------------------------------------------
type AppointmentGormModel struct {
	ID                uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	UserID            uuid.UUID `gorm:"type:uuid;not null;index"` // Chave estrangeira para UserGormModel
	User              UserGormModel `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"` // Relacionamento
	ClientID          *uuid.UUID `gorm:"type:uuid;index"` // Opcional, pode ser nulo
	ClientUser        *UserGormModel `gorm:"foreignKey:ClientID;constraint:OnUpdate:SET NULL,OnDelete:SET NULL;"` // Relacionamento opcional
	ClientName        string    `gorm:"size:255"`
	ClientEmail       string    `gorm:"size:255"`
	ClientPhone       string    `gorm:"size:50"`
	ServiceDescription string    `gorm:"type:text"`
	StartTime         time.Time `gorm:"not null;index"`
	EndTime           time.Time `gorm:"not null"`
	Status            string    `gorm:"size:50;not null;default:'PENDING'"` // Usando string para status no GORM
	Notes             string    `gorm:"type:text"`
	Price             float64
	CreatedAt         time.Time
	UpdatedAt         time.Time
	DeletedAt         gorm.DeletedAt `gorm:"index"`
}

// ToEntity converte um AppointmentGormModel para uma entity.Appointment.
func (m *AppointmentGormModel) ToEntity() *entity.Appointment {
	return &entity.Appointment{
		ID:                m.ID,
		UserID:            m.UserID,
		ClientID:          m.ClientID, // Preserva o ponteiro
		ClientName:        m.ClientName,
		ClientEmail:       m.ClientEmail,
		ClientPhone:       m.ClientPhone,
		ServiceDescription: m.ServiceDescription,
		StartTime:         m.StartTime,
		EndTime:           m.EndTime,
		Status:            entity.AppointmentStatus(m.Status), // Converte string para o tipo customizado
		Notes:             m.Notes,
		Price:             m.Price,
		CreatedAt:         m.CreatedAt,
		UpdatedAt:         m.UpdatedAt,
	}
}

// AppointmentFromEntity converte uma entity.Appointment para um AppointmentGormModel para persistência.
func AppointmentFromEntity(e *entity.Appointment) *AppointmentGormModel {
	return &AppointmentGormModel{
		ID:                e.ID, // Se e.ID for uuid.Nil, GORM (com default) irá gerar
		UserID:            e.UserID,
		ClientID:          e.ClientID,
		ClientName:        e.ClientName,
		ClientEmail:       e.ClientEmail,
		ClientPhone:       e.ClientPhone,
		ServiceDescription: e.ServiceDescription,
		StartTime:         e.StartTime,
		EndTime:           e.EndTime,
		Status:            string(e.Status), // Converte tipo customizado para string
		Notes:             e.Notes,
		Price:             e.Price,
		CreatedAt:         e.CreatedAt, // GORM pode popular se for zero
		UpdatedAt:         e.UpdatedAt, // GORM pode popular se for zero
	}
}
//func (m *AppointmentGormModel) ToEntity() *entity.Appointment { /* ... */ }
//func AppointmentFromEntity(e *entity.Appointment) *AppointmentGormModel { /* ... */ }
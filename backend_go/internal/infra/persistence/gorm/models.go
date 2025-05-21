package gorm

import (
	"time"

	"github.com/RenanGroh/SaaS_microbusiness_project/backend_go/internal/entity" // <<< AJUSTE O PATH DO MÓDULO AQUI
	"gorm.io/gorm"
)

// UserGormModel representa a estrutura do usuário para o GORM, com tags de banco.
type UserGormModel struct {
	gorm.Model        // Inclui ID uint, CreatedAt, UpdatedAt, DeletedAt
	Name       string `gorm:"size:100;not null"`
	Email      string `gorm:"size:100;uniqueIndex;not null"` // uniqueIndex garante unicidade no banco
	Password   string `gorm:"not null"`                    // Hash da senha
}

// ToEntity converte um UserGormModel para uma entity.User.
func (m *UserGormModel) ToEntity() *entity.User {
	return &entity.User{
		ID:        m.ID,
		Name:      m.Name,
		Email:     m.Email,
		Password:  m.Password, // Passa o hash
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}

// FromEntity converte uma entity.User para um UserGormModel para persistência.
func FromEntity(e *entity.User) *UserGormModel {
	return &UserGormModel{
		Model: gorm.Model{ // Se ID for 0, GORM gera. Se não, tenta usar o ID fornecido.
			ID:        e.ID,
			CreatedAt: e.CreatedAt,
			UpdatedAt: e.UpdatedAt,
		},
		Name:     e.Name,
		Email:    e.Email,
		Password: e.Password, // Salva o hash
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
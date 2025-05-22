package repository

import (
	"time"
	"github.com/google/uuid"
	"github.com/RenanGroh/SaaS_microbusiness_project/backend_go/internal/entity"
)

// AppointmentRepository define a interface para interações com o armazenamento de dados de agendamentos.
type AppointmentRepository interface {
	Create(appointment *entity.Appointment) error
	FindByID(id uuid.UUID) (*entity.Appointment, error)
	FindByUserID(userID uuid.UUID, startTimeFilter, endTimeFilter *time.Time) ([]*entity.Appointment, error) // Lista agendamentos de um usuário, com filtros de data opcionais
	Update(appointment *entity.Appointment) error
	Delete(id uuid.UUID) error // Pode ser um soft delete ou hard delete
	// Adicione outros métodos conforme necessário, ex:
	// FindByDateRangeForAllUsers(start, end time.Time) ([]*entity.Appointment, error)
}
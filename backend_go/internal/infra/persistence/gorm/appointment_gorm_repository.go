package gorm

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/RenanGroh/SaaS_microbusiness_project/backend_go/internal/entity"
	"github.com/RenanGroh/SaaS_microbusiness_project/backend_go/internal/repository"
	"gorm.io/gorm"
	// "gorm.io/gorm/clause" // Para Preload e outras cláusulas
)

type gormAppointmentRepository struct {
	db *gorm.DB
}

func NewGormAppointmentRepository(db *gorm.DB) repository.AppointmentRepository {
	return &gormAppointmentRepository{db: db}
}

func (r *gormAppointmentRepository) Create(appointmentEntity *entity.Appointment) error {
	appointmentGorm := AppointmentFromEntity(appointmentEntity)
	result := r.db.Create(appointmentGorm)
	if result.Error != nil {
		return result.Error
	}
	// Atualizar a entidade original com ID e Timestamps gerados pelo GORM
	appointmentEntity.ID = appointmentGorm.ID
	appointmentEntity.CreatedAt = appointmentGorm.CreatedAt
	appointmentEntity.UpdatedAt = appointmentGorm.UpdatedAt
	return nil
}

func (r *gormAppointmentRepository) FindByID(id uuid.UUID) (*entity.Appointment, error) {
	var appointmentGorm AppointmentGormModel
	// Usar Preload para carregar dados do usuário associado (opcional, mas útil)
	result := r.db.Preload("User").Preload("ClientUser").First(&appointmentGorm, "id = ?", id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil // Indica não encontrado
		}
		return nil, result.Error
	}
	return appointmentGorm.ToEntity(), nil
}

func (r *gormAppointmentRepository) FindByUserID(userID uuid.UUID, startTimeFilter, endTimeFilter *time.Time) ([]*entity.Appointment, error) {
	var appointmentsGorm []AppointmentGormModel
	query := r.db.Preload("User").Preload("ClientUser").Where("user_id = ?", userID)

	if startTimeFilter != nil {
		query = query.Where("start_time >= ?", *startTimeFilter)
	}
	if endTimeFilter != nil {
		query = query.Where("end_time <= ?", *endTimeFilter)
	}

	// Ordenar por data de início, por exemplo
	result := query.Order("start_time asc").Find(&appointmentsGorm)
	if result.Error != nil {
		return nil, result.Error
	}

	var appointmentEntities []*entity.Appointment
	for _, ag := range appointmentsGorm {
		appointmentEntities = append(appointmentEntities, ag.ToEntity())
	}
	return appointmentEntities, nil
}

func (r *gormAppointmentRepository) Update(appointmentEntity *entity.Appointment) error {
	if appointmentEntity.ID == uuid.Nil {
		return errors.New("ID do agendamento não pode ser nulo para atualização")
	}
	appointmentGorm := AppointmentFromEntity(appointmentEntity)

	// Garante que apenas campos não nulos da struct sejam atualizados (se usando Updates)
	// ou usa Save para atualizar todos os campos.
	// Para atualizar campos específicos e também os "zero values" (ex: bool false, int 0)
	// usar .Select("*") com .Omit("CreatedAt") é uma boa prática.
	// Ou usar .Model(&AppointmentGormModel{ID: appointmentGorm.ID}).Updates(appointmentGorm)
	// que só atualiza campos não-zero. Para nosso caso, vamos assumir que todos os campos
	// da entidade são os desejados para atualização.
	result := r.db.Model(&AppointmentGormModel{}).Where("id = ?", appointmentGorm.ID).Updates(appointmentGorm)

	// Se você quer que "UpdatedAt" seja atualizado mesmo se nenhum outro campo mudou:
	// result := r.db.Save(appointmentGorm)

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("agendamento não encontrado para atualização ou nenhum dado alterado")
	}
	// Recarregar para obter o UpdatedAt mais recente se não estiver usando Save()
	// ou se Updates() não o atualiza automaticamente em todas as configs do GORM.
	// Mas se `appointmentGorm` tem `UpdatedAt time.Time`, o GORM deve cuidar disso.
	// A entidade já foi passada por referência, então se o `appointmentGorm` for atualizado
	// pelo `Updates` (o que não acontece para todos os campos), ela estaria atualizada.
	// Para garantir, você pode recarregar:
	// updatedAppointment, err := r.FindByID(appointmentEntity.ID)
	// if err == nil && updatedAppointment != nil {
	// 	*appointmentEntity = *updatedAppointment
	// }
	return nil
}

func (r *gormAppointmentRepository) Delete(id uuid.UUID) error {
	if id == uuid.Nil {
		return errors.New("ID do agendamento não pode ser nulo para deleção")
	}
	// Se estiver usando soft delete (gorm.DeletedAt na struct), o GORM faz soft delete.
	// Para hard delete: r.db.Unscoped().Delete(&AppointmentGormModel{}, id)
	result := r.db.Delete(&AppointmentGormModel{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("agendamento não encontrado para deleção")
	}
	return nil
}
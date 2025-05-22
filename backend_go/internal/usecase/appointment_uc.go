package usecase

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/RenanGroh/SaaS_microbusiness_project/backend_go/internal/entity"
	"github.com/RenanGroh/SaaS_microbusiness_project/backend_go/internal/repository"
	// "log" // Para debug
)

// AppointmentUseCase encapsula a lógica de negócios relacionada a agendamentos.
type AppointmentUseCase struct {
	appointmentRepo repository.AppointmentRepository
	userRepo        repository.UserRepository // Para verificar se o UserID existe, se necessário
}

// NewAppointmentUseCase cria uma nova instância de AppointmentUseCase.
func NewAppointmentUseCase(appRepo repository.AppointmentRepository, userRepo repository.UserRepository) *AppointmentUseCase {
	return &AppointmentUseCase{
		appointmentRepo: appRepo,
		userRepo:        userRepo,
	}
}

// CreateAppointmentInputDTO define os dados necessários para criar um agendamento.
// É bom ter DTOs de entrada para casos de uso para desacoplar da camada de delivery.
type CreateAppointmentInputDTO struct {
	UserID            uuid.UUID // ID do usuário (profissional) que está criando o agendamento
	ClientID          *uuid.UUID
	ClientName        string
	ClientEmail       string
	ClientPhone       string
	ServiceDescription string
	StartTime         time.Time
	EndTime           time.Time
	Notes             string
	Price             float64
	// Status inicial é geralmente PENDING, não precisa ser input
}

// CreateAppointment cria um novo agendamento.
func (uc *AppointmentUseCase) CreateAppointment(input CreateAppointmentInputDTO) (*entity.Appointment, error) {
	// Validações de negócio:
	// - UserID existe? (uc.userRepo.FindByID(input.UserID))
	// - ClientID existe, se fornecido? (uc.userRepo.FindByID(*input.ClientID)))
	// - StartTime é antes de EndTime?
	// - Não há conflitos de horário para este UserID? (Lógica mais complexa)
	// - Outras validações...

	if input.UserID == uuid.Nil {
		return nil, errors.New("ID do usuário é obrigatório")
	}
	if input.StartTime.IsZero() || input.EndTime.IsZero() {
		return nil, errors.New("data/hora de início e fim são obrigatórias")
	}
	if input.EndTime.Before(input.StartTime) || input.EndTime.Equal(input.StartTime) {
		return nil, errors.New("data/hora de término deve ser após a data/hora de início")
	}
	// Exemplo: Verificar se o profissional existe
	// professional, err := uc.userRepo.FindByID(input.UserID)
	// if err != nil || professional == nil {
	// 	return nil, errors.New("profissional (usuário) não encontrado")
	// }


	appointment := &entity.Appointment{
		ID:                uuid.New(), // Gerar novo UUID para o agendamento
		UserID:            input.UserID,
		ClientID:          input.ClientID,
		ClientName:        input.ClientName,
		ClientEmail:       input.ClientEmail,
		ClientPhone:       input.ClientPhone,
		ServiceDescription: input.ServiceDescription,
		StartTime:         input.StartTime,
		EndTime:           input.EndTime,
		Status:            entity.AppointmentStatusPending, // Status inicial
		Notes:             input.Notes,
		Price:             input.Price,
		// CreatedAt e UpdatedAt serão preenchidos pelo GORM/repo
	}

	err := uc.appointmentRepo.Create(appointment)
	if err != nil {
		// log.Printf("Erro ao criar agendamento no repositório: %v", err)
		return nil, errors.New("falha ao salvar agendamento: " + err.Error())
	}

	return appointment, nil
}

// GetAppointmentByID busca um agendamento pelo seu ID.
// Verifica se o userID fornecido (do token) tem permissão para ver este agendamento.
func (uc *AppointmentUseCase) GetAppointmentByID(appointmentID, requestingUserID uuid.UUID) (*entity.Appointment, error) {
	appointment, err := uc.appointmentRepo.FindByID(appointmentID)
	if err != nil {
		return nil, errors.New("erro ao buscar agendamento: " + err.Error())
	}
	if appointment == nil {
		return nil, errors.New("agendamento não encontrado")
	}

	// Regra de negócio: Usuário só pode ver seus próprios agendamentos
	if appointment.UserID != requestingUserID {
		// Ou se o usuário for um admin, ou se o requestingUserID for o ClientID, etc.
		// log.Printf("Tentativa de acesso não autorizado ao agendamento %s pelo usuário %s", appointmentID, requestingUserID)
		return nil, errors.New("acesso não autorizado ao agendamento")
	}

	return appointment, nil
}

// ListUserAppointments lista os agendamentos de um usuário específico.
func (uc *AppointmentUseCase) ListUserAppointments(userID uuid.UUID, startTimeFilter, endTimeFilter *time.Time) ([]*entity.Appointment, error) {
	if userID == uuid.Nil {
		return nil, errors.New("ID do usuário é obrigatório para listar agendamentos")
	}
	return uc.appointmentRepo.FindByUserID(userID, startTimeFilter, endTimeFilter)
}

// UpdateAppointmentInputDTO define os dados para atualizar um agendamento.
// Todos os campos são ponteiros para que possamos distinguir entre um valor não fornecido e um valor zero.
type UpdateAppointmentInputDTO struct {
	ClientID          *uuid.UUID
	ClientName        *string
	ClientEmail       *string
	ClientPhone       *string
	ServiceDescription *string
	StartTime         *time.Time
	EndTime           *time.Time
	Status            *entity.AppointmentStatus
	Notes             *string
	Price             *float64
}

// UpdateAppointment atualiza um agendamento existente.
// Verifica se o userID fornecido (do token) tem permissão.
func (uc *AppointmentUseCase) UpdateAppointment(appointmentID, requestingUserID uuid.UUID, input UpdateAppointmentInputDTO) (*entity.Appointment, error) {
	existingAppointment, err := uc.GetAppointmentByID(appointmentID, requestingUserID) // Reutiliza a verificação de permissão
	if err != nil {
		return nil, err // Erro já tratado por GetAppointmentByID (não encontrado ou não autorizado)
	}

	// Aplicar atualizações da input para a entidade existente
	updated := false
	if input.ClientID != nil {
		existingAppointment.ClientID = input.ClientID
		updated = true
	}
	if input.ClientName != nil {
		existingAppointment.ClientName = *input.ClientName
		updated = true
	}
	if input.ClientEmail != nil {
		existingAppointment.ClientEmail = *input.ClientEmail
		updated = true
	}
    if input.ClientPhone != nil {
		existingAppointment.ClientPhone = *input.ClientPhone
		updated = true
	}
	if input.ServiceDescription != nil {
		existingAppointment.ServiceDescription = *input.ServiceDescription
		updated = true
	}
	if input.StartTime != nil {
		existingAppointment.StartTime = *input.StartTime
		updated = true
	}
	if input.EndTime != nil {
		existingAppointment.EndTime = *input.EndTime
		updated = true
	}
	if input.Status != nil {
		existingAppointment.Status = *input.Status
		updated = true
	}
	if input.Notes != nil {
		existingAppointment.Notes = *input.Notes
		updated = true
	}
	if input.Price != nil {
		existingAppointment.Price = *input.Price
		updated = true
	}

	// Validação após atualização (ex: StartTime < EndTime)
	if existingAppointment.EndTime.Before(existingAppointment.StartTime) || existingAppointment.EndTime.Equal(existingAppointment.StartTime) {
		return nil, errors.New("data/hora de término deve ser após a data/hora de início")
	}

	if !updated {
		// log.Println("Nenhum campo fornecido para atualização do agendamento.")
		return existingAppointment, nil // Ou retorne um erro "nada para atualizar"
	}
	// existingAppointment.UpdatedAt será atualizado pelo GORM

	err = uc.appointmentRepo.Update(existingAppointment)
	if err != nil {
		// log.Printf("Erro ao atualizar agendamento %s no repositório: %v", appointmentID, err)
		return nil, errors.New("falha ao atualizar agendamento: " + err.Error())
	}

	return existingAppointment, nil
}

// CancelAppointment cancela um agendamento (exemplo de mudança de status).
// Verifica se o userID fornecido (do token) tem permissão.
func (uc *AppointmentUseCase) CancelAppointment(appointmentID, requestingUserID uuid.UUID) (*entity.Appointment, error) {
	appointment, err := uc.GetAppointmentByID(appointmentID, requestingUserID)
	if err != nil {
		return nil, err
	}

	// Regra de negócio: pode cancelar apenas se estiver PENDING ou CONFIRMED?
	if appointment.Status != entity.AppointmentStatusPending && appointment.Status != entity.AppointmentStatusConfirmed {
		return nil, errors.New("agendamento não pode ser cancelado no status atual: " + string(appointment.Status))
	}

	appointment.Status = entity.AppointmentStatusCancelled
	// appointment.UpdatedAt será atualizado pelo GORM

	err = uc.appointmentRepo.Update(appointment) // Reutiliza o método Update do repo
	if err != nil {
		return nil, errors.New("falha ao cancelar agendamento: " + err.Error())
	}

	return appointment, nil
}
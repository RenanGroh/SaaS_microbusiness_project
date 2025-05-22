package entity

import (
	"time"
	"github.com/google/uuid"
)

// AppointmentStatus define os possíveis status de um agendamento.
type AppointmentStatus string

const (
	AppointmentStatusPending   AppointmentStatus = "PENDING"
	AppointmentStatusConfirmed AppointmentStatus = "CONFIRMED"
	AppointmentStatusCancelled AppointmentStatus = "CANCELLED"
	AppointmentStatusCompleted AppointmentStatus = "COMPLETED"
	// Adicione outros status conforme necessário
)

// Appointment representa a entidade de agendamento no domínio.
type Appointment struct {
	ID                uuid.UUID // Chave primária do agendamento
	UserID            uuid.UUID // Chave estrangeira para o usuário (o profissional/MEI)
	ClientID          *uuid.UUID // Opcional: Se o cliente também for um usuário registrado no sistema
	ClientName        string    // Nome do cliente (se não for um usuário registrado)
	ClientEmail       string    // Email do cliente (para contato/notificações)
	ClientPhone       string    // Telefone do cliente
	ServiceDescription string    // Descrição do serviço a ser realizado
	StartTime         time.Time // Data e hora de início do agendamento
	EndTime           time.Time // Data e hora de término do agendamento
	Status            AppointmentStatus // Status do agendamento (PENDING, CONFIRMED, etc.)
	Notes             string    // Observações adicionais sobre o agendamento
	Price             float64   // Preço do serviço (opcional, pode ser gerenciado em outro lugar)
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

// Você pode adicionar construtores ou métodos de validação aqui se necessário.
// Ex: func NewAppointment(...) (*Appointment, error)
package http

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/RenanGroh/SaaS_microbusiness_project/backend_go/internal/delivery/http/middleware" // Para GetUserIDFromContext
	"github.com/RenanGroh/SaaS_microbusiness_project/backend_go/internal/entity"
	"github.com/RenanGroh/SaaS_microbusiness_project/backend_go/internal/usecase"
	// "log" // Para debug
)

// --- DTOs para Appointment ---

// CreateAppointmentRequest define o JSON esperado para criar um agendamento.
type CreateAppointmentRequest struct {
	// UserID não é necessário no request, pois será pego do token do usuário autenticado.
	ClientID          *string   `json:"clientId"` // string UUID ou nulo
	ClientName        string    `json:"clientName" binding:"required_without=ClientID,omitempty,min=2"`
	ClientEmail       string    `json:"clientEmail" binding:"omitempty,email"`
	ClientPhone       string    `json:"clientPhone"`
	ServiceDescription string    `json:"serviceDescription" binding:"required"`
	StartTime         time.Time `json:"startTime" binding:"required" time_format:"2006-01-02T15:04:05Z07:00"` // RFC3339
	EndTime           time.Time `json:"endTime" binding:"required" time_format:"2006-01-02T15:04:05Z07:00"`   // RFC3339
	Notes             string    `json:"notes"`
	Price             float64   `json:"price"`
}

// UpdateAppointmentRequest define o JSON para atualizar um agendamento.
// Todos os campos são opcionais (ponteiros).
type UpdateAppointmentRequest struct {
	ClientID          *string   `json:"clientId"`
	ClientName        *string   `json:"clientName"`
	ClientEmail       *string   `json:"clientEmail"`
	ClientPhone       *string   `json:"clientPhone"`
	ServiceDescription *string   `json:"serviceDescription"`
	StartTime         *time.Time `json:"startTime" time_format:"2006-01-02T15:04:05Z07:00"`
	EndTime           *time.Time `json:"endTime" time_format:"2006-01-02T15:04:05Z07:00"`
	Status            *string   `json:"status"` // String para o status (PENDING, CONFIRMED, etc.)
	Notes             *string   `json:"notes"`
	Price             *float64  `json:"price"`
}

// AppointmentResponse define o JSON retornado para um agendamento.
type AppointmentResponse struct {
	ID                uuid.UUID  `json:"id"`
	UserID            uuid.UUID  `json:"userId"`
	ClientID          *uuid.UUID `json:"clientId,omitempty"`
	ClientName        string     `json:"clientName"`
	ClientEmail       string     `json:"clientEmail"`
	ClientPhone       string     `json:"clientPhone"`
	ServiceDescription string     `json:"serviceDescription"`
	StartTime         time.Time  `json:"startTime"`
	EndTime           time.Time  `json:"endTime"`
	Status            string     `json:"status"` // Status como string
	Notes             string     `json:"notes"`
	Price             float64    `json:"price"`
	CreatedAt         time.Time  `json:"createdAt"`
	UpdatedAt         time.Time  `json:"updatedAt"`
	// User *UserResponse `json:"user,omitempty"` // Opcional: incluir dados do profissional
	// ClientUser *UserResponse `json:"clientUser,omitempty"` // Opcional: incluir dados do cliente se for um usuário
}

// --- AppointmentHandler ---
type AppointmentHandler struct {
	appointmentUseCase *usecase.AppointmentUseCase
}

func NewAppointmentHandler(uc *usecase.AppointmentUseCase) *AppointmentHandler {
	return &AppointmentHandler{appointmentUseCase: uc}
}

// mapEntityToResponse converte uma entidade Appointment para AppointmentResponse DTO.
func mapAppointmentEntityToResponse(appEntity *entity.Appointment) AppointmentResponse {
	return AppointmentResponse{
		ID:                appEntity.ID,
		UserID:            appEntity.UserID,
		ClientID:          appEntity.ClientID,
		ClientName:        appEntity.ClientName,
		ClientEmail:       appEntity.ClientEmail,
		ClientPhone:       appEntity.ClientPhone,
		ServiceDescription: appEntity.ServiceDescription,
		StartTime:         appEntity.StartTime,
		EndTime:           appEntity.EndTime,
		Status:            string(appEntity.Status),
		Notes:             appEntity.Notes,
		Price:             appEntity.Price,
		CreatedAt:         appEntity.CreatedAt,
		UpdatedAt:         appEntity.UpdatedAt,
	}
}

// CreateAppointment godoc
// @Summary      Cria um novo agendamento para o usuário autenticado
// @Description  Cria um agendamento. O UserID é pego do token JWT.
// @Tags         appointments
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        appointment body CreateAppointmentRequest true "Dados do Agendamento"
// @Success      201  {object} AppointmentResponse "Agendamento criado"
// @Failure      400  {object} map[string]string "Dados inválidos"
// @Failure      401  {object} map[string]string "Não autorizado"
// @Failure      500  {object} map[string]string "Erro interno"
// @Router       /appointments [post]
func (h *AppointmentHandler) CreateAppointment(c *gin.Context) {
	requestingUserID, exists := middleware.GetUserIDFromContext(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuário não autenticado"})
		return
	}

	var req CreateAppointmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dados inválidos na requisição", "details": err.Error()})
		return
	}

	// Converter ClientID de string para *uuid.UUID se fornecido
	var clientIDPtr *uuid.UUID
	if req.ClientID != nil && *req.ClientID != "" {
		parsedClientID, err := uuid.Parse(*req.ClientID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "ClientID inválido, deve ser um UUID ou vazio"})
			return
		}
		clientIDPtr = &parsedClientID
	}


	inputDTO := usecase.CreateAppointmentInputDTO{
		UserID:            requestingUserID,
		ClientID:          clientIDPtr,
		ClientName:        req.ClientName,
		ClientEmail:       req.ClientEmail,
		ClientPhone:       req.ClientPhone,
		ServiceDescription: req.ServiceDescription,
		StartTime:         req.StartTime,
		EndTime:           req.EndTime,
		Notes:             req.Notes,
		Price:             req.Price,
	}

	appointmentEntity, err := h.appointmentUseCase.CreateAppointment(inputDTO)
	if err != nil {
		// Tratar erros específicos do caso de uso (ex: conflito de horário, etc.)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Falha ao criar agendamento: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, mapAppointmentEntityToResponse(appointmentEntity))
}

// GetAppointmentByID godoc
// @Summary      Busca um agendamento específico pelo ID
// @Description  Retorna detalhes de um agendamento se o usuário autenticado tiver permissão.
// @Tags         appointments
// @Security     BearerAuth
// @Produce      json
// @Param        id path string true "ID do Agendamento (UUID)"
// @Success      200  {object} AppointmentResponse
// @Failure      400  {object} map[string]string "ID inválido"
// @Failure      401  {object} map[string]string "Não autorizado"
// @Failure      404  {object} map[string]string "Agendamento não encontrado"
// @Failure      500  {object} map[string]string "Erro interno"
// @Router       /appointments/{id} [get]
func (h *AppointmentHandler) GetAppointmentByID(c *gin.Context) {
	requestingUserID, exists := middleware.GetUserIDFromContext(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuário não autenticado"})
		return
	}

	appointmentIDStr := c.Param("id")
	appointmentID, err := uuid.Parse(appointmentIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID do agendamento inválido"})
		return
	}

	appointmentEntity, err := h.appointmentUseCase.GetAppointmentByID(appointmentID, requestingUserID)
	if err != nil {
		if err.Error() == "agendamento não encontrado" || err.Error() == "acesso não autorizado ao agendamento" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()}) // Ou 403 para não autorizado
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar agendamento: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, mapAppointmentEntityToResponse(appointmentEntity))
}

// ListUserAppointments godoc
// @Summary      Lista os agendamentos do usuário autenticado
// @Description  Retorna uma lista de agendamentos do usuário, com filtros opcionais de data.
// @Tags         appointments
// @Security     BearerAuth
// @Produce      json
// @Param        startTime query string false "Data/Hora de Início do Filtro (RFC3339, ex: 2023-01-01T00:00:00Z)"
// @Param        endTime query string false "Data/Hora de Fim do Filtro (RFC3339, ex: 2023-01-31T23:59:59Z)"
// @Success      200  {array}  AppointmentResponse
// @Failure      400  {object} map[string]string "Parâmetro de data inválido"
// @Failure      401  {object} map[string]string "Não autorizado"
// @Failure      500  {object} map[string]string "Erro interno"
// @Router       /appointments [get]
func (h *AppointmentHandler) ListUserAppointments(c *gin.Context) {
	requestingUserID, exists := middleware.GetUserIDFromContext(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuário não autenticado"})
		return
	}

	var startTimeFilter, endTimeFilter *time.Time
	if c.Query("startTime") != "" {
		st, err := time.Parse(time.RFC3339, c.Query("startTime"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Formato de startTime inválido, use RFC3339"})
			return
		}
		startTimeFilter = &st
	}
	if c.Query("endTime") != "" {
		et, err := time.Parse(time.RFC3339, c.Query("endTime"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Formato de endTime inválido, use RFC3339"})
			return
		}
		endTimeFilter = &et
	}

	appointmentEntities, err := h.appointmentUseCase.ListUserAppointments(requestingUserID, startTimeFilter, endTimeFilter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao listar agendamentos: " + err.Error()})
		return
	}

	responses := make([]AppointmentResponse, len(appointmentEntities))
	for i, appEntity := range appointmentEntities {
		responses[i] = mapAppointmentEntityToResponse(appEntity)
	}
	c.JSON(http.StatusOK, responses)
}

// UpdateAppointment godoc
// @Summary      Atualiza um agendamento existente
// @Description  Atualiza os campos de um agendamento se o usuário autenticado tiver permissão.
// @Tags         appointments
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id path string true "ID do Agendamento (UUID)"
// @Param        appointment body UpdateAppointmentRequest true "Dados para Atualização"
// @Success      200  {object} AppointmentResponse "Agendamento atualizado"
// @Failure      400  {object} map[string]string "ID ou dados inválidos"
// @Failure      401  {object} map[string]string "Não autorizado"
// @Failure      404  {object} map[string]string "Agendamento não encontrado"
// @Failure      500  {object} map[string]string "Erro interno"
// @Router       /appointments/{id} [put]
func (h *AppointmentHandler) UpdateAppointment(c *gin.Context) {
	requestingUserID, exists := middleware.GetUserIDFromContext(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuário não autenticado"})
		return
	}

	appointmentIDStr := c.Param("id")
	appointmentID, err := uuid.Parse(appointmentIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID do agendamento inválido"})
		return
	}

	var req UpdateAppointmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dados inválidos na requisição", "details": err.Error()})
		return
	}

	// Converter UpdateAppointmentRequest para UpdateAppointmentInputDTO do usecase
	updateDTO := usecase.UpdateAppointmentInputDTO{}
	if req.ClientID != nil {
		if *req.ClientID == "" { // Permitir enviar "" para limpar ClientID
            var nilUUID *uuid.UUID // explicitamente nil
            updateDTO.ClientID = nilUUID
        } else {
            parsedClientID, err := uuid.Parse(*req.ClientID)
            if err != nil {
                c.JSON(http.StatusBadRequest, gin.H{"error": "ClientID inválido na atualização, deve ser um UUID ou vazio"})
                return
            }
            updateDTO.ClientID = &parsedClientID
        }
	}
	updateDTO.ClientName = req.ClientName
	updateDTO.ClientEmail = req.ClientEmail
    updateDTO.ClientPhone = req.ClientPhone
	updateDTO.ServiceDescription = req.ServiceDescription
	updateDTO.StartTime = req.StartTime
	updateDTO.EndTime = req.EndTime
	if req.Status != nil {
		status := entity.AppointmentStatus(*req.Status)
		// Validação do status (ex: deve ser um dos valores válidos de AppointmentStatus)
		// if status != entity.AppointmentStatusPending && status != entity.AppointmentStatusConfirmed ... {
		// c.JSON(http.StatusBadRequest, gin.H{"error": "Status inválido"})
		// return
		// }
		updateDTO.Status = &status
	}
	updateDTO.Notes = req.Notes
	updateDTO.Price = req.Price

	updatedAppointmentEntity, err := h.appointmentUseCase.UpdateAppointment(appointmentID, requestingUserID, updateDTO)
	if err != nil {
		// Tratar erros do caso de uso
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Falha ao atualizar agendamento: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, mapAppointmentEntityToResponse(updatedAppointmentEntity))
}

// CancelAppointment godoc
// @Summary      Cancela um agendamento
// @Description  Muda o status de um agendamento para CANCELLED se o usuário tiver permissão.
// @Tags         appointments
// @Security     BearerAuth
// @Produce      json
// @Param        id path string true "ID do Agendamento (UUID)"
// @Success      200  {object} AppointmentResponse "Agendamento cancelado"
// @Failure      400  {object} map[string]string "ID inválido ou agendamento não pode ser cancelado"
// @Failure      401  {object} map[string]string "Não autorizado"
// @Failure      404  {object} map[string]string "Agendamento não encontrado"
// @Failure      500  {object} map[string]string "Erro interno"
// @Router       /appointments/{id}/cancel [patch]
func (h *AppointmentHandler) CancelAppointment(c *gin.Context) {
	requestingUserID, exists := middleware.GetUserIDFromContext(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuário não autenticado"})
		return
	}

	appointmentIDStr := c.Param("id")
	appointmentID, err := uuid.Parse(appointmentIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID do agendamento inválido"})
		return
	}

	cancelledAppointmentEntity, err := h.appointmentUseCase.CancelAppointment(appointmentID, requestingUserID)
	if err != nil {
		// Tratar erros do caso de uso (ex: status inválido para cancelamento)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Falha ao cancelar agendamento: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, mapAppointmentEntityToResponse(cancelledAppointmentEntity))
}
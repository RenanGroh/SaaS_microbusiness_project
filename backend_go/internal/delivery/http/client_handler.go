package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/RenanGroh/SaaS_microbusiness_project/backend_go/internal/delivery/http/middleware"
	"github.com/RenanGroh/SaaS_microbusiness_project/backend_go/internal/entity"
	"github.com/RenanGroh/SaaS_microbusiness_project/backend_go/internal/usecase"
)

// --- DTOs para Client ---

type CreateClientRequest struct {
	Name  string `json:"name" binding:"required,min=2"`
	Email string `json:"email" binding:"omitempty,email"`
	Phone string `json:"phone"`
	Notes string `json:"notes"`
}

type UpdateClientRequest struct {
	Name  *string `json:"name"`
	Email *string `json:"email"`
	Phone *string `json:"phone"`
	Notes *string `json:"notes"`
}

type ClientResponse struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"userId"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Phone     string    `json:"phone"`
	Notes     string    `json:"notes"`
	CreatedAt string    `json:"createdAt"`
	UpdatedAt string    `json:"updatedAt"`
}

// --- ClientHandler ---
type ClientHandler struct {
	clientUseCase *usecase.ClientUseCase
}

func NewClientHandler(uc *usecase.ClientUseCase) *ClientHandler {
	return &ClientHandler{clientUseCase: uc}
}

func mapClientEntityToResponse(clientEntity *entity.Client) ClientResponse {
	return ClientResponse{
		ID:        clientEntity.ID,
		UserID:    clientEntity.UserID,
		Name:      clientEntity.Name,
		Email:     clientEntity.Email,
		Phone:     clientEntity.Phone,
		Notes:     clientEntity.Notes,
		CreatedAt: clientEntity.CreatedAt.Format(time.RFC3339),
		UpdatedAt: clientEntity.UpdatedAt.Format(time.RFC3339),
	}
}

// CreateClient godoc
// @Summary      Cria um novo cliente para o usuário autenticado
// @Description  Cria um cliente. O UserID é pego do token JWT.
// @Tags         clients
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        client body CreateClientRequest true "Dados do Cliente"
// @Success      201  {object} ClientResponse "Cliente criado"
// @Failure      400  {object} map[string]string "Dados inválidos"
// @Failure      401  {object} map[string]string "Não autorizado"
// @Failure      500  {object} map[string]string "Erro interno"
// @Router       /clients [post]
func (h *ClientHandler) CreateClient(c *gin.Context) {
	requestingUserID, exists := middleware.GetUserIDFromContext(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuário não autenticado"})
		return
	}

	var req CreateClientRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dados inválidos na requisição", "details": err.Error()})
		return
	}

	inputDTO := usecase.CreateClientInputDTO{
		UserID: requestingUserID,
		Name:   req.Name,
		Email:  req.Email,
		Phone:  req.Phone,
		Notes:  req.Notes,
	}

	clientEntity, err := h.clientUseCase.CreateClient(inputDTO)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Falha ao criar cliente: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, mapClientEntityToResponse(clientEntity))
}

// GetClientByID godoc
// @Summary      Busca um cliente específico pelo ID
// @Description  Retorna detalhes de um cliente se o usuário autenticado tiver permissão.
// @Tags         clients
// @Security     BearerAuth
// @Produce      json
// @Param        id path string true "ID do Cliente (UUID)"
// @Success      200  {object} ClientResponse
// @Failure      400  {object} map[string]string "ID inválido"
// @Failure      401  {object} map[string]string "Não autorizado"
// @Failure      404  {object} map[string]string "Cliente não encontrado"
// @Failure      500  {object} map[string]string "Erro interno"
// @Router       /clients/{id} [get]
func (h *ClientHandler) GetClientByID(c *gin.Context) {
	requestingUserID, exists := middleware.GetUserIDFromContext(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuário não autenticado"})
		return
	}

	clientIDStr := c.Param("id")
	clientID, err := uuid.Parse(clientIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID do cliente inválido"})
		return
	}

	clientEntity, err := h.clientUseCase.GetClientByID(clientID, requestingUserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar cliente: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, mapClientEntityToResponse(clientEntity))
}

// ListUserClients godoc
// @Summary      Lista os clientes do usuário autenticado
// @Description  Retorna uma lista de clientes do usuário.
// @Tags         clients
// @Security     BearerAuth
// @Produce      json
// @Success      200  {array}  ClientResponse
// @Failure      401  {object} map[string]string "Não autorizado"
// @Failure      500  {object} map[string]string "Erro interno"
// @Router       /clients [get]
func (h *ClientHandler) ListUserClients(c *gin.Context) {
	requestingUserID, exists := middleware.GetUserIDFromContext(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuário não autenticado"})
		return
	}

	clientEntities, err := h.clientUseCase.ListUserClients(requestingUserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao listar clientes: " + err.Error()})
		return
	}

	responses := make([]ClientResponse, len(clientEntities))
	for i, clientEntity := range clientEntities {
		responses[i] = mapClientEntityToResponse(clientEntity)
	}
	c.JSON(http.StatusOK, responses)
}

// UpdateClient godoc
// @Summary      Atualiza um cliente existente
// @Description  Atualiza os campos de um cliente se o usuário autenticado tiver permissão.
// @Tags         clients
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id path string true "ID do Cliente (UUID)"
// @Param        client body UpdateClientRequest true "Dados para Atualização"
// @Success      200  {object} ClientResponse "Cliente atualizado"
// @Failure      400  {object} map[string]string "ID ou dados inválidos"
// @Failure      401  {object} map[string]string "Não autorizado"
// @Failure      404  {object} map[string]string "Cliente não encontrado"
// @Failure      500  {object} map[string]string "Erro interno"
// @Router       /clients/{id} [put]
func (h *ClientHandler) UpdateClient(c *gin.Context) {
	requestingUserID, exists := middleware.GetUserIDFromContext(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuário não autenticado"})
		return
	}

	clientIDStr := c.Param("id")
	clientID, err := uuid.Parse(clientIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID do cliente inválido"})
		return
	}

	var req UpdateClientRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dados inválidos na requisição", "details": err.Error()})
		return
	}

	inputDTO := usecase.UpdateClientInputDTO{
		Name:  req.Name,
		Email: req.Email,
		Phone: req.Phone,
		Notes: req.Notes,
	}

	clientEntity, err := h.clientUseCase.UpdateClient(clientID, requestingUserID, inputDTO)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Falha ao atualizar cliente: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, mapClientEntityToResponse(clientEntity))
}

// DeleteClient godoc
// @Summary      Exclui um cliente
// @Description  Exclui um cliente se o usuário tiver permissão.
// @Tags         clients
// @Security     BearerAuth
// @Produce      json
// @Param        id path string true "ID do Cliente (UUID)"
// @Success      204  {string} string "No Content"
// @Failure      400  {object} map[string]string "ID inválido"
// @Failure      401  {object} map[string]string "Não autorizado"
// @Failure      404  {object} map[string]string "Cliente não encontrado"
// @Failure      500  {object} map[string]string "Erro interno"
// @Router       /clients/{id} [delete]
func (h *ClientHandler) DeleteClient(c *gin.Context) {
	requestingUserID, exists := middleware.GetUserIDFromContext(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuário não autenticado"})
		return
	}

	clientIDStr := c.Param("id")
	clientID, err := uuid.Parse(clientIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID do cliente inválido"})
		return
	}

	err = h.clientUseCase.DeleteClient(clientID, requestingUserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Falha ao excluir cliente: " + err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

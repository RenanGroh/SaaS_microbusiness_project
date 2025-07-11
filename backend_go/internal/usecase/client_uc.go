package usecase

import (
	"errors"
	"github.com/google/uuid"
	"github.com/RenanGroh/SaaS_microbusiness_project/backend_go/internal/entity"
	"github.com/RenanGroh/SaaS_microbusiness_project/backend_go/internal/repository"
)

// ClientUseCase encapsula a lógica de negócios relacionada a clientes.
type ClientUseCase struct {
	clientRepo repository.ClientRepository
	userRepo   repository.UserRepository // Para verificar se o UserID existe, se necessário
}

// NewClientUseCase cria uma nova instância de ClientUseCase.
func NewClientUseCase(clientRepo repository.ClientRepository, userRepo repository.UserRepository) *ClientUseCase {
	return &ClientUseCase{
		clientRepo: clientRepo,
		userRepo:   userRepo,
	}
}

// CreateClientInputDTO define os dados necessários para criar um cliente.
type CreateClientInputDTO struct {
	UserID uuid.UUID
	Name   string
	Email  string
	Phone  string
	Notes  string
}

// CreateClient cria um novo cliente.
func (uc *ClientUseCase) CreateClient(input CreateClientInputDTO) (*entity.Client, error) {
	if input.UserID == uuid.Nil {
		return nil, errors.New("ID do usuário é obrigatório")
	}
	if input.Name == "" {
		return nil, errors.New("nome do cliente é obrigatório")
	}

	// Opcional: Verificar se o usuário (profissional) existe
	// _, err := uc.userRepo.FindByID(input.UserID)
	// if err != nil {
	// 	return nil, errors.New("usuário não encontrado: " + err.Error())
	// }

	client := &entity.Client{
		ID:     uuid.New(),
		UserID: input.UserID,
		Name:   input.Name,
		Email:  input.Email,
		Phone:  input.Phone,
		Notes:  input.Notes,
	}

	err := uc.clientRepo.Create(client)
	if err != nil {
		return nil, errors.New("falha ao salvar cliente: " + err.Error())
	}

	return client, nil
}

// GetClientByID busca um cliente pelo seu ID, verificando permissão.
func (uc *ClientUseCase) GetClientByID(clientID, requestingUserID uuid.UUID) (*entity.Client, error) {
	client, err := uc.clientRepo.FindByID(clientID)
	if err != nil {
		return nil, errors.New("erro ao buscar cliente: " + err.Error())
	}
	if client == nil {
		return nil, errors.New("cliente não encontrado")
	}

	// Regra de negócio: Usuário só pode ver seus próprios clientes
	if client.UserID != requestingUserID {
		return nil, errors.New("acesso não autorizado ao cliente")
	}

	return client, nil
}

// ListUserClients lista todos os clientes de um usuário específico.
func (uc *ClientUseCase) ListUserClients(userID uuid.UUID) ([]*entity.Client, error) {
	if userID == uuid.Nil {
		return nil, errors.New("ID do usuário é obrigatório para listar clientes")
	}
	return uc.clientRepo.FindByUserID(userID)
}

// UpdateClientInputDTO define os dados para atualizar um cliente.
type UpdateClientInputDTO struct {
	Name  *string
	Email *string
	Phone *string
	Notes *string
}

// UpdateClient atualiza um cliente existente.
func (uc *ClientUseCase) UpdateClient(clientID, requestingUserID uuid.UUID, input UpdateClientInputDTO) (*entity.Client, error) {
	existingClient, err := uc.GetClientByID(clientID, requestingUserID) // Reutiliza a verificação de permissão
	if err != nil {
		return nil, err
	}

	updated := false
	if input.Name != nil {
		existingClient.Name = *input.Name
		updated = true
	}
	if input.Email != nil {
		existingClient.Email = *input.Email
		updated = true
	}
	if input.Phone != nil {
		existingClient.Phone = *input.Phone
		updated = true
	}
	if input.Notes != nil {
		existingClient.Notes = *input.Notes
		updated = true
	}

	if !updated {
		return existingClient, nil // Nada para atualizar
	}

	err = uc.clientRepo.Update(existingClient)
	if err != nil {
		return nil, errors.New("falha ao atualizar cliente: " + err.Error())
	}

	return existingClient, nil
}

// DeleteClient exclui um cliente.
func (uc *ClientUseCase) DeleteClient(clientID, requestingUserID uuid.UUID) error {
	_, err := uc.GetClientByID(clientID, requestingUserID) // Reutiliza a verificação de permissão
	if err != nil {
		return err
	}

	return uc.clientRepo.Delete(clientID)
}

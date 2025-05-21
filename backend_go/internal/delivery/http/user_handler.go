package http

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/RenanGroh/SaaS_microbusiness_project/backend_go/internal/delivery/http/middleware"
	"github.com/RenanGroh/SaaS_microbusiness_project/backend_go/internal/usecase"
)

// -----------------------------------------------------------------------------
// DTOs (Data Transfer Objects) para este handler
// -----------------------------------------------------------------------------

// CreateUserInput define a estrutura para os dados de entrada da criação de usuário.
// Os campos devem ser exportados (começar com letra maiúscula) para o binding do Gin funcionar.
type CreateUserInput struct {
	Name     string `json:"name" binding:"required,min=3,max=100"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

// UserResponse define a estrutura para os dados de saída do usuário (sem senha).
type UserResponse struct {
    ID        uuid.UUID `json:"id"` // <<< MUDOU PARA uuid.UUID
    Name      string    `json:"name"`
    Email     string    `json:"email"`
    CreatedAt time.Time `json:"createdAt,omitempty"`
    UpdatedAt time.Time `json:"updatedAt,omitempty"`
}

type LoginInput struct {
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token string       `json:"token"`
	User  UserResponse `json:"user"` // Para retornar alguns dados do usuário
}

// -----------------------------------------------------------------------------
// UserHandler e seus métodos
// -----------------------------------------------------------------------------

// UserHandler encapsula os handlers HTTP relacionados a usuários.
type UserHandler struct {
	userUseCase *usecase.UserUseCase // Depende do caso de uso
}

// NewUserHandler cria uma nova instância de UserHandler.
func NewUserHandler(uc *usecase.UserUseCase) *UserHandler {
	return &UserHandler{userUseCase: uc}
}

// CreateUser é o handler para a rota POST /users (ou /api/v1/users)
func (h *UserHandler) CreateUser(c *gin.Context) {
	var input CreateUserInput // <<< AQUI ESTÁ O USO DE CreateUserInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dados inválidos", "details": err.Error()})
		return
	}

	// Chamar o caso de uso para criar o usuário
	createdUserEntity, err := h.userUseCase.CreateUser(input.Name, input.Email, input.Password)
	if err != nil {
		// Mapear erros do caso de uso para respostas HTTP apropriadas
		if err.Error() == "email já está em uso" { // Melhorar com tipos de erro customizados no futuro
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Falha ao criar usuário: " + err.Error()})
		return
	}

	// Mapear a entidade retornada pelo caso de uso para o DTO de resposta
	response := UserResponse{
        ID:        createdUserEntity.ID,
        Name:      createdUserEntity.Name,
        Email:     createdUserEntity.Email,
        CreatedAt: createdUserEntity.CreatedAt, // <<< IMPORTANTE
        UpdatedAt: createdUserEntity.UpdatedAt, // <<< IMPORTANTE
    }
    c.JSON(http.StatusCreated, response)
}

// GetUserByEmail é o handler para a rota GET /users/by-email (ou /api/v1/users/by-email)
func (h *UserHandler) GetUserByEmail(c *gin.Context) {
	email := c.Query("email")
	if email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Parâmetro 'email' é obrigatório"})
		return
	}

	userEntity, err := h.userUseCase.GetUserByEmail(email) // Supondo que este método existe no UserUseCase
	if err != nil {
		if err.Error() == "usuário não encontrado" { // Melhorar com tipos de erro
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar usuário: " + err.Error()})
		return
	}

	response := UserResponse{
		ID:    userEntity.ID,
		Name:  userEntity.Name,
		Email: userEntity.Email,
	}
	c.JSON(http.StatusOK, response)
}

// GetUserByID (Placeholder, se você descomentar a rota que o usa)
// GetUserByID godoc
// @Summary      Busca um usuário pelo ID
// @Description  Retorna os dados de um usuário dado seu ID (UUID)
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        id path string true "ID do Usuário (UUID)"
// @Success      200  {object} UserResponse
// @Failure      400  {object} map[string]string "ID inválido (não é UUID)"
// @Failure      404  {object} map[string]string "Usuário não encontrado"
// @Failure      500  {object} map[string]string "Erro interno"
// @Router       /users/{id} [get] // Ajuste a rota se necessário no router.go
func (h *UserHandler) GetUserByID(c *gin.Context) {
	idParam := c.Param("id")
	userID, err := uuid.Parse(idParam) // Faz o parse do parâmetro da URL para UUID
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de usuário inválido, deve ser um UUID"})
		return
	}

	userEntity, err := h.userUseCase.GetUserByID(userID) // Chama o método do UserUseCase
	if err != nil {
		if err.Error() == "usuário não encontrado" { // Melhorar com tipos de erro customizados
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar usuário: " + err.Error()})
		return
	}

	response := UserResponse{
		ID:        userEntity.ID,
		Name:      userEntity.Name,
		Email:     userEntity.Email,
		CreatedAt: userEntity.CreatedAt,
		UpdatedAt: userEntity.UpdatedAt,
	}
	c.JSON(http.StatusOK, response)
}

// Login godoc
// @Summary      Autentica um usuário
// @Description  Autentica um usuário com email e senha e retorna um token (futuramente)
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        credentials body LoginInput true "Credenciais de Login"
// @Success      200  {object} map[string]string "Login bem-sucedido (futuramente com token)"
// @Failure      400  {object} map[string]string "Dados inválidos"
// @Failure      401  {object} map[string]string "Credenciais inválidas"
// @Failure      500  {object} map[string]string "Erro interno do servidor"
// @Router       /auth/login [post]
func (h *UserHandler) Login(c *gin.Context) {
	var input LoginInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dados inválidos", "details": err.Error()})
		return
	}

	// O UserCase.Login agora retorna (token, userEntity, error)
	tokenString, loggedInUserEntity, err := h.userUseCase.Login(input.Email, input.Password)
	if err != nil {
		if err.Error() == "credenciais inválidas" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro no processo de login: " + err.Error()})
		return
	}

	// Mapear a entidade para o UserResponse DTO
	userResponseData := UserResponse{
		ID:    loggedInUserEntity.ID,
		Name:  loggedInUserEntity.Name,
		Email: loggedInUserEntity.Email,
		CreatedAt: loggedInUserEntity.CreatedAt, // Adicione se quiser
		UpdatedAt: loggedInUserEntity.UpdatedAt, // Adicione se quiser
	}

	// Criar a resposta final com o token e os dados do usuário
	loginResp := LoginResponse{
		Token: tokenString,
		User:  userResponseData,
	}

	c.JSON(http.StatusOK, loginResp)
}

// GetUserProfile godoc
// @Summary      Obtém o perfil do usuário autenticado
// @Description  Retorna os dados do usuário atualmente logado (requer token JWT)
// @Tags         users
// @Security     BearerAuth
// @Produce      json
// @Success      200  {object} UserResponse "Perfil do usuário"
// @Failure      401  {object} map[string]string "Não autorizado"
// @Failure      404  {object} map[string]string "Usuário não encontrado (raro se o token é válido)"
// @Router       /users/me [get]
func (h *UserHandler) GetUserProfile(c *gin.Context) {
	// Extrair o UserID do contexto Gin, que foi colocado pelo AuthMiddleware
	// Usando a função helper:
	userID, exists := middleware.GetUserIDFromContext(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Falha ao obter ID do usuário do token"})
		return
	}

	// Alternativamente, acessando diretamente as claims:
	// authPayload := c.MustGet(middleware.AuthorizationPayloadKey).(*security.Claims)
	// userID := authPayload.UserID

	// Buscar o usuário usando o ID do token
	// Você precisará de um método no UserUseCase como GetUserByID(id uuid.UUID)
	// que por sua vez chama userRepo.FindByID(id uuid.UUID)
	userEntity, err := h.userUseCase.GetUserByID(userID) // << PRECISA CRIAR UserUseCase.GetUserByID
	if err != nil {
		// Tratar erro, ex: usuário não encontrado (embora improvável se o token é válido e recente)
		c.JSON(http.StatusNotFound, gin.H{"error": "Usuário não encontrado: " + err.Error()})
		return
	}
	if userEntity == nil { // Caso o use case retorne nil, nil para não encontrado
		c.JSON(http.StatusNotFound, gin.H{"error": "Usuário não encontrado"})
		return
	}


	response := UserResponse{
		ID:        userEntity.ID,
		Name:      userEntity.Name,
		Email:     userEntity.Email,
		CreatedAt: userEntity.CreatedAt,
		UpdatedAt: userEntity.UpdatedAt,
	}
	c.JSON(http.StatusOK, response)
}
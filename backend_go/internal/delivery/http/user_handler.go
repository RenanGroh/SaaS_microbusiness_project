package http // Pacote http (ou httpdelivery, como você preferir chamar)

import (
	"net/http" // Para códigos de status HTTP
	"time"

	"github.com/gin-gonic/gin"
	// Certifique-se que o path para 'usecase' está correto e corresponde ao seu go.mod
	"github.com/RenanGroh/SaaS_microbusiness_project/backend_go/internal/usecase"
	// "github.com/RenanGroh/SaaS_microbusiness_project/backend_go/internal/entity" // Não é mais necessário aqui se o usecase retorna a entidade
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
    ID        uint      `json:"id"`
    Name      string    `json:"name"`
    Email     string    `json:"email"`
    CreatedAt time.Time `json:"createdAt,omitempty"` // Verifique se está aqui
    UpdatedAt time.Time `json:"updatedAt"` // Verifique se está aqui
}

type LoginInput struct {
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required"`
}

// LoginResponse (placeholder por enquanto, depois adicionaremos o token aqui)
// type LoginResponse struct {
//  Token string `json:"token"`
//  User  UserResponse `json:"user"` // Opcional, se quiser retornar dados do usuário
// }

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
func (h *UserHandler) GetUserByID(c *gin.Context) {
	id := c.Param("id")
	c.JSON(http.StatusOK, gin.H{"message": "GetUserByID - placeholder", "id": id})
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

	loggedInUser, err := h.userUseCase.Login(input.Email, input.Password)
	if err != nil {
		if err.Error() == "credenciais inválidas" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro no processo de login: " + err.Error()})
		return
	}

	// Por enquanto, apenas uma mensagem de sucesso. Depois adicionaremos o token JWT.
	// Opcional: retornar alguns dados do usuário (sem a senha!)
	userResponse := UserResponse{
		ID:    loggedInUser.ID,
		Name:  loggedInUser.Name,
		Email: loggedInUser.Email,
		// CreatedAt: loggedInUser.CreatedAt, // Se quiser incluir
		// UpdatedAt: loggedInUser.UpdatedAt, // Se quiser incluir
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Login bem-sucedido!",
		"user":    userResponse,
		// "token": "AQUI_VAI_O_TOKEN_JWT_FUTURAMENTE",
	})
}
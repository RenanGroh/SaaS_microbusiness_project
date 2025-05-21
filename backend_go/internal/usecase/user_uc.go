package usecase

import (
	"errors"

	"github.com/RenanGroh/SaaS_microbusiness_project/backend_go/internal/entity" // Ajuste o path do módulo
	"github.com/RenanGroh/SaaS_microbusiness_project/backend_go/internal/infra/security"
	"github.com/RenanGroh/SaaS_microbusiness_project/backend_go/internal/repository"
	//"log"
	"github.com/google/uuid"
)

// UserUseCase encapsula a lógica de negócios relacionada a usuários.
type UserUseCase struct {
	userRepo          repository.UserRepository
	jwtSecret         string // Nova dependência
	jwtExpirationHours int    // Nova dependência
}

// NewUserUseCase cria uma nova instância de UserUseCase.
func NewUserUseCase(
	repo repository.UserRepository,
	jwtSecret string, // Adicionar
	jwtExpirationHours int, // Adicionar
) *UserUseCase {
	return &UserUseCase{
		userRepo:          repo,
		jwtSecret:         jwtSecret,
		jwtExpirationHours: jwtExpirationHours,
	}
}

// CreateUser é o caso de uso para criar um novo usuário.
// Ele lida com a validação de negócios, hashing de senha e persistência.
func (uc *UserUseCase) CreateUser(name, email, rawPassword string) (*entity.User, error) {
	// 1. Validação de Negócio (Exemplo: verificar se email já existe)
	existingUser, err := uc.userRepo.FindByEmail(email)
	if err != nil {
		// Se o erro não for "não encontrado", é um erro do repositório
		// (o repositório deve retornar nil,nil se não encontrar, ou um erro específico)
		// Aqui, assumimos que o repo retorna (nil, nil) se não encontrar.
		// Se retornar um erro genérico, esta lógica precisa mudar.
		// Vamos assumir que FindByEmail retorna (nil, nil) para not found e um erro real para outros problemas.
		// Esta lógica de tratamento de erro de "usuário já existe" pode ser mais robusta.
	}
	if existingUser != nil {
		return nil, errors.New("email já está em uso") // Erro de negócio específico
	}

	// 2. Hashear a senha
	// hashedPassword, err := uc.passwordHasher.Hash(rawPassword) // Se usar interface
	hashedPassword, err := security.HashPassword(rawPassword) // Chamando diretamente por enquanto
	if err != nil {
		return nil, errors.New("falha ao processar senha") // Erro interno
	}

	// 3. Criar a entidade User
	user := &entity.User{
		Name:     name,
		Email:    email,
		Password: hashedPassword,
		// CreatedAt e UpdatedAt serão preenchidos pelo GORM ou pelo repo
	}

	// 4. Persistir o usuário usando o repositório
	err = uc.userRepo.Create(user)
	if err != nil {
		// O repositório pode retornar um erro se, por exemplo, houver uma race condition
		// e o email foi cadastrado por outra requisição entre a verificação e o Create.
		return nil, errors.New("falha ao salvar usuário: " + err.Error())
	}

	// 5. Retornar a entidade do usuário criado (o ID e Timestamps devem ter sido populados pelo repo/GORM)
	return user, nil
}

func (uc *UserUseCase) Login(email, rawPassword string) (tokenString string, userDetails *entity.User, err error) {
	// 1. Buscar usuário pelo email
	user, err := uc.userRepo.FindByEmail(email)
	if err != nil {
		return "", nil, errors.New("erro interno ao tentar autenticar")
	}
	if user == nil {
		return "", nil, errors.New("credenciais inválidas")
	}

	// 2. Verificar a senha
	passwordMatch := security.CheckPasswordHash(rawPassword, user.Password)
	if !passwordMatch {
		return "", nil, errors.New("credenciais inválidas")
	}

	// 3. Gerar o token JWT
	token, err := security.GenerateJWT(user.ID, user.Email, uc.jwtSecret, uc.jwtExpirationHours) // A CHAMADA AQUI JÁ DEVE ESTAR CORRETA
	if err != nil {
		return "", nil, errors.New("falha ao gerar token de autenticação")
	}

	// 4. Autenticação bem-sucedida
	return token, user, nil	
}

// GetUserByEmail é um exemplo de outro caso de uso.
func (uc *UserUseCase) GetUserByEmail(email string) (*entity.User, error) {
	user, err := uc.userRepo.FindByEmail(email)
	if err != nil {
		return nil, errors.New("erro ao buscar usuário por email")
	}
	if user == nil {
		return nil, errors.New("usuário não encontrado") // Erro de negócio específico
	}
	return user, nil
}

// GetUserByID é o caso de uso para buscar um usuário pelo seu ID.
func (uc *UserUseCase) GetUserByID(id uuid.UUID) (*entity.User, error) {
	if id == uuid.Nil {
		return nil, errors.New("ID do usuário não pode ser nulo")
	}

	user, err := uc.userRepo.FindByID(id) // Chama o método do repositório
	if err != nil {
		// Erro do repositório (ex: problema de conexão)
		return nil, errors.New("erro ao buscar usuário por ID: " + err.Error())
	}
	if user == nil {
		// Repositório retornou nil, nil (usuário não encontrado)
		return nil, errors.New("usuário não encontrado") // Ou retorne (nil, nil) para o handler decidir
	}

	return user, nil
}
package usecase

import (
	"errors"

	"github.com/RenanGroh/SaaS_microbusiness_project/backend_go/internal/entity" // Ajuste o path do módulo
	"github.com/RenanGroh/SaaS_microbusiness_project/backend_go/internal/infra/security"
	"github.com/RenanGroh/SaaS_microbusiness_project/backend_go/internal/repository"
	//"log"
)

// UserUseCase encapsula a lógica de negócios relacionada a usuários.
type UserUseCase struct {
	userRepo repository.UserRepository
	// Se você criar interfaces para PasswordHasher:
	// passwordHasher security.PasswordHasher
}

// NewUserUseCase cria uma nova instância de UserUseCase.
func NewUserUseCase(repo repository.UserRepository /*, hasher security.PasswordHasher*/) *UserUseCase {
	return &UserUseCase{
		userRepo: repo,
		// passwordHasher: hasher,
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

func (uc *UserUseCase) Login(email, rawPassword string) (*entity.User, error) {
	// 1. Buscar usuário pelo email
	user, err := uc.userRepo.FindByEmail(email)
	if err != nil {
		// Erro do repositório (não "não encontrado")
		// log.Printf("Erro do repositório ao buscar email %s: %v", email, err)
		return nil, errors.New("erro interno ao tentar autenticar")
	}
	if user == nil {
		// log.Printf("Usuário não encontrado para o email: %s", email)
		return nil, errors.New("credenciais inválidas") // Mensagem genérica
	}

	// 2. Verificar a senha
	// passwordMatch := uc.passwordHasher.Check(rawPassword, user.Password) // Se usar interface
	passwordMatch := security.CheckPasswordHash(rawPassword, user.Password) // Chamando diretamente
	if !passwordMatch {
		// log.Printf("Senha incorreta para o usuário: %s", email)
		return nil, errors.New("credenciais inválidas") // Mensagem genérica
	}

	// 3. Autenticação bem-sucedida
	// log.Printf("Usuário autenticado com sucesso: %s", email)
	return user, nil
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
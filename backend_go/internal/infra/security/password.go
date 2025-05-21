package security

import "golang.org/x/crypto/bcrypt"

// Interface (opcional para agora, mas bom para o futuro)
 type PasswordHasher interface {
 	Hash(password string) (string, error)
 	Check(password, hash string) bool
}

// BcryptPasswordHasher implementa PasswordHasher usando bcrypt.
 type BcryptPasswordHasher struct{}

func NewBcryptPasswordHasher() *BcryptPasswordHasher {
 	return &BcryptPasswordHasher{}
}

// HashPassword gera um hash bcrypt para a senha fornecida.
// func (h *BcryptPasswordHasher) Hash(password string) (string, error) { // Se usar interface
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// CheckPasswordHash compara uma senha em texto puro com seu hash bcrypt.
// func (h *BcryptPasswordHasher) Check(password, hash string) bool { // Se usar interface
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
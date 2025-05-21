package security

import (
	"errors" // Para erros customizados
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	// Se você não tiver UserID como uint, ajuste o tipo em Claims
	// "github.com/RenanGroh/SaaS_microbusiness_project/backend_go/internal/entity"
)

// Claims define a estrutura das "reivindicações" que serão incluídas no token JWT.
// jwt.RegisteredClaims incorpora claims padrão como Issuer, Subject, Audience, ExpirationTime, etc.
type Claims struct {
	UserID uuid.UUID `json:"user_id"` // CORRETO: UUID
	Email  string    `json:"email"`
	// Roles []string `json:"roles"`
	jwt.RegisteredClaims
}

// GenerateJWT gera um novo token JWT para um usuário.
// Requer o ID do usuário, email, a chave secreta JWT e a duração da expiração.
func GenerateJWT(userID uuid.UUID, userEmail string, jwtSecretKey string, expirationHours int) (string, error) { // <<< userID AGORA É uuid.UUID
	if jwtSecretKey == "" {
		return "", errors.New("chave secreta JWT não pode ser vazia")
	}
	if userID == uuid.Nil { // Opcional: verificar se o UUID é válido/não nulo
		return "", errors.New("ID do usuário para JWT não pode ser nulo")
	}


	expirationTime := time.Now().Add(time.Duration(expirationHours) * time.Hour)

	claims := &Claims{
		UserID: userID, // <<< AGORA CORRESPONDE: userID é uuid.UUID, Claims.UserID é uuid.UUID
		Email:  userEmail,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			// Issuer:    "bizly.com",
			Subject:   userID.String(), // Usar o UUID como string para o Subject é uma boa prática
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(jwtSecretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ValidateJWT verifica se o token JWT fornecido é válido.
// Retorna as claims se o token for válido, ou um erro caso contrário.
func ValidateJWT(tokenString string, jwtSecretKey string) (*Claims, error) {
	if jwtSecretKey == "" {
		return nil, errors.New("chave secreta JWT não pode ser vazia")
	}

	claims := &Claims{} // Claims.UserID já é uuid.UUID

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("método de assinatura inesperado")
		}
		return []byte(jwtSecretKey), nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, errors.New("token expirado")
		}
		if errors.Is(err, jwt.ErrTokenMalformed) {
			return nil, errors.New("token malformado")
		}
		return nil, errors.New("token inválido: " + err.Error())
	}

	if !token.Valid {
		return nil, errors.New("token inválido")
	}

	return claims, nil
}
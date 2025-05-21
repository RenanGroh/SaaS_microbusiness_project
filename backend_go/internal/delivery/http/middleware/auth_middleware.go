package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/RenanGroh/SaaS_microbusiness_project/backend_go/internal/config" // Para JWTSecret
	"github.com/RenanGroh/SaaS_microbusiness_project/backend_go/internal/infra/security"
	"github.com/google/uuid" // Para o tipo UserID no contexto
)

const (
	AuthorizationHeaderKey  = "authorization"
	AuthorizationTypeBearer = "bearer"
	AuthorizationPayloadKey = "authorization_payload" // Chave para armazenar claims no contexto Gin
)

// AuthMiddleware é um middleware Gin para autenticação JWT.
func AuthMiddleware(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader(AuthorizationHeaderKey)
		if len(authHeader) == 0 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Cabeçalho de autorização não fornecido"})
			return
		}

		fields := strings.Fields(authHeader)
		if len(fields) < 2 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Formato do cabeçalho de autorização inválido"})
			return
		}

		authType := strings.ToLower(fields[0])
		if authType != AuthorizationTypeBearer {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Tipo de autorização não suportado: " + authType})
			return
		}

		accessToken := fields[1]
		claims, err := security.ValidateJWT(accessToken, cfg.JWTSecret)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token inválido ou expirado: " + err.Error()})
			return
		}

		// Adiciona as claims (ou apenas o UserID) ao contexto do Gin
		// para que os handlers subsequentes possam acessá-las.
		c.Set(AuthorizationPayloadKey, claims) // Armazena todas as claims
		// ou apenas o ID:
		// c.Set("userID", claims.UserID)

		c.Next() // Continua para o próximo handler na cadeia
	}
}

// Helper para obter UserID do contexto (opcional, mas útil)
func GetUserIDFromContext(c *gin.Context) (uuid.UUID, bool) {
    payload, exists := c.Get(AuthorizationPayloadKey)
    if !exists {
        return uuid.Nil, false
    }
    claims, ok := payload.(*security.Claims)
    if !ok {
        return uuid.Nil, false
    }
    return claims.UserID, true
}
package entity

import (
	"time"
	"github.com/google/uuid"
)

// Client representa a entidade de cliente no domínio.
type Client struct {
	ID        uuid.UUID
	UserID    uuid.UUID // ID do usuário (profissional/MEI) que é dono deste cliente
	Name      string
	Email     string
	Phone     string
	Notes     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

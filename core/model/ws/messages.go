package ws

import (
	"context"
	"time"
)

type Message struct {
	Id        int        `json:"id"`
	FromUser  string     `json:"from_user"`
	Nombre    *string    `json:"nombre"`
	Apellido  *string    `json:"apellido"`
	CasoId    string     `json:"caso_id"`
	MediaUrl  []string   `json:"media_url"`
	Content   string     `json:"content"`
	IsRead    bool       `json:"is_read"`
	CreatedOn *time.Time `json:"created_on"`
	IsDeleted bool       `json:"is_deleted"`
}

type MessageData struct {
	FromUser string `json:"from_user"`
	// ToUser   string  `json:"to_user"`
	CasoId   string   `json:"caso_id"`
	MediaUrl []string `json:"media_url"`
	Content  string   `json:"content"`
	IsRead   *bool    `json:"is_read"`
}

type WsRepository interface {
	GetMessages(ctx context.Context, casoId string) ([]Message, error)
	SaveMessage(ctx context.Context, m *MessageData) (Message, error)
}

type WsUseCase interface {
	GetMessages(ctx context.Context, casoId string) ([]Message, error)
}

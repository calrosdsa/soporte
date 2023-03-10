package ws

import (
	"context"
	"time"
)

type Message struct {
	Id              int    `json:"id"`
	ClienteId       string `json:"client_id"`
	FuncionarioId   string `json:"funcioanrio_id"`
	CasoId          string `json:"caso_id"`
	ClienteName     string `json:"client_name"`
	FuncionarioName string `json:"funcionario_name"`
	MediaUrl        string `json:"media_url"`
	Content         string `json:"content"`
	IsRead          bool   `json:"is_read"`
	CreatedOn       *time.Time `json:"created_on"`
	IsDeleted       bool   `json:"is_deleted"`
}

type MessageData struct {
	ClienteId       string  `json:"client_id"`
	FuncionarioId   string  `json:"funcionario_id"`
	CasoId          string  `json:"caso_id"`
	ClienteName     string  `json:"client_name"`
	FuncionarioName string  `json:"funcionario_name"`
	MediaUrl        *string `json:"media_url"`
	Content         string  `json:"content"`
	IsRead          *bool   `json:"is_read"`
}

type WsRepository interface {
	GetMessages(ctx context.Context, casoId string) ([]Message, error)
	SaveMessage(ctx context.Context, m *MessageData) (Message, error)
}

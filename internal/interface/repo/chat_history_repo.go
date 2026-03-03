package repo

import "context"

type ChatHistory struct {
	Role    string // "user" 或 "model"
	Message string
}

type HistoryRepository interface {
	GetHistory(ctx context.Context, userID string) ([]ChatHistory, error)
	SaveHistory(ctx context.Context, userID string, history ChatHistory) error
}

package repository

import (
	"GopherMind/internal/interface/repo"
	"context"
	"sync"
)

// 限制記憶長度，例如只保留最近 5 筆，避免 Token 爆炸
const maxHistory = 5

type ChatHistoryRepository struct {
	data map[string][]repo.ChatHistory
	mu   sync.RWMutex // 確保並發安全
}

func NewChatHistoryRepository() *ChatHistoryRepository {
	return &ChatHistoryRepository{
		data: make(map[string][]repo.ChatHistory),
	}
}

func (m *ChatHistoryRepository) GetHistory(ctx context.Context, userID string) ([]repo.ChatHistory, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	// 回傳副本，避免外部修改影響原始資料
	if h, ok := m.data[userID]; ok {
		return h, nil
	}
	return []repo.ChatHistory{}, nil
}

func (m *ChatHistoryRepository) SaveHistory(ctx context.Context, userID string, record repo.ChatHistory) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	h := m.data[userID]
	h = append(h, record)

	if len(h) > maxHistory {
		h = h[len(h)-maxHistory:]
	}

	m.data[userID] = h
	return nil
}

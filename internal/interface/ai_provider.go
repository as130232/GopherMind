package _interface

import "context"

type AiProvider interface {
	Name() string
	RegisterTools(tools map[string]BaseTool)
	Chat(ctx context.Context, userID, message string) (string, error)
	ChatByTool(ctx context.Context, userID, message string) (string, error)
	AnalyzeImage(ctx context.Context, mimeType string, imageBytes []byte, prompt string) (string, error)
}

package _interface

import "context"

// ToolMetadata 存放工具的中立定義
type ToolMetadata struct {
	Name        string
	Description string
	IsOnce      bool
	Parameters  map[string]ToolParameter
}

// ToolParameter 定義參數的中立描述
type ToolParameter struct {
	Name        string
	Type        string // string, number, boolean
	Description string
	Required    bool
}

// BaseTool 是所有工具必須實作的介面，它不依賴任何 AI SDK
type BaseTool interface {
	GetMetadata() ToolMetadata
	Execute(ctx context.Context, args map[string]any) (any, error)
}

package tool

import (
	"GopherMind/internal/interface"
)

// Registry 管理所有註冊的工具
type Registry struct {
	tools map[string]_interface.BaseTool
}

func NewRegistry() *Registry {
	r := &Registry{
		tools: make(map[string]_interface.BaseTool),
	}
	return r
}
func (r *Registry) Register(t _interface.BaseTool) {
	r.tools[t.GetMetadata().Name] = t
}

func (r *Registry) GetAllTools() map[string]_interface.BaseTool {
	return r.tools
}

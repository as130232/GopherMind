package tool

import (
	"GopherMind/internal/interface"
	"context"
	"fmt"
)

type WeatherTool struct {
}

func NewWeatherTool() *WeatherTool {
	return &WeatherTool{}
}

func (t *WeatherTool) GetMetadata() _interface.ToolMetadata {
	return _interface.ToolMetadata{
		Name:        "get_weather",
		Description: "獲取台灣天氣資訊",
		Parameters: map[string]_interface.ToolParameter{
			"location": {
				Type:        "string",
				Description: "地點，例如：台北",
				Required:    true,
			},
		},
	}
}

func (t *WeatherTool) Execute(ctx context.Context, args map[string]interface{}) (interface{}, error) {
	location, ok := args["location"].(string)
	if !ok {
		return nil, fmt.Errorf("missing argument: location")
	}
	return location, nil
}

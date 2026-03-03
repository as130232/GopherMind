package app

import (
	"GopherMind/config"
	"GopherMind/consts"
	"GopherMind/internal/interface"
	"GopherMind/internal/intergration/twse"
	"GopherMind/internal/platform/telegram"
	"GopherMind/internal/provider"
	"GopherMind/internal/repository"
	"GopherMind/internal/service"
	"GopherMind/internal/tool"
	foodTool "GopherMind/internal/tool/food"
	stockTool "GopherMind/internal/tool/stock"
	"context"
	"fmt"
	"net/http"
)

type App struct {
	cfg        *config.Config
	router     http.Handler
	aiProvider _interface.AiProvider
	foodSvc    *service.FoodService
	tgBot      *telegram.Bot
}

func New(ctx context.Context) (*App, error) {
	// 初始化配置，載入環境變數
	cfg, err := config.Load()
	if err != nil {
		return nil, fmt.Errorf("config load failed: %w", err)
	}
	// 初始化 repo
	chatHistoryRepo := repository.NewChatHistoryRepository()

	// 初始化 AI Provider，選擇AI模型，預設gemini
	var aiProvider _interface.AiProvider
	switch cfg.Ai.ActiveAI {
	case consts.TypeOpenAi:
		return nil, fmt.Errorf("openai not supported yet")
	default:
		aiProvider, err = provider.NewGeminiProvider(ctx, cfg.Ai.GeminiAPIKey, cfg.Ai.GeminiModel, chatHistoryRepo)
		if err != nil {
			return nil, err
		}
	}
	// 初始化 API Client
	twseClient := twse.NewClient()

	// 初始化 Service 層
	foodSvc := service.NewFoodService(aiProvider)

	// 初始化 Bot
	tgBot := telegram.NewBot(cfg.Bot.TelegramToken, aiProvider)

	// 初始化AI工具箱與註冊工具列表
	registry := tool.NewRegistry()
	registry.Register(stockTool.NewGetStockInfo(twseClient))
	registry.Register(foodTool.NewAnalyzeFoodNutritionTool(foodSvc))
	registry.Register(tool.NewWeatherTool())
	// 將工具列表註冊到AI中(由 Provider 內部實作翻譯邏輯)
	aiProvider.RegisterTools(registry.GetAllTools())

	// 組裝 Router
	router := NewRouter(aiProvider, tgBot, foodSvc)
	app := &App{
		cfg:        cfg,
		router:     router,
		aiProvider: aiProvider,
		tgBot:      tgBot,
		foodSvc:    foodSvc,
	}
	return app, nil
}

func (a *App) Run() error {
	fmt.Printf("🚀 GopherMind [%s] 啟動於埠號 %s\n", a.cfg.Server.Env, a.cfg.Server.Port)
	return http.ListenAndServe(":"+a.cfg.Server.Port, a.router)
}

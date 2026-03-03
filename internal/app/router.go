package app

import (
	"GopherMind/internal/handler"
	ai "GopherMind/internal/handler/ai"
	"GopherMind/internal/interface"
	"GopherMind/internal/platform/telegram"
	"GopherMind/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func NewRouter(aiProvider _interface.AiProvider, tgBot *telegram.Bot, foodSvc *service.FoodService) *chi.Mux {
	r := chi.NewRouter()

	// 注入基礎中介軟體
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// 系統路由
	r.Route("/server", func(r chi.Router) {
		r.Get("/health", handler.HealthCheck)
		// 未來可擴充 r.Post("/config", api.UpdateConfig)
	})

	// AI 邏輯路由群組
	r.Route("/ai", func(r chi.Router) {
		// 這裡注入 AI Provider，專門處理對話
		r.Post("/message", ai.MessageHandler(aiProvider))
		//r.Post("/analyze-meal", ai.AnalyzeFoodHandler(aiProvider))
		r.Post("/analyze-meal", ai.AnalyzeFoodHandler(aiProvider, foodSvc))
	})

	r.Route("/webhook", func(r chi.Router) {
		r.Post("/telegram", handler.TelegramHandler(tgBot))
	})
	return r
}

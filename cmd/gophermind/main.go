package main

import (
	"GopherMind/internal/http/handler"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
	"os"
)

func main() {
	InitServer()
}

func InitServer() {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// 註冊測試 API
	r.Get("/health", handler.HealthCheck)

	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "8080"
	}
	fmt.Printf("🚀 GopherMind 啟動於埠號 %s\n", port)

	err := http.ListenAndServe(":"+port, r)
	if err != nil {
		panic(fmt.Sprintf("服務啟動失敗: %s", err))
	}
}

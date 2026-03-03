package main

import (
	"GopherMind/internal/app"
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	// 1. 建立一個可取消的 Context，用於優雅關閉服務
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// 2. 初始化 App (包含 Config 讀取與 DI 注入)
	application, err := app.New(ctx)
	if err != nil {
		// 使用 log.Fatalf 確保初始化失敗時立即停止並記錄原因
		log.Fatalf("❌ GopherMind 初始化失敗: %v", err)
	}

	// 3. 在一個獨立的 Goroutine 中啟動服務，這樣主執行緒才能監聽系統訊號 (如 Ctrl+C 或 Heroku 的重啟訊號)
	go func() {
		if err := application.Run(); err != nil {
			log.Fatalf("❌ 服務運行異常: %v", err)
		}
	}()

	// 4. 等待關閉訊號 (Graceful Shutdown 準備)
	<-ctx.Done()
	fmt.Println("\n\n🛑 接收到關閉訊號，GopherMind 正在安全退出...")
	// 這裡可以預留 5 秒讓處理中的 Request 完成
	_, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	fmt.Println("👋 服務已停止。")
}

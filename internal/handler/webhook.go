package handler

import (
	"GopherMind/internal/platform/telegram"
	"context"
	"encoding/json"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"net/http"
	"time"
)

func TelegramHandler(tgBot *telegram.Bot) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var update tgbotapi.Update

		// 1. 解析 Telegram 傳來的 JSON
		if err := json.NewDecoder(r.Body).Decode(&update); err != nil {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		// 2. 建立一個獨立的背景 Context，並設定一個合理的超時時間（例如 60 秒），避免 API 卡住導致資源洩漏
		go func(u tgbotapi.Update) {
			bgCtx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
			defer cancel() // 確保任務結束或超時後釋放資源

			tgBot.HandleWebhook(bgCtx, &u)
		}(update) // 將 update 傳入，避免閉包 (Closure) 的變數覆蓋問題

		// 3. 關鍵：立刻回傳 200 OK 給 Telegram，如果不馬上回傳，Telegram 伺服器會以為我沒收到，並瘋狂重試
		w.WriteHeader(http.StatusOK)
	}
}

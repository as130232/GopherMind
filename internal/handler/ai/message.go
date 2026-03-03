package handler

import (
	"GopherMind/internal/interface"
	"encoding/json"
	"net/http"
	"strings"
)

// ChatRequest 定義了客戶端傳送過來的 JSON 格式
type ChatRequest struct {
	Message string `json:"message"` // 用戶輸入的原始文字
	UserID  string `json:"userId"`  // 用於未來識別用戶身份或獲取 Context
}

// ChatResponse 定義了回傳給客戶端的格式
type ChatResponse struct {
	Reply  interface{} `json:"reply"`
	Source string      `json:"source"` // 用於識別是 Gemini 還是 OpenAI 回覆
}

func MessageHandler(ai _interface.AiProvider) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req ChatRequest

		// 1. 解析 JSON 請求
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		// 2. 根據 AI Provider (Gemini)，向AI傳送訊息
		reply, err := ai.Chat(r.Context(), "", req.Message)
		if err != nil {
			// 檢查是否為 429 錯誤
			if strings.Contains(err.Error(), "429") {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusTooManyRequests)
				err = json.NewEncoder(w).Encode(map[string]string{
					"error": "AI 目前太累了，請在 40 秒後再試一次",
					"code":  "RATE_LIMIT_EXCEEDED",
				})
				if err != nil {
					return
				}
				return
			}
			http.Error(w, "AI generate failed", http.StatusInternalServerError)
			return
		}

		// 3. 回傳結果
		resp := ChatResponse{
			Reply:  reply,
			Source: "GopherMind-" + ai.Name(),
		}

		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(resp)
		if err != nil {
			return
		}
	}
}

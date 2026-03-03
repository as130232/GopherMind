package telegram

import (
	_interface "GopherMind/internal/interface"
	"GopherMind/internal/model"
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"io"
	"log"
	"net/http"
	"strconv"
)

type Bot struct {
	api *tgbotapi.BotAPI
	ai  _interface.AiProvider
}

func NewBot(token string, ai _interface.AiProvider) *Bot {
	api, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Fatalf("❌ Telegram Bot 初始化失敗: %v", err)
	}
	return &Bot{api: api, ai: ai}
}

// HandleWebhook 接收從外層 Handler 傳進來的 Update 物件
func (b *Bot) HandleWebhook(ctx context.Context, update *tgbotapi.Update) {
	// 如果不是一般文字或圖片訊息，先忽略
	if update.Message == nil {
		return
	}
	b.handleMessage(ctx, update.Message)
}

func (b *Bot) handleMessage(ctx context.Context, msg *tgbotapi.Message) {
	userID := msg.Chat.ID
	// 取得使用者輸入的文字 (Telegram 的圖片附帶文字為Caption)
	text := msg.Text
	if msg.Caption != "" {
		text = msg.Caption
	}

	// 判斷使用者是否傳送了圖片
	if msg.Photo != nil && len(msg.Photo) > 0 {
		b.api.Send(tgbotapi.NewMessage(userID, "📥 收到圖片，正在思考中..."))
		// 取得最高畫質的圖片 FileID (陣列的最後一個)
		bestPhoto := msg.Photo[len(msg.Photo)-1]
		// 取得圖片下載連結並下載成 Byte 陣列
		imageBytes, err := b.downloadFile(bestPhoto.FileID)
		if err != nil {
			b.api.Send(tgbotapi.NewMessage(userID, "❌ 圖片下載失敗"))
			return
		}
		// 動態偵測真實的 MIME Type
		mimeType := http.DetectContentType(imageBytes)
		// 把圖片緩存至 Context 裡
		attachment := &model.FileAttachment{MimeType: mimeType, Data: imageBytes}
		ctx = context.WithValue(ctx, model.ImageAttachmentKey, attachment)
	}
	// AI 會去判斷 text 內容。如果是問股價就叫 StockTool，如果是問熱量就叫 FoodTool
	reply, err := b.ai.ChatByTool(ctx, strconv.FormatInt(userID, 10), text)
	if err != nil {
		_, err = b.api.Send(tgbotapi.NewMessage(userID, "❌ 大腦當機了: "+err.Error()))
		if err != nil {
			return
		}
		return
	}
	// 5. 將最終結果回傳給 Telegram
	finalMsg := tgbotapi.NewMessage(userID, reply)
	finalMsg.ParseMode = "Markdown"
	_, err = b.api.Send(finalMsg)
	if err != nil {
		return
	}
}

// 根據 FileID 下載 Telegram 伺服器上的圖片
func (b *Bot) downloadFile(fileID string) ([]byte, error) {
	file, err := b.api.GetFile(tgbotapi.FileConfig{FileID: fileID})
	if err != nil {
		return nil, err
	}

	url := file.Link(b.api.Token)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}

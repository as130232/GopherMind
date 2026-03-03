# GopherMind
A smart AI Agent powered by Go, transforming static commands into semantic workflows with rich API integrations.

GopherMind 是一個基於 Go 語言開發的智慧型 AI Agent。它旨在將傳統的「硬編碼指令」進化為「語意理解驅動」的自動化工作流。

## 🌟 核心理念
過去我們習慣透過特定關鍵字觸發功能，而 **GopherMind** 透過大型語言模型 (LLM) 與向量資料庫 (Qdrant)，能理解用戶的自然語言意圖，並自動調度對應的工具來完成任務。

## 🛠 關鍵技術
- **Language:** Go (Golang)
- **AI Engine:** Gemini/ OpenAI / Anthropic / Ollama (Function Calling)
- **Vector DB:** Qdrant (用於語意路由與長期記憶)
- **Messaging:** LineBot SDK
- **Data:** Redis, PostgreSQL, Kafka/RabbitMQ

## 🚀 已整合功能 (Tools)
GopherMind 整合了多樣化的生活與開發工具：
- 📈 **金融資訊**：即時股價、殖利率、匯率換算、三大法人統計。
- 🌦 **環境感知**：天氣預報與降雨提醒、即時油價查詢。
- 🏢 **辦公自動化**：Femas 打卡時間管理與下班提醒。
- 🔍 **社群情報**：PTT 各大版面文章與圖片資源爬取。
- 🌐 **語言工具**：多國語言即時翻譯。

## 📂 專案架構
本專案遵循標準 Go 專案結構，將 AI 決策層 (Agent) 與執行層 (Tools) 完美解耦，便於快速擴充新功能。

## ⚙️ 快速開始
1. **複製專案**
   ```bash
   git clone [https://github.com/as130232/GopherMind.git](https://github.com/as130232/GopherMind.git)

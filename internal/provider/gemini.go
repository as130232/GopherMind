package provider

import (
	"GopherMind/consts"
	"GopherMind/internal/interface"
	"GopherMind/internal/interface/repo"
	"GopherMind/internal/model"
	"GopherMind/internal/repository"
	"context"
	"errors"
	"fmt"
	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
	"log"
)

type GeminiProvider struct {
	client          *genai.Client
	model           *genai.GenerativeModel
	tools           map[string]_interface.BaseTool
	chatHistoryRepo *repository.ChatHistoryRepository
}

var (
	//ModelName = "gemini-2.0-flash"
	ModelName = "gemini-2.5-flash"
)

func NewGeminiProvider(ctx context.Context, apiKey, geminiModel string,
	chatHistoryRepo *repository.ChatHistoryRepository) (*GeminiProvider, error) {
	if len(geminiModel) != 0 {
		ModelName = geminiModel
	}
	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return nil, err
	}
	// 檢查現階段提供哪些模型
	iter := client.ListModels(ctx)
	for {
		m, err := iter.Next()
		if errors.Is(err, iterator.Done) {
			break
		}
		if err != nil {
			log.Printf("列出模型失敗: %v", err)
			break
		}
		fmt.Printf("Gemini可用模型: %s, 支援方法: %v\n", m.Name, m.SupportedGenerationMethods)
	}
	clientModel := client.GenerativeModel(ModelName)
	return &GeminiProvider{
		client:          client,
		model:           clientModel,
		chatHistoryRepo: chatHistoryRepo,
	}, nil
}
func (g *GeminiProvider) Name() string {
	return consts.TypeGemini
}

func (g *GeminiProvider) RegisterTools(tools map[string]_interface.BaseTool) {
	var fns []*genai.FunctionDeclaration
	g.tools = tools
	for _, t := range tools {
		meta := t.GetMetadata()
		// 將中立的 model.ToolMetadata 轉換為 genai 的格式
		properties := make(map[string]*genai.Schema)
		var required []string

		for name, detail := range meta.Parameters {
			var schemaType genai.Type
			switch detail.Type {
			case "object":
				schemaType = genai.TypeObject
			case "number":
				schemaType = genai.TypeNumber
			case "boolean":
				schemaType = genai.TypeBoolean
			default:
				schemaType = genai.TypeString
			}
			properties[name] = &genai.Schema{
				Type:        schemaType,
				Description: detail.Description,
			}
			if detail.Required {
				required = append(required, name)
			}
		}

		fns = append(fns, &genai.FunctionDeclaration{
			Name:        meta.Name,
			Description: meta.Description,
			Parameters: &genai.Schema{
				Type:       genai.TypeObject,
				Properties: properties,
				Required:   required,
			},
		})
	}
	g.model.Tools = []*genai.Tool{{FunctionDeclarations: fns}}
	g.model.SystemInstruction = &genai.Content{
		Parts: []genai.Part{genai.Text(`
	   你是一個全方位的智能助理 Go通寶。
	   1. 股市分析：當用戶詢問股票時，優先呼叫 get_stock_price 工具。
	   2. 飲食管理：當用戶上傳食物照片、詢問熱量、營養素或是否能吃時，請呼叫 analyze_food_nutrition 工具。
	   3. 一般閒聊：如果用戶只是打招呼或問與上述無關的問題，請直接用親切的語氣回覆，不用呼叫工具。
	`)},
	}
}

func (g *GeminiProvider) Chat(ctx context.Context, userID, message string) (string, error) {
	clientModel := g.client.GenerativeModel(ModelName)
	resp, err := clientModel.GenerateContent(ctx, genai.Text(message))
	if err != nil {
		fmt.Printf("Gemini GenerateContent err: %+v", err)
		return "", err
	}
	if len(resp.Candidates) > 0 {
		return fmt.Sprintf("%v", resp.Candidates[0].Content.Parts[0]), nil
	}
	return "No response", nil
}

func (g *GeminiProvider) AnalyzeImage(ctx context.Context, mimeType string, imgBytes []byte, prompt string) (string, error) {
	clientModel := g.client.GenerativeModel(ModelName)
	// 強制要求模型回傳 JSON 格式
	clientModel.ResponseMIMEType = "application/json"
	promptPart := genai.Text(prompt)
	imgPart := genai.Blob{
		MIMEType: mimeType,
		Data:     imgBytes,
	}
	resp, err := clientModel.GenerateContent(ctx, promptPart, imgPart)
	if err != nil {
		fmt.Printf("Gemini GenerateContent err: %+v", err)
		return "", err
	}
	if len(resp.Candidates) > 0 && len(resp.Candidates[0].Content.Parts) > 0 {
		if textPart, ok := resp.Candidates[0].Content.Parts[0].(genai.Text); ok {
			return fmt.Sprintf("%v", textPart), nil
		}
	}
	return "No response", nil
}

func (g *GeminiProvider) ChatByTool(ctx context.Context, userID, message string) (string, error) {
	// 檢查 Context 中是否被 Telegram (或未來 Line) 塞入了圖片
	if attachment, ok := ctx.Value(model.ImageAttachmentKey).(*model.FileAttachment); ok && attachment != nil {
		// 偷偷修改發給 AI 的文字，讓它知道有圖片存在
		message = fmt.Sprintf("【系統提示：此請求已附帶一張圖片 (MIME: %s)，安全存放於系統記憶體中待工具讀取】\n用戶原始訊息：%s", attachment.MimeType, message)
	}
	// 1. 建立乾淨的模型副本 (併發安全)，避免並發請求互相干擾 SystemInstruction，需再把在註冊時準備好的 Tools 和指令塞進去
	clientModel := g.client.GenerativeModel(ModelName)
	clientModel.Tools = g.model.Tools
	clientModel.SystemInstruction = g.model.SystemInstruction

	// 讀取並載入歷史紀錄
	history, _ := g.chatHistoryRepo.GetHistory(ctx, userID)
	var geminiHistory []*genai.Content
	for _, h := range history {
		geminiHistory = append(geminiHistory, &genai.Content{
			Role:  h.Role, // "user" 或 "model"
			Parts: []genai.Part{genai.Text(h.Message)},
		})
	}
	session := clientModel.StartChat()
	if len(geminiHistory) > 0 {
		session.History = geminiHistory // 載入用戶的歷史記憶
	}
	resp, err := session.SendMessage(ctx, genai.Text(message))
	if err != nil {
		return "", err
	}

	part := resp.Candidates[0].Content.Parts[0]

	funcCall, ok := part.(genai.FunctionCall)
	if !ok {
		// AI 決定回覆文字
		return fmt.Sprintf("%v", part), nil
	}

	fmt.Printf("🎯 AI 呼叫工具: %s\n", funcCall.Name)

	// 執行工具 (此時 FoodTool 會從 ctx 中拿出圖片)
	tool := g.tools[funcCall.Name]
	toolResult, err := tool.Execute(ctx, funcCall.Args)
	if err != nil {
		return "", fmt.Errorf("工具執行失敗: %w", err)
	}
	finalReply := ""
	if tool.GetMetadata().IsOnce { //是否將第一次結果直接回傳
		finalReply = fmt.Sprintf("%v", toolResult)
	} else {
		// 將 Service 分析完的結果餵回給 AI
		secondResp, err := session.SendMessage(ctx, genai.FunctionResponse{
			Name:     funcCall.Name,
			Response: map[string]any{"content": toolResult},
		})
		if err != nil {
			return "", err
		}
		finalReply = fmt.Sprintf("%v", secondResp.Candidates[0].Content.Parts[0])
	}
	// 緩存本次對話
	g.chatHistoryRepo.SaveHistory(ctx, userID, repo.ChatHistory{Role: "user", Message: message})
	g.chatHistoryRepo.SaveHistory(ctx, userID, repo.ChatHistory{Role: "model", Message: finalReply})
	return finalReply, nil
}

func (g *GeminiProvider) GetStockInfo(ctx context.Context, message string) (string, error) {
	clientModel := g.client.GenerativeModel(ModelName)
	clientModel.SystemInstruction = &genai.Content{
		Parts: []genai.Part{genai.Text(`
	   你是一個專業的股市助理。當用戶詢問股票資訊時：
	   1. 優先嘗試將公司名稱轉換為股票代碼（例如：台積電 -> 2330, 鴻海 -> 2317）。
	   2. 使用 get_stock_info 工具來獲取即時數據。
	   3. 除非你完全無法確定代碼，否則不要詢問用戶，請直接呼叫工具。
	`)},
	}

	session := clientModel.StartChat()
	resp, err := session.SendMessage(ctx, genai.Text(message))
	if err != nil {
		return "", err
	}

	part := resp.Candidates[0].Content.Parts[0]

	funcCall, ok := part.(genai.FunctionCall)
	if !ok {
		// 🟢 成功路徑：AI 回傳的是文字，代表思考結束
		return fmt.Sprintf("%v", part), nil
	}

	// 執行 Go 工具
	toolResult, err := g.tools[funcCall.Name].Execute(ctx, funcCall.Args)
	if err != nil {
		// 🔴 錯誤處理：如果工具執行失敗，把錯誤餵給 AI 讓它知道，或是直接結束
		return "", fmt.Errorf("工具執行異常且已停止: %w", err)
	}

	// 將結果丟回 AI，讓它進行下一次思考
	secondResp, err := session.SendMessage(ctx, genai.FunctionResponse{
		Name:     funcCall.Name,
		Response: map[string]any{"content": toolResult},
	})
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%v", secondResp.Candidates[0].Content.Parts[0]), nil
}

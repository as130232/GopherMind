package food

import (
	_interface "GopherMind/internal/interface"
	"GopherMind/internal/model"
	"GopherMind/internal/service"
	"context"
	"fmt"
)

type AnalyzeFoodNutritionTool struct {
	foodSvc *service.FoodService
}

func NewAnalyzeFoodNutritionTool(svc *service.FoodService) *AnalyzeFoodNutritionTool {
	return &AnalyzeFoodNutritionTool{foodSvc: svc}
}

func (t *AnalyzeFoodNutritionTool) GetMetadata() _interface.ToolMetadata {
	return _interface.ToolMetadata{
		Name:        "analyze_food_nutrition",
		Description: "當使用者上傳食物照片，並要求分析熱量、營養素、卡路里或判斷能不能吃時，呼叫此工具。系統會自動讀取使用者最新上傳的圖片進行分析。",
		Parameters:  map[string]_interface.ToolParameter{}, // 不需要參數，因為圖片在 Context 裡
		IsOnce:      true,
	}
}

func (t *AnalyzeFoodNutritionTool) Execute(ctx context.Context, args map[string]interface{}) (interface{}, error) {
	// 從 Context 中取出圖片附件
	attachment, ok := ctx.Value(model.ImageAttachmentKey).(*model.FileAttachment)
	if !ok || attachment == nil {
		// 如果沒有圖片，請 AI 回覆提醒使用者
		return "分析失敗：找不到圖片。提醒使用者必須先上傳一張食物的照片。", nil
	}
	// 分析食物
	resp, err := t.foodSvc.AnalyzeFoodNutrition(ctx, attachment.MimeType, attachment.Data)
	if err != nil {
		return nil, fmt.Errorf("分析過程發生錯誤: %w", err)
	}

	replyText := fmt.Sprintf(
		"🍱 **%s** 營養分析\n\n"+
			"🔥 總熱量: %d kcal\n"+
			"🥩 蛋白質: %d g\n"+
			"🍚 碳水: %d g\n"+
			"🥑 脂肪: %d g\n\n"+
			"💡 分析:\n%s",
		resp.Name, resp.Calories, resp.Protein, resp.Carbs, resp.Fat, resp.Analysis,
	)
	return replyText, nil
}

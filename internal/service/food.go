package service

import (
	"GopherMind/internal/interface"
	"context"
	"encoding/json"
	"fmt"
)

type FoodService struct {
	ai _interface.AiProvider
}

func NewFoodService(ai _interface.AiProvider) *FoodService {
	return &FoodService{ai: ai}
}

// FoodNutrition 定義你的系統要的熱量資料結構
type FoodNutrition struct {
	Name        string   `json:"name"`
	Ingredients []string `json:"ingredients"`
	Calories    int      `json:"calories"`
	Protein     int      `json:"protein"`
	Fat         int      `json:"fat"`
	Carbs       int      `json:"carbs"`
	Analysis    string   `json:"analysis"`
}

func (s *FoodService) AnalyzeFoodNutrition(ctx context.Context, mimeType string, imageBytes []byte) (*FoodNutrition, error) {
	prompt := `請作為一位專業的運動營養師，分析這張圖片中的食物。
		請精確估算熱量與巨量營養素，並回傳純 JSON 格式，若是多個餐點，則幫我集合成同一個JSON物件，必須包含以下欄位：
		- name (字串): 食物名稱
		- ingredients (字串陣列): 辨識出的主要食材
		- calories (正整數): 預估總熱量 (大卡)
		- protein (正整數): 預估總蛋白質 (克)
		- fat (正整數): 預估總脂肪 (克)
		- carbs (正整數): 預估總碳水化合物 (克)
		- analysis (字串): 簡短說明你的估算依據，或指出圖片中無法精準判斷的部分`
	resp, err := s.ai.AnalyzeImage(ctx, mimeType, imageBytes, prompt)
	if err != nil {
		return nil, fmt.Errorf("AI 分析圖片失敗: %w", err)
	}
	var nutrition FoodNutrition
	if err := json.Unmarshal([]byte(resp), &nutrition); err != nil {
		return nil, fmt.Errorf("解析資料失敗: err: %w, result: %s", err, string(resp))
	}
	return &nutrition, nil
}

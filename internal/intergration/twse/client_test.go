package twse

import (
	"context"
	"testing"
	"time"
)

func TestFetchStockPrice_Real(t *testing.T) {
	client := NewClient()

	stockID := "2330" // 台積電
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	t.Logf("正在測試抓取股票 %s 的即時資料...", stockID)
	result, err := client.FetchStockInfo(ctx, stockID)

	// 斷言結果 (Assertions)
	if err != nil {
		t.Fatalf("❌ 測試失敗：無法抓取資料, 錯誤: %v", err)
	}

	if result == nil {
		t.Fatalf("❌ 測試失敗：回傳結果不應為空")
	}

	if result.Code != stockID {
		t.Errorf("❌ 預期代碼 2330, 得到 %s", result.Code)
	}

	t.Logf("✅ 測試成功！股票名稱: %s (%s)\n當前價格: %s\n昨收價: %s\n今日最高: %s\n今日最低: %s\n漲幅: %.2f%%\n成交量: %s 股",
		result.Name, result.Code, result.Close, result.Yesterday, result.High, result.Low, result.GetIncrease(), result.Volume)
}

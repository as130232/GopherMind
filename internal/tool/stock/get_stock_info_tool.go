package tool

import (
	"GopherMind/internal/interface"
	"GopherMind/internal/intergration/twse"
	"context"
	"fmt"
)

type GetStockInfoTool struct {
	twseClient *twse.Client
}

func NewGetStockInfo(twseClient *twse.Client) *GetStockInfoTool {
	return &GetStockInfoTool{twseClient: twseClient}
}

func (t *GetStockInfoTool) GetMetadata() _interface.ToolMetadata {
	return _interface.ToolMetadata{
		Name:        "get_stock_info",
		Description: "獲取台灣股市最新價格。參數 stock_id 支援數字代碼(如2330)或知名公司名稱(如台積電)",
		Parameters: map[string]_interface.ToolParameter{
			"stock_id": {
				Type:        "string",
				Description: "股票代碼，例如：2330",
				Required:    true,
			},
		},
	}
}

func (t *GetStockInfoTool) Execute(ctx context.Context, args map[string]interface{}) (interface{}, error) {
	stockID, ok := args["stock_id"].(string)
	if !ok {
		return nil, fmt.Errorf("missing argument: stock_id")
	}
	info, err := t.twseClient.FetchStockInfo(ctx, stockID)
	if err != nil {
		return nil, err
	}
	// 提供豐富的上下文，讓 AI 能根據漲跌幅、成交量進行判斷
	return fmt.Sprintf(
		"數據獲取成功！該股票為 %s (%s)。目前成交價 %s 元，今日漲跌幅為 %.2f%%，成交量為 %s。"+
			"今日開盤 %s，最高 %s，最低 %s。請根據這些數據為用戶提供簡短分析。",
		info.Name, info.Code, info.Close, info.GetIncrease(), info.Volume,
		info.Open, info.High, info.Low,
	), nil
}

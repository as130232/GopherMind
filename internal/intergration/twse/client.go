package twse

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"time"
)

type Client struct {
	restyClient *resty.Client
	baseURL     string
}

func NewClient() *Client {
	c := resty.New()
	c.
		SetBaseURL(baseURL).
		SetTimeout(5 * time.Second).
		SetRetryCount(3).
		SetRetryWaitTime(200 * time.Millisecond).
		SetRetryMaxWaitTime(2 * time.Second).
		AddRetryCondition(func(r *resty.Response, err error) bool {
			if err != nil {
				return true
			}
			return r.StatusCode() >= 500
		})

	return &Client{
		restyClient: c,
	}
}

const (
	baseURL = "https://mis.twse.com.tw/stock/api"
)

func (c *Client) FetchStockInfo(ctx context.Context, stockID string) (*StockInfo, error) {
	// 這裡使用你熟悉的證交所 API 格式
	// 例如：https://mis.twse.com.tw/stock/api/getStockInfo.jsp?ex_ch=tse_2330.tw
	// "https://www.twse.com.tw/exchangeReport/STOCK_DAY_AVG_ALL"
	//fullURL := fmt.Sprintf("%s/getStockInfo.jsp", c.baseURL)

	resp, err := c.restyClient.R().
		SetContext(ctx).
		SetQueryParam("ex_ch", fmt.Sprintf("tse_%s.tw", stockID)).
		Get("/getStockInfo.jsp")
	if err != nil {
		return nil, fmt.Errorf("twse api request failed: %w", err)
	}

	var result Response
	if err = json.Unmarshal(resp.Body(), &result); err != nil {
		return nil, fmt.Errorf("JSON 解析失敗: %w, 原始內容: %s", err, resp.String())
	}

	if len(result.MsgArray) == 0 {
		return nil, fmt.Errorf("找不到股票代碼: %s", stockID)
	}
	return &result.MsgArray[0], nil
}

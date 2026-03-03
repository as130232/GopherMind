package twse

import (
	"math"
	"strconv"
)

// Response 證交所 API 的外層包裝
type Response struct {
	MsgArray  []StockInfo `json:"msgArray"`
	Referer   string      `json:"referer"`
	UserDelay int         `json:"userDelay"`
	Rtcode    string      `json:"rtcode"`
	QueryTime struct {
		SysDate           string `json:"sysDate"`
		StockInfoItem     int    `json:"stockInfoItem"`
		StockInfo         int    `json:"stockInfo"`
		SessionStr        string `json:"sessionStr"`
		SysTime           string `json:"sysTime"`
		ShowChart         bool   `json:"showChart"`
		SessionFromTime   int    `json:"sessionFromTime"`
		SessionLatestTime int    `json:"sessionLatestTime"`
	} `json:"queryTime"`
	Rtmessage   string `json:"rtmessage"`
	ExKey       string `json:"exKey"`
	CachedAlive int    `json:"cachedAlive"`
}

// StockInfo 對應你的 Java MsgPO
type StockInfo struct {
	Code     string `json:"c"`  // 股票編號
	FullCode string `json:"ch"` // 股票編號.tw
	Exchange string `json:"ex"` // 證交所代號: tse
	Name     string `json:"n"`  // 股票中文名稱: 台積電
	FullName string `json:"nf"` // 公司詳細名稱: 台灣積體電路製造股份有限公司

	Open      string `json:"o"` // 開盤價(09:00時的價格)
	High      string `json:"h"` // 當日最高價
	Low       string `json:"l"` // 當日最低價
	Close     string `json:"z"` // 現價or收盤價，有時會給"-"
	Yesterday string `json:"y"` // 平盤價 (昨收)
	LimitUp   string `json:"u"` // 漲停價
	LimitDown string `json:"w"` // 跌停價

	Date     string `json:"d"`     // 日期:20210510
	Time     string `json:"t"`     // 當日盤中時間:13:30:00
	OverTime string `json:"ot"`    // 當日盤後時間:14:30:00
	Tlong    string `json:"tlong"` //當日盤後時間轉為毫秒

	Volume string `json:"v"`  //累積總交易量(13:30以前不包含14:30)
	Tv     string `json:"tv"` //單量-盤中結束時交易量(13:30)
	Fv     string `json:"fv"` //單量-盤後時交易量(14:30)

	A string `json:"a"` //五檔賣:1920.0000_1925.0000_1930.0000_1935.0000_1940.0000_
	B string `json:"b"` //五檔買:1915.0000_1910.0000_1905.0000_1900.0000_1895.0000_
	F string `json:"f"` //五檔賣對應委託數量:429_989_1558_1020_631_
	G string `json:"g"` //五檔買對應委託數量:249_484_419_2828_393_

	Field18 string `json:"%"`
	Field1  string `json:"@"`
	Field11 string `json:"^"`
	Field16 string `json:"#"`
	Ps      string `json:"ps"`
	Pid     string `json:"pid"`
	Pz      string `json:"pz"`
	Bp      string `json:"bp"`
	Oa      string `json:"oa"`
	Ob      string `json:"ob"`
	M       string `json:"m%"`
	Key     string `json:"key"`
	Ip      string `json:"ip"`
	Mt      string `json:"mt"`
	Ov      string `json:"ov"`
	I       string `json:"i"`
	It      string `json:"it"`
	Oz      string `json:"oz"`
	P       string `json:"p"`
	Ts      string `json:"ts"`
}

// GetIncrease 計算漲幅 (收盤 - 平盤) / 平盤 * 100
func (s *StockInfo) GetIncrease() float64 {
	closePrice, errZ := strconv.ParseFloat(s.Close, 64)
	yestPrice, errY := strconv.ParseFloat(s.Yesterday, 64)

	if errZ != nil || errY != nil || yestPrice == 0 {
		return -1.0
	}

	result := ((closePrice - yestPrice) / yestPrice) * 100
	return math.Round(result*100) / 100 // 四捨五入到第二位
}

// GetAmplitude 計算振幅 (最高 - 最低) / 平盤 * 100
func (s *StockInfo) GetAmplitude() float64 {
	high, errH := strconv.ParseFloat(s.High, 64)
	low, errL := strconv.ParseFloat(s.Low, 64)
	yest, errY := strconv.ParseFloat(s.Yesterday, 64)

	if errH != nil || errL != nil || errY != nil || yest == 0 {
		return -1.0
	}

	result := ((high - low) / yest) * 100
	return math.Round(result*100) / 100
}

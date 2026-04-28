package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// open-meteo API のレスポンス構造体
type WeatherResponse struct {
	Daily struct {
		Time              []string  `json:"time"`
		WeatherCode       []int     `json:"weathercode"`
		TemperatureMax    []float64 `json:"temperature_2m_max"`
		TemperatureMin    []float64 `json:"temperature_2m_min"`
		PrecipitationSum  []float64 `json:"precipitation_sum"`
		PrecipitationProb []int     `json:"precipitation_probability_max"`
	} `json:"daily"`
}

// WMO天気コードを日本語の説明に変換
func weatherCodeToDesc(code int) string {
	switch {
	case code == 0:
		return "快晴"
	case code == 1:
		return "おおむね晴れ"
	case code == 2:
		return "一部曇り"
	case code == 3:
		return "曇り"
	case code >= 45 && code <= 48:
		return "霧"
	case code >= 51 && code <= 55:
		return "霧雨"
	case code >= 61 && code <= 65:
		return "雨"
	case code >= 71 && code <= 75:
		return "雪"
	case code == 77:
		return "あられ"
	case code >= 80 && code <= 82:
		return "にわか雨"
	case code >= 85 && code <= 86:
		return "にわか雪"
	case code >= 95 && code <= 99:
		return "雷雨"
	default:
		return "不明"
	}
}

func main() {
	// 東京の緯度・経度
	url := "https://api.open-meteo.com/v1/forecast" +
		"?latitude=35.6895&longitude=139.6917" +
		"&daily=weathercode,temperature_2m_max,temperature_2m_min,precipitation_sum,precipitation_probability_max" +
		"&timezone=Asia%2FTokyo" +
		"&forecast_days=2"

	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("APIリクエスト失敗: %v\n", err)
		return
	}
	defer resp.Body.Close()

	var weather WeatherResponse
	if err := json.NewDecoder(resp.Body).Decode(&weather); err != nil {
		fmt.Printf("JSONデコード失敗: %v\n", err)
		return
	}

	if len(weather.Daily.Time) < 2 {
		fmt.Println("明日の天気データが取得できませんでした")
		return
	}

	// インデックス1が明日
	tomorrow := weather.Daily.Time[1]
	t, _ := time.Parse("2006-01-02", tomorrow)

	fmt.Println("=== 東京の明日の天気予報 ===")
	fmt.Printf("日付　　　　: %s (%s)\n", t.Format("2006年01月02日"), weekdayJP(t.Weekday()))
	fmt.Printf("天気　　　　: %s\n", weatherCodeToDesc(weather.Daily.WeatherCode[1]))
	fmt.Printf("最高気温　　: %.1f°C\n", weather.Daily.TemperatureMax[1])
	fmt.Printf("最低気温　　: %.1f°C\n", weather.Daily.TemperatureMin[1])
	fmt.Printf("降水量　　　: %.1f mm\n", weather.Daily.PrecipitationSum[1])
	fmt.Printf("降水確率　　: %d%%\n", weather.Daily.PrecipitationProb[1])
	fmt.Println("=========================")
	fmt.Println("=========================")
}

func weekdayJP(w time.Weekday) string {
	days := []string{"日", "月", "火", "水", "木", "金", "土"}
	return days[w] + "曜日"
}

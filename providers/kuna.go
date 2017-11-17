package providers

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"
)

const KunaTickerURL string = "https://kuna.io/api/v2/tickers/btcuah"

// GetTicker Make HTTP query and return ticker info
func GetKunaTicker() (Ticker, error) {
	resp, err := http.Get(KunaTickerURL)
	if err != nil {
		log.Fatal("Received error", err.Error())
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err.Error())
	}
	return getKunaResp(body)
}

type kunaTickerData struct {
	Buy   string `json:"buy"`
	Sell  string `json:"sell"`
	Low   string `json:"low"`
	High  string `json:"high"`
	Last  string `json:"last"`
	Vol   string `json:"vol"`
	Price string `json:"price"`
}

type kunaTickerResponse struct {
	At     int64          `json:"at"`
	Ticker kunaTickerData `json:"ticker"`
}

// Parse JSON API body response to KunaTicker
func getKunaResp(body []byte) (Ticker, error) {
	var s = new(kunaTickerResponse)
	err := json.Unmarshal(body, &s)
	if err != nil {
		log.Fatalln("whoops:", err)
	}
	return getTickerFromResponse(s)
}

func getTickerFromResponse(response *kunaTickerResponse) (Ticker, error) {
	at := time.Unix(response.At, 0)
	buy, err := strconv.ParseFloat(response.Ticker.Buy, 64)
	if err != nil {
		return nil, err
	}
	sell, err := strconv.ParseFloat(response.Ticker.Sell, 64)
	if err != nil {
		return nil, err
	}
	ticker := newKunaTicker(at, buy, sell)
	return ticker, nil
}

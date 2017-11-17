package providers

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"
)

const CoinbaseBuyURL string = "https://api.coinbase.com/v2/prices/BTC-USD/buy"
const CoinbaseSellURL string = "https://api.coinbase.com/v2/prices/BTC-USD/sell"

func GetCoinbaseTicker() (Ticker, error) {
	buy, err := loadData(CoinbaseBuyURL)
	if err != nil {
		log.Panic(err)
	}
	sell, err := loadData(CoinbaseSellURL)
	if err != nil {
		log.Panic(err)
	}
	now := time.Now()
	ticker := newCoinbaseTicker(now, buy, sell)
	return ticker, nil
}

type dataStruct struct {
	Base     string `json:"base"`
	Currency string `json:"currency"`
	Amount   string `json:"amount"`
}

type responseStruct struct {
	Data dataStruct `json:"data"`
}

func loadData(url string) (float64, error) {
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal("Received error", err.Error())
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err.Error())
	}
	return getAmount(body)
}

func getAmount(body []byte) (float64, error) {
	var s = new(responseStruct)
	err := json.Unmarshal(body, &s)
	if err != nil {
		log.Panic(err)
	}
	amount, err := strconv.ParseFloat(s.Data.Amount, 64)
	if err != nil {
		log.Panic(err)
	}
	return amount, nil
}

package providers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"
)

const CoinbaseBuyURL string = "https://api.coinbase.com/v2/prices/BTC-USD/buy"
const CoinbaseSellURL string = "https://api.coinbase.com/v2/prices/BTC-USD/sell"
const CoinbaseSpotURL string = "https://api.coinbase.com/v2/prices/BTC-USD/spot"

func GetCoinbaseTicker() (Ticker, error) {
	buy, err := loadData(CoinbaseBuyURL)
	if err != nil {
		panic(err)
	}
	sell, err := loadData(CoinbaseSellURL)
	if err != nil {
		panic(err)
	}
	now := time.Now()
	ticker := newCoinbaseTicker(now, buy, sell)
	return ticker, nil
}

func GetCoinbaseSpotPrice(date time.Time) float64 {
	url := fmt.Sprintf("%s?date=%d-%d-%d", CoinbaseSpotURL, date.Year(), date.Month(), date.Day())
	price, err := loadData(url)
	if err != nil {
		panic(err)
	}
	return price
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
	log.Output(1, fmt.Sprintf("Coinbase Request: %s\n", url))
	resp, err := http.Get(url)
	if err != nil {
		log.Panic("Received error ", err.Error())
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
		panic(err)
	}
	amount, err := strconv.ParseFloat(s.Data.Amount, 64)
	if err != nil {
		panic(err)
	}
	return amount, nil
}

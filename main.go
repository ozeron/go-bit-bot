package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/ozeron/go-bit-bot/providers"

	"github.com/yanzay/tbot"
)

var balance float64 = 0.0678
var exchangeRate float64 = 5000

func main() {
	botApiKey := os.Getenv("BIT_BOT_TOKEN")
	// Create new telegram bot server using token
	bot, err := tbot.NewServer(botApiKey)
	if err != nil {
		log.Fatal(err)
	}

	// Yo handler works without slash, simple text response
	bot.Handle("yo", "YO!")

	// Handle with HiHandler function
	bot.HandleFunc("/start", HiHandler)
	bot.HandleFunc("/add {amount}", AddHandler)
	bot.HandleFunc("/withdraw {amount}", WithdrawHandler)
	bot.HandleFunc("/exchange {amount}", ExchangeRateHandler)
	bot.HandleFunc("/wallet {address}", WalletRateHandler)

	bot.HandleFunc("/substrace", HiHandler)

	// Handler can accept varialbes
	// bot.HandleFunc("/say {text}", SayHandler)
	// Bot can send stickers, photos, music
	// bot.HandleFunc("/sticker", StickerHandler)
	// bot.HandleFunc("/photo", PhotoHandler)
	// bot.HandleFunc("/keyboard", KeyboardHandler)

	// Use file handler to handle user uploads
	// bot.HandleFile(FileHandler)

	// Set default handler if you want to process unmatched input
	bot.HandleDefault(EchoHandler)

	// Start listening for messages
	err = bot.ListenAndServe()
	log.Fatal(err)
}

func getTickerState() string {
	kunaTicker, err := providers.GetKunaTicker()
	if err != nil {
		log.Fatal(err)
	}
	coinbaseTicker, err := providers.GetCoinbaseTicker()
	if err != nil {
		log.Fatal(err)
	}
	roi := Roi(coinbaseTicker.Sell(), balance)
	log.Output(1, fmt.Sprintf("Bal: %f EXCH: %F SELL: %f ROI: %f\n", balance, exchangeRate, coinbaseTicker.Sell(), roi))
	gain := exchangeRate * balance * (1 + roi)
	return fmt.Sprintf("State:\n%s\n%s\nBalance: %.8f BTC\nROI: %.f%% Gain: %.f$", kunaTicker.String(), coinbaseTicker.String(), balance, 100*(1+roi), gain)
}

func AddHandler(message *tbot.Message) {
	stringAmount, ok := message.Vars["amount"]
	if !ok {
		log.Panicf("Received: %s, wrong amount", message.Data)
	}
	amount, error := strconv.ParseFloat(stringAmount, 32)
	if error != nil {
		log.Panic(error)
	}
	balance += amount
	message.Replyf("New balance: %.8f", balance)
}

func WithdrawHandler(message *tbot.Message) {
	stringAmount, ok := message.Vars["amount"]
	if !ok {
		log.Panicf("Received: %s, wrong amount", message.Data)
	}
	amount, error := strconv.ParseFloat(stringAmount, 32)
	if error != nil {
		log.Panic(error)
	}
	balance -= amount
	message.Replyf("New balance: %.8f", balance)
}

func ExchangeRateHandler(message *tbot.Message) {
	stringAmount, ok := message.Vars["amount"]
	if !ok {
		log.Panicf("Received: %s, wrong amount", message.Data)
	}
	amount, error := strconv.ParseFloat(stringAmount, 32)
	if error != nil {
		log.Panic(error)
	}
	exchangeRate = amount
	message.Replyf("Wallet buying exchange rate: %.3f$", exchangeRate)
}

func WalletRateHandler(message *tbot.Message) {
	address, ok := message.Vars["address"]
	if !ok {
		log.Panicf("Received: %s, wrong amount", message.Data)
	}
	wallet := providers.LoadWallet(address)
	invested := wallet.InvestedAmount()
	coinbaseTicker, err := providers.GetCoinbaseTicker()
	if err != nil {
		log.Fatal(err)
	}
	roi := Roi(coinbaseTicker.Sell(), invested)
	capital := invested * (1 + roi)
	response := fmt.Sprintf("%s\nROI: %.f%% Invested*: %.2f$ Capital*: %.2f$\n", wallet.String(), (1+roi)*100, invested, capital)
	message.Reply(response)
}

func Roi(sellPrice float64, balance float64) float64 {
	gainFromInvestment := sellPrice * balance
	costOfInvestment := exchangeRate * balance
	log.Output(2, fmt.Sprintf("Gain: %f, Cost: %f", gainFromInvestment, costOfInvestment))
	return (gainFromInvestment - costOfInvestment) / costOfInvestment
}

func HiHandler(message *tbot.Message) {
	// Handler can reply with several messages
	buttons := [][]string{
		{"/start", "/exchane"},
		{"/add", "/withdraw"},
	}
	message.Replyf("Hello, %s!", message.From)
	time.Sleep(1 * time.Second)
	message.ReplyKeyboard(getTickerState(), buttons)
}

func EchoHandler(message *tbot.Message) {
	message.Reply(fmt.Sprintf("Received: '%s'", message.Data))
	message.Reply("Start with /start command")
}

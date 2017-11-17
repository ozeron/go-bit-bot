package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/ozeron/go-bit-bot/providers"

	"github.com/yanzay/tbot"
)

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
	return fmt.Sprintf("State:\n%s\n%s\n", kunaTicker.String(), coinbaseTicker.String())
}

func WalletRateHandler(message *tbot.Message) {
	address, ok := message.Vars["address"]
	if !ok {
		log.Panicf("Received: %s, wrong amount", message.Data)
	}
	wallet := providers.LoadWallet(address)
	coinbaseTicker, err := providers.GetCoinbaseTicker()
	if err != nil {
		log.Fatal(err)
	}

	invested := wallet.InvestedAmount()
	roi := Roi(coinbaseTicker.Sell(), wallet)
	capital := invested * (1 + roi)

	response := fmt.Sprintf("%s\n%s\nROI: %.f%% Invested*: %.2f$ Capital*: %.2f$\n", coinbaseTicker.String(), wallet.String(), (1+roi)*100, invested, capital)
	message.Reply(response)
}

func Roi(sellPrice float64, wallet *providers.Wallet) float64 {
	costOfInvestment := wallet.InvestedAmount()
	gainFromInvestment := sellPrice * wallet.BTC()
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

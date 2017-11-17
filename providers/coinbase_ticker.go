package providers

import (
	"fmt"
	"time"
)

type coinbaseTicket struct {
	At   time.Time
	Buy  float64
	sell float64
}

func newCoinbaseTicker(t time.Time, buy float64, sell float64) *coinbaseTicket {
	return &coinbaseTicket{At: t, Buy: buy, sell: sell}
}

func (t coinbaseTicket) String() string {
	return fmt.Sprintf("Coinbase: %.f – %.f$", t.Buy, t.sell)
}

func (t coinbaseTicket) Sell() float64 {
	return t.sell
}

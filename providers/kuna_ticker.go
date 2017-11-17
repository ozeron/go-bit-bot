package providers

import (
	"fmt"
	"time"
)

type kunaTicker struct {
	At   time.Time
	Buy  float64
	sell float64
}

func newKunaTicker(at time.Time, buy float64, sell float64) *kunaTicker {
	return &kunaTicker{At: at, Buy: buy, sell: sell}
}

func (t kunaTicker) String() string {
	return fmt.Sprintf("Kuna: %.f – %.f UAH", t.Buy, t.sell)
}

func (t kunaTicker) Sell() float64 {
	return t.sell
}

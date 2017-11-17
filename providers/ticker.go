package providers

// basic ticker interface
type Ticker interface {
	String() string
	Sell() float64
}

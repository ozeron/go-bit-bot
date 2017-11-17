package providers_test

import (
	"testing"
	"time"

	"github.com/ozeron/go-bit-bot/base_provider"
	"github.com/ozeron/go-bit-bot/providers"
	"github.com/stretchr/testify/assert"
	httpmock "gopkg.in/jarcoal/httpmock.v1"
)

const buyResponse string = `{
	"data": {
		"base": "BTC",
		"currency": "USD",
		"amount": "7323.00"
	}
}`

const sellResponse string = `{
	"data": {
		"base": "BTC",
		"currency": "USD",
		"amount": "7182.32"
	}
}`

func TestGetCoinbaseTicker(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", providers.CoinbaseBuyURL,
		httpmock.NewStringResponder(200, buyResponse))

	httpmock.RegisterResponder("GET", providers.CoinbaseSellURL,
		httpmock.NewStringResponder(200, sellResponse))

	base := &base_provider.Ticker{
		At:   time.Unix(1510790630, 0),
		Buy:  7323,
		Sell: 7182,
	}

	ticker, err := providers.GetCoinbaseTicker()
	assert.Equal(t, err, nil, "error should be bil")
	assert.Equal(t, ticker, base, "tickers should be equal")
}

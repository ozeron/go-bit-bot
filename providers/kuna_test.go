package providers_test

import (
	"testing"

	"github.com/ozeron/go-bit-bot/providers"
	"github.com/stretchr/testify/assert"
	httpmock "gopkg.in/jarcoal/httpmock.v1"
)

const ticketResponse string = `
{
	"at": 1510790630,
	"ticker": {
		"buy": "196503.0",
		"sell": "198000.0",
		"low": "179500.0",
		"high": "198600.0",
		"last": "198000.0",
		"vol": "22.29281",
		"price": "4268959.161497"
	}
}
`

func TestGetKunaTicker(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", providers.KunaTickerURL,
		httpmock.NewStringResponder(200, ticketResponse))

	expected := ""

	ticker, err := providers.GetKunaTicker()
	assert.Equal(t, err, nil, "error should be bil")
	assert.Equal(t, ticker.String(), expected, "tickers should be equal")
}

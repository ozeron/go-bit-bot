package providers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

const Satoshi float64 = 100000000
const GetWalletHistoryURL string = "https://blockchain.info/rawaddr/"

type Transaction struct {
	spent        bool
	amount       int64
	Time         time.Time
	exchangeRate float64
}

func (t *Transaction) BTC() float64 {
	return float64(t.amount) / Satoshi
}

func (t *Transaction) String() string {
	var spent string
	if t.spent {
		spent = "📥"
	} else {
		spent = "📥"
	}
	return fmt.Sprintf("%s %s - %.8f BTC | %.f $", spent, t.Time.Format("2006 Jan 2"), t.BTC(), t.ExchangeRate())
}

func (t *Transaction) ExchangeRate() float64 {
	if t.exchangeRate == 0 {
		t.exchangeRate = t.loadPrice()
	}
	return t.exchangeRate
}

func (t *Transaction) loadPrice() float64 {
	return GetCoinbaseSpotPrice(t.Time)
}

type Wallet struct {
	address      string
	balance      int64
	transactions []Transaction
	invested     float64
}

func (w *Wallet) String() string {
	str := fmt.Sprintf("Balance: %.8f BTC\n", w.BTC())
	for i, _ := range w.transactions {
		t := &w.transactions[i]
		str += "\n"
		str += fmt.Sprintf("%s", t.String())
	}
	return str
}

func (w *Wallet) BTC() float64 {
	return float64(w.balance) / Satoshi
}

func (w *Wallet) InvestedAmount() float64 {
	if w.invested > 0 {
		return w.invested
	}
	w.invested = w.loadInvestedAmount()
	return w.invested
}

func (w *Wallet) loadInvestedAmount() float64 {
	var amount float64
	for i, _ := range w.transactions {
		t := &w.transactions[i]
		price := t.ExchangeRate()
		amount += float64(t.amount) * price / Satoshi
	}
	return amount
}

type outTransaction struct {
	Addr  string  `json:"addr"`
	Value float64 `json:"value"`
	Spent bool    `json:"boolean"`
}

type infoTransaction struct {
	Out  []outTransaction `json:"out"`
	Time float64          `json:"time"`
}

func (data *infoTransaction) findAndBuildTransaction(address string) (*Transaction, error) {
	transaction, err := data.findTransaction(address)
	if err != nil {
		return nil, err
	}
	time := time.Unix(int64(data.Time), 0)
	return &Transaction{Time: time, amount: int64(transaction.Value), spent: transaction.Spent}, nil
}

func (transaction *infoTransaction) findTransaction(address string) (*outTransaction, error) {
	for _, out := range transaction.Out {
		if out.Addr == address {
			return &out, nil
		}
	}
	return nil, fmt.Errorf("Transaction for address not found")
}

type walletInfo struct {
	Txs          []infoTransaction `json:"txs"`
	Address      string            `json:"address"`
	FinalBalance float64           `json:"final_balance"`
}

func (wallet *walletInfo) findTransactions(address string) ([]Transaction, error) {
	data := make([]Transaction, len(wallet.Txs))
	for i, info := range wallet.Txs {
		transaction, err := info.findAndBuildTransaction(address)
		if err != nil {
			return nil, err
		}
		data[i] = *transaction
	}
	return data, nil
}

func LoadWallet(address string) *Wallet {
	resp, err := http.Get(GetWalletHistoryURL + address)
	if err != nil {
		log.Fatal("Received error", err.Error())
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err.Error())
	}
	data := &walletInfo{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		panic(err.Error())
	}
	wallet := &Wallet{}
	wallet.balance = int64(data.FinalBalance)
	wallet.address = data.Address
	wallet.transactions, err = data.findTransactions(address)
	return wallet
}

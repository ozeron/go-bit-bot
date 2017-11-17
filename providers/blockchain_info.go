package providers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

const Satoshi int64 = 100000000
const GetWalletHistoryURL string = "https://blockchain.info/rawaddr/"

type Transaction struct {
	spent  bool
	amount int64
	Time   time.Time
}

func (t *Transaction) String() string {
	var spent string
	if t.spent {
		spent = "📥"
	} else {
		spent = "📥"
	}
	return fmt.Sprintf("%s %s - %d", spent, t.Time.Format("2006 Jan 2"), t.amount)
}

type Wallet struct {
	address      string
	balance      int64
	transactions []Transaction
}

func (w *Wallet) String() string {
	str := fmt.Sprintf("Address: %s\nBalance: %d BTC\n", w.address, w.BTC())
	for _, t := range w.transactions {
		str += "\n"
		str += fmt.Sprintf("%s", t.String())
	}
	return str
}

func (w *Wallet) BTC() float64 {
	return float64(w.balance) / float64(Satoshi)
}

func (w *Wallet) InvestedAmount() float64 {
	var amount float64
	for _, t := range w.transactions {
		price := GetCoinbaseSpotPrice(t.Time)
		amount += float64(t.amount) * price / float64(Satoshi)
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

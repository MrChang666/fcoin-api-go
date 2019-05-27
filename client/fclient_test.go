package client

import (
	"encoding/json"
	"fmt"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
	"testing"
)

func initClient() *FCoinClient {
	baseUrl := "https://api.fcoin.com/v2"
	assKey := ""
	secretKey := ""
	fc := NewFCoinClient(secretKey, assKey, baseUrl)
	return fc
}

func TestFcoinClient_GetBalance(t *testing.T) {
	fc := initClient()
	balance, err := fc.GetBalance()
	require.NoError(t, err)

	var total decimal.Decimal
	for _, v := range balance.Data {
		if res := statistic(v.Currency, v.Balance, fc); res.GreaterThan(decimal.Zero) {
			total = total.Add(res)
			fmt.Println(v.Currency, res)
		}
	}

	fmt.Println(total)
}

func statistic(symbol, balance string, fc *FCoinClient) decimal.Decimal {
	usdt := "usdt"
	var usBal decimal.Decimal
	usBal, _ = decimal.NewFromString(balance)

	if symbol != usdt {
		ticker, err := fc.GetLatestTickerBySymbol(symbol + usdt)
		if err != nil {
			fmt.Println(err)
		}
		if ticker != nil && len(ticker.Data.Ticker) > 0 {
			price := decimal.NewFromFloat(ticker.Data.Ticker[0])
			usBal = price.Mul(usBal)
		}

	}
	return usBal
}

func TestFcoinClient_GetUSDTBalance(t *testing.T) {
	fc := initClient()
	b, err := fc.GetUSDTBalance()
	require.NoError(t, err)
	fmt.Println(b)
}

//submitted	已提交
//partial_filled	部分成交
//partial_canceled	部分成交已撤销
//filled	完全成交
//canceled	已撤销
//pending_cancel	撤销已提交
func TestFcoinClient_GetOrders(t *testing.T) {
	fc := initClient()
	order := &Order{
		symbol: "eosusdt",
		states: "filled",
	}
	orderList, err := fc.GetOrders(order)
	require.NoError(t, err)
	fmt.Println(orderList)
}

func TestFcoinClient_CreateOrder(t *testing.T) {
	fc := initClient()
	newOrder := &NewOrder{
		Amount:    "1000",
		OrderType: "limit", //限价limit 市价 market
		Exchange:  "main",  //主板
		Side:      "buy",   //sell buy
		Symbol:    "trxusdt",
		Price:     "0.030",
	}
	res, err := fc.CreateOrder(newOrder)
	require.NoError(t, err)
	t.Log(res)
}
func TestFcoinClient_CreateOrder2(t *testing.T) {
	fc := initClient()
	newOrder := &NewOrder{
		Amount:    "1",
		OrderType: "limit", //限价limit 市价 market
		Exchange:  "main",  //主板
		Side:      "sell",  //sell buy
		Symbol:    "eosusdt",
		Price:     "3.7",
	}
	res, err := fc.CreateOrder(newOrder)
	require.NoError(t, err)
	t.Log(res)
}
func TestFcoinClient_GetTickerBySymbol(t *testing.T) {
	symbol := "ethusdt"
	fc := initClient()
	ticker, err := fc.GetLatestTickerBySymbol(symbol)
	require.NoError(t, err)
	fmt.Println(ticker)
}

func TestFcoinClient_CancelOrder(t *testing.T) {
	fc := initClient()
	id := "97bgUNZTnSfwjv3L5E3su550HmNSZafRsMf0ONY-MnvyHoSkYse3My6T2bTp-pp-eWeTeEQ9KD3NGZBrjY3qbg=="
	res, err := fc.CancelOrder(id)
	require.NoError(t, err)
	t.Log(res)
}

func TestFcoinClient_GetOrderById(t *testing.T) {
	fc := initClient()
	id := "RO__dSP9fx75hmLWczJ6-URrYRkr_1axyGIvbutt2UIGRMCtT_6JyYEzODvzjJgcadJgoyE8mENEu0CcxOd-qA"
	res, err := fc.GetOrderById(id)
	require.NoError(t, err)
	fmt.Println(res)
	require.ElementsMatch(t, res.Status, ORDER_STATES_SUCCESS)

}

func TestFcoinClient_GetCandle(t *testing.T) {
	fc := initClient()
	symbol := "eosbtc"
	resolution := "M1"
	limit := "20"
	c, err := fc.GetCandle(symbol, resolution, limit)
	require.NoError(t, err)
	fmt.Println(c)
}

func TestMax(t *testing.T) {
	d := decimal.New(3, -1)
	fmt.Println(d.String())
}

func TestGetTrades(t *testing.T) {
	fc := initClient()
	//https://api.fcoin.com/v2/market/depth/$level/$symbol
	url := fc.baseUrl + "/market/depth/L20/ethusdt"
	res, err := fc.getOpenResponse(url)
	require.NoError(t, err)
	//fmt.Println(string(res))
	d := &Depth{}
	err = json.Unmarshal(res, d)
	require.NoError(t, err)
	fmt.Println(len(d.Data.Asks))
	fmt.Println(d.Data.Asks[30])
	fmt.Println(d.Data.Asks[31])
	fmt.Println(d.Data.Asks[20])
	fmt.Println(d.Data.Asks[10])
	fmt.Println(d.Data.Asks[0])

	fmt.Println(len(d.Data.Bids))
	fmt.Println(d.Data.Bids[30])
	fmt.Println(d.Data.Bids[20])
	fmt.Println(d.Data.Bids[8])
	fmt.Println(d.Data.Bids[0])

}

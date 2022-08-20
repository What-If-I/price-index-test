package exchange

import (
	"context"
	"fmt"
	"math/rand"
	"time"
)

type Ticker string

const (
	BTCUSDTicker Ticker = "BTC_USD"
)

type TickerPrice struct {
	Ticker Ticker
	Time   time.Time
	Price  string // decimal value. example: "0", "10", "12.2", "13.2345122"
}

// Exchange represents some service for streaming data.
// You can subscribe to exchange via SubscribePriceStream method and receive price updates every pushInterval.
type Exchange struct {
	name         string
	ticker       Ticker
	pushInterval time.Duration
}

func NewExchange(name string, ticker Ticker, pushInterval time.Duration) Exchange {
	return Exchange{name: name, ticker: ticker, pushInterval: pushInterval}
}

func (e Exchange) SubscribePriceStream(ctx context.Context) (chan TickerPrice, chan error) {
	ch := make(chan TickerPrice)
	errCh := make(chan error)

	var currencyPrice float64 = 23270
	go func() {
		defer close(ch)
		defer close(errCh)
		for {
			select {
			case <-ctx.Done():
				return
			default:
				volatility := rand.Float64()
				price := TickerPrice{
					Ticker: e.ticker,
					Time:   time.Now(),
					Price:  fmt.Sprintf("%.2f", currencyPrice+volatility),
				}

				networkLatency := time.Duration(rand.Intn(3000)) * time.Millisecond
				time.Sleep(e.pushInterval + networkLatency)

				// fmt.Printf("%s: sending price %s from time %d\n", e.name, price.Price, price.Time.Unix())
				ch <- price
			}
			currencyPrice += 1 // bitcoin price is slowly going up \( ﾟヮﾟ)/
		}
	}()
	return ch, errCh
}

package main

import (
	"context"
	"fmt"
	"os/signal"
	"syscall"
	"time"

	"index-price/internal/exchange"
	"index-price/internal/priceagg"
)

const (
	amountOfExchanges  = 100
	pricesSendInterval = time.Second * 5
	tickInterval       = time.Second * 60
)

func main() {
	ctx, _ := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	// subscribe to BTC Price stream
	for price := range subscribePriceStream(ctx, exchange.BTCUSDTicker, pricesSendInterval, tickInterval) {
		// for each tick output timestamp and price
		fmt.Printf("%d, %s\n", price.Time.Unix(), price.Amount)
	}
}

func subscribePriceStream(ctx context.Context, ticker exchange.Ticker, sendInterval, tickInterval time.Duration) chan priceagg.Price {
	pricesAggregator := priceagg.New(tickInterval)

	// subscribe to exchanges
	for i := 0; i < amountOfExchanges; i++ {
		exchangeName := fmt.Sprintf("exchange_%d", i)
		priceStream := exchange.NewExchange(exchangeName, ticker, sendInterval)
		pricesAggregator.AddSource(ctx, priceStream)
	}

	ch := make(chan priceagg.Price)
	// get price from price pricesAggregator each interval
	go func() {
		defer close(ch)
		for {
			select {
			case <-ctx.Done():
				return
			default:
				time.Sleep(tickInterval)
				ch <- pricesAggregator.CalcCurrentPrice()
			}
		}
	}()
	return ch
}

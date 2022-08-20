package priceagg

import (
	"context"
	"fmt"
	"time"

	"index-price/internal/exchange"
	"index-price/internal/storage"
)

// Aggregator allows to receive prices from different exchanges and calculate current price
type Aggregator struct {
	storage        storage.PricesStorage
	timeResolution time.Duration
}

type Price struct {
	Time   time.Time
	Amount string // decimal value. example: "0", "10", "12.2", "13.2345122"
}

type priceStreamSubscriber interface {
	SubscribePriceStream(ctx context.Context) (chan exchange.TickerPrice, chan error)
}

func New(timeResolution time.Duration) Aggregator {
	pb := storage.New(timeResolution)
	return Aggregator{storage: pb}
}

func (a *Aggregator) AddSource(ctx context.Context, stream priceStreamSubscriber) {
	// subscribe to price stream
	priceStream, errStream := stream.SubscribePriceStream(ctx)

	// send data from price stream to aggregator
	go func() {
		for {
			select {
			case price := <-priceStream:
				a.storage.Add(price)
			case err := <-errStream:
				if err != nil {
					fmt.Printf("error from sream: %s", err)
				}
				// reconnection to closed stream is considered but not implemented
				return
			}
		}
	}()
	return
}

func (a *Aggregator) CalcCurrentPrice() Price {
	prices := a.storage.Get()

	t := time.Now().Truncate(a.timeResolution)
	if len(prices) == 0 {
		return Price{Time: t, Amount: "0"}
	}

	averagedAmount := calcWAM(prices)

	return Price{
		Time:   t,
		Amount: averagedAmount.StringFixed(2),
	}
}

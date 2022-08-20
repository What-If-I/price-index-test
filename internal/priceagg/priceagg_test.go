package priceagg

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"index-price/internal/exchange"
)

func TestAggregator_CalcCurrentPrice(t *testing.T) {
	aggregator := New(1 * time.Minute)

	// given recent prices for current minute
	aggregator.storage.Add(exchange.TickerPrice{Price: "199", Time: time.Now().Add(-1 * time.Second)})
	aggregator.storage.Add(exchange.TickerPrice{Price: "201", Time: time.Now().Add(10 * time.Second)})
	aggregator.storage.Add(exchange.TickerPrice{Price: "200", Time: time.Now().Add(8 * time.Second)})

	// assert price
	want := Price{Amount: "200.00", Time: time.Now().Truncate(time.Minute)}
	got := aggregator.CalcCurrentPrice()
	assert.Equal(t, want, got)
}

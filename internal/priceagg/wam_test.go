package priceagg

import (
	"testing"
	"time"

	"github.com/shopspring/decimal"

	"index-price/internal/exchange"
)

func Test_calcWMA(t *testing.T) {
	tests := []struct {
		name   string
		prices []exchange.TickerPrice
		want   decimal.Decimal
	}{
		{
			name:   "no prices",
			prices: nil,
			want:   decimal.Zero,
		},
		{
			name:   "one price",
			prices: []exchange.TickerPrice{{Price: "199", Time: unixSec(1)}},
			want:   decimal.NewFromFloat32(199),
		},
		{
			name:   "two prices with same weight",
			prices: []exchange.TickerPrice{{Price: "200", Time: unixSec(1)}, {Price: "300", Time: unixSec(1)}},
			want:   decimal.NewFromFloat32(250),
		},
		{
			name:   "two prices with different weight",
			prices: []exchange.TickerPrice{{Price: "200", Time: unixSec(1)}, {Price: "300", Time: unixSec(10)}},
			want:   decimal.NewFromFloat32(290),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := calcWAM(tt.prices).Round(2)
			if !got.Equal(tt.want) {
				t.Errorf("want %s, got %s", tt.want, got)
			}
		})
	}
}

func unixSec(s uint64) time.Time {
	return time.Unix(int64(s), 0)
}

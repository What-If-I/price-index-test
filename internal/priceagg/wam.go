package priceagg

import (
	"sort"
	"time"

	"github.com/shopspring/decimal"

	"index-price/internal/exchange"
)

// calcWAM calculates weighted average price.
func calcWAM(prices []exchange.TickerPrice) decimal.Decimal {
	if len(prices) == 0 {
		return decimal.Zero
	}

	if len(prices) == 1 {
		return decimal.RequireFromString(prices[0].Price)
	}

	prices = append(make([]exchange.TickerPrice, 0, len(prices)), prices...)
	sort.Slice(prices, func(i, j int) bool { return prices[i].Time.Before(prices[j].Time) })
	minTime := prices[0].Time

	numerator, denominator := decimal.Zero, decimal.Zero
	for _, price := range prices {
		priceWeight := calcPriceWeight(price.Time, minTime)
		numerator = numerator.Add(decimal.RequireFromString(price.Price).Mul(decimal.NewFromInt(priceWeight)))
		denominator = denominator.Add(decimal.NewFromInt(priceWeight))
	}

	return numerator.Div(denominator)
}

// calcPriceWeight calculates price weight.
// min value is needed to increase distance between prices timestamps.
// Because distance between raw unix timestamps is negligible and doesn't affect calculation of WAM.
func calcPriceWeight(priceTime time.Time, min time.Time) int64 {
	w := priceTime.Sub(min).Seconds()
	if w <= 0 {
		return 1
	}
	return int64(w)
}

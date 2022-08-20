package storage

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"index-price/internal/exchange"
)

func TestPricesStorage_Add(t *testing.T) {
	tests := []struct {
		name           string
		timeStart      time.Time
		checkAfterSec  int
		prices         []exchange.TickerPrice
		expectedPrices []exchange.TickerPrice
	}{
		{
			name:           "price added within time frame",
			timeStart:      unixSec(0),
			checkAfterSec:  10,
			prices:         []exchange.TickerPrice{{Price: "199", Time: unixSec(1)}},
			expectedPrices: []exchange.TickerPrice{{Price: "199", Time: unixSec(1)}},
		},
		{
			name:           "no prices added",
			timeStart:      unixSec(0),
			checkAfterSec:  10,
			prices:         nil,
			expectedPrices: []exchange.TickerPrice{},
		},
		{
			name:           "price added before time frame",
			timeStart:      unixSec(120),
			checkAfterSec:  1,
			prices:         []exchange.TickerPrice{{Price: "199", Time: unixSec(10)}},
			expectedPrices: []exchange.TickerPrice{},
		},
		{
			name:           "added price became outdated, but that the only price in the storage",
			timeStart:      unixSec(0),
			checkAfterSec:  71,
			prices:         []exchange.TickerPrice{{Price: "199", Time: unixSec(10)}},
			expectedPrices: []exchange.TickerPrice{{Price: "199", Time: unixSec(10)}},
		},
		{
			name:           "added price became outdated and there is relevant price",
			timeStart:      unixSec(0),
			checkAfterSec:  71,
			prices:         []exchange.TickerPrice{{Price: "199", Time: unixSec(10)}, {Price: "199", Time: unixSec(72)}},
			expectedPrices: []exchange.TickerPrice{{Price: "199", Time: unixSec(72)}},
		},
		{
			name:           "first price is expired, other two are not",
			timeStart:      unixSec(0),
			checkAfterSec:  62,
			prices:         []exchange.TickerPrice{{Price: "199", Time: unixSec(1)}, {Price: "200", Time: unixSec(2)}, {Price: "201", Time: unixSec(62)}},
			expectedPrices: []exchange.TickerPrice{{Price: "200", Time: unixSec(2)}, {Price: "201", Time: unixSec(62)}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			timeResolution := 1 * time.Minute
			p := New(timeResolution)
			p.tickLine = calcTickLine(tt.timeStart, timeResolution)
			p.timeNow = timeMock(tt.timeStart)

			for _, price := range tt.prices {
				p.Add(price)
			}

			p.timeNow = timeMock(tt.timeStart.Add(time.Second * time.Duration(tt.checkAfterSec)))

			got := p.Get()
			require.Equal(t, tt.expectedPrices, got)
		})
	}
}

func Test_calcTickLine(t *testing.T) {
	tests := []struct {
		name       string
		now        time.Time
		resolution time.Duration
		want       time.Time
	}{
		{
			name:       "time is 70, tick line is 60",
			now:        unixSec(70),
			resolution: 1 * time.Minute,
			want:       unixSec(60),
		},
		{
			name:       "time is 120, tick line is 120",
			now:        unixSec(120),
			resolution: 1 * time.Minute,
			want:       unixSec(120),
		},
		{
			name:       "time is 11, tick line is 10 when resolution is 10 seconds",
			now:        unixSec(11),
			resolution: 10 * time.Second,
			want:       unixSec(10),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := calcTickLine(tt.now, tt.resolution)
			require.Equal(t, tt.want, got)
		})
	}
}

func timeMock(t time.Time) func() time.Time {
	return func() time.Time {
		return t
	}
}

func unixSec(s uint64) time.Time {
	return time.Unix(int64(s), 0)
}

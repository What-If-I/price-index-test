package storage

import (
	"sync"
	"time"

	"index-price/internal/exchange"
)

// PricesStorage is a thread safe in-memory storage that allows to store only relevant prices.
// What price is relevant is determined by time resolution.
//
// Example:
// if time resolution is 1 minute and current time is 70 seconds then prices older than 10 seconds are no longer relevant.
// in this example prices p1 will be ignored and only p2 and p3 will be used.
//      start          tickLine
//   p1  | p2            || p3                           |
// ----------------------------------------------------->t
// 0 5  10               60      70                   120 time
//                               ^
//                               current time

type PricesStorage struct {
	mutex *sync.RWMutex

	// prices are stored in two different arrays to faster delete outdated prices
	beforeTick []exchange.TickerPrice
	afterTick  []exchange.TickerPrice

	tickLine       time.Time
	timeResolution time.Duration

	timeNow func() time.Time // allows to mock time
}

func New(timeResolution time.Duration) PricesStorage {
	tickLine := calcTickLine(time.Now(), timeResolution).Add(timeResolution)
	m := sync.RWMutex{}
	s := PricesStorage{
		beforeTick:     make([]exchange.TickerPrice, 0),
		afterTick:      make([]exchange.TickerPrice, 0),
		timeResolution: timeResolution,
		tickLine:       tickLine,
		mutex:          &m,
		timeNow:        time.Now,
	}
	return s
}

func (p *PricesStorage) Add(price exchange.TickerPrice) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	if price.Time.Before(p.tickLine.Add(-p.timeResolution)) {
		// price is too old, ignore it
		return
	}

	if price.Time.Before(p.tickLine) {
		p.beforeTick = append(p.beforeTick, price)
	} else {
		p.afterTick = append(p.afterTick, price)
	}

	currentTick := calcTickLine(p.timeNow(), p.timeResolution)
	clearTickLine := p.tickLine.Add(p.timeResolution)
	if currentTick.Equal(clearTickLine) || currentTick.After(clearTickLine) {
		// the future has come. old prices can be deleted as they are no longer needed
		p.beforeTick = p.afterTick
		p.afterTick = nil
		p.tickLine = currentTick
	}

	return
}

func (p *PricesStorage) Get() []exchange.TickerPrice {
	p.mutex.RLock()
	allPrices := append(p.beforeTick, p.afterTick...)
	p.mutex.RUnlock()

	if len(allPrices) == 0 {
		return allPrices
	}

	// everything before this line is too old to be used
	startLine := p.timeNow().Add(-p.timeResolution)

	relevantPriceStartIdx := -1
	for i, price := range allPrices {
		if price.Time.Before(startLine) {
			continue
		}
		relevantPriceStartIdx = i
		break
	}

	if relevantPriceStartIdx == -1 {
		// no relevant prices found, but its better to at least return latest price than nothing
		return []exchange.TickerPrice{allPrices[len(allPrices)-1]}
	}

	return allPrices[relevantPriceStartIdx:]
}

func calcTickLine(t time.Time, resolution time.Duration) time.Time {
	return t.Truncate(resolution)
}

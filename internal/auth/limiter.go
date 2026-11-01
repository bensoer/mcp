package auth

import (
	"math"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

// Limiter maintains independent token buckets for rate-limit identities.
type Limiter struct {
	mu       sync.Mutex
	limiters map[string]*rate.Limiter
	rps      int
}

// NewLimiter creates an in-memory token-bucket limiter with the given sustained
// requests per second and a burst of twice that value, rounded up.
func NewLimiter(rps int) *Limiter {
	if rps < 1 {
		rps = 1
	}
	return &Limiter{
		limiters: make(map[string]*rate.Limiter),
		rps:      rps,
	}
}

// Allow consumes one token for identity. When the bucket is exhausted, it
// returns the delay until the next token is available.
func (l *Limiter) Allow(identity string) (bool, time.Duration) {
	l.mu.Lock()
	limiter, ok := l.limiters[identity]
	if !ok {
		limiter = rate.NewLimiter(rate.Limit(l.rps), burstSize(l.rps))
		l.limiters[identity] = limiter
	}
	l.mu.Unlock()

	reservation := limiter.ReserveN(time.Now(), 1)
	delay := reservation.Delay()
	if delay == 0 {
		return true, 0
	}
	reservation.Cancel()
	return false, delay
}

func burstSize(rps int) int {
	return max(1, int(math.Ceil(float64(rps)*2)))
}

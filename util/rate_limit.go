package util

import (
	"time"

	"github.com/go-mojito/mojito"
)

func RateLimit(key string, limit int, decay time.Duration) (ok bool, err error) {
	count := 0
	if err = mojito.DefaultCache().GetOrDefault(key, &count, 0); err != nil {
		return
	}

	if count >= limit {
		// Reset rate limit decay to avoid spamming
		if err = mojito.DefaultCache().Set(key, count); err != nil {
			return
		}
		if err = mojito.DefaultCache().ExpireAfter(key, decay); err != nil {
			return
		}
		return
	}

	if err = mojito.DefaultCache().Set(key, count+1); err != nil {
		return
	}
	if err = mojito.DefaultCache().ExpireAfter(key, decay); err != nil {
		return
	}

	return true, nil
}

func RateLimitSingle(key string, decay time.Duration) (ok bool, err error) {
	return RateLimit(key, 1, decay)
}

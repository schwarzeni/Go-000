package Week06

import "errors"

type RateLimiter interface {
	// Allow 询问 limiter 是否可以访问，失败返回 error
	Allow() error
}

var ErrExceededLimit = errors.New("限流")

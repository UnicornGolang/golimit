package limiter

import (
	"sync"
	"time"
)

type TokenBuket struct {
	// 锁
	lock sync.Mutex
	// 令牌桶的大小
	size int
	// 当前桶中令牌数
	count int
	// token 填充速率
	rate time.Duration
	// 上次请求的时间
	lastRequestTime time.Time
}

func (tb *TokenBuket) fillToken() {
	tb.count += tb.getFillTokenCount()
}

// 填充 token
func (tb *TokenBuket) getFillTokenCount() int {
	if tb.count >= tb.size {
		return 0
	}
	if tb.lastRequestTime.IsZero() {
		return 0
	}
	duration := time.Since(tb.lastRequestTime)
	count := int(duration / tb.rate)
	if tb.size-tb.count > count {
		return count
	} else {
		return (tb.size - tb.count)
	}

}

// 是否限流
func (tb *TokenBuket) allow() bool {
	// 先填充 token
	tb.fillToken()

	// 满足 token 则放行
	if tb.count > 0 {
		tb.count--
		tb.lastRequestTime = time.Now()
		return true
	}
	return false
}

type Limiter struct {
	tokenbuket *TokenBuket
}

func NewLimiter(r time.Duration, size int) *Limiter {
	return &Limiter{
		tokenbuket: &TokenBuket{
			rate:  r,
			size:  size,
			count: size,
		},
	}
}

func (l *Limiter) Allow() bool {
	l.tokenbuket.lock.Lock()
	defer l.tokenbuket.lock.Unlock()
	// 计算补充 token 数
	// 当前 token 是否满足本次消耗
	return l.tokenbuket.allow()
}

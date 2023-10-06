package limiter

import (
	"fmt"
	"sync/atomic"
	"time"

	"golang.org/x/time/rate"
)

var totalQuery int32

func handler() {
	atomic.AddInt32(&totalQuery, 1)
	time.Sleep(50 * time.Millisecond)
}

func callHandler() {

	limiter := rate.NewLimiter(rate.Every(100*time.Millisecond), 1)
	for {
		// 1. Wait / WaitN, 阻塞等到了令牌就去执行
		// 如果只需要申请一个可以使用 limiter.wait(context.Background())
		// limiter.WaitN(context.Background(), 1)

		// 2. Reserve / ReserveN : 可以接收时间参数，返回需要等待的时间
		// cost := limiter.ReserveN(time.Now(), 1)
		// time.Sleep(cost.Delay())

		// 3. 非阻塞式等待, Allow / AllowN, 返回 boolean 表示当前是否有令牌
		if limiter.AllowN(time.Now(), 1) {
			handler()
		}
	}
}

func EmbedLimiter() {
	go callHandler()
	for i := 0; i < 5; i++ {
		time.Sleep(time.Second)
		fmt.Printf("过去 1 秒钟接口调用了%d次\n", atomic.LoadInt32(&totalQuery))
		atomic.StoreInt32(&totalQuery, 0) // 每隔 1 秒清零一次
	}
}

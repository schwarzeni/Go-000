package Week06

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestRing(t *testing.T) {
	size := 4
	r := &ring{
		data:  make([]int, size),
		headP: 0,
	}

	// [0 1 2 3]
	// p = 0
	for i := 0; i < r.size(); i++ {
		r.access(i, i, func(v *int) { *v = i })
	}
	if s := r.sum(); s != 6 {
		t.Errorf("expected 6, nut got %d", s)
	}

	// [0 0 2 3]
	// p = 2
	r.move(2)
	if s := r.sum(); s != 5 {
		t.Errorf("expected 5, nut got %d", s)
	}

	// [10 5 2 3]
	// p = 2
	r.access(2, 2, func(v *int) { *v = 10 })
	r.access(3, 3, func(v *int) { *v = 5 })
	if s := r.sum(); s != 20 {
		t.Errorf("expected 20, nut got %d", s)
	}

	// 不移动
	r.move(0)
	if s := r.sum(); s != 20 {
		t.Errorf("expected 20, nut got %d", s)
	}

	// [0 0 0 0]
	// p = 2
	r.move(r.size())
	if s := r.sum(); s != 0 {
		t.Errorf("expected 0, nut got %d", s)
	}
}

func TestSlidingWindowLimiter_Init(t *testing.T) {
	// 限速器每秒接受10次访问
	// 并发访问 12 次，失败两次
	var (
		s1        = NewSlidingWindowLimiter(10)
		errcount  int64
		wg1       sync.WaitGroup
		taskCount = 12
	)
	wg1.Add(taskCount)
	for i := 0; i < taskCount; i++ {
		go func() {
			if err := s1.Allow(); err != nil {
				atomic.AddInt64(&errcount, 1)
			}
			wg1.Done()
		}()
	}
	wg1.Wait()
	if errcount != 2 {
		t.Errorf("expect 2, but got %d", errcount)
	}
}

func TestSlidingWindowLimiter_WithInterval(t *testing.T) {
	// 限速器每秒可接受3个访问
	// 第一个 100ms，并发访问3次，都能成功访问
	// 过去 100ms 后，并发访问4次，失败一次
	var (
		s1        = NewSlidingWindowLimiter(3)
		errcount  int64
		wg1       sync.WaitGroup
		taskCount = 3
	)
	wg1.Add(taskCount)
	for i := 0; i < taskCount; i++ {
		go func() {
			if err := s1.Allow(); err != nil {
				atomic.AddInt64(&errcount, 1)
			}
			wg1.Done()
		}()
	}
	wg1.Wait()
	if errcount != 0 {
		t.Errorf("errcount should be 0, but got %d", errcount)
	}
	time.Sleep(time.Millisecond * 100)
	taskCount += 1
	wg1.Add(taskCount)
	for i := 0; i < taskCount; i++ {
		go func() {
			if err := s1.Allow(); err != nil {
				atomic.AddInt64(&errcount, 1)
			}
			wg1.Done()
		}()
	}
	wg1.Wait()
	if errcount != 1 {
		t.Errorf("errcount should be 1, but got %d", errcount)
	}
}

func TestSlidingWindowLimiter_LongInterval(t *testing.T) {
	// 限速器每秒可访问 10 次
	// 测试梅 100ms 访问一次，时长 2s ，共 20 个请求，不应该报错
	var (
		s1 = NewSlidingWindowLimiter(10)
	)
	for i := 0; i < 20; i++ {
		if err := s1.Allow(); err != nil {
			t.Errorf("unexpect err in loop %d", i)
		}
		time.Sleep(time.Millisecond * 100)
	}
}

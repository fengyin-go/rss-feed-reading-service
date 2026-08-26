package feedstate

import (
	"context"
	"sync"
)

// Waiter 维护按 key 注册的退出请求上下文。
// 关键约束：每个退出请求必须各自响应自己的取消信号，
// 不同请求之间不得共享取消状态。因此 Bind/Wait 按 key 独立
// 存储上下文，Wait 只阻塞在本次调用注册的 ctx 上。
type Waiter struct {
	mu   sync.Mutex
	wait map[string]context.Context
}

func NewWaiter() *Waiter { return &Waiter{} }

// Bind 为给定 key 注册本次退出请求的可取消 ctx。
// 同一 key 重复 Bind 会覆盖前次注册——后到的请求必须用自己的 ctx。
func (w *Waiter) Bind(key string, ctx context.Context) {
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.wait == nil {
		w.wait = make(map[string]context.Context)
	}
	w.wait[key] = ctx
}

// Wait 阻塞在 key 对应的 ctx 上，直到该 ctx 被取消。
// 注意：ctx 必须是本次退出请求自己注册的、可取消的上下文，
// 而非 context.Background()，否则取消信号永远到不了这里。
func (w *Waiter) Wait(ctx context.Context, key string, started chan<- struct{}) error {
	w.mu.Lock()
	registered := w.wait[key]
	w.mu.Unlock()

	// 取本次请求注册的 ctx；若未注册则回退到调用方传入的 ctx，
	// 保证永远等待在调用方可取消的上下文上。
	if registered != nil {
		ctx = registered
	}

	// 通知调用方：等待已开始，可以安排取消。
	started <- struct{}{}

	<-ctx.Done()

	// 清理本次注册，避免已取消的 ctx 被后续同 key 请求复用
	// 而导致新请求立即结束。
	w.mu.Lock()
	delete(w.wait, key)
	w.mu.Unlock()

	return ctx.Err()
}

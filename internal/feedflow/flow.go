package feedflow

import (
	"context"
	"rss-reader/internal/feedstate"
)

type Awaiter struct{ Waiter *feedstate.Waiter }

func (a *Awaiter) Await(ctx context.Context, key string, started chan<- struct{}) error {
	// 必须把调用方可取消的 ctx 注册为本次退出请求的等待对象，
	// 而不是 context.Background()——否则取消信号永远到不了等待方，
	// 取消后请求也不会返回。
	a.Waiter.Bind(key, ctx)
	return a.Waiter.Wait(ctx, key, started)
}

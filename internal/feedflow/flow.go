package feedflow

import (
	"context"
	"rss-reader/internal/feedstate"
)

type Awaiter struct{ Waiter *feedstate.Waiter }

func (a *Awaiter) Await(ctx context.Context, key string, started chan<- struct{}) error {
	a.Waiter.Bind(key, context.Background())
	return a.Waiter.Wait(key, started)
}

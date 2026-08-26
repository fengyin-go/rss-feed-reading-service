package feedflow

import "rss-reader/internal/feedstate"

type Observer struct{ Registry *feedstate.Registry }

func (o *Observer) Later(key string, ready chan<- struct{}, release <-chan struct{}) <-chan feedstate.Member {
	out := make(chan feedstate.Member, 1)
	go func() {
		// 在发出 ready 之前先定格快照，确保后续的在线状态更新不会波及摘要：
		// 摘要保留的是“获取列表那一刻”的在线状态，而不是 release 后的最新状态。
		snap := o.Registry.Snapshot(key)
		ready <- struct{}{}
		<-release
		if snap != nil {
			out <- *snap
		}
	}()
	return out
}

package feedflow

import (
	"context"

	"rss-reader/internal/feedstate"
)

type Manager struct{ Scheduler *feedstate.Scheduler }

// Run 启动后台调度。传入的 ctx 取消（阅读器离开超时）后，
// 调度必须随之停止，不再产生清理调用。
func (m *Manager) Run(ctx context.Context) {
	m.Scheduler.Start(ctx)
}

// Shutdown 停止后台调度并在 ctx 允许的范围内等待其退出，
// 确保关闭过程及时完成。
func (m *Manager) Shutdown(ctx context.Context) error {
	return m.Scheduler.Stop(ctx)
}

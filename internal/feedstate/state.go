package feedstate

import "errors"

var ErrLimit = errors.New("stream resource limit")
var ErrWrite = errors.New("stream write failed")

type Manager struct {
	Open      int
	Limit     int
	Committed bool
}
type Lease struct {
	manager *Manager
	closed  bool
}

func (m *Manager) Acquire() (*Lease, error) {
	if m.Open >= m.Limit {
		return nil, ErrLimit
	}
	m.Open++
	return &Lease{manager: m}, nil
}
func (l *Lease) Close() {
	if !l.closed {
		l.closed = true
		l.manager.Open--
	}
}
// Finish 提交本次导出任务的结果状态。仅在写入成功时才标记 Committed，
// 写入失败时保留未完成状态，调用方可据此区分成功与失败。
func (m *Manager) Finish(err error) error {
	if err != nil {
		// 写入失败：不提交，结果状态保持失败。
		return err
	}
	m.Committed = true
	return nil
}

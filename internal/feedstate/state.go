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
func (m *Manager) Finish(err error) error {
	if err != nil {
		m.Committed = true
		return nil
	}
	m.Committed = true
	return nil
}

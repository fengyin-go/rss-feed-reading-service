package feedstate

import "sync"

type Session struct {
	Identity string
	Labels   []string
}
type Pool struct{ pool sync.Pool }

func NewPool() *Pool { return &Pool{pool: sync.Pool{New: func() any { return &Session{} }}} }

// Get 取出复用 Session。Session 在多个发布源的阅读会话之间复用，
// 取出时必须清空上一个持有者残留的 Identity/Labels，否则会把别人的身份或标签串到本次会话。
func (p *Pool) Get() *Session {
	s := p.pool.Get().(*Session)
	s.Identity = ""
	s.Labels = nil
	return s
}

func (p *Pool) Put(value *Session) { p.pool.Put(value) }

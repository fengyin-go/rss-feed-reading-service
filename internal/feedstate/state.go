package feedstate

import "sync"

type Session struct {
	Identity string
	Labels   []string
}
type Pool struct{ pool sync.Pool }

func NewPool() *Pool               { return &Pool{pool: sync.Pool{New: func() any { return &Session{} }}} }
func (p *Pool) Get() *Session      { return p.pool.Get().(*Session) }
func (p *Pool) Put(value *Session) { p.pool.Put(value) }

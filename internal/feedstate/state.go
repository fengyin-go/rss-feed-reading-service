package feedstate

import "sync"

type Member struct {
	Name   string
	Active bool
}
type Registry struct {
	mu      sync.RWMutex
	members map[string]*Member
}

func NewRegistry() *Registry { return &Registry{members: map[string]*Member{}} }
func (r *Registry) Put(key string, value Member) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if current := r.members[key]; current != nil {
		*current = value
		return
	}
	copy := value
	r.members[key] = &copy
}
func (r *Registry) Snapshot(key string) *Member {
	r.mu.RLock()
	defer r.mu.RUnlock()
	value := r.members[key]
	if value == nil {
		return nil
	}
	return value
}

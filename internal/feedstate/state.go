package feedstate

type Policy interface {
	Allow(string) bool
	Add(string)
}
type MapPolicy struct{ rules map[string]bool }

func (p *MapPolicy) Allow(value string) bool { return p != nil && p.rules[value] }
func (p *MapPolicy) Add(value string) {
	if p == nil || p.rules == nil {
		return
	}
	p.rules[value] = true
}
func LoadPolicy(enabled bool) Policy {
	if !enabled {
		return nil
	}
	return &MapPolicy{rules: make(map[string]bool)}
}

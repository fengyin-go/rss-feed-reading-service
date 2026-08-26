package feedstate

type Policy interface {
	Allow(string) bool
	Add(string)
}
type MapPolicy struct{ rules map[string]bool }

func (p *MapPolicy) Allow(value string) bool { return p != nil && p.rules[value] }
func (p *MapPolicy) Add(value string)        { p.rules[value] = true }
func LoadPolicy(enabled bool) Policy {
	if !enabled {
		var policy *MapPolicy
		return policy
	}
	return &MapPolicy{}
}

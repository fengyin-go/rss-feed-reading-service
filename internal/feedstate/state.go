package feedstate

type Tracker struct {
	Version int
	Status  string
	effects map[string]bool
}

func NewTracker() *Tracker                          { return &Tracker{effects: map[string]bool{}} }
func (t *Tracker) Apply(version int, status string) { t.Version, t.Status = version, status }
func (t *Tracker) Effect(key string)                { t.effects[key] = true }
func (t *Tracker) EffectCount() int                 { return len(t.effects) }

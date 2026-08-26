package feedstate

type Tracker struct {
	Version int
	Status  string
	effects map[string]bool
}

func NewTracker() *Tracker { return &Tracker{effects: map[string]bool{}} }

// Apply 推进订阅源状态。版本号必须严格递增：低于或等于当前版本的写入一律忽略，
// 这样迟到到达的旧回调无法把已经推进到成功（更高版本）的状态改回恢复中。
func (t *Tracker) Apply(version int, status string) {
	if version <= t.Version {
		return
	}
	t.Version, t.Status = version, status
}

// Effect 记录某次操作产生的效果。失败与成功复用同一个操作标识：同一 key 多次
// 写入只会保留一条效果记录，恢复通知因此只发出一次。
func (t *Tracker) Effect(key string) { t.effects[key] = true }

func (t *Tracker) EffectCount() int { return len(t.effects) }

package feedflow

import "rss-reader/internal/feedstate"

type Reconnector struct{ Tracker *feedstate.Tracker }

// FirstFailed 记录首次恢复失败，并把订阅源置为恢复中（版本 1）。
func (r *Reconnector) FirstFailed(operation string) {
	r.Tracker.Effect(operation)
	r.Tracker.Apply(1, "retrying")
}

// RetrySucceeded 记录重试恢复成功，并把订阅源置为已恢复（版本 2）。
// 复用与 FirstFailed 相同的操作标识，避免新增带 "-retry" 后缀的效果记录，
// 从而保证恢复通知只发出一次。
func (r *Reconnector) RetrySucceeded(operation string) {
	r.Tracker.Effect(operation)
	r.Tracker.Apply(2, "connected")
}

// LateFirstCallback 处理迟到的首次回调。由于 Tracker.Apply 现在要求版本号
// 严格递增，当重试已经成功（版本 2）后，这条版本号为 1 的旧回调会被安全忽略，
// 不会把已恢复状态退回恢复中。
func (r *Reconnector) LateFirstCallback() {
	r.Tracker.Apply(1, "retrying")
}

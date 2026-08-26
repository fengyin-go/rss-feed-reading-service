package feedflow

import "rss-reader/internal/feedstate"

type Recorder struct{ Pool *feedstate.Pool }

// Later 注册一个延迟事件：ready 信号就绪后等待 release 放行，
// 再把 value 当前的 Identity 发出。
//
// value 是池化复用对象，甲结束会话 Put 回池后、release 放行前，
// 乙可能复用同一对象并覆盖 Identity，导致甲的延迟事件串带上乙的身份。
// 因此在记录时立即快照 Identity，而非在放行时读取指针字段，
// 延迟事件要保留甲的信息。
func (r *Recorder) Later(value *feedstate.Session, ready chan<- struct{}, release <-chan struct{}) <-chan string {
	out := make(chan string, 1)
	identity := value.Identity // 记录时快照，解耦后续对复用对象的篡改
	go func() {
		ready <- struct{}{}
		<-release
		out <- identity
	}()
	return out
}

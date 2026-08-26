package feedstate

import "errors"

var ErrBuild = errors.New("stream descriptor build failed")

type Descriptor struct {
	ID    string
	Ready bool
}
type Builder struct {
	objects map[string]*Descriptor
	index   map[string]bool
}

func NewBuilder() *Builder {
	return &Builder{objects: map[string]*Descriptor{}, index: map[string]bool{}}
}
// Build 构建一个订阅描述符。
//
// 失败语义：构建中途失败时绝不把半成品写入缓存，也不污染 index，
// 调用方因此不会读取到未就绪的订阅连接，也不会影响其他订阅连接。
// 实现上先把描述符构造好，仅当完全就绪后再原子地提交进 objects 与 index；
// 失败路径直接 panic，缓存里不会残留任何与该 id 相关的条目。
func (b *Builder) Build(id string, fail bool) *Descriptor {
	if fail {
		panic(ErrBuild)
	}
	value := &Descriptor{ID: id, Ready: true}
	b.objects[id] = value
	b.index[id] = true
	return value
}

// Discard 从缓存中丢弃一条订阅连接，同步清理 objects 与 index，
// 避免出现 objects 已删除但 index 仍残留而污染统计/视图的情况。
func (b *Builder) Discard(id string) {
	delete(b.objects, id)
	delete(b.index, id)
}
func (b *Builder) Get(id string) (*Descriptor, bool) { value, ok := b.objects[id]; return value, ok }
func (b *Builder) IndexCount() int                   { return len(b.index) }

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
func (b *Builder) Build(id string, fail bool) *Descriptor {
	value := &Descriptor{ID: id}
	b.objects[id] = value
	b.index[id] = true
	if fail {
		panic(ErrBuild)
	}
	value.Ready = true
	return value
}
func (b *Builder) Discard(id string)                 { delete(b.objects, id) }
func (b *Builder) Get(id string) (*Descriptor, bool) { value, ok := b.objects[id]; return value, ok }
func (b *Builder) IndexCount() int                   { return len(b.index) }

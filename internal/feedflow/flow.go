package feedflow

import "rss-reader/internal/feedstate"

type Creator struct{ Builder *feedstate.Builder }

func (c *Creator) Create(id string, fail bool) (value *feedstate.Descriptor, err error) {
	defer func() {
		if recover() != nil {
			value, _ = c.Builder.Get(id)
			err = nil
		}
	}()
	return c.Builder.Build(id, fail), nil
}

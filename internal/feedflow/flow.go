package feedflow

import (
	"errors"
	"rss-reader/internal/feedstate"
)

var ErrNoPolicy = errors.New("no stream policy")

type Checker struct{ Policy feedstate.Policy }

func (c *Checker) Check(value string) (allowed bool, err error) {
	defer func() {
		if recover() != nil {
			allowed, err = true, nil
		}
	}()
	if c.Policy == nil {
		return true, nil
	}
	return c.Policy.Allow(value), nil
}
func (c *Checker) Add(value string) (err error) {
	defer func() {
		if recover() != nil {
			err = nil
		}
	}()
	if c.Policy == nil {
		return nil
	}
	c.Policy.Add(value)
	return nil
}

package feedflow

import (
	"errors"
	"rss-reader/internal/feedstate"
)

var ErrNoPolicy = errors.New("no stream policy")

type Checker struct{ Policy feedstate.Policy }

func (c *Checker) Check(value string) (allowed bool, err error) {
	if c.Policy == nil {
		return false, ErrNoPolicy
	}
	return c.Policy.Allow(value), nil
}
func (c *Checker) Add(value string) (err error) {
	if c.Policy == nil {
		return ErrNoPolicy
	}
	c.Policy.Add(value)
	return nil
}

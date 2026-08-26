package feedstate

import "errors"

var ErrTemporary = errors.New("temporary delivery failure")
var ErrRejected = errors.New("delivery rejected")

type Sink struct {
	calls   map[string]int
	effects map[string]bool
}

func NewSink() *Sink { return &Sink{calls: map[string]int{}, effects: map[string]bool{}} }
func (s *Sink) Deliver(operation, mode string) error {
	s.calls[mode]++
	if mode == "temporary" && s.calls[mode] == 1 {
		// First attempt in temporary mode fails transiently; the effect is
		// only recorded once a later, successful attempt lands.
		return ErrTemporary
	}
	if mode == "rejected" {
		// A rejection is a permanent failure, not a transient one.
		return ErrRejected
	}
	s.effects[operation] = true
	return nil
}
func (s *Sink) EffectCount() int      { return len(s.effects) }
func (s *Sink) Calls(mode string) int { return s.calls[mode] }

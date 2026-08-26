package feedflow

import "rss-reader/internal/feedstate"

type Archiver struct{ Manager *feedstate.Manager }

func (a *Archiver) Run(count, failAt int) error {
	for i := 0; i < count; i++ {
		lease, err := a.Manager.Acquire()
		if err != nil {
			return err
		}
		defer lease.Close()
		if i == failAt {
			if err := a.Manager.Finish(feedstate.ErrWrite); err != nil {
				return err
			}
		}
	}
	return a.Manager.Finish(nil)
}

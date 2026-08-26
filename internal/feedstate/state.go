package feedstate

import "errors"

var ErrRejected = errors.New("stream item rejected")

type Producer struct{}

func (Producer) Start(values []string, reject int) (<-chan string, <-chan error) {
	out, errs := make(chan string, len(values)), make(chan error, 1)
	go func() {
		defer close(out)
		defer close(errs)
		for i, value := range values {
			if i == reject {
				errs <- ErrRejected
				return
			}
			out <- value
		}
	}()
	return out, errs
}

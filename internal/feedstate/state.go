package feedstate

import "errors"

var ErrRejected = errors.New("stream item rejected")

type Producer struct{}

func (Producer) Start(values []string, reject int) (<-chan string, <-chan error) {
	out, errs := make(chan string, len(values)), make(chan error, 1)
	go func() {

		for i, value := range values {
			if i == reject {
				return
			}
			out <- value
		}
	}()
	return out, errs
}

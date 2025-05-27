package message

import (
	"sync"
	"time"

	"messager/domain/message"
)

type Job interface {
	Start()
	Stop()
}

type job struct {
	service  message.Service
	ticker   *time.Ticker
	duration time.Duration
	stop     chan struct{}
	start    bool
	wg       *sync.WaitGroup
	onError  func(err error)
}

func New(service message.Service, interval time.Duration, onError func(err error)) Job {
	j := job{
		service:  service,
		ticker:   nil,
		duration: interval,
		stop:     make(chan struct{}),
		start:    false,
		wg:       new(sync.WaitGroup),
		onError:  onError,
	}

	if j.onError == nil {
		j.onError = func(err error) {}
	}

	return &j
}

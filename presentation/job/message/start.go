package message

import (
	"context"
	"time"
)

func (j *job) Start() {
	if j.start {
		return
	}

	j.ticker = time.NewTicker(j.duration)
	j.stop = make(chan struct{})
	j.start = true

	go func() {
		for {
			select {
			case <-j.ticker.C:
				j.wg.Add(1)
				defer j.wg.Done()

				if err := j.service.Process(context.Background()); err != nil {
					j.onError(err)
				}
			case <-j.stop:
				j.ticker.Stop()

				return
			}
		}
	}()
}

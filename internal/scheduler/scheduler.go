package scheduler

import (
	"errandboi/internal/publisher"
	"errors"
	"fmt"
	"time"

	"github.com/gammazero/workerpool"
)

type scheduler struct {
	Publisher *publisher.Publisher
	Stop      chan struct{}
}

func NewScheduler(pb *publisher.Publisher) (*scheduler, error) {
	sch := &scheduler{Stop: make(chan struct{})}
	if pb != nil {
		sch.Publisher = pb
		return sch, nil
	}
	return nil, errors.New("Publisher cannot be null")
}

func (sch *scheduler) WorkInIntervals(d time.Duration) {
	ticker := time.NewTicker(d)
	go func() {
		for {
			select {
			case <-ticker.C:
				sch.Publisher.GetEvents()
				sch.Publisher.Work()
			case <-sch.Stop:
				err := sch.Publisher.Cancel()

				if err != nil {
					fmt.Println("publisher was not cancelled")
				}
				sch.Publisher.Wp = workerpool.New(sch.Publisher.WorkerSize)
				ticker.Stop()
				return
			}
		}
	}()
}

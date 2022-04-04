package scheduler

import (
	"fmt"
	"time"

	"errandboi/internal/publisher"

	"github.com/gammazero/workerpool"
	"go.uber.org/zap"
)

type scheduler struct {
	Publisher *publisher.Publisher
	Stop      chan struct{}
	Logger    *zap.Logger
}

func NewScheduler(pb *publisher.Publisher, logger *zap.Logger) (*scheduler, error) {
	sch := &scheduler{Stop: make(chan struct{}), Logger: logger}
	if pb != nil {
		sch.Publisher = pb

		return sch, nil
	}

	pubisherNilErr := fmt.Errorf("publisher cannot be nil")

	return nil, pubisherNilErr
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
					sch.Logger.Error("publisher was not cancelled", zap.Error(err))
				}

				sch.Publisher.Wp = workerpool.New(sch.Publisher.WorkerSize)

				ticker.Stop()

				return
			}
		}
	}()
}

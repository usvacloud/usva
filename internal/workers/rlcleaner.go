package workers

import (
	"context"
	"time"

	"github.com/usvacloud/usva/pkg/ratelimit"
)

type RatelimitCleaner struct {
	Ratelimiter ratelimit.Ratelimiter
	Interval    time.Duration
	Running     bool
}

func NewRatelimitCleaner(rl ratelimit.Ratelimiter, it time.Duration) *RatelimitCleaner {
	return &RatelimitCleaner{
		Ratelimiter: rl,
		Interval:    it,
	}
}

func (r *RatelimitCleaner) worker(ctx context.Context, ch chan struct{}) {
	r.Running = true
	ticker := time.NewTicker(r.Interval)
	for r.Running {
		r.Ratelimiter.Clean()
		<-ticker.C
	}
	ch <- struct{}{}
}

func (r *RatelimitCleaner) Run(ctx context.Context) chan struct{} {
	ch := make(chan struct{}, 1)
	go r.worker(ctx, ch)
	return ch
}

package updater

import (
	"context"
	"time"
)

type UpdateKey struct{}

type Updater struct {
	*time.Ticker

	reset chan struct{}

	f func()
	d time.Duration
}

func New(d time.Duration, f func()) *Updater {
	return &Updater{
		Ticker: time.NewTicker(d),
		f: f,
		reset: make(chan struct{}, 1),
		d: d,
	}
}

func (u *Updater) Update() {
	for {
		select{
		case <-u.C:
			u.f()
		case <-u.reset:
			<-u.reset
			continue
		}
	}
}

func (u *Updater) Reset() {
	u.Ticker.Stop()
	u.reset <- struct{}{}
	u.Ticker.Reset(u.d)
	u.reset <- struct{}{}
}

func WithUpdateContext(ctx context.Context) context.Context {
	return context.WithValue(
		ctx,
		UpdateKey{},
		struct{}{},
	)
}

func IsUpdateContext(ctx context.Context) bool {
	val := ctx.Value(UpdateKey{})

	_, ok := val.(struct{})

	return ok
}
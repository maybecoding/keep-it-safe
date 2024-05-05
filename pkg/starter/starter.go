package starter

import (
	"context"

	"golang.org/x/sync/errgroup"
)

type Fn func(ctx context.Context) error

type Starter struct {
	g          *errgroup.Group
	ctx        context.Context
	onRun      []Fn
	onShutdown []Fn
}

func New(ctx context.Context) *Starter {
	onRun := make([]Fn, 0, 10)
	onShutdown := make([]Fn, 0, 10)
	g, ctx := errgroup.WithContext(ctx)
	return &Starter{
		g:          g,
		ctx:        ctx,
		onRun:      onRun,
		onShutdown: onShutdown,
	}
}

func (s *Starter) OnRun(fns ...Fn) *Starter {
	s.onRun = append(s.onRun, fns...)
	return s
}

func (s *Starter) OnShutdown(fns ...Fn) *Starter {
	s.onShutdown = append(s.onShutdown, fns...)
	return s
}

func (s *Starter) Run() error {
	for _, fn := range s.onRun {
		fn := fn
		s.g.Go(func() error {
			return fn(s.ctx)
		})
	}
	for _, fn := range s.onShutdown {
		fn := fn
		s.g.Go(func() error {
			<-s.ctx.Done()
			return fn(s.ctx)
		})
	}
	return s.g.Wait()
}

func (s *Starter) Cancel(err error) {
	s.g.Go(func() error {
		return err
	})
}

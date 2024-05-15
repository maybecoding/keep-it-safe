// Package starter used for starting and shuting down application components.
package starter

import (
	"context"
	"fmt"

	"golang.org/x/sync/errgroup"
)

// Fn type for functions can be injected into starter.
type Fn func(ctx context.Context) error

// Starter struct for starting components.
type Starter struct {
	g          *errgroup.Group
	ctx        context.Context
	onRun      []Fn
	onShutdown []Fn
}

// New creates a new starter.
func New(ctx context.Context) *Starter {
	const defaultStartSliceLen = 10
	onRun := make([]Fn, 0, defaultStartSliceLen)
	onShutdown := make([]Fn, 0, defaultStartSliceLen)
	g, ctx := errgroup.WithContext(ctx)
	return &Starter{
		g:          g,
		ctx:        ctx,
		onRun:      onRun,
		onShutdown: onShutdown,
	}
}

// OnRun mark function for run on run.
func (s *Starter) OnRun(fns ...Fn) *Starter {
	s.onRun = append(s.onRun, fns...)
	return s
}

// OnShutdown mark function for run on termination.
func (s *Starter) OnShutdown(fns ...Fn) *Starter {
	s.onShutdown = append(s.onShutdown, fns...)
	return s
}

// Run runs all functions.
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
	err := s.g.Wait()
	if err != nil {
		return fmt.Errorf("starter - Run s.g.Wait: %w", err)
	}
	return nil
}

// Cancel canceles execution.
func (s *Starter) Cancel(err error) {
	s.g.Go(func() error {
		return err
	})
}

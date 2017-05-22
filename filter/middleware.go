package filter

import (
	"github.com/prizem-io/gateway/backend"
	"github.com/prizem-io/gateway/context"
)

type filterMiddleware struct {
	currentFilter  int
	filters        []Execution
	backendHandler backend.Handler
	nextCalled     bool
	stopped        bool
}

func (m *filterMiddleware) Execute(ctx context.Context) error {
	for !m.stopped {
		next := m.currentFilter
		if next < len(m.filters) {
			m.currentFilter++
			execution := &m.filters[next]
			m.nextCalled = false
			err := execution.Filter.Evaluate(ctx, execution.Configuration)
			if err != nil {
				m.stopped = true
				return err
			}
		} else {
			m.stopped = true
			if m.backendHandler != nil {
				err := m.backendHandler.Handle(ctx)
				if err != nil {
					m.stopped = true
					return err
				}
			}
		}
		m.nextCalled = true
	}

	return nil
}

func (m *filterMiddleware) Next(ctx context.Context) error {
	if m.nextCalled {
		// Prevent from calling more than once from a single filter
		return nil
	}

	return m.Execute(ctx)
}

func (m *filterMiddleware) Stop() {
	m.stopped = true
}

func (m *filterMiddleware) IsStopped() bool {
	return m.stopped
}

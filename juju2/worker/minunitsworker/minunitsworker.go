// Copyright 2012, 2013 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package minunitsworker

import (
	"github.com/juju/loggo"
	worker "gopkg.in/juju/worker.v1"

	"github.com/juju/1.25-upgrade/juju2/state"
	"github.com/juju/1.25-upgrade/juju2/watcher/legacy"
)

var logger = loggo.GetLogger("juju.worker.minunitsworker")

// MinUnitsWorker ensures the minimum number of units for applications is respected.
type MinUnitsWorker struct {
	st *state.State
}

// NewMinUnitsWorker returns a Worker that runs service.EnsureMinUnits()
// if the number of alive units belonging to a service decreases, or if the
// minimum required number of units for a service is increased.
func NewMinUnitsWorker(st *state.State) worker.Worker {
	mu := &MinUnitsWorker{st: st}
	return legacy.NewStringsWorker(mu)
}

func (mu *MinUnitsWorker) SetUp() (state.StringsWatcher, error) {
	return mu.st.WatchMinUnits(), nil
}

func (mu *MinUnitsWorker) handleOneService(serviceName string) error {
	service, err := mu.st.Application(serviceName)
	if err != nil {
		return err
	}
	return service.EnsureMinUnits()
}

func (mu *MinUnitsWorker) Handle(serviceNames []string) error {
	for _, name := range serviceNames {
		logger.Infof("processing service %q", name)
		if err := mu.handleOneService(name); err != nil {
			logger.Errorf("failed to process service %q: %v", name, err)
			return err
		}
	}
	return nil
}

func (mu *MinUnitsWorker) TearDown() error {
	// Nothing to do here.
	return nil
}

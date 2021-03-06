// Copyright 2015 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package txnpruner

import (
	"time"

	"github.com/juju/errors"
	"github.com/juju/utils/clock"
	"gopkg.in/juju/worker.v1"

	jworker "github.com/juju/1.25-upgrade/juju2/worker"
)

// TransactionPruner defines the interface for types capable of
// pruning transactions.
type TransactionPruner interface {
	MaybePruneTransactions() error
}

// New returns a worker which periodically prunes the data for
// completed transactions.
func New(tp TransactionPruner, interval time.Duration, clock clock.Clock) worker.Worker {
	return jworker.NewSimpleWorker(func(stopCh <-chan struct{}) error {
		for {
			select {
			case <-clock.After(interval):
				err := tp.MaybePruneTransactions()
				if err != nil {
					return errors.Annotate(err, "pruning failed, txnpruner stopping")
				}
			case <-stopCh:
				return nil
			}
		}
	})
}

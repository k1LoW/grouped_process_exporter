package grouper

import (
	"github.com/k1LoW/grouped_process_exporter/grouped_proc"
	"github.com/k1LoW/grouped_process_exporter/metric"
	"golang.org/x/sync/semaphore"
)

type Grouper interface {
	Name() string
	SetNormalizeRegexp(nReStr string) error
	SetExcludeRegexp(eReStr string) error
	Collect(gprocs *grouped_proc.GroupedProcs, enabled map[metric.MetricKey]bool, sem *semaphore.Weighted) error
}

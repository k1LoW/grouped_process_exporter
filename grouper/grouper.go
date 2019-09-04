package grouper

import (
	"github.com/k1LoW/grouped_process_exporter/grouped_proc"
	"github.com/k1LoW/grouped_process_exporter/metric"
)

type Grouper interface {
	Name() string
	SetNormalizeRegexp(nReStr string) error
	Collect(gprocs *grouped_proc.GroupedProcs, enabled map[metric.MetricKey]bool) error
}

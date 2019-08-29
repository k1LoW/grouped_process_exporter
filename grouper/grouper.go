package grouper

import "github.com/k1LoW/grouped_process_exporter/grouped_proc"

type Grouper interface {
	Name() string
	Collect(gpMap map[string]*grouped_proc.GroupedProc, enabled map[grouped_proc.MetricKey]bool) error
}

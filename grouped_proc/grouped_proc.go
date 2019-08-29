package grouped_proc

import (
	"sync"

	"github.com/prometheus/procfs"
)

type GroupedProc struct {
	sync.Mutex
	Metrics map[MetricKey]Metric
	Enabled map[MetricKey]bool
}

func NewGroupedProc(enabled map[MetricKey]bool) *GroupedProc {
	return &GroupedProc{
		Enabled: enabled,
		Metrics: AvairableMetrics(),
	}
}

func DefaultEnabledMetrics() map[MetricKey]bool {
	enabled := make(map[MetricKey]bool)
	for _, k := range MetricKeys {
		enabled[k] = false
	}
	enabled[ProcCount] = true
	return enabled
}

func (g *GroupedProc) AppendPid(pid int) error {
	proc, err := procfs.NewProc(pid)
	if err != nil {
		return err
	}

	for _, k := range MetricKeys {
		if g.Enabled[k] {
			g.Lock()
			err := g.Metrics[k].CollectFromProc(proc)
			g.Unlock()
			if err != nil {
				return err
			}
		}
	}

	return nil
}

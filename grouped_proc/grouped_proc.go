package grouped_proc

import (
	"sync"

	"github.com/k1LoW/grouped_process_exporter/metric"
	"github.com/prometheus/procfs"
)

type GroupedProc struct {
	sync.Mutex
	Metrics        map[metric.MetricKey]metric.Metric
	Enabled        map[metric.MetricKey]bool
	Exists         bool
	ProcMountPoint string
}

func NewGroupedProc(enabled map[metric.MetricKey]bool) *GroupedProc {
	return &GroupedProc{
		Enabled:        enabled,
		Metrics:        metric.AvairableMetrics(),
		Exists:         true,
		ProcMountPoint: procfs.DefaultMountPoint,
	}
}

func DefaultEnabledMetrics() map[metric.MetricKey]bool {
	enabled := make(map[metric.MetricKey]bool)
	for _, k := range metric.MetricKeys {
		enabled[k] = false
	}
	enabled[metric.ProcProcs] = true
	return enabled
}

func (g *GroupedProc) AppendAndCollectFromProc(pid int) error {
	fs, err := procfs.NewFS(g.ProcMountPoint)
	if err != nil {
		return err
	}
	proc, err := fs.Proc(pid)
	if err != nil {
		return err
	}

	for _, k := range metric.MetricKeys {
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

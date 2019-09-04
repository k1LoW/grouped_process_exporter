package grouped_proc

import (
	"os"
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
	procMountPoint := os.Getenv("GROUPED_PROCESS_PROC_MOUNT_POINT")
	if procMountPoint == "" {
		procMountPoint = procfs.DefaultMountPoint
	}
	return &GroupedProc{
		Enabled:        enabled,
		Metrics:        metric.AvairableMetrics(),
		Exists:         true,
		ProcMountPoint: procMountPoint,
	}
}

func (g *GroupedProc) AppendProcAndCollect(pid int) error {
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

type GroupedProcs struct {
	sm sync.Map
}

func (m *GroupedProcs) Load(group string) (*GroupedProc, bool) {
	gproc, ok := m.sm.Load(group)
	if !ok {
		return nil, false
	}
	return gproc.(*GroupedProc), true
}

func (m *GroupedProcs) Store(group string, gproc *GroupedProc) {
	m.sm.Store(group, gproc)
}

func (m *GroupedProcs) Delete(group string) {
	m.sm.Delete(group)
}

func (m *GroupedProcs) Range(f func(group string, gproc *GroupedProc) bool) {
	m.sm.Range(func(key, value interface{}) bool {
		return f(key.(string), value.(*GroupedProc))
	})
}

func (m *GroupedProcs) Length() int {
	l := 0
	m.Range(func(group string, gproc *GroupedProc) bool {
		l = l + 1
		return true
	})
	return l
}

func NewGroupedProcs() *GroupedProcs {
	return &GroupedProcs{}
}

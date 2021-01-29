package metric

import (
	"sync"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/procfs"
)

// ProcProcsMetric is metric
type ProcProcsMetric struct {
	sync.Mutex
	metrics map[int]struct{}
}

func (m *ProcProcsMetric) Describe() map[string]*prometheus.Desc {
	descs := map[string]*prometheus.Desc{
		"grouped_process_num_procs": prometheus.NewDesc(
			"grouped_process_num_procs",
			"Number of processes in the group",
			[]string{"grouper", "group"}, nil,
		),
	}
	return descs
}

func (m *ProcProcsMetric) String() string {
	return "proc_count"
}

func (m *ProcProcsMetric) CollectFromProc(proc procfs.Proc) error {
	m.Lock()
	m.metrics[proc.PID] = struct{}{}
	m.Unlock()
	return nil
}

func (m *ProcProcsMetric) PushCollected(ch chan<- prometheus.Metric, descs map[string]*prometheus.Desc, grouper string, group string) error {
	m.Lock()
	if d, ok := descs["grouped_process_num_procs"]; ok {
		ch <- prometheus.MustNewConstMetric(d, prometheus.GaugeValue, float64(len(m.metrics)), grouper, group)
	}
	m.metrics = make(map[int]struct{}) // clear
	m.Unlock()
	return nil
}

func (m *ProcProcsMetric) RequiredWeight() int64 {
	return 0
}

func NewProcProcsMetric() *ProcProcsMetric {
	return &ProcProcsMetric{
		metrics: make(map[int]struct{}),
	}
}

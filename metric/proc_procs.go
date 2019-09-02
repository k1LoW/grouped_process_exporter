package metric

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/procfs"
)

// ProcProcsMetric is metric
type ProcProcsMetric struct {
	metrics map[int]struct{}
}

func (m *ProcProcsMetric) Describe() map[string]*prometheus.Desc {
	descs := map[string]*prometheus.Desc{
		"grouped_process_procs": prometheus.NewDesc(
			"grouped_process_procs",
			"Amount of grouped procs",
			[]string{"grouper", "group"}, nil,
		),
	}
	return descs
}

func (m *ProcProcsMetric) String() string {
	return "proc_count"
}

func (m *ProcProcsMetric) CollectFromProc(proc procfs.Proc) error {
	m.metrics[proc.PID] = struct{}{}
	return nil
}

func (m *ProcProcsMetric) PushCollected(ch chan<- prometheus.Metric, descs map[string]*prometheus.Desc, grouper string, group string) error {
	ch <- prometheus.MustNewConstMetric(descs["grouped_process_procs"], prometheus.GaugeValue, float64(len(m.metrics)), grouper, group)
	m.metrics = make(map[int]struct{}) // clear
	return nil
}

func NewProcProcsMetric() *ProcProcsMetric {
	return &ProcProcsMetric{
		metrics: make(map[int]struct{}),
	}
}

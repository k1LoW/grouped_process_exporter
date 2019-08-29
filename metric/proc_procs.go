package metric

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/procfs"
)

// ProcProcsMetric is metric
type ProcProcsMetric struct {
	Pids []int
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
	m.Pids = append(m.Pids, proc.PID)
	return nil
}

func (m *ProcProcsMetric) SetCollectedMetric(ch chan<- prometheus.Metric, descs map[string]*prometheus.Desc, grouper string, group string) error {
	ch <- prometheus.MustNewConstMetric(descs["grouped_process_processes"], prometheus.GaugeValue, float64(len(m.Pids)), grouper, group)
	return nil
}

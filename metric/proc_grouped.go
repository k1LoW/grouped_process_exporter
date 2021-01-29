package metric

import (
	"sync"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/procfs"
)

// ProcGroupedMetric is metric
type ProcGroupedMetric struct {
	sync.Mutex
	count float64
}

func (m *ProcGroupedMetric) Describe() map[string]*prometheus.Desc {
	descs := map[string]*prometheus.Desc{
		"grouped_process_num_grouped": prometheus.NewDesc(
			"grouped_process_num_grouped",
			"Number of grouped",
			[]string{"grouper", "group"}, nil,
		),
	}
	return descs
}

func (m *ProcGroupedMetric) String() string {
	return "grouped_count"
}

func (m *ProcGroupedMetric) Collect(k string) error {
	m.Lock()
	m.count = m.count + 1
	m.Unlock()
	return nil
}

func (m *ProcGroupedMetric) CollectFromProc(proc procfs.Proc) error {
	return nil
}

func (m *ProcGroupedMetric) PushCollected(ch chan<- prometheus.Metric, descs map[string]*prometheus.Desc, grouper string, group string) error {
	m.Lock()
	if d, ok := descs["grouped_process_num_grouped"]; ok {
		ch <- prometheus.MustNewConstMetric(d, prometheus.GaugeValue, m.count, grouper, group)
	}
	m.count = 0 // clear
	m.Unlock()
	return nil
}

func (m *ProcGroupedMetric) RequiredWeight() int64 {
	return 0
}

func NewProcGroupedMetric() *ProcGroupedMetric {
	return &ProcGroupedMetric{
		count: 0,
	}
}

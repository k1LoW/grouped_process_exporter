package metric

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/procfs"
)

type Metric interface {
	Describe() map[string]*prometheus.Desc
	String() string
	CollectFromProc(proc procfs.Proc) error
	PushCollected(ch chan<- prometheus.Metric, descs map[string]*prometheus.Desc, grouper string, group string) error
}

type MetricKey string

var (
	ProcProcs MetricKey = "proc_procs"
	ProcStat  MetricKey = "proc_stat"
	ProcIO    MetricKey = "proc_io"
)

var MetricKeys = []MetricKey{
	ProcProcs,
	ProcStat,
	ProcIO,
}

func AvairableMetrics() map[MetricKey]Metric {
	metrics := map[MetricKey]Metric{}

	// procs
	metrics[ProcProcs] = NewProcProcsMetric()

	// stat
	metrics[ProcStat] = NewProcStatMetric()

	// io
	metrics[ProcIO] = NewProcIOMetric()

	return metrics
}

func DefaultEnabledMetrics() map[MetricKey]bool {
	enabled := make(map[MetricKey]bool)
	for _, k := range MetricKeys {
		enabled[k] = false
	}
	enabled[ProcProcs] = true
	return enabled
}

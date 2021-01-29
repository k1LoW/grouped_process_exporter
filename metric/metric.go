package metric

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/procfs"
)

type Metric interface {
	Describe() map[string]*prometheus.Desc
	String() string
	Collect(k string) error
	CollectFromProc(proc procfs.Proc) error
	PushCollected(ch chan<- prometheus.Metric, descs map[string]*prometheus.Desc, grouper string, group string) error
	RequiredWeight() int64
}

type MetricKey string

var (
	ProcProcs   MetricKey = "proc_procs"
	ProcGrouped MetricKey = "proc_grouped"
	ProcStat    MetricKey = "proc_stat"
	ProcIO      MetricKey = "proc_io"
	ProcStatus  MetricKey = "proc_status"
)

var MetricKeys = []MetricKey{
	ProcProcs,
	ProcGrouped,
	ProcStat,
	ProcIO,
	ProcStatus,
}

func AvairableMetrics() map[MetricKey]Metric {
	metrics := map[MetricKey]Metric{}

	// procs
	metrics[ProcProcs] = NewProcProcsMetric()

	// grouped
	metrics[ProcGrouped] = NewProcGroupedMetric()

	// stat
	metrics[ProcStat] = NewProcStatMetric()

	// io
	metrics[ProcIO] = NewProcIOMetric()

	// status
	metrics[ProcStatus] = NewProcStatusMetric()

	return metrics
}

func DefaultEnabledMetrics() map[MetricKey]bool {
	enabled := make(map[MetricKey]bool)
	for _, k := range MetricKeys {
		enabled[k] = false
	}
	enabled[ProcProcs] = true
	enabled[ProcGrouped] = true
	return enabled
}

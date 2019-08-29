package metric

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/procfs"
)

type Metric interface {
	Describe() map[string]*prometheus.Desc
	String() string
	CollectFromProc(proc procfs.Proc) error
	SetCollectedMetric(ch chan<- prometheus.Metric, descs map[string]*prometheus.Desc, grouper string, group string) error
}

type MetricKey string

var (
	ProcProcs MetricKey = "proc_procs"
	ProcIO    MetricKey = "proc_io"
)

var MetricKeys = []MetricKey{
	ProcProcs,
	ProcIO,
}

func AvairableMetrics() map[MetricKey]Metric {
	metrics := map[MetricKey]Metric{}

	// procs
	metrics[ProcProcs] = &ProcProcsMetric{}

	// io
	metrics[ProcIO] = &ProcIOMetric{}

	return metrics
}

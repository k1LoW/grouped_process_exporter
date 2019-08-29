package grouped_proc

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
	ProcCount MetricKey = "proc_count"
	ProcIO    MetricKey = "proc_io"
)

var MetricKeys = []MetricKey{
	ProcCount,
	ProcIO,
}

func AvairableMetrics() map[MetricKey]Metric {
	metrics := map[MetricKey]Metric{}

	// count
	metrics[ProcCount] = &ProcCountMetric{}

	// io
	metrics[ProcIO] = &ProcIOMetric{}

	return metrics
}

// ProcCountMetric is metric
type ProcCountMetric struct {
	Pids []int
}

func (m *ProcCountMetric) Describe() map[string]*prometheus.Desc {
	descs := map[string]*prometheus.Desc{
		"grouped_process_processes": prometheus.NewDesc(
			"grouped_process_processes",
			"TODO",
			[]string{"grouper", "group"}, nil,
		),
	}
	return descs
}

func (m *ProcCountMetric) String() string {
	return "proc_count"
}

func (m *ProcCountMetric) CollectFromProc(proc procfs.Proc) error {
	m.Pids = append(m.Pids, proc.PID)
	return nil
}

func (m *ProcCountMetric) SetCollectedMetric(ch chan<- prometheus.Metric, descs map[string]*prometheus.Desc, grouper string, group string) error {
	ch <- prometheus.MustNewConstMetric(descs["grouped_process_processes"], prometheus.GaugeValue, float64(len(m.Pids)), grouper, group)
	return nil
}

// ProcIOMetric is metric
type ProcIOMetric struct {
	procfs.ProcIO
}

func (m *ProcIOMetric) Describe() map[string]*prometheus.Desc {
	descs := map[string]*prometheus.Desc{
		"grouped_process_io_r_char": prometheus.NewDesc(
			"grouped_process_io_r_char",
			"TODO",
			[]string{"grouper", "group"}, nil,
		),
		"grouped_process_io_w_char": prometheus.NewDesc(
			"grouped_process_io_w_char",
			"TODO",
			[]string{"grouper", "group"}, nil,
		),
		"grouped_process_io_sysc_r": prometheus.NewDesc(
			"grouped_process_io_sysc_r",
			"TODO",
			[]string{"grouper", "group"}, nil,
		),
		"grouped_process_io_sysc_w": prometheus.NewDesc(
			"grouped_process_io_sysc_w",
			"TODO",
			[]string{"grouper", "group"}, nil,
		),
		"grouped_process_io_read_bytes": prometheus.NewDesc(
			"grouped_process_io_read_bytes",
			"TODO",
			[]string{"grouper", "group"}, nil,
		),
		"grouped_process_io_write_bytes": prometheus.NewDesc(
			"grouped_process_io_write_bytes",
			"TODO",
			[]string{"grouper", "group"}, nil,
		),
		"grouped_process_io_cancelled_write_bytes": prometheus.NewDesc(
			"grouped_process_io_cancelled_write_bytes",
			"TODO",
			[]string{"grouper", "group"}, nil,
		),
	}
	return descs
}

func (m *ProcIOMetric) String() string {
	return "proc_io"
}

func (m *ProcIOMetric) CollectFromProc(proc procfs.Proc) error {
	pio, err := proc.IO()
	if err != nil {
		return err
	}
	m.RChar = m.RChar + pio.RChar
	m.WChar = m.WChar + pio.WChar
	m.SyscR = m.SyscR + pio.SyscR
	m.SyscW = m.SyscW + pio.SyscW
	m.ReadBytes = m.ReadBytes + pio.ReadBytes
	m.WriteBytes = m.WriteBytes + pio.WriteBytes
	m.CancelledWriteBytes = m.CancelledWriteBytes + pio.CancelledWriteBytes
	return nil
}

func (m *ProcIOMetric) SetCollectedMetric(ch chan<- prometheus.Metric, descs map[string]*prometheus.Desc, grouper string, group string) error {
	ch <- prometheus.MustNewConstMetric(descs["grouped_process_io_r_char"], prometheus.GaugeValue, float64(m.RChar), grouper, group)
	ch <- prometheus.MustNewConstMetric(descs["grouped_process_io_w_char"], prometheus.GaugeValue, float64(m.WChar), grouper, group)
	ch <- prometheus.MustNewConstMetric(descs["grouped_process_io_sysc_r"], prometheus.GaugeValue, float64(m.SyscW), grouper, group)
	ch <- prometheus.MustNewConstMetric(descs["grouped_process_io_sysc_w"], prometheus.GaugeValue, float64(m.SyscW), grouper, group)
	ch <- prometheus.MustNewConstMetric(descs["grouped_process_io_read_bytes"], prometheus.GaugeValue, float64(m.ReadBytes), grouper, group)
	ch <- prometheus.MustNewConstMetric(descs["grouped_process_io_write_bytes"], prometheus.GaugeValue, float64(m.WriteBytes), grouper, group)
	ch <- prometheus.MustNewConstMetric(descs["grouped_process_io_cancelled_write_bytes"], prometheus.GaugeValue, float64(m.CancelledWriteBytes), grouper, group)
	return nil
}

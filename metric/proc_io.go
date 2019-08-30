package metric

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/procfs"
)

// ProcIOMetric is metric
type ProcIOMetric struct {
	metrics map[int]procfs.ProcIO
}

func (m *ProcIOMetric) Describe() map[string]*prometheus.Desc {
	descs := map[string]*prometheus.Desc{
		"grouped_process_io_r_char": prometheus.NewDesc(
			"grouped_process_io_r_char",
			"Grouped /proc/[PID]/io.rchar",
			[]string{"grouper", "group"}, nil,
		),
		"grouped_process_io_w_char": prometheus.NewDesc(
			"grouped_process_io_w_char",
			"Grouped /proc/[PID]/io.wchar",
			[]string{"grouper", "group"}, nil,
		),
		"grouped_process_io_sysc_r": prometheus.NewDesc(
			"grouped_process_io_sysc_r",
			"Grouped /proc/[PID]/io.syscr",
			[]string{"grouper", "group"}, nil,
		),
		"grouped_process_io_sysc_w": prometheus.NewDesc(
			"grouped_process_io_sysc_w",
			"Grouped /proc/[PID]/io.syscw",
			[]string{"grouper", "group"}, nil,
		),
		"grouped_process_io_read_bytes": prometheus.NewDesc(
			"grouped_process_io_read_bytes",
			"Grouped /proc/[PID]/io.read_bytes",
			[]string{"grouper", "group"}, nil,
		),
		"grouped_process_io_write_bytes": prometheus.NewDesc(
			"grouped_process_io_write_bytes",
			"Grouped /proc/[PID]/io.write_bytes",
			[]string{"grouper", "group"}, nil,
		),
		"grouped_process_io_cancelled_write_bytes": prometheus.NewDesc(
			"grouped_process_io_cancelled_write_bytes",
			"Grouped /proc/[PID]/io.cancelled_write_bytes",
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
	m.metrics[proc.PID] = pio
	return nil
}

func (m *ProcIOMetric) SetCollectedMetric(ch chan<- prometheus.Metric, descs map[string]*prometheus.Desc, grouper string, group string) error {
	var (
		rChar               float64
		wChar               float64
		syscR               float64
		syscW               float64
		readBytes           float64
		writeBytes          float64
		cancelledWriteBytes float64
	)

	for _, metric := range m.metrics {
		rChar = rChar + float64(metric.RChar)
		wChar = wChar + float64(metric.WChar)
		syscR = syscR + float64(metric.SyscR)
		syscW = syscW + float64(metric.SyscW)
		readBytes = readBytes + float64(metric.ReadBytes)
		writeBytes = writeBytes + float64(metric.WriteBytes)
		cancelledWriteBytes = cancelledWriteBytes + float64(metric.CancelledWriteBytes)
	}

	ch <- prometheus.MustNewConstMetric(descs["grouped_process_io_r_char"], prometheus.CounterValue, float64(rChar), grouper, group)
	ch <- prometheus.MustNewConstMetric(descs["grouped_process_io_w_char"], prometheus.CounterValue, float64(wChar), grouper, group)
	ch <- prometheus.MustNewConstMetric(descs["grouped_process_io_sysc_r"], prometheus.CounterValue, float64(syscW), grouper, group)
	ch <- prometheus.MustNewConstMetric(descs["grouped_process_io_sysc_w"], prometheus.CounterValue, float64(syscW), grouper, group)
	ch <- prometheus.MustNewConstMetric(descs["grouped_process_io_read_bytes"], prometheus.CounterValue, float64(readBytes), grouper, group)
	ch <- prometheus.MustNewConstMetric(descs["grouped_process_io_write_bytes"], prometheus.CounterValue, float64(writeBytes), grouper, group)
	ch <- prometheus.MustNewConstMetric(descs["grouped_process_io_cancelled_write_bytes"], prometheus.CounterValue, float64(cancelledWriteBytes), grouper, group)
	return nil
}

func NewProcIOMetric() *ProcIOMetric {
	return &ProcIOMetric{
		metrics: make(map[int]procfs.ProcIO),
	}
}

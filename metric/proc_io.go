package metric

import (
	"sync"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/procfs"
)

// ProcIOMetric is metric
type ProcIOMetric struct {
	sync.Mutex
	metrics map[int]procfs.ProcIO
}

func (m *ProcIOMetric) Describe() map[string]*prometheus.Desc {
	descs := map[string]*prometheus.Desc{
		"grouped_process_io_r_char_total": prometheus.NewDesc(
			"grouped_process_io_r_char_total",
			"Total number of grouped /proc/[PID]/io.rchar",
			[]string{"grouper", "group"}, nil,
		),
		"grouped_process_io_w_char_total": prometheus.NewDesc(
			"grouped_process_io_w_char_total",
			"Total number of grouped /proc/[PID]/io.wchar",
			[]string{"grouper", "group"}, nil,
		),
		"grouped_process_io_sysc_r_total": prometheus.NewDesc(
			"grouped_process_io_sysc_r_total",
			"Total number of grouped /proc/[PID]/io.syscr",
			[]string{"grouper", "group"}, nil,
		),
		"grouped_process_io_sysc_w_total": prometheus.NewDesc(
			"grouped_process_io_sysc_w_total",
			"Total number of grouped /proc/[PID]/io.syscw",
			[]string{"grouper", "group"}, nil,
		),
		"grouped_process_io_read_bytes_total": prometheus.NewDesc(
			"grouped_process_io_read_bytes_total",
			"Total number of grouped /proc/[PID]/io.read_bytes",
			[]string{"grouper", "group"}, nil,
		),
		"grouped_process_io_write_bytes_total": prometheus.NewDesc(
			"grouped_process_io_write_bytes_total",
			"Total number of grouped /proc/[PID]/io.write_bytes",
			[]string{"grouper", "group"}, nil,
		),
		"grouped_process_io_cancelled_write_bytes_total": prometheus.NewDesc(
			"grouped_process_io_cancelled_write_bytes_total",
			"Total number of grouped /proc/[PID]/io.cancelled_write_bytes",
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
	m.Lock()
	m.metrics[proc.PID] = pio
	m.Unlock()
	return nil
}

func (m *ProcIOMetric) PushCollected(ch chan<- prometheus.Metric, descs map[string]*prometheus.Desc, grouper string, group string) error {
	var (
		rChar               float64
		wChar               float64
		syscR               float64
		syscW               float64
		readBytes           float64
		writeBytes          float64
		cancelledWriteBytes float64
	)

	m.Lock()
	for _, metric := range m.metrics {
		rChar = rChar + float64(metric.RChar)
		wChar = wChar + float64(metric.WChar)
		syscR = syscR + float64(metric.SyscR)
		syscW = syscW + float64(metric.SyscW)
		readBytes = readBytes + float64(metric.ReadBytes)
		writeBytes = writeBytes + float64(metric.WriteBytes)
		cancelledWriteBytes = cancelledWriteBytes + float64(metric.CancelledWriteBytes)
	}
	m.Unlock()

	ch <- prometheus.MustNewConstMetric(descs["grouped_process_io_r_char_total"], prometheus.CounterValue, rChar, grouper, group)
	ch <- prometheus.MustNewConstMetric(descs["grouped_process_io_w_char_total"], prometheus.CounterValue, wChar, grouper, group)
	ch <- prometheus.MustNewConstMetric(descs["grouped_process_io_sysc_r_total"], prometheus.CounterValue, syscR, grouper, group)
	ch <- prometheus.MustNewConstMetric(descs["grouped_process_io_sysc_w_total"], prometheus.CounterValue, syscW, grouper, group)
	ch <- prometheus.MustNewConstMetric(descs["grouped_process_io_read_bytes_total"], prometheus.CounterValue, readBytes, grouper, group)
	ch <- prometheus.MustNewConstMetric(descs["grouped_process_io_write_bytes_total"], prometheus.CounterValue, writeBytes, grouper, group)
	ch <- prometheus.MustNewConstMetric(descs["grouped_process_io_cancelled_write_bytes_total"], prometheus.CounterValue, cancelledWriteBytes, grouper, group)

	return nil
}

func (m *ProcIOMetric) RequiredWeight() int64 {
	return 1
}

func NewProcIOMetric() *ProcIOMetric {
	return &ProcIOMetric{
		metrics: make(map[int]procfs.ProcIO),
	}
}

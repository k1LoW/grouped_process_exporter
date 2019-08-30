package metric

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/procfs"
)

// ProcIOMetric is metric
type ProcIOMetric struct {
	procfs.ProcIO
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
	ch <- prometheus.MustNewConstMetric(descs["grouped_process_io_r_char"], prometheus.CounterValue, float64(m.RChar), grouper, group)
	ch <- prometheus.MustNewConstMetric(descs["grouped_process_io_w_char"], prometheus.CounterValue, float64(m.WChar), grouper, group)
	ch <- prometheus.MustNewConstMetric(descs["grouped_process_io_sysc_r"], prometheus.CounterValue, float64(m.SyscW), grouper, group)
	ch <- prometheus.MustNewConstMetric(descs["grouped_process_io_sysc_w"], prometheus.CounterValue, float64(m.SyscW), grouper, group)
	ch <- prometheus.MustNewConstMetric(descs["grouped_process_io_read_bytes"], prometheus.CounterValue, float64(m.ReadBytes), grouper, group)
	ch <- prometheus.MustNewConstMetric(descs["grouped_process_io_write_bytes"], prometheus.CounterValue, float64(m.WriteBytes), grouper, group)
	ch <- prometheus.MustNewConstMetric(descs["grouped_process_io_cancelled_write_bytes"], prometheus.CounterValue, float64(m.CancelledWriteBytes), grouper, group)
	return nil
}

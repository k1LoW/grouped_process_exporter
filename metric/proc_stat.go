package metric

import (
	"os/exec"
	"strconv"
	"strings"
	"sync"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/procfs"
)

// ProcStatMetric is metric
type ProcStatMetric struct {
	sync.Mutex
	metrics map[int]procfs.ProcStat
	clkTck  float64
}

func (m *ProcStatMetric) Describe() map[string]*prometheus.Desc {
	descs := map[string]*prometheus.Desc{
		"grouped_process_stat_minflt_total": prometheus.NewDesc(
			"grouped_process_stat_minflt_total",
			"Total number of grouped /proc/[PID]/stat.minflt",
			[]string{"grouper", "group"}, nil,
		),
		"grouped_process_stat_cminflt_total": prometheus.NewDesc(
			"grouped_process_stat_cminflt_total",
			"Total number of grouped /proc/[PID]/stat.rchar",
			[]string{"grouper", "group"}, nil,
		),
		"grouped_process_stat_majflt_total": prometheus.NewDesc(
			"grouped_process_stat_majflt_total",
			"Total number of grouped /proc/[PID]/stat.majflt",
			[]string{"grouper", "group"}, nil,
		),
		"grouped_process_stat_cmajflt_total": prometheus.NewDesc(
			"grouped_process_stat_cmajflt_total",
			"Total number of grouped /proc/[PID]/stat.cmajflt",
			[]string{"grouper", "group"}, nil,
		),
		"grouped_process_stat_utime_total": prometheus.NewDesc(
			"grouped_process_stat_utime_total",
			"Total number of grouped /proc/[PID]/stat.utime",
			[]string{"grouper", "group"}, nil,
		),
		"grouped_process_stat_stime_total": prometheus.NewDesc(
			"grouped_process_stat_stime_total",
			"Total number of grouped /proc/[PID]/stat.stime",
			[]string{"grouper", "group"}, nil,
		),
		"grouped_process_stat_cutime_total": prometheus.NewDesc(
			"grouped_process_stat_cutime_total",
			"Total number of grouped /proc/[PID]/stat.cutime",
			[]string{"grouper", "group"}, nil,
		),
		"grouped_process_stat_cstime_total": prometheus.NewDesc(
			"grouped_process_stat_cstime_total",
			"Total number of grouped /proc/[PID]/stat.cstime",
			[]string{"grouper", "group"}, nil,
		),
		"grouped_process_stat_numthreads": prometheus.NewDesc(
			"grouped_process_stat_numthreads",
			"Grouped /proc/[PID]/stat.numthreads",
			[]string{"grouper", "group"}, nil,
		),
		"grouped_process_stat_vsize_bytes": prometheus.NewDesc(
			"grouped_process_stat_vsize_bytes",
			"Grouped /proc/[PID]/stat.vsize",
			[]string{"grouper", "group"}, nil,
		),
		"grouped_process_stat_rss": prometheus.NewDesc(
			"grouped_process_stat_rss",
			"Grouped /proc/[PID]/stat.rss",
			[]string{"grouper", "group"}, nil,
		),
		"grouped_process_stat_clk_tck": prometheus.NewDesc(
			"grouped_process_stat_clk_tck",
			"clock ticks (divide by sysconf(_SC_CLK_TCK))",
			[]string{"grouper", "group"}, nil,
		),
	}
	return descs
}

func (m *ProcStatMetric) String() string {
	return "proc_stat"
}

func (m *ProcStatMetric) CollectFromProc(proc procfs.Proc) error {
	stat, err := proc.Stat()
	if err != nil {
		return err
	}
	m.Lock()
	m.metrics[proc.PID] = stat
	m.Unlock()
	return nil
}

func (m *ProcStatMetric) PushCollected(ch chan<- prometheus.Metric, descs map[string]*prometheus.Desc, grouper string, group string) error {
	var (
		minFlt     float64
		cMinFlt    float64
		majFlt     float64
		cMajFlt    float64
		uTime      float64
		sTime      float64
		cUTime     float64
		cSTime     float64
		numThreads float64
		vSize      float64
		rss        float64
	)
	m.Lock()
	for _, metric := range m.metrics {
		minFlt = minFlt + float64(metric.MinFlt)
		cMinFlt = cMinFlt + float64(metric.CMinFlt)
		majFlt = majFlt + float64(metric.MajFlt)
		cMajFlt = cMajFlt + float64(metric.CMajFlt)
		uTime = uTime + float64(metric.UTime)
		sTime = sTime + float64(metric.STime)
		cUTime = cUTime + float64(metric.CUTime)
		cSTime = cSTime + float64(metric.CSTime)
		numThreads = numThreads + float64(metric.NumThreads)
		vSize = vSize + float64(metric.VSize)
		rss = rss + float64(metric.RSS)
	}
	m.Unlock()

	ch <- prometheus.MustNewConstMetric(descs["grouped_process_stat_minflt_total"], prometheus.CounterValue, minFlt, grouper, group)
	ch <- prometheus.MustNewConstMetric(descs["grouped_process_stat_cminflt_total"], prometheus.CounterValue, cMinFlt, grouper, group)
	ch <- prometheus.MustNewConstMetric(descs["grouped_process_stat_majflt_total"], prometheus.CounterValue, majFlt, grouper, group)
	ch <- prometheus.MustNewConstMetric(descs["grouped_process_stat_cmajflt_total"], prometheus.CounterValue, cMajFlt, grouper, group)
	ch <- prometheus.MustNewConstMetric(descs["grouped_process_stat_utime_total"], prometheus.CounterValue, uTime, grouper, group)
	ch <- prometheus.MustNewConstMetric(descs["grouped_process_stat_stime_total"], prometheus.CounterValue, sTime, grouper, group)
	ch <- prometheus.MustNewConstMetric(descs["grouped_process_stat_cutime_total"], prometheus.CounterValue, cUTime, grouper, group)
	ch <- prometheus.MustNewConstMetric(descs["grouped_process_stat_cstime_total"], prometheus.CounterValue, cSTime, grouper, group)
	ch <- prometheus.MustNewConstMetric(descs["grouped_process_stat_numthreads"], prometheus.GaugeValue, numThreads, grouper, group)
	ch <- prometheus.MustNewConstMetric(descs["grouped_process_stat_vsize_bytes"], prometheus.GaugeValue, vSize, grouper, group)
	ch <- prometheus.MustNewConstMetric(descs["grouped_process_stat_rss"], prometheus.GaugeValue, rss, grouper, group)
	ch <- prometheus.MustNewConstMetric(descs["grouped_process_stat_clk_tck"], prometheus.GaugeValue, m.clkTck, grouper, group)

	return nil
}

// NewProcStatMetric
func NewProcStatMetric() *ProcStatMetric {
	return &ProcStatMetric{
		metrics: make(map[int]procfs.ProcStat),
		clkTck:  ClkTck(),
	}
}

func (m *ProcStatMetric) RequiredWeight() int64 {
	return 1
}

// ClkTck return clocks per sec (CLK_TCK)
func ClkTck() float64 {
	tck := float64(128)
	out, err := exec.Command("/usr/bin/getconf", "CLK_TCK").Output() // #nosec
	if err == nil {
		i, err := strconv.ParseFloat(strings.TrimSuffix(string(out), "\n"), 64)
		if err == nil {
			tck = float64(i)
		}
	}
	return tck
}

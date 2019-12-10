package metric

import (
	"sync"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/procfs"
)

// ProcStatusMetric is metric
type ProcStatusMetric struct {
	sync.Mutex
	metrics map[int]procfs.ProcStatus
}

func (m *ProcStatusMetric) Describe() map[string]*prometheus.Desc {
	descs := map[string]*prometheus.Desc{
		"grouped_process_status_VmPeak_bytes_total": prometheus.NewDesc(
			"grouped_process_status_VmPeak_bytes_total",
			"Total size of grouped /proc/[PID]/status.VmPeak // Peak virtual memory size.",
			[]string{"grouper", "group"}, nil,
		),
		"grouped_process_status_VmSize_bytes_total": prometheus.NewDesc(
			"grouped_process_status_VmSize_bytes_total",
			"Total size of grouped /proc/[PID]/status.VmSize // Virtual memory size.",
			[]string{"grouper", "group"}, nil,
		),
		"grouped_process_status_VmLck_bytes_total": prometheus.NewDesc(
			"grouped_process_status_VmLck_bytes_total",
			"Total size of grouped /proc/[PID]/status.VmLck // Locked memory size.",
			[]string{"grouper", "group"}, nil,
		),
		"grouped_process_status_VmPin_bytes_total": prometheus.NewDesc(
			"grouped_process_status_VmPin_bytes_total",
			"Total size of grouped /proc/[PID]/status.VmPin // Pinned memory size.",
			[]string{"grouper", "group"}, nil,
		),
		"grouped_process_status_VmHWM_bytes_total": prometheus.NewDesc(
			"grouped_process_status_VmHWM_bytes_total",
			"Total size of grouped /proc/[PID]/status.VmHWM // Peak resident set size.",
			[]string{"grouper", "group"}, nil,
		),
		"grouped_process_status_VmRSS_bytes_total": prometheus.NewDesc(
			"grouped_process_status_VmRSS_bytes_total",
			"Total size of grouped /proc/[PID]/status.VmRSS // Resident set size (sum of RssAnnon RssFile and RssShmem).",
			[]string{"grouper", "group"}, nil,
		),
		"grouped_process_status_RssAnon_bytes_total": prometheus.NewDesc(
			"grouped_process_status_RssAnon_bytes_total",
			"Total size of grouped /proc/[PID]/status.RssAnon // Size of resident anonymous memory.",
			[]string{"grouper", "group"}, nil,
		),
		"grouped_process_status_RssFile_bytes_total": prometheus.NewDesc(
			"grouped_process_status_RssFile_bytes_total",
			"Total size of grouped /proc/[PID]/status.RssFile // Size of resident file mappings.",
			[]string{"grouper", "group"}, nil,
		),
		"grouped_process_status_RssShmem_bytes_total": prometheus.NewDesc(
			"grouped_process_status_RssShmem_bytes_total",
			"Total size of grouped /proc/[PID]/status.RssShmem // Size of resident shared memory.",
			[]string{"grouper", "group"}, nil,
		),
		"grouped_process_status_VmData_bytes_total": prometheus.NewDesc(
			"grouped_process_status_VmData_bytes_total",
			"Total size of grouped /proc/[PID]/status.VmData // Size of data segments.",
			[]string{"grouper", "group"}, nil,
		),
		"grouped_process_status_VmStk_bytes_total": prometheus.NewDesc(
			"grouped_process_status_VmStk_bytes_total",
			"Total size of grouped /proc/[PID]/status.VmStk // Size of stack segments.",
			[]string{"grouper", "group"}, nil,
		),
		"grouped_process_status_VmExe_bytes_total": prometheus.NewDesc(
			"grouped_process_status_VmExe_bytes_total",
			"Total size of grouped /proc/[PID]/status.VmExe // Size of text segments.",
			[]string{"grouper", "group"}, nil,
		),
		"grouped_process_status_VmLib_bytes_total": prometheus.NewDesc(
			"grouped_process_status_VmLib_bytes_total",
			"Total size of grouped /proc/[PID]/status.VmLib // Shared library code size.",
			[]string{"grouper", "group"}, nil,
		),
		"grouped_process_status_VmPTE_bytes_total": prometheus.NewDesc(
			"grouped_process_status_VmPTE_bytes_total",
			"Total size of grouped /proc/[PID]/status.VmPTE // Page table entries size.",
			[]string{"grouper", "group"}, nil,
		),
		"grouped_process_status_VmPMD_bytes_total": prometheus.NewDesc(
			"grouped_process_status_VmPMD_bytes_total",
			"Total size of grouped /proc/[PID]/status.VmPMD // Size of second-level page tables.",
			[]string{"grouper", "group"}, nil,
		),
		"grouped_process_status_VmSwap_bytes_total": prometheus.NewDesc(
			"grouped_process_status_VmSwap_bytes_total",
			"Total size of grouped /proc/[PID]/status.VmSwap // Swapped-out virtual memory size by anonymous private.",
			[]string{"grouper", "group"}, nil,
		),
		"grouped_process_status_HugetlbPages_bytes_total": prometheus.NewDesc(
			"grouped_process_status_HugetlbPages_bytes_total",
			"Total size of grouped /proc/[PID]/status.HugetlbPages // Size of hugetlb memory portions",
			[]string{"grouper", "group"}, nil,
		),
		"grouped_process_status_VoluntaryCtxtSwitches_total": prometheus.NewDesc(
			"grouped_process_status_VoluntaryCtxtSwitches_total",
			"Total number of grouped /proc/[PID]/status.VoluntaryCtxtSwitches // Number of voluntary context switches.",
			[]string{"grouper", "group"}, nil,
		),
		"grouped_process_status_NonVoluntaryCtxtSwitches_total": prometheus.NewDesc(
			"grouped_process_status_NonVoluntaryCtxtSwitches_total",
			"Total number of grouped /proc/[PID]/status.NonVoluntaryCtxtSwitches // Number of involuntary context switches.",
			[]string{"grouper", "group"}, nil,
		),
	}
	return descs
}

func (m *ProcStatusMetric) String() string {
	return "proc_status"
}

func (m *ProcStatusMetric) CollectFromProc(proc procfs.Proc) error {
	status, err := proc.NewStatus()
	if err != nil {
		return err
	}
	m.Lock()
	m.metrics[proc.PID] = status
	m.Unlock()
	return nil
}

func (m *ProcStatusMetric) PushCollected(ch chan<- prometheus.Metric, descs map[string]*prometheus.Desc, grouper string, group string) error {
	var (
		vmPeak                   float64
		vmSize                   float64
		vmLck                    float64
		vmPin                    float64
		vmHWM                    float64
		vmRSS                    float64
		rssAnon                  float64
		rssFile                  float64
		rssShmem                 float64
		vmData                   float64
		vmStk                    float64
		vmExe                    float64
		vmLib                    float64
		vmPTE                    float64
		vmPMD                    float64
		vmSwap                   float64
		hugetlbPages             float64
		voluntaryCtxtSwitches    float64
		nonVoluntaryCtxtSwitches float64
	)
	m.Lock()
	for _, metric := range m.metrics {
		vmPeak = vmPeak + float64(metric.VmPeak)
		vmSize = vmSize + float64(metric.VmSize)
		vmLck = vmLck + float64(metric.VmLck)
		vmPin = vmPin + float64(metric.VmPin)
		vmHWM = vmHWM + float64(metric.VmHWM)
		vmRSS = vmRSS + float64(metric.VmRSS)
		rssAnon = rssAnon + float64(metric.RssAnon)
		rssFile = rssFile + float64(metric.RssFile)
		rssShmem = rssShmem + float64(metric.RssShmem)
		vmData = vmData + float64(metric.VmData)
		vmStk = vmStk + float64(metric.VmStk)
		vmExe = vmExe + float64(metric.VmExe)
		vmLib = vmLib + float64(metric.VmLib)
		vmPTE = vmPTE + float64(metric.VmPTE)
		vmPMD = vmPMD + float64(metric.VmPMD)
		vmSwap = vmSwap + float64(metric.VmSwap)
		hugetlbPages = hugetlbPages + float64(metric.HugetlbPages)
		voluntaryCtxtSwitches = voluntaryCtxtSwitches + float64(metric.VoluntaryCtxtSwitches)
		nonVoluntaryCtxtSwitches = nonVoluntaryCtxtSwitches + float64(metric.NonVoluntaryCtxtSwitches)
	}
	m.Unlock()

	ch <- prometheus.MustNewConstMetric(descs["grouped_process_status_VmPeak_bytes_total"], prometheus.GaugeValue, vmPeak, grouper, group)
	ch <- prometheus.MustNewConstMetric(descs["grouped_process_status_VmSize_bytes_total"], prometheus.GaugeValue, vmSize, grouper, group)
	ch <- prometheus.MustNewConstMetric(descs["grouped_process_status_VmLck_bytes_total"], prometheus.GaugeValue, vmLck, grouper, group)
	ch <- prometheus.MustNewConstMetric(descs["grouped_process_status_VmPin_bytes_total"], prometheus.GaugeValue, vmPin, grouper, group)
	ch <- prometheus.MustNewConstMetric(descs["grouped_process_status_VmHWM_bytes_total"], prometheus.GaugeValue, vmHWM, grouper, group)
	ch <- prometheus.MustNewConstMetric(descs["grouped_process_status_VmRSS_bytes_total"], prometheus.GaugeValue, vmRSS, grouper, group)
	ch <- prometheus.MustNewConstMetric(descs["grouped_process_status_RssAnon_bytes_total"], prometheus.GaugeValue, rssAnon, grouper, group)
	ch <- prometheus.MustNewConstMetric(descs["grouped_process_status_RssFile_bytes_total"], prometheus.GaugeValue, rssFile, grouper, group)
	ch <- prometheus.MustNewConstMetric(descs["grouped_process_status_RssShmem_bytes_total"], prometheus.GaugeValue, rssShmem, grouper, group)
	ch <- prometheus.MustNewConstMetric(descs["grouped_process_status_VmData_bytes_total"], prometheus.GaugeValue, vmData, grouper, group)
	ch <- prometheus.MustNewConstMetric(descs["grouped_process_status_VmStk_bytes_total"], prometheus.GaugeValue, vmStk, grouper, group)
	ch <- prometheus.MustNewConstMetric(descs["grouped_process_status_VmExe_bytes_total"], prometheus.GaugeValue, vmExe, grouper, group)
	ch <- prometheus.MustNewConstMetric(descs["grouped_process_status_VmLib_bytes_total"], prometheus.GaugeValue, vmLib, grouper, group)
	ch <- prometheus.MustNewConstMetric(descs["grouped_process_status_VmPTE_bytes_total"], prometheus.GaugeValue, vmPTE, grouper, group)
	ch <- prometheus.MustNewConstMetric(descs["grouped_process_status_VmPMD_bytes_total"], prometheus.GaugeValue, vmPMD, grouper, group)
	ch <- prometheus.MustNewConstMetric(descs["grouped_process_status_VmSwap_bytes_total"], prometheus.GaugeValue, vmSwap, grouper, group)
	ch <- prometheus.MustNewConstMetric(descs["grouped_process_status_HugetlbPages_bytes_total"], prometheus.GaugeValue, hugetlbPages, grouper, group)
	ch <- prometheus.MustNewConstMetric(descs["grouped_process_status_VoluntaryCtxtSwitches_total"], prometheus.CounterValue, voluntaryCtxtSwitches, grouper, group)
	ch <- prometheus.MustNewConstMetric(descs["grouped_process_status_NonVoluntaryCtxtSwitches_total"], prometheus.CounterValue, nonVoluntaryCtxtSwitches, grouper, group)

	return nil
}

// NewProcStatusMetric
func NewProcStatusMetric() *ProcStatusMetric {
	return &ProcStatusMetric{
		metrics: make(map[int]procfs.ProcStatus),
	}
}

func (m *ProcStatusMetric) RequiredWeight() int64 {
	return 1
}

package collector

import (
	"sync"

	"github.com/k1LoW/grouped_process_exporter/grouped_proc"
	"github.com/k1LoW/grouped_process_exporter/grouper"
	"github.com/k1LoW/grouped_process_exporter/metric"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
)

type GroupedProcCollector struct {
	sync.Mutex
	GroupedProcs *grouped_proc.GroupedProcs
	Metrics      map[metric.MetricKey]metric.Metric
	Enabled      map[metric.MetricKey]bool
	Grouper      grouper.Grouper
	descs        map[string]*prometheus.Desc
}

func (c *GroupedProcCollector) Describe(ch chan<- *prometheus.Desc) {
	for _, key := range metric.MetricKeys {
		if c.Enabled[key] {
			descs := c.Metrics[key].Describe()
			for name, desc := range descs {
				c.descs[name] = desc
				ch <- desc
			}
		}
	}
}

func (c *GroupedProcCollector) Collect(ch chan<- prometheus.Metric) {
	c.Lock()
	log.Debugln("Start collecting")
	_ = c.Grouper.Collect(c.GroupedProcs, c.Enabled)
	c.GroupedProcs.Range(func(group string, gproc *grouped_proc.GroupedProc) bool {
		log.Debugf("Collect grouped process: %s\n", group)
		if !gproc.Exists {
			c.GroupedProcs.Delete(group)
			log.Debugf("Delete grouped process: %s\n", group)
			return true
		}
		for key, metric := range gproc.Metrics {
			if gproc.Enabled[key] {
				err := metric.PushCollected(ch, c.descs, c.Grouper.Name(), group)
				if err != nil {
					log.Errorf("Failed to push collected metrics: %v\n", err)
					// TODO: metric.PushDefaultMetric(ch, c.descs, c.Grouper.Name(), group)
					return true
				}
			}
		}
		gproc.Exists = false
		return true
	})
	c.Unlock()
}

func (c *GroupedProcCollector) EnableMetric(metric metric.MetricKey) {
	c.Enabled[metric] = true
}

func (c *GroupedProcCollector) DisableMetric(metric metric.MetricKey) {
	c.Enabled[metric] = false
}

// NewGroupedProcCollector
func NewGroupedProcCollector(g grouper.Grouper) (*GroupedProcCollector, error) {
	return &GroupedProcCollector{
		GroupedProcs: grouped_proc.NewGroupedProcs(),
		Metrics:      metric.AvairableMetrics(),
		Enabled:      metric.DefaultEnabledMetrics(),
		Grouper:      g,
		descs:        make(map[string]*prometheus.Desc),
	}, nil
}

package collector

import (
	"fmt"

	"github.com/k1LoW/grouped_process_exporter/grouped_proc"
	"github.com/k1LoW/grouped_process_exporter/grouper"
	"github.com/k1LoW/grouped_process_exporter/metric"
	"github.com/prometheus/client_golang/prometheus"
)

type GroupedProcCollector struct {
	GroupedProcs map[string]*grouped_proc.GroupedProc
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
	c.GroupedProcs = make(map[string]*grouped_proc.GroupedProc) // clear
	_ = c.Grouper.Collect(c.GroupedProcs, c.Enabled)
	for group, proc := range c.GroupedProcs {
		for key, metric := range proc.Metrics {
			if proc.Enabled[key] {
				err := metric.SetCollectedMetric(ch, c.descs, c.Grouper.Name(), group)
				if err != nil {
					// TODO: metric.SetDefaultMetric(ch, c.descs, c.Grouper.Name(), group)
					continue
				}
			}
		}
	}
}

func (c *GroupedProcCollector) Debug() {
	_ = c.Grouper.Collect(c.GroupedProcs, c.Enabled)
	for group, proc := range c.GroupedProcs {
		fmt.Printf("%s: %#v\n", group, proc)
	}
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
		GroupedProcs: make(map[string]*grouped_proc.GroupedProc),
		Metrics:      metric.AvairableMetrics(),
		Enabled:      grouped_proc.DefaultEnabledMetrics(),
		Grouper:      g,
		descs:        make(map[string]*prometheus.Desc),
	}, nil
}

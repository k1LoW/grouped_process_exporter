package proc_status_name

import (
	"errors"
	"regexp"
	"sync"

	"github.com/k1LoW/grouped_process_exporter/grouped_proc"
	"github.com/k1LoW/grouped_process_exporter/metric"
	"github.com/prometheus/procfs"
)

type ProcStatusName struct {
	nRe *regexp.Regexp
}

func (g *ProcStatusName) Name() string {
	return "proc_status_name"
}

func (g *ProcStatusName) Collect(gpMap map[string]*grouped_proc.GroupedProc, enabled map[metric.MetricKey]bool) error {
	wg := &sync.WaitGroup{}

	procs, err := procfs.AllProcs()
	if err != nil {
		return err
	}
	for _, proc := range procs {
		status, err := proc.NewStatus()
		if err != nil {
			// TODO: Log
			continue
		}
		pid := proc.PID
		name := status.Name
		if g.nRe != nil {
			matches := g.nRe.FindStringSubmatch(name)
			if len(matches) > 1 {
				name = matches[1]
			}
		}
		_, ok := gpMap[name]
		if !ok {
			gpMap[name] = grouped_proc.NewGroupedProc(enabled)
		}
		wg.Add(1)
		go func(wg *sync.WaitGroup, pid int, g *grouped_proc.GroupedProc) {
			_ = g.AppendPid(pid)
			wg.Done()
		}(wg, pid, gpMap[name])
	}
	wg.Wait()
	return nil
}

func (g *ProcStatusName) SetNormalizeRegexp(nReStr string) error {
	if nReStr == "" {
		return nil
	}
	nRe, err := regexp.Compile(nReStr)
	if err != nil {
		return err
	}
	if nRe.NumSubexp() != 1 {
		return errors.New("number of parenthesized subexpressions in this regexp should be 1")
	}
	g.nRe = nRe
	return nil
}

// NewProcStatusName
func NewProcStatusName() *ProcStatusName {
	return &ProcStatusName{}
}

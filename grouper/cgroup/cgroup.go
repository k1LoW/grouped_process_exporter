package cgroup

import (
	"bufio"
	"context"
	"errors"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"

	"github.com/k1LoW/grouped_process_exporter/grouped_proc"
	"github.com/k1LoW/grouped_process_exporter/metric"
	"github.com/prometheus/common/log"
	"golang.org/x/sync/semaphore"
)

// DefaultSubsystems cgroups subsystems default list
var DefaultSubsystems = []string{
	"cpuset",
	"cpu",
	"cpuacct",
	"blkio",
	"memory",
	"devices",
	"freezer",
	"net_cls",
	"net_prio",
	"perf_event",
	"hugetlb",
	"pids",
	"rdma",
}

type Cgroup struct {
	fsPath     string
	subsystems []string
	nRe        *regexp.Regexp // normalize regexp
	eRe        *regexp.Regexp // exclude regexp
}

func (c *Cgroup) Name() string {
	return "cgroup"
}

func (c *Cgroup) Collect(gprocs *grouped_proc.GroupedProcs, enabled map[metric.MetricKey]bool, sem *semaphore.Weighted) error {
	wg := &sync.WaitGroup{}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	log.Debugf("Cgroup subsystems %s\n", c.subsystems)
	realSubsystems := []string{}
	for _, s := range c.subsystems {
		path := filepath.Clean(filepath.Join(c.fsPath, s))
		f, err := os.Lstat(path)
		if err != nil {
			log.Debugf("%s\n", err)
			continue
		}
		if f.Mode()&os.ModeSymlink == os.ModeSymlink {
			realpath, err := filepath.EvalSymlinks(path)
			if err != nil {
				return err
			}
			f, err = os.Lstat(realpath)
			if err != nil {
				return err
			}
			path = realpath
		}
		if f.IsDir() && !contains(realSubsystems, filepath.Base(path)) {
			realSubsystems = append(realSubsystems, filepath.Base(path))
		}
	}
	log.Debugf("Resolve symlinks %s\n", realSubsystems)

	for _, s := range realSubsystems {
		searchDir := filepath.Clean(filepath.Join(c.fsPath, s))
		err := filepath.Walk(searchDir, func(path string, f os.FileInfo, err error) error {
			if err != nil {
				return nil
			}
			if f == nil {
				return nil
			}
			if err := sem.Acquire(ctx, 2); err != nil {
				return err
			}
			defer sem.Release(2)
			if !f.IsDir() {
				return nil
			}
			cPath := strings.Replace(path, searchDir, "", 1)
			if c.eRe != nil {
				if c.eRe.MatchString(cPath) {
					return nil
				}
			}
			if c.nRe != nil {
				matches := c.nRe.FindStringSubmatch(cPath)
				if len(matches) > 1 {
					cPath = matches[1]
				}
			}
			if cPath == "" {
				return nil
			}
			{
				f, err := os.Open(filepath.Clean(filepath.Join(path, "cgroup.procs")))
				if err != nil {
					_ = f.Close()
					return nil
				}
				var (
					gproc *grouped_proc.GroupedProc
					ok    bool
				)
				gproc, ok = gprocs.Load(cPath)
				if !ok {
					gproc = grouped_proc.NewGroupedProc(enabled)
					gprocs.Store(cPath, gproc)
				}
				if err := gproc.Collect(cPath); err != nil {
					return err
				}
				gproc.Exists = true
				reader := bufio.NewReaderSize(f, 1028)
				for {
					line, _, err := reader.ReadLine()
					if err == io.EOF {
						break
					} else if err != nil {
						_ = f.Close()
						return err
					}
					pid, err := strconv.Atoi(string(line))
					if err != nil {
						_ = f.Close()
						return err
					}
					err = sem.Acquire(ctx, gproc.RequiredWeight)
					if err != nil {
						_ = f.Close()
						return err
					}
					wg.Add(1)
					go func(wg *sync.WaitGroup, pid int, gproc *grouped_proc.GroupedProc) {
						_ = gproc.AppendProcAndCollect(pid)
						sem.Release(gproc.RequiredWeight)
						wg.Done()
					}(wg, pid, gproc)
				}
			}
			return nil
		})
		if err != nil {
			return err
		}
	}

	wg.Wait()
	return nil
}

func (c *Cgroup) SetNormalizeRegexp(nReStr string) error {
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
	c.nRe = nRe
	return nil
}

func (c *Cgroup) SetExcludeRegexp(eReStr string) error {
	if eReStr == "" {
		return nil
	}
	eRe, err := regexp.Compile(eReStr)
	if err != nil {
		return err
	}
	c.eRe = eRe
	return nil
}

// NewCgroup
func NewCgroup(fsPath string, subsystems []string) *Cgroup {
	if len(subsystems) == 0 {
		subsystems = DefaultSubsystems
	}

	return &Cgroup{
		fsPath:     fsPath,
		subsystems: subsystems,
	}
}

func contains(s []string, e string) bool {
	for _, v := range s {
		if e == v {
			return true
		}
	}
	return false
}

package cgroup

import (
	"bufio"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	"github.com/k1LoW/grouped_process_exporter/grouped_procs"
)

// Subsystems cgroups subsystems list
var Subsystems = []string{
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

func ListCgroupedProcs(fsPath string) (map[string]*grouped_procs.GroupedProcs, error) {
	gpMap := make(map[string]*grouped_procs.GroupedProcs)

	wg := &sync.WaitGroup{}

	for _, s := range Subsystems {
		searchDir := filepath.Clean(filepath.Join(fsPath, s))

		err := filepath.Walk(searchDir, func(path string, f os.FileInfo, err error) error {
			if err != nil {
				return nil
			}
			if f == nil {
				return nil
			}
			if f.IsDir() {
				cPath := strings.Replace(path, searchDir, "", 1)
				if cPath != "" {
					f, err := os.Open(filepath.Clean(filepath.Join(path, "cgroup.procs")))
					if err != nil {
						_ = f.Close()
						return nil
					}
					_, ok := gpMap[cPath]
					if ok {
						_ = f.Close()
						return nil
					}
					gpMap[cPath] = new(grouped_procs.GroupedProcs)
					gpMap[cPath].Path = cPath
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

						wg.Add(1)
						go func(wg *sync.WaitGroup, pid int, g *grouped_procs.GroupedProcs) {
							defer wg.Done()
							_ = g.AppendPid(pid)
						}(wg, pid, gpMap[cPath])
					}
				}
				return nil
			}
			return nil
		})
		if err != nil {
			return nil, err
		}
	}

	wg.Wait()
	return gpMap, nil
}

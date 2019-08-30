/*
Copyright © 2019 Ken'ichiro Oyama <k1lowxb@gmail.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/k1LoW/grouped_process_exporter/collector"
	"github.com/k1LoW/grouped_process_exporter/grouper"
	"github.com/k1LoW/grouped_process_exporter/grouper/cgroup"
	"github.com/k1LoW/grouped_process_exporter/metric"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/cobra"
)

var (
	address   string
	endpoint  string
	group     string
	collectIO bool
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "grouped_process_exporter",
	Short: "Exporter for grouped process",
	Long:  `Exporter for grouped process.`,
	Run: func(cmd *cobra.Command, args []string) {
		status, err := runRoot(args, address, endpoint, group, collectIO)
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "%s\n", err)
		}
		os.Exit(status)
	},
}

func runRoot(args []string, address, endpoint, group string, collectIO bool) (int, error) {
	var g grouper.Grouper
	switch group {
	case "cgroup":
		fsPath := "/sys/fs/cgroup"
		g = cgroup.NewCgroup(fsPath)
	default:
		return 1, errors.New("invalid grouping type")
	}

	collector, err := collector.NewGroupedProcCollector(g)
	if err != nil {
		return 1, err
	}
	if collectIO {
		collector.EnableMetric(metric.ProcIO)
	}
	prometheus.MustRegister(collector)
	http.Handle(endpoint, promhttp.Handler())
	log.Fatal(http.ListenAndServe(address, nil))
	return 0, nil
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringVarP(&address, "telemetry.address", "", ":9629", "Address on which to expose metrics.")
	rootCmd.Flags().StringVarP(&endpoint, "telemetry.endpoint", "", "/metrics", "Path under which to expose metrics.")
	rootCmd.Flags().StringVarP(&group, "group.type", "", "cgroup", "Grouping type.")
	rootCmd.Flags().BoolVarP(&collectIO, "collector.io", "", false, "Enable collecting /proc/[PID]/io.")
}

# Exporter for grouped process

## Supported grouping

- control group v1 ( `cgroup`, default )
- /proc/[PID]/status.Name: ( `name`, **this is PoC** )

### :unlock: Advanced grouping

Exporter normalize group names using regexp `--group.normalize` option.

For example, by setting `--group.normalize='^.+(tcpdp).+$'`, Exporter normalized the group names `/path/to/tcpdp-eth0` and `/path/to/tcpdp-eth1` to `tcpdp`.

## Avairable Metrics

### Amount of grouped procs ( default on )

| Name | Type | Description |
| --- | --- | --- |
| grouped_process_procs | Gauge | Amount of grouped procs |

### Grouped /proc/[PID]/stat ( `--collector.stat` )

| Name | Type | Description |
| --- | --- | --- |
| grouped_process_stat_minflt_total | Counter | Total number of grouped /proc/[PID]/stat.minflt |
| grouped_process_stat_cminflt_total | Counter | Total number of grouped /proc/[PID]/stat.rchar |
| grouped_process_stat_majflt_total | Counter | Total number of grouped /proc/[PID]/stat.majflt |
| grouped_process_stat_cmajflt_total | Counter | Total number of grouped /proc/[PID]/stat.cmajflt |
| grouped_process_stat_utime_total | Counter | Total number of grouped /proc/[PID]/stat.utime |
| grouped_process_stat_stime_total | Counter | Total number of grouped /proc/[PID]/stat.stime |
| grouped_process_stat_cutime_total | Counter | Total number of grouped /proc/[PID]/stat.cutime |
| grouped_process_stat_cstime_total | Counter | Total number of grouped /proc/[PID]/stat.cstime |
| grouped_process_stat_numthreads | Gauge | Grouped /proc/[PID]/stat.numthreads |
| grouped_process_stat_vsize_bytes | Gauge | Grouped /proc/[PID]/stat.vsize |
| grouped_process_stat_rss | Gauge | Grouped /proc/[PID]/stat.rss |
| grouped_process_stat_clk_tck | Gauge | clock ticks (divide by sysconf(_SC_CLK_TCK)) |

### Grouped /proc/[PID]/io ( `--collector.io` )

| Name | Type | Description |
| --- | --- | --- |
| grouped_process_io_r_char_total | Counter | Total number of grouped /proc/[PID]/io.rchar |
| grouped_process_io_w_char_total | Counter | Total number of grouped /proc/[PID]/io.wchar |
| grouped_process_io_sysc_r_total | Counter | Total number of grouped /proc/[PID]/io.syscr |
| grouped_process_io_sysc_w_total | Counter | Total number of grouped /proc/[PID]/io.syscw |
| grouped_process_io_read_bytes_total | Counter | Total number of grouped /proc/[PID]/io.read_bytes |
| grouped_process_io_write_bytes_total | Counter | Total number of grouped /proc/[PID]/io.write_bytes |
| grouped_process_io_cancelled_write_bytes_total | Counter | Total number of grouped /proc/[PID]/io.cancelled_write_bytes |

## TODO

- [ ] Test
- [ ] More metrics
- [ ] Logging

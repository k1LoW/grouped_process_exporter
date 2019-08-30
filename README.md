# Exporter for grouped process

## Supported grouping

- control group v1 ( `cgroup`, default )
- /proc/[PID]/status.Name: ( `name`, **this is PoC** )

### Advanced grouping

Exporter normalize group names using regexp `--group.normalize` option.

For example, by setting `--group.normalize='^(/path/to/tcpdp).*$'`, Exporter normalized the group names `/path/to/tcpdp-eth0` and `/path/to/tcpdp-eth1` to `/path/to/tcpdp`.

## Avairable Metrics

### Amount of grouped procs ( default on )

| Name | Type | Description |
| --- | --- | --- |
| grouped_process_procs | Gauge | Amount of grouped procs |

### Grouped /proc/[PID]/io ( `--collector.io` )

| Name | Type | Description |
| --- | --- | --- |
| grouped_process_io_r_char | Gauge | Grouped /proc/[PID]/io.rchar |
| grouped_process_io_w_char | Gauge | Grouped /proc/[PID]/io.wchar |
| grouped_process_io_sysc_r | Gauge | Grouped /proc/[PID]/io.syscr |
| grouped_process_io_sysc_w | Gauge | Grouped /proc/[PID]/io.syscw |
| grouped_process_io_read_bytes | Gauge | Grouped /proc/[PID]/io.read_bytes |
| grouped_process_io_write_bytes | Gauge | Grouped /proc/[PID]/io.write_bytes |
| grouped_process_io_cancelled_write_bytes | Gauge | Grouped /proc/[PID]/io.cancelled_write_bytes |

## TODO

- [ ] Test
- [ ] Logging

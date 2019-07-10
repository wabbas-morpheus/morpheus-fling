# morpheus-fling

A small, command line based binary for aggregating useful statistics about large environments prior to deploying Morpheus.

## Functions

1. Port Scanning
2. OS Statistics

### Port Scanning

`morpheus-fling` reads ips and associated ports and performs a scan of these ports to inspect openness. 

### OS Statistics

`morpheus-fling` aggregates statistics about the OS it is installed on utilizing Linux kernel syscalls including memory, CPU, and available disk.

## Usage

Place binary on filesystem and give posix permissions to be executed.  Binary expects a file to exist called `network.txt`, however user can specify a separate infile.  Similarly, binary will by default write output to a file called `output.txt` but user can flag a separate outfile to be created and appended to.

Run:

```bash
slimshady@morpheus1:~# ./morpheus-fling
```

Or run:

```bash
slimshady@morpheus1:~# ./morpheus-fling --infile=/path/to/foo.txt --outfile=/path/to/bar.txt
```

By default `morpheus-fling` only allows `1024` semaphores to be used based on standard ulimit expectations in environments.  This can also be adjusted by passing a `ulimit` flag.  To use an accurate representation from the system `morpheus-fling` is deployed to run `ulimit -n` on your system and qualify your binary run with this value flagged.
Example:

```bash
slimshady@morpheus1:~# ulimit -n
204800
slimshady@morpheus1:~# ./morpheus-fling --ulimit=204800
```

### Inputs

As described, `morpheus-fling` defaults to looking for an input file called `network.txt`.  Format for the entries in this file should follow `ip:port` notation as below.

```text
10.30.21.100:10092
10.30.21.100:3306
10.30.21.193:22
10.30.21.100:15672
10.30.21.100:5672
```

### Report

`morpheus-fling` will generate a report during the run.  Below is an example of the contents.  The deafult is a file in the same directory as the binary called `output.txt` but this can be adjusted by making use of a flag during the run like `--outfile=/home/slimshady/foobar.txt`

```text
PORT SCANS:
10.30.21.100:3306 open
10.30.21.100:15672 open
10.30.21.100:5672 closed
10.30.21.100:10092 closed
10.30.21.193:22 closed


OS STATS:
{
  "sysinfo": {
    "version": "0.9.2",
    "timestamp": "2019-07-09T16:12:52.316092459-06:00"
  },
  "node": {
    "hostname": "labs-den-demo-morpheus",
    "machineid": "2f67b055ae2d1078d70401de58a63a28",
    "timezone": "America/Denver"
  },
  "os": {
    "name": "Ubuntu 14.04.5 LTS",
    "vendor": "ubuntu",
    "version": "14.04",
    "release": "14.04.5",
    "architecture": "amd64"
  },
  "kernel": {
    "release": "4.2.0-42-generic",
    "version": "#49~14.04.1-Ubuntu SMP Wed Jun 29 20:22:11 UTC 2016",
    "architecture": "x86_64"
  },
  "product": {
    "name": "MBI-6219G-T-Pack",
    "vendor": "Supermicro",
    "version": "0123456789",
    "serial": "S215034X6B37474"
  },
  "board": {
    "name": "B2SS2-F",
    "vendor": "Supermicro",
    "version": "1.01",
    "serial": "ZD16BS000265",
    "assettag": "To be filled by O.E.M."
  },
  "chassis": {
    "type": 1,
    "vendor": "Supermicro",
    "version": "0123456789",
    "serial": "0123456789",
    "assettag": "To be filled by O.E.M."
  },
  "bios": {
    "vendor": "American Megatrends Inc.",
    "version": "1.0c",
    "date": "04/29/2016"
  },
  "cpu": {
    "vendor": "GenuineIntel",
    "model": "Intel(R) Xeon(R) CPU E3-1240 v5 @ 3.50GHz",
    "speed": 3500,
    "cache": 8192,
    "cpus": 1,
    "cores": 4,
    "threads": 8
  },
  "memory": {
    "type": "DDR4",
    "speed": 2400,
    "size": 65536
  },
  "storage": [
    {
      "name": "sda",
      "driver": "sd",
      "vendor": "ATA",
      "model": "SanDisk SD8SB8U2",
      "serial": "163047802208",
      "size": 256
    },
    {
      "name": "sdb",
      "driver": "sd",
      "vendor": "ATA",
      "model": "SanDisk SD8SB8U1",
      "serial": "164103801795",
      "size": 1024
    }
  ],
  "network": [
    {
      "name": "eth0",
      "driver": "igb",
      "macaddress": "0c:c4:7a:98:ba:8a",
      "port": "fibre",
      "speed": 1000
    },
    {
      "name": "eth1",
      "driver": "igb",
      "macaddress": "0c:c4:7a:98:ba:8b",
      "port": "fibre",
      "speed": 1000
    }
  ]
}
```

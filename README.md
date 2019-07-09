# morpheus-fling

A small, command line based binary for aggregating useful statistics about large environments prior to deploying Morpheus.

### Port Scanning

`morpheus-fling` reads ips and associated ports and performs a scan of these ports to inspect openness. 

### OS Statistics

`morpheus-fling` aggregates statistics about the OS it is installed on utilizing Linux kernel syscalls including memory, CPU, and available disk.

## Usage

Place binary on filesystem and give posix permissions to be executed.  Binary expects a file to exist called `network.txt`, however user can specify a separate infile.  Similarly, binary will by default write output to a file called `output.txt` but user can flag a separate outfile to be created and appended to.

Run:
```
$ ./morpheus-fling
```

Or run:
```
$ ./morpheus-fling --infile=foo.txt --outfile=bar.txt
```

By default `morpheus-fling` only allows `1024` semaphores to be used based on standard ulimit expectations in environments.  This can also be adjusted by passing a `ulimit` flag.  To use an accurate representation from the system `morpheus-fling` is deployed to run `ulimit -n` on your system and qualify your binary run with this value flagged.
Example:
```
slimshady@morpheus1:~# ulimit -n
204800
slimshady@morpheus1:~# ./morpheus-fling --ulimit=204800
```

### Inputs

As described, `morpheus-fling` defaults to looking for an input file called `network.txt`.  Format for the entries in this file should follow `ip:port` notation as below.
```
10.30.21.100:10092
10.30.21.100:3306
10.30.21.193:22
10.30.21.100:15672
10.30.21.100:5672
```

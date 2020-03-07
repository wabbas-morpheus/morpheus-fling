# morpheus-fling

A small, command line based binary for aggregating useful statistics about large environments prior to deploying Morpheus.

## Functions

1. Port Scanning
2. OS Statistics
3. Elasticsearch Health
4. Elasticsearch Indices
5. Bundling and encrypting of output files and log file

### Port Scanning

`morpheus-fling` optionally reads ips from a file and associated ports and performs a scan of these ports to inspect openness.

### OS Statistics

`morpheus-fling` aggregates statistics about the OS it is installed on utilizing Linux kernel syscalls including memory, CPU, and available disk.

### Elasticsearch Health

`morpheus-fling` queries the localhost provided Elasticsearch REST API to get a health output.

### Elasticsearch Indices

`morpheus-fling` queries the localhost provided Elasticsearch REST API to get a breakdown of all indices and their respective health.

### Bundling

By default `morpheus-fling` will read your `current` log file for your running `morpheus-ui` service and place it as a json value in the
master file.  It encrypts the content and places a key into a .zip file with the output bundle.

## Usage

Download the binary directly to your server.

```bash
wget https://github.com/gomorpheus/morpheus-fling/releases/download/v2.1.4/morpheus-fling
```

Give posix permissions to be executed.  Binary allows the specification of an `-infile` for port scanning.  Binary will by default write output to a file called `output.txt` but user can flag a separate outfile to be created and appended to.

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

### Usage

```text
morpheus-fling [options]
Options:
-infile     The source file for network port scanning.  If none is provided port scans will be skipped.
-secfile    The morpheus secrets file.  Defaults to "/etc/morpheus/morpheus-secrets.json".
-outfile    The destination directory of the generated package, "output.txt" by default.
-ulimit     Ulimit of the system, defaults to 1024.
-logfile    Logfile to add to the bundle.  Defaults to "/var/log/morpheus/morpheus-ui/current".
-bundler    Path and file to bundle into.  Defaults to "/tmp/bundler.zip".
-keyfile    Path and file to put the public key encrypted AES-GCM key into.  Defaults to "/tmp/bundlerkey.enc"

-help    Prints this text.
Examples:
Generates a bundle with port scans, system stats, elasticsearch results and morpheus logs
   $ ./morpheus-fling -infile="/home/slimshady/network.txt"

Generates a bundle with no portscans in it at /tmp/bundler.zip
   $ ./morpheus-fling
```

### Other Options

`morpheus-fling` also allows arguments to be pass for specific bundling of a logfile and specific naming of the archive.  Without specification the default is to look for `/var/log/morpheus-morpheus-ui/current` and bundle that with the outfile into `/tmp/bundler.zip`.

If you like you can do this, however

```bash
slimshady@morpheus1:~# ./morpheus-fling --logfile=/path/to/file --bundler=/path/to/archive_name.zip
```

### Inputs

As described, `morpheus-fling` allows an argument for a file of `ip:port` to be scanned for openness.  Format for the entries in this file should follow `ip:port` notation as below.

```text
10.30.21.100:10092
10.30.21.100:3306
10.30.21.193:22
10.30.21.100:15672
10.30.21.100:5672
```

### Report

`morpheus-fling` will generate a report during the run. The default is a file in the same directory as the binary called `output.txt` but this can be adjusted by making use of a flag during the run like `--outfile=/home/slimshady/foobar.txt`

The contents of the output file is encrypted, while the plaintext is displayed via standard out.  The encrypted output file and the 
public key encrypted AES256-GCM key are then bundled by default in `/tmp/bundler.zip`

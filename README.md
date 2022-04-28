# Intro

Simple multi-threaded cli ICMP ping swip tool for ad-hoc host availability checking.

# Config file

Configuration file should be provided in JSON format with simple structure:

```json
{
    "targets": [
        "<host1.fqdn>",
        "<host2.ip>",
        "google.com",
        "8.8.8.8"
        ],
    "params": {
        "numofpackets": 4
    }
}
```
Targets sections is self-descriptive.
In params section has only ## of packets that should send against of each target.


# Usage
Should run with admin priviliges via sudo or admin console.

Usage of pswip:
  
  -c string
        application config file. (default "pswip.conf")
  
  -l string
        log file. (default "pswip.log")

## Config and usage example 

```json
{
    "targets": [
        "8.8.8.8",
        "google.com"
        ],
    "params": {
        "numofpackets": 2
    }
}
```

```sh
HOST% sudo ./pswip -c example.json
Password:
Target: 8.8.8.8
 2 packets transmitted, 2 packets received, 0% packet loss
 round-trip min/avg/max/stddev = 23.36ms/23.532ms/23.704ms/172µs
================================================================================
Target: 142.250.74.142
 2 packets transmitted, 2 packets received, 0% packet loss
 round-trip min/avg/max/stddev = 20.499ms/20.8915ms/21.284ms/392.5µs
================================================================================
```


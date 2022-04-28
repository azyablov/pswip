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



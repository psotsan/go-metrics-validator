[![Go Report Card](https://goreportcard.com/badge/github.com/psotsan/go-metrics-validator)](https://goreportcard.com/report/github.com/psotsan/go-metrics-validator)
# Metrics Validator

Go CLI tool to validate metrics against thresholds.  

---

## Quick start

```bash
# Build
go build -o metrics-validator


# Run
./metrics-validator metrics.txt thresholds.txt
```

---

## File formats

- *Empty lines and lines starting with # are ignored.*

- *Metric names are case‑insensitive.*

- *Lines with wrong format are skipped (stderr).*



**thresolds.conf** (3 fields)
```
# name,value,limit_type  (limit_type = max or min)
cpu_usage,80,max
disk_free,20,min
```

**metrics.conf**
```
name,value,unit,timestamp
cpu_usage,75.2,%,2025-05-15T10:00:00Z
```

---

## Example output
```
cpu_usage = 75.20 (threshold 80.00) [WARNING]
mem_usage = 96.00 (threshold 90.00) [WARNING]
```

---

## Tests
```bash
go test -v
```

# portwatch

A lightweight CLI daemon that monitors open ports and alerts on unexpected listeners with configurable rules.

---

## Installation

```bash
go install github.com/yourname/portwatch@latest
```

Or build from source:

```bash
git clone https://github.com/yourname/portwatch.git && cd portwatch && go build -o portwatch .
```

---

## Usage

Start the daemon with a config file:

```bash
portwatch --config portwatch.yaml
```

Example `portwatch.yaml`:

```yaml
interval: 30s
alert:
  method: log
rules:
  allow:
    - port: 22
    - port: 80
    - port: 443
  deny:
    - port: 0-1024
      except: [22, 80, 443]
```

When an unexpected listener is detected, portwatch logs an alert:

```
[ALERT] Unexpected listener on port 4444 (PID 8821, process: nc)
```

Run as a background service:

```bash
portwatch --config portwatch.yaml --daemon
```

---

## Flags

| Flag | Default | Description |
|------|---------|-------------|
| `--config` | `portwatch.yaml` | Path to config file |
| `--interval` | `30s` | Poll interval |
| `--daemon` | `false` | Run as background daemon |
| `--verbose` | `false` | Enable verbose logging |

---

## License

MIT © yourname
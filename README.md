# hardware_exporter
hardware_exporter is a prometheus exporter that gathers information about installed hardware.
However, it does not export metrics for the installed hardware, for that use [node_exporter](https://github.com/prometheus/node_exporter)

## Supported information
hardware-exporter gathers information about
- cpu
- memory
- block devices
- nics
- gpu

## Usage
hardware_exporter doesn't expose information directly as prometheus metrics but only in json format,
therefore you also need [json_exporter](https://github.com/prometheus-community/json_exporter) to turn the json into valid metrics.

### Install the agent
Install the agent on every host using this script (install the binary and a systemd service)

```shell
curl https://raw.githubusercontent.com/danielr1996/hardware_exporter/refs/heads/main/deploy/systemd/install.sh | bash
```

## Development
```shell
cd agent && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o dist/hardware_exporter .

```
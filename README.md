# Piab = Prometheus in a box

<br>
## Introduction

This repo contains a dockerized setup for getting started with the Prometheus monitoring tool along with a small REST api to manage the alerts and the receivers.

It basically includes everything you need to set up a monitoring/alerting/graphing infrastructure.
Just follow the usage instructions from below and you're ready to go. You only need to add more hosts that you want to monitor and then tell Prometheus about it by adding them into the *targets* section in *prometheus.yml*.

The included small API is very handy for managing the alerts and the receivers of the alerts. It's very simple and it uses a MongoDB backend to store the alerts and receivers.

<br>
## Usage

```bash
git clone https://github.com/mariusmilea/piab.git

make up
```

<br>
## Web GUI

```
Prometheus: http://localhost:9090

Alertmanager: http://localhost:9093
```

<br>
## Methods

**GET an alert**
```bash
curl http://localhost:12345/v1/alerts
```
**GET a receiver**
```bash
curl http://localhost:12345/v1/receivers
```
**POST an alert**
```bash
curl -XPOST http://localhost:12345/v1/alerts \
-d \
'{
    "name": "alert_load1",
    "expression": "node_load1 > 1",
    "duration": "5m",
    "label": {"team": "\"admins\""},
    "summary": "Instance has a high load.",
    "description": "High 1 minute load",
    "runbook": "https://confluence/wiki/alerts"
}'
```
**POST a receiver**
```bash
curl -XPOST http://localhost:12345/v1/receivers \
-d \
'{
	"email": "joe@company.com",
	"label": {
	"team": "\"admins\""
  }
}' 
```
**Reload Alertmanager after adding an alert**
```
curl -XPOST http://localhost:12345/v1/alerts/generate
```
**Reload Prometheus after adding a receiver**
```
curl -XPOST http://localhost:12345/v1/receivers/generate
```

<br>
## Note

The localhost hostname can differ depending on the OS where this will be run from.
This address can be determined by running for example:

```bash
export | grep DOCKER_HOST | cut -f3 -d/ | cut -f1 -d:
```

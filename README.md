# prometheus-mylar3-exporter

<p align="center">
<img src="https://raw.githubusercontent.com/esanchezm/prometheus-mylar3-exporter/master/logo.png" height="230">
</p>

A [Mylar3](https://github.com/mylar3/mylar3) prometheus exporter written in Go.

## Pre-requisites

This tool requires Mylar3 API to be enabled and an API key to access it. You can do it under `Settings` -> `Web interface` and set `Enable API` option. Then generate an API key to be used for the exporter.

## Metrics

These are the metrics this program exports:


| Metric name                                         | Type     | Description      |
| --------------------------------------------------- | -------- | ---------------- |
| `mylar3_up`                                    | gauge    | Whether the Mylar3 server is answering requests from this exporter. A `version` label with the server version is added. |

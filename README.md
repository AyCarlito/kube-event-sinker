# kube-event-sinker

kube-event-sinker watches Kubernetes events and pushes them to a specified sink.

## Prerequisites

Install the following:

1. [Go](https://go.dev/dl/)
2. [Docker](https://docs.docker.com/engine/install/)
3. [Helm](https://helm.sh/docs/intro/install/)

## Install

- Build the binary:

```shell
git clone https://github.com/AyCarlito/kube-event-sinker.git
cd kube-event-sinker
make build
```

- Move the binary to the desired location e.g:

```shell
mv bin/kube-event-sinker /usr/local/bin/
```

- Alternatively, the application is containerised through the `Dockerfile` at the root of the repository, which can
be built and run through:

```shell
make docker-build docker-run
```

- The application is packaged as a Helm chart which can built and installed through:

```shell
make helm helm-install
```

## Usage

- `kube-event-sinker` is a [Cobra](https://github.com/spf13/cobra) CLI application built on a structure of commands,
arguments & flags:

```shell
./bin/kube-event-sinker --help
Watch and push Kubernetes events to a sink.

Usage:
  kube-event-sinker [flags]

Flags:
  -h, --help                          help for kube-event-sinker
      --kubeconfig string             Path to a kubeconfig file.
      --metrics-bind-address string   The address the metric endpoint binds to. (default ":9111")
      --sink string                   Sink that events should be pushed to. (default "null")
```

## Examples

- Watch Kubernetes events across all namespaces and log them.

```shell
./bin/kube-event-sinker --sink zap
{"level":"info","time":"2025-06-11T15:10:33.007+0100","caller":"sinker/sinker.go:84","msg":"Starting sinker"}
{"level":"info","time":"2025-06-11T15:10:33.007+0100","caller":"sinker/sinker.go:88","msg":"Waiting for cache to sync"}
{"level":"info","time":"2025-06-11T15:10:33.107+0100","caller":"sinker/sinker.go:92","msg":"Cache synced"}
{"level":"info","time":"2025-06-11T15:10:37.877+0100","caller":"sinks/zap.go:42","msg":"Handling","kind":"Deployment","name":"busybox","namespace":"kube-event-sinker","reason":"ScalingReplicaSet","type":"Normal"}
{"level":"info","time":"2025-06-11T15:10:37.892+0100","caller":"sinks/zap.go:42","msg":"Handling","kind":"ReplicaSet","name":"busybox-69fdfd884c","namespace":"kube-event-sinker","reason":"SuccessfulCreate","type":"Normal"}
{"level":"info","time":"2025-06-11T15:10:37.897+0100","caller":"sinks/zap.go:42","msg":"Handling","kind":"Pod","name":"busybox-69fdfd884c-hhzt9","namespace":"kube-event-sinker","reason":"Scheduled","type":"Normal"}
{"level":"info","time":"2025-06-11T15:10:38.364+0100","caller":"sinks/zap.go:42","msg":"Handling","kind":"Pod","name":"busybox-69fdfd884c-hhzt9","namespace":"kube-event-sinker","reason":"Pulling","type":"Normal"}
{"level":"info","time":"2025-06-11T15:10:39.093+0100","caller":"sinks/zap.go:42","msg":"Handling","kind":"Pod","name":"busybox-69fdfd884c-hhzt9","namespace":"kube-event-sinker","reason":"Pulled","type":"Normal"}
{"level":"info","time":"2025-06-11T15:10:39.191+0100","caller":"sinks/zap.go:42","msg":"Handling","kind":"Pod","name":"busybox-69fdfd884c-hhzt9","namespace":"kube-event-sinker","reason":"Created","type":"Normal"}
{"level":"info","time":"2025-06-11T15:10:39.280+0100","caller":"sinks/zap.go:42","msg":"Handling","kind":"Pod","name":"busybox-69fdfd884c-hhzt9","namespace":"kube-event-sinker","reason":"Started","type":"Normal"}
```

- Prometheus metrics are generated for all sinks, including the above zap sink.
- If metrics are the only information of interest, the events can be pushed to a null sink.

```shell
 ./bin/kube-event-sinker
{"level":"info","time":"2025-06-11T15:27:13.890+0100","caller":"sinker/sinker.go:84","msg":"Starting sinker"}
{"level":"info","time":"2025-06-11T15:27:13.890+0100","caller":"sinker/sinker.go:88","msg":"Waiting for cache to sync"}
{"level":"info","time":"2025-06-11T15:27:13.991+0100","caller":"sinker/sinker.go:92","msg":"Cache synced"}

curl "http://127.0.0.1:9111/metrics" -s |  grep kube
# HELP kube_event_sinker_events_total Total of Kubernetes events handled.
# TYPE kube_event_sinker_events_total counter
kube_event_sinker_events_total{kind="Pod",reason="BackOff",type="Warning"} 1
kube_event_sinker_events_total{kind="Pod",reason="Created",type="Normal"} 1
kube_event_sinker_events_total{kind="Pod",reason="Pulled",type="Normal"} 1
kube_event_sinker_events_total{kind="Pod",reason="Pulling",type="Normal"} 1
kube_event_sinker_events_total{kind="Pod",reason="Scheduled",type="Normal"} 1
kube_event_sinker_events_total{kind="Pod",reason="Started",type="Normal"} 1
kube_event_sinker_events_total{kind="ReplicaSet",reason="SuccessfulCreate",type="Normal"} 1
```

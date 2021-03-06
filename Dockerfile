FROM registry.svc.ci.openshift.org/openshift/release:golang-1.11 AS builder

WORKDIR /go/src/github.com/prometheus/node_exporter
COPY . .
RUN make build

FROM  registry.svc.ci.openshift.org/openshift/origin-v4.0:base
LABEL io.k8s.display-name="OpenShift Prometheus Node Exporter" \
      io.k8s.description="Prometheus exporter for machine metrics" \
      io.openshift.tags="prometheus,monitoring" \
      maintainer="The Prometheus Authors <prometheus-developers@googlegroups.com>"

COPY --from=builder /go/src/github.com/prometheus/node_exporter/node_exporter /bin/node_exporter

EXPOSE      9100
USER        nobody
ENTRYPOINT  [ "/bin/node_exporter" ]

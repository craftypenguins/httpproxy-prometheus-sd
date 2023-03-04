# httpproxy-prometheus-sd

K8S HttpProxies (Concourse) Prometheus Scrape Targets Exporter

Simple tool query the K8S API for HttpProxies (Concourse) objects and expose them as a Prometheus scrape target.

For each one that is annotated with blackbox-monitor: "true" it will
generate as a JSON response to a query a Prometheus [http_sd_config|https://prometheus.io/docs/prometheus/latest/configuration/configuration/#http_sd_config] formatted list of scrape targets to use with Blackbox Exporter

## RBAC

See the /k8s folder.  A non default service account and cluster role needs to be setup
so that the program can collect all httpproxy instances over all namespaces.

## Tagging

Tag an httpproxy like this or add to it's manifest

```bash
kubectl annotate httpproxy <target_name> blackbox-monitor="true"
```

## http / https

If a certificate is assigned in the httpproxy, then the url will be rendered with `https://`, `http://` otherwise.

## url paths

Right now, route paths are not considered.  Only the virtualhost name.
This is a todo if needed.

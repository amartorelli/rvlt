#!/bin/bash

# Pre-requisite: helm plugin install https://github.com/ContainerSolutions/helm-monitor

helm upgrade -i -f values.yaml rvlt-charts/helloworld --version $HELLOWORLD_VER helloworld

# Automatic roll back if 5xx are returned

helm monitor prometheus --prometheus=http://prometheus.rvlt.com helloworld 'rate(helloworld_http_requests_total{code=~"^5.*$"}[5m]) > 0'

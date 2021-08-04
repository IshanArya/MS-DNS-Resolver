#!/bin/bash

GOOS=linux go build -o ./app cmd/coredns_diagnostics/main.go

az acr build --image coredns_diagnostics:v0.0.$1 \
  --registry ishanRegistry \
  --file Dockerfile .
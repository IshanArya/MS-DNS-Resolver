#!/bin/bash

GOOS=linux go build -o ./app cmd/dns_resolver/main.go

az acr build --image dns_resolver:v0.0.$1 \
  --registry ishanRegistry \
  --file Dockerfile .
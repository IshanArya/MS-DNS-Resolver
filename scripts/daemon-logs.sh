#!/bin/bash

kubectl logs -l name=dnsresolver-spread -n kube-system --prefix --max-log-requests=10 --tail=20
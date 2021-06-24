# MS-DNS-Resolver

## Hints

See logs:
```shell
kubectl logs $(kubectl get pods -n kube-system --selector=job-name=dnsresolver --output=jsonpath='{.items[*].metadata.name}') -n kube-system
```
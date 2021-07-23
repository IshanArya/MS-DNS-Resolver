# MS-DNS-Resolver

## Hints

See logs:
```shell
kubectl logs $(kubectl get pods -A --selector=job-name=dnsresolver --output=jsonpath='{.items[*].metadata.name}') --all-namespaces
k logs $(k get pods -A -l name=dnsresolver-spread --output=jsonpath='{.items[*].metadata.name}') -n kubes-system
```

```shell
az acr build --image dns_resolver:v0.0.x \
  --registry ishanRegistry \
  --file Dockerfile . 
```
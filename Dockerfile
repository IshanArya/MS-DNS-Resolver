FROM golang:1.16.5

WORKDIR /app

COPY . .

RUN apt-get update && apt-get install -y dnsutils
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build ./...

ENTRYPOINT /app/dns_resolver -path ./configs/kube.txt


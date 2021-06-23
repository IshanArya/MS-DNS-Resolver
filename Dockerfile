FROM golang:1.16.5

WORKDIR /app

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go install -v ./...

ENTRYPOINT app/dns_resolver


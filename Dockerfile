# syntax=docker/dockerfile:1

FROM golang:1.19 as builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY cmd ./
COPY pkg ./

RUN cd cmd/exporter && \
    CGO_ENABLED=0 GOOS=linux go build -o /prometheus-nginx-exporter

FROM golang:1.19 as worker

COPY --from=builder /prometheus-nginx-exporter /app/prometheus-nginx-exporter

CMD ["/app/prometheus-nginx-exporter"]
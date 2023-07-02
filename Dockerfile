# syntax=docker/dockerfile:1

FROM golang:1.19 as builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY cmd ./
COPY pkg ./


RUN CGO_ENABLED=0 GOOS=linux go build ./cmd/exporter -o /prometheus-nginx-exporter

FROM golang:1.19 as worker

COPY --from=builder /prometheus-nginx-exporter /app/prometheus-nginx-exporter

CMD ["/app/prometheus-nginx-exporter"]
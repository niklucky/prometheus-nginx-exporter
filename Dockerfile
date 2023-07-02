# syntax=docker/dockerfile:1

FROM golang:1.19-alpine as builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY cmd ./cmd
COPY pkg ./pkg


RUN CGO_ENABLED=0 GOOS=linux go build -o /prometheus-nginx-exporter ./cmd/exporter

FROM golang:1.19-alpine as worker

COPY --from=builder /prometheus-nginx-exporter /app/prometheus-nginx-exporter

CMD ["/app/prometheus-nginx-exporter"]
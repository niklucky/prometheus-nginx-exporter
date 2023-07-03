package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/hpcloud/tail"
	"github.com/niklucky/prometheus-nginx-exporter/pkg/exporter"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	config := newConfig()

	// Called on each collect request
	basicStats := func() ([]exporter.NginxStats, error) {
		var httpClient = &http.Client{
			Timeout: time.Second * 10,
		}
		resp, err := httpClient.Get(config.nginx.uri)
		if err != nil {
			fmt.Printf("request to basic_stats failed: %s: %s", config.nginx.uri, err)
			return []exporter.NginxStats{}, nil
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Printf("body read fail: %s\n", err)
			return []exporter.NginxStats{}, nil
		}
		r := bytes.NewReader(body)

		return exporter.ScanBasicStats(r)
	}

	bc := exporter.NewBasicCollector(basicStats)
	reg := prometheus.NewRegistry()
	// reg.MustRegister(collectors.NewGoCollector())
	reg.MustRegister(bc)

	m := exporter.NewMetrics(reg)
	go tailAccessLogFile(m, config.nginx.accessLogPath)

	mux := http.NewServeMux()
	promHandler := promhttp.HandlerFor(reg, promhttp.HandlerOpts{})
	mux.Handle("/metrics", promHandler)

	// Start listening for HTTP connections
	port := fmt.Sprintf(":%d", config.promPort)
	log.Printf("Starting nginx exporter on %s/metrics", port)
	if err := http.ListenAndServe(port, mux); err != nil {
		log.Fatalf("Error starting nginx exporter: %s", err)
	}
}

type nginxAccessLog struct {
	Request              string  `json:"request"`
	RequestTime          float64 `json:"request_time"`
	RequestMethod        string  `json:"request_method"`
	RequestUri           string  `json:"request_uri"`
	Connection           string  `json:"connection"`
	LogTime              float64 `json:"time"`
	RemoteAddr           string  `json:"remote_addr"`
	BodyBytesSent        float64 `json:"body_bytes_sent"`
	HttpReferer          string  `json:"http_referer"`
	StatusCode           int     `json:"status"`
	UserAgent            string  `json:"user_agent"`
	UpstreamAddr         string  `json:"upstream_addr"`
	UpstreamStatus       int     `json:"upstream_status"`
	UpstreamResponseTime float64 `json:"upstream_response_time"`
	UpstreamConnectTime  float64 `json:"upstream_connect_time"`
	UpstreamHeaderTime   float64 `json:"upstream_header_time"`
}

func toAccessLog(accessLogRequest []byte) (*nginxAccessLog, error) {

	const substr = `{"time":`
	start := strings.Index(string(accessLogRequest), substr)
	if start < 0 {
		msg := fmt.Sprintf("failed to find access-log request JSON '%s' starting with '%s'", string(accessLogRequest), substr)
		fmt.Println(msg)
		return nil, errors.New(msg)
	}
	var ret nginxAccessLog
	err := json.Unmarshal(bytes.Trim([]byte(string(accessLogRequest)[start:]), "\x00"), &ret)
	if err != nil {
		log.Printf("failed to unmarshal access-log '%s' with '%v'", string(accessLogRequest)[start:], err)
		return nil, err
	}

	return &ret, nil
}

func tailAccessLogFile(m *exporter.Metrics, path string) {
	t, err := tail.TailFile(path, tail.Config{Follow: true, ReOpen: true})
	if err != nil {
		log.Fatalf("tail.TailFile failed: %s", err)
	}

	for line := range t.Lines {
		accessLog, err := toAccessLog([]byte(line.Text))
		if err != nil {
			continue
		}

		m.Size.Add(accessLog.BodyBytesSent)

		m.Requests.With(prometheus.Labels{
			"method":      accessLog.RequestMethod,
			"status_code": strconv.Itoa(accessLog.StatusCode),
			"path":        accessLog.RequestUri,
		}).Add(1)

		m.Duration.With(prometheus.Labels{
			"method":      accessLog.RequestMethod,
			"status_code": strconv.Itoa(accessLog.StatusCode),
			"path":        accessLog.RequestUri,
		}).Observe(accessLog.RequestTime)

	}

}

package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type config struct {
	promPort int
	nginx    nginxConfig
}

type nginxConfig struct {
	uri           string
	accessLogPath string
	allowedPaths  []string
}

func newConfig() *config {
	var (
		targetHost    = flag.String("target.host", "localhost", "nginx address with basic_status")
		targetPort    = flag.Int("target.port", 8080, "nginx port with basic_status")
		targetPath    = flag.String("target.path", "/status", "nginx path  with basic_status")
		promPort      = flag.Int("prom port", 9150, "port to expose promitheus metrics")
		accessLogPath = flag.String("target.log", "/var/log/nginx/access.log", "path to nginx logs")
		allowedPaths  = flag.String("allowed.paths", "/", "allowed HTTP path for metrics. Other will be ignored")
	)

	flag.Parse()

	nginxUri := fmt.Sprintf("http://%s:%d%s", *targetHost, *targetPort, *targetPath)
	nginxAccessLogPath := *accessLogPath

	if os.Getenv("NGINX_URI") != "" {
		nginxUri = os.Getenv("NGINX_URI")
	}
	if os.Getenv("NGINX_ACCESS_LOG_PATH") != "" {
		nginxAccessLogPath = os.Getenv("NGINX_ACCESS_LOG_PATH")
	}
	if os.Getenv("PROM_PORT") != "" {
		if p, e := strconv.Atoi(os.Getenv("PROM_PORT")); e == nil {
			promPort = &p
			log.Printf("error parsing PROM_PORT: %s", e)
		}
	}
	var nginxAllowedPaths []string
	if allowedPaths != nil {
		nginxAllowedPaths = strings.Split(*allowedPaths, ",")
	}
	if os.Getenv("ALLOWED_PATHS") != "" {
		paths := os.Getenv("ALLOWED_PATHS")
		nginxAllowedPaths = strings.Split(paths, ",")
	}

	return &config{
		nginx: nginxConfig{
			uri:           nginxUri,
			accessLogPath: nginxAccessLogPath,
			allowedPaths:  nginxAllowedPaths,
		},
		promPort: *promPort,
	}
}

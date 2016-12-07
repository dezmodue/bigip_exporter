package main

import (
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/ExpressenAB/bigip_exporter/collector"
	"github.com/ExpressenAB/bigip_exporter/config"
	"github.com/pr8kerl/f5er/f5"
	"github.com/prometheus/client_golang/prometheus"
)

func listen(exporter_bind_address string, exporter_bind_port int) {
	http.Handle("/metrics", prometheus.Handler())
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
			<head><title>BIG-IP Exporter</title></head>
			<body>
			<h1>BIG-IP Exporter</h1>
			<p><a href="/metrics">Metrics</a></p>
			</body>
			</html>`))
	})
	exporter_bind := exporter_bind_address + ":" + strconv.Itoa(exporter_bind_port)
	log.Fatal(http.ListenAndServe(exporter_bind, nil))
}

func main() {
	config := config.GetConfig()

	if config.Exporter_config.Exporter_debug {
		log.Printf("Config: %v", config)
	}

	bigip_endpoint := config.Bigip_config.Bigip_host + ":" + strconv.Itoa(config.Bigip_config.Bigip_port)
	var exporter_partitions_list []string
	if config.Exporter_config.Exporter_partitions != "" {
		exporter_partitions_list = strings.Split(config.Exporter_config.Exporter_partitions, ",")
	} else {
		exporter_partitions_list = nil
	}
	auth_method := f5.TOKEN
	if config.Bigip_config.Bigip_basic_auth {
		auth_method = f5.BASIC_AUTH
	}

	bigip := f5.New(bigip_endpoint, config.Bigip_config.Bigip_username, config.Bigip_config.Bigip_password, auth_method)

	_, bigipCollector := collector.NewBigIpCollector(bigip, config.Exporter_config.Exporter_namespace, exporter_partitions_list)

	prometheus.MustRegister(bigipCollector)
	listen(config.Exporter_config.Exporter_bind_address, config.Exporter_config.Exporter_bind_port)
}

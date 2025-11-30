package main

import (
	"log"
	"net"
	"net/http"

	"github.com/danielr1996/hardware_exporter/internal/collector"
	"github.com/danielr1996/hardware_exporter/internal/util"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func collectAll() collector.Inventory {
	uuid := util.DeriveUUID()

	return collector.Inventory{
		UUID:      uuid,
		System:    collector.CollectSystemInfo(),
		Baseboard: collector.CollectBaseboardInfo(),
		Bios:      collector.CollectBiosInfo(),
		CPU:       collector.CollectCPUInfo(),
		Memory:    collector.CollectMemory(),
		Disks:     collector.CollectBlockDevices(),
		Nics:      collector.CollectNics(),
		PCI:       collector.CollectPCIDevices(),
	}
}

func main() {
	inv := collectAll()
	reg := prometheus.NewRegistry()
	reg.MustRegister(collector.NewInventoryCollector(inv))

	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{}))

	// Try IPv4 listener
	ipv4, err4 := net.Listen("tcp4", "0.0.0.0:9105")
	if err4 != nil {
		log.Printf("IPv4 unavailable: %v", err4)
	} else {
		log.Println("Listening on 0.0.0.0:9105 (IPv4)")
		go func() {
			if err := http.Serve(ipv4, mux); err != nil {
				log.Fatalf("IPv4 server error: %v", err)
			}
		}()
	}

	// Try IPv6 listener
	ipv6, err6 := net.Listen("tcp6", "[::]:9105")
	if err6 != nil {
		log.Printf("IPv6 unavailable: %v", err6)
	} else {
		log.Println("Listening on [::]:9105 (IPv6)")
		if err := http.Serve(ipv6, mux); err != nil {
			log.Fatalf("IPv6 server error: %v", err)
		}
	}
}

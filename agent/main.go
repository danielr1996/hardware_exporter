package main

import (
	"log"
	"net"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Inventory struct {
	UUID      string
	System    SystemInfo
	Baseboard BaseboardInfo
	Bios      BiosInfo
	CPU       CPUInfo
	Memory    []MemoryModule
	Disks     []BlockDevice
	Nics      []NicInfo
	PCI       []PCIDevice
}

func collectAll() Inventory {
	uuid := deriveUUID()

	return Inventory{
		UUID:      uuid,
		System:    collectSystemInfo(),
		Baseboard: collectBaseboardInfo(),
		Bios:      collectBiosInfo(),
		CPU:       collectCPUInfo(),
		Memory:    collectMemory(),
		Disks:     collectBlockDevices(),
		Nics:      collectNics(),
		PCI:       collectPCIDevices(),
	}
}

func main() {
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())

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

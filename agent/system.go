package main

import (
	"os"
)

type SystemInfo struct {
	Hostname      string
	OS            string
	Kernel        string
	Chassis       string
	ProductVendor string
	ProductName   string
	ProductFamily string
	ProductSerial string
}

func collectSystemInfo() SystemInfo {
	hostname, _ := os.Hostname()

	return SystemInfo{
		Hostname:      hostname,
		OS:            readOSReleaseField("PRETTY_NAME"),
		Kernel:        readFirstLine("/proc/sys/kernel/osrelease"),
		Chassis:       readFirstLine("/sys/class/dmi/id/chassis_type"),
		ProductVendor: readFirstLine("/sys/class/dmi/id/sys_vendor"),
		ProductName:   readFirstLine("/sys/class/dmi/id/product_name"),
		ProductFamily: readFirstLine("/sys/class/dmi/id/product_family"),
		ProductSerial: readFirstLine("/sys/class/dmi/id/product_serial"),
	}
}

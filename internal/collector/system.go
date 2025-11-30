package collector

import (
	"os"

	"github.com/danielr1996/hardware_exporter/internal/util"
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

func CollectSystemInfo() SystemInfo {
	hostname, _ := os.Hostname()

	return SystemInfo{
		Hostname:      hostname,
		OS:            util.ReadOSReleaseField("PRETTY_NAME"),
		Kernel:        util.ReadFirstLine("/proc/sys/kernel/osrelease"),
		Chassis:       util.ReadFirstLine("/sys/class/dmi/id/chassis_type"),
		ProductVendor: util.ReadFirstLine("/sys/class/dmi/id/sys_vendor"),
		ProductName:   util.ReadFirstLine("/sys/class/dmi/id/product_name"),
		ProductFamily: util.ReadFirstLine("/sys/class/dmi/id/product_family"),
		ProductSerial: util.ReadFirstLine("/sys/class/dmi/id/product_serial"),
	}
}

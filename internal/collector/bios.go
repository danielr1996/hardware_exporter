package collector

import "github.com/danielr1996/hardware_exporter/internal/util"

type BiosInfo struct {
	Vendor  string
	Version string
	Date    string
}

func CollectBiosInfo() BiosInfo {
	return BiosInfo{
		Vendor:  util.ReadFirstLine("/sys/class/dmi/id/bios_vendor"),
		Version: util.ReadFirstLine("/sys/class/dmi/id/bios_version"),
		Date:    util.ReadFirstLine("/sys/class/dmi/id/bios_date"),
	}
}

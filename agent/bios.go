package main

type BiosInfo struct {
	Vendor  string
	Version string
	Date    string
}

func collectBiosInfo() BiosInfo {
	return BiosInfo{
		Vendor:  readFirstLine("/sys/class/dmi/id/bios_vendor"),
		Version: readFirstLine("/sys/class/dmi/id/bios_version"),
		Date:    readFirstLine("/sys/class/dmi/id/bios_date"),
	}
}

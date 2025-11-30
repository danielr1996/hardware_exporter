package main

type BaseboardInfo struct {
	Vendor  string
	Name    string
	Version string
	Serial  string
}

func collectBaseboardInfo() BaseboardInfo {
	return BaseboardInfo{
		Vendor:  readFirstLine("/sys/class/dmi/id/board_vendor"),
		Name:    readFirstLine("/sys/class/dmi/id/board_name"),
		Version: readFirstLine("/sys/class/dmi/id/board_version"),
		Serial:  readFirstLine("/sys/class/dmi/id/board_serial"),
	}
}

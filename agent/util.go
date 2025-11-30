package main

import (
	"bufio"
	"crypto/sha1"
	"encoding/hex"
	"os"
	"path/filepath"
	"strings"
)

func readFirstLine(path string) string {
	data, err := os.ReadFile(path)
	if err != nil {
		return ""
	}
	s := strings.TrimSpace(string(data))
	if i := strings.IndexByte(s, '\n'); i != -1 {
		return s[:i]
	}
	return s
}

func readOSReleaseField(key string) string {
	f, err := os.Open("/etc/os-release")
	if err != nil {
		return ""
	}
	defer f.Close()

	sc := bufio.NewScanner(f)
	prefix := key + "="
	for sc.Scan() {
		line := sc.Text()
		if strings.HasPrefix(line, prefix) {
			val := strings.TrimPrefix(line, prefix)
			val = strings.Trim(val, `"`)
			return val
		}
	}
	return ""
}

func cleanSerial(s string) string {
	switch strings.TrimSpace(s) {
	case "", "System Serial Number", "To Be Filled By O.E.M.", "Default string", "None", "Not Specified":
		return ""
	default:
		return s
	}
}

func deriveUUID() string {
	serial := cleanSerial(readFirstLine("/sys/class/dmi/id/product_serial"))
	productUUID := cleanSerial(readFirstLine("/sys/class/dmi/id/product_uuid"))

	// lowest MAC
	lowestMAC := ""
	if ifs, err := os.ReadDir("/sys/class/net"); err == nil {
		for _, iface := range ifs {
			m := readFirstLine(filepath.Join("/sys/class/net", iface.Name(), "address"))
			if m == "" {
				continue
			}
			if lowestMAC == "" || m < lowestMAC {
				lowestMAC = m
			}
		}
	}

	// fallback
	hostID := serial
	if hostID == "" {
		hostID = productUUID
	}
	if hostID == "" {
		hostID = lowestMAC
	}
	if hostID == "" {
		hostID = "unknown"
	}

	sum := sha1.Sum([]byte(hostID))
	hash := hex.EncodeToString(sum[:])

	timeLow := hash[0:8]
	timeMid := hash[8:12]
	timeHi := hash[12:16]
	clockSeq := hash[16:20]
	node := hash[20:32]

	timeHi = "5" + timeHi[1:]     // version
	clockSeq = "a" + clockSeq[1:] // variant

	return timeLow + "-" + timeMid + "-" + timeHi + "-" + clockSeq + "-" + node
}

package main

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

type InventoryCollector struct {
	uuid      string
	system    SystemInfo
	baseboard BaseboardInfo
	bios      BiosInfo
	cpu       CPUInfo
	memory    []MemoryModule
	disks     []BlockDevice
	nics      []NicInfo
	pci       []PCIDevice
	timestamp float64

	// descriptors
	systemInfoDesc    *prometheus.Desc
	baseboardInfoDesc *prometheus.Desc
	biosInfoDesc      *prometheus.Desc

	cpuInfoDesc     *prometheus.Desc
	cpuPhysicalDesc *prometheus.Desc
	cpuLogicalDesc  *prometheus.Desc
	cpuSocketsDesc  *prometheus.Desc
	cpuMinHzDesc    *prometheus.Desc
	cpuMaxHzDesc    *prometheus.Desc

	memInfoDesc        *prometheus.Desc
	memSizeDesc        *prometheus.Desc
	memSpeedDesc       *prometheus.Desc
	memConfigSpeedDesc *prometheus.Desc
	memRankDesc        *prometheus.Desc
	memCLDesc          *prometheus.Desc

	blockInfoDesc *prometheus.Desc
	blockSizeDesc *prometheus.Desc

	nicInfoDesc  *prometheus.Desc
	nicSpeedDesc *prometheus.Desc

	pciInfoDesc *prometheus.Desc

	scrapeTSDesc *prometheus.Desc
}

func NewInventoryCollector(inv Inventory) *InventoryCollector {
	u := inv.UUID

	return &InventoryCollector{
		uuid:      u,
		system:    inv.System,
		baseboard: inv.Baseboard,
		bios:      inv.Bios,
		cpu:       inv.CPU,
		memory:    inv.Memory,
		disks:     inv.Disks,
		nics:      inv.Nics,
		pci:       inv.PCI,
		timestamp: float64(time.Now().Unix()),

		systemInfoDesc: prometheus.NewDesc(
			"inventory_system_info",
			"System information",
			[]string{"uuid", "hostname", "os", "kernel", "chassis", "product_vendor", "product_name", "product_family", "product_serial"},
			nil,
		),
		baseboardInfoDesc: prometheus.NewDesc(
			"inventory_baseboard_info",
			"Baseboard info",
			[]string{"uuid", "vendor", "name", "version", "serial"},
			nil,
		),
		biosInfoDesc: prometheus.NewDesc(
			"inventory_bios_info",
			"Bios info",
			[]string{"uuid", "vendor", "version", "date"},
			nil,
		),

		cpuInfoDesc: prometheus.NewDesc(
			"inventory_cpu_info",
			"CPU info",
			[]string{"uuid", "vendor", "model", "arch"},
			nil,
		),
		cpuPhysicalDesc: prometheus.NewDesc("inventory_cpu_physical_cores", "Physical CPU cores", []string{"uuid"}, nil),
		cpuLogicalDesc:  prometheus.NewDesc("inventory_cpu_logical_cores", "Logical CPU cores", []string{"uuid"}, nil),
		cpuSocketsDesc:  prometheus.NewDesc("inventory_cpu_sockets", "CPU sockets", []string{"uuid"}, nil),
		cpuMinHzDesc:    prometheus.NewDesc("inventory_cpu_min_frequency_hz", "CPU minimum frequency", []string{"uuid"}, nil),
		cpuMaxHzDesc:    prometheus.NewDesc("inventory_cpu_max_frequency_hz", "CPU maximum frequency", []string{"uuid"}, nil),

		memInfoDesc: prometheus.NewDesc(
			"inventory_memory_device_info",
			"Memory module info",
			[]string{"uuid", "locator", "bank", "channel", "type", "form", "manufacturer", "serial", "part", "rank", "cl", "speed", "configured_speed"},
			nil,
		),
		memSizeDesc:        prometheus.NewDesc("inventory_memory_device_size_bytes", "Memory size in bytes", []string{"uuid", "locator"}, nil),
		memSpeedDesc:       prometheus.NewDesc("inventory_memory_device_speed_hz", "Memory speed (Hz)", []string{"uuid", "locator"}, nil),
		memConfigSpeedDesc: prometheus.NewDesc("inventory_memory_device_configured_speed_hz", "Configured memory speed (Hz)", []string{"uuid", "locator"}, nil),
		memRankDesc:        prometheus.NewDesc("inventory_memory_device_ranks", "Memory rank count", []string{"uuid", "locator"}, nil),
		memCLDesc:          prometheus.NewDesc("inventory_memory_device_cl_cycles", "CL latency (cycles)", []string{"uuid", "locator"}, nil),

		blockInfoDesc: prometheus.NewDesc(
			"inventory_block_device_info",
			"Block device info",
			[]string{"uuid", "device", "rotational"},
			nil,
		),
		blockSizeDesc: prometheus.NewDesc(
			"inventory_block_device_size_bytes",
			"Block device size in bytes",
			[]string{"uuid", "device"},
			nil,
		),

		nicInfoDesc: prometheus.NewDesc(
			"inventory_nic_info",
			"NIC info",
			[]string{"uuid", "ifname", "mac"},
			nil,
		),
		nicSpeedDesc: prometheus.NewDesc(
			"inventory_nic_speed_bits_per_second",
			"NIC link speed",
			[]string{"uuid", "ifname"},
			nil,
		),

		pciInfoDesc: prometheus.NewDesc(
			"inventory_pci_device_info",
			"PCI/PCIe device info",
			[]string{"uuid", "address", "vendor_id", "device_id", "class", "driver", "name"},
			nil,
		),

		scrapeTSDesc: prometheus.NewDesc(
			"inventory_scrape_timestamp_seconds",
			"Exporter timestamp",
			[]string{"uuid"},
			nil,
		),
	}
}

func (c *InventoryCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.systemInfoDesc
	ch <- c.baseboardInfoDesc
	ch <- c.biosInfoDesc

	ch <- c.cpuInfoDesc
	ch <- c.cpuPhysicalDesc
	ch <- c.cpuLogicalDesc
	ch <- c.cpuSocketsDesc
	ch <- c.cpuMinHzDesc
	ch <- c.cpuMaxHzDesc

	ch <- c.memInfoDesc
	ch <- c.memSizeDesc
	ch <- c.memSpeedDesc
	ch <- c.memConfigSpeedDesc
	ch <- c.memRankDesc
	ch <- c.memCLDesc

	ch <- c.blockInfoDesc
	ch <- c.blockSizeDesc

	ch <- c.nicInfoDesc
	ch <- c.nicSpeedDesc

	ch <- c.pciInfoDesc

	ch <- c.scrapeTSDesc
}

func (c *InventoryCollector) Collect(ch chan<- prometheus.Metric) {
	u := c.uuid

	// system
	s := c.system
	ch <- prometheus.MustNewConstMetric(
		c.systemInfoDesc,
		prometheus.GaugeValue,
		1,
		u, s.Hostname, s.OS, s.Kernel, s.Chassis, s.ProductVendor, s.ProductName, s.ProductFamily, s.ProductSerial,
	)

	// baseboard
	b := c.baseboard
	ch <- prometheus.MustNewConstMetric(
		c.baseboardInfoDesc,
		prometheus.GaugeValue,
		1,
		u, b.Vendor, b.Name, b.Version, b.Serial,
	)

	// bios
	bi := c.bios
	ch <- prometheus.MustNewConstMetric(
		c.biosInfoDesc,
		prometheus.GaugeValue,
		1,
		u, bi.Vendor, bi.Version, bi.Date,
	)

	// cpu
	cpu := c.cpu
	ch <- prometheus.MustNewConstMetric(c.cpuInfoDesc, prometheus.GaugeValue, 1, u, cpu.Vendor, cpu.ModelName, cpu.Architecture)
	ch <- prometheus.MustNewConstMetric(c.cpuPhysicalDesc, prometheus.GaugeValue, float64(cpu.PhysicalCores), u)
	ch <- prometheus.MustNewConstMetric(c.cpuLogicalDesc, prometheus.GaugeValue, float64(cpu.LogicalCores), u)
	ch <- prometheus.MustNewConstMetric(c.cpuSocketsDesc, prometheus.GaugeValue, float64(cpu.Sockets), u)
	ch <- prometheus.MustNewConstMetric(c.cpuMinHzDesc, prometheus.GaugeValue, cpu.MinHz, u)
	ch <- prometheus.MustNewConstMetric(c.cpuMaxHzDesc, prometheus.GaugeValue, cpu.MaxHz, u)

	// memory
	for _, m := range c.memory {
		ch <- prometheus.MustNewConstMetric(
			c.memInfoDesc,
			prometheus.GaugeValue,
			1,
			u, m.Locator, m.Bank, m.Channel, m.Type, m.Form, m.Manufacturer, m.Serial, m.Part, m.Rank, m.CLLatency, m.Speed, m.ConfigSpeed,
		)

		ch <- prometheus.MustNewConstMetric(c.memSizeDesc, prometheus.GaugeValue, float64(m.SizeBytes), u, m.Locator)
		ch <- prometheus.MustNewConstMetric(c.memSpeedDesc, prometheus.GaugeValue, m.SpeedHz, u, m.Locator)
		ch <- prometheus.MustNewConstMetric(c.memConfigSpeedDesc, prometheus.GaugeValue, m.ConfigSpeedHz, u, m.Locator)
		ch <- prometheus.MustNewConstMetric(c.memRankDesc, prometheus.GaugeValue, m.RankN, u, m.Locator)
		ch <- prometheus.MustNewConstMetric(c.memCLDesc, prometheus.GaugeValue, m.CLCycles, u, m.Locator)
	}

	// block devices
	for _, d := range c.disks {
		rot := "0"
		if d.Rotational {
			rot = "1"
		}
		ch <- prometheus.MustNewConstMetric(c.blockInfoDesc, prometheus.GaugeValue, 1, u, d.Name, rot)
		ch <- prometheus.MustNewConstMetric(c.blockSizeDesc, prometheus.GaugeValue, float64(d.SizeBytes), u, d.Name)
	}

	// nics
	for _, n := range c.nics {
		ch <- prometheus.MustNewConstMetric(c.nicInfoDesc, prometheus.GaugeValue, 1, u, n.Name, n.MAC)
		if n.Speed > 0 {
			ch <- prometheus.MustNewConstMetric(c.nicSpeedDesc, prometheus.GaugeValue, float64(n.Speed), u, n.Name)
		}
	}

	// pci devices
	for _, p := range c.pci {
		ch <- prometheus.MustNewConstMetric(
			c.pciInfoDesc,
			prometheus.GaugeValue,
			1,
			u, p.Address, p.VendorID, p.DeviceID, p.Class, p.Driver, p.DeviceName,
		)
	}

	// scrape timestamp
	ch <- prometheus.MustNewConstMetric(c.scrapeTSDesc, prometheus.GaugeValue, c.timestamp, u)
}

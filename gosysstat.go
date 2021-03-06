package main

import (
	"bytes"
	"fmt"
	"log"
	"os"

	"github.com/ossareh/libgosysstat/core"
	"github.com/ossareh/libgosysstat/processor/cpu"
	"github.com/ossareh/libgosysstat/processor/mem"
)

const (
	TICK_INTERVAL       = 1
	CPU_STAT_FMT        = "%s:(user:%.2f%%, sys:%.2f%%, idle:%.2f%%, io:%.2f%%)\n"
	CPU_SINGLE_STAT_FMT = "%s:%d\n"
	MEM_STAT_FMT        = "RAM:(total:%d KB, used:%.2f%%, cache:%.2f%%)\n"
	SWAP_STAT_FMT       = "SWAP:(total:%d KB, used:%.2f%%)\n"
)

func prepareCpuValues(values []uint64) (user, sys, idle, io float64) {
	var total float64
	for _, v := range values {
		total += float64(v)
	}
	user = (float64(values[0]+values[1]) / total) * 100
	sys = (float64(values[2]) / total) * 100
	idle = (float64(values[3]) / total) * 100
	io = (float64(values[4]) / total) * 100
	return
}

func formatCpuStat(data []core.Stat) string {
	var buf bytes.Buffer

	for _, d := range data {
		values := d.Values()
		var s string
		switch d.Type() {
		case cpu.TOTAL:
			user, sys, idle, io := prepareCpuValues(values)
			s = fmt.Sprintf(CPU_STAT_FMT, "Total", user, sys, idle, io)
		case cpu.INTR:
			s = fmt.Sprintf(CPU_SINGLE_STAT_FMT, "Interrupts", int(values[0]))
		case cpu.CTXT:
			s = fmt.Sprintf(CPU_SINGLE_STAT_FMT, "Context Switches", int(values[0]))
		case cpu.PROCS:
			s = fmt.Sprintf(CPU_SINGLE_STAT_FMT, "Processes", values[0])
		case cpu.PROCS_RUNNING:
			s = fmt.Sprintf(CPU_SINGLE_STAT_FMT, "Processes Running", values[0])
		case cpu.PROCS_BLOCKED:
			// intentionally left blank
		default:
			user, sys, idle, io := prepareCpuValues(values)
			s = fmt.Sprintf(CPU_STAT_FMT, "CPU"+d.Type(), user, sys, idle, io)
		}
		buf.WriteString(s)
	}
	return buf.String()
}

func formatMemStat(data []core.Stat) string {
	total := data[0].Values()[0]
	used := float64(data[1].Values()[0])
	cached := float64(data[2].Values()[0])
	swapTotal := data[3].Values()[0]
	swapUsed := float64(data[4].Values()[0])
	usedPct := (used / float64(total)) * 100
	cachedPct := (cached / float64(total)) * 100

	memStr := fmt.Sprintf(MEM_STAT_FMT, total, usedPct, cachedPct)
	if swapTotal > 0 {
		swapUsedPct := (swapUsed / float64(swapTotal)) * 100
		return memStr + fmt.Sprintf(SWAP_STAT_FMT, swapTotal, swapUsedPct)
	} else {
		return memStr
	}
}

func main() {
	// fetch /proc/diskstats (disk)
	// fetch /proc/meminfo (mem)

	cpuFh, err := os.Open("/proc/stat")
	if err != nil {
		log.Fatalf(err.Error())
	}
	defer cpuFh.Close()
	cpuStatProcessor := cpu.NewProcessor(cpuFh)
	cpuStatResults := make(chan []core.Stat)

	memFh, err := os.Open("/proc/meminfo")
	if err != nil {
		log.Fatalf(err.Error())
	}
	defer memFh.Close()
	memStatProcessor := mem.NewProcessor(memFh)
	memStatResults := make(chan []core.Stat)

	go core.StatProcessor(cpuStatProcessor, TICK_INTERVAL, cpuStatResults)
	go core.StatProcessor(memStatProcessor, TICK_INTERVAL, memStatResults)
	for {
		select {
		case c := <-cpuStatResults:
			fmt.Println(formatCpuStat(c))

		case c := <-memStatResults:
			fmt.Println(formatMemStat(c))
		}
	}
}

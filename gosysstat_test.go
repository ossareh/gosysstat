package main

import (
	"testing"

	"github.com/ossareh/libgosysstat/core"
	"github.com/ossareh/libgosysstat/processor/cpu"
)

func TestPrepareCpuValues(t *testing.T) {
	user, sys, idle, io := prepareCpuValues([]uint64{50, 50, 50, 50, 50})
	if user != 40.0 {
		t.Fatalf("Expected user to be %f, got %f", 40.0, user)
	}
	if sys != 20.0 {
		t.Fatalf("Expected sys to be %f, got %f", 20.0, user)
	}
	if idle != 20.0 {
		t.Fatalf("Expected idle to be %f, got %f", 20.0, user)
	}
	if io != 20.0 {
		t.Fatalf("Expected io to be %f, got %f", 20.0, user)
	}
}

type TestCpuStat struct {
	t string
}

func (t *TestCpuStat) Type() string {
	return t.t
}

func (t *TestCpuStat) Values() []uint64 {
	switch t.t {
	case cpu.TOTAL:
		return []uint64{50, 50, 50, 50, 50}
	case cpu.INTR:
		return []uint64{5000}
	case cpu.CTXT:
		return []uint64{6000}
	case cpu.PROCS:
		return []uint64{50}
	case cpu.PROCS_RUNNING:
		return []uint64{100}
	case cpu.PROCS_BLOCKED:
		return []uint64{10}
	default:
		return []uint64{50, 50, 50, 50, 50}
	}
}

func TestFormatCpuStat(t *testing.T) {
	str := formatCpuStat([]core.Stat{
		&TestCpuStat{cpu.TOTAL},
		&TestCpuStat{"0"},
		&TestCpuStat{"1"},
		&TestCpuStat{cpu.INTR},
		&TestCpuStat{cpu.CTXT},
		&TestCpuStat{cpu.PROCS},
		&TestCpuStat{cpu.PROCS_RUNNING},
		&TestCpuStat{cpu.PROCS_BLOCKED},
	})

	expected := "Total:(user:40.00%, sys:20.00%, idle:20.00%, io:20.00%)\n" +
		"CPU0:(user:40.00%, sys:20.00%, idle:20.00%, io:20.00%)\n" +
		"CPU1:(user:40.00%, sys:20.00%, idle:20.00%, io:20.00%)\n" +
		"Interrupts:5000\n" +
		"Context Switches:6000\n" +
		"Processes:50\n" +
		"Processes Running:100\n"

	if str != expected {
		t.Fatalf("Expected %s, got %s", expected, str)
	}
}

type TestMemStat struct {
	t string
	v []uint64
}

func (t *TestMemStat) Type() string {
	return t.t
}

func (t *TestMemStat) Values() []uint64 {
	return t.v
}

func TestFormatMemStatWithSwap(t *testing.T) {
	str := formatMemStat([]core.Stat{
		&TestMemStat{"total", []uint64{10000}},
		&TestMemStat{"used", []uint64{3000}},
		&TestMemStat{"cached", []uint64{4000}},
		&TestMemStat{"swap_total", []uint64{1000}},
		&TestMemStat{"swap_used", []uint64{500}},
	})

	expected := "RAM:(total:10000 KB, used:30.00%, cache:40.00%)\n" +
		"SWAP:(total:1000 KB, used:50.00%)\n"

	if str != expected {
		t.Fatalf("Expected %s, got %s", expected, str)
	}
}

func TestFormatMemStatNoSwap(t *testing.T) {
	str := formatMemStat([]core.Stat{
		&TestMemStat{"total", []uint64{10000}},
		&TestMemStat{"used", []uint64{3000}},
		&TestMemStat{"cached", []uint64{4000}},
		&TestMemStat{"swap_total", []uint64{0}},
		&TestMemStat{"swap_used", []uint64{0}},
	})

	expected := "RAM:(total:10000 KB, used:30.00%, cache:40.00%)\n"

	if str != expected {
		t.Fatalf("Expected %s, got %s", expected, str)
	}
}

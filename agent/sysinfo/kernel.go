package sysinfo

import (
	"bytes"
	"math"
	"os/exec"

	"github.com/lodastack/agent/agent/common"

	"github.com/lodastack/log"
	"github.com/lodastack/nux"
)

func FsKernelMetrics() (L []*common.Metric) {
	maxFiles, err := nux.KernelMaxFiles()
	if err != nil {
		log.Error("failed collect kernel metrics:", err)
		return
	}

	L = append(L, toMetric("kernel.files.max", maxFiles, nil))

	allocateFiles, err := nux.KernelAllocateFiles()
	if err != nil {
		log.Error("failed to call KernelAllocateFiles:", err)
		return
	}

	v := math.Ceil(float64(allocateFiles) * 100 / float64(maxFiles))
	L = append(L, toMetric("kernel.files.allocated", allocateFiles, nil))
	L = append(L, toMetric("kernel.files.allocated.percent", v, nil))
	L = append(L, toMetric("kernel.files.left", maxFiles-allocateFiles, nil))
	return
}

// exec `ps` to get all process states
func PsMetrics() (L []*common.Metric) {
	out, err := execPS()
	if err != nil {
		log.Error("failed to call ps command:", err)
		return
	}
	fields := make(map[string]int64)
	for i, status := range bytes.Fields(out) {
		if i == 0 && string(status) == "STAT" {
			// This is a header, skip it
			continue
		}
		switch status[0] {
		case 'W':
			fields["wait"] = fields["wait"] + int64(1)
		case 'U', 'D', 'L':
			// Also known as uninterruptible sleep or disk sleep
			fields["blocked"] = fields["blocked"] + int64(1)
		case 'Z':
			fields["zombies"] = fields["zombies"] + int64(1)
		case 'T':
			fields["stopped"] = fields["stopped"] + int64(1)
		case 'R':
			fields["running"] = fields["running"] + int64(1)
		case 'S':
			fields["sleeping"] = fields["sleeping"] + int64(1)
		case 'I':
			fields["idle"] = fields["idle"] + int64(1)
		case 'X':
			fields["exit"] = fields["exit"] + int64(1)
		case '?':
			fields["unknown"] = fields["unknown"] + int64(1)
		default:
			log.Errorf("processes: Unknown state [ %s ] from ps",
				string(status[0]))
		}
		fields["total"] = fields["total"] + int64(1)
	}
	L = append(L, toMetric("ps.zombies.num", fields["zombies"], nil))
	L = append(L, toMetric("ps.running.num", fields["running"], nil))
	L = append(L, toMetric("ps.total.num", fields["total"], nil))
	return
}

func execPS() ([]byte, error) {
	bin, err := exec.LookPath("ps")
	if err != nil {
		return nil, err
	}

	out, err := exec.Command(bin, "axo", "state").Output()
	if err != nil {
		return nil, err
	}

	return out, err
}

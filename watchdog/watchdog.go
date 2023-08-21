package watchdog

import (
	"os"
	"strconv"
	"time"

	"github.com/charlesbases/logger"
	"github.com/raulk/go-watchdog"
)

const envMaxHeapMem = "MAX_HEAP"

var (
	log = logger.Named("watchdog")

	watermarks = []float64{0.50, 0.60, 0.70, 0.80, 0.85, 0.90, 0.925, 0.95}
)

// heapMem 从环境变量 envMaxHeapMem 中获取最大内存限制
func heapMem() uint64 {
	val := os.Getenv(envMaxHeapMem)
	if len(val) == 0 {
		return 0
	}

	mem, err := strconv.ParseUint(val, 10, 64)
	if err != nil {
		logger.Warnf(`env %s[%s] cannot be converted to uint64.`, envMaxHeapMem, val)
		return 0
	}

	return mem
}

// Memory .
func Memory() (onstop func()) {
	// watchdog.HeapProfileDir = "heapprof"
	// watchdog.HeapProfileThreshold = 0.9
	// watchdog.HeapProfileMaxCaptures = 10
	watchdog.Logger = log

	policy := watchdog.NewWatermarkPolicy(watermarks...)

	// firs priority.
	if mem := heapMem(); mem != 0 {
		err, onstop := watchdog.HeapDriven(mem, 10, policy)
		if err == nil {
			log.Infof("initialized heap-driven watchdog. MaxHeapMem: %d bytes", mem)
			return onstop
		}

		log.Warnf("failed to initialize heap-driven watchdog. %s.", err.Error())
		log.Warn("trying a cgroup-driven watchdog ...")
	}

	// second priority.
	if err, onstop := watchdog.CgroupDriven(30*time.Second, policy); err == nil {
		log.Info("initialized cgroup-driven watchdog.")
		return onstop
	} else {
		log.Warnf("failed to initialize cgroup-driven watchdog. %s.", err.Error())
		log.Warnf("trying a system-driven watchdog ...")
	}

	// third priority.
	if err, onstop := watchdog.SystemDriven(0, 30*time.Second, policy); err == nil {
		log.Info("initialized system-driven watchdog.")
		return onstop
	} else {
		log.Warnf("failed to initialize system-driven watchdog. %s.", err)
		log.Warnf("system running without a memory watchdog.")
	}

	return func() {}
}

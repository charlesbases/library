package watchdog

import (
	"path/filepath"

	"github.com/raulk/go-watchdog"
)

// New .
func New() {
	watchdog.HeapProfileDir = filepath.Join("lr.Path()", "heapprof")
	watchdog.HeapProfileMaxCaptures = 10
	watchdog.HeapProfileThreshold = 0.9
}

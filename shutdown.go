package library

import (
	"os"
	"syscall"
)

// Shutdown returns all the singals that are being watched for to shut down services.
func Shutdown() []os.Signal {
	return []os.Signal{
		syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGSTOP, syscall.SIGKILL, syscall.SIGTERM,
	}
}

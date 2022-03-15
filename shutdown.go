package library

import (
	"os"
	"syscall"
)

// Shutdown returns all the singals that are being watched for to shut down services.
func Shutdown() []os.Signal {
	return []os.Signal{
		syscall.SIGHUP,  // 终端连接断开
		syscall.SIGINT,  // ctrl + C
		syscall.SIGQUIT, // ctrl + /
		syscall.SIGTERM, // kill pid
		// syscall.SIGKILL, // kill -9 pid // 不可捕获，立即退出
		// syscall.SIGSTOP, // 不可捕获
	}
}

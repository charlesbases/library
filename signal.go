package library

import (
	"os"
	"os/signal"
	"syscall"
)

// exited 退出信号
var exited = []os.Signal{
	syscall.SIGHUP,  // 终端连接断开
	syscall.SIGINT,  // ctrl + C
	syscall.SIGQUIT, // ctrl + /
	syscall.SIGTERM, // kill pid
	// syscall.SIGKILL, // kill -9 pid // 不可捕获，立即退出
	// syscall.SIGSTOP, // 不可捕获
}

// Shutdown returns all the singals that are being watched for to shut down services.
func Shutdown() chan os.Signal {
	var c = make(chan os.Signal)
	signal.Notify(c, exited...)
	return c
}

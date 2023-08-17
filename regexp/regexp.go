package regexp

import "regexp"

var (
	// IP ip 地址
	IP = regexp.MustCompile(`^(\d+\.){3}\d+$`)
	// HEX 十六进制
	HEX = regexp.MustCompile(`^(0x)?[0-9a-fA-F]+$`)
	// ServerName server name
	ServerName = regexp.MustCompile(`^[a-z][a-z0-9_]*$`)
)

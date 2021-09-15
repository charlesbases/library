package color

// BlackString .
func BlackString(format string, a ...interface{}) string {
	return New(FgBlack).Sprintf(format, a...)
}

// RedString .
func RedString(format string, a ...interface{}) string {
	return New(FgRed).Sprintf(format, a...)
}

// GreenString .
func GreenString(format string, a ...interface{}) string {
	return New(FgGreen).Sprintf(format, a...)
}

// YellowString .
func YellowString(format string, a ...interface{}) string {
	return New(FgYellow).Sprintf(format, a...)
}

// BlueString .
func BlueString(format string, a ...interface{}) string {
	return New(FgBlue).Sprintf(format, a...)
}

// MagentaString .
func MagentaString(format string, a ...interface{}) string {
	return New(FgMagenta).Sprintf(format, a...)
}

// CyanString .
func CyanString(format string, a ...interface{}) string {
	return New(FgCyan).Sprintf(format, a...)
}

// WhiteString .
func WhiteString(format string, a ...interface{}) string {
	return New(FgWhite).Sprintf(format, a...)
}

package library

import "regexp"

type regular struct {
	r *regexp.Regexp
}

var (
	// REGEXP_HEX 十六进制正则
	REGEXP_HEX *regular = &regular{regexp.MustCompile(`^(0x)?[0-9a-fA-F]+$`)}
)

// Match .
func (r *regular) Match(data *[]byte) bool {
	return r.r.Match(*data)
}

// MatchString .
func (r *regular) MatchString(data *string) bool {
	return r.r.MatchString(*data)
}

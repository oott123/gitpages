package mregex

import "regexp"

type Regexp struct {
	*regexp.Regexp
}

func (r *Regexp) UnmarshalText(data []byte) error {
	reg, err := regexp.Compile(string(data))
	if err == nil {
		*r = Regexp{reg}
	}
	return nil
}

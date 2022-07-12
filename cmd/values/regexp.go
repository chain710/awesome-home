package values

import (
	"github.com/spf13/pflag"
	"regexp"
)

type Regexp struct {
	impl *regexp.Regexp
}

var _ pflag.Value = &Regexp{}

func (r *Regexp) String() string {
	if r.impl == nil {
		return ""
	}
	return r.impl.String()
}

func (r *Regexp) Set(s string) error {
	if exp, err := regexp.Compile(s); err != nil {
		return err
	} else {
		r.impl = exp
		return nil
	}
}

func (r *Regexp) Type() string {
	return "Regexp"
}

func (r *Regexp) Regexp() *regexp.Regexp {
	return r.impl
}

func (r *Regexp) MatchString(s string) bool {
	if r.impl == nil {
		return true
	}

	return r.impl.MatchString(s)
}

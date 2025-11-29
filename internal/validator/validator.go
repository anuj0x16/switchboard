package validator

import "regexp"

var EmailRe = regexp.MustCompile(`^[a-zA-Z0-9.!#$%&'*+/=?^_` + "`" + `{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$`)

type Validator struct {
	Errors map[string]string
}

func New() *Validator {
	return &Validator{
		Errors: make(map[string]string),
	}
}

func (v *Validator) Valid() bool {
	return len(v.Errors) == 0
}

func (v *Validator) AddError(key string, err string) {
	if _, exists := v.Errors[key]; !exists {
		v.Errors[key] = err
	}
}

func (v *Validator) Check(ok bool, key string, err string) {
	if !ok {
		v.AddError(key, err)
	}
}

func Match(re *regexp.Regexp, s string) bool {
	return re.MatchString(s)
}

package validator

import (
	"regexp"
	"strings"
)

// source - https://gist.github.com/brydavis/0c7da92bd508195744708eeb2b54ac96
var URLRX = regexp.MustCompile("^(http://www.|https://www.|http://|https://|/|//)?[A-z0-9_-]*?[:]?[A-z0-9_-]*?[@]?[A-z0-9]+([-.]{1}[a-z0-9]+)*.[a-z]{2,5}(:[0-9]{1,5})?(/.*)?$")

type Validator struct {
	FieldError map[string]string
}

func (v *Validator) Valid() bool {
	return len(v.FieldError) == 0
}

func (v *Validator) AddFieldError(key, message string) {
	if v.FieldError == nil {
		v.FieldError = make(map[string]string)
	}
	if _, ok := v.FieldError[key]; !ok {
		v.FieldError[key] = message
	}
}

func (v *Validator) CheckField(ok bool, key, message string) {
	if !ok {
		v.AddFieldError(key, message)
	}
}

func NotBlank(value string) bool {
	return strings.TrimSpace(value) != ""
}

func Matches(value string, rx *regexp.Regexp) bool {
	return rx.MatchString(value)
}

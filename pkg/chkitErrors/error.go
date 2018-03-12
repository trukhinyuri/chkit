package chkitErrors

import (
	"bytes"
	"fmt"
)

type Err string

func (err Err) Error() string {
	return string(err)
}

func (err Err) Errors() []error {
	return []error{err}
}

type Wrapper struct {
	main          error
	reasons       []error
	cachedMessage string
}

func (wrapper *Wrapper) AddReasons(reasons ...error) *Wrapper {
	for _, reason := range reasons {
		wrapper.reasons = append(wrapper.reasons, reason)
	}
	return wrapper
}

func (wrapper *Wrapper) AddReasonF(f string, vals ...interface{}) *Wrapper {
	return wrapper.AddReasons(fmt.Errorf(f, vals...))
}
func (wrapper *Wrapper) Error() string {
	if wrapper.cachedMessage != "" {
		return wrapper.cachedMessage
	}
	buf := bytes.NewBufferString(wrapper.main.Error())
	if len(wrapper.reasons) > 0 {
		buf.WriteString(": ")
	}
	for i, reason := range wrapper.reasons {
		if i != 0 {
			buf.WriteString(", " + reason.Error())
		} else {
			buf.WriteString(reason.Error())
		}
	}
	wrapper.cachedMessage = buf.String()
	return wrapper.cachedMessage
}

func (wrapper *Wrapper) Match(errs ...error) bool {
	for _, err := range errs {
		if wrapper.main == err {
			return true
		}
	}
	return false
}

func (wrapper *Wrapper) Errors() []error {
	errs := make([]error, len(wrapper.reasons))
	copy(errs, wrapper.reasons)
	return errs
}

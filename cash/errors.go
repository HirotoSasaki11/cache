package cash

import (
	"strings"
)

type MultiError []error

func (e *MultiError) Append(err error) {
	if err == nil {
		return
	}
	if e == nil {
		*e = MultiError{}
	}
	*e = append(*e, err)
}

func (e MultiError) Error() string {
	msgs := make([]string, len(e))
	for i := range e {
		msgs[i] = e[i].Error()
	}
	return strings.Join(msgs, "\n")
}

func (e MultiError) ErrorOrNil() error {
	if len(e) == 0 {
		return nil
	} else {
		return e
	}
}

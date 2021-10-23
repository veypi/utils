package utils

import "errors"

func MultiErr(errs ...error) error {
	var s = ""
	for _, e := range errs {
		if e != nil {
			s += e.Error() + "\n"
		}
	}
	if s == "" {
		return nil
	}
	return errors.New(s)
}

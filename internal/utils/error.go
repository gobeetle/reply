package utils

import "errors"

func ErrorsToStrings(errors []error) []string {
	strSlice := make([]string, len(errors))
	for i, err := range errors {
		if err != nil {
			strSlice[i] = err.Error()
		}
	}
	return strSlice
}

func StringsToErrors(strSlice []string) []error {
	errs := make([]error, len(strSlice))
	for i, str := range strSlice {
		errs[i] = errors.New(str)
	}
	return errs
}

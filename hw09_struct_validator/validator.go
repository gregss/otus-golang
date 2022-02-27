package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

type ValidationError struct {
	Field string
	Err   error
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	errstr := ""
	for _, e := range v {
		if e.Err != nil {
			errstr += fmt.Sprintf("field: %v, err: %v", e.Field, e.Err)
		}
	}

	return errstr
}

type OtherError struct {
	Err error
}

func (v OtherError) Error() string {
	return v.Err.Error()
}

func Validate(v interface{}) error { //nolint:gocognit
	st := reflect.TypeOf(v)
	if st.Kind() != reflect.Struct {
		return OtherError{errors.New("is`t type struct")}
	}

	stv := reflect.ValueOf(v)
	verrors := make(ValidationErrors, 0)
	for i := 0; i < st.NumField(); i++ {
		f := st.Field(i)
		vst := f.Tag.Get("validate")
		if vst == "" {
			continue
		}

		for _, vrule := range strings.Split(vst, "|") {
			idx := strings.Index(vrule, ":")
			switch vrule[:idx] {
			case "len":
				ruleval, _ := strconv.Atoi(vrule[idx+1:])
				if ruleval != len(stv.Field(i).String()) {
					verrors = append(verrors, ValidationError{f.Name, errors.New(vrule[:idx])})
				}
			case "min":
				ruleval, _ := strconv.Atoi(vrule[idx+1:])
				if stv.Field(i).Int() < int64(ruleval) {
					verrors = append(verrors, ValidationError{f.Name, errors.New(vrule[:idx])})
				}
			case "max":
				ruleval, _ := strconv.Atoi(vrule[idx+1:])
				if stv.Field(i).Int() > int64(ruleval) {
					verrors = append(verrors, ValidationError{f.Name, errors.New(vrule[:idx])})
				}
			case "in":
				for _, ruleval := range strings.Split(vrule[idx+1:], ",") {
					if ruleval == stv.Field(i).String() {
						break
					}
					verrors = append(verrors, ValidationError{f.Name, errors.New(vrule[:idx])})
				}
			case "regexp":
				rg, _ := regexp.Compile(vrule[idx+1:])
				rg.MatchString(stv.Field(i).String())
				if !rg.MatchString(stv.Field(i).String()) {
					verrors = append(verrors, ValidationError{f.Name, errors.New(vrule[:idx])})
				}
			default:
			}
		}
	}

	if len(verrors) > 0 {
		return verrors
	}

	return nil
}

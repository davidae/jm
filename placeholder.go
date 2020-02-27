package jm

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"time"
)

// Placeholder is a function used to validate partially unknown values
// in the actual JSON by assigning .
type Placeholder func() (string, func(actualValue interface{}) error)

// WithNotEmpty looks for a placeholder value in the expected JSON
// and confirms the actual value is not empty.
func WithNotEmpty(placeholder string) Placeholder {
	return func() (string, func(interface{}) error) {
		return placeholder, func(val interface{}) error {
			if isEmpty(val) {
				return errors.New("expected value to be not empty, but it was")
			}

			return nil
		}
	}
}

// WithRegexp looks for the placeholder value in the expected JSON and confirms the
// actual value matches the given regular expression.
func WithRegexp(placeholder string, re *regexp.Regexp) Placeholder {
	return func() (string, func(interface{}) error) {
		return placeholder, func(val interface{}) error {
			str, ok := val.(string)
			if !ok {
				return fmt.Errorf("cannot match, value is of type %T - not a string", val)
			}

			if re.MatchString(str) {
				return nil
			}

			return fmt.Errorf("value %s does not match with regexp %s", str, re)
		}
	}
}

// WithTimeLayout looks for the placeholder value in the expected JSON
// and confirms the actual value matches the given time layout format
func WithTimeLayout(placeholder, layout string) Placeholder {
	return func() (string, func(interface{}) error) {
		return placeholder, func(val interface{}) error {
			str, ok := val.(string)
			if !ok {
				return fmt.Errorf("cannot parse time, value is of type %T - not a string", val)
			}

			if _, err := time.Parse(layout, str); err != nil {
				return fmt.Errorf("cannot parse layout %q with %q", layout, str)
			}

			return nil
		}
	}
}

func isEmpty(i interface{}) bool {
	if i == nil {
		return true
	}

	val := reflect.ValueOf(i)

	switch val.Kind() {
	case reflect.Array, reflect.Chan, reflect.Map, reflect.Slice:
		return val.Len() == 0
	case reflect.Ptr:
		if val.IsNil() {
			return true
		}
		deref := val.Elem().Interface()
		return isEmpty(deref)
	default:
		zero := reflect.Zero(val.Type())
		return reflect.DeepEqual(i, zero.Interface())
	}
}

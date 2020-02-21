package jm

import (
	"fmt"
	"reflect"
	"regexp"
	"time"
)

// Placeholder is a function used to validate partially unknown values
// in the actual JSON by assigning .
type Placeholder func() (string, func(actualValue interface{}) error)

// WithNotEmpty is a Placeholder function that looks for '$NOT_EMPTY' as a
// placeholder value in the expected JSON and confirms the actual value is
// not empty.
func WithNotEmpty() Placeholder {
	return func() (string, func(interface{}) error) {
		return "$NOT_EMPTY", func(val interface{}) error {
			if isEmpty(val) {
				return fmt.Errorf("unexpected empty value for placeholder $NOT_EMPTY")
			}

			return nil
		}
	}
}

// WithRegexp is a Placeholder function that looks for '$MATCHES_REGEXP' as a
// placeholder value in the expected JSON and confirms the actual value is
// matches the given regular expression.
func WithRegexp(re *regexp.Regexp) Placeholder {
	return func() (string, func(interface{}) error) {
		return "$MATCHES_REGEXP", func(val interface{}) error {
			str, ok := val.(string)
			if !ok {
				return fmt.Errorf("cannot match string, value %#v is not a string, it is of type %T", val, val)
			}

			if re.MatchString(str) {
				return nil
			}

			return fmt.Errorf("")
		}
	}
}

// WithTimeLayout is a Placeholder function that looks for '$TIME_LAYOUT' as a
// placeholder value in the expected JSON and confirms the actual value
// matches the given time layout format
func WithTimeLayout(layout string) Placeholder {
	return func() (string, func(interface{}) error) {
		return "$TIME_LAYOUT", func(val interface{}) error {
			str, ok := val.(string)
			if !ok {
				return fmt.Errorf("cannot parse time string, %#v is not a string; it is of type %T", val, val)
			}

			_, err := time.Parse(layout, str)
			return err
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

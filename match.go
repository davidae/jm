package jm

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
)

var (
	// ErrMissingKey is an error when an expecting key is missing from the JSON
	ErrMissingKey = errors.New("missing key")
	// ErrUnknownKey is an error when an unexpected key is present in the JSON
	ErrUnknownKey = errors.New("unknown key")
	// ErrArrayLengths is an error when comparing two arrys and their lengths are not equal
	ErrArrayLengths = errors.New("array lengths are not equal")
	// ErrNotEqualValues is an error when two expected equal values are not equal
	ErrNotEqualValues = errors.New("values are not equal")
)

// Match will compare two JSONs, expected and actual, and will return an erorr if
// there are any discrepancies between them. It will return the first mismatch that is found
// (depth-first search), additional mismatches is always plausible, but not reported.
//
// Placeholders can be used to validate, infer, or allow partially unknown values in the actual JSON.
// refer to the README or tests to see how they can be applied
func Match(expected, actual []byte, placeholders ...Placeholder) error {
	var exp, act interface{}

	if err := json.Unmarshal(expected, &exp); err != nil {
		return err
	}

	if err := json.Unmarshal(actual, &act); err != nil {
		return err
	}

	return isEqual(exp, act, "", placeholders...)
}

func isEqual(expected, actual interface{}, key string, ph ...Placeholder) error {
	if ea, aa, ok := isArray(expected, actual); ok {
		l, err := areLenEqual(ea, aa)
		if err != nil {
			return errUnderKey(err, key)
		}

		for i := 0; i < l; i++ {
			if err := isEqual(ea[i], aa[i], key, ph...); err != nil {
				return err
			}
		}

		return nil
	}

	if eo, ao, ok := isObject(expected, actual); ok {
		if err := areKeysEqual(eo, ao); err != nil {
			return err
		}
		for k := range eo {
			if err := isEqual(eo[k], ao[k], k, ph...); err != nil {
				return err
			}
		}

		return nil
	}

	if err := isEqualValue(expected, actual, ph); err != nil {
		return err
	}

	return nil
}

func isEqualValue(expected, actual interface{}, ph []Placeholder) error {
	for _, p := range ph {
		expectedStr, ok := expected.(string)
		if !ok {
			continue
		}

		if key, fn := p(); key == expectedStr {
			if err := fn(actual); err != nil {
				return fmt.Errorf("placeholder %s match failed: %w", key, err)
			}

			return nil
		}
	}

	if !reflect.DeepEqual(expected, actual) {
		return fmt.Errorf("value %#v and %#v: %w", expected, actual, ErrNotEqualValues)
	}

	return nil
}

func isArray(i, j interface{}) ([]interface{}, []interface{}, bool) {
	ia, ok := i.([]interface{})
	if !ok {
		return []interface{}{}, []interface{}{}, false
	}

	ja, ok := j.([]interface{})
	if !ok {
		return []interface{}{}, []interface{}{}, false
	}

	return ia, ja, true
}

func isObject(i, j interface{}) (map[string]interface{}, map[string]interface{}, bool) {
	io, ok := i.(map[string]interface{})
	if !ok {
		return map[string]interface{}{}, map[string]interface{}{}, false
	}
	jo, ok := j.(map[string]interface{})
	if !ok {
		return map[string]interface{}{}, map[string]interface{}{}, false
	}

	return io, jo, true
}

func areLenEqual(i, j []interface{}) (int, error) {
	il, jl := len(i), len(j)
	if il != jl {
		return 0, fmt.Errorf("mismatch array length %d and %d: %w", il, jl, ErrArrayLengths)
	}

	return il, nil
}

type keyMatcher struct {
	expected, actual bool
}

func areKeysEqual(expected, actual map[string]interface{}) error {
	matches := make(map[string]*keyMatcher)

	for k := range expected {
		matches[k] = &keyMatcher{expected: true}
	}

	for k := range actual {
		m, ok := matches[k]
		if !ok {
			return fmt.Errorf("key %q is only present in actual JSON: %w", k, ErrUnknownKey)
		}

		m.actual = true
	}

	for k, match := range matches {
		if !match.actual {
			return fmt.Errorf("key %q is not present in actual JSON: %w", k, ErrMissingKey)
		}
	}

	return nil
}

func errUnderKey(err error, key string) error {
	if key != "" {
		return fmt.Errorf("mismatch under key %s: %w", key, err)
	}

	return err
}

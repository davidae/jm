package jm

import (
	"fmt"
	"io/ioutil"
	"testing"
)

func TestEqualJSON(t *testing.T) {
	out := mustReadFile(t, "test/stubs/match/expected.json")

	if err := Match(out, out); err != nil {
		t.Errorf("unexpected error, expected JSON to be equal itself: %s", err)
	}
}

func TestNotEqualJSONWithAdditionalArrayItem(t *testing.T) {
	var (
		expected       = mustReadFile(t, "test/stubs/match/expected.json")
		actual         = mustReadFile(t, "test/stubs/match/additional_array_item.json")
		expectedErrMsg = "mismatch under key arr_1: mismatch array length 2 and 3: array lengths are not equal"
	)

	if err := Match(expected, actual); err == nil {
		t.Error("expected an error, but nil was returned")
	} else if err.Error() != expectedErrMsg {
		t.Errorf("expected error message %s to be, but got %s", expectedErrMsg, err)
	}
}

func TestNotEqualJSONWithMismatchValue(t *testing.T) {
	var (
		expected       = mustReadFile(t, "test/stubs/match/expected.json")
		actual         = mustReadFile(t, "test/stubs/match/mismatch_value.json")
		expectedErrMsg = `mismatch under key arr_2: value 2 and "hex": values are not equal`
	)

	if err := Match(expected, actual); err == nil {
		t.Error("expected an error, but nil was returned")
	} else if err.Error() != expectedErrMsg {
		t.Errorf("expected error message %s to be, but got %s", expectedErrMsg, err)
	}
}

func TestNotEqualJSONWithMissingArrayItem(t *testing.T) {
	var (
		expected       = mustReadFile(t, "test/stubs/match/expected.json")
		actual         = mustReadFile(t, "test/stubs/match/missing_array_item.json")
		expectedErrMsg = "mismatch under key arr_2: mismatch array length 3 and 2: array lengths are not equal"
	)

	if err := Match(expected, actual); err == nil {
		t.Error("expected an error, but nil was returned")
	} else if err.Error() != expectedErrMsg {
		t.Errorf("expected error message %s to be, but got %s", expectedErrMsg, err)
	}
}

func TestNotEqualJSONWithMissingKey(t *testing.T) {
	var (
		expected       = mustReadFile(t, "test/stubs/match/expected.json")
		actual         = mustReadFile(t, "test/stubs/match/missing_key.json")
		expectedErrMsg = `mismatch under key double: key "nothing" is not present in actual JSON: missing key`
	)

	if err := Match(expected, actual); err == nil {
		t.Error("expected an error, but nil was returned")
	} else if err.Error() != expectedErrMsg {
		t.Errorf("expected error message %s to be, but got %s", expectedErrMsg, err)
	}
}

func TestNotEqualJSONWithUnexpectedKey(t *testing.T) {
	var (
		expected       = mustReadFile(t, "test/stubs/match/expected.json")
		actual         = mustReadFile(t, "test/stubs/match/unexpected_key.json")
		expectedErrMsg = `mismatch under key double: key "wait_for_me" is only present in actual JSON: unknown key`
	)

	if err := Match(expected, actual); err == nil {
		t.Error("expected an error, but nil was returned")
	} else if err.Error() != expectedErrMsg {
		t.Errorf("expected error message %s to be, but got %s", expectedErrMsg, err)
	}
}

func TestNotEqualJSONWithNotEmptyPlaceholder(t *testing.T) {
	var (
		expected = mustReadFile(t, "test/stubs/placeholder/with_not_empty.json")
		actual   = mustReadFile(t, "test/stubs/placeholder/actual.json")
	)

	if err := Match(expected, actual, WithNotEmpty()); err != nil {
		t.Errorf("unexpected error, expected JSON to be equal itself: %s", err)
	}
}

func TestEqualWithCustomPlaceholder(t *testing.T) {
	gte3 := func() (string, func(interface{}) error) {
		return "$GTE_3", func(val interface{}) error {
			valFloat, ok := val.(float64)
			if !ok {
				return fmt.Errorf("expected value to compare to be an float64 but got: %T", val)
			}

			if valFloat >= 3 {
				return nil
			}

			return fmt.Errorf("%f is not greater or equal than 3", valFloat)
		}
	}

	if err := Match(
		[]byte(`{"value": "$GTE_3"}`),
		[]byte(`{"value": 3}`), gte3); err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	if err := Match(
		[]byte(`{"value": "$GTE_3"}`),
		[]byte(`{"value": 2}`), gte3); err == nil {
		t.Error("expected an error, but it was nil")
	}
}

func TestEmptyBytes(t *testing.T) {
	expectedErrMsg := "unexpected end of JSON input"

	if err := Match([]byte{}, []byte{}); err == nil {
		t.Error("expected an error, but nil was returned")
	} else if err.Error() != expectedErrMsg {
		t.Errorf("expected error message %s to be, but got %s", expectedErrMsg, err)
	}
}

func mustReadFile(t *testing.T, filename string) []byte {
	out, err := ioutil.ReadFile(filename)
	if err != nil {
		t.Error("failed to read JSON")
	}

	return out
}

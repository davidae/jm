package jm

import (
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
		expectedErrMsg = `value 2 and "hex": values are not equal`
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
		expectedErrMsg = `key "nothing" is not present in actual JSON: missing key`
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
		expectedErrMsg = `key "wait_for_me" is only present in actual JSON: unknown key`
	)

	if err := Match(expected, actual); err == nil {
		t.Error("expected an error, but nil was returned")
	} else if err.Error() != expectedErrMsg {
		t.Errorf("expected error message %s to be, but got %s", expectedErrMsg, err)
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

package jm

import (
	"io/ioutil"
	"testing"
	"time"
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
		t.Errorf("expected error message to be %q, but got %q", expectedErrMsg, err)
	}
}

func TestNotEqualJSONWithMismatchValue(t *testing.T) {
	var (
		expected       = mustReadFile(t, "test/stubs/match/expected.json")
		actual         = mustReadFile(t, "test/stubs/match/mismatch_value.json")
		expectedErrMsg = "mismatch under key arr_2: element [{},[2]] appears more times in [\"A\",1,[{},[2]]] than in [\"A\",1,[{},[\"hex\"]]]"
	)

	if err := Match(expected, actual); err == nil {
		t.Error("expected an error, but nil was returned")
	} else if err.Error() != expectedErrMsg {
		t.Errorf("expected error message to be %q, but got %q", expectedErrMsg, err)
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

func TestEqualWithIntegerMismatch(t *testing.T) {
	var (
		expected       = `{"id": "1","count":2}`
		actual         = `{"id": "1","count":4}`
		expectedErrMsg = "value 2 and 4: values are not equal"
	)

	if err := Match([]byte(expected), []byte(actual), WithTimeLayout("$TIME_RFC3339", time.RFC3339)); err == nil {
		t.Error("expected a mismatch on 'count'")
	} else if err.Error() != expectedErrMsg {
		t.Errorf("expected error message %s to be, but got %s", expectedErrMsg, err)
	}
}

func TestEqualWithArrayIntegerMismatch(t *testing.T) {
	var (
		expected       = `{"id": "1","count":[1,2,3,4,5,6]}`
		actual         = `{"id": "1","count":[1,2,3,4,"s",6]}`
		expectedErrMsg = "mismatch under key count: element 5 appears more times in [1,2,3,4,5,6] than in [1,2,3,4,\"s\",6]"
	)

	if err := Match([]byte(expected), []byte(actual), WithTimeLayout("$TIME_RFC3339", time.RFC3339)); err == nil {
		t.Error("expected a mismatch on 'count'")
	} else if err.Error() != expectedErrMsg {
		t.Errorf("expected error message to be %q, but got %q", expectedErrMsg, err)
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

func TestEqualWithUnorderedArrayWhenIsEqualWithPrimitiveDataTypes(t *testing.T) {
	var (
		expected = `{"id": "1","count":[1,2,["a","c","b"]]}`
		actual   = `{"id": "1","count":[2,1,["a","b","c"]]}`
	)

	if err := Match([]byte(expected), []byte(actual)); err != nil {
		t.Errorf("unexpected mismatch: %s", err)
	}
}

func TestEqualWithUnorderedArrayWhenIsEqualWithStruct(t *testing.T) {
	var (
		expected = `{"id": "1","count":[1,2,["a","c",{"flag":true}]]}`
		actual   = `{"id": "1","count":[2,1,["a",{"flag":true},"c"]]}`
	)

	if err := Match([]byte(expected), []byte(actual)); err != nil {
		t.Errorf("unexpected mismatch: %s", err)
	}
}

func TestEqualWithUnorderedArrayWhenIsNotEqual(t *testing.T) {
	var (
		expected       = `{"id": "1","count":[1,2,["a","c","e"]]}`
		actual         = `{"id": "1","count":[2,1,["a","b","c"]]}`
		expectedErrMsg = "mismatch under key count: element [\"a\",\"c\",\"e\"] appears more times in [1,2,[\"a\",\"c\",\"e\"]] than in [2,1,[\"a\",\"b\",\"c\"]]"
	)

	if err := Match([]byte(expected), []byte(actual)); err == nil {
		t.Error("expected an error, but nil was returned")
	} else if err.Error() != expectedErrMsg {
		t.Errorf("expected error message to be %q, but got %q", expectedErrMsg, err)
	}
}

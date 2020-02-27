package jm

import (
	"fmt"
	"regexp"
	"testing"
	"time"
)

const (
	nestedActualJSON = `{
		"sample": {
			"an_int": 12,
			"a_float": 32.3,
			"nested": {
				"a_bool": false,
				"nothing": null,
				"double": {
					"arr_1": ["A", "B"],
					"arr_2": ["A", 2.9998, [{},[2]]]}}}}`

	nestedExpectedJSON = `{
	"sample": {
		"an_int": 12,
		"a_float": 32.3,
		"nested": {
            "a_bool": false,
			"nothing": null,
			"double": {
                "arr_1": ["A", "B"],
				"arr_2": ["A", "$GTE_3", [{},[2]]]}}}}`
)

func TestNestedJSONAndCustomPlaceholder(t *testing.T) {
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

	expectedErrStr := "placeholder $GTE_3 match failed: 2.999800 is not greater or equal than 3"

	err := Match([]byte(nestedExpectedJSON), []byte(nestedActualJSON), gte3)
	if err == nil {
		t.Errorf("expected an error, not none")
	} else if err.Error() != expectedErrStr {
		t.Errorf("expected error string %q, but got %q", expectedErrStr, err.Error())
	}
}

func TestWithNotEmptyPlaceholderWhenMatch(t *testing.T) {
	if err := Match(
		[]byte(`{"name": "$NOT_EMPTY"}`),
		[]byte(`{"name": "John"}`),
		WithNotEmpty("$NOT_EMPTY"),
	); err != nil {
		t.Errorf("unexpected error %s, expected JSONs to be equal", err)
	}
}

func TestWithNotEmptyPlaceholderWhenEmpty(t *testing.T) {
	expectedErrStr := "placeholder $NOT_EMPTY match failed: expected value to be not empty, but it was"
	err := Match(
		[]byte(`{"name": "$NOT_EMPTY"}`),
		[]byte(`{"name": ""}`),
		WithNotEmpty("$NOT_EMPTY"),
	)
	if err == nil {
		t.Error("expected an error, got none")
	} else if err.Error() != expectedErrStr {
		t.Errorf("expected error string %q, but got %q", expectedErrStr, err.Error())
	}
}

func TestWithTimeLayoutWhenMatch(t *testing.T) {
	if err := Match(
		[]byte(`{"created_at": "$TIME_RFC3339"}`),
		[]byte(`{"created_at": "2009-11-10T23:00:00Z"}`),
		WithTimeLayout("$TIME_RFC3339", time.RFC3339),
	); err != nil {
		t.Errorf("unexpected error %s, expected JSONs to be equal", err)
	}
}

func TestWithTimeLayoutWhenInvalidTimeString(t *testing.T) {
	expectedErrStr := "placeholder $TIME_RFC3339 match failed: cannot parse time, value is of type bool - not a string"
	err := Match(
		[]byte(`{"created_at": "$TIME_RFC3339"}`),
		[]byte(`{"created_at": false}`),
		WithTimeLayout("$TIME_RFC3339", time.RFC3339),
	)
	if err == nil {
		t.Error("expected an error, got none")
	} else if err.Error() != expectedErrStr {
		t.Errorf("expected error string %q, but got %q", expectedErrStr, err.Error())
	}
}

func TestWithTimeLayoutWhenInvalidTimeFormat(t *testing.T) {
	expectedErrStr := "placeholder $TIME_RFC3339 match failed: cannot parse layout \"2006-01-02T15:04:05Z07:00\" with \"hello\""
	err := Match(
		[]byte(`{"created_at": "$TIME_RFC3339"}`),
		[]byte(`{"created_at": "hello"}`),
		WithTimeLayout("$TIME_RFC3339", time.RFC3339),
	)
	if err == nil {
		t.Error("expected an error, got none")
	} else if err.Error() != expectedErrStr {
		t.Errorf("expected error string %q, but got %q", expectedErrStr, err.Error())
	}
}

func TestWithRegexpWhenMatch(t *testing.T) {
	re := regexp.MustCompile("^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}$")

	if err := Match([]byte(`{"uuid": "$UUID"}`), []byte(`{"uuid": "4d6d4d98-996a-4364-8514-e936f1db552b"}`),
		WithRegexp("$UUID", re),
	); err != nil {
		t.Errorf("unexpected error %s, expected JSONs to be equal", err)
	}
}
func TestWithRegexpWhenMatchWithInvalidRegexString(t *testing.T) {
	re := regexp.MustCompile("^2")
	expectedErrStr := "placeholder $UUID match failed: cannot match, value is of type float64 - not a string"

	err := Match([]byte(`{"uuid": "$UUID"}`), []byte(`{"uuid": 1}`), WithRegexp("$UUID", re))
	if err == nil {
		t.Error("expected an error, got none")
	} else if err.Error() != expectedErrStr {
		t.Errorf("expected error string %q, but got %q", expectedErrStr, err.Error())
	}
}

func TestWithRegexpWhenMatchWithNoMatch(t *testing.T) {
	re := regexp.MustCompile("^2")
	expectedErrStr := "placeholder $UUID match failed: value z does not match with regexp ^2"

	err := Match(
		[]byte(`{"uuid": "$UUID"}`),
		[]byte(`{"uuid": "z"}`),
		WithRegexp("$UUID", re))
	if err == nil {
		t.Error("expected an error, got none")
	} else if err.Error() != expectedErrStr {
		t.Errorf("expected error string %q, but got %q", expectedErrStr, err.Error())
	}
}

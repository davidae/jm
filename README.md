# jm - jsonMatch
A small package to compare JSONs with the advantage of using placeholders for possible unknown values. 
It should work with any type of valid JSON - however nested and tangled it may be. The package only uses 
golangs standard library, no dependencies.

# Installation
```Shell
$ go get github.com/davidae/jm
```

# Usage
## Without placeholders
You can compare two JSON by simply doing
```go
err := jm.Match(
    []byte(`{"hello": "world"}`),
    []byte(`{"hello": "friend"}`))

fmt.Printf("we should have an ErrNotEqualValues here: %t", errors.Is(err, jm.ErrNotEqualValues))
```
If there is an error, it implies that the two JSONs are not equal. The error should indicate why they were not
equal. 

## With placeholders
You can use placeholders in the expected JSON when it is hard to determine the actual value, such as
```go
err := jm.Match(
    []byte(`{"created_at": "$TIME_LAYOUT"}`),
    []byte(`{"created_at": "2009-11-10T23:00:00Z"}`),
    jm.WithTimeLayout(time.RFC3339))

fmt.Println("we should not have an error here: %t", err == nil)
```
It is also possible to define your own custom placeholders,
```go
err := jm.Match(
    []byte(`{"value": "$GTE_3"}`),
    []byte(`{"value": 3}`),
    func() (string, func(interface{}) error) {
        return "$GTE_3", func(val interface{}) error {
            valFloat, ok := val.(float64)
            if !ok {
                return fmt.Errorf("expected value be a float64 but got: %T", val)
            }

            if valFloat >= 3 {
                return nil
            }

            return fmt.Errorf("%f is not greater or equal than 3", valFloat)
        }
    },
)

fmt.Println("we should not have an error here: %t", err == nil)	
```

## Errors
This package uses go1.13+ errors and each error should be possible to unwrap into one of the following,
```go
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
```

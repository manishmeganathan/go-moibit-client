package moibit

import (
	"errors"
	"fmt"
	"reflect"
	"testing"
)

func ExampleNewFilePath() {
	directory := "appdata/users"
	filename := "userdata.json"

	fp, _ := NewFilePath(directory, filename)
	fmt.Println("Path:", fp)

	// Output:
	// Path: /appdata/users/userdata.json
}

func TestNewFilePath(t *testing.T) {
	tests := []struct {
		input []string
		path  string
		error error
	}{
		{[]string{"data", "devices"}, "/data/devices", nil},
		{[]string{"//data", "devices/", "hello.txt"}, "/data/devices/hello.txt", nil},
		{[]string{"data/devices", "hello.txt"}, "/data/devices/hello.txt", nil},
		{[]string{"/data"}, "/data", nil},
		{[]string{"/data", ""}, "/data", nil},
		{[]string{"data", "/devices"}, "/data/devices", nil},
		{[]string{"data/", "/devices/hello.txt"}, "/data/devices/hello.txt", nil},
		{
			[]string{"data", "/devices.txt/hello"}, "/",
			errors.New("failed to construct filepath: slash detected after period or missing extension"),
		},
		{
			[]string{"data", "/devices.txt/hello.txt"}, "/",
			errors.New("failed to construct filepath: multiple periods in final element"),
		},
		{
			[]string{"dat.txt", "/devices.txt/hello.txt"}, "/",
			errors.New("failed to construct filepath: non-final element '0' contains period"),
		},
	}

	for _, test := range tests {
		fp, err := NewFilePath(test.input...)

		if !reflect.DeepEqual(err, test.error) {
			t.Fatalf("error mismatch. expected: %v. got: %v", test.error, err)
		}

		if !reflect.DeepEqual(test.path, fp.Path()) {
			t.Fatalf("path mismatch. expected: %v. got: %v", test.path, fp.Path())
		}
	}
}

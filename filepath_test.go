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

func TestFilePath_Check(t *testing.T) {
	tests := []struct {
		input  []string
		isroot bool
		isdir  bool
		isfile bool
	}{
		{[]string{"data", "devices"}, false, true, false},
		{[]string{""}, true, true, false},
		{[]string{"data", "hello.txt"}, false, false, true},
	}

	for _, test := range tests {
		fp, _ := NewFilePath(test.input...)

		if fp.IsRoot() != test.isroot {
			t.Fatalf("root check mismatch. expected: %v. got: %v", test.isroot, fp.IsRoot())
		}

		if fp.IsFile() != test.isfile {
			t.Fatalf("file check mismatch. expected: %v. got: %v", test.isfile, fp.IsFile())
		}

		if fp.IsDirectory() != test.isdir {
			t.Fatalf("dir check mismatch. expected: %v. got: %v", test.isdir, fp.IsDirectory())
		}
	}
}

func TestFilePath_Parent(t *testing.T) {
	tests := []struct {
		input  []string
		parent string
	}{
		{[]string{"data", "devices"}, "/data"},
		{[]string{"data"}, "/"},
		{[]string{"data/hello.txt"}, "/data"},
		{[]string{""}, "/"},
		{[]string{"data", "devices", "laptop.txt"}, "/data/devices"},
	}

	for _, test := range tests {
		fp, _ := NewFilePath(test.input...)

		if !reflect.DeepEqual(test.parent, fp.Parent().Path()) {
			t.Fatalf("parent path mismatch. expected: %v. got: %v", test.parent, fp.Parent().Path())
		}
	}
}

func TestFilePath_Grow(t *testing.T) {
	tests := []struct {
		input  []string
		growth []string
		error  error
		path   string
	}{
		{
			[]string{"data", "devices"}, []string{"laptop"},
			nil, "/data/devices/laptop",
		},
		{
			[]string{"data", "devices"}, []string{"laptop.txt"},
			nil, "/data/devices/laptop.txt",
		},
		{
			[]string{"data", "devices.txt"}, []string{"laptop"},
			errors.New("cannot grow file path: already pointing to a file"), "/data/devices.txt",
		},
	}

	for _, test := range tests {
		fp, _ := NewFilePath(test.input...)

		err := fp.Grow(test.growth...)
		if !reflect.DeepEqual(err, test.error) {
			t.Fatalf("error mismath. expected: %v. got: %v", test.error, err)
		}

		if !reflect.DeepEqual(fp.Path(), test.path) {
			t.Fatalf("path mismatch. expected: %v. got: %v", test.path, fp.Path())
		}
	}
}

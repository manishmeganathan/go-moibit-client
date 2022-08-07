package moibit

import (
	"fmt"
	"strings"
)

// FilePath represents the path to a file/directory in MOIBit.
type FilePath struct {
	elements  []string
	extension string
}

// Root returns a FilePath object that points to file system root i.e "/"
func Root() FilePath {
	return FilePath{}
}

// NewFilePath generates a new FilePath from a given variadic set of path elements.
// Each path element represents one level on the file system as a directory
// and can contain any character except for periods (.) and slashes (/).
//
// The last element may contain a single period to represent a file with a name and an extension like
// "file.json" or "hello.txt". If the last element has no period, the path is constructed for a directory
//
// An error is returned if any path element is invalid by not satisfying the above conditions.
// If no elements are given, the returned FilePath references the root directory.
func NewFilePath(elements ...string) (FilePath, error) {
	// Create a FilePath with initialized elements slice
	fp := FilePath{make([]string, 0, len(elements)), ""}

	// Iterate through the elements
	for idx, element := range elements {
		// If the element contains slash (/), throw an error
		if strings.Contains(element, "/") {
			return Root(), fmt.Errorf("failed to construct filepath: element '%v' contains slash", idx)
		}

		// If the element contains period (.)
		if strings.Contains(element, ".") {
			// Check if element is the last one
			if idx == len(element)-1 {
				// Split the element along the period and check if it contains only 2 split elements
				split := strings.Split(element, ".")
				if len(split) != 2 {
					return Root(), fmt.Errorf("failed to construct filepath: last element contains multiple periods")
				}

				// Append the first split element into the fp elements (file name)
				fp.elements = append(fp.elements, split[0])
				// Set the file second split element as file extension
				fp.extension = split[1]
				continue
			}

			// Throw error for non-final element with period (.)
			return Root(), fmt.Errorf("failed to construct filepath: non-final element '%v' contains period", idx)
		}

		// Append the element into the fp elements
		fp.elements = append(fp.elements, element)
	}

	return fp, nil
}

func NewFilePathFromString(path string) (FilePath, error) {
	return FilePath{}, nil
}

func (fp FilePath) Path() string {
	return ""
}

func (fp FilePath) IsRoot() bool {
	return fp.IsDirectory() && len(fp.elements) == 0
}

func (fp FilePath) IsDirectory() bool {
	return fp.extension == ""
}

func (fp FilePath) IsFile() bool {
	return fp.extension != ""
}

func (fp *FilePath) Grow(elements ...string) error {
	return nil
}

func (fp FilePath) Parent() FilePath {
	return FilePath{}
}

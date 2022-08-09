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
// Each path element is a string that may contain multiple filesystem path levels.
// If a path element contains a slash (/), it will be split and cleaned.
//
// The last element may contain a single period to represent a file with a name and an extension like
// "file.json" or "hello.txt". If the last element has no period, the path is assumed to be a directory.
// If an element apart from the last element contains a period, function returns an error.
//
// An error is returned if any path element is invalid by not satisfying the above conditions.
// If no elements are given, the returned FilePath references the root directory.
func NewFilePath(elements ...string) (FilePath, error) {
	// Create a FilePath with initialized elements slice
	fp := FilePath{make([]string, 0, len(elements)), ""}

	// Iterate through the elements
	for idx, element := range elements {
		// If the element is last one -> allow for possibility of a file extension
		if idx == len(elements)-1 {
			// If a period exists in the element, then a file extension might exist
			if strings.Contains(element, ".") {
				// Split the element along the period, but if there is more than one period in
				// the element, then it is not possible to determine the file extension correctly
				split := strings.Split(element, ".")
				if len(split) != 2 {
					return FilePath{}, fmt.Errorf("failed to construct filepath: multiple periods in final element")
				}

				// Prune slashes from the split elements (filename and extension)
				prunedFilename := cleanPath(split[0])
				prunedExtension := cleanPath(split[1])

				// If number of slash-pruned extension elements is more than 1, it means that there
				// is a slash after the period (in the extension) -> malformed file extension
				if len(prunedExtension) != 1 {
					return FilePath{}, fmt.Errorf("failed to construct filepath: slash detected after period or missing extension")
				}

				// Set filepath extension and append pruned elements
				fp.extension = prunedExtension[0]
				fp.elements = append(fp.elements, prunedFilename...)
				continue
			}
		}

		// Clean slashes from the element and check if any
		// of the sub-elements contains periods -> throw error,
		// otherwise, append it to the filepath elements
		pruned := cleanPath(element)
		for _, p := range pruned {
			if strings.Contains(p, ".") {
				return FilePath{}, fmt.Errorf("failed to construct filepath: non-final element '%v' contains period", idx)
			}

			fp.elements = append(fp.elements, p)
		}
	}

	return fp, nil
}

// Path returns a string representing the full path of represented by the FilePath
func (fp FilePath) Path() string {
	// Initialize a string builder
	var s strings.Builder

	// Add a slash and then all the filepath elements separated by a /
	s.WriteString(fmt.Sprintf("/%v", strings.Join(fp.elements, "/")))
	if fp.IsFile() {
		// If the filepath points to a file, add its extension to the path
		s.WriteString(fmt.Sprintf(".%v", fp.extension))
	}

	// Return the string from the builder
	return s.String()
}

// String implements the Stringer interface for FilePath
// Returns the path representation of the FilePath.
func (fp FilePath) String() string {
	return fp.Path()
}

// IsRoot returns whether the FilePath points to the "/" directory
func (fp FilePath) IsRoot() bool {
	return fp.IsDirectory() && len(fp.elements) == 0
}

// IsDirectory returns whether the FilePath points to a directory
func (fp FilePath) IsDirectory() bool {
	return fp.extension == ""
}

// IsFile returns whether the FilePath points to a file
func (fp FilePath) IsFile() bool {
	return fp.extension != ""
}

// Grow accepts a variadic set of path elements to grow the FilePath with.
// Returns an error if the elements are invalid or if the FilePath points to a file.
func (fp *FilePath) Grow(elements ...string) error {
	// FilePath cannot be grown if it is a file
	if fp.IsFile() {
		return fmt.Errorf("cannot grow file path: already pointing to a file")
	}

	// Append the given elements into the filepath elements
	elems := append(fp.elements, elements...)
	// Create a new FilePath from the full set of elements
	newfp, err := NewFilePath(elems...)
	if err != nil {
		return fmt.Errorf("cannot grow file path: bad elements: %w", err)
	}

	// Set the new FilePath to the method caller
	*fp = newfp
	return nil
}

// Parent returns a FilePath that points to the parent directory of the FilePath.
// If the FilePath is points to the root, the returned FilePath also points to the root.
func (fp FilePath) Parent() FilePath {
	// If the filepath points to a directory and contains less than 1 path element -> Return Root()
	if fp.IsDirectory() && len(fp.elements) <= 1 {
		return Root()
	}

	// Return a new filepath with one element popped from the end, with no extension
	return FilePath{fp.elements[:len(fp.elements)-1], ""}
}

// cleanPath is utility function that accepts a string path and returns
// a slice of strings which are free of slashes. Empty path elements will
// be discarded and an already clean path will be returned within the slice
func cleanPath(path string) []string {
	// If the path contains slashes, they need to be pruned
	if strings.Contains(path, "/") {
		// Initialize a slice of strings
		elements := make([]string, 0)

		// Split along the slashes and iterate through the split elements
		split := strings.Split(path, "/")
		for _, splitElement := range split {
			// Clean the split element and append them to elements
			pruned := cleanPath(splitElement)
			elements = append(elements, pruned...)
		}

		// Return the collected elements
		return elements

	} else {
		// If the path is an empty string, return an empty slice
		if path == "" {
			return nil
		}

		// Wrap path in a slice and return -> needs no cleaning
		return []string{path}
	}
}

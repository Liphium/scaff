package scaffui

import "strings"

func NewError(node Node, err error) *Error {
	if err == nil {
		return nil
	}

	if cerr, ok := err.(*Error); ok {
		if cerr == nil {
			return nil
		}

		cerr.add(node)
		return cerr
	}

	return &Error{
		path: []string{node.ID()},
		err:  err,
	}
}

var _ error = &Error{}

// An error that actually traces the path of all the nodes hit by the error (for easier error readability)
type Error struct {
	path []string
	err  error
}

func (e *Error) add(node Node) {
	e.path = append([]string{node.ID()}, e.path...)
}

func (e *Error) Error() string {
	if e == nil {
		return "<nil>"
	}

	formattedPath := strings.Join(e.path, " -> ")
	if e.err == nil {
		return formattedPath + ": <nil error>"
	}

	return formattedPath + ": " + e.err.Error()
}

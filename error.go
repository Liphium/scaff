package scaff

import "strings"

// An interface that attaches a simple string identifer to
type Identifiable interface {
	ID() string
}

// Create a new error that actually traces the path of different identifiable objects. This makes it visible where an error happened for more easy debugging.
//
// To use this, just wrap all of the times you return an error with this function.
func NewTracedError(identifiable Identifiable, err error) *TracedError {
	if err == nil {
		return nil
	}

	if cerr, ok := err.(*TracedError); ok {
		if cerr == nil {
			return nil
		}

		cerr.add(identifiable)
		return cerr
	}

	return &TracedError{
		path: []string{identifiable.ID()},
		err:  err,
	}
}

var _ error = &TracedError{}

// An error that actually traces the path of all the nodes hit by the error (for easier error readability)
type TracedError struct {
	path []string
	err  error
}

func (e *TracedError) add(identifiable Identifiable) {
	e.path = append([]string{identifiable.ID()}, e.path...)
}

// Get the actual error that happened. This will append all collected identifiables to the path of the error for debugging.
//
// If the error is nil, this will just return <nil> to prevent crashes.
func (e *TracedError) Error() string {
	if e == nil {
		return "<nil>"
	}

	formattedPath := strings.Join(e.path, " -> ")
	if e.err == nil {
		return formattedPath + ": <nil error>"
	}

	return formattedPath + ": " + e.err.Error()
}

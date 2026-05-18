package scaff

import "strings"

// An interface that attaches a simple string identifer to
type Identifiable interface {
	ID() string
}

// Create a new error that actually traces the path of different identifiable objects. This makes it visible where an error happened for more easy debugging.
//
// To use this, just wrap all of the times you return an error with this function.
func NewTracedError(identifiable Identifiable, err error) TracedError {
	if cerr, ok := err.(errorTrace); ok {

		cerr.path = append([]string{identifiable.ID()}, cerr.path...)
		return cerr
	}

	return errorTrace{
		path: []string{identifiable.ID()},
		err:  err,
	}
}

type TracedError interface {
	error
	diufhdsufhuidshuif()
}

var _ TracedError = errorTrace{}

// An error that actually traces the path of all the nodes hit by the error (for easier error readability)
type errorTrace struct {
	path []string
	err  error
}

func (et errorTrace) diufhdsufhuidshuif() {}

// Get the actual error that happened. This will append all collected identifiables to the path of the error for debugging.
//
// If the error is nil, this will just return <nil> to prevent crashes.
func (e errorTrace) Error() string {
	formattedPath := strings.Join(e.path, " -> ")
	if e.err == nil {
		return formattedPath + ": <nil error>"
	}

	return formattedPath + ": " + e.err.Error()
}

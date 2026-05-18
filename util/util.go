// Scaff's utility package. Contains a few useful things like loggers, etc.
package sutil

func Ptr[T any](value T) *T {
	return &value
}

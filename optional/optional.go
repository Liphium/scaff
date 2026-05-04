package optional

type O[T any] struct {
	value    T
	hasValue bool
}

func With[T any](value T) O[T] {
	return O[T]{value: value, hasValue: true}
}

func None[T any]() O[T] {
	return O[T]{
		hasValue: false,
	}
}

func (o *O[T]) SetValue(value T) {
	o.value = value
	o.hasValue = true
}

func (o *O[T]) Clear() {
	var val T
	o.value = val
	o.hasValue = false
}

func (o O[T]) Or(other T) T {
	if o.hasValue {
		return o.value
	}
	return other
}

func (o O[T]) HasValue() bool {
	return o.hasValue
}

func (o O[T]) Value() (T, bool) {
	return o.value, o.hasValue
}

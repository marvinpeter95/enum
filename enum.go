package enum

import "fmt"

// Enum is the interface that all generated enums implement. It provides methods for
// string representation, text marshaling, type identification, and validation.
type Enum[E Enum[E]] interface {
	fmt.Stringer
	fmt.GoStringer
	comparable

	// Type returns the type name of the enum in order to make enums usable as pflags.
	Type() string

	// EnumType returns the type name of the enum in order to identify them as enums.
	EnumType() string

	// IsValid provides a quick way to determine if the typed value is part of the allowed enumerated values.
	IsValid() bool

	// Validate checks if the typed value is part of the allowed enumerated values. If not it returns an error.
	Validate() error

	// EnumValues identifies the enum as an enum by returning a list of the possible values.
	EnumValues() []E
}

// EnumPointer is the interface that all generated enums implement. It extends Enum with additional methods
// that require modification of the enum value, such as text unmarshaling and setting the value from a string.
type EnumPointer[E Enum[E]] interface {
	*E

	// Set sets the value of the enum from a string. This is required to make enums usable as pflags.
	// It returns an error if the provided string is not a valid value for the enum.
	Set(string) error

	// SetValue sets the value of the enum from another value of the same type. This is a more efficient
	// way to set the value of the enum, as it avoids the need to parse a string. This is not required
	// for pflags, but it can be useful in other contexts.
	// It returns an error if the provided value is not a valid value for the enum.
	SetValue(E) error
}

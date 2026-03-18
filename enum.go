package enum

import (
	"fmt"
)

type BaseType interface {
	~string | ~int
}

// Enum is the interface that all generated enums implement. It provides methods for
// string representation, text marshaling, type identification, and validation.
type Enum interface {
	fmt.Stringer
	fmt.GoStringer

	// Type returns the type name of the enum in order to make enums usable as pflags.
	Type() string

	// EnumType returns the type name of the enum in order to identify them as enums.
	EnumType() string

	// IsValid provides a quick way to determine if the typed value is part of the allowed enumerated values.
	IsValid() bool

	// Validate checks if the typed value is  part of the allowed enumerated values. If not it returns an error.
	Validate() error
}

// EnumPointer is the interface that all generated enums implement. It extends Enum with additional methods
// that require modification of the enum value, such as text unmarshaling and setting the value from a string.
type EnumPointer interface {
	Enum

	// Set sets the value of the enum from a string. This is required to make enums usable as pflags.
	Set(string) error
}

// EnumOf is the interface that all generated enums implement. It extends Enum with an additional method
// that returns a list of the possible values for the enum.
type EnumOf[T BaseType] interface {
	Enum

	// EnumValues identifies the enum as an enum by returning a list of the possible values.
	EnumValues() []T
}

// EnumOfPointer is the interface that all generated enums implement. It extends EnumPointer with an additional method
// that returns a list of the possible values for the enum.
type EnumOfPointer[T BaseType] interface {
	EnumPointer

	// EnumValues identifies the enum as an enum by returning a list of the possible values.
	EnumValues() []T
}

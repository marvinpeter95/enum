# enum

Generate common methods for Golang enums.

## Usage

### Installation

```bash
go install github.com/marvinpeter95/enum/cmd/enum@latest
```

### Command Line

```bash
enum -types Type1,Type2 my_file.go
```

### Go Generate

```go
//go:generate enum -types Type1,Type2
package mypackage

type Type1 string

type Type2 int

// ...
```

## Generated Code

### Variables

```go
// ErrInvalid[Type] is returned when an invalid value is provided for [Type].
var ErrInvalid[Type] = errors.New("invalid [Type]")
```

### Functions

```go
// MustParse[Type] converts a string to a [Type], and panics if is not valid.
func MustParse[Type](string) [Type]

// Parse[Type] attempts to convert a string to a [Type].
func Parse[Type](string) ([Type], error)
```

### Methods

```go
// EnumValues identifies [Type] as an enum by returning a list of the possible values.
func ([Type]) EnumValues() [][Type]

// Type returns the type name of [Type].
func ([Type]) Type() string

// EnumType returns the type name of [Type].
func ([Type]) EnumType() string

// Set sets the value of the enum.
func (*[Type]) Set(v string) error

// IsValid provides a quick way to determine if the typed value is
// part of the allowed enumerated values.
func ([Type]) IsValid() bool

// Validate checks if the typed value is  part of the allowed enumerated
// values. If not it returns an error.
func ([Type]) Validate() error

// String implements the Stringer interface.
func ([Type]) String() string

// GoString implements the GoStringer interface.
func ([Type]) GoString() string

// MarshalText implements the text marshaller method.
func ([Type]) MarshalText() ([]byte, error)

// UnmarshalText implements the text unmarshaller method.
func (*[Type]) UnmarshalText([]byte) error
```

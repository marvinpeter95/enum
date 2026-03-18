package parser

import (
	"fmt"
	"strings"
)

// Enum represents the base type of an enum as either an int or a string.
type EnumType string

const (
	EnumTypeInvalid EnumType = "" // Use this to represent an invalid or uninitialized enum type
	EnumTypeInt     EnumType = "int"
	EnumTypeString  EnumType = "string"
)

// Enum represents a parsed enum type with its name, base type, values, and aliases.
type Enum struct {
	Name     string              // The name of the enum type
	Type     EnumType            // The base type of the enum (int or string)
	Values   []EnumValue         // The list of values defined for this enum
	Aliases  []EnumAlias         // The list of aliases defined for this enum
	valueSet map[string]struct{} // Used to ensure unique values
	aliasSet map[string]struct{} // Used to ensure unique aliases
}

// String returns a string representation of the Enum, including its name, type, values, and aliases.
func (e *Enum) String() string {
	sb := new(strings.Builder)
	fmt.Fprintf(sb, "%s(%s)\n", e.Name, e.Type)

	for _, value := range e.Values {
		fmt.Fprintf(sb, "  | %s\n", value.String())
	}

	for _, alias := range e.Aliases {
		fmt.Fprintf(sb, "  | %s\n", alias.String())
	}

	return strings.TrimSpace(sb.String())
}

// EnumValue represents a single value of an enum, including its name and the corresponding value.
type EnumValue struct {
	Name  string
	Value string
}

// String returns a string representation of the EnumValue, showing its name and value.
func (v EnumValue) String() string {
	return fmt.Sprintf("%s(%s)", v.Name, v.Value)
}

// EnumAlias represents an alias for an enum value, including its name and the corresponding EnumValue.
type EnumAlias struct {
	Name  string
	Value *EnumValue
}

// String returns a string representation of the EnumAlias, showing its name and the value it aliases.
func (a EnumAlias) String() string {
	return fmt.Sprintf("%s = %s", a.Name, a.Value.String())
}

// newEnum creates a new Enum with the given type name and initializes its fields.
func newEnum(name string) *Enum {
	return &Enum{
		Name:     name,
		Type:     EnumTypeInvalid,
		Values:   []EnumValue{},
		Aliases:  []EnumAlias{},
		valueSet: make(map[string]struct{}),
		aliasSet: make(map[string]struct{}),
	}
}

// AddValue adds a new value to the Enum if it doesn't already exist, ensuring that each value is unique.
func (e *Enum) AddValue(name string, value string) {
	if e.valueSet == nil {
		e.valueSet = make(map[string]struct{})
	}

	if _, exists := e.valueSet[name]; !exists {
		e.Values = append(e.Values, EnumValue{Name: name, Value: value})
		e.valueSet[name] = struct{}{}
	}
}

// AddAlias adds a new alias to the Enum if it doesn't already exist, ensuring that each alias is unique
// and corresponds to an existing value.
func (e *Enum) AddAlias(name string, alias string) bool {
	if e.aliasSet == nil {
		e.aliasSet = make(map[string]struct{})
	}
	if _, exists := e.aliasSet[name]; !exists {
		for i := range e.Values {
			if v := &e.Values[i]; v.Name == alias {
				e.Aliases = append(e.Aliases, EnumAlias{Name: name, Value: v})
				e.aliasSet[name] = struct{}{}
				return true
			}
		}

		return false
	}

	return true
}

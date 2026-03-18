package parser

import "strconv"

// IotaState keeps track of the current value of iota and the type it is associated with.
// It is used to generate the correct values for constants that use iota in Go.
type IotaState struct {
	Counter int
	Type    string
}

// NextValue returns the next value of iota as a string. It increments the counter and returns its value as a string.
func (i *IotaState) NextValue() string {
	i.Counter++
	return strconv.Itoa(i.Counter)
}

// Reset resets the iota counter to -1 and sets the type for the current iota sequence.
// This is called when a new iota sequence starts, typically when a new type is defined.
func (i *IotaState) Reset(typeName string) {
	i.Counter = -1
	i.Type = typeName
}

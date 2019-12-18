package cli

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

// IntSlice wraps []int to satisfy flag.Value
type IntSlice struct {
	slice      []int
	hasBeenSet bool
}

// NewIntSlice makes an *IntSlice with default values
func NewIntSlice(defaults ...int) *IntSlice {
	return &IntSlice{slice: append([]int{}, defaults...)}
}

// SetInt directly adds an integer to the list of values
func (i *IntSlice) SetInt(value int) {
	if !i.hasBeenSet {
		i.slice = []int{}
		i.hasBeenSet = true
	}

	i.slice = append(i.slice, value)
}

// Set parses the value into an integer and appends it to the list of values
func (i *IntSlice) Set(value string) error {
	if !i.hasBeenSet {
		i.slice = []int{}
		i.hasBeenSet = true
	}

	if strings.HasPrefix(value, slPfx) {
		_ = json.Unmarshal([]byte(strings.Replace(value, slPfx, "", 1)), &i.slice)
		i.hasBeenSet = true
		return nil
	}

	tmp, err := strconv.ParseInt(value, 0, 64)
	if err != nil {
		return err
	}

	i.slice = append(i.slice, int(tmp))
	return nil
}

// String returns a readable representation of this value(for usage defaults)
func (i *IntSlice) String() string {
	return fmt.Sprintf("%#v", i.slice)
}

// Serialize allows IntSlice to fulfill Serializer
func (i *IntSlice) Serialize() string {
	jsonBytes, _ := json.Marshal(i.slice)
	return fmt.Sprintf("%s%s", slPfx, string(jsonBytes))
}

// Value returns the slice of ints set by this flag
func (i *IntSlice) Value() []int {
	return i.slice
}

// Get returns the slice of ints set by this flag
func (i *IntSlice) Get() interface{} {
	return *i
}

package validator

import "reflect"

// FieldLevel contais all the information and helper functions to validate a field
type FieldLevel interface {

	// returns the top level struct, if any
	Top() reflect.Value

	// returns the current fields parent strut, of any or
	// the comparison value if called 'VarWithValue'
	Parent() reflect.Value

	// returns current field for validation
	Field() reflect.Value

	// returns the field's name with the tag
	// name taking precedence over the fields actual name.
	FieldName() string

	// returns the struct field's name
	StructFieldName() string

	// returns param for validation against current field
	Param() string

	// ExtractType gets the actual underlying type of field value.
	ExtractType(field reflect.Value) (value reflect.Value, kind reflect.Kind, nullable bool)

	GetStructFieldOK() (reflect.Value, reflect.Kind, bool)

	GetStructFieldOKAdvanced(val reflect.Value, namespace string) (reflect.Value, reflect.Kind, bool)

	GetStructFieldOK2() (reflect.Value, reflect.Kind, bool, bool)

	GetStructFieldOKAdvanced2(val reflect.Value, namespace string) (reflect.Value, reflect.Kind, bool, bool)
}

var _ FieldLevel = new(validate)

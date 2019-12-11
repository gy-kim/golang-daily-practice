package validator

import (
	"bytes"
	"fmt"
	"reflect"
	"strings"

	ut "github.com/go-playground/universal-translator"
)

const (
	fieldErrMsg = "Key: '%s' Error:Field validation for '%s' failed on the '%s' tag"
)

// ValidationErrorsTranslations is the translation return type
type ValidationErrorsTranslations map[string]string

// InvalidValidationError describes an invalid argument passed to
// `Struct`, `StructExecpt`, `StructPartial` or `Field`
type InvalidValidationError struct {
	Type reflect.Type
}

func (e *InvalidValidationError) Error() string {
	if e.Type == nil {
		return "validator: (nil)"
	}
	return "validator: (nil " + e.Type.String() + ")"
}

// ValidationErrors is an array of FieldError's
type ValidationErrors []FieldError

func (ve ValidationErrors) Error() string {
	buff := bytes.NewBufferString("")

	var fe *fieldError

	for i := 0; i < len(ve); i++ {
		fe = ve[i].(*fieldError)
		buff.WriteString(fe.Error())
		buff.WriteString("\n")
	}

	return strings.TrimSpace(buff.String())
}

// Translate transaltes all of the ValidationErrors
func (ve ValidationErrors) Translate(ut ut.Translator) ValidationErrorsTranslations {
	trans := make(ValidationErrorsTranslations)
	var fe *fieldError

	for i := 0; i < len(ve); i++ {
		fe = ve[i].(*fieldError)

		trans[fe.ns] = fe.Translate(ut)
	}
	return trans
}

// FieldError contains all functions to get error details
type FieldError interface {
	Tag() string

	ActualTag() string

	Namespace() string

	StructNamespace() string

	Field() string

	StructField() string

	Value() interface{}

	Param() string

	Kind() reflect.Kind

	Type() reflect.Type

	Translate(ut ut.Translator) string
}

// compile time interface checks
var _ FieldError = new(fieldError)
var _ error = new(fieldError)

type fieldError struct {
	v              *Validate
	tag            string
	actualTag      string
	ns             string
	structNs       string
	fieldLen       uint8
	structfieldLen uint8
	value          interface{}
	param          string
	kind           reflect.Kind
	typ            reflect.Type
}

func (fe *fieldError) Tag() string {
	return fe.tag
}

func (fe *fieldError) ActualTag() string {
	return fe.actualTag
}

func (fe *fieldError) Namespace() string {
	return fe.ns
}

func (fe *fieldError) StructNamespace() string {
	return fe.structNs
}

func (fe *fieldError) Field() string {
	return fe.ns[len(fe.ns)-int(fe.fieldLen):]
}

func (fe *fieldError) StructField() string {
	return fe.structNs[len(fe.structNs)-int(fe.structfieldLen):]
}

func (fe *fieldError) Value() interface{} {
	return fe.value
}

func (fe *fieldError) Param() string {
	return fe.param
}

func (fe *fieldError) Kind() reflect.Kind {
	return fe.kind
}

func (fe *fieldError) Type() reflect.Type {
	return fe.typ
}

func (fe *fieldError) Error() string {
	return fmt.Sprintf(fieldErrMsg, fe.ns, fe.Field(), fe.tag)
}

func (fe *fieldError) Translate(ut ut.Translator) string {
	m, ok := fe.v.transTagFunc[ut]
	if !ok {
		return fe.Error()
	}

	fn, ok := m[fe.tag]
	if !ok {
		return fe.Error()
	}
	return fn(ut, fe)
}

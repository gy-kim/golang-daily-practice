package validator

import (
	"context"
	"reflect"
)

// StructLevelFunc accepts all values needed for struct level validation
type StructLevelFunc func(sl StructLevel)

// StructLevelFuncCtx accepts all values needed for struct level validation
type StructLevelFuncCtx func(ctx context.Context, sl StructLevel)

func wrapStrutLevelFunc(fn StructLevelFunc) StructLevelFuncCtx {
	return func(ctx context.Context, sl StructLevel) {
		fn(sl)
	}
}

// StructLevel contains all the information and helper functions
type StructLevel interface {
	Validator() *Validate

	Top() reflect.Value

	Parent() reflect.Value

	Current() reflect.Value

	ExtractType(field reflect.Value) (value reflect.Value, kind reflect.Kind, nullable bool)

	ReportError(field interface{}, fieldName, structFieldName string, tag, param string)

	ReportValidationErrors(relativeNamespace, relativeActualNamespace string, errs ValidationErrors)
}

var _ StructLevel = new(validate)

func (v *validate) Top() reflect.Value {
	return v.top
}

func (v *validate) Current() reflect.Value {
	return v.slCurrent
}

func (v *validate) Validator() *Validate {
	return v.v
}

func (v *validate) ExtractType(field reflect.Value) (reflect.Value, reflect.Kind, bool) {
	return v.extractTypeInternal(field, false)
}

func (v *validate) ReportError(field interface{}, fieldName, structFieldName, tag, param string) {
	fv, kind, _ := v.extractTypeInternal(reflect.ValueOf(field), false)

	if len(structFieldName) == 0 {
		structFieldName = fieldName
	}

	v.str1 = string(append(v.ns, fieldName...))

	if v.v.hasTagNameFunc || fieldName != structFieldName {
		v.str2 = string(append(v.actualNs, structFieldName...))
	} else {
		v.str2 = v.str1
	}

	if kind == reflect.Invalid {
		v.errs = append(v.errs,
			&fieldError{
				v:              v.v,
				tag:            tag,
				actualTag:      tag,
				ns:             v.str1,
				structNs:       v.str2,
				fieldLen:       uint8(len(fieldName)),
				structfieldLen: uint8(len(structFieldName)),
				param:          param,
				kind:           kind,
			},
		)
		return
	}

	v.errs = append(v.errs,
		&fieldError{
			v:              v.v,
			tag:            tag,
			actualTag:      tag,
			ns:             v.str1,
			structNs:       v.str2,
			fieldLen:       uint8(len(fieldName)),
			structfieldLen: uint8(len(structFieldName)),
			value:          fv.Interface(),
			param:          param,
			kind:           kind,
			typ:            fv.Type(),
		},
	)
}

func (v *validate) ReportValidationErrors(relativeNamespace, relativeStructNamespace string, errs ValidationErrors) {
	var err *FieldError
	for i := 0; i < len(errs); i++ {
		err = errs[i].(*fieldError)
		err.ns = string(append(append(v.ns, relativeNamespace...), err.ns...))
		err.structNS = string(append(append(v.actualNs, relativeStructNamespace...), err.structNs...))
		v.errs = append(v.errs, err)
	}
}

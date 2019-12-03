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

package validator

import (
	"reflect"
	"strconv"
	"strings"
)

func (v *validate) extractTypeInternal(current reflect.Value, nullable bool) (reflect.Value, reflect.Kind, bool) {
BEGIN:
	switch current.Kind() {
	case reflect.Ptr:
		nullable = true

		if current.IsNil() {
			return current, reflect.Ptr, nullable
		}
		current = current.Elem()
		goto BEGIN
	case reflect.Interface:
		nullable = true

		if current.IsNil() {
			return current, reflect.Interface, nullable
		}
		current = current.Elem()
		goto BEGIN
	case reflect.Invalid:
		return current, reflect.Invalid, nullable
	default:
		if v.v.hasCustomFuncs {
			if fn, ok := v.v.customFuncs[current.Type()]; ok {
				current = reflect.ValueOf(fn(current))
				goto BEGIN
			}
		}
		return current, current.Kind(), nullable
	}
}

func (v *validate) getStructFieldOKInternal(val reflect.Value, namespace string) (current reflect.Value, kind reflect.Kind, nullable bool, found bool) {
BEGIN:
	current, kind, nullable = v.ExtractType(val)
	if kind == reflect.Invalid {
		return
	}

	if namespace == "" {
		found = true
		return
	}

	switch kind {
	case reflect.Ptr, reflect.Interface:
		typ := current.Type()
		fld := namespace
		var ns string

		if typ != timeType {
			idx := strings.Index(namespace, namespaceSeparator)

			if idx != -1 {
				fld = namespace[:idx]
				ns = namespace[idx+1:]
			} else {
				ns = ""
			}
			bracketIdx := strings.Index(fld, leftBracket)
			if bracketIdx != -1 {
				fld = fld[:bracketIdx]

				ns = namespace[bracketIdx:]
			}

			val = current.FieldByName(fld)
			namespace = ns
			goto BEGIN
		}
	case reflect.Array, reflect.Slice:
		idx := strings.Index(namespace, leftBracket)
		idx2 := strings.Index(namespace, rightBracket)

		arrIdx, _ := strconv.Atoi(namespace[idx+1 : idx2])

		if arrIdx >= current.Len() {
			return
		}

		startIdx := idx2 + 1
		if startIdx < len(namespace) {
			if namespace[startIdx:startIdx+1] == namespaceSeparator {
				startIdx++
			}
		}

		val = current.Index(arrIdx)
		namespace = namespace[startIdx:]
		goto BEGIN
	}
}

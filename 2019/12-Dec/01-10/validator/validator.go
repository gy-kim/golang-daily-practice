package validator

import (
	"context"
	"reflect"
)

type validate struct {
	v              *Validate
	top            reflect.Value
	ns             []byte
	actualNs       []byte
	errs           ValidationErrors
	includeExclude map[string]struct{}
	ffn            FilterFunc
	slflParent     reflect.Value
	slCurrent      reflect.Value
	flField        reflect.Value
	cf             *cField
	ct             *cTag
	misc           []byte
	str1           string
	str2           string
	fldIsPointer   bool
	isPartial      bool
	hasExcludes    bool
}

func (v *validate) validateStruct(ctx context.Context, parent reflect.Value, current reflect.Value, typ reflect.Type, ns []byte, structNs []byte, ct *cTag) {
	cs, ok := v.v.structCache.Get(typ)
	if !ok {
		cs = v.v.extractStructCache(current, typ.Name())
	}

	if len(ns) == 0 && len(cs.name) != 0 {
		ns = append(ns, cs.name...)
		ns = append(ns, '.')

		structNs = append(structNs, cs.name...)
		structNs = append(structNs, '.')
	}

	if ct == nil || ct.typeof != typeStructOnly {
		var f *cField

		for i := 0; i < len(cs.fields); i++ {
			f = cs.fields[i]

			if v.isPartial {
				if v.ffn != nil {
					if v.ffn(append(structNs, f.name...)) {
						continue
					}
				} else {
					_, ok := v.includeExclude[string(append(structNs, f.name...))]
					if (ok && v.hasExcludes) || (!ok && !v.hasExcludes) {
						continue
					}
				}
			}
			v.traverseField(ctx, parent, current.Field(f.idx), ns, structNs, f, f.cTags)
		}
	}

	if cs.fn != nil {
		v.slflParent = parent
		v.slCurrent = current
		v.ns = ns
		v.actualNs = structNs

		cs.fn(ctx, v)
	}
}

// func (v *validate) traverseField(ctx context.Context, parent reflect.Value, current reflect.Value, ns []byte, structNs []byte, cf *cField, ct *cTag) {
// 	var typ reflect.Type
// 	var kind reflect.Kind
// }

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

// 	current, kind, v.fldIsPointer = v.extractTypeInternal(current, false)

// 	switch kind {
// 	case reflect.Ptr, reflect.Interface, reflect.Invalid:
// 		if ct == nil {
// 			return
// 		}
// 		if ct.typeof == typeOmitEmpty || ct.typeof == typeIsDefault {
// 			return
// 		}

// 		if ct.hasTag {
// 			if kind == reflect.Invalid {
// 				v.str1 = string(append(ns, cf.altName...))
// 				if v.v.hasTagNameFunc {
// 					v.str2 = string(append(structNs, cf.name...))
// 				} else {
// 					v.str2 = v.str1
// 				}
// 				v.errs = append(v.errs,
// 					&fieldError{
// 						v:              v.v,
// 						tag:            ct.aliasTag,
// 						actualTag:      ct.tag,
// 						ns:             v.str1,
// 						structNs:       v.str2,
// 						fieldLen:       uint8(len(cf.altName)),
// 						structfieldLen: uint8(len(cf.name)),
// 						param:          ct.param,
// 						kind:           kind,
// 					},
// 				)
// 				return
// 			}

// 			v.str1 = string(append(ns, cf.altName...))
// 			if v.v.hasTagNameFunc {
// 				v.str2 = string(append(structNs, cf.name...))
// 			} else {
// 				v.str2 = v.str1
// 			}
// 			if !ct.runValidationWhenNil {
// 				v.errs = append(v.errs,
// 					&fieldError{
// 						v:              v.v,
// 						tag:            ct.tag,
// 						actualTag:      ct.tag,
// 						ns:             v.str1,
// 						structNs:       v.str2,
// 						fieldLen:       uint8(len(cf.altName)),
// 						structfieldLen: uint8(len(cf.name)),
// 						value:          current.Interface(),
// 						param:          ct.param,
// 						kind:           kind,
// 						typ:            current.Type(),
// 					},
// 				)
// 				return
// 			}
// 		}
// 	case reflect.Struct:
// 		typ = current.Type()
// 		if typ != timeType {
// 			if ct != nil {
// 				if ct.typeof == typeStructOnly {
// 					goto CONTINUE
// 				} else if ct.typeof == typeIsDefault {
// 					v.slflParent = parent
// 					v.flField = current
// 					v.cf = cf
// 					v.ct = ct

// 					if !ct.fn(ctx, v) {
// 						v.str1 = string(append(ns, cf.altName...))

// 						if v.v.hasTagNameFunc {
// 							v.str2 = string(append(structNs, cf.name...))
// 						} else {
// 							v.str2 = v.str1
// 						}

// 						v.errs = append(v.errs,
// 							&fieldError{
// 								v:              v.v,
// 								tag:            ct.aliasTag,
// 								ns:             v.str1,
// 								structNs:       v.str2,
// 								fieldLen:       uint8(len(cf.altName)),
// 								structfieldLen: uint8(len(cf.name)),
// 								value:          current.Interface(),
// 								param:          ct.param,
// 								kind:           kind,
// 								typ:            typ,
// 							},
// 						)
// 						return
// 					}
// 				}
// 				ct = ct.next
// 			}
// 			if ct != nil && ct.typeof == typeNoStructLevel {
// 				return
// 			}
// 		CONTINUE:
// 			if len(cf.name) > 0 {
// 				ns = append(append(ns, cf.altName...), '.')
// 				structNs = append(append(structNs, cf.name...), '.')
// 			}
// 			v.validateStruct(ctx, current, current, typ, ns, structNs, ct)
// 			return
// 		}
// 	}
// 	if !ct.hasTag {
// 		return
// 	}
// 	typ = current.Type()

// OUTER:
// 	for {
// 		if ct == nil {
// 			return
// 		}
// 		switch ct.typeof {
// 		case typeOmitEmpty:
// 			v.slflParent = parent
// 			v.flField = current
// 			v.cf = cf
// 			v.ct = ct

// 			if !v.fldIsPointer && !hasValue(v) {
// 				return
// 			}
// 			ct = ct.next
// 			continue
// 		}
// 	}
// }

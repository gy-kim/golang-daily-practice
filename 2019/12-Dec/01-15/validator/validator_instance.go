package validator

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"sync"
	"time"

	ut "github.com/go-playground/universal-translator"
)

const (
	defaultTagName        = "validate"
	utf8HexComma          = "0x2C"
	utf8Pipe              = "0x7C"
	tagSeparator          = ","
	orSeparator           = "|"
	tagKeySeparator       = "="
	structOnlyTag         = "structonly"
	noStructLevelTag      = "nostructlevel"
	omitempty             = "omitempty"
	isdefault             = "isdefault"
	requiredWithoutAllTag = "required_without_all"
	requiredWithoutTag    = "required_without"
	requiredWithTag       = "required_with"
	requiredWithAllTag    = "required_with_all"
	skipValidationTag     = "-"
	diveTag               = "dive"
	keysTag               = "keys"
	endKeysTag            = "endkeys"
	requiredTag           = "required"
	namespaceSeparator    = "."
	leftBracket           = "["
	rightBracket          = "]"
	restrictedTagChars    = ".[],|=+()`~!@#$%^&*\\\"/?<>{}"
	restrictedAliasErr    = "Alias '%s' either contains restricted characters or is the same as a restricted tag needed for normal operation"
	restrictedTagErr      = "Tag '%s' either contains restricted characters or is the same as a restricted tag needed for normal operation"
)

var (
	timeType      = reflect.TypeOf(time.Time{})
	defaultCField = &cField{namesEqual: true}
)

// FilterFunc is the type used to filter fields using StuctFieltered(...) function.
type FilterFunc func(ns []byte) bool

// CustomTypeFunc allows for overriding or adding custom field type handler functions
// field = field value of the type to return a value to be validated
type CustomTypeFunc func(field reflect.Value) interface{}

// TagNameFunc allows for adding of custom tag name parser
type TagNameFunc func(field reflect.StructField) string

type internalValidationFuncWrapper struct {
	fn                 FuncCtx
	runValidationOnNil bool
}

// Validate contains the validator settings and cache
type Validate struct {
	tagName          string
	pool             *sync.Pool
	hasCustomFuncs   bool
	hasTagNameFunc   bool
	tagNameFunc      TagNameFunc
	structLevelFuncs map[reflect.Type]StructLevelFuncCtx
	customFuncs      map[reflect.Type]CustomTypeFunc
	aliases          map[string]string
	validations      map[string]internalValidationFuncWrapper
	transTagFunc     map[ut.Translator]map[string]TranslationFunc
	tagCache         *tagCache
	structCache      *structCache
}

func New() *Validate {
	tc := new(tagCache)
	tc.m.Store(make(map[string]*cTag))

	sc := new(structCache)
	sc.m.Store(make(map[reflect.Type]*cStruct))

	v := &Validate{
		tagName:     defaultTagName,
		aliases:     make(map[string]string, len(bakedInAliases)),
		validations: make(map[string]internalValidationFuncWrapper, len(bakedInValidators)),
		tagCache:    tc,
		structCache: sc,
	}

	for k, val := range bakedInAliases {
		v.RegisterAlias(k, val)
	}

	for k, val := range bakedInValidators {
		switch k {
		case requiredWithTag, requiredWithAllTag, requiredWithoutTag, requiredWithoutAllTag:
			_ = v.registerValidation(k, wrapFunc(val), true, true)
		default:
			_ = v.registerValidation(k, wrapFunc(val), true, false)
		}
	}

	v.pool = &sync.Pool{
		New: func() interface{} {
			return &validate{
				v:        v,
				ns:       make([]byte, 0, 64),
				actualNs: make([]byte, 0, 64),
				misc:     make([]byte, 32),
			}
		},
	}
	return v
}

func (v *Validate) SetTagName(name string) {
	v.tagName = name
}

func (v *Validate) RegisterTagNameFunc(fn TagNameFunc) {
	v.tagNameFunc = fn
	v.hasTagNameFunc = true
}

func (v *Validate) RegisterValidation(tag string, fn Func, callValidationEventIfNull ...bool) error {
	return v.RegisterValidationCtx(tag, wrapFunc(fn), callValidationEventIfNull...)
}

func (v *Validate) RegisterValidationCtx(tag string, fn FuncCtx, callValidationEventIfNull ...bool) error {
	var nilCheckable bool
	if len(callValidationEventIfNull) > 0 {
		nilCheckable = callValidationEventIfNull[0]
	}

	return v.registerValidation(tag, fn, false, nilCheckable)
}

func (v *Validate) registerValidation(tag string, fn FuncCtx, bakedIn bool, nilCheckable bool) error {
	if len(tag) == 0 {
		return errors.New("Function Key cannot be empty")
	}
	if fn == nil {
		return errors.New("Function cannot be empty")
	}

	_, ok := restrictedTags[tag]
	if !bakedIn && (ok || strings.ContainsAny(tag, restrictedTagChars)) {
		panic(fmt.Sprintf(restrictedTagErr, tag))
	}

	v.validations[tag] = internalValidationFuncWrapper{fn: fn, runValidationOnNil: nilCheckable}
	return nil
}

func (v *Validate) RegisterAlias(alias, tags string) {
	_, ok := restrictedTags[alias]

	if ok || strings.ContainsAny(alias, restrictedTagChars) {
		panic(fmt.Sprintf(restrictedAliasErr, alias))
	}

	v.aliases[alias] = tags
}

func (v *Validate) RegisterStructValidation(fn StructLevelFunc, types ...interface{}) {
	v.RegisterStructValidationCtx(wrapStructLevelFunc(fn), types...)
}

func (v *Validate) RegisterStructValidationCtx(fn StructLevelFuncCtx, types ...interface{}) {
	if v.structLevelFuncs == nil {
		v.structLevelFuncs = make(map[reflect.Type]StructLevelFuncCtx)
	}

	for _, t := range types {
		tv := reflect.ValueOf(t)
		if tv.Kind() == reflect.Ptr {
			t = reflect.Indirect(tv).Interface()
		}
		v.structLevelFuncs[reflect.TypeOf(t)] = fn
	}
}

func (v *Validate) RegisterCustomTypeFunc(fn CustomTypeFunc, types ...interface{}) {
	if v.customFuncs == nil {
		v.customFuncs = make(map[reflect.Type]CustomTypeFunc)
	}

	for _, t := range types {
		v.customFuncs[reflect.TypeOf(t)] = fn
	}
	v.hasCustomFuncs = true
}

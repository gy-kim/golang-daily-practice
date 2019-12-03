package validator

import (
	"reflect"
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
	namespeceSeparator    = "."
	leftBracket           = "["
	rightBracket          = "]"
	restrictedTagChars    = ".[],|=+()`~!@#$%^&*\\\"/?<>{}"
	restrictedAliasErr    = "Alias '%s' either contains restricted characters or is the same as a restricted tag needed for normal operation"
	restructedTagErr      = "Tag '%s' either contains restricted characters or is the same as a restricted tag needed for normal operation"
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

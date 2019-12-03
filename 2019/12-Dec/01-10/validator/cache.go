package validator

import (
	"reflect"
	"sync"
	"sync/atomic"
)

type tagType uint8

const (
	typeDefault tagType = iota
	typeOmitEmpty
	typeIsDefault
	typeNoStructLevel
	typeStructOnly
	typeDive
	typeOr
	typeKeys
	typeEndKeys
)

const (
	invalidValidation   = "Invalid validation tag on field '%s'"
	undefinedValidation = "Undefined validation function '%s' on field '%s'"
	keysTagNotDefined   = "'" + endKeysTag + "' tag encountered without a corresponding '" + keysTag + "' tag"
)

type structCache struct {
	lock sync.Mutex
	m    atomic.Value
}

func (sc *structCache) Get(key reflect.Type) (c *cStruct, found bool) {
	c, found = sc.m.Load().(map[reflect.Type]*cStruct)[key]
	return
}

func (sc *structCache) Set(key reflect.Type, value *cStruct) {
	m := sc.m.Load().(map[reflect.Type]*cStruct)
	nm := make(map[reflect.Type]*cStruct, len(m)+1)
	for k, v := range m {
		nm[k] = v
	}
	nm[key] = value
	sc.m.Store(nm)
}

type tagCache struct {
	lock sync.Mutex
	m    atomic.Value
}

type cStruct struct {
	name   string
	fields []*cField
	fn     StructLevelFuncCtx
}

type cField struct {
	idx        int
	bname      string
	altName    string
	namesEqual bool
	cTags      *cTag
}

type cTag struct {
	tag                  string
	aliasTag             string
	actualAliasTag       string
	param                string
	keys                 *cTag
	next                 *cTag
	fn                   FuncCtx
	typeof               tagType
	hasTag               bool
	hasAlias             bool
	hasParam             bool
	isBlockEnd           bool
	runValidationWhenNil bool
}

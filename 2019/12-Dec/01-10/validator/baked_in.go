package validator

import (
	"context"
	"strings"
	"sync"
)

// Func accepts a FieldLevel interface for all validation needs.
type Func func(fl FieldLevel) bool

// FuncCtx accepts a context.Context and FieldLevel interface for all validation needs.
type FuncCtx func(ctx context.Context, fl FieldLevel) bool

// wrapFunc wraps noramal Func makes it compatible with FuncCtx
func wrapFun(fn Func) FuncCtx {
	if fn == nil {
		return nil
	}
	return func(ctx context.Context, fl FieldLevel) bool {
		return fn(fl)
	}
}

var (
	restrictedTags = map[string]struct{}{
		diveTag:           {},
		keysTag:           {},
		endKeysTag:        {},
		structOnlyTag:     {},
		omitempty:         {},
		skipValidationTag: {},
		utf8HexComma:      {},
		noStructLevelTag:  {},
		requiredTag:       {},
		isdefault:         {},
	}

	bakedInAliases = map[string]string{
		"iscolor": "hexcolor|rgb|rgba|hsl|hsla",
	}

	bakedInValidators = map[string]Func{
		// "required":             hasValue,
		// "required_with":        requiredWith,
		// "required_with_all":    requiredWithAll,
		// "required_without":     requiredWithout,
		// "required_without_all": requiredWithoutAll,
		// "isdefault":            isDefault,
		// "len":                  hasLengthOf,
		// "min":                  hasMinOf,
		// "max":                  hasMaxOf,
		// "eq":                   isEq,
		// "ne":                   isNe,
		// "lt":                   isLt,
		// "lte":                  isLte,
		// "gt":                   isGt,
		// "gte":                  isGte,
		// "eqfield":              isEqField,
		// "eqcsfield":            isEqCrossStructField,
		// "necsfield":            isNeCrossStructField,
		// "gtcsfield":            isGtCrossStructField,
		// "gtecsfield":           isGteCrossStructField,
		// "ltcsfield":            isLtCrossStructField,
		// "ltecsfield":           isLteCrossStructField,
		// "nefield":              isNeField,
		// "gtefield":             isGteField,
		// "gtfield":              isGtField,
		// "ltefield":             isLteField,
		// "ltfield":              isLtField,
		// "fieldcontains":        fieldContains,
		// "fieldexcludes":        fieldExcludes,
		// "alpha":                isAlpha,
		// "alphanum":             isAlphanum,
		// "alphaunicode":         isAlphaUnicode,
		// "alphanumunicode":      isAlphanumUnicode,
		// "numeric":              isNumeric,
		// "number":               isNumber,
		// "hexadecimal":          isHexadecimal,
		// "hexcolor":             isHEXColor,
		// "rgb":                  isRGB,
		// "rgba":                 isRGBA,
		// "hsl":                  isHSL,
		// "hsla":                 isHSLA,
		// "email":                isEmail,
		// "url":                  isURL,
		// "uri":                  isURI,
		// "urn_rfc2141":          isUrnRFC2141, // RFC 2141
		// "file":                 isFile,
		// "base64":               isBase64,
		// "base64url":            isBase64URL,
		// "contains":             contains,
		// "containsany":          containsAny,
		// "containsrune":         containsRune,
		// "excludes":             excludes,
		// "excludesall":          excludesAll,
		// "excludesrune":         excludesRune,
		// "startswith":           startsWith,
		// "endswith":             endsWith,
		// "isbn":                 isISBN,
		// "isbn10":               isISBN10,
		// "isbn13":               isISBN13,
		// "eth_addr":             isEthereumAddress,
		// "btc_addr":             isBitcoinAddress,
		// "btc_addr_bech32":      isBitcoinBech32Address,
		// "uuid":                 isUUID,
		// "uuid3":                isUUID3,
		// "uuid4":                isUUID4,
		// "uuid5":                isUUID5,
		// "uuid_rfc4122":         isUUIDRFC4122,
		// "uuid3_rfc4122":        isUUID3RFC4122,
		// "uuid4_rfc4122":        isUUID4RFC4122,
		// "uuid5_rfc4122":        isUUID5RFC4122,
		// "ascii":                isASCII,
		// "printascii":           isPrintableASCII,
		// "multibyte":            hasMultiByteCharacter,
		// "datauri":              isDataURI,
		// "latitude":             isLatitude,
		// "longitude":            isLongitude,
		// "ssn":                  isSSN,
		// "ipv4":                 isIPv4,
		// "ipv6":                 isIPv6,
		// "ip":                   isIP,
		// "cidrv4":               isCIDRv4,
		// "cidrv6":               isCIDRv6,
		// "cidr":                 isCIDR,
		// "tcp4_addr":            isTCP4AddrResolvable,
		// "tcp6_addr":            isTCP6AddrResolvable,
		// "tcp_addr":             isTCPAddrResolvable,
		// "udp4_addr":            isUDP4AddrResolvable,
		// "udp6_addr":            isUDP6AddrResolvable,
		// "udp_addr":             isUDPAddrResolvable,
		// "ip4_addr":             isIP4AddrResolvable,
		// "ip6_addr":             isIP6AddrResolvable,
		// "ip_addr":              isIPAddrResolvable,
		// "unix_addr":            isUnixAddrResolvable,
		// "mac":                  isMAC,
		// "hostname":             isHostnameRFC952,  // RFC 952
		// "hostname_rfc1123":     isHostnameRFC1123, // RFC 1123
		// "fqdn":                 isFQDN,
		// "unique":               isUnique,
		// "oneof":                isOneOf,
		// "html":                 isHTML,
		// "html_encoded":         isHTMLEncoded,
		// "url_encoded":          isURLEncoded,
		// "dir":                  isDir,
	}
)

var oneofValsCache = map[string][]string{}
var oneofValsCacheRWLock = sync.RWMutex{}

func parseOneOfParam2(s string) []string {
	oneofValsCacheRWLock.RLock()
	vals, ok := oneofValsCache[s]
	oneofValsCacheRWLock.RUnlock()
	if !ok {
		oneofValsCacheRWLock.Lock()
		vals = strings.Fields(s)
		oneofValsCache[s] = vals
		oneofValsCacheRWLock.Unlock()
	}
	return vals
}

func isURLEncoded(fl FieldLevel) bool {
	return uRLEncodedRegex.MatchString(fl.Field().String())
}

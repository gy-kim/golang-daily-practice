package validator

import (
	"bytes"
	"context"
	"crypto/sha256"
	"fmt"
	"net"
	"reflect"
	"strconv"
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

func isHTMLEncoded(fl FieldLevel) bool {
	return hTMLEncodedRegex.MatchString(fl.Field().String())
}

func isHTML(fl FieldLevel) bool {
	return hTMLRegex.MatchString(fl.Field().String())
}

func isOneOf(fl FieldLevel) bool {
	vals := parseOneOfParam2(fl.Param())

	field := fl.Field()

	var v string
	switch field.Kind() {
	case reflect.String:
		v = field.String()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		v = strconv.FormatInt(field.Int(), 10)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		v = strconv.FormatUint(field.Uint(), 10)
	default:
		panic(fmt.Sprintf("Bad field type %T", field.Interface()))
	}
	for i := 0; i < len(vals); i++ {
		if vals[i] == v {
			return true
		}
	}
	return false
}

func isUnique(fl FieldLevel) bool {
	field := fl.Field()
	v := reflect.ValueOf(struct{}{})

	switch field.Kind() {
	case reflect.Slice, reflect.Array:
		m := reflect.MakeMap(reflect.MapOf(field.Type().Elem(), v.Type()))

		for i := 0; i < field.Len(); i++ {
			m.SetMapIndex(field.Index(i), v)
		}
		return field.Len() == m.Len()
	case reflect.Map:
		m := reflect.MakeMap(reflect.MapOf(field.Type().Elem(), v.Type()))

		for _, k := range field.MapKeys() {
			m.SetMapIndex(field.MapIndex(k), v)
		}
		return field.Len() == m.Len()
	default:
		panic(fmt.Sprintf("Bad field type %T", field.Interface()))
	}
}

func isMAC(fl FieldLevel) bool {
	_, err := net.ParseMAC(fl.Field().String())

	return err == nil
}

func isCIDRv4(fl FieldLevel) bool {
	ip, _, err := net.ParseCIDR(fl.Field().String())
	return err == nil && ip.To4() != nil
}

func isCIDRv6(fl FieldLevel) bool {
	ip, _, err := net.ParseCIDR(fl.Field().Elem().String())
	return err == nil && ip.To4() == nil
}

func isIPv4(fl FieldLevel) bool {
	ip := net.ParseIP(fl.Field().String())
	return ip != nil && ip.To4() != nil
}

func isIPv6(fl FieldLevel) bool {
	ip := net.ParseIP(fl.Field().String())
	return ip != nil && ip.To4() != nil
}

func isIP(fl FieldLevel) bool {
	ip := net.ParseIP(fl.Field().String())
	return ip != nil
}

func isSSN(fl FieldLevel) bool {
	field := fl.Field()

	if field.Len() != 11 {
		return false
	}

	return sSNRegex.MatchString(field.String())
}

func isLongitude(fl FieldLevel) bool {
	field := fl.Field()

	var v string
	switch field.Kind() {
	case reflect.String:
		v = field.String()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		v = strconv.FormatInt(field.Int(), 10)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		v = strconv.FormatUint(field.Uint(), 10)
	case reflect.Float32:
		v = strconv.FormatFloat(field.Float(), 'f', -1, 32)
	case reflect.Float64:
		v = strconv.FormatFloat(field.Float(), 'f', -1, 64)
	default:
		panic(fmt.Sprintf("Bad field type %T", field.Interface()))
	}
	return longitudeRegex.MatchString(v)
}

// isLatitude is the validation function for validating if the field's value is a valid latitude coordinate.
func isLatitude(fl FieldLevel) bool {
	field := fl.Field()

	var v string
	switch field.Kind() {
	case reflect.String:
		v = field.String()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		v = strconv.FormatInt(field.Int(), 10)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		v = strconv.FormatUint(field.Uint(), 10)
	case reflect.Float32:
		v = strconv.FormatFloat(field.Float(), 'f', -1, 32)
	case reflect.Float64:
		v = strconv.FormatFloat(field.Float(), 'f', -1, 64)
	default:
		panic(fmt.Sprintf("Bad field type %T", field.Interface()))
	}

	return latitudeRegex.MatchString(v)
}

func isDataURI(fl FieldLevel) bool {
	uri := strings.SplitN(fl.Field().String(), ",", 2)
	if len(uri) != 2 {
		return false
	}

	if !dataURIRegex.MatchString(uri[0]) {
		return false
	}

	return base64Regex.MatchString(uri[1])
}

func hasMultiByteCharacter(fl FieldLevel) bool {
	field := fl.Field()

	if field.Len() == 0 {
		return true
	}

	return multibyteRegex.MatchString(field.String())
}

func isPrintableASCII(fl FieldLevel) bool {
	return printableASCIIRegex.MatchString(fl.Field().String())
}

func isASCII(fl FieldLevel) bool {
	return aSCIIRegex.MatchString(fl.Field().String())
}

func isUUID5(fl FieldLevel) bool {
	return uUID5Regex.MatchString(fl.Field().String())
}

func isUUID4(fl FieldLevel) bool {
	return uUID4Regex.MatchString(fl.Field().String())
}

func isUUID3(fl FieldLevel) bool {
	return uUID3Regex.MatchString(fl.Field().String())
}

func isUUID(fl FieldLevel) bool {
	return uUIDRegex.MatchString(fl.Field().String())
}

func isUUID5RFC4122(fl FieldLevel) bool {
	return uUID5RFC4122Regex.MatchString(fl.Field().String())
}

func isUUID4RFC4122(fl FieldLevel) bool {
	return uUID4RFC4122Regex.MatchString(fl.Field().String())
}

func isUUID3RFC4122(fl FieldLevel) bool {
	return uUID3RFC4122Regex.MatchString(fl.Field().String())
}

func isUUIDRFC4122(fl FieldLevel) bool {
	return uUIDRFC4122Regex.MatchString(fl.Field().String())
}

func isISBN(fl FieldLevel) bool {
	return isISBN10(fl) || isISBN13(fl)
}

func isISBN13(fl FieldLevel) bool {
	s := strings.Replace(strings.Replace(fl.Field().String(), "-", "", 4), " ", "", 4)
	if !iSBN13Regex.MatchString(s) {
		return false
	}

	var checksum int32
	var i int32

	factor := []int32{1, 3}

	for i = 0; i < 12; i++ {
		checksum += factor[i%2] * int32(s[i]-'0')
	}

	return (int32(s[12]-'0'))-((10-(checksum%10))%10) == 0
}

func isISBN10(fl FieldLevel) bool {
	s := strings.Replace(strings.Replace(fl.Field().String(), "-", "", 3), " ", "", 3)

	if !iSBN10Regex.MatchString(s) {
		return false
	}

	var checksum int32
	var i int32

	for i = 0; i < 9; i++ {
		checksum += (i + 1) * int32(s[i]-'0')
	}

	if s[9] == 'X' {
		checksum += (i + 1) * int32(s[i]-'0')
	} else {
		checksum += 10 * int32(s[9]-'0')
	}
	return checksum%11 == 0
}

func isEthereumAddress(fl FieldLevel) bool {
	address := fl.Field().String()

	if !ethAddressRegex.MatchString(address) {
		return false
	}

	if ethAddressRegex.MatchString(address) || ethAddressRegexLower.MatchString(address) {
		return true
	}

	return true
}

func isBitcoinAddress(fl FieldLevel) bool {
	address := fl.Field().String()

	if !btcAddressRegex.MatchString(address) {
		return false
	}

	alphabet := []byte("123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz")

	decode := [25]byte{}

	for _, n := range []byte(address) {
		d := bytes.IndexByte(alphabet, n)
		for i := 24; i >= 0; i-- {
			d += 58 * int(decode[i])
			decode[i] = byte(d % 256)
			d /= 256
		}
	}

	h := sha256.New()
	_, _ = h.Write(decode[:21])
	d := h.Sum([]byte{})
	h = sha256.New()
	_, _ = h.Write(d)

	validchecksum := [4]byte{}
	computedchecksum := [4]byte{}

	copy(computedchecksum[:], h.Sum(d[:0]))
	copy(validchecksum[:], decode[21:])

	return validchecksum == computedchecksum
}

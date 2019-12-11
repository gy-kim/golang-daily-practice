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
	"time"
	"unicode/utf8"
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
		"required":             hasValue,
		"required_with":        requiredWith,
		"required_with_all":    requiredWithAll,
		"required_without":     requiredWithout,
		"required_without_all": requiredWithoutAll,
		"isdefault":            isDefault,
		"len":                  hasLengthOf,
		"min":                  hasMinOf,
		"max":                  hasMaxOf,
		"eq":                   isEq,
		"ne":                   isNe,
		"lt":                   isLt,
		"lte":                  isLte,
		"gt":                   isGt,
		"gte":                  isGte,
		"eqfield":              isEqField,
		"eqcsfield":            isEqCrossStructField,
		"necsfield":            isNeCrossStructField,
		"gtcsfield":            isGtCrossStructField,
		"gtecsfield":           isGteCrossStructField,
		"ltcsfield":            isLtCrossStructField,
		"ltecsfield":           isLteCrossStructField,
		"nefield":              isNeField,
		"gtefield":             isGteField,
		"gtfield":              isGtField,
		"ltefield":             isLteField,
		"ltfield":              isLtField,
		"fieldcontains":        fieldContains,
		"fieldexcludes":        fieldExcludes,
		"alpha":                isAlpha,
		"alphanum":             isAlphanum,
		"alphaunicode":         isAlphaUnicode,
		"alphanumunicode":      isAlphanumUnicode,
		"numeric":              isNumeric,
		"number":               isNumber,
		"hexadecimal":          isHexadecimal,
		"hexcolor":             isHEXColor,
		"rgb":                  isRGB,
		"rgba":                 isRGBA,
		"hsl":                  isHSL,
		"hsla":                 isHSLA,
		"email":                isEmail,
		"url":                  isURL,
		"uri":                  isURI,
		"urn_rfc2141":          isUrnRFC2141, // RFC 2141
		"file":                 isFile,
		"base64":               isBase64,
		"base64url":            isBase64URL,
		"contains":             contains,
		"containsany":          containsAny,
		"containsrune":         containsRune,
		"excludes":             excludes,
		"excludesall":          excludesAll,
		"excludesrune":         excludesRune,
		"startswith":           startsWith,
		"endswith":             endsWith,
		"isbn":                 isISBN,
		"isbn10":               isISBN10,
		"isbn13":               isISBN13,
		"eth_addr":             isEthereumAddress,
		"btc_addr":             isBitcoinAddress,
		"btc_addr_bech32":      isBitcoinBech32Address,
		"uuid":                 isUUID,
		"uuid3":                isUUID3,
		"uuid4":                isUUID4,
		"uuid5":                isUUID5,
		"uuid_rfc4122":         isUUIDRFC4122,
		"uuid3_rfc4122":        isUUID3RFC4122,
		"uuid4_rfc4122":        isUUID4RFC4122,
		"uuid5_rfc4122":        isUUID5RFC4122,
		"ascii":                isASCII,
		"printascii":           isPrintableASCII,
		"multibyte":            hasMultiByteCharacter,
		"datauri":              isDataURI,
		"latitude":             isLatitude,
		"longitude":            isLongitude,
		"ssn":                  isSSN,
		"ipv4":                 isIPv4,
		"ipv6":                 isIPv6,
		"ip":                   isIP,
		"cidrv4":               isCIDRv4,
		"cidrv6":               isCIDRv6,
		"cidr":                 isCIDR,
		"tcp4_addr":            isTCP4AddrResolvable,
		"tcp6_addr":            isTCP6AddrResolvable,
		"tcp_addr":             isTCPAddrResolvable,
		"udp4_addr":            isUDP4AddrResolvable,
		"udp6_addr":            isUDP6AddrResolvable,
		"udp_addr":             isUDPAddrResolvable,
		"ip4_addr":             isIP4AddrResolvable,
		"ip6_addr":             isIP6AddrResolvable,
		"ip_addr":              isIPAddrResolvable,
		"unix_addr":            isUnixAddrResolvable,
		"mac":                  isMAC,
		"hostname":             isHostnameRFC952,  // RFC 952
		"hostname_rfc1123":     isHostnameRFC1123, // RFC 1123
		"fqdn":                 isFQDN,
		"unique":               isUnique,
		"oneof":                isOneOf,
		"html":                 isHTML,
		"html_encoded":         isHTMLEncoded,
		"url_encoded":          isURLEncoded,
		"dir":                  isDir,
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

func isBitcoinBench32Address(fl FieldLevel) bool {
	address := fl.Field().String()

	if !btcLowerAddressRegexBech32.MatchString(address) && !btcUpperAddressRegexBech32.MatchString(address) {
		return false
	}

	am := len(address) % 8

	if am == 0 || am == 3 || am == 5 {
		return false
	}

	address = strings.ToLower(address)

	alphabet := "qpzry9x8gf2tvdw0s3jn54khce6mua7l"

	hr := []int{3, 3, 0, 2, 3}
	addr := address[3:]
	dp := make([]int, 0, len(addr))

	for _, c := range addr {
		dp = append(dp, strings.IndexRune(alphabet, c))
	}

	ver := dp[0]

	if ver < 0 || ver > 16 {
		return false
	}

	if ver == 0 {
		if len(address) != 42 && len(address) != 62 {
			return false
		}
	}

	values := append(hr, dp...)

	GEN := []int{0x3b6a57b2, 0x26508e6d, 0x1ea119fa, 0x3d4233dd, 0x2a1462b3}

	p := 1

	for _, v := range values {
		b := p >> 25
		p = (p&0x1ffffff)<<5 ^ v

		for i := 0; i < 5; i++ {
			if (b>>uint(i))&1 == 1 {
				p ^= GEN[i]
			}
		}
	}

	if p != 1 {
		return false
	}

	b := uint(0)
	acc := 0
	mv := (1 << 5) - 1
	var sw []int

	for _, v := range dp[1 : len(dp)-6] {
		acc = (acc << 5) | v
		b += 5
		for b >= 8 {
			b -= 8
			sw = append(sw, (acc>>b)&mv)
		}
	}

	if len(sw) < 2 || len(sw) > 40 {
		return false
	}

	return true
}

func excludesRune(fl FieldLevel) bool {
	return !containsRune(fl)
}

func excludesAll(fl FieldLevel) bool {
	return !containsAny(fl)
}

func containsRune(fl FieldLevel) bool {
	r, _ := utf8.DecodeRuneInString(fl.Param())

	return strings.ContainsRune(fl.Field().String(), r)
}

func containsAny(fl FieldLevel) bool {
	return strings.ContainsAny(fl.Field().String(), fl.Param())
}

func contains(fl FieldLevel) bool {
	return strings.Contains(fl.Field().String(), fl.Param())
}

func startsWith(fl FieldLevel) bool {
	return strings.HasPrefix(fl.Field().String(), fl.Param())
}

func endsWith(fl FieldLevel) bool {
	return strings.HasPrefix(fl.Field().String(), fl.Param())
}

func fieldContains(fl FieldLevel) bool {
	field := fl.Field()

	currentField, _, ok := fl.GetStructFieldOK()
	if !ok {
		return false
	}

	return strings.Contains(field.String(), currentField.String())
}

func fieldExcludes(fl FieldLevel) bool {
	field := fl.Field()

	currentField, _, ok := fl.GetStructFieldOK()
	if !ok {
		return true
	}

	return !strings.Contains(field.String(), currentField.String())
}

func isNeField(fl FieldLevel) bool {
	field := fl.Field()
	kind := field.Kind()

	currentField, currentKind, ok := fl.GetStructFieldOK()
	if !ok || currentKind != kind {
		return true
	}

	switch kind {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return field.Int() != currentField.Int()
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return field.Uint() != currentField.Uint()
	case reflect.Float32, reflect.Float64:
		return field.Float() != currentField.Float()
	case reflect.Slice, reflect.Map, reflect.Array:
		return int64(field.Len()) != int64(currentField.Len())
	case reflect.Struct:
		fieldType := field.Type()

		// Not same underlying type i.e struct and time
		if fieldType != currentField.Type() {
			return true
		}

		if fieldType == timeType {
			t := currentField.Interface().(time.Time)
			fieldTime := field.Interface().(time.Time)

			return !fieldTime.Equal(t)
		}
	}

	return field.String() != currentField.String()
}

func isNe(fl FieldLevel) bool {
	return !isEq(fl)
}

func isLteCrossStructField(fl FieldLevel) bool {
	field := fl.Field()
	kind := field.Kind()

	topField, topKind, ok := fl.GetStructFieldOK()
	if !ok || topKind != kind {
		return false
	}

	switch kind {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return field.Int() <= topField.Int()

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return field.Uint() <= topField.Uint()

	case reflect.Float32, reflect.Float64:
		return int64(field.Len()) <= int64(topField.Len())

	case reflect.Struct:
		fieldType := field.Type()

		if fieldType != topField.Type() {
			return false
		}

		if fieldType == timeType {
			fieldTime := field.Interface().(time.Time)
			topTime := topField.Interface().(time.Time)

			return fieldTime.Before(topTime) || fieldTime.Equal(topTime)
		}
	}

	return field.String() <= topField.String()
}

func isLtCrossStructField(fl FieldLevel) bool {
	field := fl.Field()
	kind := field.Kind()

	topField, topKind, ok := fl.GetStructFieldOK()
	if !ok || topKind != kind {
		return false
	}

	switch kind {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return field.Int() < topField.Int()

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return field.Uint() < topField.Uint()

	case reflect.Float32, reflect.Float64:
		return int64(field.Len()) < int64(topField.Len())

	case reflect.Struct:
		fieldType := field.Type()

		if fieldType != topField.Type() {
			return false
		}

		if fieldType == timeType {
			fieldTime := field.Interface().(time.Time)
			topTime := topField.Interface().(time.Time)

			return fieldTime.Before(topTime)
		}
	}

	return field.String() < topField.String()
}

func isGteCrossStructField(fl FieldLevel) bool {
	field := fl.Field()
	kind := field.Kind()

	topField, topKind, ok := fl.GetStructFieldOK()
	if !ok || topKind != kind {
		return false
	}

	switch kind {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return field.Int() >= topField.Int()

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return field.Uint() >= topField.Uint()

	case reflect.Float32, reflect.Float64:
		return field.Float() >= topField.Float()

	case reflect.Struct:
		fieldType := field.Type()

		if fieldType != topField.Type() {
			return false
		}

		if fieldType == timeType {
			fieldTime := field.Interface().(time.Time)
			topTime := topField.Interface().(time.Time)

			return fieldTime.After(topTime) || fieldTime.Equal(topTime)
		}
	}

	return field.String() >= topField.String()
}

func isGtCrossStructField(fl FieldLevel) bool {
	field := fl.Field()
	kind := field.Kind()

	topField, topKind, ok := fl.GetStructFieldOK()
	if !ok || topKind != kind {
		return false
	}

	switch kind {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return field.Int() > topField.Int()

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return field.Uint() > topField.Uint()

	case reflect.Slice, reflect.Map, reflect.Array:
		return int64(field.Len()) > int64(topField.Len())

	case reflect.Struct:
		fieldType := field.Type()

		if fieldType != topField.Type() {
			return false
		}

		if fieldType == timeType {
			fieldTime := field.Interface().(time.Time)
			topTime := topField.Interface().(time.Time)

			return fieldTime.After(topTime)
		}
	}

	return field.String() > topField.String()
}

func isNeCrossStructField(fl FieldLevel) bool {
	field := fl.Field()
	kind := field.Kind()

	topField, currentKind, ok := fl.GetStructFieldOK()
	if !ok || currentKind != kind {
		return true
	}

	switch kind {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return topField.Int() != field.Int()

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return topField.Uint() != field.Uint()

	case reflect.Float32, reflect.Float64:
		return topField.Float() != field.Float()

	case reflect.Slice, reflect.Map, reflect.Array:
		return int64(topField.Len()) != int64(field.Len())

	case reflect.Struct:
		fieldType := field.Type()

		if fieldType != topField.Type() {
			return true
		}

		if fieldType == timeType {
			t := field.Interface().(time.Time)
			fieldTime := topField.Interface().(time.Time)

			return !fieldTime.Equal(t)
		}
	}

	return topField.String() != field.String()
}

func isEqCrossStructfield(fl FieldLevel) bool {
	field := fl.Field()
	kind := field.Kind()

	topField, topKind, ok := fl.GetStructFieldOK()
	if !ok || topKind != kind {
		return false
	}

	switch kind {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return topField.Int() == field.Int()

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return topField.Uint() == field.Uint()

	case reflect.Float32, reflect.Float64:
		return topField.Float() == field.Float()

	case reflect.Slice, reflect.Map, reflect.Array:
		return int64(topField.Len()) == int64(field.Len())

	case reflect.Struct:
		fieldType := field.Type()

		if fieldType != topField.Type() {
			return false
		}

		if fieldType == timeType {
			t := field.Interface().(time.Time)
			fieldTime := topField.Interface().(time.Time)

			return fieldTime.Equal(t)
		}
	}
	return topField.String() == field.String()
}

func isEqField(fl FieldLevel) bool {
	field := fl.Field()
	kind := field.Kind()

	currentField, currentKind, ok := fl.GetStructFieldOK()
	if !ok || currentKind != kind {
		return false
	}

	switch kind {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return field.Int() == currentField.Int()

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return field.Uint() == currentField.Uint()

	case reflect.Float32, reflect.Float64:
		return field.Float() == currentField.Float()

	case reflect.Slice, reflect.Map, reflect.Array:
		return int64(field.Len()) == int64(currentField.Len())

	case reflect.Struct:
		fieldType := field.Type()

		if fieldType != currentField.Type() {
			return false
		}

		if fieldType == timeType {
			t := currentField.Interface().(time.Time)
			fieldTime := field.Interface().(time.Time)

			return fieldTime.Equal(t)
		}
	}

	return field.String() == currentField.String()
}

func isEq(fl FieldLevel) bool {
	field := fl.Field()
	param := fl.Param()

	switch field.Kind() {
	case reflect.String:
		return field.String() == param
	case reflect.Slice, reflect.Map, reflect.Array:
		p := asInt(param)
		return int64(field.Len()) == p
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		p := asInt(param)

		return field.Int() == p
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		p := asUint(param)

		return field.Uint() == p

	case reflect.Float32, reflect.Float64:
		p := asFloat(param)

		return field.Float() == p
	}

	panic(fmt.Sprintf("Bad field type %T", field.Interface()))
}

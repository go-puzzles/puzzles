// File:		validator.go
// Created by:	Hoven
// Created on:	2024-09-19
//
// This file is part of the Example Project.
//
// (c) 2024 Example Corp. All rights reserved.

package putils

import (
	"encoding/json"
	"net"
	"net/mail"
	"net/url"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

var (
	alphaMatcher       *regexp.Regexp = regexp.MustCompile(`^[a-zA-Z]+$`)
	letterRegexMatcher *regexp.Regexp = regexp.MustCompile(`[a-zA-Z]`)
	numberRegexMatcher *regexp.Regexp = regexp.MustCompile(`\d`)
	intStrMatcher      *regexp.Regexp = regexp.MustCompile(`^[\+-]?\d+$`)
	urlMatcher         *regexp.Regexp = regexp.MustCompile(`^((ftp|http|https?):\/\/)?(\S+(:\S*)?@)?((([1-9]\d?|1\d\d|2[01]\d|22[0-3])(\.(1?\d{1,2}|2[0-4]\d|25[0-5])){2}(?:\.([0-9]\d?|1\d\d|2[0-4]\d|25[0-4]))|(([a-zA-Z0-9]+([-\.][a-zA-Z0-9]+)*)|((www\.)?))?(([a-z\x{00a1}-\x{ffff}0-9]+-?-?)*[a-z\x{00a1}-\x{ffff}0-9]+)(?:\.([a-z\x{00a1}-\x{ffff}]{2,}))?))(:(\d{1,5}))?((\/|\?|#)[^\s]*)?$`)
	dnsMatcher         *regexp.Regexp = regexp.MustCompile(`^(?:[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?\.)+[a-zA-Z]{2,}$`)
	chineseMatcher     *regexp.Regexp = regexp.MustCompile("[\u4e00-\u9fa5]")
	base64Matcher      *regexp.Regexp = regexp.MustCompile(`^(?:[A-Za-z0-9+\\/]{4})*(?:[A-Za-z0-9+\\/]{2}==|[A-Za-z0-9+\\/]{3}=|[A-Za-z0-9+\\/]{4})$`)
	base64URLMatcher   *regexp.Regexp = regexp.MustCompile(`^([A-Za-z0-9_-]{4})*([A-Za-z0-9_-]{2}(==)?|[A-Za-z0-9_-]{3}=?)?$`)
	binMatcher         *regexp.Regexp = regexp.MustCompile(`^(0b)?[01]+$`)
	hexMatcher         *regexp.Regexp = regexp.MustCompile(`^(#|0x|0X)?[0-9a-fA-F]+$`)
)

// IsAlpha checks if the string contains only letters (a-zA-Z).
func IsAlpha(str string) bool {
	return alphaMatcher.MatchString(str)
}

// IsASCII checks if string is all ASCII char.
func IsASCII(str string) bool {
	for i := 0; i < len(str); i++ {
		if str[i] > unicode.MaxASCII {
			return false
		}
	}
	return true
}

// IsPrintable checks if string is all printable chars.
func IsPrintable(str string) bool {
	for _, r := range str {
		if !unicode.IsPrint(r) {
			if r == '\n' || r == '\r' || r == '\t' || r == '`' {
				continue
			}
			return false
		}
	}
	return true
}

// ContainLetter check if the string contain at least one letter.
func ContainLetter(str string) bool {
	return letterRegexMatcher.MatchString(str)
}

// ContainNumber check if the string contain at least one number.
func ContainNumber(input string) bool {
	return numberRegexMatcher.MatchString(input)
}

// IsJSON checks if the string is valid JSON.
func IsJSON(str string) bool {
	var js json.RawMessage
	return json.Unmarshal([]byte(str), &js) == nil
}

// IsNumberStr check if the string can convert to a number.
func IsNumberStr(s string) bool {
	return IsIntStr(s) || IsFloatStr(s)
}

// IsFloatStr check if the string can convert to a float.
func IsFloatStr(str string) bool {
	_, e := strconv.ParseFloat(str, 64)
	return e == nil
}

// IsIntStr check if the string can convert to a integer.
func IsIntStr(str string) bool {
	return intStrMatcher.MatchString(str)
}

// IsIp check if the string is a ip address.
func IsIp(ipstr string) bool {
	ip := net.ParseIP(ipstr)
	return ip != nil
}

// IsIpV4 check if the string is a ipv4 address.
func IsIpV4(ipstr string) bool {
	ip := net.ParseIP(ipstr)
	if ip == nil {
		return false
	}
	return ip.To4() != nil
}

// IsIpV6 check if the string is a ipv6 address.
func IsIpV6(ipstr string) bool {
	ip := net.ParseIP(ipstr)
	if ip == nil {
		return false
	}
	return ip.To4() == nil && len(ip) == net.IPv6len
}

// IsPort check if the string is a valid net port.
func IsPort(str string) bool {
	if i, err := strconv.ParseInt(str, 10, 64); err == nil && i > 0 && i < 65536 {
		return true
	}
	return false
}

// IsUrl check if the string is url.
func IsUrl(str string) bool {
	if str == "" || len(str) >= 2083 || len(str) <= 3 || strings.HasPrefix(str, ".") {
		return false
	}
	u, err := url.Parse(str)
	if err != nil {
		return false
	}
	if strings.HasPrefix(u.Host, ".") {
		return false
	}
	if u.Host == "" && (u.Path != "" && !strings.Contains(u.Path, ".")) {
		return false
	}

	return urlMatcher.MatchString(str)
}

// IsDns check if the string is dns.
func IsDns(dns string) bool {
	return dnsMatcher.MatchString(dns)
}

// IsEmail check if the string is a email address.
func IsEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

// ContainChinese check if the string contain mandarin chinese.
func ContainChinese(s string) bool {
	return chineseMatcher.MatchString(s)
}

// IsBase64 check if the string is base64 string.
func IsBase64(base64 string) bool {
	return base64Matcher.MatchString(base64)
}

// IsEmptyString check if the string is empty.
func IsEmptyString(str string) bool {
	return len(str) == 0
}

// IsRegexMatch check if the string match the regexp.
func IsRegexMatch(str, regex string) bool {
	reg := regexp.MustCompile(regex)
	return reg.MatchString(str)
}

// IsStrongPassword check if the string is strong password, if len(password) is less than the length param, return false
// Strong password: alpha(lower+upper) + number + special chars(!@#$%^&*()?><).
func IsStrongPassword(password string, length int) bool {
	if len(password) < length {
		return false
	}
	var num, lower, upper, special bool
	for _, r := range password {
		switch {
		case unicode.IsDigit(r):
			num = true
		case unicode.IsUpper(r):
			upper = true
		case unicode.IsLower(r):
			lower = true
		case unicode.IsSymbol(r), unicode.IsPunct(r):
			special = true
		}
	}

	return num && lower && upper && special
}

// IsWeakPassword check if the string is weak password
// Weak password: only letter or only number or letter + number.
func IsWeakPassword(password string) bool {
	var num, letter, special bool
	for _, r := range password {
		switch {
		case unicode.IsDigit(r):
			num = true
		case unicode.IsLetter(r):
			letter = true
		case unicode.IsSymbol(r), unicode.IsPunct(r):
			special = true
		}
	}

	return (num || letter) && !special
}

// IsZeroValue checks if value is a zero value.
func IsZeroValue(value any) bool {
	if value == nil {
		return true
	}

	rv := reflect.ValueOf(value)
	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}

	if !rv.IsValid() {
		return true
	}

	switch rv.Kind() {
	case reflect.String:
		return rv.Len() == 0
	case reflect.Bool:
		return !rv.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return rv.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return rv.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return rv.Float() == 0
	case reflect.Ptr, reflect.Chan, reflect.Func, reflect.Interface, reflect.Slice, reflect.Map:
		return rv.IsNil()
	}

	return reflect.DeepEqual(rv.Interface(), reflect.Zero(rv.Type()).Interface())
}

// IsGBK check if data encoding is gbk
// Note: this function is implemented by whether double bytes fall within the encoding range of gbk,
// while each byte of utf-8 encoding format falls within the encoding range of gbk.
// Therefore, utf8.valid() should be called first to check whether it is not utf-8 encoding,
// and then call IsGBK() to check gbk encoding. like below
/**
	data := []byte("你好")
	if utf8.Valid(data) {
		fmt.Println("data encoding is utf-8")
	}else if(IsGBK(data)) {
		fmt.Println("data encoding is GBK")
	}
	fmt.Println("data encoding is unknown")
**/
func IsGBK(data []byte) bool {
	i := 0
	for i < len(data) {
		if data[i] <= 0xff {
			i++
			continue
		} else {
			if data[i] >= 0x81 &&
				data[i] <= 0xfe &&
				data[i+1] >= 0x40 &&
				data[i+1] <= 0xfe &&
				data[i+1] != 0xf7 {
				i += 2
				continue
			} else {
				return false
			}
		}
	}

	return true
}

// IsNumber check if the value is number(integer, float) or not.
func IsNumber(v any) bool {
	return IsInt(v) || IsFloat(v)
}

// IsFloat check if the value is float(float32, float34) or not.
func IsFloat(v any) bool {
	switch v.(type) {
	case float32, float64:
		return true
	}
	return false
}

// IsInt check if the value is integer(int, unit) or not.
func IsInt(v any) bool {
	switch v.(type) {
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, uintptr:
		return true
	}
	return false
}

// IsBin check if a give string is a valid binary value or not.
func IsBin(v string) bool {
	return binMatcher.MatchString(v)
}

// IsHex check if a give string is a valid hexadecimal value or not.
func IsHex(v string) bool {
	return hexMatcher.MatchString(v)
}

// IsBase64URL check if a give string is a valid URL-safe Base64 encoded string.
func IsBase64URL(v string) bool {
	return base64URLMatcher.MatchString(v)
}

// IsJWT check if a give string is a valid JSON Web Token (JWT).
func IsJWT(v string) bool {
	strings := strings.Split(v, ".")
	if len(strings) != 3 {
		return false
	}

	for _, s := range strings {
		if !IsBase64URL(s) {
			return false
		}
	}

	return true
}

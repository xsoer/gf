// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gstr

import (
	"strings"

	"github.com/gogf/gf/v2/internal/utils"
)

// IsNumeric tests whether the given string s is numeric.
func IsNumeric(s string) bool {
	return utils.IsNumeric(s)
}

// IsExported checks and returns whether given `name` is exported in Golang.
func IsExported(name string) bool {
	if name == "" {
		return false
	}
	return IsLetterUpper(name[0])
}

// IsLetterLower tests whether the given byte b is in lower case.
func IsLetterLower(b byte) bool {
	return utils.IsLetterLower(b)
}

// IsLetterUpper tests whether the given byte b is in upper case.
func IsLetterUpper(b byte) bool {
	return utils.IsLetterUpper(b)
}

// IsGNUVersion checks and returns whether given `version` is valid GNU version string.
func IsGNUVersion(version string) bool {
	if version != "" && (version[0] == 'v' || version[0] == 'V') {
		version = version[1:]
	}
	if version == "" {
		return false
	}
	var array = strings.Split(version, ".")
	if len(array) > 3 {
		return false
	}
	for _, v := range array {
		if v == "" {
			return false
		}
		if !IsNumeric(v) {
			return false
		}
		if v[0] == '-' || v[0] == '+' {
			return false
		}
	}
	return true
}

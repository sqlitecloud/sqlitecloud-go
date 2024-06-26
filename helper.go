//
//                    ////              SQLite Cloud
//        ////////////  ///
//      ///             ///  ///        Product     : SQLite Cloud GO SDK
//     ///             ///  ///         Version     : 1.0.0
//     //             ///   ///  ///    Date        : 2021/08/31
//    ///             ///   ///  ///    Author      : Andreas Pfeil
//   ///             ///   ///  ///
//   ///     //////////   ///  ///      Description : GO Functions for parsing
//   ////                ///  ///                     SQLite Cloud values and
//     ////     //////////   ///                      enquoting strings.
//        ////            ////
//          ////     /////
//             ///                      Copyright   : 2021 by SQLite Cloud Inc.
//
// -----------------------------------------------------------------------TAB=2

package sqlitecloud

import (
	"errors"
	"fmt"
	"strings"
)

// Helper functions

func escapeString(s, quote string) string {
	repeatedquote := quote + quote
	s = strings.Replace(s, repeatedquote, quote, -1)
	s = strings.Replace(s, quote, repeatedquote, -1)
	return s
}

// SQCloudEnquoteString enquotes the given string if necessary and returns the result as a news created string.
// If the given string contains a '"', the '"' is properly escaped.
// If the given string contains one or more spaces, the whole string is enquoted with '"'
func SQCloudEnquoteString(s string) string {
	s = strings.TrimSpace(s)

	singlequote := "'"
	doublequote := "\""

	if !strings.Contains(s, singlequote) && !strings.Contains(s, doublequote) {
		if strings.Contains(s, " ") {
			return fmt.Sprintf("'%s'", s)
		}

		return s
	}

	if l := len(s); l > 1 {
		for _, q := range []string{singlequote, doublequote} {
			if strings.HasPrefix(s, q) && strings.HasSuffix(s, q) {
				unquoted := s[1 : l-1]
				if !strings.Contains(unquoted, q) {
					return s
				}
				return q + escapeString(unquoted, q) + q
			}
		}
	}

	for _, q := range []string{singlequote, doublequote} {
		if !strings.Contains(s, q) {
			return q + s + q
		}
	}

	return singlequote + escapeString(s, singlequote) + singlequote
}

/*
// SQCloudEnquoteString enquotes the given string if necessary and returns the result as a news created string.
// If the given string contains a '"', the '"' is properly escaped.
// If the given string contains one or more spaces, the whole string is enquoted with '"'
func SQCloudEnquoteString(Token string) string {
	Token = strings.Replace(Token, "\"", "\"\"", -1)
	Token = strings.Replace(Token, "'", "''", -1)
	if strings.Contains(Token, " ") {
		return fmt.Sprintf("\"%s\"", Token)
	}
	return Token
}
*/

// parseBool parses the given string value and tries to interpret the value as a bool.
// true is returned if the value is "TRUE", "ENABLED" or 1.
// false is returned if the value is "FALSE", "DISABLED" or 0.
// The specified defaultValue is returned, if the given string value was an emptry string.
// An error is returned in any other case.
func parseBool(value string, defaultValue bool) (bool, error) {
	switch strings.ToUpper(strings.TrimSpace(value)) {
	case "FALSE", "DISABLED", "0":
		return false, nil
	case "TRUE", "ENABLED", "1":
		return true, nil
	case "":
		return defaultValue, nil
	default:
		return false, errors.New("ERROR: Not a Boolean value")
	}
}

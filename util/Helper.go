package util

import (
	"errors"
	"fmt"
	"regexp"
)

var (
	ptnLetter   = regexp.MustCompile("^[a-zA-Z]$")
	ptnNumber   = regexp.MustCompile("^[0-9]$")
	ptnLiteral  = regexp.MustCompile("^[_a-zA-Z0-9]$")
	ptnOperator = regexp.MustCompile("^[*+\\-<>=!&|^%/,]$")
)

func IsLetter(s string) bool {
	return ptnLetter.MatchString(s)
}

func IsNumber(s string) bool {
	return ptnNumber.MatchString(s)
}

func IsLiteral(s string) bool {
	return ptnLiteral.MatchString(s)
}

func IsOperator(s string) bool {
	return ptnOperator.MatchString(s)
}

func ErrWithArgs(format string, args ...string) error {
	errMsg := fmt.Sprintf(format, args)
	return errors.New(errMsg)
}

func MapToString(elements []interface{}) []string {
	var results []string
	for _, element := range elements {
		results = append(results, element.(string))
	}
	return results
}

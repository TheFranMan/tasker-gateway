package common

import (
	"regexp"
	"strconv"
)

var regexID = regexp.MustCompile(`^[1-9]+[0-9]*$`)
var regexToken = regexp.MustCompile(`^[0-9a-z]{8}-[0-9a-z]{4}-[0-9a-z]{4}-[0-9a-z]{4}-[0-9a-z]{12}$`)

func ValidID(id int) bool {
	return regexID.MatchString(strconv.Itoa(id))
}

func ValidToken(token string) bool {
	return regexToken.MatchString(token)
}

package common

import (
	"regexp"
	"strconv"
)

var regexID = regexp.MustCompile(`^[1-9]+[0-9]*$`)

func ValidID(id int) bool {
	return regexID.MatchString(strconv.Itoa(id))
}

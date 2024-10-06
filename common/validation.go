package common

import "regexp"

var regexID = regexp.MustCompile(`^[1-9]+[0-9]*$`)

func ValidID(id string) bool {
	return regexID.MatchString(id)
}

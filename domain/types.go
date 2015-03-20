package domain

import (
	"fmt"
	"regexp"
)

const (
	NameRegex = "^[A-Za-z0-9]+$"

	NOT_FOUND_ERROR = "not found"
	MALFORMED_ERROR = "malformed"
)

var nameRegex *regexp.Regexp

func init() {
	// A Regexp is safe for concurrent use by multiple goroutines.
	nameRegex = regexp.MustCompile(NameRegex)
	nameRegex.Longest()
}

type Validator interface {
	Validate() error
}

type Name string

func (t Name) Validate() error {

	if len(t) == 0 {
		return fmt.Errorf("malformed name")
	}

	if !nameRegex.MatchString(string(t)) {
		return fmt.Errorf("malformed name [%s]", t)
	}

	return nil
}

type URL string

func (t URL) Validate() error {
	return nil
}

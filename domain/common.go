package domain

import (
	"fmt"
	"regexp"
)

const (
	NameRegex = "^[A-Za-z0-9]+$"
	URLRegex  = `^((ftp|http|https):\/\/)?(\S+(:\S*)?@)?((([1-9]\d?|1\d\d|2[01]\d|22[0-3])(\.(1?\d{1,2}|2[0-4]\d|25[0-5])){2}(?:\.([0-9]\d?|1\d\d|2[0-4]\d|25[0-4]))|((www\.)?)?(([a-z\x{00a1}-\x{ffff}0-9]+-?-?_?)*[a-z\x{00a1}-\x{ffff}0-9]+)(?:\.([a-z\x{00a1}-\x{ffff}]{2,}))?)|localhost)(:(\d{1,5}))?((\/|\?|#)[^\s]*)?$`

	NOT_FOUND_ERROR = "not found"
	MALFORMED_ERROR = "malformed"
)

type NamesRepository interface {
	Get(name Name) (URL, error)
	Put(name Name, url URL) error
	DeleteAll() error
}

var (
	nameRegex *regexp.Regexp
	urlRegex  *regexp.Regexp
)

func init() {
	// A Regexp is safe for concurrent use by multiple goroutines.
	nameRegex = regexp.MustCompile(NameRegex)
	nameRegex.Longest()

	urlRegex = regexp.MustCompile(URLRegex)
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

	if len(t) == 0 {
		return fmt.Errorf("malformed url")
	}

	if !urlRegex.MatchString(string(t)) {
		return fmt.Errorf("malformed url [%s]", t)
	}
	return nil
}

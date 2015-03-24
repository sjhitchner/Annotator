package domain

import (
	"testing"
)

var NAME_VALID = []string{
	"steve",
	"alex",
	"Steve",
	"Alex23",
	"steve00",
	"aleX",
}
var NAME_INVALID = []string{
	"",
	"steve h",
	"alex$",
	"Steve.me",
	"Alex.com",
}

func TestNameValidate(t *testing.T) {

	for _, str := range NAME_VALID {
		name := Name(str)
		if err := name.Validate(); err != nil {
			t.Fatalf("[%s] should be valid name [%v]", name, err)
		}
	}

	for _, str := range NAME_INVALID {
		name := Name(str)
		if err := name.Validate(); err == nil {
			t.Fatalf("[%s] should be invalid name [%v]", name, err)
		}
	}

}

func TestURLValidate(t *testing.T) {

}

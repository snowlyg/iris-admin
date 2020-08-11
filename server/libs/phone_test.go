package libs

import (
	"regexp"
	"testing"
)

func TestGeneratePhoneNumber_CreatePhoneNumber(t *testing.T) {

	ge := &GeneratePhoneNumber{
		CacheData: nil,
	}

	got := ge.CreatePhoneNumber()

	pattern := "^((13[0-9])|(14[5,7])|(15[0-3,5-9])|(17[0,3,5-8])|(18[0-9])|166|198|199|(147))\\d{8}$"
	match, err := regexp.Match(pattern, []byte(got))
	if err != nil {
		t.Errorf("regexp.Match get error: %v", err)
	}

	if !match {
		t.Errorf("phone number %v not match pattern %v", got, pattern)
	}

}

func TestGeneratePhoneNumber_CreateUniquePhoneNumber(t *testing.T) {

	ge := &GeneratePhoneNumber{
		CacheData: []string{},
	}

	got := ge.CreateUniquePhoneNumber()

	pattern := "^((13[0-9])|(14[5,7])|(15[0-3,5-9])|(17[0,3,5-8])|(18[0-9])|166|198|199|(147))\\d{8}$"
	match, err := regexp.Match(pattern, []byte(got))
	if err != nil {
		t.Errorf("regexp.Match get error: %v", err)
	}

	if !match {
		t.Errorf("phone number %v not match pattern %v", got, pattern)
	}

}

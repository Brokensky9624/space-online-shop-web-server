package types

import (
	"fmt"
	"reflect"
)

type Member struct { // User structure
	ID       string     `json:"id"`
	Username string     `json:"username"`
	Password string     `json:"password"`
	Role     MemberRole `json:"role"`
}

type MemberRole int

const (
	Admin MemberRole = iota
	Normal
)

type MemberCreateParam struct {
	Username string `json:"username" required:"true"`
	Password string `json:"password" required:"true"`
	Email    string `json:"email"`
}

func (param MemberCreateParam) Check() error {
	return checkRequiredFields(param)
}

type MemberEditParam struct {
	Username string `json:"username" required:"true"`
	Email    string `json:"email" required:"true"`
}

func (param MemberEditParam) Check() error {
	return checkRequiredFields(param)
}

type MemberDeleteParam struct {
	Username string `json:"username" required:"true"`
	Password string `json:"password" required:"true"`
}

func (param MemberDeleteParam) Check() error {
	return checkRequiredFields(param)
}

type MemberInfoParam struct {
	Username string `json:"username" required:"true"`
	Email    string `json:"email"`
	Birth    string `json:"birth"`
	Phone    string `json:"phone"`
}

func (param MemberInfoParam) Check() error {
	return checkRequiredFields(param)
}

func checkRequiredFields(param interface{}) error {
	t := reflect.TypeOf(param)
	v := reflect.ValueOf(param)

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i).String()
		jsonTag := field.Tag.Get("json")
		requiredTag := field.Tag.Get("required")
		if requiredTag == "true" && value == "" {
			return fmt.Errorf("%s is required but empty", jsonTag)
		}
	}
	return nil
}

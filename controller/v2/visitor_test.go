package v2

import (
	"fmt"
	"reflect"
	"testing"
)

func TestStructTag(t *testing.T) {
	form := VisitorLoginForm{
		VisitorId: "121212",
		ReferUrl:  "http://",
	}
	formRef := reflect.TypeOf(form)
	fmt.Println("Type:", formRef.Name())
	fmt.Println("Kind:", formRef.Kind())
	for i := 0; i < formRef.NumField(); i++ {
		field := formRef.Field(i)
		tag := field.Tag.Get("json")
		fmt.Printf("%d. %v (%v), tag: '%v'\n", i+1, field.Name, field.Type.Name(), tag)
	}
}

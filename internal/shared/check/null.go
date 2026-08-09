package check

import (
	"fmt"
	"reflect"
)

type NullError struct {
	Fields []string
}

func (e *NullError) Error() string {
	return fmt.Sprintf("empty fields: %v", e.Fields)
}

func Null(s any) error {
	v := reflect.ValueOf(s)
	if v.Kind() == reflect.Pointer {
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		return fmt.Errorf("check.Null: expected struct, got %s", v.Kind())
	}

	var empty []string
	t := v.Type()
	for i := 0; i < v.NumField(); i++ {
		f := v.Field(i)
		switch f.Kind() {
		case reflect.String:
			if len(f.String()) == 0 {
				empty = append(empty, t.Field(i).Name)
			}
		case reflect.Slice, reflect.Map, reflect.Pointer, reflect.Interface:
			if f.IsNil() {
				empty = append(empty, t.Field(i).Name)
			}
		default:
			if f.IsZero() {
				empty = append(empty, t.Field(i).Name)
			}
		}
	}

	if len(empty) > 0 {
		return &NullError{Fields: empty}
	}
	return nil
}
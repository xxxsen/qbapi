package qbapi

import (
	"fmt"
	"reflect"
	"strings"
)

func ToMap(in interface{}, tagName string) (map[string]string, error) {
	out := make(map[string]string)

	v := reflect.ValueOf(in)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return nil, fmt.Errorf("ToMap only accepts struct or struct pointer; got %T", v)
	}

	t := v.Type()
	for i := 0; i < v.NumField(); i++ {
		vf := reflect.ValueOf(v.Field(i).Interface())
		//skip nil pointer
		if vf.Kind() == reflect.Ptr && vf.IsNil() {
			continue
		}
		if vf.Kind() == reflect.Ptr {
			vf = vf.Elem()
		}

		fi := t.Field(i)
		tagValue := fi.Tag.Get(tagName)
		if len(tagValue) == 0 {
			return nil, fmt.Errorf("contains non tag field:%s", fi.Name)
		}
		itemList := strings.Split(tagValue, ",")
		value := itemList[0]
		out[value] = fmt.Sprintf("%v", vf)
	}
	return out, nil
}

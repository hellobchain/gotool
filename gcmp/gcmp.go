package gcmp

import (
	"fmt"
	"reflect"
)

func Equal(a, b interface{}) bool {
	return reflect.DeepEqual(a, b)
}

func Diff(a, b interface{}) string {
	va := reflect.ValueOf(a)
	vb := reflect.ValueOf(b)
	if va.Type() != vb.Type() {
		return fmt.Sprintf("type mismatch: %v vs %v", va.Type(), vb.Type())
	}
	return diffValue("", va, vb)
}

func diffValue(path string, va, vb reflect.Value) string {
	if !va.IsValid() || !vb.IsValid() {
		if va.IsValid() != vb.IsValid() {
			return fmt.Sprintf("%s: invalidity mismatch\n", path)
		}
		return ""
	}
	if va.Type().Comparable() && va.CanInterface() {
		if va.Interface() == vb.Interface() {
			return ""
		}
	}
	switch va.Kind() {
	case reflect.Struct:
		for i := 0; i < va.NumField(); i++ {
			f := va.Type().Field(i)
			if sub := diffValue(path+"."+f.Name, va.Field(i), vb.Field(i)); sub != "" {
				return sub
			}
		}
	case reflect.Slice, reflect.Array:
		if va.Len() != vb.Len() {
			return fmt.Sprintf("%s: len %d vs %d\n", path, va.Len(), vb.Len())
		}
		for i := 0; i < va.Len(); i++ {
			if sub := diffValue(fmt.Sprintf("%s[%d]", path, i), va.Index(i), vb.Index(i)); sub != "" {
				return sub
			}
		}
	case reflect.Map:
		if va.Len() != vb.Len() {
			return fmt.Sprintf("%s: map len %d vs %d\n", path, va.Len(), vb.Len())
		}
		for _, k := range va.MapKeys() {
			v1 := va.MapIndex(k)
			v2 := vb.MapIndex(k)
			if !v2.IsValid() {
				return fmt.Sprintf("%s: missing key %v\n", path, k)
			}
			if sub := diffValue(fmt.Sprintf("%s[%v]", path, k), v1, v2); sub != "" {
				return sub
			}
		}
	default:
		if va.Interface() != vb.Interface() {
			return fmt.Sprintf("%s: %v vs %v\n", path, va.Interface(), vb.Interface())
		}
	}
	return ""
}

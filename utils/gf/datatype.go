package gf

import (
	"reflect"
)

func InArray(array []interface{}, element interface{}) bool {
	if element == nil || array == nil {
		return false
	}
	for _, value := range array {
		if reflect.TypeOf(value).Kind() == reflect.TypeOf(element).Kind() {
			if value == element {
				return true
			}
		}
	}
	return false
}

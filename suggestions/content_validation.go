package suggestions

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

func ValidateFieldLengths(obj interface{}) (bool, string) {
	t := reflect.ValueOf(obj)
	for i := 0; i < t.NumField(); i++ {
		fieldValue, ok := t.Field(i).Interface().(string)
		if !ok {
			continue
		}

		field := t.Type().Field(i)

		splitTag := strings.Split(field.Tag.Get("length"), "-")
		if len(splitTag) < 2 {
			continue
		}

		min, _ := strconv.Atoi(splitTag[0])
		max, _ := strconv.Atoi(splitTag[1])
		if len(fieldValue) > max {
			return false, fmt.Sprintf("%v longer than max length %v", field.Name, max)
		}

		if len(fieldValue) < min {
			return false, fmt.Sprintf("%v shorter than min length %v", field.Name, min)
		}
	}

	return true, ""
}

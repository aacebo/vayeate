package serialize

import (
	"encoding/binary"
	"reflect"
)

func Unmarshall[T any](b []byte, v T) {
	value := reflect.Indirect(reflect.ValueOf(v))
	j := 0

	for i := 0; i < value.NumField(); i++ {
		j = parse(j, b, value.Field(i))
	}
}

func parse(startIdx int, b []byte, v reflect.Value) int {
	kind := v.Kind()
	length := binary.BigEndian.Uint32(b[startIdx : startIdx+4])
	value := b[startIdx+4 : startIdx+4+int(length)]

	if kind == reflect.String {
		v.SetString(string(value))
	}

	return startIdx + 4 + int(length)
}

package serialize

import (
	"encoding/binary"
	"reflect"
)

func Marshall[T any](v T) []byte {
	b := []byte{}
	value := reflect.Indirect(reflect.ValueOf(v))

	for i := 0; i < value.NumField(); i++ {
		b = append(b, serialize(value.Field(i))...)
	}

	return b
}

func serialize(v reflect.Value) []byte {
	b := make([]byte, 4)

	if v.Kind() == reflect.String {
		binary.BigEndian.PutUint32(b, uint32(v.Len()))
		b = append(b, []byte(v.String())...)
	}

	return b
}

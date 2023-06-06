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
	kind := v.Kind()
	b := make([]byte, 4)

	if kind == reflect.Slice {
		binary.BigEndian.PutUint32(b, uint32(v.Len()))
		b = append(b, v.Bytes()...)
	} else if kind == reflect.String {
		binary.BigEndian.PutUint32(b, uint32(v.Len()))
		b = append(b, []byte(v.String())...)
	}

	return b
}

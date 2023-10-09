package packet

import "reflect"

func Sizeof(t reflect.Value) int {
	switch t.Kind() {
	case reflect.Array:
		if s := Sizeof(t.Index(0)); s >= 0 {
			return s * t.Len()
		}

	case reflect.Struct:
		sum := 0
		for i, n := 0, t.NumField(); i < n; i++ {
			s := Sizeof(t.Field(i))
			if s < 0 {
				return -1
			}
			sum += s
		}
		return sum

	case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Float32, reflect.Float64, reflect.Complex64, reflect.Complex128:
		return int(t.Type().Size())
	case reflect.Slice:
		return 0
	}

	return -1
}

package utils

import (
	"math/rand"
	"reflect"
	"time"
	"unsafe"
)

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

// BytesToString accepts bytes and returns their string presentation
// instead of string() this method doesn't generate memory allocations,
// BUT it is not safe to use everywhere because it points to the same
// byte slice.  This helps on 0 memory allocations
func BytesToString(b []byte) string {
	bh := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	sh := reflect.StringHeader{bh.Data, bh.Len}
	return *(*string)(unsafe.Pointer(&sh))
}

func StringToBytes(s string) []byte {
	sh := (*reflect.StringHeader)(unsafe.Pointer(&s))
	bh := reflect.SliceHeader{sh.Data, sh.Len, 0}
	return *(*[]byte)(unsafe.Pointer(&bh))
}

func StringDefault(value *string, defaultValue string) string {
	if value != nil {
		return *value
	}

	return defaultValue
}

func RandomString(strlen int) string {
	const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, strlen)
	for i := 0; i < strlen; i++ {
		result[i] = chars[rand.Intn(len(chars))]
	}
	return string(result)
}

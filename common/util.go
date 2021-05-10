package common

import (
	"reflect"
	"strings"
	"time"
)

// interface 转换为 []interface
func ToSlice(array interface{}) []interface{} {
	v := reflect.ValueOf(array)
	if v.Kind() != reflect.Slice {
		panic("ToSlice array not slice")
	}
	l := v.Len()
	ret := make([]interface{}, l)
	for i := 0; i < l; i++ {
		ret[i] = v.Index(i).Interface()
	}
	return ret
}

// 判断一个变量是否为空
func IsEmpty(object interface{}) bool {

	if object == nil {
		return true
	} else if object == "" {
		return true
	} else if object == false {
		return true
	}

	for _, v := range numericZeros {
		if object == v {
			return true
		}
	}

	objValue := reflect.ValueOf(object)

	switch objValue.Kind() {
	case reflect.Map:
		fallthrough
	case reflect.Slice, reflect.Chan:
		{
			return objValue.Len() == 0
		}
	case reflect.Struct:
		switch object.(type) {
		case time.Time:
			return object.(time.Time).IsZero()
		}
	case reflect.Ptr:
		{
			if objValue.IsNil() {
				return true
			}
			switch object.(type) {
			case *time.Time:
				return object.(*time.Time).IsZero()
			default:
				return false
			}
		}
	}
	return false
}

var numericZeros = []interface{}{
	int(0),
	int8(0),
	int16(0),
	int32(0),
	int64(0),
	uint(0),
	uint8(0),
	uint16(0),
	uint32(0),
	uint64(0),
	float32(0),
	float64(0),
}

// 下划线转驼峰，首字母大写
func StrToTF1(str string) (result string) {
	tmp := strings.Split(str, "_")
	for _, tmpS := range tmp {
		if len(tmpS) > 0 {
			result += strings.ToUpper(tmpS[0:1]) + tmpS[1:]
		}
	}
	return
}

var ExcelChar = []string{"", "A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z"}

func ConvertNumToCols(num int) string {
	var cols string
	v := num
	for v > 0 {
		k := v % 26
		if k == 0 {
			k = 26
		}
		v = (v - k) / 26
		cols = ExcelChar[k] + cols
	}
	return cols
}
package common

import (
	"reflect"
	"strings"
)

// 迭代一个对象的所有字段名
func EnumAnObjFieldNames(rv reflect.Type, cb func(f reflect.StructField)) {
	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}
	num := rv.NumField()
	for i := 0; i < num; i++ {
		tmpF := rv.Field(i)
		tmpType := tmpF.Type
		// 如果是时间就不能迭代了
		if tmpType.Kind() == reflect.Struct && !strings.Contains(tmpType.Name(), "Time") && tmpF.Tag.Get("skip") != "true" {
			EnumAnObjFieldNames(tmpType, cb)
		} else {
			cb(tmpF)
		}

	}
}

// 迭代一个对象的所有字段名 todo (json规则) 有点问题, 可能并不需要这个方法
func EnumAnObjFieldNamesByJson(rv reflect.Type, cb func(f reflect.StructField)) {
	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}
	num := rv.NumField()
	for i := 0; i < num; i++ {
		tmpF := rv.Field(i)
		if IsJsonSkip(tmpF) {
			continue
		}
		tmpType := tmpF.Type
		// 如果是时间就不能迭代了
		if tmpType.Kind() == reflect.Struct && !strings.Contains(tmpType.Name(), "Time") && tmpF.Tag.Get("skip") != "true" {
			EnumAnObjFieldNamesByJson(tmpType, cb)
		} else {
			cb(tmpF)
		}

	}
}

// 迭代一个对象的所有字段名(可以返回深度), 深度规则与json转换相同, 深度默认最外层为0, currentDepth 不需要传值; 新增加结构体指针类型迭代(树形结构体不适用[无限循环]),
// 可以自定义 skip tag, 用于跳过
func EnumAnObjFieldNamesWithDepthByJson(rv reflect.Type, cb func(f reflect.StructField, kind reflect.Kind, depth uint), currentDepth ...uint) {
	var depth uint = 0
	if len(currentDepth) > 0 {
		depth = currentDepth[0]
	}
	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}
	num := rv.NumField()
	for i := 0; i < num; i++ {
		tmpF := rv.Field(i)
		if tmpF.Tag.Get("json") == "-" {
			continue
		}
		tmpType := tmpF.Type
		// 如果是时间就不能迭代了
		if ty := tmpType.Kind(); (ty == reflect.Struct || (ty == reflect.Ptr && tmpType.Elem().Kind() == reflect.Struct)) && !strings.Contains(tmpType.Name(), "Time") && tmpF.Tag.Get("skip") != "true" {
			if tmpF.Tag.Get("json") == "" {
				if ty == reflect.Struct {
					EnumAnObjFieldNamesWithDepthByJson(tmpType, cb, depth)
				}
				if ty == reflect.Ptr && tmpType.Elem().Kind() == reflect.Struct {
					EnumAnObjFieldNamesWithDepthByJson(tmpType.Elem(), cb, depth)
				}
			} else {
				if ty == reflect.Struct {
					//cb(tmpF, reflect.Struct, depth+1)
					EnumAnObjFieldNamesWithDepthByJson(tmpType, cb, depth+1)
				}
				if ty == reflect.Ptr && tmpType.Elem().Kind() == reflect.Struct {
					cb(tmpF, reflect.Struct, depth+1)
					EnumAnObjFieldNamesWithDepthByJson(tmpType.Elem(), cb, depth+1)
				}
			}
			// 如果是数组结构体
		} else if ty == reflect.Slice && tmpType.Elem().Kind() == reflect.Struct {
			cb(tmpF, reflect.Slice, depth+1)
			EnumAnObjFieldNamesWithDepthByJson(tmpType.Elem(), cb, depth+1)
			// 如果是指针数组结构体
		} else if ty == reflect.Ptr && tmpType.Elem().Kind() == reflect.Slice && tmpType.Elem().Kind() == reflect.Struct {
			cb(tmpF, reflect.Slice, depth+1)
			EnumAnObjFieldNamesWithDepthByJson(tmpType.Elem().Elem(), cb, depth+1)
		} else {
			cb(tmpF, 0, depth)
		}
	}
}

// 迭代一个对象的所有字段的值    -------- alternated by Zebreay
func EnumAnObjFieldValues(rv reflect.Value, cb func(f reflect.Value)) {
	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}
	num := rv.NumField()
	for i := 0; i < num; i++ {
		tmpF := rv.Field(i)
		// 如果是时间就不能迭代了
		if tmpF.Kind() == reflect.Struct && !strings.Contains(tmpF.Type().Name(), "Time") && rv.Type().Field(i).Tag.Get("skip") != "true" {
			EnumAnObjFieldValues(tmpF, cb)
		} else {
			cb(tmpF)
		}

	}
}

// 终极迭代大法    -------- alternated by Zebreay
func EnumAnStruct(rv reflect.Value, cb func(f reflect.StructField, v reflect.Value)) {
	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}
	num := rv.NumField()
	for i := 0; i < num; i++ {
		tmpT := rv.Type().Field(i)
		tmpV := rv.Field(i)
		// 如果是时间就不能迭代了
		if tmpV.Kind() == reflect.Struct && !strings.Contains(tmpV.Type().Name(), "Time") && tmpT.Tag.Get("skip") != "true" {
			EnumAnStruct(tmpV, cb)
		} else if tmpV.Kind() == reflect.Ptr && tmpV.Type().Elem().Kind() == reflect.Struct && !strings.Contains(tmpV.Type().String(), "Time") && tmpT.Tag.Get("skip") != "true" {
			EnumAnStruct(NewAddr(tmpV.Type().Elem()), cb)
		} else if tmpV.Kind() == reflect.Slice {
			var nonStructSlice bool
			for i := 0; i < tmpV.Len(); i++ {
				if item := tmpV.Index(i); item.Kind() == reflect.Struct && !strings.Contains(tmpV.Type().Name(), "Time") && tmpT.Tag.Get("skip") != "true" {
					EnumAnStruct(item, cb)
				} else {
					nonStructSlice = true
				}
			}
			if nonStructSlice || tmpV.Len() == 0 {
				cb(tmpT, tmpV)
			}
		} else {
			cb(tmpT, tmpV)
		}
	}
}

// 新建一个reflect Value
func NewAddr(rt reflect.Type) reflect.Value {
	n := reflect.New(rt)
	return n
}

func IsUpperFirstLetter(letter string) bool {
	return letter[0] >= 'A' && letter[0] <= 'Z'
}

func IsJsonSkip(f reflect.StructField) bool {
	return f.Tag.Get("json") == "-" || (f.Anonymous && !IsUpperFirstLetter(f.Type.Name())) || (!f.Anonymous && !IsUpperFirstLetter(f.Name))
}

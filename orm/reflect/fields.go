package reflect

import (
	"errors"
	"reflect"
)

// InterateFields 遍历字段
// 可以注释里面说明，这里只接受 XXX 之类的数据
func InterateFields(entity any) (map[string]any, error) {
	if entity == nil {
		return nil, errors.New("不支持 nil")
	}
	typ := reflect.TypeOf(entity) // entity 为 nil 的话，typ 为 nil，val 不为 nil
	val := reflect.ValueOf(entity)
	if val.IsZero() {
		return nil, errors.New("不支持零值")
	}

	// 此处如果用 if 的话，在多级指针的时候，只会解一级指针，会忽略多级指针
	// for 会一直解引用，直到类型不为指针
	for typ.Kind() == reflect.Pointer {
		// Elem 返回指针指向的值，切片里面的值
		// Elem returns a type's element type.
		// It panics if the type's Kind is not Array, Chan, Map, Pointer, or Slice.
		typ = typ.Elem()
		val = val.Elem()
	}

	if typ.Kind() != reflect.Struct {
		return nil, errors.New("不支持类型")
	}

	// NumField returns a struct type's field count.
	// It panics if the type's Kind is not Struct.
	numField := typ.NumField()
	res := make(map[string]any, numField)
	for i := 0; i < numField; i++ {
		// Field returns a struct type's i'th field.
		// Field 类型
		fieldType := typ.Field(i)
		// Field 的值
		fieldVal := val.Field(i)

		// IsExported 确定字段类型是否可导出
		if fieldType.IsExported() {
			// Interface returns v's current value as an interface{}.
			// Interface 返回的值是真正可读的值
			res[fieldType.Name] = fieldVal.Interface()
		} else {
			// 返回的字段的零值
			res[fieldType.Name] = reflect.Zero(fieldType.Type).Interface()
		}
	}
	return res, nil
}

func SetField(entity any, field string, newValue any) error {
	val := reflect.ValueOf(entity)
	for val.Type().Kind() == reflect.Pointer {
		val = val.Elem()
	}

	fieldVal := val.FieldByName(field)
	if !fieldVal.CanSet() {
		return errors.New("字段不可设置")
	}
	fieldVal.Set(reflect.ValueOf(newValue))
	return nil

}

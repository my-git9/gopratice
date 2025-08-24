package reflect

import "reflect"

func IterateArrayOrSlice(entity any) ([]any, error) {
	val := reflect.ValueOf(entity)
	res := make([]any, 0, val.Len())
	for i := 0; i < val.Len(); i++ {
		ele := val.Index(i)
		res = append(res, ele.Interface())
	}
	return res, nil
}

// 返回值是 keys, values ,error
func IterateMap(entity any) ([]any, []any, error) {
	// 实操中需要先检测类型，是否可遍历
	val := reflect.ValueOf(entity)
	resKeys := make([]any, 0, val.Len())
	resValues := make([]any, 0, val.Len())

	// 第一种遍历方式
	itr := val.MapRange()
	for itr.Next() {
		resKeys = append(resKeys, itr.Key().Interface())
		resValues = append(resValues, itr.Value().Interface())
	}

	/*
		// 第二种遍历方式
		keys := val.MapKeys()
		for _, keys := range keys {
			v := val.MapIndex(keys)
			resKeys = append(resKeys, keys.Interface())
			resValues = append(resValues, v.Interface())
		}

	*/
	return resKeys, resValues, nil
}

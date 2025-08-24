package reflect

import (
	"github.com/stretchr/testify/assert"
	"gopratice/orm/reflect/types"
	"reflect"
	"testing"
)

func TestInterateFunc(t *testing.T) {
	testCases := []struct {
		name    string
		entity  any
		wantRes map[string]FuncInfo
		wantErr error
	}{
		{
			name:   "struct",
			entity: types.NewUser("Tom", 18),
			wantRes: map[string]FuncInfo{
				"GetAge": {
					// 下标 0 指向接收器
					InputTypes: []reflect.Type{reflect.TypeOf(types.User{})},
					Name:       "GetAge",
					OutPutTypes: []reflect.Type{
						// Int 类型
						reflect.TypeOf(0),
					},
					Result: []any{18},
				},
			},
		},
		{
			// 方法接收器
			//   1. 以结构体作为输入，那么只能访问到结构体作为接收器的方法
			//   2. 以指针作为输入，那么能访问到任何接收器的方法
			// 输入的第一个参数，永远是接收器本身，如 func (u User) test(test string) {} --> 第一个参数为 u
			name:   "pointer",
			entity: types.NewUserPtr("Tom", 18),
			wantRes: map[string]FuncInfo{
				"GetAge": {
					Name: "GetAge",
					// 下标 0 指向接收器
					InputTypes: []reflect.Type{reflect.TypeOf(&types.User{})},
					OutPutTypes: []reflect.Type{
						// Int 类型
						reflect.TypeOf(0),
					},
					Result: []any{18},
				},
				"ChangeName": {
					Name:        "ChangeName",
					InputTypes:  []reflect.Type{reflect.TypeOf(&types.User{}), reflect.TypeOf("")},
					OutPutTypes: []reflect.Type{},
					Result:      []any{},
				},
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			res, err := IterateFunc(tc.entity)
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantRes, res)
		})
	}
}

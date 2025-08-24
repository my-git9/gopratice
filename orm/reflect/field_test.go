package reflect

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestInterateFields(t *testing.T) {
	type User struct {
		Name string
		age  int
	}

	testCases := []struct {
		name      string
		entity    any
		wantError error
		wantRes   map[string]any
	}{
		{
			name: "test",
			entity: User{
				Name: "Tom",
				age:  18,
			},
			wantError: nil,
			wantRes: map[string]any{
				"Name": "Tom",
				"age":  0,
			},
		},
		{
			name: "pointer",
			entity: &User{
				Name: "Tom",
				age:  18,
			},
			wantError: nil,
			wantRes: map[string]any{
				"Name": "Tom",
				"age":  0,
			},
		},
		{
			name: "multiple pointer",
			entity: func() **User {
				res := &User{
					Name: "Tom",
					age:  18,
				}
				return &res
			}(),
			wantError: nil,
			wantRes: map[string]any{
				"Name": "Tom",
				"age":  0,
			},
		},
		{
			name:      "basis type",
			entity:    18,
			wantError: errors.New("不支持类型"),
			wantRes:   nil,
		},
		{
			name:      "nil",
			entity:    nil,
			wantError: errors.New("不支持 nil"),
			wantRes:   nil,
		},
		{
			name:      "user il",
			entity:    (*User)(nil), // entity 的 type 就不会为 nil 了，值是 nil 了
			wantError: errors.New("不支持零值"),
			wantRes:   nil,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			fields, err := InterateFields(tc.entity)
			assert.Equal(t, tc.wantError, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantRes, fields)
		})
	}
}

func TestSetField(t *testing.T) {
	type User struct {
		Name string
		age  int
	}

	testCases := []struct {
		name string

		entity   any
		field    string
		newValue any
		wantErr  error
		// 修改后的 entity
		wantEntity any
	}{
		{
			name: "struct",
			entity: User{
				Name: "Tom",
			},
			field:    "Name",
			newValue: "Jerry",
			wantErr:  errors.New("字段不可设置"),
		},
		{
			name: "pointer",
			entity: &User{
				Name: "Tom",
			},
			field:    "Name",
			newValue: "Jerry",
			wantEntity: &User{
				Name: "Jerry",
			},
			wantErr: nil,
		},
		{
			name: "pointer exported", // 私有字段不可修改
			entity: &User{
				age: 18,
			},
			field:    "age",
			newValue: 19,
			wantErr:  errors.New("字段不可设置"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := SetField(tc.entity, tc.field, tc.newValue)
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantEntity, tc.entity)
		})
	}

	//var i = 0
	//ptr := &i
	//reflect.ValueOf(ptr).Elem().Set(reflect.ValueOf(12))
	//assert.Equal(t, 12, i)
}

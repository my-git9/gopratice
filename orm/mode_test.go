package orm

import (
	"github.com/stretchr/testify/assert"
	"gopratice/orm/internal/errs"
	"reflect"
	"testing"
)

func Test_parseModel(t *testing.T) {
	testCases := []struct {
		name   string
		entity any

		wantModel *model
		wantErr   error
	}{
		{
			name:    "struct",
			entity:  TestModel{},
			wantErr: errs.ErrPointerOnly,
		},
		{
			name:   "pointer",
			entity: &TestModel{},
			wantModel: &model{
				tableName: "test_model",
				fields: map[string]*field{
					"Id": {
						colName: "id",
					},
					"Age": {
						colName: "age",
					},
					"FirstName": {
						colName: "first_name",
					},
					"LastName": {
						colName: "last_name",
					},
				},
			},
		},
	}
	r := &registry{
		models: make(map[reflect.Type]*model),
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			m, err := r.parseModel(tc.entity)
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantModel, m)
		},
		)
	}
}

func TestRegistry_get(t *testing.T) {
	testCases := []struct {
		name   string
		entity any

		wantModel *model
		wantErr   error
		cacheSize int
	}{
		{
			name:   "pointer",
			entity: &TestModel{},
			wantModel: &model{
				tableName: "test_model",
				fields: map[string]*field{
					"Id": {
						colName: "id",
					},
					"Age": {
						colName: "age",
					},
					"FirstName": {
						colName: "first_name",
					},
					"LastName": {
						colName: "last_name",
					},
				},
			},
			cacheSize: 1,
		},
		{
			name:   "tag",
			entity: func() any {
				type TagTable struct {
					FirstName string `orm:"column=first_name_t"`
				}
				return &TagTable{}
			}(),
			wantModel: &model{
				tableName: "tag_table",
				fields: map[string]*field{
					"FirstName": {
						colName: "first_name_t",
					},
				},
			},
			cacheSize: 1,
		},
		{
			name:   "empty column",
			entity: func() any {
				type TagTable struct {
					FirstName string `orm:"column="`
				}
				return &TagTable{}
			}(),
			wantModel: &model{
				tableName: "tag_table",
				fields: map[string]*field{
					"FirstName": {
						colName: "first_name",
					},
				},
			},
			cacheSize: 1,
		},
		{
			name:   "column only",
			entity: func() any {
				type TagTable struct {
					FirstName string `orm:"column"`
				}
				return &TagTable{}
			}(),
			wantErr: errs.NewErrInvaildTagContent("column"),
			cacheSize: 1,
		},
		{
			name:   "ignore tag",
			entity: func() any {
				type TagTable struct {
					FirstName string `orm:"foo=bar"`
				}
				return &TagTable{}
			}(),
			wantModel: &model{
				tableName: "tag_table",
				fields: map[string]*field{
					"FirstName": {
						colName: "first_name",
					},
				},
			},
			cacheSize: 1,
		},
	}

	r := newRegistry()
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			m, err := r.get(tc.entity)
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantModel, m)

			// 检测数据有没有缓存到 models 里面
			assert.Equal(t, tc.cacheSize, len(r.models))
			typ := reflect.TypeOf(tc.entity)
			m, ok := r.models[typ]
			assert.True(t, ok)
			assert.Equal(t, tc.wantModel, m)
		},
		)
	}
}

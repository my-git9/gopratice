package orm

import (
	"github.com/stretchr/testify/assert"
	"gopratice/orm/internal/errs"
	"testing"
)

func TestDeleter_Build(t *testing.T) {
	testCases := []struct {
		name    string
		builder QueryBuilder

		wantQuery *Query
		wantErr   error
	}{
		{
			name:    "no where",
			builder: &Deleter[TestModel]{},
			wantErr: errs.NewErrNoCondition(),
		},
		{
			name:    "no from",
			builder: (&Deleter[TestModel]{}).Where(C("Age").Eq(18)),
			wantQuery: &Query{
				SQL:  "DELETE FROM `test_model` WHERE `age` = ?;",
				Args: []any{18},
			},
		},
		{
			name:    "with from",
			builder: (&Deleter[TestModel]{}).From("`test_model2`").Where(C("Age").Eq(19)),
			wantQuery: &Query{
				SQL:  "DELETE FROM `test_model2` WHERE `age` = ?;",
				Args: []any{19},
			},
		},
		{
			name:    "empty from",
			builder: (&Deleter[TestModel]{}).From("").Where(C("Age").Eq(19)),
			wantQuery: &Query{
				SQL:  "DELETE FROM `test_model` WHERE `age` = ?;",
				Args: []any{19},
			},
		},
		{
			name:    "from db",
			builder: (&Deleter[TestModel]{}).From("`test_db`.`test_model`").Where(C("Age").Eq(19)),
			wantQuery: &Query{
				SQL:  "DELETE FROM `test_db`.`test_model` WHERE `age` = ?;",
				Args: []any{19},
			},
		},
		// empty where
		{
			name:    "empty where",
			builder: (&Deleter[TestModel]{}).Where(),
			wantErr: errs.NewErrNoCondition(),
		},
		{
			name:    "where not",
			builder: (&Deleter[TestModel]{}).Where(Not(C("Age").Eq(18))),
			wantQuery: &Query{
				SQL:  "DELETE FROM `test_model` WHERE  NOT (`age` = ?);",
				Args: []any{18},
			},
		},
		{
			name:    "where and",
			builder: (&Deleter[TestModel]{}).Where(C("Age").Eq(18).And(C("FirstName").Eq("Tom"))),
			wantQuery: &Query{
				SQL:  "DELETE FROM `test_model` WHERE (`age` = ?) AND (`first_name` = ?);",
				Args: []any{18, "Tom"},
			},
		},
		{
			name:    "where or",
			builder: (&Deleter[TestModel]{}).Where(C("Age").Eq(18).Or(C("FirstName").Eq("Tom"))),
			wantQuery: &Query{
				SQL:  "DELETE FROM `test_model` WHERE (`age` = ?) OR (`first_name` = ?);",
				Args: []any{18, "Tom"},
			},
		},
		{
			name:    "invalid column",
			builder: (&Deleter[TestModel]{}).Where(C("InvalidColumn").Eq(18)),
			wantErr: errs.NewErrUnknownField("InvalidColumn"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			q, err := tc.builder.Build()
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantQuery, q)
		})
	}
}

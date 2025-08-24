package orm

import (
	"context"
	"gopratice/orm/internal/errs"
	"strings"
)

type Deleter[T any] struct {
	table   string
	where   []Predicate
	builder builder
}

func (d *Deleter[T]) Build() (*Query, error) {
	//s.sb = &strings.Builder{}
	//sb := s.sb
	d.builder.sb = &strings.Builder{}
	var err error
	d.builder.model, err = parseModel(new(T))
	if err != nil {
		return nil, err
	}
	sb := d.builder.sb
	sb.WriteString("DELETE FROM ")
	// 怎么拿到表名
	// 通过反射拿到表名
	if d.table == "" {
		//var t T
		//typ := reflect.TypeOf(t)
		sb.WriteByte('`')
		sb.WriteString(d.builder.model.tableName)
		//sb.WriteString(typ.Name())
		sb.WriteByte('`')
	} else {
		//sb.WriteByte('`')
		sb.WriteString(d.table)
		//sb.WriteByte('`')
	}

	if len(d.where) == 0 {
		return nil, errs.NewErrNoCondition()
	}

	//args := make([]any, 0, len(s.where))
	if len(d.where) > 0 {
		sb.WriteString(" WHERE ")
		if err := d.builder.buildPredicates(d.where); err != nil {
			return nil, err
		}
	}

	sb.WriteByte(';')
	return &Query{
		SQL:  sb.String(),
		Args: d.builder.args,
	}, nil
}

func (d *Deleter[T]) From(table string) *Deleter[T] {
	d.table = table
	return d
}

// ids:=[]int{1,2,3}
func (d *Deleter[T]) Where(ps ...Predicate) *Deleter[T] {
	d.where = ps
	return d
}

func (d *Deleter[T]) Get(ctx context.Context) (*T, error) {
	panic("implement me")
}

func (d *Deleter[T]) GetMulti(ctx context.Context) ([]*T, error) {
	panic("implement me")
}

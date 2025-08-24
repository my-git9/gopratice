package orm

import (
	"context"
	"strings"
)

type Selector[T any] struct {
	table string
	where []Predicate
	builder builder
}

func (s *Selector[T]) Build() (*Query, error) {
	//s.sb = &strings.Builder{}
	//sb := s.sb
	s.builder.sb = &strings.Builder{}
	var err error
	s.builder.model, err = parseModel(new(T))
	if err != nil {
		return nil, err
	}
	sb := s.builder.sb
	sb.WriteString("SELECT * FROM ")
	// 怎么拿到表名
	// 通过反射拿到表名
	if s.table == "" {
		//var t T
		//typ := reflect.TypeOf(t)
		sb.WriteByte('`')
		sb.WriteString(s.builder.model.tableName)
		//sb.WriteString(typ.Name())
		sb.WriteByte('`')
	} else {
		//sb.WriteByte('`')
		sb.WriteString(s.table)
		//sb.WriteByte('`')
	}

	//args := make([]any, 0, len(s.where))
	if len(s.where) > 0 {
		sb.WriteString(" WHERE ")
		if err := s.builder.buildPredicates(s.where); err != nil{
			return nil, err
		}
	}

	sb.WriteByte(';')
	return &Query{
		SQL:  sb.String(),
		Args: s.builder.args,
	}, nil
}

func (s *Selector[T]) From(table string) *Selector[T] {
	s.table = table
	return s
}

// ids:=[]int{1,2,3}
func (s *Selector[T]) Where(ps ...Predicate) *Selector[T] {
	s.where = ps
	return s
}

func (s *Selector[T]) Get(ctx context.Context) (*T, error) {
	panic("implement me")
}

func (s *Selector[T]) GetMulti(ctx context.Context) ([]*T, error) {
	panic("implement me")
}

package orm

import (
	"context"
	"fmt"
	"gopratice/orm/internal/errs"
	"strings"
)

type Deleter[T any] struct {
	table string
	model *model
	where []Predicate
	sb    *strings.Builder
	args  []any

	r *registry
}

func (d *Deleter[T]) Build() (*Query, error) {
	//s.sb = &strings.Builder{}
	//sb := s.sb
	d.sb = &strings.Builder{}
	var err error
	d.model, err = d.r.get(new(T))
	if err != nil {
		return nil, err
	}
	sb := d.sb
	sb.WriteString("DELETE * FROM ")
	// 怎么拿到表名
	// 通过反射拿到表名
	if d.table == "" {
		//var t T
		//typ := reflect.TypeOf(t)
		sb.WriteByte('`')
		sb.WriteString(d.model.tableName)
		//sb.WriteString(typ.Name())
		sb.WriteByte('`')
	} else {
		//sb.WriteByte('`')
		sb.WriteString(d.table)
		//sb.WriteByte('`')
	}

	//args := make([]any, 0, len(s.where))
	if len(d.where) > 0 {
		sb.WriteString(" WHERE ")
		p := d.where[0]
		for i := 1; i < len(d.where); i++ {
			p = p.And(d.where[i])
		}
		// 在这里处理 p
		// p.left 构建好
		// p.op 构建好
		// p.right 构建好
		if err := d.buildExpression(p); err != nil{
			return nil, err
		}
	}

	sb.WriteByte(';')
	return &Query{
		SQL:  sb.String(),
		Args: d.args,
	}, nil
}

func (s *Deleter[T]) buildExpression(expr Expression) error {
	switch exp := expr.(type) {
	case nil:
	case Predicate:
		_, ok := exp.left.(Predicate)
		if ok {
			s.sb.WriteByte('(')
		}
		if err := s.buildExpression(exp.left); err != nil {
			return err
		}
		if ok {
			s.sb.WriteByte(')')
		}
		s.sb.WriteByte(' ')
		s.sb.WriteString(exp.op.String())
		s.sb.WriteByte(' ')

		_, ok = exp.right.(Predicate)
		if ok {
			s.sb.WriteByte('(')
		}
		if err := s.buildExpression(exp.right); err != nil {
			return err
		}
		if ok {
			s.sb.WriteByte(')')
		}
	case Column:

		fd, ok := s.model.fields[exp.name]
		if !ok {
			return errs.NewErrUnknownField(exp.name)
		}
		s.sb.WriteByte('`')
		s.sb.WriteString(fd.colName)
		s.sb.WriteByte('`')
	case value:
		s.sb.WriteByte('?')
		s.addArg(exp.val)
	default:
		return fmt.Errorf("unexpected expression type: %v", exp)
	}

	return nil
	/*
		switch left := p.left.(type) {
		case Column:
			sb.WriteByte('`')
			sb.WriteString(left.name)
			sb.WriteByte('`')
		}
		sb.WriteString(p.op.String())
		switch right := p.right.(type) {
		case value:
			sb.WriteByte('?')
			args = append(args, right.val)
		}
	*/
}

func (d *Deleter[T]) addArg(val any) {
	if d.args == nil {
		d.args = make([]any, 0, 8)
	}
	d.args = append(d.args, val)
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

func (s *Deleter[T]) Get(ctx context.Context) (*T, error) {
	panic("implement me")
}

func (s *Deleter[T]) GetMulti(ctx context.Context) ([]*T, error) {
	panic("implement me")
}

package orm

import (
	"fmt"
	"gopratice/orm/internal/errs"
	"strings"
)

type builder struct {
	sb    *strings.Builder
	args  []any
	model *model
}

func (b *builder) buildPredicates(ps []Predicate) error {
	p := ps[0]
	for i := 1; i < len(ps); i++ {
		p = p.And(ps[i])
	}
	// 在这里处理 p
	// p.left 构建好
	// p.op 构建好
	// p.right 构建好

	return b.buildExpression(p)
}

func (b *builder) buildExpression(expr Expression) error {
	switch exp := expr.(type) {
	case nil:
	case Predicate:
		_, ok := exp.left.(Predicate)
		if ok {
			b.sb.WriteByte('(')
		}
		if err := b.buildExpression(exp.left); err != nil {
			return err
		}
		if ok {
			b.sb.WriteByte(')')
		}
		b.sb.WriteByte(' ')
		b.sb.WriteString(exp.op.String())
		b.sb.WriteByte(' ')

		_, ok = exp.right.(Predicate)
		if ok {
			b.sb.WriteByte('(')
		}
		if err := b.buildExpression(exp.right); err != nil {
			return err
		}
		if ok {
			b.sb.WriteByte(')')
		}
	case Column:

		fd, ok := b.model.fields[exp.name]
		if !ok {
			return errs.NewErrUnknownField(exp.name)
		}
		b.sb.WriteByte('`')
		b.sb.WriteString(fd.colName)
		b.sb.WriteByte('`')
	case value:
		b.sb.WriteByte('?')
		b.addArg(exp.val)
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

func (b *builder) addArg(val any) {
	if b.args == nil {
		b.args = make([]any, 0, 8)
	}
	b.args = append(b.args, val)
}
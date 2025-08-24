package orm

// 这种叫衍生类型
type op string

// 这种叫别名
// type op=string
const (
	opEq  op = "="
	opGt  op = ">"
	opLt  op = "<"
	opNot op = "NOT"
	opAnd op = "AND"
	opOr  op = "OR"
)

func (o op) String() string {
	return string(o)
}

// Expression 是一个标记接口，代表表达式
type Expression interface {
	expr()
}

type Predicate struct {
	left  Expression
	op    op
	right Expression
}

/*
// 另一种设计方式
func Eq(column string, arg any) Predicate {
    return Predicate{
        Column: column,
        Op:     "=",
        Arg:    arg,
    }
}
*/

type Column struct {
	name string
}

// 调用：C("id").Eq(1)
func C(name string) Column {
	return Column{
		name: name,
	}
}

func (c Column) Eq(arg any) Predicate {
	return Predicate{
		left:  c,
		op:    opEq,
		right: value{val: arg},
	}
}

func (c Column) Gt(arg any) Predicate {
	return Predicate{
		left:  c,
		op:    opGt,
		right: value{val: arg},
	}
}

func (c Column) Lt(arg any) Predicate {
	return Predicate{
		left:  c,
		op:    opLt,
		right: value{val: arg},
	}
}

func (c Column) expr() {}

// 用法：Not(C("id").Eq(12))
func Not(p Predicate) Predicate {
	return Predicate{
		op: opNot,
		right: Predicate{
			left:  p.left,
			op:    p.op,
			right: p.right,
		},
	}
}

// 用法：C("id").Eq(12).And(C("name").Eq("tom"))
func (left Predicate) And(right Predicate) Predicate {
	return Predicate{
		left:  left,
		op:    opAnd,
		right: right,
	}
}

// 用法：C("id").Eq(12).Or(C("name").Eq("tom"))
func (left Predicate) Or(right Predicate) Predicate {
	return Predicate{
		left:  left,
		op:    opOr,
		right: right,
	}
}

func (p Predicate) expr() {}

type value struct {
	val any
}

func (value) expr() {}

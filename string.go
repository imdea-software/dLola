package dLola

import (
	"errors"
	"fmt"
	//	"strconv"
)

type StrExpr interface {
	AcceptStr(StrExprVisitor)
	Sprint() string
	GetPos() Position
	InstantiateStrExpr(int, int) InstStrExpr
}

type StrExprVisitor interface {
	VisitStringLiteralExpr(StringLiteralExpr) //
	VisitStreamOffsetExpr(StreamOffsetExpr)   //same method with same arguments as in ExprVisitor
	VisitConstExpr(ConstExpr)                 //same method with same arguments as in ExprVisitor
	VisitIfThenElseExpr(IfThenElseExpr)       //same method with same arguments as in ExprVisitor
	VisitStrConcatExpr(StrConcatExpr)         //
}

type StringLiteralExpr struct {
	S   string
	Pos Position
}
type StrConcatExpr struct {
	Left  StrExpr
	Right StrExpr
}

// Accept
func (this StringLiteralExpr) AcceptStr(visitor StrExprVisitor) {
	visitor.VisitStringLiteralExpr(this)
}
func (this StrConcatExpr) AcceptStr(visitor StrExprVisitor) {
	visitor.VisitStrConcatExpr(this)
}
func (this StreamOffsetExpr) AcceptStr(visitor StrExprVisitor) {
	visitor.VisitStreamOffsetExpr(this)
}
func (this ConstExpr) AcceptStr(visitor StrExprVisitor) {
	visitor.VisitConstExpr(this)
}

// Sprint()
func (this StringLiteralExpr) Sprint() string {
	return this.S
}
func (this StrConcatExpr) Sprint() string {
	return fmt.Sprintf("%s strConcat %s", this.Left.Sprint(), this.Right.Sprint())
}

func NewStringLiteralExpr(a, p interface{}) StringLiteralExpr {
	return StringLiteralExpr{a.(string), NewPosition(p)}
}
func NewStrConcatExpr(l, r interface{}) StrConcatExpr {
	return StrConcatExpr{l.(StrExpr), r.(StrExpr)}
}

func (this StringLiteralExpr) GetPos() Position {
	return this.Pos
}
func (this StrConcatExpr) GetPos() Position {
	return this.Left.GetPos()
}

func getStrExpr(e interface{}) (StrExpr, error) {
	switch v := e.(type) {
	case StringExpr:
		return v.StExpr, nil
	case StreamOffsetExpr:
		return v, nil
	case ConstExpr:
		return v, nil
	case StrExpr:
		return v, nil
	default:
		str := fmt.Sprintf("cannot convert to num \"%s\"\n", e.(Expr).Sprint())
		return nil, errors.New(str)
	}

}
func StrExprToExpr(expr StrExpr) Expr {
	switch v := expr.(type) {
	case StreamOffsetExpr:
		return v
	case ConstExpr:
		return v
	default:
		return NewStringExpr(expr) //note: here we use expr, not its type
	}

}

// FLATTEN
type RightSubexprStr interface {
	buildBinaryExpr(left StrExpr, right StrExpr) StrExpr
	getInner() StrExpr
}
type RightStrConcatExpr struct{ E StrExpr }

func (r RightStrConcatExpr) buildBinaryExpr(left StrExpr, right StrExpr) StrExpr {
	return StrConcatExpr{left, right}
}
func (r RightStrConcatExpr) getInner() StrExpr { return r.E }

func NewRightStrConcatExpr(a interface{}) RightStrConcatExpr {
	n, _ := getStrExpr(a)
	return RightStrConcatExpr{n}
}
func FlattenStr(a, b interface{}) Expr {
	exprs := ToSlice(b)
	first, _ := getStrExpr(a)
	if len(exprs) == 0 {
		return StrExprToExpr(first)
	}
	right := exprs[len(exprs)-1].(RightSubexprStr)
	curr := right.getInner()
	for i := len(exprs) - 2; i >= 0; i-- {
		left := exprs[i].(RightSubexprStr)
		curr = right.buildBinaryExpr(left.getInner(), curr)
		right = left
	}
	ret := StrExprToExpr(right.buildBinaryExpr(first, curr))
	return ret
}

type StrComparison interface {
	Sprint() string
	AcceptStrComp(StrComparisonVisitor)
	GetPos() Position
	InstantiateStrCompExpr(int, int) InstStrComparison
}

type StrComparisonVisitor interface {
	VisitStrEqExpr(StrEqExpr)
}

type StrEqExpr struct {
	Left  StrExpr
	Right StrExpr
}

func (this StrEqExpr) AcceptStrComp(visitor StrComparisonVisitor) {
	visitor.VisitStrEqExpr(this)
}
func (this StrEqExpr) Sprint() string {
	return fmt.Sprintf("%s strEq %s", this.Left.Sprint(), this.Right.Sprint())
}
func (this StrEqExpr) GetPos() Position {
	return this.Left.GetPos()
}
func NewStrEqExpr(l, r interface{}) StrEqExpr {
	return StrEqExpr{l.(StrExpr), r.(StrExpr)}
}

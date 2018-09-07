package dLola

import (
	//	"errors"
	"fmt"
	//	"strconv"
)

type StrExpr interface {
	AcceptStr(StrExprVisitor)
	Sprint() string
}

type StrExprVisitor interface {
	VisitStringLiteralExpr(StringLiteralExpr) //
	VisitStreamOffsetExpr(StreamOffsetExpr)   //same method with same arguments as in ExprVisitor
	VisitConstExpr(ConstExpr)                 //same method with same arguments as in ExprVisitor
	VisitIfThenElseExpr(IfThenElseExpr)       //same method with same arguments as in ExprVisitor
	VisitStrConcatExpr(StrConcatExpr)         //
	VisitStrEqExpr(StrEqExpr)                 //
}

type StringLiteralExpr struct {
	S   string
	Pos Position
}
type StrConcatExpr struct {
	Left  StrExpr
	Right StrExpr
}
type StrEqExpr struct {
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
func (this StrEqExpr) AcceptStr(visitor StrExprVisitor) {
	visitor.VisitStrEqExpr(this)
}

// Sprint()
func (this StringLiteralExpr) Sprint() string {
	return this.S
}
func (this StrConcatExpr) Sprint() string {
	return fmt.Sprintf("%s strConcat %s", this.Left.Sprint(), this.Right.Sprint())
}
func (this StrEqExpr) Sprint() string {
	return fmt.Sprintf("%s strEq %s", this.Left.Sprint(), this.Right.Sprint())
}

func NewStringLiteralExpr(a, p interface{}) StringLiteralExpr {
	return StringLiteralExpr{a.(string), NewPosition(p)}
}
func NewStrConcatExpr(l, r interface{}) StrConcatExpr {
	return StrConcatExpr{l.(StrExpr), r.(StrExpr)}
}
func NewStrEqExpr(l, r interface{}) StrEqExpr {
	return StrEqExpr{l.(StrExpr), r.(StrExpr)}
}

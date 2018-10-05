package dLola

import (
	"errors"
	"fmt"
)

type BoolExpr interface {
	Sprint() string
	AcceptBool(BoolExprVisitor)
	GetPos() Position
	InstantiateBoolExpr(int, int) InstBoolExpr
}

type BoolExprVisitor interface {
	VisitTruePredicate(TruePredicate)
	VisitFalsePredicate(FalsePredicate)
	VisitNotPredicate(NotPredicate)
	VisitAndPredicate(AndPredicate)
	VisitOrPredicate(OrPredicate)
	VisitIfThenElseExpr(IfThenElseExpr)     //same method with same arguments as in ExprVisitor
	VisitConstExpr(ConstExpr)               //same method with same arguments as in ExprVisitor
	VisitStreamOffsetExpr(StreamOffsetExpr) //same method with same arguments as in ExprVisitor
	//
	VisitNumComparisonPredicate(NumComparisonPredicate)
	VisitStrComparisonPredicate(StrComparisonPredicate)
	//	VisitPathPredicate(PathPredicate)
	//	VisitStrPredicate(StrPredicate)
	//	VisitTagPredicate(TagPredicate)
}

/*var (
	True      TruePredicate
	False     FalsePredicate
	TrueExpr  BoolExpr = BoolExpr{True}
	FalseExpr BoolExpr = BoolExpr{False}
)*/

type TruePredicate struct{ Pos Position }
type FalsePredicate struct{ Pos Position }

type NotPredicate struct {
	Inner BoolExpr
}
type AndPredicate struct {
	Left  BoolExpr
	Right BoolExpr
}
type OrPredicate struct {
	Left  BoolExpr
	Right BoolExpr
}

/*type IfThenElsePredicate struct {
	If   BoolExpr
	Then BoolExpr
	Else BoolExpr
}*/
type NumComparisonPredicate struct {
	Comp NumComparison
}
type StrComparisonPredicate struct {
	Comp StrComparison
}

func BoolExprToExpr(p BoolExpr) Expr {
	if s, ok := p.(StreamOffsetExpr); ok {
		return s
	} else if k, ok := p.(ConstExpr); ok {
		return k
	} else {
		return NewBooleanExpr(p)
	}
}

func getBoolExpr(e interface{}) (BoolExpr, error) {
	switch v := e.(type) {
	case BooleanExpr:
		return v.BExpr, nil
	case StreamOffsetExpr:
		return v, nil
	case ConstExpr:
		return v, nil
	case BoolExpr:
		return v, nil
	case TruePredicate:
		return v, nil
	case FalsePredicate:
		return v, nil
	default:
		//		fmt.Printf("Is error \n", e)
		str := fmt.Sprintf("cannot convert to bool \"%s\"\n", e.(Expr).Sprint()) //here v has type interface{}
		return nil, errors.New(str)
	}
}

func NewAndPredicate(a, b interface{}) BoolExpr {
	preds := ToSlice(b)
	first, _ := getBoolExpr(a)
	if len(preds) == 0 {
		return first
	}
	right, _ := getBoolExpr(preds[len(preds)-1])
	for i := len(preds) - 2; i >= 0; i-- {
		left, _ := getBoolExpr(preds[i])
		right = AndPredicate{left, right}
	}
	ret := AndPredicate{first, right}
	return ret
}
func NewOrPredicate(a, b interface{}) BoolExpr {
	preds := ToSlice(b)
	first, _ := getBoolExpr(a)
	if len(preds) == 0 {
		return first
	}
	right, _ := getBoolExpr(preds[len(preds)-1])
	for i := len(preds) - 2; i >= 0; i-- {
		left, _ := getBoolExpr(preds[i])
		right = OrPredicate{left, right}
	}
	return OrPredicate{first, right}
}

func NewNotPredicate(n interface{}) NotPredicate {
	return NotPredicate{n.(BoolExpr)}
}

func NewTruePredicate(p interface{}) TruePredicate {
	return TruePredicate{NewPosition(p)}
}

func NewFalsePredicate(p interface{}) FalsePredicate {
	return FalsePredicate{NewPosition(p)}
}

//
// sprint() functions of the different Predicates
//
func (p AndPredicate) Sprint() string {
	return fmt.Sprintf("(%s /\\ %s)", p.Left.Sprint(), p.Right.Sprint())
}
func (p OrPredicate) Sprint() string {
	return fmt.Sprintf("(%s \\/  %s)", p.Left.Sprint(), p.Right.Sprint())
}
func (p NotPredicate) Sprint() string {
	return fmt.Sprintf("~ %s", p.Inner.Sprint())
}
func (p TruePredicate) Sprint() string {
	return fmt.Sprintf("true")
}
func (p FalsePredicate) Sprint() string {
	return fmt.Sprintf("false")
}
func (p NumComparisonPredicate) Sprint() string {
	return p.Comp.Sprint()
}
func (p StrComparisonPredicate) Sprint() string {
	return p.Comp.Sprint()
}

func (this TruePredicate) AcceptBool(v BoolExprVisitor) {
	v.VisitTruePredicate(this)
}
func (this FalsePredicate) AcceptBool(v BoolExprVisitor) {
	v.VisitFalsePredicate(this)
}
func (this NotPredicate) AcceptBool(v BoolExprVisitor) {
	v.VisitNotPredicate(this)
}
func (this AndPredicate) AcceptBool(v BoolExprVisitor) {
	v.VisitAndPredicate(this)
}
func (this OrPredicate) AcceptBool(v BoolExprVisitor) {
	v.VisitOrPredicate(this)
}

func (this IfThenElseExpr) AcceptBool(v BoolExprVisitor) {
	v.VisitIfThenElseExpr(this)
}

// ConstExpr implement AcceptBool so StreamExpr are BoolExpr

func (this ConstExpr) AcceptBool(v BoolExprVisitor) {
	v.VisitConstExpr(this)
}

// StreamExpr impleemnts AcceptBool so StreamExpr are Boolexpr

func (this StreamOffsetExpr) AcceptBool(v BoolExprVisitor) {
	v.VisitStreamOffsetExpr(this)
}

func (this NumComparisonPredicate) AcceptBool(v BoolExprVisitor) {
	v.VisitNumComparisonPredicate(this)
}

func (this StrComparisonPredicate) AcceptBool(v BoolExprVisitor) {
	v.VisitStrComparisonPredicate(this)
}

func NewNumComparisonPredicate(a interface{}) NumComparisonPredicate {
	return NumComparisonPredicate{a.(NumComparison)}
}

func NewStrComparisonPredicate(a interface{}) StrComparisonPredicate {
	return StrComparisonPredicate{a.(StrComparison)}
}

func (this AndPredicate) GetPos() Position {
	return this.Left.GetPos()
}
func (this OrPredicate) GetPos() Position {
	return this.Left.GetPos()
}
func (this NotPredicate) GetPos() Position {
	return this.Inner.GetPos()
}
func (this TruePredicate) GetPos() Position {
	return this.Pos
}
func (this FalsePredicate) GetPos() Position {
	return this.Pos
}

/*func (this IfThenElsePredicate) GetPos() Position {
	return this.If.GetPos()
}*/
func (this NumComparisonPredicate) GetPos() Position {
	return this.Comp.GetPos()
}
func (this StrComparisonPredicate) GetPos() Position {
	return this.Comp.GetPos()
}

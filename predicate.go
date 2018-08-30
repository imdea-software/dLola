package dLola

import (
	"errors"
	"fmt"
)

type BooleanExpr interface {
	Sprint() string
	AcceptBool(BooleanExprVisitor)
}

type BooleanExprVisitor interface {
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
	//	VisitPathPredicate(PathPredicate)
	//	VisitStrPredicate(StrPredicate)
	//	VisitTagPredicate(TagPredicate)
}

var (
	True      TruePredicate
	False     FalsePredicate
	TrueExpr  BoolExpr = BoolExpr{True}
	FalseExpr BoolExpr = BoolExpr{False}
)

type TruePredicate struct{}
type FalsePredicate struct{}

type NotPredicate struct {
	Inner BooleanExpr
}
type AndPredicate struct {
	Left  BooleanExpr
	Right BooleanExpr
}
type OrPredicate struct {
	Left  BooleanExpr
	Right BooleanExpr
}
type IfThenElsePredicate struct {
	If   BooleanExpr
	Then BooleanExpr
	Else BooleanExpr
}
type NumComparisonPredicate struct {
	Comp NumComparison
}

func BooleanExprToExpr(p BooleanExpr) Expr {
	if s, ok := p.(StreamOffsetExpr); ok {
		return s
	} else if k, ok := p.(ConstExpr); ok {
		return k
	} else {
		return NewBoolExpr(p)
	}
}

func getBoolExpr(e interface{}) (BooleanExpr, error) {
	//	fmt.Printf("Converting %s to bool\n", e)
	/*if v, ok := e.(BoolExpr); ok {
		fmt.Printf("Is BoolExpr \n", e)
		return v.BExpr, nil
	} else if v, ok := e.(StreamOffsetExpr); ok {
		fmt.Printf("Is Stream \n", e)
		return v, nil
	} else if k, ok := e.(ConstExpr); ok {
		fmt.Printf("Is ConstExpr \n", e)
		return k, nil
	} else {
		fmt.Printf("Is error \n", e)
		str := fmt.Sprintf("cannot convert to bool \"%s\"\n", e.(Expr).Sprint())
		return nil, errors.New(str)
	}*/
	switch v := e.(type) {
	case BoolExpr:
		return v.BExpr, nil
	case StreamOffsetExpr:
		return v, nil
	case ConstExpr:
		return v, nil
	case BooleanExpr:
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

func NewAndPredicate(a, b interface{}) BooleanExpr {
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
func NewOrPredicate(a, b interface{}) BooleanExpr {
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

func NewNotPredicate(p interface{}) NotPredicate {
	return NotPredicate{p.(BooleanExpr)}
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

func (this TruePredicate) AcceptBool(v BooleanExprVisitor) {
	v.VisitTruePredicate(this)
}
func (this FalsePredicate) AcceptBool(v BooleanExprVisitor) {
	v.VisitFalsePredicate(this)
}
func (this NotPredicate) AcceptBool(v BooleanExprVisitor) {
	v.VisitNotPredicate(this)
}
func (this AndPredicate) AcceptBool(v BooleanExprVisitor) {
	v.VisitAndPredicate(this)
}
func (this OrPredicate) AcceptBool(v BooleanExprVisitor) {
	v.VisitOrPredicate(this)
}

func (this IfThenElseExpr) AcceptBool(v BooleanExprVisitor) {
	v.VisitIfThenElseExpr(this)
}

// ConstExpr implement AcceptBool so StreamExpr are BooleanExpr

func (this ConstExpr) AcceptBool(v BooleanExprVisitor) {
	v.VisitConstExpr(this)
}

// StreamExpr impleemnts AcceptBool so StreamExpr are Booleanexpr

func (this StreamOffsetExpr) AcceptBool(v BooleanExprVisitor) {
	v.VisitStreamOffsetExpr(this)
}

func (this NumComparisonPredicate) AcceptBool(v BooleanExprVisitor) {
	v.VisitNumComparisonPredicate(this)
}

func NewNumComparisonPredicate(a interface{}) NumComparisonPredicate {
	return NumComparisonPredicate{a.(NumComparison)}
}

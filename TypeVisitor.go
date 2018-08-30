package dLola

import (
	"fmt"
	//	"strings"
)

type TypeVisitor struct { //implements ExprVisitor, BooleanExprVisitor, NumExprVisitor, NumComparisonVisitor and StreamExprVisitor
	symTab  map[StreamName]int //symbol table containing the declared variables and their type, StreamName is just string
	errors  []string           //list of all the errors found
	reqType int                //requested type for the subexpressions 0:Bool, 1:Num
}

var (
	boool int = 0
	num   int = 1
)

func (v *TypeVisitor) VisitConstExpr(c ConstExpr) {

}

/*func (v *TypeVisitor) VisitLetExpr(l LetExpr) {

}*/

func (v *TypeVisitor) VisitIfThenElseExpr(ite IfThenElseExpr) {
	checkTypeIf(v, ite)
}

func (v *TypeVisitor) VisitStreamOffsetExpr(s StreamOffsetExpr) {
	s.SExpr.AcceptStream(v)
}

func (v *TypeVisitor) VisitBoolExpr(b BoolExpr) {
	b.BExpr.AcceptBool(v) //will check bool type
}

func (v *TypeVisitor) VisitNumericExpr(n NumericExpr) {
	n.NExpr.AcceptNum(v) //will check num type
}

/*BoolExprVisitor methods*/
func (v *TypeVisitor) VisitTruePredicate(t TruePredicate) {
	if v.reqType != boool {
		s := fmt.Sprintf("Line X: Cannot use True in a non-boolean expression")
		v.errors = append(v.errors, s)
	}
}
func (v *TypeVisitor) VisitFalsePredicate(f FalsePredicate) {
	if v.reqType != boool {
		s := fmt.Sprintf("Line X: Cannot use False in a non-boolean expression")
		v.errors = append(v.errors, s)
	}
}
func (v *TypeVisitor) VisitNotPredicate(n NotPredicate) {
	v.reqType = boool
	n.Inner.AcceptBool(v)
}
func (v *TypeVisitor) VisitAndPredicate(a AndPredicate) {
	checkTypeBoolOp(v, a.Left, a.Right)
}
func (v *TypeVisitor) VisitOrPredicate(o OrPredicate) {
	checkTypeBoolOp(v, o.Left, o.Right)
}

func (v *TypeVisitor) VisitNumComparisonPredicate(n NumComparisonPredicate) {
	n.Comp.AcceptNumComp(v)
}

/*END BoolExprVisitor methods*/

/*NumComparisonVisitor methods*/
func (v *TypeVisitor) VisitNumLess(e NumLess) {
	checkTypeNumOp(v, e.Left, e.Right)
}
func (v *TypeVisitor) VisitNumLessEq(e NumLessEq) {
	checkTypeNumOp(v, e.Left, e.Right)
}
func (v *TypeVisitor) VisitNumEq(e NumEq) {
	checkTypeNumOp(v, e.Left, e.Right)
}
func (v *TypeVisitor) VisitNumGreater(e NumGreater) {
	checkTypeNumOp(v, e.Left, e.Right)
}
func (v *TypeVisitor) VisitNumGreaterEq(e NumGreaterEq) {
	checkTypeNumOp(v, e.Left, e.Right)
}
func (v *TypeVisitor) VisitNumNotEq(e NumNotEq) {
	checkTypeNumOp(v, e.Left, e.Right)
}

/*END NumComparisonVisitor methods*/

/*NumExprVisitor methods*/
func (v *TypeVisitor) VisitIntLiteralExpr(i IntLiteralExpr) {
	if v.reqType != boool {
		s := fmt.Sprintf("Line X: Cannot use Int Literal in a non-numeric expression")
		v.errors = append(v.errors, s)
	}

}

func (v *TypeVisitor) VisitFloatLiteralExpr(f FloatLiteralExpr) {
	if v.reqType != boool {
		s := fmt.Sprintf("Line X: Cannot use Float Literal in a non-numeric expression")
		v.errors = append(v.errors, s)
	}

}

func (v *TypeVisitor) VisitNumMulExpr(e NumMulExpr) {
	checkTypeNumOp(v, e.Left, e.Right)
}

func (v *TypeVisitor) VisitNumDivExpr(e NumDivExpr) {
	checkTypeNumOp(v, e.Left, e.Right)
}

func (v *TypeVisitor) VisitNumPlusExpr(e NumPlusExpr) {
	checkTypeNumOp(v, e.Left, e.Right)
}

func (v *TypeVisitor) VisitNumMinusExpr(e NumMinusExpr) {
	checkTypeNumOp(v, e.Left, e.Right)
}

/*END NumExprVisitor methods*/

/*StreamExprVisitor methods*/
func (v *TypeVisitor) VisitStreamFetchExpr(s StreamFetchExpr) {
	streamname := s.Name
	streamoffset := s.Offset
	if _, err := streamoffset.(FloatLiteralExpr); err {
		err := fmt.Sprintf("line X: Stream %s cannot have a non integer offset ", streamname)
		v.errors = append(v.errors, err)
	}

	elem, ok := v.symTab[streamname]
	if ok && elem != v.reqType { //declared and types match

	}
	if ok { //declared but types do not match
		err := fmt.Sprintf("line X: Stream %s is of type %s but it is required to have type %s", streamname, intToType(elem), intToType(v.reqType))
		v.errors = append(v.errors, err)

	} else { //not declared
		err := fmt.Sprintf("line X: Stream %s not declared in line X", streamname)
		v.errors = append(v.errors, err)
	}

}

/*END StreamExprVisitor methods*/

/*Not exported functions*/
func checkTypeNumOp(v *TypeVisitor, left NumExpr, right NumExpr) {
	v.reqType = num
	left.AcceptNum(v)  //will check the left expression
	right.AcceptNum(v) //will check the right expression
}

func checkTypeBoolOp(v *TypeVisitor, left BooleanExpr, right BooleanExpr) {
	v.reqType = boool
	left.AcceptBool(v)  //will check the left expression
	right.AcceptBool(v) //will check the right expression
}

func checkTypeIf(v *TypeVisitor, ite IfThenElseExpr) {
	v.reqType = boool
	ite.If.Accept(v)   //will check the left expression
	ite.Then.Accept(v) //will check the right expression
	ite.Else.Accept(v) //will check the right expression
}

/*func prettyLet(v *TypeVisitor, l LetExpr) {
	l.Bind.Accept(v) //will check the right expression
	l.Body.Accept(v) //will check the right expression
}
*/

func intToType(i int) string {
	switch i {
	case boool:
		return "bool"
	case num:
		return "num"
	default:
		return "Unknown type"
	}

}

/*END Not exported functions*/

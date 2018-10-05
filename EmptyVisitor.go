package dLola

/*import (
	"fmt"
	"strings"
)*/

/*Empty Visitor that traverses recursively the AST, CODE MAY NEED TO BE ADDED IN THEIR SPECIFIC PLACES*/
type EmptyVisitor struct { //implements ExprVisitor, BoolExprVisitor, NumExprVisitor, NumComparisonVisitor and StreamExprVisitor
	//fields, may be stateful
}

/*ExprVisitor methods*/
func (v *EmptyVisitor) VisitConstExpr(c ConstExpr) {

}

func (v *EmptyVisitor) VisitLetExpr(l LetExpr) {
	emptyLet(v, l)
}

func (v *EmptyVisitor) VisitIfThenElseExpr(ite IfThenElseExpr) {
	emptyIf(v, ite)
}

func (v *EmptyVisitor) VisitStringExpr(s StringExpr) {
	s.StExpr.AcceptStr(v)
}

func (v *EmptyVisitor) VisitStreamOffsetExpr(s StreamOffsetExpr) {
	s.SExpr.AcceptStream(v)
}

func (v *EmptyVisitor) VisitBooleanExpr(b BooleanExpr) {
	b.BExpr.AcceptBool(v)
}

func (v *EmptyVisitor) VisitNumericExpr(n NumericExpr) {
	n.NExpr.AcceptNum(v)
}

/*END ExprVisitor methods*/

/*BoolExprVisitor methods*/
func (v *EmptyVisitor) VisitTruePredicate(t TruePredicate) {

}
func (v *EmptyVisitor) VisitFalsePredicate(f FalsePredicate) {

}
func (v *EmptyVisitor) VisitNotPredicate(n NotPredicate) {
	n.Inner.AcceptBool(v)
}
func (v *EmptyVisitor) VisitAndPredicate(a AndPredicate) {
	emptyBoolOp(v, a.Left, a.Right)
}
func (v *EmptyVisitor) VisitOrPredicate(o OrPredicate) {
	emptyBoolOp(v, o.Left, o.Right)
}

func (v *EmptyVisitor) VisitNumComparisonPredicate(n NumComparisonPredicate) {
	n.Comp.AcceptNumComp(v)
}

func (v *EmptyVisitor) VisitStrComparisonPredicate(s StrComparisonPredicate) {
	s.Comp.AcceptStrComp(v)
}

/*END BoolExprVisitor methods*/

/*NumComparisonVisitor methods*/
func (v *EmptyVisitor) VisitNumLess(e NumLess) {
	emptyNumOp(v, e.Left, e.Right)
}
func (v *EmptyVisitor) VisitNumLessEq(e NumLessEq) {
	emptyNumOp(v, e.Left, e.Right)
}
func (v *EmptyVisitor) VisitNumEq(e NumEq) {
	emptyNumOp(v, e.Left, e.Right)
}
func (v *EmptyVisitor) VisitNumGreater(e NumGreater) {
	emptyNumOp(v, e.Left, e.Right)
}
func (v *EmptyVisitor) VisitNumGreaterEq(e NumGreaterEq) {
	emptyNumOp(v, e.Left, e.Right)
}
func (v *EmptyVisitor) VisitNumNotEq(e NumNotEq) {
	emptyNumOp(v, e.Left, e.Right)
}

/*END NumComparisonVisitor methods*/

/*NumExprVisitor methods*/
func (v *EmptyVisitor) VisitIntLiteralExpr(i IntLiteralExpr) {

}

func (v *EmptyVisitor) VisitFloatLiteralExpr(f FloatLiteralExpr) {

}

func (v *EmptyVisitor) VisitNumMulExpr(e NumMulExpr) {
	emptyNumOp(v, e.Left, e.Right)
}

func (v *EmptyVisitor) VisitNumDivExpr(e NumDivExpr) {
	emptyNumOp(v, e.Left, e.Right)
}

func (v *EmptyVisitor) VisitNumPlusExpr(e NumPlusExpr) {
	emptyNumOp(v, e.Left, e.Right)
}

func (v *EmptyVisitor) VisitNumMinusExpr(e NumMinusExpr) {
	emptyNumOp(v, e.Left, e.Right)
}

/*END NumExprVisitor methods*/

/*StreamExprVisitor methods*/
func (v *EmptyVisitor) VisitStreamFetchExpr(s StreamFetchExpr) {

}

/*END StreamExprVisitor methods*/

/*StrExprVisitor methods: strings*/

func (v *EmptyVisitor) VisitStringLiteralExpr(s StringLiteralExpr) {

}

func (v *EmptyVisitor) VisitStrConcatExpr(s StrConcatExpr) {
	emptyStrOp(v, s.Left, s.Right)
}

func (v *EmptyVisitor) VisitStrEqExpr(s StrEqExpr) {
	emptyStrOp(v, s.Left, s.Right)
}

/*END StrExprVisitor methods*/

/*Not exported functions*/
func emptyNumOp(v *EmptyVisitor, left NumExpr, right NumExpr) {
	left.AcceptNum(v)  //will treat the left expression
	right.AcceptNum(v) //will treat the right expression
}

func emptyBoolOp(v *EmptyVisitor, left BoolExpr, right BoolExpr) {
	left.AcceptBool(v)  //will treat the left expression
	right.AcceptBool(v) //will treat the right expression
}

func emptyIf(v *EmptyVisitor, ite IfThenElseExpr) {
	ite.If.Accept(v)   //will treat the left expression
	ite.Then.Accept(v) //will treat the right expression
	ite.Else.Accept(v) //will treat the right expression
}

func emptyLet(v *EmptyVisitor, l LetExpr) {
	l.Bind.Accept(v) //will treat the right expression
	l.Body.Accept(v) //will treat the right expression
}

func emptyStrOp(v *EmptyVisitor, left StrExpr, right StrExpr) {
	left.AcceptStr(v)  //will treat the right expression
	right.AcceptStr(v) //will treat the right expression
}

/*END Not exported functions*/

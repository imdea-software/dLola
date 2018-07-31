package dLola

import (
	"fmt"
	"strings"
)

/*ExprVisitor will pretty print the ast for a generic expression, receiving the current number of tabs(depth) and will increase it*/

type PrettyPrinterVisitor struct { //implements ExprVisitor, BooleanExprVisitor, NumExprVisitor and NumComparisonVisitor
	layer int    //depth of the numeric expression in the overall AST, also number of tabs to print, need correct initialization(stateful)
	s     string //string containing the formatted AST so far, need correct initialization (stateful)
}

func (v *PrettyPrinterVisitor) VisitConstExpr(c ConstExpr) {
	v.s += c.Sprint()
}

func (v *PrettyPrinterVisitor) VisitLetExpr(l LetExpr) {
	prettyLet(v, l)
}

func (v *PrettyPrinterVisitor) VisitIfThenElseExpr(ite IfThenElseExpr) {
	prettyIf(v, ite)
}

func (v *PrettyPrinterVisitor) VisitStreamOffsetExpr(s StreamOffsetExpr) {
	v.s += s.Sprint()
}

func (v *PrettyPrinterVisitor) VisitBoolExpr(b BoolExpr) {
	b.BExpr.AcceptBool(v)
}

func (v *PrettyPrinterVisitor) VisitNumericExpr(n NumericExpr) {
	n.NExpr.AcceptNum(v)
}

/*BoolExprVisitor methods*/
func (v *PrettyPrinterVisitor) VisitTruePredicate(t TruePredicate) {
	v.s += "True"
}
func (v *PrettyPrinterVisitor) VisitFalsePredicate(f FalsePredicate) {
	v.s += "False"
}
func (v *PrettyPrinterVisitor) VisitNotPredicate(n NotPredicate) {
	v.s += "Not "
	n.Inner.AcceptBool(v)
}
func (v *PrettyPrinterVisitor) VisitAndPredicate(a AndPredicate) {
	prettyBoolOp(v, "AND", a.Left, a.Right)
}
func (v *PrettyPrinterVisitor) VisitOrPredicate(o OrPredicate) {
	prettyBoolOp(v, "OR", o.Left, o.Right)
}

func (v *PrettyPrinterVisitor) VisitNumComparisonPredicate(n NumComparisonPredicate) {
	n.Comp.AcceptNumComp(v)
}

/*END BoolExprVisitor methods*/

/*NumComparisonVisitor methods*/
func (v *PrettyPrinterVisitor) VisitNumLess(e NumLess) {
	prettyNumOp(v, "<", e.Left, e.Right)
}
func (v *PrettyPrinterVisitor) VisitNumLessEq(e NumLessEq) {
	prettyNumOp(v, "<=", e.Left, e.Right)
}
func (v *PrettyPrinterVisitor) VisitNumEq(e NumEq) {
	prettyNumOp(v, "==", e.Left, e.Right)
}
func (v *PrettyPrinterVisitor) VisitNumGreater(e NumGreater) {
	prettyNumOp(v, ">", e.Left, e.Right)
}
func (v *PrettyPrinterVisitor) VisitNumGreaterEq(e NumGreaterEq) {
	prettyNumOp(v, ">=", e.Left, e.Right)
}
func (v *PrettyPrinterVisitor) VisitNumNotEq(e NumNotEq) {
	prettyNumOp(v, "!=", e.Left, e.Right)
}

/*END NumComparisonVisitor methods*/

/*NumExprVisitor methods*/
func (v *PrettyPrinterVisitor) VisitIntLiteralExpr(i IntLiteralExpr) {
	v.s += string(i.Num)
}

func (v *PrettyPrinterVisitor) VisitFloatLiteralExpr(f FloatLiteralExpr) {
	v.s += fmt.Sprintf("%f", f.Num)
}

func (v *PrettyPrinterVisitor) VisitNumMulExpr(e NumMulExpr) {
	prettyNumOp(v, "*", e.Left, e.Right)
}

func (v *PrettyPrinterVisitor) VisitNumDivExpr(e NumDivExpr) {
	prettyNumOp(v, "/", e.Left, e.Right)
}

func (v *PrettyPrinterVisitor) VisitNumPlusExpr(e NumPlusExpr) {
	prettyNumOp(v, "+", e.Left, e.Right)
}

func (v *PrettyPrinterVisitor) VisitNumMinusExpr(e NumMinusExpr) {
	prettyNumOp(v, "-", e.Left, e.Right)
}

/*END NumExprVisitor methods*/

func prettyNumOp(v *PrettyPrinterVisitor, op string, left NumExpr, right NumExpr) {
	tabsNow := strings.Repeat("\t", v.layer)
	tabs := tabsNow + "\t"
	v.layer++ //NEXT LAYER
	v.s += op + "{\n" + tabs
	left.AcceptNum(v) //will append the left expression string
	v.s += "\n" + tabs
	right.AcceptNum(v) //will append the right expression string
	v.s += "\n" + tabsNow + "}\n"
	v.layer-- //NOW LAYER
}

func prettyBoolOp(v *PrettyPrinterVisitor, op string, left BooleanExpr, right BooleanExpr) {
	tabsNow := strings.Repeat("\t", v.layer)
	tabs := tabsNow + "\t"
	v.layer++ //NEXT LAYER
	v.s += op + "{\n" + tabs
	left.AcceptBool(v) //will append the left expression string
	v.s += "\n" + tabs
	right.AcceptBool(v) //will append the right expression string
	v.s += "\n" + tabsNow + "}\n"
	v.layer-- //NOW LAYER
}

func prettyIf(v *PrettyPrinterVisitor, ite IfThenElseExpr) {
	tabsNow := strings.Repeat("\t", v.layer)
	tabs := tabsNow + "\t"
	v.layer++ //NEXT LAYER
	v.s += "If{\n" + tabs
	ite.If.Accept(v) //will append the left expression string
	v.s += "\n}Then{\n" + tabs
	ite.Then.Accept(v) //will append the right expression string
	v.s += "\n}Else{\n" + tabs
	ite.Else.Accept(v) //will append the right expression string
	v.s += "\n}\n"
	v.layer-- //NOW LAYER
}

func prettyLet(v *PrettyPrinterVisitor, l LetExpr) {
	tabsNow := strings.Repeat("\t", v.layer)
	tabs := tabsNow + "\t"
	v.layer++ //NEXT LAYER
	v.s += "Let{\n" + tabs + l.Name.Sprint() + "\n}Bind{\n" + tabs
	l.Bind.Accept(v) //will append the right expression string
	v.s += "\n}Body{\n" + tabs
	l.Body.Accept(v) //will append the right expression string
	v.s += "\n}" + tabs + "}\n"
	v.layer-- //NOW LAYER
}

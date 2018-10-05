package dLola

import (
	//	"fmt"
	"strings"
)

/*ExprVisitor will pretty print the ast for a generic expression, receiving the current number of tabs(depth) and will increase it*/

type PrettyPrinterVisitor struct { //implements ExprVisitor, BooleanExprVisitor, NumExprVisitor, NumComparisonVisitor and StrVisitor
	layer int    //depth of the expression in the overall AST, also number of tabs to print, need correct initialization(stateful)
	s     string //string containing the formatted AST so far, need correct initialization (stateful)
}

/*ExprVisitor methods*/
func (v *PrettyPrinterVisitor) VisitConstExpr(c ConstExpr) {
	v.s += c.Sprint() + "\n"
}

func (v *PrettyPrinterVisitor) VisitLetExpr(l LetExpr) {
	prettyLet(v, l)
}

func (v *PrettyPrinterVisitor) VisitIfThenElseExpr(ite IfThenElseExpr) {
	prettyIf(v, ite)
}

func (v *PrettyPrinterVisitor) VisitStringExpr(s StringExpr) {
	tabsNow := strings.Repeat("\t", v.layer)
	v.s += "StringExpr\n" + tabsNow
	s.StExpr.AcceptStr(v)
}

func (v *PrettyPrinterVisitor) VisitStreamOffsetExpr(s StreamOffsetExpr) {
	v.s += s.Sprint() + "\n"
}

func (v *PrettyPrinterVisitor) VisitBooleanExpr(b BooleanExpr) {
	tabsNow := strings.Repeat("\t", v.layer)
	v.s += "BoolExpr\n" + tabsNow
	b.BExpr.AcceptBool(v)
}

func (v *PrettyPrinterVisitor) VisitNumericExpr(n NumericExpr) {
	tabsNow := strings.Repeat("\t", v.layer)
	v.s += "NumericExpr\n" + tabsNow
	n.NExpr.AcceptNum(v)
}

/*END ExprVisitor methods*/

/*BoolExprVisitor methods*/
func (v *PrettyPrinterVisitor) VisitTruePredicate(t TruePredicate) {
	v.s += "True\n"
}
func (v *PrettyPrinterVisitor) VisitFalsePredicate(f FalsePredicate) {
	v.s += "False\n"
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

func (v *PrettyPrinterVisitor) VisitStrComparisonPredicate(s StrComparisonPredicate) {
	s.Comp.AcceptStrComp(v)
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
	v.s += i.Sprint() + "\n"
}

func (v *PrettyPrinterVisitor) VisitFloatLiteralExpr(f FloatLiteralExpr) {
	v.s += f.Sprint() + "\n"
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

/*StreamExprVisitor methods*/
func (v *PrettyPrinterVisitor) VisitStreamFetchExpr(s StreamFetchExpr) {
	//not needed for PrettyPrinter
}

/*END StreamExprVisitor methods*/

/*StrExprVisitor methods: strings*/

func (v *PrettyPrinterVisitor) VisitStringLiteralExpr(s StringLiteralExpr) {
	v.s += "\"" + s.Sprint() + "\"" + "\n"
}

func (v *PrettyPrinterVisitor) VisitStrConcatExpr(s StrConcatExpr) {
	prettyStrOp(v, "StrConcat", s.Left, s.Right)
}

func (v *PrettyPrinterVisitor) VisitStrEqExpr(s StrEqExpr) {
	prettyStrOp(v, "StrEq", s.Left, s.Right)
}

/*END StrExprVisitor methods*/

func prettyNumOp(v *PrettyPrinterVisitor, op string, left NumExpr, right NumExpr) {
	tabsNow := strings.Repeat("\t", v.layer)
	tabs := tabsNow + "\t"
	v.layer++                //NEXT LAYER
	v.s += op + "{\n" + tabs //higher layers put the correct indentation of lower layers
	left.AcceptNum(v)        //will append the left expression string
	v.s += tabs
	right.AcceptNum(v) //will append the right expression string
	v.s += tabsNow + "}\n"
	v.layer-- //NOW LAYER
}

func prettyBoolOp(v *PrettyPrinterVisitor, op string, left BoolExpr, right BoolExpr) {
	tabsNow := strings.Repeat("\t", v.layer)
	tabs := tabsNow + "\t"
	v.layer++ //NEXT LAYER
	v.s += op + "{\n" + tabs
	left.AcceptBool(v) //will append the left expression string
	v.s += tabs
	right.AcceptBool(v) //will append the right expression string
	v.s += tabsNow + "}\n"
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

func prettyStrOp(v *PrettyPrinterVisitor, op string, left StrExpr, right StrExpr) {
	tabsNow := strings.Repeat("\t", v.layer)
	tabs := tabsNow + "\t"
	v.layer++                //NEXT LAYER
	v.s += op + "{\n" + tabs //higher layers put the correct indentation of lower layers
	left.AcceptStr(v)        //will append the left expression string
	v.s += tabs
	right.AcceptStr(v) //will append the right expression string
	v.s += tabsNow + "}\n"
	v.layer-- //NOW LAYER
}

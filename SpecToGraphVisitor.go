package dLola

/*import (
	"fmt"
	"strings"
)*/

/*SpecToGraph Visitor that traverses recursively the AST, to build the Dependency Graph as in well-formed.go */
type SpecToGraphVisitor struct { //implements ExprVisitor, BooleanExprVisitor, NumExprVisitor, NumComparisonVisitor and StreamExprVisitor
	//fields, may be stateful
	g DepGraphAdj //stateful, need correct initialization
	s StreamName  //streamname of the expression
}

/*ExprVisitor methods*/
func (v *SpecToGraphVisitor) VisitConstExpr(c ConstExpr) {

}

func (v *SpecToGraphVisitor) VisitLetExpr(l LetExpr) {
	specToGraphLet(v, l)
}

func (v *SpecToGraphVisitor) VisitIfThenElseExpr(ite IfThenElseExpr) {
	specToGraphIf(v, ite)
}

func (v *SpecToGraphVisitor) VisitStringExpr(s StringExpr) {
	s.StExpr.AcceptStr(v)
}

func (v *SpecToGraphVisitor) VisitStreamOffsetExpr(s StreamOffsetExpr) {
	s.SExpr.AcceptStream(v)
}

func (v *SpecToGraphVisitor) VisitBooleanExpr(b BooleanExpr) {
	b.BExpr.AcceptBool(v)
}

func (v *SpecToGraphVisitor) VisitNumericExpr(n NumericExpr) {
	n.NExpr.AcceptNum(v)
}

/*END ExprVisitor methods*/

/*BoolExprVisitor methods*/
func (v *SpecToGraphVisitor) VisitTruePredicate(t TruePredicate) {

}
func (v *SpecToGraphVisitor) VisitFalsePredicate(f FalsePredicate) {

}
func (v *SpecToGraphVisitor) VisitNotPredicate(n NotPredicate) {
	n.Inner.AcceptBool(v)
}
func (v *SpecToGraphVisitor) VisitAndPredicate(a AndPredicate) {
	specToGraphBoolOp(v, a.Left, a.Right)
}
func (v *SpecToGraphVisitor) VisitOrPredicate(o OrPredicate) {
	specToGraphBoolOp(v, o.Left, o.Right)
}

func (v *SpecToGraphVisitor) VisitNumComparisonPredicate(n NumComparisonPredicate) {
	n.Comp.AcceptNumComp(v)
}

func (v *SpecToGraphVisitor) VisitStrComparisonPredicate(s StrComparisonPredicate) {
	s.Comp.AcceptStrComp(v)
}

/*END BoolExprVisitor methods*/

/*NumComparisonVisitor methods*/
func (v *SpecToGraphVisitor) VisitNumLess(e NumLess) {
	specToGraphNumOp(v, e.Left, e.Right)
}
func (v *SpecToGraphVisitor) VisitNumLessEq(e NumLessEq) {
	specToGraphNumOp(v, e.Left, e.Right)
}
func (v *SpecToGraphVisitor) VisitNumEq(e NumEq) {
	specToGraphNumOp(v, e.Left, e.Right)
}
func (v *SpecToGraphVisitor) VisitNumGreater(e NumGreater) {
	specToGraphNumOp(v, e.Left, e.Right)
}
func (v *SpecToGraphVisitor) VisitNumGreaterEq(e NumGreaterEq) {
	specToGraphNumOp(v, e.Left, e.Right)
}
func (v *SpecToGraphVisitor) VisitNumNotEq(e NumNotEq) {
	specToGraphNumOp(v, e.Left, e.Right)
}

/*END NumComparisonVisitor methods*/

/*NumExprVisitor methods*/
func (v *SpecToGraphVisitor) VisitIntLiteralExpr(i IntLiteralExpr) {

}

func (v *SpecToGraphVisitor) VisitFloatLiteralExpr(f FloatLiteralExpr) {

}

func (v *SpecToGraphVisitor) VisitNumMulExpr(e NumMulExpr) {
	specToGraphNumOp(v, e.Left, e.Right)
}

func (v *SpecToGraphVisitor) VisitNumDivExpr(e NumDivExpr) {
	specToGraphNumOp(v, e.Left, e.Right)
}

func (v *SpecToGraphVisitor) VisitNumPlusExpr(e NumPlusExpr) {
	specToGraphNumOp(v, e.Left, e.Right)
}

func (v *SpecToGraphVisitor) VisitNumMinusExpr(e NumMinusExpr) {
	specToGraphNumOp(v, e.Left, e.Right)
}

/*END NumExprVisitor methods*/

/*StreamExprVisitor methods*/
func (v *SpecToGraphVisitor) VisitStreamFetchExpr(s StreamFetchExpr) {
	a := Adj{v.s, s.Offset.val, s.Name}
	adjs, ok := v.g[v.s]
	if ok {
		if !elem(adjs, a, EqAdj) {
			v.g[v.s] = append(adjs, a) // add the dependency: v.s depends on the value of s
		}
	} else {
		v.g[v.s] = []Adj{a}
	}

}

/*END StreamExprVisitor methods*/

/*StrExprVisitor methods: strings*/

func (v *SpecToGraphVisitor) VisitStringLiteralExpr(s StringLiteralExpr) {

}

func (v *SpecToGraphVisitor) VisitStrConcatExpr(s StrConcatExpr) {
	specToGraphStrOp(v, s.Left, s.Right)
}

func (v *SpecToGraphVisitor) VisitStrEqExpr(s StrEqExpr) {
	specToGraphStrOp(v, s.Left, s.Right)
}

/*END StrExprVisitor methods*/

/*Not exported functions*/
func specToGraphNumOp(v *SpecToGraphVisitor, left NumExpr, right NumExpr) {
	left.AcceptNum(v)  //will treat the left expression
	right.AcceptNum(v) //will treat the right expression
}

func specToGraphBoolOp(v *SpecToGraphVisitor, left BoolExpr, right BoolExpr) {
	left.AcceptBool(v)  //will treat the left expression
	right.AcceptBool(v) //will treat the right expression
}

func specToGraphIf(v *SpecToGraphVisitor, ite IfThenElseExpr) {
	ite.If.Accept(v)   //will treat the left expression
	ite.Then.Accept(v) //will treat the right expression
	ite.Else.Accept(v) //will treat the right expression
}

func specToGraphLet(v *SpecToGraphVisitor, l LetExpr) {
	l.Bind.Accept(v) //will treat the right expression
	l.Body.Accept(v) //will treat the right expression
}

func specToGraphStrOp(v *SpecToGraphVisitor, left StrExpr, right StrExpr) {
	left.AcceptStr(v)  //will treat the right expression
	right.AcceptStr(v) //will treat the right expression
}

/*END Not exported functions*/

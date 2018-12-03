package dLola

/*import (
	"fmt"
	"strings"
)*/

func (this ConstExpr) ConstantSubs(spec *Spec) Expr {
	constDecl, ok := spec.Const[this.Name]
	if ok {
		return constDecl.Val
	}
	return this //othw let it be
}

func (this LetExpr) ConstantSubs(spec *Spec) Expr {
	return LetExpr{this.Name, this.Params, this.Result, this.Bind.ConstantSubs(spec), this.Body.ConstantSubs(spec)}
}
func (this IfThenElseExpr) ConstantSubs(spec *Spec) Expr {
	return IfThenElseExpr{this.If.ConstantSubs(spec), this.Then.ConstantSubs(spec), this.Else.ConstantSubs(spec)}
}
func (this StreamOffsetExpr) ConstantSubs(spec *Spec) Expr { // expr = a[x|d] we will use the default value to infer the type
	return this.SExpr.ConstantSubsStreamExpr(spec) //note it does not follow the pattern of the rest
}
func (this BooleanExpr) ConstantSubs(spec *Spec) Expr {
	return BooleanExpr{this.BExpr.ConstantSubsBoolExpr(spec)}
}
func (this NumericExpr) ConstantSubs(spec *Spec) Expr {
	return NumericExpr{this.NExpr.ConstantSubsNumExpr(spec)}
}
func (this StringExpr) ConstantSubs(spec *Spec) Expr {
	return StringExpr{this.StExpr.ConstantSubsStrExpr(spec)}
}

//Boolean
func (this TruePredicate) ConstantSubsBoolExpr(spec *Spec) BoolExpr {
	return this
}
func (this FalsePredicate) ConstantSubsBoolExpr(spec *Spec) BoolExpr {
	return this
}
func (this NotPredicate) ConstantSubsBoolExpr(spec *Spec) BoolExpr {
	return NotPredicate{this.Inner.ConstantSubsBoolExpr(spec)}
}
func (this StreamOffsetExpr) ConstantSubsBoolExpr(spec *Spec) BoolExpr {
	return this.SExpr.ConstantSubsBoolStreamExpr(spec)
}
func (this ConstExpr) ConstantSubsBoolExpr(spec *Spec) BoolExpr {
	return ConstExpr{this.Name, this.Pos}
}
func (this AndPredicate) ConstantSubsBoolExpr(spec *Spec) BoolExpr {
	return AndPredicate{this.Left.ConstantSubsBoolExpr(spec), this.Right.ConstantSubsBoolExpr(spec)}
}
func (this OrPredicate) ConstantSubsBoolExpr(spec *Spec) BoolExpr {
	return OrPredicate{this.Left.ConstantSubsBoolExpr(spec), this.Right.ConstantSubsBoolExpr(spec)}
}

/*func (this IfThenElsePredicate) ConstantSubsBoolExpr(spec *Spec) BoolExpr {
	return IfThenElsePredicate{this.If.ConstantSubsBoolExpr(spec), this.Then.ConstantSubsBoolExpr(spec), this.Else.ConstantSubsBoolExpr(spec)}
}*/
func (this NumComparisonPredicate) ConstantSubsBoolExpr(spec *Spec) BoolExpr {
	return NumComparisonPredicate{this.Comp.ConstantSubsNumCompExpr(spec)}
}
func (this StrComparisonPredicate) ConstantSubsBoolExpr(spec *Spec) BoolExpr {
	return StrComparisonPredicate{this.Comp.ConstantSubsStrCompExpr(spec)}
}

//Stream
func (this StreamFetchExpr) ConstantSubsStreamExpr(spec *Spec) Expr {
	constDecl, ok := spec.Const[this.Name]
	if ok {
		return constDecl.Val
	}
	return StreamOffsetExpr{this} //othw let it be
}
func (this StreamFetchExpr) ConstantSubsBoolStreamExpr(spec *Spec) BoolExpr {
	constDecl, ok := spec.Const[this.Name]
	if ok {
		return constDecl.Val.(BooleanExpr).BExpr
	}
	return StreamOffsetExpr{this} //othw let it be
}
func (this StreamFetchExpr) ConstantSubsNumStreamExpr(spec *Spec) NumExpr {
	constDecl, ok := spec.Const[this.Name]
	if ok {
		return constDecl.Val.(NumericExpr).NExpr
	}
	return StreamOffsetExpr{this} //othw let it be
}
func (this StreamFetchExpr) ConstantSubsStrStreamExpr(spec *Spec) StrExpr {
	constDecl, ok := spec.Const[this.Name]
	if ok {
		return constDecl.Val.(StringExpr).StExpr
	}
	return StreamOffsetExpr{this} //othw let it be
}

//Num
func (this NumLess) ConstantSubsNumCompExpr(spec *Spec) NumComparison {
	return NumLess{this.Left.ConstantSubsNumExpr(spec), this.Right.ConstantSubsNumExpr(spec)}
}
func (this NumLessEq) ConstantSubsNumCompExpr(spec *Spec) NumComparison {
	return NumLessEq{this.Left.ConstantSubsNumExpr(spec), this.Right.ConstantSubsNumExpr(spec)}
}
func (this NumGreater) ConstantSubsNumCompExpr(spec *Spec) NumComparison {
	return NumGreater{this.Left.ConstantSubsNumExpr(spec), this.Right.ConstantSubsNumExpr(spec)}
}
func (this NumGreaterEq) ConstantSubsNumCompExpr(spec *Spec) NumComparison {
	return NumGreaterEq{this.Left.ConstantSubsNumExpr(spec), this.Right.ConstantSubsNumExpr(spec)}
}
func (this NumEq) ConstantSubsNumCompExpr(spec *Spec) NumComparison {
	return NumEq{this.Left.ConstantSubsNumExpr(spec), this.Right.ConstantSubsNumExpr(spec)}
}
func (this NumNotEq) ConstantSubsNumCompExpr(spec *Spec) NumComparison {
	return NumNotEq{this.Left.ConstantSubsNumExpr(spec), this.Right.ConstantSubsNumExpr(spec)}
}

func (this IntLiteralExpr) ConstantSubsNumExpr(spec *Spec) NumExpr {
	return this
}
func (this FloatLiteralExpr) ConstantSubsNumExpr(spec *Spec) NumExpr {
	return this
}
func (this NumMulExpr) ConstantSubsNumExpr(spec *Spec) NumExpr {
	return NumMulExpr{this.Left.ConstantSubsNumExpr(spec), this.Right.ConstantSubsNumExpr(spec)}
}
func (this NumDivExpr) ConstantSubsNumExpr(spec *Spec) NumExpr {
	return NumDivExpr{this.Left.ConstantSubsNumExpr(spec), this.Right.ConstantSubsNumExpr(spec)}
}
func (this NumPlusExpr) ConstantSubsNumExpr(spec *Spec) NumExpr {
	return NumPlusExpr{this.Left.ConstantSubsNumExpr(spec), this.Right.ConstantSubsNumExpr(spec)}
}
func (this NumMinusExpr) ConstantSubsNumExpr(spec *Spec) NumExpr {
	return NumMinusExpr{this.Left.ConstantSubsNumExpr(spec), this.Right.ConstantSubsNumExpr(spec)}
}
func (this StreamOffsetExpr) ConstantSubsNumExpr(spec *Spec) NumExpr {
	return this.SExpr.ConstantSubsNumStreamExpr(spec)
}
func (this ConstExpr) ConstantSubsNumExpr(spec *Spec) NumExpr {
	constDecl, ok := spec.Const[this.Name]
	if ok {
		return constDecl.Val.(NumericExpr).NExpr
	}
	return this //othw let it be
}

//String
func (this StringLiteralExpr) ConstantSubsStrExpr(spec *Spec) StrExpr {
	return this
}
func (this StrConcatExpr) ConstantSubsStrExpr(spec *Spec) StrExpr {
	return StrConcatExpr{this.Left.ConstantSubsStrExpr(spec), this.Right.ConstantSubsStrExpr(spec)}
}
func (this StreamOffsetExpr) ConstantSubsStrExpr(spec *Spec) StrExpr {
	return this.SExpr.ConstantSubsStrStreamExpr(spec)
}
func (this ConstExpr) ConstantSubsStrExpr(spec *Spec) StrExpr {
	constDecl, ok := spec.Const[this.Name]
	if ok {
		return constDecl.Val.(StringExpr).StExpr
	}
	return this //othw let it be
}
func (this StrEqExpr) ConstantSubsStrCompExpr(spec *Spec) StrComparison {
	return StrEqExpr{this.Left.ConstantSubsStrExpr(spec), this.Right.ConstantSubsStrExpr(spec)}
}

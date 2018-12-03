package dLola

/*import (
	"fmt"
)*/

func getUdependencies(stream InstStreamExpr, adjacencies []Adj, uExpr InstExpr) []InstStreamExpr {
	r := make([]InstStreamExpr, 0)
	if len(adjacencies) > 0 {
		depen := make(map[InstStreamExpr]struct{})
		uExpr.Udependencies(adjacencies, depen)
		for stream, _ := range depen {
			r = append(r, stream)
		}
	}
	return r
}

//Simplify
func (this InstConstExpr) Udependencies(adj []Adj, depen map[InstStreamExpr]struct{}) {

}
func (this InstLetExpr) Udependencies(adj []Adj, depen map[InstStreamExpr]struct{}) {
}
func (this InstIfThenElseExpr) Udependencies(adj []Adj, depen map[InstStreamExpr]struct{}) {
	if len(adj) != len(depen) {
		this.If.Udependencies(adj, depen)
	}
	if len(adj) != len(depen) {
		this.Then.Udependencies(adj, depen)
	}
	if len(adj) != len(depen) {
		this.Else.Udependencies(adj, depen)
	}
}
func (this InstStreamOffsetExpr) Udependencies(adj []Adj, depen map[InstStreamExpr]struct{}) {
	this.SExpr.Udependencies(adj, depen)
}
func (this InstBooleanExpr) Udependencies(adj []Adj, depen map[InstStreamExpr]struct{}) {
	if len(adj) != len(depen) {
		this.BExpr.UdependenciesBool(adj, depen)
	}
}
func (this InstNumericExpr) Udependencies(adj []Adj, depen map[InstStreamExpr]struct{}) {
	if len(adj) != len(depen) {
		this.NExpr.UdependenciesNum(adj, depen)
	}
}
func (this InstStringExpr) Udependencies(adj []Adj, depen map[InstStreamExpr]struct{}) {
	if len(adj) != len(depen) {
		this.StExpr.UdependenciesStr(adj, depen)
	}
}

//Boolean
func (this InstTruePredicate) UdependenciesBool(adj []Adj, depen map[InstStreamExpr]struct{}) {

}
func (this InstFalsePredicate) UdependenciesBool(adj []Adj, depen map[InstStreamExpr]struct{}) {

}
func (this InstNotPredicate) UdependenciesBool(adj []Adj, depen map[InstStreamExpr]struct{}) {
	if len(adj) != len(depen) {
		this.Inner.UdependenciesBool(adj, depen)
	}
}
func (this InstStreamOffsetExpr) UdependenciesBool(adj []Adj, depen map[InstStreamExpr]struct{}) {
	if len(adj) != len(depen) {
		this.SExpr.Udependencies(adj, depen)
	}
}
func (this InstConstExpr) UdependenciesBool(adj []Adj, depen map[InstStreamExpr]struct{}) {

}
func (this InstAndPredicate) UdependenciesBool(adj []Adj, depen map[InstStreamExpr]struct{}) {
	binaryOpBoolDependencies(this.Left, this.Right, adj, depen)
}
func (this InstOrPredicate) UdependenciesBool(adj []Adj, depen map[InstStreamExpr]struct{}) {
	binaryOpBoolDependencies(this.Left, this.Right, adj, depen)
}

/*func (this InstIfThenElsePredicate) InstantiateBoolExpr(tick, tlen int) InstBoolExpr {
	return InstIfThenElsePredicate{this.If.Simplify(), this.Then.Simplify(), this.Else.Simplify()}
}*/
func (this InstNumComparisonPredicate) UdependenciesBool(adj []Adj, depen map[InstStreamExpr]struct{}) {
	this.Comp.UdependenciesNumComp(adj, depen)
}
func (this InstStrComparisonPredicate) UdependenciesBool(adj []Adj, depen map[InstStreamExpr]struct{}) {
	this.Comp.UdependenciesStrComp(adj, depen)
}

//Stream
func (this InstStreamFetchExpr) Udependencies(adj []Adj, depen map[InstStreamExpr]struct{}) {
	if len(adj) != len(depen) {
		depen[this] = struct{}{}
	}
}

//Num
func (this InstNumLess) UdependenciesNumComp(adj []Adj, depen map[InstStreamExpr]struct{}) {
	binaryOpNumDependencies(this.Left, this.Right, adj, depen)
}
func (this InstNumLessEq) UdependenciesNumComp(adj []Adj, depen map[InstStreamExpr]struct{}) {
	binaryOpNumDependencies(this.Left, this.Right, adj, depen)
}
func (this InstNumGreater) UdependenciesNumComp(adj []Adj, depen map[InstStreamExpr]struct{}) {
	binaryOpNumDependencies(this.Left, this.Right, adj, depen)
}
func (this InstNumGreaterEq) UdependenciesNumComp(adj []Adj, depen map[InstStreamExpr]struct{}) {
	binaryOpNumDependencies(this.Left, this.Right, adj, depen)
}
func (this InstNumEq) UdependenciesNumComp(adj []Adj, depen map[InstStreamExpr]struct{}) {
	binaryOpNumDependencies(this.Left, this.Right, adj, depen)
}
func (this InstNumNotEq) UdependenciesNumComp(adj []Adj, depen map[InstStreamExpr]struct{}) {
	binaryOpNumDependencies(this.Left, this.Right, adj, depen)
}

func (this InstIntLiteralExpr) UdependenciesNum(adj []Adj, depen map[InstStreamExpr]struct{}) {

}
func (this InstFloatLiteralExpr) UdependenciesNum(adj []Adj, depen map[InstStreamExpr]struct{}) {

}
func (this InstNumMulExpr) UdependenciesNum(adj []Adj, depen map[InstStreamExpr]struct{}) {
	binaryOpNumDependencies(this.Left, this.Right, adj, depen)
}
func (this InstNumDivExpr) UdependenciesNum(adj []Adj, depen map[InstStreamExpr]struct{}) {
	binaryOpNumDependencies(this.Left, this.Right, adj, depen)
}
func (this InstNumPlusExpr) UdependenciesNum(adj []Adj, depen map[InstStreamExpr]struct{}) {
	binaryOpNumDependencies(this.Left, this.Right, adj, depen)
}
func (this InstNumMinusExpr) UdependenciesNum(adj []Adj, depen map[InstStreamExpr]struct{}) {
	binaryOpNumDependencies(this.Left, this.Right, adj, depen)
}
func (this InstStreamOffsetExpr) UdependenciesNum(adj []Adj, depen map[InstStreamExpr]struct{}) {
	this.SExpr.Udependencies(adj, depen)
}
func (this InstConstExpr) UdependenciesNum(adj []Adj, depen map[InstStreamExpr]struct{}) {

}

//String
func (this InstStringLiteralExpr) UdependenciesStr(adj []Adj, depen map[InstStreamExpr]struct{}) {

}
func (this InstStrConcatExpr) UdependenciesStr(adj []Adj, depen map[InstStreamExpr]struct{}) {
	binaryOpStrDependencies(this.Left, this.Right, adj, depen)
}
func (this InstStreamOffsetExpr) UdependenciesStr(adj []Adj, depen map[InstStreamExpr]struct{}) {
	if len(adj) != len(depen) {
		this.SExpr.Udependencies(adj, depen)
	}
}
func (this InstConstExpr) UdependenciesStr(adj []Adj, depen map[InstStreamExpr]struct{}) {

}
func (this InstStrEqExpr) UdependenciesStrComp(adj []Adj, depen map[InstStreamExpr]struct{}) {
	binaryOpStrDependencies(this.Left, this.Right, adj, depen)
}

//Literals need to implement InstExpr to compile, implementation of Simplify (should not be needed at runtime)
//will be used as the result value of the expression, note they are InstExpr, not the corresponding subtype
func (this InstTruePredicate) Udependencies(adj []Adj, depen map[InstStreamExpr]struct{}) {

}
func (this InstFalsePredicate) Udependencies(adj []Adj, depen map[InstStreamExpr]struct{}) {

}
func (this InstIntLiteralExpr) Udependencies(adj []Adj, depen map[InstStreamExpr]struct{}) {

}
func (this InstFloatLiteralExpr) Udependencies(adj []Adj, depen map[InstStreamExpr]struct{}) {

}
func (this InstStringLiteralExpr) Udependencies(adj []Adj, depen map[InstStreamExpr]struct{}) {

}

func binaryOpBoolDependencies(l, r InstBoolExpr, adj []Adj, depen map[InstStreamExpr]struct{}) {
	if len(adj) != len(depen) {
		l.UdependenciesBool(adj, depen)
	}
	if len(adj) != len(depen) {
		r.UdependenciesBool(adj, depen)
	}
}
func binaryOpNumDependencies(l, r InstNumExpr, adj []Adj, depen map[InstStreamExpr]struct{}) {
	if len(adj) != len(depen) {
		l.UdependenciesNum(adj, depen)
	}
	if len(adj) != len(depen) {
		r.UdependenciesNum(adj, depen)
	}
}
func binaryOpStrDependencies(l, r InstStrExpr, adj []Adj, depen map[InstStreamExpr]struct{}) {
	if len(adj) != len(depen) {
		l.UdependenciesStr(adj, depen)
	}
	if len(adj) != len(depen) {
		r.UdependenciesStr(adj, depen)
	}
}

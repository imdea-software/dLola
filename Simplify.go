package dLola

/*import (
	"errors"
	"fmt"
	"strconv"
)*/

//Simplify
func (this InstConstExpr) Simplify() InstExpr {
	return this
}
func (this InstLetExpr) Simplify() InstExpr {
	return InstLetExpr{this.Name, this.Bind.Simplify(), this.Body.Simplify()}
}
func (this InstIfThenElseExpr) Simplify() InstExpr {
	i := this.If.Simplify()
	_, tbranch := i.(InstTruePredicate)
	_, fbranch := i.(InstFalsePredicate)
	if tbranch {
		return this.Then.Simplify()
	} else {
		if fbranch {
			return this.Else.Simplify()
		}
	}
	return InstIfThenElseExpr{i, this.Then.Simplify(), this.Else.Simplify()}
}
func (this InstStreamOffsetExpr) Simplify() InstExpr {
	return this //note it is not the same pattern as with Substitute
}
func (this InstBooleanExpr) Simplify() InstExpr {
	return InstBooleanExpr{this.BExpr.SimplifyBool()}
}
func (this InstNumericExpr) Simplify() InstExpr {
	return InstNumericExpr{this.NExpr.SimplifyNum()}
}
func (this InstStringExpr) Simplify() InstExpr {
	return InstStringExpr{this.StExpr.SimplifyStr()}
}

//Boolean
func (this InstTruePredicate) SimplifyBool() InstBoolExpr {
	return this
}
func (this InstFalsePredicate) SimplifyBool() InstBoolExpr {
	return this
}
func (this InstNotPredicate) SimplifyBool() InstBoolExpr {
	if _, t := this.Inner.(InstTruePredicate); t {
		return InstFalsePredicate{}
	}
	if _, f := this.Inner.(InstFalsePredicate); f {
		return InstTruePredicate{}
	}
	return InstNotPredicate{this.Inner.SimplifyBool()}
}
func (this InstStreamOffsetExpr) SimplifyBool() InstBoolExpr {
	return this //note it is not the same pattern as with Substitute
}
func (this InstConstExpr) SimplifyBool() InstBoolExpr {
	return this
}
func (this InstAndPredicate) SimplifyBool() InstBoolExpr {
	if _, f := this.Left.(InstFalsePredicate); f {
		return InstFalsePredicate{}
	}
	if _, t := this.Left.(InstTruePredicate); t {
		return this.Right.SimplifyBool()
	}
	if _, t := this.Right.(InstFalsePredicate); t {
		return InstFalsePredicate{}
	}
	if _, t := this.Right.(InstTruePredicate); t {
		return this.Left.SimplifyBool()
	}
	return InstAndPredicate{this.Left.SimplifyBool(), this.Right.SimplifyBool()}
}
func (this InstOrPredicate) SimplifyBool() InstBoolExpr {
	if _, f := this.Left.(InstFalsePredicate); f {
		return this.Right.SimplifyBool()
	}
	if _, t := this.Left.(InstTruePredicate); t {
		return InstTruePredicate{}
	}
	if _, t := this.Right.(InstFalsePredicate); t {
		return this.Left.SimplifyBool()
	}
	if _, t := this.Right.(InstTruePredicate); t {
		return InstTruePredicate{}
	}
	return InstOrPredicate{this.Left.SimplifyBool(), this.Right.SimplifyBool()}
}

/*func (this InstIfThenElsePredicate) InstantiateBoolExpr(tick, tlen int) InstBoolExpr {
	return InstIfThenElsePredicate{this.If.Simplify(), this.Then.Simplify(), this.Else.Simplify()}
}*/
func (this InstNumComparisonPredicate) SimplifyBool() InstBoolExpr {
	return this.Comp.SimplifyNumComp() //does not follow same pattern than Substitute
}
func (this InstStrComparisonPredicate) SimplifyBool() InstBoolExpr {
	return this.Comp.SimplifyStrComp()
}

//Stream

//Num
func (this InstNumLess) SimplifyNumComp() InstBoolExpr {
	if v, ok := operateComp(this.Left, this.Right, lessInt, lessFloat); ok {
		return v
	}
	return InstNumComparisonPredicate{InstNumLess{this.Left.SimplifyNum(), this.Right.SimplifyNum()}}
}
func (this InstNumLessEq) SimplifyNumComp() InstBoolExpr {
	if v, ok := operateComp(this.Left, this.Right, lesseqInt, lesseqFloat); ok {
		return v
	}
	return InstNumComparisonPredicate{InstNumLessEq{this.Left.SimplifyNum(), this.Right.SimplifyNum()}}
}
func (this InstNumGreater) SimplifyNumComp() InstBoolExpr {
	if v, ok := operateComp(this.Left, this.Right, greaterInt, greaterFloat); ok {
		return v
	}
	return InstNumComparisonPredicate{InstNumGreater{this.Left.SimplifyNum(), this.Right.SimplifyNum()}}
}
func (this InstNumGreaterEq) SimplifyNumComp() InstBoolExpr {
	if v, ok := operateComp(this.Left, this.Right, greateqInt, greateqFloat); ok {
		return v
	}
	return InstNumComparisonPredicate{InstNumGreaterEq{this.Left.SimplifyNum(), this.Right.SimplifyNum()}}
}
func (this InstNumEq) SimplifyNumComp() InstBoolExpr {
	if v, ok := operateComp(this.Left, this.Right, eqInt, eqFloat); ok {
		return v
	}
	return InstNumComparisonPredicate{InstNumEq{this.Left.SimplifyNum(), this.Right.SimplifyNum()}}
}
func (this InstNumNotEq) SimplifyNumComp() InstBoolExpr {
	if v, ok := operateComp(this.Left, this.Right, neqInt, neqFloat); ok {
		return v
	}
	return InstNumComparisonPredicate{InstNumNotEq{this.Left.SimplifyNum(), this.Right.SimplifyNum()}}
}

func (this InstIntLiteralExpr) SimplifyNum() InstNumExpr {
	return this
}
func (this InstFloatLiteralExpr) SimplifyNum() InstNumExpr {
	return this
}
func (this InstNumMulExpr) SimplifyNum() InstNumExpr {
	if v, ok := checkNeutralOperate(this.Left, this.Right, 1, multInt, multFloat); ok {
		return v
	}
	return InstNumMulExpr{this.Left.SimplifyNum(), this.Right.SimplifyNum()}
}
func (this InstNumDivExpr) SimplifyNum() InstNumExpr {
	vir, ir := this.Right.(InstIntLiteralExpr)
	vfr, fr := this.Right.(InstFloatLiteralExpr)
	neutralR := (ir && vir.Num == 1) || (fr && vfr.Num == float32(1))
	if neutralR { //divisor is 1
		return this.Left.SimplifyNum()
	}
	return InstNumDivExpr{this.Left.SimplifyNum(), this.Right.SimplifyNum()}
}
func (this InstNumPlusExpr) SimplifyNum() InstNumExpr {
	if v, ok := checkNeutralOperate(this.Left, this.Right, 0, plusInt, plusFloat); ok {
		return v
	}
	return InstNumPlusExpr{this.Left.SimplifyNum(), this.Right.SimplifyNum()}
}
func (this InstNumMinusExpr) SimplifyNum() InstNumExpr {
	if v, ok := checkNeutralOperate(this.Left, this.Right, 0, minusInt, minusFloat); ok {
		return v
	}
	return InstNumMinusExpr{this.Left.SimplifyNum(), this.Right.SimplifyNum()}
}
func (this InstStreamOffsetExpr) SimplifyNum() InstNumExpr {
	return this //note it is not the same pattern as with Substitute
}
func (this InstConstExpr) SimplifyNum() InstNumExpr {
	return this
}

//String
func (this InstStringLiteralExpr) SimplifyStr() InstStrExpr {
	return this
}
func (this InstStrConcatExpr) SimplifyStr() InstStrExpr {
	if v, ok := checkEmptyOperate(this.Left, this.Right, "", concatStr); ok {
		return v
	}
	return InstStrConcatExpr{this.Left.SimplifyStr(), this.Right.SimplifyStr()}
}
func (this InstStreamOffsetExpr) SimplifyStr() InstStrExpr {
	return this //note it is not the same pattern as with Substitute
}
func (this InstConstExpr) SimplifyStr() InstStrExpr {
	return this
}
func (this InstStrEqExpr) SimplifyStrComp() InstBoolExpr {
	if v, ok := operateCompStr(this.Left, this.Right, eqStr); ok {
		return v
	}
	return InstStrComparisonPredicate{InstStrEqExpr{this.Left.SimplifyStr(), this.Right.SimplifyStr()}}
}

//Literals need to implement InstExpr to compile, implementation of Simplify (should not be needed at runtime)
func (this InstTruePredicate) Simplify() InstExpr {
	return this
}
func (this InstFalsePredicate) Simplify() InstExpr {
	return this
}
func (this InstIntLiteralExpr) Simplify() InstExpr {
	return this
}
func (this InstFloatLiteralExpr) Simplify() InstExpr {
	return this
}
func (this InstStringLiteralExpr) Simplify() InstExpr {
	return this
}

//Num comparison auxiliary funcs
func operateComp(left, right InstNumExpr, fcompi func(int, int) bool, fcompf func(float32, float32) bool) (InstBoolExpr, bool) {
	vil, il := left.(InstIntLiteralExpr)
	vfl, fl := left.(InstFloatLiteralExpr)
	vir, ir := right.(InstIntLiteralExpr)
	vfr, fr := right.(InstFloatLiteralExpr)
	if il && ir { //both are int literals, operate
		return convertToInst(fcompi(vil.Num, vir.Num)), true
	}
	if il && fr { //int op float
		return convertToInst(fcompf(float32(vil.Num), vfr.Num)), true
	}
	if fl && ir { //float op int
		return convertToInst(fcompf(vfl.Num, float32(vir.Num))), true
	}
	if fl && fr { //float op float
		return convertToInst(fcompf(vfl.Num, vfr.Num)), true
	}
	return nil, false
}

func lessInt(a, b int) bool {
	return a < b
}
func lesseqInt(a, b int) bool {
	return a <= b
}
func greaterInt(a, b int) bool {
	return a > b
}
func greateqInt(a, b int) bool {
	return a >= b
}
func eqInt(a, b int) bool {
	return a == b
}
func neqInt(a, b int) bool {
	return a != b
}
func lessFloat(a, b float32) bool {
	return a < b
}
func lesseqFloat(a, b float32) bool {
	return a <= b
}
func greaterFloat(a, b float32) bool {
	return a > b
}
func greateqFloat(a, b float32) bool {
	return a >= b
}
func eqFloat(a, b float32) bool {
	return a == b
}
func neqFloat(a, b float32) bool {
	return a != b
}
func convertToInst(b bool) InstBoolExpr {
	if b {
		return InstTruePredicate{}
	}
	return InstFalsePredicate{}
}

//Num expr
func checkNeutralOperate(left, right InstNumExpr, neutral int, fint func(int, int) int, ffloat func(float32, float32) float32) (InstNumExpr, bool) {
	vil, il := left.(InstIntLiteralExpr)
	vfl, fl := left.(InstFloatLiteralExpr)
	neutralL := (il && vil.Num == neutral) || (fl && vfl.Num == float32(neutral))
	vir, ir := right.(InstIntLiteralExpr)
	vfr, fr := right.(InstFloatLiteralExpr)
	neutralR := (ir && vir.Num == neutral) || (fr && vfr.Num == float32(neutral))
	if neutralL { //left operand is neutral of the operation
		return right.SimplifyNum(), true
	}
	if neutralR { //right operand is neutral of the operation
		return left.SimplifyNum(), true
	}
	if il && ir { //both are int literals, operate
		return InstIntLiteralExpr{fint(vil.Num, vir.Num)}, true
	}
	if il && fr { //int op float
		return InstFloatLiteralExpr{ffloat(float32(vil.Num), vfr.Num)}, true
	}
	if fl && ir { //float op int
		return InstFloatLiteralExpr{ffloat(vfl.Num, float32(vir.Num))}, true
	}
	if fl && fr { //float op float
		return InstFloatLiteralExpr{ffloat(vfl.Num, vfr.Num)}, true
	}
	return nil, false
}

func multInt(a, b int) int {
	return a * b
}
func divInt(a, b int) int {
	return a / b
}
func plusInt(a, b int) int {
	return a + b
}
func minusInt(a, b int) int {
	return a - b
}
func multFloat(a, b float32) float32 {
	return a * b
}
func divFloat(a, b float32) float32 {
	return a / b
}
func plusFloat(a, b float32) float32 {
	return a + b
}
func minusFloat(a, b float32) float32 {
	return a - b
}

//String Expr
func checkEmptyOperate(left, right InstStrExpr, neutral string, fstr func(string, string) string) (InstStrExpr, bool) {
	vsl, sl := left.(InstStringLiteralExpr)
	neutralL := sl && vsl.S == neutral
	vsr, sr := right.(InstStringLiteralExpr)
	neutralR := sr && vsr.S == neutral
	if neutralL {
		return right.SimplifyStr(), true
	}
	if neutralR {
		return left.SimplifyStr(), true
	}
	if sl && sr {
		return InstStringLiteralExpr{fstr(vsl.S, vsr.S)}, true
	}
	return nil, false
}

func concatStr(s, r string) string {
	return s + r
}

//String Comp
func operateCompStr(left, right InstStrExpr, fcomps func(string, string) bool) (InstBoolExpr, bool) {
	vsl, sl := left.(InstStringLiteralExpr)
	vsr, sr := right.(InstStringLiteralExpr)
	if sl && sr {
		return convertToInst(fcomps(vsl.S, vsr.S)), true
	}
	return nil, false
}
func eqStr(a, b string) bool {
	return a == b
}

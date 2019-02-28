package dLola

import (
	//	"errors"
	"fmt"
	//	"strconv"
)

func SimplifyExpr(exp InstExpr) InstExpr {
	expSimpl := true
	for expSimpl { //while something in the expression get simplified try to simplify further
		//fmt.Printf("simplifying again: %s\n", exp.Sprint())
		exp, expSimpl = exp.Simplify()
	}
	//fmt.Printf("Done simplifying: %s\n", exp.Sprint())
	return exp
}

//returns the constant if it is indeed a constant and a flag that marks if the expression can be further simplified
func getConstantBool(exp InstBoolExpr) (InstExpr, bool) {
	switch c := exp.(type) {
	case InstTruePredicate:
		return c, false
	case InstFalsePredicate:
		return c, false
	}
	return InstBooleanExpr{exp}, true
}
func getConstantNum(exp InstNumExpr) (InstExpr, bool) {
	switch c := exp.(type) {
	case InstIntLiteralExpr:
		return c, false
	case InstFloatLiteralExpr:
		return c, false
	}
	return InstNumericExpr{exp}, true
}
func getConstantStr(exp InstStrExpr) (InstExpr, bool) {
	switch c := exp.(type) {
	case InstStringLiteralExpr:
		return c, false
	}
	return InstStringExpr{exp}, true
}

func isGround(exp InstExpr) bool {
	keepsimp := false
	if bexp, isbool := exp.(InstBooleanExpr); isbool {
		_, keepsimp = getConstantBool(bexp.BExpr)
	}
	if nexp, isnum := exp.(InstNumericExpr); isnum {
		_, keepsimp = getConstantNum(nexp.NExpr)
	}
	if sexp, isstr := exp.(InstStringExpr); isstr {
		_, keepsimp = getConstantStr(sexp.StExpr)
	}
	return !keepsimp
}

//Simplify
func (this InstConstExpr) Simplify() (InstExpr, bool) {
	return this, false
}
func (this InstLetExpr) Simplify() (InstExpr, bool) {
	bind, simplbind := this.Bind.Simplify()
	if isGround(bind) {
		fmt.Printf("Bind is ground %s\n", bind.Sprint())
		bodySub := this.Body.Substitute(InstStreamFetchExpr{this.Name, 0}, bind) //tick of the binding in a let expression should be 0
		return bodySub.Simplify()                                                //we remove a bind of the let expr, simplifying it
	}
	body, simplbody := this.Body.Simplify()
	return InstLetExpr{this.Name, bind, body}, simplbind || simplbody
}
func (this InstIfThenElseExpr) Simplify() (InstExpr, bool) {
	//fmt.Printf("Simplifying IF: %s\n", this.Sprint())
	i, _ := this.If.Simplify() //will be InstTrue/FalsePredicate or InstBooleanExpr if the If could not be resolved
	_, tbranch := i.(InstTruePredicate)
	_, fbranch := i.(InstFalsePredicate)
	//fmt.Printf("Simplifying IF condition simplified of type: %T\n", i)
	if tbranch {
		//fmt.Printf("Simplifying IF then branch: %s\n", i.Sprint())
		e, _ := this.Then.Simplify()
		return e, true
	} else {
		if fbranch {
			//fmt.Printf("Simplifying IF then branch: %s\n", i.Sprint())
			e, _ := this.Else.Simplify()
			return e, true
		}
	}
	//fmt.Printf("Simplifying IF condition could not be resolved: %s\n", i.Sprint())
	then, simplthen := this.Then.Simplify()
	elsse, simplelsse := this.Else.Simplify()
	return InstIfThenElseExpr{i, then, elsse}, simplthen || simplelsse
}
func (this InstStreamOffsetExpr) Simplify() (InstExpr, bool) {
	return this, false //note it is not the same pattern as with Substitute
}
func (this InstBooleanExpr) Simplify() (InstExpr, bool) {
	b, simpl := this.BExpr.SimplifyBool()
	c, keepsimpl := getConstantBool(b)
	return c, simpl && keepsimpl
}
func (this InstNumericExpr) Simplify() (InstExpr, bool) {
	//fmt.Printf("Simplifying Numeric expression: %s\n", this.Sprint())
	n, simpl := this.NExpr.SimplifyNum()
	c, keepsimpl := getConstantNum(n)
	return c, simpl && keepsimpl
}
func (this InstStringExpr) Simplify() (InstExpr, bool) {
	s, simpl := this.StExpr.SimplifyStr()
	c, keepsimpl := getConstantStr(s)
	return c, simpl && keepsimpl
}

//Boolean
func (this InstTruePredicate) SimplifyBool() (InstBoolExpr, bool) {
	return this, false
}
func (this InstFalsePredicate) SimplifyBool() (InstBoolExpr, bool) {
	return this, false
}
func (this InstNotPredicate) SimplifyBool() (InstBoolExpr, bool) {
	if _, t := this.Inner.(InstTruePredicate); t {
		return InstFalsePredicate{}, true
	}
	if _, f := this.Inner.(InstFalsePredicate); f {
		return InstTruePredicate{}, true
	}
	n, simpl := this.Inner.SimplifyBool()
	return InstNotPredicate{n}, simpl
}
func (this InstStreamOffsetExpr) SimplifyBool() (InstBoolExpr, bool) {
	return this, false //note it is not the same pattern as with Substitute
}
func (this InstConstExpr) SimplifyBool() (InstBoolExpr, bool) {
	return this, false
}
func (this InstAndPredicate) SimplifyBool() (InstBoolExpr, bool) {
	if _, f := this.Left.(InstFalsePredicate); f {
		return InstFalsePredicate{}, true
	}
	if _, t := this.Left.(InstTruePredicate); t {
		return this.Right.SimplifyBool()
	}
	if _, t := this.Right.(InstFalsePredicate); t {
		return InstFalsePredicate{}, true
	}
	if _, t := this.Right.(InstTruePredicate); t {
		return this.Left.SimplifyBool()
	}
	l, lsimpl := this.Left.SimplifyBool()
	r, rsimpl := this.Right.SimplifyBool()
	return InstAndPredicate{l, r}, lsimpl || rsimpl
}
func (this InstOrPredicate) SimplifyBool() (InstBoolExpr, bool) {
	if _, f := this.Left.(InstFalsePredicate); f {
		return this.Right.SimplifyBool()
	}
	if _, t := this.Left.(InstTruePredicate); t {
		return InstTruePredicate{}, true
	}
	if _, t := this.Right.(InstFalsePredicate); t {
		return this.Left.SimplifyBool()
	}
	if _, t := this.Right.(InstTruePredicate); t {
		return InstTruePredicate{}, true
	}
	l, lsimpl := this.Left.SimplifyBool()
	r, rsimpl := this.Right.SimplifyBool()
	return InstOrPredicate{l, r}, lsimpl || rsimpl

}

/*func (this InstIfThenElsePredicate) InstantiateBoolExpr(tick, tlen int) InstBoolExpr {
	return InstIfThenElsePredicate{this.If.Simplify(), this.Then.Simplify(), this.Else.Simplify()}
}*/
func (this InstNumComparisonPredicate) SimplifyBool() (InstBoolExpr, bool) {
	return this.Comp.SimplifyNumComp() //does not follow same pattern than Substitute
}
func (this InstStrComparisonPredicate) SimplifyBool() (InstBoolExpr, bool) {
	return this.Comp.SimplifyStrComp()
}

//Stream

//Num
func (this InstNumLess) SimplifyNumComp() (InstBoolExpr, bool) {
	if v, ok := operateComp(this.Left, this.Right, lessInt, lessFloat); ok {
		return v, true
	}
	l, lsimpl := this.Left.SimplifyNum()
	r, rsimpl := this.Right.SimplifyNum()
	return InstNumComparisonPredicate{InstNumLess{l, r}}, lsimpl || rsimpl
}
func (this InstNumLessEq) SimplifyNumComp() (InstBoolExpr, bool) {
	if v, ok := operateComp(this.Left, this.Right, lesseqInt, lesseqFloat); ok {
		return v, true
	}
	l, lsimpl := this.Left.SimplifyNum()
	r, rsimpl := this.Right.SimplifyNum()
	return InstNumComparisonPredicate{InstNumLessEq{l, r}}, lsimpl || rsimpl
}
func (this InstNumGreater) SimplifyNumComp() (InstBoolExpr, bool) {
	if v, ok := operateComp(this.Left, this.Right, greaterInt, greaterFloat); ok {
		return v, true
	}
	l, lsimpl := this.Left.SimplifyNum()
	r, rsimpl := this.Right.SimplifyNum()
	return InstNumComparisonPredicate{InstNumGreater{l, r}}, lsimpl || rsimpl
}
func (this InstNumGreaterEq) SimplifyNumComp() (InstBoolExpr, bool) {
	if v, ok := operateComp(this.Left, this.Right, greateqInt, greateqFloat); ok {
		return v, true
	}
	l, lsimpl := this.Left.SimplifyNum()
	r, rsimpl := this.Right.SimplifyNum()
	return InstNumComparisonPredicate{InstNumGreaterEq{l, r}}, lsimpl || rsimpl
}
func (this InstNumEq) SimplifyNumComp() (InstBoolExpr, bool) {
	if v, ok := operateComp(this.Left, this.Right, eqInt, eqFloat); ok {
		return v, true
	}
	l, lsimpl := this.Left.SimplifyNum()
	r, rsimpl := this.Right.SimplifyNum()
	return InstNumComparisonPredicate{InstNumEq{l, r}}, lsimpl || rsimpl
}
func (this InstNumNotEq) SimplifyNumComp() (InstBoolExpr, bool) {
	if v, ok := operateComp(this.Left, this.Right, neqInt, neqFloat); ok {
		return v, true
	}
	l, lsimpl := this.Left.SimplifyNum()
	r, rsimpl := this.Right.SimplifyNum()
	return InstNumComparisonPredicate{InstNumNotEq{l, r}}, lsimpl || rsimpl
}

func (this InstIntLiteralExpr) SimplifyNum() (InstNumExpr, bool) {
	return this, false
}
func (this InstFloatLiteralExpr) SimplifyNum() (InstNumExpr, bool) {
	return this, false
}
func (this InstNumMulExpr) SimplifyNum() (InstNumExpr, bool) {
	if v, ok := checkNeutralOperate(this.Left, this.Right, 1, multInt, multFloat); ok {
		return v, true
	}
	l, lsimpl := this.Left.SimplifyNum()
	r, rsimpl := this.Right.SimplifyNum()
	return InstNumMulExpr{l, r}, lsimpl || rsimpl
}
func (this InstNumDivExpr) SimplifyNum() (InstNumExpr, bool) {
	vil, il := this.Left.(InstIntLiteralExpr)
	vfl, fl := this.Left.(InstFloatLiteralExpr)
	vir, ir := this.Right.(InstIntLiteralExpr)
	vfr, fr := this.Right.(InstFloatLiteralExpr)
	neutralR := (ir && vir.Num == 1) || (fr && vfr.Num == float32(1))
	if il && ir { //both are int literals, operate
		//fmt.Printf("Both ints \n")
		return InstIntLiteralExpr{divInt(vil.Num, vir.Num)}, true
	}
	if il && fr { //int op float
		//fmt.Printf("int op float \n")
		return InstFloatLiteralExpr{divFloat(float32(vil.Num), vfr.Num)}, true
	}
	if fl && ir { //float op int
		//fmt.Printf("float op int \n")
		return InstFloatLiteralExpr{divFloat(vfl.Num, float32(vir.Num))}, true
	}
	if fl && fr { //float op float
		//fmt.Printf("float op float \n")
		return InstFloatLiteralExpr{divFloat(vfl.Num, vfr.Num)}, true
	}
	if neutralR { //divisor is 1
		return this.Left.SimplifyNum()
	}
	l, lsimpl := this.Left.SimplifyNum()
	r, rsimpl := this.Right.SimplifyNum()
	return InstNumDivExpr{l, r}, lsimpl || rsimpl
}
func (this InstNumPlusExpr) SimplifyNum() (InstNumExpr, bool) {
	//fmt.Printf("Simplifying Sum expression: %s\n", this.Sprint())
	if v, ok := checkNeutralOperate(this.Left, this.Right, 0, plusInt, plusFloat); ok {
		//fmt.Printf("Neutral expression:\n")
		return v, true
	}
	l, lsimpl := this.Left.SimplifyNum()
	r, rsimpl := this.Right.SimplifyNum()
	return InstNumPlusExpr{l, r}, lsimpl || rsimpl
}
func (this InstNumMinusExpr) SimplifyNum() (InstNumExpr, bool) {
	if v, ok := checkNeutralOperate(this.Left, this.Right, 0, minusInt, minusFloat); ok {
		return v, true
	}
	l, lsimpl := this.Left.SimplifyNum()
	r, rsimpl := this.Right.SimplifyNum()
	return InstNumMinusExpr{l, r}, lsimpl || rsimpl
}
func (this InstStreamOffsetExpr) SimplifyNum() (InstNumExpr, bool) {
	return this, false //note it is not the same pattern as with Substitute
}
func (this InstConstExpr) SimplifyNum() (InstNumExpr, bool) {
	return this, false
}

//String
func (this InstStringLiteralExpr) SimplifyStr() (InstStrExpr, bool) {
	return this, false
}
func (this InstStrConcatExpr) SimplifyStr() (InstStrExpr, bool) {
	if v, ok := checkEmptyOperate(this.Left, this.Right, "", concatStr); ok {
		return v, true
	}
	l, lsimpl := this.Left.SimplifyStr()
	r, rsimpl := this.Right.SimplifyStr()
	return InstStrConcatExpr{l, r}, lsimpl || rsimpl
}
func (this InstStreamOffsetExpr) SimplifyStr() (InstStrExpr, bool) {
	return this, false //note it is not the same pattern as with Substitute
}
func (this InstConstExpr) SimplifyStr() (InstStrExpr, bool) {
	return this, false
}
func (this InstStrEqExpr) SimplifyStrComp() (InstBoolExpr, bool) {
	if v, ok := operateCompStr(this.Left, this.Right, eqStr); ok {
		return v, true
	}
	l, lsimpl := this.Left.SimplifyStr()
	r, rsimpl := this.Right.SimplifyStr()
	return InstStrComparisonPredicate{InstStrEqExpr{l, r}}, lsimpl || rsimpl
}

//Literals need to implement InstExpr to compile, implementation of Simplify (should not be needed at runtime)
//will be used as the result value of the expression, note they are InstExpr, not the corresponding subtype
func (this InstTruePredicate) Simplify() (InstExpr, bool) {
	return this, false
}
func (this InstFalsePredicate) Simplify() (InstExpr, bool) {
	return this, false
}
func (this InstIntLiteralExpr) Simplify() (InstExpr, bool) {
	return this, false
}
func (this InstFloatLiteralExpr) Simplify() (InstExpr, bool) {
	return this, false
}
func (this InstStringLiteralExpr) Simplify() (InstExpr, bool) {
	return this, false
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
	//fmt.Printf("Check Neutral op: il: %t, fl: %t, ir:%t, fr:%t\n", il, fl, ir, fr)
	if il && ir { //both are int literals, operate
		//fmt.Printf("Both ints \n")
		return InstIntLiteralExpr{fint(vil.Num, vir.Num)}, true
	}
	if il && fr { //int op float
		//fmt.Printf("int op float \n")
		return InstFloatLiteralExpr{ffloat(float32(vil.Num), vfr.Num)}, true
	}
	if fl && ir { //float op int
		//fmt.Printf("float op int \n")
		return InstFloatLiteralExpr{ffloat(vfl.Num, float32(vir.Num))}, true
	}
	if fl && fr { //float op float
		//fmt.Printf("float op float \n")
		return InstFloatLiteralExpr{ffloat(vfl.Num, vfr.Num)}, true
	}
	if neutralL { //left operand is neutral of the operation
		//fmt.Printf("Left was neutral \n")
		e, _ := right.SimplifyNum()
		return e, true
	}
	if neutralR { //right operand is neutral of the operation
		//fmt.Printf("Right was neutral \n")
		e, _ := left.SimplifyNum()
		return e, true
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
	if sl && sr {
		return InstStringLiteralExpr{fstr(vsl.S, vsr.S)}, true
	}
	if neutralL {
		return right.SimplifyStr()
	}
	if neutralR {
		return left.SimplifyStr()
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

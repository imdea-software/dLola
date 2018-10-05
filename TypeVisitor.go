package dLola

import (
	"fmt"
	//	"strings"
)

type TypeVisitor struct { //implements ExprVisitor, BooleanExprVisitor, NumExprVisitor, NumComparisonVisitor and StreamExprVisitor
	symTab           map[StreamName]StreamType //symbol table containing the declared variables and their type, StreamName is just string
	errors           []string                  //list of all the errors found
	reqType          StreamType                //requested type for the subexpressions /*IMPORTANT the types used here will be StreamType, in ast.go , see the constants below the type definition!!!
	streamsForbidden bool                      //for constants expressions that must not include stream references
}

func (v *TypeVisitor) VisitConstExpr(c ConstExpr) { //this is the usage of a constant in an expression
	constType, ok := v.symTab[c.Name]
	if !ok {
		s := fmt.Sprintf("Line %d(%d): Constant %s has not been declared", c.Pos.Line, c.Pos.Col, c.Name)
		v.errors = append(v.errors, s)
	} else {
		if v.reqType != constType {
			s := fmt.Sprintf("Line %d(%d): Cannot use constant %s of type %s in a %s expression", c.Pos.Line, c.Pos.Col, c.Name, constType.Sprint(), v.reqType.Sprint())
			v.errors = append(v.errors, s)
		}
	}
}

func (v *TypeVisitor) VisitLetExpr(l LetExpr) {
	//TODO maybe resolve the syntactic sugar substituting the variables in the expression, and then theck types normally(limitation: in a body that uses the same variable repeatedly, it will be computed more than once) e.g. input num v; let a = 3*v in a + a + a
	//TODO for the bindings type inference of the expression to know the type of the variable, and check there are no other streams with the same name
	//TODO for the body its type must match the type of the output stream

}

func (v *TypeVisitor) VisitIfThenElseExpr(ite IfThenElseExpr) {
	checkTypeIf(v, ite)
}

func (v *TypeVisitor) VisitStringExpr(s StringExpr) {
	if v.reqType != StringT {
		s := fmt.Sprintf("Line %d(%d): Cannot assign a String expression to a %s stream", s.GetPos().Line, s.GetPos().Col, v.reqType.Sprint())
		v.errors = append(v.errors, s)
	}
	s.StExpr.AcceptStr(v)
}

func (v *TypeVisitor) VisitStreamOffsetExpr(s StreamOffsetExpr) {
	s.SExpr.AcceptStream(v)
}

func (v *TypeVisitor) VisitBooleanExpr(b BooleanExpr) {
	if v.reqType != BoolT {
		s := fmt.Sprintf("Line %d(%d): Cannot assign a Boolean expression to a %s stream", b.GetPos().Line, b.GetPos().Col, v.reqType.Sprint())
		v.errors = append(v.errors, s)
	}
	b.BExpr.AcceptBool(v) //will check bool type
}

func (v *TypeVisitor) VisitNumericExpr(n NumericExpr) {
	if v.reqType != NumT {
		s := fmt.Sprintf("Line %d(%d): Cannot assign a Numeric expression to a %s stream", n.GetPos().Line, n.GetPos().Col, v.reqType.Sprint())
		v.errors = append(v.errors, s)
	}
	n.NExpr.AcceptNum(v) //will check num type
}

/*BoolExprVisitor methods*/
func (v *TypeVisitor) VisitTruePredicate(t TruePredicate) {
	if v.reqType != BoolT {
		s := fmt.Sprintf("Line %d(%d): Cannot use True in a non-boolean expression", t.Pos.Line, t.Pos.Col)
		v.errors = append(v.errors, s)
	}
}
func (v *TypeVisitor) VisitFalsePredicate(f FalsePredicate) {
	if v.reqType != BoolT {
		s := fmt.Sprintf("Line %d(%d): Cannot use False in a non-boolean expression", f.Pos.Line, f.Pos.Col)
		v.errors = append(v.errors, s)
	}
}
func (v *TypeVisitor) VisitNotPredicate(n NotPredicate) {
	v.reqType = BoolT
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

func (v *TypeVisitor) VisitStrComparisonPredicate(s StrComparisonPredicate) {
	s.Comp.AcceptStrComp(v)
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
	if v.reqType != NumT {
		s := fmt.Sprintf("Line %d(%d): Cannot use Int Literal in a non-numeric expression", i.Pos.Line, i.Pos.Col)
		v.errors = append(v.errors, s)
	}
}

func (v *TypeVisitor) VisitFloatLiteralExpr(f FloatLiteralExpr) {
	if v.reqType != NumT {
		s := fmt.Sprintf("Line %d(%d): Cannot use Float Literal in a non-numeric expression", f.Pos.Line, f.Pos.Col)
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
	streamdef := s.Default
	if v.streamsForbidden { //this stream is being used in a constant declaration
		err := fmt.Sprintf("Line %d(%d): Stream %s cannot be used in a constant declaration ", s.Pos.Line, s.Pos.Col, streamname)
		v.errors = append(v.errors, err)
	}

	if streamoffset.err { //streamoffset is a float if _, ok := n.(FloatLiteralExpr); ok {
		err := fmt.Sprintf("Line %d(%d): Stream %s cannot have a non integer offset ", s.Pos.Line, s.Pos.Col, streamname)
		v.errors = append(v.errors, err)
	}

	streamtype, ok := v.symTab[streamname]
	if !ok { //not declared
		err := fmt.Sprintf("Line %d(%d): Stream %s not declared", s.Pos.Line, s.Pos.Col, streamname)
		v.errors = append(v.errors, err)
	} else { //declared
		if streamdef.typp == Unknown { //offset = 0
			if streamtype != v.reqType {
				err := fmt.Sprintf("line %d(%d): Stream %s is of type %s but it is required to have type %s", s.Pos.Line, s.Pos.Col, streamname, streamtype.Sprint(), v.reqType.Sprint())
				v.errors = append(v.errors, err)
			}
		} else { //offset != 0
			if streamtype != v.reqType || streamtype != streamdef.typp {
				err := fmt.Sprintf("line %d(%d): Stream %s is of type %s, it is required to have type %s, but its default value has type %s",
					s.Pos.Line, s.Pos.Col, streamname, streamtype.Sprint(), v.reqType.Sprint(), streamdef.typp.Sprint())
				v.errors = append(v.errors, err)
			}
		}
	}
}

/*END StreamExprVisitor methods*/

/*StrExprVisitor methods: strings*/

func (v *TypeVisitor) VisitStringLiteralExpr(s StringLiteralExpr) {
	if v.reqType != StringT {
		s := fmt.Sprintf("Line %d(%d): Cannot use a StringLiteral in a non-string expression", s.Pos.Line, s.Pos.Col)
		v.errors = append(v.errors, s)
	}
}

func (v *TypeVisitor) VisitStrConcatExpr(s StrConcatExpr) {
	checkTypeStrOp(v, s.Left, s.Right)
}

/*END StrExprVisitor methods*/

/*StrComparisonVisitor methods: strings*/

func (v *TypeVisitor) VisitStrEqExpr(s StrEqExpr) {
	checkTypeStrOp(v, s.Left, s.Right)
}

/*END StrComparisonVisitor methods*/

/*Not exported functions*/
func checkTypeNumOp(v *TypeVisitor, left NumExpr, right NumExpr) {
	v.reqType = NumT
	left.AcceptNum(v) //will check the left expression
	v.reqType = NumT
	right.AcceptNum(v) //will check the right expression
}

func checkTypeBoolOp(v *TypeVisitor, left BoolExpr, right BoolExpr) {
	v.reqType = BoolT
	left.AcceptBool(v)  //will check the left expression
	v.reqType = BoolT   //IMPORTANT:this is needed because there are expressions that return a boolean but their operands are not boolean, and therefore they will request appropiate types, e.g. "a" SEq "b" and true
	right.AcceptBool(v) //will check the right expression
}

func checkTypeIf(v *TypeVisitor, ite IfThenElseExpr) {
	outputType := v.reqType // type of the stream, set in spec.go before calling Accept
	v.reqType = BoolT
	ite.If.Accept(v)       //will check the left expression
	v.reqType = outputType //v.getStreamType(ite.Then)
	ite.Then.Accept(v)     //will check the right expression
	v.reqType = outputType // v.getStreamType(ite.Then) //so the Accept on the Else branch will check if it returns the same type as the Then branch
	ite.Else.Accept(v)     //will check the right expression
}

func checkTypeStrOp(v *TypeVisitor, left StrExpr, right StrExpr) {
	v.reqType = StringT
	left.AcceptStr(v) //will check the left expression
	v.reqType = StringT
	right.AcceptStr(v) //will check the right expression
}

/*END Not exported functions*/

/*VisitStreamFetch*/
/*
	if ok && streamtype == v.reqType && streamtype == streamdef.typp { //declared and types match with superexpression and the default value

	} else {
		if ok { //declared but types do not match
			err := fmt.Sprintf("line %d(%d): Stream %s is of type %s but it is required to have type %s, and its default value has type %s", streamname, streamtype.Sprint(), v.reqType.Sprint(), streamdef.typp.Sprint())
			v.errors = append(v.errors, err)

		} else { //not declared
			err := fmt.Sprintf("line %d(%d): Stream %s not declared", streamname)
			v.errors = append(v.errors, err)
		}
	}
*/

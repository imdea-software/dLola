package dLola

import (
	//	"errors"
	"bufio"
	"fmt"
	"os"
	"strconv"
)

func generateInput(s StreamName, t StreamType, eval bool, c chan Resolved, tlen int, ttlMap map[StreamName]Time) {
	//go readEventFile(s,t,c, onComma, makeResolved)
	go produceEvent(s, t, eval, c, tlen, ttlMap)
}

/*Produce event instead of reading it*/
func produceEvent(s StreamName, t StreamType, eval bool, c chan Resolved, tlen int, ttlMap map[StreamName]Time) {
	var v InstExpr
	for i := 0; i < tlen; i++ {
		inststream := InstStreamFetchExpr{s, i}
		switch t {
		case BoolT:
			p := Position{0, 0, 0}
			if i%2 == 0 {
				v = InstTruePredicate{p}
			} else {
				v = InstFalsePredicate{p}
			}
		case NumT:
			v = InstIntLiteralExpr{i + 1}
		case StringT:
			v = InstStringLiteralExpr{string(i)}
		default:

		}
		//fmt.Printf("Producing event %v\n", Resolved{inststream, Resp{v, eval, i, 0, ttlMap[s]}})
		c <- Resolved{inststream, Resp{v, eval, i, 0, ttlMap[s]}}
	}
}

//needs a function to split tokens in the input file and another to parse the token and produce a Resolved
func readEventFile(s StreamName, t StreamType, c chan Resolved, tokenSeparator func(data []byte, atEOF bool) (advance int, token []byte, err error), tokenToResolved func(s StreamName, t StreamType, token string, tick int, ttlMap map[StreamName]Time) Resolved, ttlMap map[StreamName]Time) {
	f, err := os.Open(fmt.Sprintf("traces/%s_%s.in", s.Sprint(), t.Sprint()))
	if err != nil {
		panic(err)
	}
	reader := bufio.NewReader(f)
	scanner := bufio.NewScanner(reader)
	scanner.Split(tokenSeparator)
	// Scan.
	i := 0
	for scanner.Scan() { //for each token
		crude := scanner.Text()
		//fmt.Printf("%q ", crude)
		c <- tokenToResolved(s, t, crude, i, ttlMap)
		i++
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading input:", err)
	}
	f.Close()
}

/*Self made format of input traces*/
// Define a split function that separates on commas.
func onComma(data []byte, atEOF bool) (advance int, token []byte, err error) {
	for i := 0; i < len(data); i++ {
		if data[i] == ',' {
			return i + 1, data[:i], nil
		}
	}
	// There is one final token to be delivered, which may be the empty string.
	// Returning bufio.ErrFinalToken here tells Scan there are no more tokens after this
	// but does not trigger an error to be returned from Scan itself.
	return 0, data, bufio.ErrFinalToken
}

func makeResolved(s StreamName, t StreamType, token string, tick int, ttlMap map[StreamName]Time) Resolved {
	fmt.Printf("token %s\n", token)
	inststream := InstStreamFetchExpr{s, tick}
	var r Resolved
	switch t {
	case BoolT:
		b, err := convertToBoolLiteral(token)
		if !err {
			r = Resolved{inststream, Resp{b, false, 0, 0, ttlMap[s]}}
		}
	case NumT:
		i, err := strconv.Atoi(token)
		if err == nil {
			r = Resolved{inststream, Resp{InstIntLiteralExpr{i}, false, 0, 0, ttlMap[s]}}
		}
	case StringT:
		r = Resolved{inststream, Resp{InstStringLiteralExpr{token}, false, 0, 0, ttlMap[s]}}
	default:

	}
	fmt.Printf("input resolved, %v\n", r)
	return r
}

/*remember that the literals of every type implement InstExpr(substitute, simplify)*/
func convertToBoolLiteral(token string) (InstExpr, bool) {
	p := Position{0, 0, 0}
	switch token {
	case "true":
		return InstTruePredicate{p}, false
	case "false":
		return InstTruePredicate{p}, false
	default:
		return InstFalsePredicate{p}, true
	}
}

/*END Self made format of input traces*/

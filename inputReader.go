package dLola

import (
	//	"errors"
	"bufio"
	"fmt"
	"os"
	"strconv"
)

func generateInput(s StreamName, t StreamType, c chan Resolved, tlen int) {
	//go readEventFile(s, t, c)
	go produceEvent(s, t, c, tlen)
}

func readEventFile(s StreamName, t StreamType, c chan Resolved) {
	f, err := os.Open(s.Sprint() + ".in")
	if err != nil {
		panic(err)
	}
	reader := bufio.NewReader(f)
	scanner := bufio.NewScanner(reader)
	// Define a split function that separates on commas.
	onComma := func(data []byte, atEOF bool) (advance int, token []byte, err error) {
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
	scanner.Split(onComma)
	// Scan.
	i := 0
	for scanner.Scan() { //for each token
		crude := scanner.Text()
		//fmt.Printf("%q ", crude)
		c <- makeResolved(s, t, crude, i)
		i++
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading input:", err)
	}
	f.Close()
}

func makeResolved(s StreamName, t StreamType, token string, tick int) Resolved {
	fmt.Printf("token %s\n", token)
	inststream := InstStreamFetchExpr{s, tick}
	var r Resolved
	switch t {
	case BoolT:
		b, err := convertToBoolLiteral(token)
		if !err {
			r = Resolved{inststream, Resp{b, false, 0, 0}}
		}
	case NumT:
		i, err := strconv.Atoi(token)
		if err == nil {
			r = Resolved{inststream, Resp{InstIntLiteralExpr{i}, false, 0, 0}}
		}
	case StringT:
		r = Resolved{inststream, Resp{InstStringLiteralExpr{token}, false, 0, 0}}
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

func produceEvent(s StreamName, t StreamType, c chan Resolved, tlen int) {
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
			v = InstIntLiteralExpr{i}
		case StringT:
			v = InstStringLiteralExpr{string(i)}
		default:

		}
		c <- Resolved{inststream, Resp{v, false, 0, 0}}
	}
}

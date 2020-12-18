package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/danp/adventofcode/scaffold"
)

func main() {
	lines := scaffold.Lines()

	var sum int
	for _, l := range lines {
		res := eval(l)
		fmt.Println(l, "=", res)
		sum += res
	}
	fmt.Println("sum with left to right", sum)

	sum = 0
	for _, l := range lines {
		res := evalPlus(l)
		fmt.Println(l, "=", res)
		sum += res
	}
	fmt.Println("sum with left to right", sum)
}

func eval(s string) int {
	var val int
	op := '+'

	for len(s) > 0 {
		var next int
		if s[0] == '(' {
			var depth int
			for i := 1; i < len(s); i++ {
				if s[i] == ')' {
					if depth == 0 {
						next = eval(s[1:i])
						s = s[i+1:]
						break
					}
					depth--
				}
				if s[i] == '(' {
					depth++
				}
			}
		} else {
			nondig := strings.IndexFunc(s, func(r rune) bool {
				return !(byte(r) >= '0' && byte(r) <= '9')
			})
			if nondig == -1 {
				nondig = len(s)
			}
			ln, err := strconv.Atoi(s[:nondig])
			if err != nil {
				panic(err)
			}
			next = ln

			s = s[nondig:]
		}

		switch op {
		case '+':
			val += next
		case '*':
			val *= next
		}

		for len(s) > 0 && s[0] == ' ' {
			s = s[1:]
		}
		if len(s) == 0 {
			break
		}

		op = rune(s[0])
		if op != '+' && op != '*' {
			panic("unknown op")
		}
		s = s[1:]

		for len(s) > 0 && s[0] == ' ' {
			s = s[1:]
		}
	}

	return val
}

type expr interface {
	node()
}

type infixExpr struct {
	left, right expr
	op          string
}

func (i infixExpr) node() {}

type integerLit struct {
	val int
}

func (i integerLit) node() {}

const (
	_ = iota
	tokenNum
	tokenOP
	tokenCP
	tokenPlus
	tokenStar
)

const (
	_ = iota
	preLowest
	preProduct
	preSum
	preParen
)

type token struct {
	kind int
	val  string
}

func evalPlus(s string) int {
	var toks []token
	for len(s) > 0 {
		if s[0] == ' ' {
			s = s[1:]
			continue
		}

		if s[0] >= '0' && s[0] <= '9' {
			eon := strings.IndexFunc(s, func(r rune) bool { return !(r >= '0' && r <= '9') })
			if eon == -1 {
				eon = len(s)
			}
			tok := s[:eon]
			toks = append(toks, token{kind: tokenNum, val: tok})
			s = s[eon:]
			continue
		}

		tok := token{val: string(s[0])}
		switch s[0] {
		case '(':
			tok.kind = tokenOP
		case ')':
			tok.kind = tokenCP
		case '+':
			tok.kind = tokenPlus
		case '*':
			tok.kind = tokenStar
		default:
			panic(fmt.Sprintf("unknown token %q", s[0]))
		}

		toks = append(toks, tok)
		s = s[1:]

	}
	if len(toks) < 3 {
		panic("expect at least three tokens")
	}
	ct := toks[0]

	precedences := map[int]int{
		tokenPlus: preSum,
		tokenStar: preProduct,
	}

	pparsers := make(map[int]func() expr)
	iparsers := make(map[int]func(expr) expr)

	parseExpression := func(prec int) expr {
		pre := pparsers[ct.kind]
		if pre == nil {
			panic("no prefix parser for " + ct.val)
		}

		left := pre()
		for {
			nprec := preLowest
			if len(toks) > 1 {
				nprec = precedences[toks[1].kind]
			}
			if len(toks) <= 1 || prec > nprec {
				return left
			}

			inf := iparsers[toks[1].kind]
			if inf == nil {
				return left
			}

			toks = toks[1:]
			ct = toks[0]
			left = inf(left)
		}
	}

	pparsers[tokenNum] = func() expr {
		v, err := strconv.Atoi(ct.val)
		if err != nil {
			panic(err)
		}
		return integerLit{val: v}
	}
	pparsers[tokenOP] = func() expr {
		toks = toks[1:]
		ct = toks[0]

		exp := parseExpression(preLowest)

		toks = toks[1:]
		ct = toks[0]

		if ct.kind != tokenCP {
			panic("expected close paren")
		}

		return exp
	}
	parseInfix := func(left expr) expr {
		exp := infixExpr{
			left: left,
			op:   ct.val,
		}

		cpre := precedences[ct.kind]
		toks = toks[1:]
		ct = toks[0]
		exp.right = parseExpression(cpre)
		return exp
	}
	iparsers[tokenPlus] = parseInfix
	iparsers[tokenStar] = parseInfix

	exp := parseExpression(preLowest)

	var evalExp func(e expr) int
	evalExp = func(e expr) int {
		switch e := e.(type) {
		case infixExpr:
			left, right := evalExp(e.left), evalExp(e.right)
			if e.op == "+" {
				return left + right
			} else if e.op == "*" {
				return left * right
			}
		case integerLit:
			return e.val
		}
		panic("unknown")
	}

	return evalExp(exp)
}
